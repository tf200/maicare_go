package schedule

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/service/notification"
	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *scheduleService) CreateSchedule(ctx context.Context, employeeID uuid.UUID, req *CreateScheduleRequest) (*CreateScheduleResponse, error) {
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

func (s *scheduleService) GetMonthlySchedulesByLocation(ctx context.Context, locationID int64, req *GetMonthlySchedulesByLocationRequest) ([]GetMonthlySchedulesByLocationResponse, error) {
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

	calendar := make(map[string][]Shift)

	for _, schedule := range schedules {
		day := schedule.Day.Time.Format("2006-01-02")
		shift := Shift{
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
		if schedule.ShiftName != nil {
			shift.ShiftName = schedule.ShiftName
		} else {
			shift.ShiftName = util.StringPtr("Custom Shift")
		}
		calendar[day] = append(calendar[day], shift)
	}

	response := []GetMonthlySchedulesByLocationResponse{}
	for date, shifts := range calendar {
		response = append(response, GetMonthlySchedulesByLocationResponse{
			Date:   date,
			Shifts: shifts,
		})
	}

	return response, nil
}

func (s *scheduleService) GetDailySchedulesByLocation(ctx context.Context, locationID int64, req *GetDailySchedulesByLocationRequest) (*GetDailySchedulesByLocationResponse, error) {
	arg := db.GetDailySchedulesByLocationParams{
		Year:       req.Year,
		Month:      req.Month,
		Day:        req.Day,
		LocationID: locationID,
	}
	schedules, err := s.Store.GetDailySchedulesByLocation(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetDailySchedulesByLocation", "Failed to get daily schedules by location", zap.Error(err))
		return nil, fmt.Errorf("failed to get daily schedules by location: %w", err)
	}

	shifts := []Shift{}
	var targetDate string
	for _, schedule := range schedules {
		if targetDate == "" {
			targetDate = schedule.Day.Time.Format("2006-01-02")
		}
		shift := Shift{
			ShiftID:           schedule.ShiftID,
			EmployeeID:        schedule.EmployeeID,
			EmployeeFirstName: schedule.EmployeeFirstName,
			EmployeeLastName:  schedule.EmployeeLastName,
			StartTime:         schedule.StartDatetime.Time,
			EndTime:           schedule.EndDatetime.Time,
			LocationID:        schedule.LocationID,
			Color:             schedule.Color,
			LocationShiftID:   schedule.LocationShiftID,
			ShiftName:         schedule.ShiftName,
			IsCustom:          schedule.IsCustom,
		}
		if schedule.ShiftName != nil {
			// If shift name is provided, add it to the shift
			shift.ShiftName = schedule.ShiftName
		} else {
			shift.ShiftName = util.StringPtr("Custom Shift")
		}
		shifts = append(shifts, shift)
	}
	if targetDate == "" {
		targetDate = time.Date(int(req.Year), time.Month(req.Month), int(req.Day), 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}
	return &GetDailySchedulesByLocationResponse{
		Date:   targetDate,
		Shifts: shifts,
	}, nil
}

func (s *scheduleService) GetScheduleByID(ctx context.Context, scheduleID uuid.UUID) (*GetScheduleByIdResponse, error) {
	schedule, err := s.Store.GetScheduleById(ctx, scheduleID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetScheduleByID", "Failed to get schedule by ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get schedule by ID: %w", err)
	}
	return &GetScheduleByIdResponse{
		ID:                schedule.ID,
		EmployeeID:        schedule.EmployeeID,
		EmployeeFirstName: schedule.EmployeeFirstName,
		EmployeeLastName:  schedule.EmployeeLastName,
		LocationID:        schedule.LocationID,
		LocationName:      schedule.LocationName,
		StartDatetime:     schedule.StartDatetime.Time,
		EndDatetime:       schedule.EndDatetime.Time,
		CreatedAt:         schedule.CreatedAt.Time,
		UpdatedAt:         schedule.UpdatedAt.Time,
		Color:             schedule.Color,
		LocationShiftID:   schedule.LocationShiftID,
		LocationShiftName: schedule.LocationShiftName,
		IsCustom:          schedule.IsCustom,
	}, nil
}

func (s *scheduleService) UpdateSchedule(ctx context.Context, scheduleID uuid.UUID, updaterEmployeeID uuid.UUID, req *UpdateScheduleRequest) (*UpdateScheduleResponse, error) {
	// Get existing schedule
	existingSchedule, err := s.Store.GetScheduleById(ctx, scheduleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("schedule not found")
		}
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateSchedule", "Failed to fetch existing schedule", zap.Error(err))
		return nil, fmt.Errorf("failed to fetch existing schedule: %w", err)
	}

	// Determine if this is a custom schedule update
	isCustom := s.determineScheduleType(req, &existingSchedule)

	var res *UpdateScheduleResponse
	if isCustom {
		res, err = s.updateCustomSchedule(ctx, scheduleID, &existingSchedule, req)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = s.updatePresetSchedule(ctx, scheduleID, &existingSchedule, req)
		if err != nil {
			return nil, err
		}
	}

	// Send notification to the employee about the updated schedule
	s.sendNotificationForUpdatedSchedule(ctx, res.ID, updaterEmployeeID, res.EmployeeID, res.StartDatetime, res.EndDatetime, res.LocationName)

	return res, nil
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, scheduleID uuid.UUID) error {
	err := s.Store.DeleteSchedule(ctx, scheduleID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteSchedule", "Failed to delete schedule", zap.Error(err))
		return fmt.Errorf("failed to delete schedule: %w", err)
	}
	return nil
}

// ================== Private Methods ==================

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

func (s *scheduleService) createCustomSchedule(ctx context.Context, employeeID uuid.UUID, req *CreateScheduleRequest) (*CreateScheduleResponse, error) {
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

func (s *scheduleService) createPresetSchedule(ctx context.Context, employeeID uuid.UUID, req *CreateScheduleRequest) (*CreateScheduleResponse, error) {
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

func (s *scheduleService) sendNotificationForNewSchedule(ctx context.Context, scheduleID uuid.UUID, creatorID, recipientID uuid.UUID, startTime, endTime time.Time, locationName string) {
	notifData := &notification.NewScheduleNotificationData{
		ScheduleID: scheduleID,
		CreatedBy:  creatorID,
		StartTime:  startTime,
		EndTime:    endTime,
		Location:   locationName,
	}
	err := s.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
		RecipientUserIDs: []uuid.UUID{recipientID},
		Type:             notification.TypeNewScheduleNotification,
		Data:             notification.NotificationData{NewScheduleNotification: notifData},
		CreatedAt:        time.Now(),
		Message:          notifData.NewScheduleMessage(),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "sendNotificationForNewSchedule", "Failed to enqueue new schedule notification", zap.Error(err), zap.String("schedule_id", scheduleID.String()))
	}
}

func (s *scheduleService) determineScheduleType(req *UpdateScheduleRequest, existingSchedule *db.GetScheduleByIdRow) bool {
	// is_custom must be explicitly provided to use start_datetime and end_datetime
	if req.IsCustom != nil {
		isCustom := *req.IsCustom
		// If is_custom is false, ignore start_datetime and end_datetime
		if !isCustom {
			req.StartDatetime = nil
			req.EndDatetime = nil
		}
		return isCustom
	}

	// If is_custom is not provided, infer from request fields
	if req.LocationShiftID != nil || req.ShiftDate != nil {
		// Preset schedule fields provided, so it's a preset schedule
		req.StartDatetime = nil
		req.EndDatetime = nil
		return false
	}

	// No schedule type fields provided, keep existing type
	// Check if existing schedule has location_shift_id to determine type
	return existingSchedule.LocationShiftID == nil
}

func (s *scheduleService) updateCustomSchedule(ctx context.Context, scheduleID uuid.UUID, existingSchedule *db.GetScheduleByIdRow, req *UpdateScheduleRequest) (*UpdateScheduleResponse, error) {
	// Validate custom schedule update
	err := s.validateCustomScheduleUpdate(req)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "updateCustomSchedule", "Custom schedule validation failed", zap.Error(err))
		return nil, err
	}

	// Prepare update parameters with existing values as defaults
	employeeID := existingSchedule.EmployeeID
	locationID := existingSchedule.LocationID
	color := existingSchedule.Color
	startDatetime := existingSchedule.StartDatetime.Time
	endDatetime := existingSchedule.EndDatetime.Time

	// Update fields if provided
	if req.EmployeeID != nil {
		employeeID = *req.EmployeeID
	}
	if req.LocationID != nil {
		locationID = *req.LocationID
	}
	if req.Color != nil {
		color = req.Color
	}
	if req.StartDatetime != nil {
		startDatetime = *req.StartDatetime
	}
	if req.EndDatetime != nil {
		endDatetime = *req.EndDatetime
	}

	// Validate datetime order
	if startDatetime.After(endDatetime) {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "updateCustomSchedule", "start_datetime must be before end_datetime", zap.Time("start", startDatetime), zap.Time("end", endDatetime))
		return nil, fmt.Errorf("start_datetime must be before end_datetime")
	}

	// Update the schedule
	arg := db.UpdateScheduleParams{
		ID:              scheduleID,
		EmployeeID:      employeeID,
		LocationID:      locationID,
		LocationShiftID: nil, // Clear location_shift_id for custom schedules
		StartDatetime:   pgtype.Timestamp{Time: startDatetime, Valid: true},
		EndDatetime:     pgtype.Timestamp{Time: endDatetime, Valid: true},
		Color:           color,
	}

	schedule, err := s.Store.UpdateSchedule(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "updateCustomSchedule", "Failed to update custom schedule", zap.Error(err))
		return nil, fmt.Errorf("failed to update custom schedule: %w", err)
	}

	return &UpdateScheduleResponse{
		ID:              schedule.ID,
		EmployeeID:      schedule.EmployeeID,
		LocationID:      schedule.LocationID,
		LocationName:    schedule.LocationName,
		StartDatetime:   schedule.StartDatetime.Time,
		EndDatetime:     schedule.EndDatetime.Time,
		Color:           schedule.Color,
		CreatedAt:       schedule.CreatedAt.Time,
		UpdatedAt:       schedule.UpdatedAt.Time,
		LocationShiftID: nil,
		ShiftName:       nil,
	}, nil
}

func (s *scheduleService) updatePresetSchedule(ctx context.Context, scheduleID uuid.UUID, existingSchedule *db.GetScheduleByIdRow, req *UpdateScheduleRequest) (*UpdateScheduleResponse, error) {
	// Prepare update parameters with existing values as defaults
	employeeID := existingSchedule.EmployeeID
	locationID := existingSchedule.LocationID
	color := existingSchedule.Color

	// Update fields if provided
	if req.EmployeeID != nil {
		employeeID = *req.EmployeeID
	}
	if req.LocationID != nil {
		locationID = *req.LocationID
	}
	if req.Color != nil {
		color = req.Color
	}

	// Determine shift ID and date to use
	var shiftIDToUse int64
	var shiftDateToUse string

	// Use existing values if not provided
	if req.LocationShiftID != nil {
		shiftIDToUse = *req.LocationShiftID
	} else if existingSchedule.LocationShiftID != nil {
		shiftIDToUse = *existingSchedule.LocationShiftID
	} else {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "updatePresetSchedule", "location_shift_id is required for preset shift schedules")
		return nil, fmt.Errorf("location_shift_id is required for preset shift schedules")
	}

	if req.ShiftDate != nil {
		shiftDateToUse = *req.ShiftDate
	} else {
		// Extract date from existing start_datetime
		shiftDateToUse = existingSchedule.StartDatetime.Time.Format("2006-01-02")
	}

	// Get the location_shift details
	locationShift, err := s.Store.GetShiftByID(ctx, shiftIDToUse)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "updatePresetSchedule", "Failed to fetch location shift", zap.Error(err))
		return nil, fmt.Errorf("invalid location_shift_id: %w", err)
	}

	// Verify the shift belongs to the specified location
	if locationShift.LocationID != locationID {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "updatePresetSchedule", "Location shift does not belong to the specified location", zap.Int64("location_shift_id", shiftIDToUse), zap.Int64("location_id", locationID))
		return nil, fmt.Errorf("location_shift_id does not belong to the specified location")
	}

	// Parse the shift date
	shiftDate, err := time.Parse("2006-01-02", shiftDateToUse)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "updatePresetSchedule", "Invalid shift_date format", zap.String("shift_date", shiftDateToUse), zap.Error(err))
		return nil, fmt.Errorf("invalid shift_date format, expected YYYY-MM-DD: %w", err)
	}

	// Convert pgtype.Time (microseconds since midnight) to time components
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

	// Update the schedule
	arg := db.UpdateScheduleParams{
		ID:              scheduleID,
		EmployeeID:      employeeID,
		LocationID:      locationID,
		LocationShiftID: &shiftIDToUse,
		StartDatetime:   pgtype.Timestamp{Time: startDatetime, Valid: true},
		EndDatetime:     pgtype.Timestamp{Time: endDatetime, Valid: true},
		Color:           color,
	}

	schedule, err := s.Store.UpdateSchedule(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "updatePresetSchedule", "Failed to update preset schedule", zap.Error(err))
		return nil, fmt.Errorf("failed to update preset schedule: %w", err)
	}

	return &UpdateScheduleResponse{
		ID:              schedule.ID,
		EmployeeID:      schedule.EmployeeID,
		LocationID:      schedule.LocationID,
		StartDatetime:   schedule.StartDatetime.Time,
		EndDatetime:     schedule.EndDatetime.Time,
		Color:           schedule.Color,
		CreatedAt:       schedule.CreatedAt.Time,
		UpdatedAt:       schedule.UpdatedAt.Time,
		LocationShiftID: &shiftIDToUse,
		ShiftName:       &locationShift.ShiftName,
	}, nil
}

func (s *scheduleService) validateCustomScheduleUpdate(req *UpdateScheduleRequest) error {
	// For custom schedules, location_shift_id and shift_date should not be provided
	if req.LocationShiftID != nil || req.ShiftDate != nil {
		return fmt.Errorf("location_shift_id and shift_date should not be provided for custom schedules")
	}
	return nil
}

func (s *scheduleService) sendNotificationForUpdatedSchedule(ctx context.Context, scheduleID uuid.UUID, updaterEmployeeID, recipientEmployeeID uuid.UUID, startTime, endTime time.Time, locationName string) {
	notifData := &notification.NewScheduleNotificationData{
		ScheduleID: scheduleID,
		CreatedBy:  updaterEmployeeID,
		StartTime:  startTime,
		EndTime:    endTime,
		Location:   locationName,
	}
	err := s.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
		RecipientUserIDs: []uuid.UUID{recipientEmployeeID},
		Type:             notification.TypeNewScheduleNotification,
		Data:             notification.NotificationData{NewScheduleNotification: notifData},
		CreatedAt:        time.Now(),
		Message:          notifData.UpdatedScheduleMessage(),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "sendNotificationForUpdatedSchedule", "Failed to enqueue notification task", zap.Error(err))
	}
}
