package schedule

import (
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/service/notification"
	"maicare_go/util"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *scheduleService) CreateSchedule(ctx context.Context, employeeID int64, req *CreateScheduleRequest) (*CreateScheduleResponse, error) {
	res := &CreateScheduleResponse{}
	var err error
	if req.IsCustom {
		res, err = s.createCustomSchedule(ctx, employeeID, req)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = s.createPresetSchedule(ctx, employeeID, req)
		if err != nil {
			return nil, err
		}

	}
	// Send notification to the employee about the new schedule
	s.sendNotificationForNewSchedule(ctx, res.ID, employeeID, req.EmployeeID, res.StartDatetime, res.EndDatetime, res.LocationName)
	return res, nil
}



func (s *scheduleService) GetMonthlySchedulesByLocation(ctx context.Context, locationID int64, req *GetMonthlySchedulesByLocationRequest) ([]Shift, error) {
	atg := db.GetMonthlySchedulesByLocationParams{
		LocationID: locationID,
		Year:       req.Year,
		Month:      req.Month,
	}
	schedules, err := s.Store.GetMonthlySchedulesByLocation(ctx, atg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetMonthlySchedulesByLocation", "Failed to get monthly schedules by location", zap.Error(err))
		return nil, fmt.Errorf("failed to get monthly schedules by location: %w", err)
	}

	calendar:= make(map[string][]Shift)

	for _, schedule := range schedules {
		day := schedule.Day.Time.Format("2006-01-02")
		shitf := Shift{
			ShiftID:           schedule.ShiftID,
			EmployeeID:        schedule.EmployeeID,
			EmployeeFirstName: schedule.EmployeeFirstName,
			EmployeeLastName:  schedule.EmployeeLastName,
			StartTime:         schedule.StartDatetime.Time,
			EndTime:           schedule.EndDatetime.Time,
			LocationID:        schedule.LocationID,
			Color:             schedule.Color,
			ShiftName:         schedule.ShiftName,
			LocationShiftID:   schedule.LocationShiftID,
			IsCustom:          schedule.IsCustom,
		}
		


















func (s *scheduleService) validateCustomSchedule(req *CreateScheduleRequest) error {
	if req.StartDatetime == nil || req.EndDatetime == nil {
		return fmt.Errorf("start_datetime and end_datetime are required for custom schedules")
	}
	if req.StartDatetime.After(*req.EndDatetime) {
		return fmt.Errorf("start_datetime must be before end_datetime")
	}
	if req.LocationShiftID != nil || req.ShiftDate != nil {
		return fmt.Errorf("location_shift_id and shift_date should not be provided for custom schedules")
	}
	return nil
}

func (s *scheduleService) validatePresetSchedule(req *CreateScheduleRequest) error {
	if req.LocationShiftID == nil || req.ShiftDate == nil {
		return fmt.Errorf("location_shift_id and shift_date are required for preset shift schedules")
	}
	if req.StartDatetime != nil || req.EndDatetime != nil {
		return fmt.Errorf("start_datetime and end_datetime should not be provided for preset shift schedules")
	}
	return nil
}

func (s *scheduleService) createCustomSchedule(ctx context.Context, employeeID int64, req *CreateScheduleRequest) (*CreateScheduleResponse, error) {
	err := s.validateCustomSchedule(req)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "createCustomSchedule", "Custom schedule validation failed", zap.Error(err))
		return nil, err
	}
	arg := db.CreateScheduleParams{
		EmployeeID:          req.EmployeeID,
		LocationID:          req.LocationID,
		IsCustom:            req.IsCustom,
		LocationShiftID:     nil,
		Color:               req.Color,
		CreatedByEmployeeID: employeeID,
		StartDatetime:       pgtype.Timestamp{Time: *req.StartDatetime, Valid: true},
		EndDatetime:         pgtype.Timestamp{Time: *req.EndDatetime, Valid: true},
	}
	schedule, err := s.Store.CreateSchedule(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "createCustomSchedule", "Failed to create custom schedule", zap.Error(err))
		return nil, fmt.Errorf("failed to create custom schedule: %w", err)
	}
	return &CreateScheduleResponse{
		ID:              schedule.ID,
		EmployeeID:      schedule.EmployeeID,
		LocationID:      schedule.LocationID,
		LocationName:    schedule.LocationName,
		StartDatetime:   schedule.StartDatetime.Time,
		EndDatetime:     schedule.EndDatetime.Time,
		Color:           schedule.Color,
		CreatedAt:       schedule.CreatedAt.Time,
		UpdatedAt:       schedule.UpdatedAt.Time,
		LocationShiftID: schedule.LocationShiftID,
	}, nil
}

func (s *scheduleService) createPresetSchedule(ctx context.Context, employeeID int64, req *CreateScheduleRequest) (*CreateScheduleResponse, error) {
	err := s.validatePresetSchedule(req)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "createPresetSchedule", "Preset schedule validation failed", zap.Error(err))
		return nil, err
	}
	locationShift, err := s.Store.GetShiftByID(ctx, *req.LocationShiftID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "createPresetSchedule", "Failed to fetch location shift", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch location shift: %w", err)
	}

	if locationShift.LocationID != req.LocationID {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "createPresetSchedule", "Location shift does not belong to the specified location", zap.Int64("location_shift_id", *req.LocationShiftID), zap.Int64("location_id", req.LocationID))
		return nil, fmt.Errorf("location shift does not belong to the specified location")
	}

	shiftDate, err := time.Parse("2006-01-02", *req.ShiftDate)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "createPresetSchedule", "Invalid shift_date format", zap.String("shift_date", *req.ShiftDate), zap.Error(err))
		return nil, fmt.Errorf("invalid shift_date format: %w", err)
	}

	startHour, startMin, startSec, startNano := util.MicrosecondsToTimeComponents(locationShift.StartTime.Microseconds)
	endHour, endMin, endSec, endNano := util.MicrosecondsToTimeComponents(locationShift.EndTime.Microseconds)

	// Combine date with shift times to create full datetime
	startDatetime := time.Date(
		shiftDate.Year(), shiftDate.Month(), shiftDate.Day(),
		startHour, startMin, startSec, startNano,
		shiftDate.Location(),
	)

	endDatetime := time.Date(
		shiftDate.Year(), shiftDate.Month(), shiftDate.Day(),
		endHour, endMin, endSec, endNano,
		shiftDate.Location(),
	)

	// Handle shifts that cross midnight (end time is before start time)
	if locationShift.EndTime.Microseconds < locationShift.StartTime.Microseconds {
		endDatetime = endDatetime.AddDate(0, 0, 1)
	}

	arg := db.CreateScheduleParams{
		EmployeeID:          req.EmployeeID,
		LocationID:          req.LocationID,
		IsCustom:            req.IsCustom,
		LocationShiftID:     req.LocationShiftID,
		CreatedByEmployeeID: employeeID,
		Color:               req.Color,
		StartDatetime:       pgtype.Timestamp{Time: startDatetime, Valid: true},
		EndDatetime:         pgtype.Timestamp{Time: endDatetime, Valid: true},
	}
	schedule, err := s.Store.CreateSchedule(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "createPresetSchedule", "Failed to create preset schedule", zap.Error(err))
		return nil, fmt.Errorf("failed to create preset schedule: %w", err)
	}
	return &CreateScheduleResponse{
		ID:              schedule.ID,
		EmployeeID:      schedule.EmployeeID,
		LocationID:      schedule.LocationID,
		LocationName:    schedule.LocationName,
		StartDatetime:   schedule.StartDatetime.Time,
		EndDatetime:     schedule.EndDatetime.Time,
		Color:           schedule.Color,
		CreatedAt:       schedule.CreatedAt.Time,
		UpdatedAt:       schedule.UpdatedAt.Time,
		LocationShiftID: schedule.LocationShiftID,
		ShiftName:       &locationShift.ShiftName,
	}, nil
}

func (s *scheduleService) sendNotificationForNewSchedule(ctx context.Context, scheduleID uuid.UUID, creatorID, recipientID int64, startTime, endTime time.Time, locationName string) {
	notifData := &notification.NewScheduleNotificationData{
		ScheduleID: scheduleID,
		CreatedBy:  creatorID,
		StartTime:  startTime,
		EndTime:    endTime,
		Location:   locationName,
	}
	err := s.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
		RecipientUserIDs: []int64{recipientID},
		Type:             notification.TypeNewScheduleNotification,
		Data:             notification.NotificationData{NewScheduleNotification: notifData},
		CreatedAt:        time.Now(),
		Message:          notifData.NewScheduleMessage(),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "sendNotificationForNewSchedule", "Failed to enqueue new schedule notification", zap.Error(err), zap.String("schedule_id", scheduleID.String()))
	}
}
