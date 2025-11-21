package organization

import (
	"context"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/util"

	"go.uber.org/zap"
)

func (s *organizationService) CreateShift(ctx context.Context, req *CreateShiftApiRequest, locationID int64) (*CreateShiftApiResponse, error) {
	startTime, err := util.StringToPgTime(req.StartTime)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateShift", "Invalid start time format", zap.Error(err))
		return nil, err
	}

	endTime, err := util.StringToPgTime(req.EndTime)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateShift", "Invalid end time format", zap.Error(err))
		return nil, err
	}

	shift, err := s.Store.CreateShift(ctx, db.CreateShiftParams{
		LocationID: locationID,
		ShiftName:  req.ShiftName,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateShift", "Failed to create shift", zap.Error(err))
		return nil, err
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "CreateShift", "Shift created successfully", zap.Int64("shift_id", shift.ID))

	return &CreateShiftApiResponse{
		ID:         shift.ID,
		LocationID: shift.LocationID,
		ShiftName:  shift.ShiftName,
		StartTime:  util.PgTimeToString(shift.StartTime),
		EndTime:    util.PgTimeToString(shift.EndTime),
	}, nil
}

func (s *organizationService) UpdateShift(ctx context.Context, shiftID int64, req *UpdateShiftApiRequest) (*UpdateShiftApiResponse, error) {
	startTime, err := util.StringToPgTime(req.StartTime)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateShift", "Invalid start time format", zap.Error(err))
		return nil, err
	}

	endTime, err := util.StringToPgTime(req.EndTime)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateShift", "Invalid end time format", zap.Error(err))
		return nil, err
	}

	shift, err := s.Store.UpdateShift(ctx, db.UpdateShiftParams{
		ID:        shiftID,
		ShiftName: req.ShiftName,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateShift", "Failed to update shift", zap.Error(err))
		return nil, err
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "UpdateShift", "Shift updated successfully", zap.Int64("shift_id", shift.ID))

	return &UpdateShiftApiResponse{
		ID:         shift.ID,
		LocationID: shift.LocationID,
		ShiftName:  shift.ShiftName,
		StartTime:  util.PgTimeToString(shift.StartTime),
		EndTime:    util.PgTimeToString(shift.EndTime),
	}, nil
}

func (s *organizationService) DeleteShift(ctx context.Context, shiftID int64) error {
	err := s.Store.DeleteShift(ctx, shiftID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteShift", "Failed to delete shift", zap.Error(err))
		return err
	}
	return nil
}

func (s *organizationService) ListShiftsByLocationID(ctx context.Context, locationID int64) ([]ListShiftsByLocationIDResponse, error) {
	shifts, err := s.Store.GetShiftsByLocationID(ctx, locationID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListShiftsByLocationID", "Failed to list shifts", zap.Error(err))
		return nil, err
	}

	response := []ListShiftsByLocationIDResponse{}
	for _, shift := range shifts {
		response = append(response, ListShiftsByLocationIDResponse{
			ID:         shift.ID,
			LocationID: shift.LocationID,
			ShiftName:  shift.ShiftName,
			StartTime:  util.PgTimeToString(shift.StartTime),
			EndTime:    util.PgTimeToString(shift.EndTime),
		})
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListShiftsByLocationID", "Shifts listed successfully", zap.Int64("location_id", locationID))

	return response, nil
}
