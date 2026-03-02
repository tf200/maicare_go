package organization

import (
	"context"
	"errors"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/util"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

var ErrLocationShiftLimitReached = errors.New("location shift limit reached: max 4 shifts per location")

func (s *organizationService) CreateShift(ctx context.Context, req *CreateShiftApiRequest, locationID uuid.UUID) (*CreateShiftApiResponse, error) {
	existingShifts, err := s.Store.GetShiftsByLocationID(ctx, locationID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateShift", "Failed to check existing shifts", zap.Error(err))
		return nil, err
	}

	if len(existingShifts) >= 4 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateShift", "Location shift limit reached", zap.String("location_id", locationID.String()))
		return nil, ErrLocationShiftLimitReached
	}

	usedSlots := make(map[int16]struct{}, len(existingShifts))
	for _, existingShift := range existingShifts {
		usedSlots[existingShift.Slot] = struct{}{}
	}

	var selectedSlot int16
	for slot := int16(1); slot <= 4; slot++ {
		if _, exists := usedSlots[slot]; !exists {
			selectedSlot = slot
			break
		}
	}

	if selectedSlot == 0 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateShift", "No available shift slot found", zap.String("location_id", locationID.String()))
		return nil, ErrLocationShiftLimitReached
	}

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
		Slot:       selectedSlot,
		ShiftName:  req.ShiftName,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateShift", "Failed to create shift", zap.Error(err))
		return nil, err
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "CreateShift", "Shift created successfully", zap.String("shift_id", shift.ID.String()))

	return &CreateShiftApiResponse{
		ID:         shift.ID,
		LocationID: shift.LocationID,
		Slot:       shift.Slot,
		ShiftName:  shift.ShiftName,
		StartTime:  util.PgTimeToString(shift.StartTime),
		EndTime:    util.PgTimeToString(shift.EndTime),
	}, nil
}

func (s *organizationService) UpdateShift(ctx context.Context, shiftID uuid.UUID, req *UpdateShiftApiRequest) (*UpdateShiftApiResponse, error) {
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

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "UpdateShift", "Shift updated successfully", zap.String("shift_id", shift.ID.String()))

	return &UpdateShiftApiResponse{
		ID:         shift.ID,
		LocationID: shift.LocationID,
		Slot:       shift.Slot,
		ShiftName:  shift.ShiftName,
		StartTime:  util.PgTimeToString(shift.StartTime),
		EndTime:    util.PgTimeToString(shift.EndTime),
	}, nil
}

func (s *organizationService) DeleteShift(ctx context.Context, shiftID uuid.UUID) error {
	err := s.Store.DeleteShift(ctx, shiftID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteShift", "Failed to delete shift", zap.Error(err))
		return err
	}
	return nil
}

func (s *organizationService) ListShiftsByLocationID(ctx context.Context, locationID uuid.UUID) ([]ListShiftsByLocationIDResponse, error) {
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
			Slot:       shift.Slot,
			ShiftName:  shift.ShiftName,
			StartTime:  util.PgTimeToString(shift.StartTime),
			EndTime:    util.PgTimeToString(shift.EndTime),
		})
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListShiftsByLocationID", "Shifts listed successfully", zap.String("location_id", locationID.String()))

	return response, nil
}
