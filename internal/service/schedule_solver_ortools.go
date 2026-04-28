//go:build ortools

package service

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"maicare_go/internal/domain"

	"github.com/google/or-tools/ortools/sat/go/cpmodel"
	cmpb "github.com/google/or-tools/ortools/sat/proto/cpmodel"
	sppb "github.com/google/or-tools/ortools/sat/proto/satparameters"
	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

func (s *ScheduleService) generateSchedulesWithORTools(ctx context.Context, locationID uuid.UUID, timezone string, locationTZ *time.Location, employees []domain.ScheduleEmployeeContractHours, locationShifts []domain.ScheduleLocationShift, week int32, year int32) (*domain.AutoGenerateSchedulesResponse, error) {
	const (
		maxSolveSeconds = 90.0
		minRestMinutes  = int64(8 * 60)
	)

	weekStart, err := isoWeekStartDate(int(year), int(week), locationTZ)
	if err != nil {
		return nil, err
	}

	inputsEmp, inputsShifts, maxShiftMinutes, err := buildAutoGenInputs(employees, locationShifts)
	if err != nil {
		return nil, err
	}

	baseResponse := buildAutoGenResponseSkeleton(locationID, timezone, week, year, weekStart, inputsEmp, inputsShifts)

	minEmployeesPerDay := int(autoGenMinStaffPerShift) * len(inputsShifts)
	if len(inputsEmp) < minEmployeesPerDay {
		baseResponse.Warnings = []domain.SchedulePlanWarning{{
			Code:    "INFEASIBLE",
			Message: "Not enough employees to staff all shifts (1 shift/day per employee).",
		}}
		return baseResponse, nil
	}

	model := cpmodel.NewCpModelBuilder()
	eCount := len(inputsEmp)
	dCount := 7
	sCount := len(inputsShifts)
	assign := make([][][]cpmodel.BoolVar, eCount)
	for eIdx := 0; eIdx < eCount; eIdx++ {
		assign[eIdx] = make([][]cpmodel.BoolVar, dCount)
		for dIdx := 0; dIdx < dCount; dIdx++ {
			assign[eIdx][dIdx] = make([]cpmodel.BoolVar, sCount)
			for shIdx := 0; shIdx < sCount; shIdx++ {
				assign[eIdx][dIdx][shIdx] = model.NewBoolVar().WithName(fmt.Sprintf("a_%d_%d_%d", eIdx, dIdx, shIdx))
			}
		}
	}

	for dIdx := 0; dIdx < dCount; dIdx++ {
		for shIdx := 0; shIdx < sCount; shIdx++ {
			expr := cpmodel.NewLinearExpr()
			for eIdx := 0; eIdx < eCount; eIdx++ {
				expr.Add(assign[eIdx][dIdx][shIdx])
			}
			model.AddLinearConstraint(expr, autoGenMinStaffPerShift, autoGenMaxStaffPerShift)
		}
	}

	for eIdx := 0; eIdx < eCount; eIdx++ {
		for dIdx := 0; dIdx < dCount; dIdx++ {
			expr := cpmodel.NewLinearExpr()
			for shIdx := 0; shIdx < sCount; shIdx++ {
				expr.Add(assign[eIdx][dIdx][shIdx])
			}
			model.AddLinearConstraint(expr, 0, 1)
		}
	}

	for eIdx := 0; eIdx < eCount; eIdx++ {
		for dIdx := 0; dIdx < dCount-1; dIdx++ {
			for shA := 0; shA < sCount; shA++ {
				endA := inputsShifts[shA].EndMinutes
				startA := inputsShifts[shA].StartMinutes
				endAbs := int64(dIdx*1440 + endA)
				if endA < startA {
					endAbs = int64(dIdx*1440 + endA + 1440)
				}
				for shB := 0; shB < sCount; shB++ {
					startAbs := int64((dIdx+1)*1440 + inputsShifts[shB].StartMinutes)
					if startAbs-endAbs < minRestMinutes {
						model.AddLinearConstraint(cpmodel.NewLinearExpr().Add(assign[eIdx][dIdx][shA]).Add(assign[eIdx][dIdx+1][shB]), 0, 1)
					}
				}
			}
		}
	}

	maxTotalMinutes := int64(7) * maxShiftMinutes
	overtimeWeight := maxTotalMinutes*int64(eCount) + 1
	objective := cpmodel.NewLinearExpr()
	for eIdx, emp := range inputsEmp {
		total := model.NewIntVar(0, maxTotalMinutes).WithName(fmt.Sprintf("total_%d", eIdx))
		totalExpr := cpmodel.NewLinearExpr()
		for dIdx := 0; dIdx < dCount; dIdx++ {
			for shIdx := 0; shIdx < sCount; shIdx++ {
				totalExpr.AddTerm(assign[eIdx][dIdx][shIdx], inputsShifts[shIdx].DurationMinutes)
			}
		}
		model.AddEquality(total, totalExpr)

		overtime := model.NewIntVar(0, maxTotalMinutes).WithName(fmt.Sprintf("ot_%d", eIdx))
		model.AddGreaterOrEqual(cpmodel.NewLinearExpr().Add(overtime), cpmodel.NewLinearExpr().Add(total).AddConstant(-emp.TargetMinutes))
		objective.AddTerm(overtime, overtimeWeight)
		objective.AddTerm(total, 1)
	}
	model.Minimize(objective)

	modelProto, err := model.Model()
	if err != nil {
		return nil, err
	}

	params := &sppb.SatParameters{
		MaxTimeInSeconds: proto.Float64(maxSolveSeconds),
		NumSearchWorkers: proto.Int32(int32(runtime.NumCPU())),
	}
	res, err := cpmodel.SolveCpModelWithParameters(modelProto, params)
	if err != nil {
		return nil, err
	}

	statusStr := "infeasible"
	switch res.GetStatus() {
	case cmpb.CpSolverStatus_OPTIMAL:
		statusStr = "optimal"
	case cmpb.CpSolverStatus_FEASIBLE:
		statusStr = "feasible"
	}
	if statusStr == "infeasible" {
		baseResponse.Warnings = []domain.SchedulePlanWarning{{
			Code:    "INFEASIBLE",
			Message: "Solver could not find a feasible schedule within the time limit.",
		}}
		return baseResponse, nil
	}

	baseResponse.Status = statusStr
	baseResponse.Slots = make([]domain.SchedulePlanSlot, 0, dCount*sCount)

	assignedMinutesByEmp := make(map[uuid.UUID]int64, eCount)
	shiftCountsByEmp := make(map[uuid.UUID]map[uuid.UUID]int, eCount)
	for _, emp := range inputsEmp {
		shiftCountsByEmp[emp.ID] = make(map[uuid.UUID]int)
	}

	for dIdx := 0; dIdx < dCount; dIdx++ {
		dateStr := weekStart.AddDate(0, 0, dIdx).Format("2006-01-02")
		for shIdx, sh := range inputsShifts {
			employeeIDs := make([]uuid.UUID, 0, int(autoGenMaxStaffPerShift))
			for eIdx, emp := range inputsEmp {
				if cpmodel.SolutionBooleanValue(res, assign[eIdx][dIdx][shIdx]) {
					employeeIDs = append(employeeIDs, emp.ID)
					assignedMinutesByEmp[emp.ID] += sh.DurationMinutes
					shiftCountsByEmp[emp.ID][sh.ID]++
				}
			}
			baseResponse.Slots = append(baseResponse.Slots, domain.SchedulePlanSlot{
				Date:        dateStr,
				ShiftID:     sh.ID,
				EmployeeIDs: employeeIDs,
			})
		}
	}

	baseResponse.Summary = make([]domain.ScheduleEmployeeSummary, 0, eCount)
	for _, emp := range inputsEmp {
		assigned := assignedMinutesByEmp[emp.ID]
		overtime := int64(0)
		if assigned > emp.TargetMinutes {
			overtime = assigned - emp.TargetMinutes
		}
		baseResponse.Summary = append(baseResponse.Summary, domain.ScheduleEmployeeSummary{
			EmployeeID:      emp.ID,
			TargetMinutes:   emp.TargetMinutes,
			AssignedMinutes: assigned,
			OvertimeMinutes: overtime,
			ShiftCounts:     shiftCountsByEmp[emp.ID],
		})
	}

	return baseResponse, nil
}
