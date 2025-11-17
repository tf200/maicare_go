package schedule

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/logger"
	"maicare_go/util"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *scheduleService) AutoGenerateSchedules(ctx context.Context, req *AutoGenerateSchedulesRequest) (*AutoGenerateSchedulesResponse, error) {
	employees, err := s.Store.ListEmployeesWithContractHours(ctx, req.EmployeeIDs)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AutoGenerateSchedules", "Failed to fetch employee contract hours", zap.Error(err))
		return nil, err
	}
	if len(employees) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AutoGenerateSchedules", "No employees found with contract hours", zap.Int("EmployeeCount", len(req.EmployeeIDs)))
		return nil, fmt.Errorf("no employees found with contract hours")
	}

	locationShifts, err := s.Store.GetShiftsByLocationID(ctx, req.LocationID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AutoGenerateSchedules", "Failed to fetch location shifts", zap.Error(err))
		return nil, err
	}

	var emp []*grpclient.Employee
	for _, e := range employees {
		emp = append(emp, &grpclient.Employee{
			Id:          e.ID.String(),
			FirstName:   e.FirstName,
			LastName:    e.LastName,
			TargetHours: *e.ContractHours,
		})
	}
	var shifts []*grpclient.Shift
	for _, ls := range locationShifts {
		shifts = append(shifts, &grpclient.Shift{
			Id:        int32(ls.ID),
			ShiftName: ls.ShiftName,
			StartTime: util.MicrosecondsToTimeString(ls.StartTime.Microseconds),
			EndTime:   util.MicrosecondsToTimeString(ls.EndTime.Microseconds),
		})
	}
	response, err := s.GrpcClient.AutoGenerateSchedules(ctx, &grpclient.GenerateScheduleRequest{
		Employees: emp,
		Shifts:    shifts,
		Week:      req.Week,
		Year:      req.Year,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AutoGenerateSchedules", "gRPC call to AutoGenerateSchedules failed", zap.Error(err))
		return nil, err
	}

	// Convert gRPC response to service response
	scheduledShifts := make([]ScheduledShift, len(response.Shifts))
	for i, shift := range response.Shifts {
		employees := make([]AssignedEmployee, len(shift.Employees))
		for j, emp := range shift.Employees {
			empId, err := uuid.Parse(emp.Id)
			if err != nil {
				s.Logger.LogBusinessEvent(logger.LogLevelError, "AutoGenerateSchedules", "Failed to parse employee ID", zap.Error(err))
				return nil, err
			}
			employees[j] = AssignedEmployee{
				EmployeeID:   empId,
				EmployeeName: emp.Name,
			}
		}
		scheduledShifts[i] = ScheduledShift{
			Date:      shift.Date,
			DayName:   shift.DayName,
			ShiftId:   shift.ShiftId,
			ShiftName: shift.ShiftName,
			StartTime: shift.StartTime,
			EndTime:   shift.EndTime,
			Hours:     shift.Hours,
			Employees: employees,
		}
	}

	// Convert GridView
	gridViewData := GridView{
		Days:        response.GridView.Days,
		Dates:       response.GridView.Dates,
		ShiftsByDay: make(map[string][]GridDay),
	}

	for date, gridDay := range response.GridView.ShiftsByDay {
		shiftsMap := make(map[string][]GridShift)
		for shiftName, gridShift := range gridDay.Shifts {
			shiftsMap[shiftName] = []GridShift{
				{
					Employees: gridShift.Employees,
					Hours:     gridShift.Hours,
					Start:     gridShift.Start,
					End:       gridShift.End,
				},
			}
		}
		gridViewData.ShiftsByDay[date] = []GridDay{
			{
				Date:   gridDay.Date,
				Shifts: shiftsMap,
			},
		}
	}

	// Convert Summary
	summaries := make([]EmployeeSummary, len(response.Summary))
	for i, summary := range response.Summary {
		shiftsMap := make(map[string]int)
		for shiftType, count := range summary.Shifts {
			shiftsMap[shiftType] = int(count)
		}
		summaries[i] = EmployeeSummary{
			ID:          summary.Id,
			FirstName:   summary.FirstName,
			LastName:    summary.LastName,
			TargetHours: summary.Target,
			ActualHours: summary.Actual,
			Deviation:   summary.Deviation,
			Status:      summary.Status,
			Shifts:      shiftsMap,
		}
	}

	return &AutoGenerateSchedulesResponse{
		Status:   response.Status,
		Week:     response.Week,
		Year:     response.Year,
		Shifts:   scheduledShifts,
		GridView: gridViewData,
		Summary:  summaries,
	}, nil
}

func (s *scheduleService) SaveGeneratedSchedules(ctx context.Context, creatorID uuid.UUID, req *SaveGeneratedSchedulesRequest) error {
	if len(req.ScheduledShifts) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "SaveGeneratedSchedules", "No schedules to save", zap.String("Empty", "true"))
		return nil
	}

	shiftIDSet := make(map[int32]struct{})
	shiftIDs := make([]int32, 0)

	for _, sch := range req.ScheduledShifts {
		if _, exists := shiftIDSet[sch.ShiftId]; !exists {
			shiftIDSet[sch.ShiftId] = struct{}{}
			shiftIDs = append(shiftIDs, sch.ShiftId)
		}
	}

	// verify shifts exist
	exist, err := s.Store.CheckAllShiftsExist(ctx, db.CheckAllShiftsExistParams{
		Ids:           shiftIDs,
		ExpectedCount: int32(len(shiftIDs)),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SaveGeneratedSchedules", "Failed to verify shift IDs", zap.Error(err))
		return err
	}
	if !exist {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "SaveGeneratedSchedules", "One or more shift IDs do not exist", zap.Int32s("ShiftIDs", shiftIDs))
		return fmt.Errorf("one or more shift IDs do not exist")
	}

	// proceed to save schedules
	for _, sch := range req.ScheduledShifts {
		for _, emp := range sch.Employees {
			startTime, err := time.Parse(time.RFC3339Nano, sch.StartTime)
			if err != nil {
				// Try parsing without timezone if RFC3339Nano fails
				startTime, err = time.Parse("2006-01-02T15:04:05", sch.StartTime)
				if err != nil {
					s.Logger.LogBusinessEvent(logger.LogLevelError, "SaveGeneratedSchedules", "Failed to parse start time", zap.Error(err))
					return err
				}
				// Set to UTC
				startTime = startTime.UTC()
			}
			endTime, err := time.Parse(time.RFC3339Nano, sch.EndTime)
			if err != nil {
				// Try parsing without timezone if RFC3339Nano fails
				endTime, err = time.Parse("2006-01-02T15:04:05", sch.EndTime)
				if err != nil {
					s.Logger.LogBusinessEvent(logger.LogLevelError, "SaveGeneratedSchedules", "Failed to parse end time", zap.Error(err))
					return err
				}
				// Set to UTC
				endTime = endTime.UTC()
			}
			_, err = s.Store.CreateSchedule(ctx, db.CreateScheduleParams{
				EmployeeID:          emp.EmployeeID,
				LocationShiftID:     util.IntPtr(int64(sch.ShiftId)),
				LocationID:          req.LocationID,
				IsCustom:            false,
				CreatedByEmployeeID: creatorID,
				StartDatetime:       pgtype.Timestamp{Time: startTime, Valid: true},
				EndDatetime:         pgtype.Timestamp{Time: endTime, Valid: true},
			})
			if err != nil {
				s.Logger.LogBusinessEvent(logger.LogLevelError, "SaveGeneratedSchedules", "Failed to save schedule", zap.Error(err))
				return err
			}
		}
	}
	return nil
}
