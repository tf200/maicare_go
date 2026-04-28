package service

import (
	"fmt"
	"math"
	"time"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

const (
	autoGenMinStaffPerShift = int64(1)
	autoGenMaxStaffPerShift = int64(2)
)

func buildAutoGenInputs(employees []domain.ScheduleEmployeeContractHours, locationShifts []domain.ScheduleLocationShift) ([]autoGenEmployee, []autoGenShift, int64, error) {
	inputsEmp := make([]autoGenEmployee, 0, len(employees))
	for _, e := range employees {
		if e.ContractHours == nil || *e.ContractHours <= 0 {
			continue
		}
		targetMinutes := int64(math.Round(*e.ContractHours * 60.0))
		inputsEmp = append(inputsEmp, autoGenEmployee{
			ID:            e.ID,
			FirstName:     e.FirstName,
			LastName:      e.LastName,
			TargetMinutes: targetMinutes,
		})
	}
	if len(inputsEmp) == 0 {
		return nil, nil, 0, fmt.Errorf("no employees with positive contract hours")
	}

	inputsShifts := make([]autoGenShift, 0, len(locationShifts))
	var maxShiftMinutes int64
	for _, ls := range locationShifts {
		startMin := int(ls.StartMicroseconds / (60 * 1_000_000))
		endMin := int(ls.EndMicroseconds / (60 * 1_000_000))
		dur := int64(endMin - startMin)
		if endMin < startMin {
			dur = int64(endMin + 1440 - startMin)
		}
		if dur <= 0 {
			continue
		}
		if dur > maxShiftMinutes {
			maxShiftMinutes = dur
		}
		inputsShifts = append(inputsShifts, autoGenShift{
			ID:              ls.ID,
			Name:            ls.ShiftName,
			StartMinutes:    startMin,
			EndMinutes:      endMin,
			DurationMinutes: dur,
		})
	}
	if len(inputsShifts) == 0 {
		return nil, nil, 0, fmt.Errorf("no valid shifts found for location")
	}

	return inputsEmp, inputsShifts, maxShiftMinutes, nil
}

func buildAutoGenResponseSkeleton(locationID uuid.UUID, timezone string, week int32, year int32, weekStart time.Time, employees []autoGenEmployee, shifts []autoGenShift) *domain.AutoGenerateSchedulesResponse {
	planEmployees := make([]domain.SchedulePlanEmployee, 0, len(employees))
	for _, emp := range employees {
		planEmployees = append(planEmployees, domain.SchedulePlanEmployee{
			ID:            emp.ID,
			FirstName:     emp.FirstName,
			LastName:      emp.LastName,
			TargetMinutes: emp.TargetMinutes,
		})
	}

	shiftTemplates := make([]domain.ScheduleShiftTemplate, 0, len(shifts))
	for _, sh := range shifts {
		shiftTemplates = append(shiftTemplates, domain.ScheduleShiftTemplate{
			ShiftID:         sh.ID,
			Name:            sh.Name,
			StartMinute:     int32(sh.StartMinutes),
			EndMinute:       int32(sh.EndMinutes),
			DurationMinutes: sh.DurationMinutes,
			Overnight:       sh.EndMinutes < sh.StartMinutes,
		})
	}

	slots := make([]domain.SchedulePlanSlot, 0, 7*len(shifts))
	for day := 0; day < 7; day++ {
		dateStr := weekStart.AddDate(0, 0, day).Format("2006-01-02")
		for _, sh := range shifts {
			slots = append(slots, domain.SchedulePlanSlot{
				Date:        dateStr,
				ShiftID:     sh.ID,
				EmployeeIDs: []uuid.UUID{},
			})
		}
	}

	return &domain.AutoGenerateSchedulesResponse{
		Status:         "infeasible",
		PlanID:         uuid.New(),
		LocationID:     locationID,
		Timezone:       timezone,
		Week:           week,
		Year:           year,
		WeekStartDate:  weekStart.Format("2006-01-02"),
		Constraints:    domain.SchedulePlanConstraints{MaxStaffPerShift: int32(autoGenMaxStaffPerShift), AllowEmptyShift: true},
		Employees:      planEmployees,
		ShiftTemplates: shiftTemplates,
		Slots:          slots,
		Summary:        []domain.ScheduleEmployeeSummary{},
	}
}
