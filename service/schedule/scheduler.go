package schedule

import (
	"context"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/logger"

	"go.uber.org/zap"
)

func (s *scheduleService) AutoGenerateSchedules(ctx context.Context, req *AutoGenerateSchedulesRequest) (*AutoGenerateSchedulesResponse, error) {
	employees, err := s.Store.ListEmployeesWithContractHours(ctx, req.EmployeeIDs)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "AutoGenerateSchedules", "Failed to fetch employee contract hours", zap.Error(err))
		return nil, err
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
			StartTime: ls.StartTime.,
			EndTime:   ls.EndTime,
			Role:      ls.Role,
		})
	}

}
