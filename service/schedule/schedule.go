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

func (s *scheduleService) CreateSchedule(
	ctx context.Context,
	creatorID uuid.UUID,
	req *CreateScheduleRequest,
) ([]CreateScheduleResponse, error) {
	if err := s.validateCreateScheduleRequest(req); err != nil {
		return nil, err
	}

	recurrence, err := s.resolveRecurrence(req)
	if err != nil {
		return nil, err
	}

	results := make([]CreateScheduleResponse, 0)

	if req.IsCustom {
		dates := s.buildCustomScheduleDates(*req.StartDatetime, recurrence)
		duration := req.EndDatetime.Sub(*req.StartDatetime)

		for _, date := range dates {
			start := time.Date(
				date.Year(),
				date.Month(),
				date.Day(),
				req.StartDatetime.Hour(),
				req.StartDatetime.Minute(),
				req.StartDatetime.Second(),
				req.StartDatetime.Nanosecond(),
				req.StartDatetime.Location(),
			)
			end := start.Add(duration)
			for _, assigneeID := range req.EmployeeIDs {
				res, createErr := s.createCustomSchedule(
					ctx,
					creatorID,
					assigneeID,
					req.LocationID,
					start,
					end,
				)
				if createErr != nil {
					return nil, createErr
				}
				results = append(results, *res)
				s.sendNotificationForNewSchedule(
					ctx,
					res.ID,
					creatorID,
					assigneeID,
					res.StartDatetime,
					res.EndDatetime,
					res.LocationName,
				)
			}
		}

		return results, nil
	}

	locationShift, locationTZ, err := s.getPresetScheduleContext(ctx, req)
	if err != nil {
		return nil, err
	}

	baseDate, err := time.ParseInLocation("2006-01-02", *req.ShiftDate, locationTZ)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"CreateSchedule",
			"Invalid shift_date format",
			zap.String("shift_date", *req.ShiftDate),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid shift_date format: %w", err)
	}

	dates := s.buildCustomScheduleDates(baseDate, recurrence)
	for _, date := range dates {
		for _, assigneeID := range req.EmployeeIDs {
			res, createErr := s.createPresetScheduleForDate(
				ctx,
				creatorID,
				assigneeID,
				req.LocationID,
				req.LocationShiftID,
				locationShift,
				date,
				locationTZ,
			)
			if createErr != nil {
				return nil, createErr
			}
			results = append(results, *res)
			s.sendNotificationForNewSchedule(
				ctx,
				res.ID,
				creatorID,
				assigneeID,
				res.StartDatetime,
				res.EndDatetime,
				res.LocationName,
			)
		}
	}

	return results, nil
}

func (s *scheduleService) GetSchedulesByLocationInRange(
	ctx context.Context,
	locationID uuid.UUID,
	req *GetSchedulesByLocationInRangeRequest,
) ([]GetSchedulesByLocationInRangeResponse, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetSchedulesByLocationInRange",
			"Invalid start_date format",
			zap.String("start_date", req.StartDate),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid start_date format, expected YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetSchedulesByLocationInRange",
			"Invalid end_date format",
			zap.String("end_date", req.EndDate),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid end_date format, expected YYYY-MM-DD")
	}

	if endDate.Before(startDate) {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetSchedulesByLocationInRange",
			"end_date is before start_date",
			zap.String("start_date", req.StartDate),
			zap.String("end_date", req.EndDate),
		)
		return nil, fmt.Errorf("end_date must be on or after start_date")
	}

	rows, err := s.Store.GetSchedulesByLocationInRange(ctx, db.GetSchedulesByLocationInRangeParams{
		LocationID: locationID,
		StartDate:  pgtype.Date{Time: startDate, Valid: true},
		EndDate:    pgtype.Date{Time: endDate, Valid: true},
	})
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetSchedulesByLocationInRange",
			"Failed to list schedules by range",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list schedules by range: %w", err)
	}

	dayMap := make(map[string][]Shift, len(rows))
	for _, row := range rows {
		day := row.Day.Time.Format("2006-01-02")
		dayMap[day] = append(dayMap[day], Shift{
			ScheduleID:        row.ShiftID,
			EmployeeID:        row.EmployeeID,
			EmployeeFirstName: row.EmployeeFirstName,
			EmployeeLastName:  row.EmployeeLastName,
			StartTime:         row.StartDatetime.Time,
			EndTime:           row.EndDatetime.Time,
			LocationID:        row.LocationID,
			ShiftName:         util.StringPtr(row.ShiftName),
			LocationShiftID:   row.LocationShiftID,
			IsCustom:          row.IsCustom,
		})
	}

	response := make([]GetSchedulesByLocationInRangeResponse, 0)
	for day := startDate; !day.After(endDate); day = day.AddDate(0, 0, 1) {
		key := day.Format("2006-01-02")
		response = append(response, GetSchedulesByLocationInRangeResponse{
			Date:   key,
			Shifts: append([]Shift{}, dayMap[key]...),
		})
	}

	return response, nil
}

func (s *scheduleService) GetScheduleByID(
	ctx context.Context,
	scheduleID uuid.UUID,
) (*GetScheduleByIdResponse, error) {
	schedule, err := s.Store.GetScheduleById(ctx, scheduleID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"GetScheduleByID",
			"Failed to get schedule by ID",
			zap.Error(err),
		)
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
		LocationShiftID:   schedule.LocationShiftID,
		LocationShiftName: util.StringPtr(schedule.LocationShiftName),
		IsCustom:          schedule.IsCustom,
	}, nil
}

func (s *scheduleService) UpdateSchedule(
	ctx context.Context,
	scheduleID uuid.UUID,
	updaterEmployeeID uuid.UUID,
	req *UpdateScheduleRequest,
) (*UpdateScheduleResponse, error) {
	// Get existing schedule
	existingSchedule, err := s.Store.GetScheduleById(ctx, scheduleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("schedule not found")
		}
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"UpdateSchedule",
			"Failed to fetch existing schedule",
			zap.Error(err),
		)
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
	s.sendNotificationForUpdatedSchedule(
		ctx,
		res.ID,
		updaterEmployeeID,
		res.EmployeeID,
		res.StartDatetime,
		res.EndDatetime,
		res.LocationName,
	)

	return res, nil
}

func (s *scheduleService) DeleteSchedule(ctx context.Context, scheduleID uuid.UUID) error {
	err := s.Store.DeleteSchedule(ctx, scheduleID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"DeleteSchedule",
			"Failed to delete schedule",
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete schedule: %w", err)
	}
	return nil
}

// ================== Private Methods ==================

func (s *scheduleService) validateCreateScheduleRequest(req *CreateScheduleRequest) error {
	if len(req.EmployeeIDs) == 0 {
		return fmt.Errorf("employee_ids is required")
	}

	seen := make(map[uuid.UUID]struct{}, len(req.EmployeeIDs))
	for _, employeeID := range req.EmployeeIDs {
		if employeeID == uuid.Nil {
			return fmt.Errorf("employee_ids contains invalid uuid")
		}
		if _, exists := seen[employeeID]; exists {
			return fmt.Errorf("employee_ids must not contain duplicates")
		}
		seen[employeeID] = struct{}{}
	}

	if req.LocationID == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}

	if req.IsCustom {
		return s.validateCustomSchedule(req)
	}

	return s.validatePresetSchedule(req)
}

func (s *scheduleService) resolveRecurrence(req *CreateScheduleRequest) (string, error) {
	if req.Recurrence == nil || *req.Recurrence == "" {
		return CreateScheduleRecurrenceNone, nil
	}

	switch *req.Recurrence {
	case CreateScheduleRecurrenceNone,
		CreateScheduleRecurrenceEndOfWeek,
		CreateScheduleRecurrenceEndOfMonth:
		return *req.Recurrence, nil
	default:
		return "", fmt.Errorf("recurrence must be one of: none, end_of_week, end_of_month")
	}
}

func (s *scheduleService) buildCustomScheduleDates(
	baseDate time.Time,
	recurrence string,
) []time.Time {
	dayStart := time.Date(
		baseDate.Year(),
		baseDate.Month(),
		baseDate.Day(),
		0,
		0,
		0,
		0,
		baseDate.Location(),
	)
	endDate := dayStart

	switch recurrence {
	case CreateScheduleRecurrenceEndOfWeek:
		daysUntilSunday := (7 - int(dayStart.Weekday())) % 7
		endDate = dayStart.AddDate(0, 0, daysUntilSunday)
	case CreateScheduleRecurrenceEndOfMonth:
		endDate = time.Date(dayStart.Year(), dayStart.Month()+1, 0, 0, 0, 0, 0, dayStart.Location())
	}

	dates := make([]time.Time, 0)
	for date := dayStart; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		dates = append(dates, date)
	}

	return dates
}

func (s *scheduleService) validateCustomSchedule(req *CreateScheduleRequest) error {
	if req.StartDatetime == nil || req.EndDatetime == nil {
		return fmt.Errorf("start_datetime and end_datetime are required for custom schedules")
	}
	if req.StartDatetime.After(*req.EndDatetime) {
		return fmt.Errorf("start_datetime must be before end_datetime")
	}
	if req.LocationShiftID != nil || req.ShiftDate != nil {
		return fmt.Errorf(
			"location_shift_id and shift_date should not be provided for custom schedules",
		)
	}
	return nil
}

func (s *scheduleService) validatePresetSchedule(req *CreateScheduleRequest) error {
	if req.LocationShiftID == nil || req.ShiftDate == nil {
		return fmt.Errorf(
			"location_shift_id and shift_date are required for preset shift schedules",
		)
	}
	if req.StartDatetime != nil || req.EndDatetime != nil {
		return fmt.Errorf(
			"start_datetime and end_datetime should not be provided for preset shift schedules",
		)
	}
	return nil
}

func (s *scheduleService) createCustomSchedule(
	ctx context.Context,
	creatorID, assigneeID, locationID uuid.UUID,
	startDatetime, endDatetime time.Time,
) (*CreateScheduleResponse, error) {
	arg := db.CreateScheduleParams{
		EmployeeID:             assigneeID,
		LocationID:             locationID,
		IsCustom:               true,
		LocationShiftID:        nil,
		ShiftNameSnapshot:      nil,
		ShiftStartTimeSnapshot: pgtype.Time{Valid: false},
		ShiftEndTimeSnapshot:   pgtype.Time{Valid: false},
		CreatedByEmployeeID:    creatorID,
		StartDatetime:          pgtype.Timestamptz{Time: startDatetime, Valid: true},
		EndDatetime:            pgtype.Timestamptz{Time: endDatetime, Valid: true},
	}

	schedule, err := s.Store.CreateSchedule(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"createCustomSchedule",
			"Failed to create custom schedule",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create custom schedule: %w", err)
	}
	return &CreateScheduleResponse{
		ID:              schedule.ID,
		EmployeeID:      schedule.EmployeeID,
		LocationID:      schedule.LocationID,
		LocationName:    schedule.LocationName,
		StartDatetime:   schedule.StartDatetime.Time,
		EndDatetime:     schedule.EndDatetime.Time,
		CreatedAt:       schedule.CreatedAt.Time,
		UpdatedAt:       schedule.UpdatedAt.Time,
		LocationShiftID: schedule.LocationShiftID,
	}, nil
}

func (s *scheduleService) getPresetScheduleContext(
	ctx context.Context,
	req *CreateScheduleRequest,
) (db.LocationShift, *time.Location, error) {
	err := s.validatePresetSchedule(req)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"getPresetScheduleContext",
			"Preset schedule validation failed",
			zap.Error(err),
		)
		return db.LocationShift{}, nil, err
	}

	locationShift, err := s.Store.GetShiftByID(ctx, *req.LocationShiftID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"getPresetScheduleContext",
			"Failed to fetch location shift",
			zap.Error(err),
		)
		return db.LocationShift{}, nil, fmt.Errorf("failed to fetch location shift: %w", err)
	}

	if locationShift.LocationID != req.LocationID {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"getPresetScheduleContext",
			"Location shift does not belong to the specified location",
			zap.String("location_shift_id", req.LocationShiftID.String()),
			zap.String("location_id", req.LocationID.String()),
		)
		return db.LocationShift{}, nil, fmt.Errorf(
			"location shift does not belong to the specified location",
		)
	}

	location, err := s.Store.GetLocation(ctx, req.LocationID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"getPresetScheduleContext",
			"Failed to fetch location",
			zap.Error(err),
			zap.String("location_id", req.LocationID.String()),
		)
		return db.LocationShift{}, nil, fmt.Errorf("failed to fetch location: %w", err)
	}

	locationTZ, err := time.LoadLocation(location.Timezone)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"getPresetScheduleContext",
			"Invalid location timezone",
			zap.Error(err),
			zap.String("location_timezone", location.Timezone),
		)
		return db.LocationShift{}, nil, fmt.Errorf("invalid location timezone: %w", err)
	}

	return locationShift, locationTZ, nil
}

func (s *scheduleService) createPresetScheduleForDate(
	ctx context.Context,
	creatorID, assigneeID, locationID uuid.UUID,
	locationShiftID *uuid.UUID,
	locationShift db.LocationShift,
	shiftDate time.Time,
	locationTZ *time.Location,
) (*CreateScheduleResponse, error) {
	shiftDate = time.Date(
		shiftDate.Year(),
		shiftDate.Month(),
		shiftDate.Day(),
		0,
		0,
		0,
		0,
		locationTZ,
	)

	startHour, startMin, startSec, startNano := util.MicrosecondsToTimeComponents(
		locationShift.StartTime.Microseconds,
	)
	endHour, endMin, endSec, endNano := util.MicrosecondsToTimeComponents(
		locationShift.EndTime.Microseconds,
	)

	// Combine date with shift times to create full datetime
	startDatetime := time.Date(
		shiftDate.Year(), shiftDate.Month(), shiftDate.Day(),
		startHour, startMin, startSec, startNano,
		locationTZ,
	)

	endDatetime := time.Date(
		shiftDate.Year(), shiftDate.Month(), shiftDate.Day(),
		endHour, endMin, endSec, endNano,
		locationTZ,
	)

	// Handle shifts that cross midnight (end time is before start time)
	if locationShift.EndTime.Microseconds < locationShift.StartTime.Microseconds {
		endDatetime = endDatetime.AddDate(0, 0, 1)
	}

	arg := db.CreateScheduleParams{
		EmployeeID:             assigneeID,
		LocationID:             locationID,
		IsCustom:               false,
		LocationShiftID:        locationShiftID,
		ShiftNameSnapshot:      &locationShift.ShiftName,
		ShiftStartTimeSnapshot: locationShift.StartTime,
		ShiftEndTimeSnapshot:   locationShift.EndTime,
		CreatedByEmployeeID:    creatorID,
		StartDatetime:          pgtype.Timestamptz{Time: startDatetime, Valid: true},
		EndDatetime:            pgtype.Timestamptz{Time: endDatetime, Valid: true},
	}
	schedule, err := s.Store.CreateSchedule(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"createPresetSchedule",
			"Failed to create preset schedule",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to create preset schedule: %w", err)
	}
	return &CreateScheduleResponse{
		ID:              schedule.ID,
		EmployeeID:      schedule.EmployeeID,
		LocationID:      schedule.LocationID,
		LocationName:    schedule.LocationName,
		StartDatetime:   schedule.StartDatetime.Time,
		EndDatetime:     schedule.EndDatetime.Time,
		CreatedAt:       schedule.CreatedAt.Time,
		UpdatedAt:       schedule.UpdatedAt.Time,
		LocationShiftID: schedule.LocationShiftID,
		ShiftName:       &locationShift.ShiftName,
	}, nil
}

func (s *scheduleService) sendNotificationForNewSchedule(
	ctx context.Context,
	scheduleID uuid.UUID,
	creatorID, recipientID uuid.UUID,
	startTime, endTime time.Time,
	locationName string,
) {
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
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"sendNotificationForNewSchedule",
			"Failed to enqueue new schedule notification",
			zap.Error(err),
			zap.String("schedule_id", scheduleID.String()),
		)
	}
}

func (s *scheduleService) determineScheduleType(
	req *UpdateScheduleRequest,
	existingSchedule *db.GetScheduleByIdRow,
) bool {
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

func (s *scheduleService) updateCustomSchedule(
	ctx context.Context,
	scheduleID uuid.UUID,
	existingSchedule *db.GetScheduleByIdRow,
	req *UpdateScheduleRequest,
) (*UpdateScheduleResponse, error) {
	// Validate custom schedule update
	err := s.validateCustomScheduleUpdate(req)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"updateCustomSchedule",
			"Custom schedule validation failed",
			zap.Error(err),
		)
		return nil, err
	}

	// Prepare update parameters with existing values as defaults
	employeeID := existingSchedule.EmployeeID
	locationID := existingSchedule.LocationID
	startDatetime := existingSchedule.StartDatetime.Time
	endDatetime := existingSchedule.EndDatetime.Time

	// Update fields if provided
	if req.EmployeeID != nil {
		employeeID = *req.EmployeeID
	}
	if req.LocationID != nil {
		locationID = *req.LocationID
	}
	if req.StartDatetime != nil {
		startDatetime = *req.StartDatetime
	}
	if req.EndDatetime != nil {
		endDatetime = *req.EndDatetime
	}

	// Validate datetime order
	if startDatetime.After(endDatetime) {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"updateCustomSchedule",
			"start_datetime must be before end_datetime",
			zap.Time("start", startDatetime),
			zap.Time("end", endDatetime),
		)
		return nil, fmt.Errorf("start_datetime must be before end_datetime")
	}

	// Update the schedule
	arg := db.UpdateScheduleParams{
		ID:                     scheduleID,
		EmployeeID:             employeeID,
		LocationID:             locationID,
		LocationShiftID:        nil, // Clear location_shift_id for custom schedules
		IsCustom:               true,
		StartDatetime:          pgtype.Timestamptz{Time: startDatetime, Valid: true},
		EndDatetime:            pgtype.Timestamptz{Time: endDatetime, Valid: true},
		ShiftNameSnapshot:      nil,
		ShiftStartTimeSnapshot: pgtype.Time{Valid: false},
		ShiftEndTimeSnapshot:   pgtype.Time{Valid: false},
	}

	schedule, err := s.Store.UpdateSchedule(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"updateCustomSchedule",
			"Failed to update custom schedule",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update custom schedule: %w", err)
	}

	return &UpdateScheduleResponse{
		ID:              schedule.ID,
		EmployeeID:      schedule.EmployeeID,
		LocationID:      schedule.LocationID,
		LocationName:    schedule.LocationName,
		StartDatetime:   schedule.StartDatetime.Time,
		EndDatetime:     schedule.EndDatetime.Time,
		CreatedAt:       schedule.CreatedAt.Time,
		UpdatedAt:       schedule.UpdatedAt.Time,
		LocationShiftID: nil,
		ShiftName:       nil,
	}, nil
}

func (s *scheduleService) updatePresetSchedule(
	ctx context.Context,
	scheduleID uuid.UUID,
	existingSchedule *db.GetScheduleByIdRow,
	req *UpdateScheduleRequest,
) (*UpdateScheduleResponse, error) {
	// Prepare update parameters with existing values as defaults
	employeeID := existingSchedule.EmployeeID
	locationID := existingSchedule.LocationID

	// Update fields if provided
	if req.EmployeeID != nil {
		employeeID = *req.EmployeeID
	}
	if req.LocationID != nil {
		locationID = *req.LocationID
	}

	// Determine shift ID and date to use
	var shiftIDToUse uuid.UUID
	var shiftDateToUse string

	// Use existing values if not provided
	if req.LocationShiftID != nil {
		shiftIDToUse = *req.LocationShiftID
	} else if existingSchedule.LocationShiftID != nil {
		shiftIDToUse = *existingSchedule.LocationShiftID
	} else {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"updatePresetSchedule",
			"location_shift_id is required for preset shift schedules",
		)
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
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"updatePresetSchedule",
			"Failed to fetch location shift",
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid location_shift_id: %w", err)
	}

	// Verify the shift belongs to the specified location
	if locationShift.LocationID != locationID {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"updatePresetSchedule",
			"Location shift does not belong to the specified location",
			zap.String("location_shift_id", shiftIDToUse.String()),
			zap.String("location_id", locationID.String()),
		)
		return nil, fmt.Errorf("location_shift_id does not belong to the specified location")
	}

	// Parse the shift date
	shiftDate, err := time.Parse("2006-01-02", shiftDateToUse)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"updatePresetSchedule",
			"Invalid shift_date format",
			zap.String("shift_date", shiftDateToUse),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid shift_date format, expected YYYY-MM-DD: %w", err)
	}

	location, err := s.Store.GetLocation(ctx, locationID)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"updatePresetSchedule",
			"Failed to fetch location",
			zap.Error(err),
			zap.String("location_id", locationID.String()),
		)
		return nil, fmt.Errorf("failed to fetch location: %w", err)
	}

	locationTZ, err := time.LoadLocation(location.Timezone)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"updatePresetSchedule",
			"Invalid location timezone",
			zap.Error(err),
			zap.String("location_timezone", location.Timezone),
		)
		return nil, fmt.Errorf("invalid location timezone: %w", err)
	}

	shiftDate = time.Date(
		shiftDate.Year(),
		shiftDate.Month(),
		shiftDate.Day(),
		0,
		0,
		0,
		0,
		locationTZ,
	)

	// Convert pgtype.Time (microseconds since midnight) to time components
	startHour, startMin, startSec, startNano := util.MicrosecondsToTimeComponents(
		locationShift.StartTime.Microseconds,
	)
	endHour, endMin, endSec, endNano := util.MicrosecondsToTimeComponents(
		locationShift.EndTime.Microseconds,
	)

	// Combine date with shift times to create full datetime
	startDatetime := time.Date(
		shiftDate.Year(), shiftDate.Month(), shiftDate.Day(),
		startHour, startMin, startSec, startNano,
		locationTZ,
	)

	endDatetime := time.Date(
		shiftDate.Year(), shiftDate.Month(), shiftDate.Day(),
		endHour, endMin, endSec, endNano,
		locationTZ,
	)

	// Handle shifts that cross midnight (end time is before start time)
	if locationShift.EndTime.Microseconds < locationShift.StartTime.Microseconds {
		endDatetime = endDatetime.AddDate(0, 0, 1)
	}

	// Update the schedule
	arg := db.UpdateScheduleParams{
		ID:                     scheduleID,
		EmployeeID:             employeeID,
		LocationID:             locationID,
		LocationShiftID:        &shiftIDToUse,
		IsCustom:               false,
		StartDatetime:          pgtype.Timestamptz{Time: startDatetime, Valid: true},
		EndDatetime:            pgtype.Timestamptz{Time: endDatetime, Valid: true},
		ShiftNameSnapshot:      &locationShift.ShiftName,
		ShiftStartTimeSnapshot: locationShift.StartTime,
		ShiftEndTimeSnapshot:   locationShift.EndTime,
	}

	schedule, err := s.Store.UpdateSchedule(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"updatePresetSchedule",
			"Failed to update preset schedule",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to update preset schedule: %w", err)
	}

	return &UpdateScheduleResponse{
		ID:              schedule.ID,
		EmployeeID:      schedule.EmployeeID,
		LocationID:      schedule.LocationID,
		StartDatetime:   schedule.StartDatetime.Time,
		EndDatetime:     schedule.EndDatetime.Time,
		CreatedAt:       schedule.CreatedAt.Time,
		UpdatedAt:       schedule.UpdatedAt.Time,
		LocationShiftID: &shiftIDToUse,
		ShiftName:       &locationShift.ShiftName,
	}, nil
}

func (s *scheduleService) validateCustomScheduleUpdate(req *UpdateScheduleRequest) error {
	// For custom schedules, location_shift_id and shift_date should not be provided
	if req.LocationShiftID != nil || req.ShiftDate != nil {
		return fmt.Errorf(
			"location_shift_id and shift_date should not be provided for custom schedules",
		)
	}
	return nil
}

func (s *scheduleService) sendNotificationForUpdatedSchedule(
	ctx context.Context,
	scheduleID uuid.UUID,
	updaterEmployeeID, recipientEmployeeID uuid.UUID,
	startTime, endTime time.Time,
	locationName string,
) {
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
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"sendNotificationForUpdatedSchedule",
			"Failed to enqueue notification task",
			zap.Error(err),
		)
	}
}
