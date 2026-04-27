package service

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"time"

	"maicare_go/async/aclient"
	"maicare_go/internal/domain"
	"maicare_go/service/notification"

	"github.com/google/or-tools/ortools/sat/go/cpmodel"
	cmpb "github.com/google/or-tools/ortools/sat/proto/cpmodel"
	sppb "github.com/google/or-tools/ortools/sat/proto/satparameters"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type ScheduleService struct {
	repository  domain.ScheduleRepository
	asynqClient aclient.AsynqClientInterface
	logger      domain.Logger
}

func NewScheduleService(repository domain.ScheduleRepository, asynqClient aclient.AsynqClientInterface, logger domain.Logger) domain.ScheduleService {
	return &ScheduleService{
		repository:  repository,
		asynqClient: asynqClient,
		logger:      logger,
	}
}

func (s *ScheduleService) CreateSchedule(ctx context.Context, creatorID uuid.UUID, req *domain.CreateScheduleRequest) ([]domain.CreateScheduleResponse, error) {
	if err := s.validateCreateScheduleRequest(req); err != nil {
		return nil, err
	}

	recurrence, err := s.resolveRecurrence(req)
	if err != nil {
		return nil, err
	}

	results := make([]domain.CreateScheduleResponse, 0)
	if req.IsCustom {
		dates := s.buildCustomScheduleDates(*req.StartDatetime, recurrence)
		duration := req.EndDatetime.Sub(*req.StartDatetime)
		for _, date := range dates {
			start := time.Date(
				date.Year(), date.Month(), date.Day(),
				req.StartDatetime.Hour(), req.StartDatetime.Minute(), req.StartDatetime.Second(), req.StartDatetime.Nanosecond(),
				req.StartDatetime.Location(),
			)
			end := start.Add(duration)
			for _, assigneeID := range req.EmployeeIDs {
				res, createErr := s.createCustomSchedule(ctx, creatorID, assigneeID, req.LocationID, start, end)
				if createErr != nil {
					return nil, createErr
				}
				results = append(results, *res)
				s.sendNotificationForNewSchedule(ctx, res.ID, creatorID, assigneeID, res.StartDatetime, res.EndDatetime, res.LocationName)
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
		s.logError(ctx, "CreateSchedule", "invalid shift_date format", err, zap.String("shift_date", *req.ShiftDate))
		return nil, fmt.Errorf("invalid shift_date format: %w", err)
	}

	dates := s.buildCustomScheduleDates(baseDate, recurrence)
	for _, date := range dates {
		for _, assigneeID := range req.EmployeeIDs {
			res, createErr := s.createPresetScheduleForDate(ctx, creatorID, assigneeID, req.LocationID, req.LocationShiftID, locationShift, date, locationTZ)
			if createErr != nil {
				return nil, createErr
			}
			results = append(results, *res)
			s.sendNotificationForNewSchedule(ctx, res.ID, creatorID, assigneeID, res.StartDatetime, res.EndDatetime, res.LocationName)
		}
	}
	return results, nil
}

func (s *ScheduleService) GetSchedulesByLocationInRange(ctx context.Context, locationID uuid.UUID, req *domain.GetSchedulesByLocationInRangeRequest) ([]domain.GetSchedulesByLocationInRangeResponse, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		s.logError(ctx, "GetSchedulesByLocationInRange", "invalid start_date format", err, zap.String("start_date", req.StartDate))
		return nil, fmt.Errorf("invalid start_date format, expected YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		s.logError(ctx, "GetSchedulesByLocationInRange", "invalid end_date format", err, zap.String("end_date", req.EndDate))
		return nil, fmt.Errorf("invalid end_date format, expected YYYY-MM-DD")
	}

	if endDate.Before(startDate) {
		s.logError(ctx, "GetSchedulesByLocationInRange", "end_date is before start_date", nil, zap.String("start_date", req.StartDate), zap.String("end_date", req.EndDate))
		return nil, fmt.Errorf("end_date must be on or after start_date")
	}

	rows, err := s.repository.GetSchedulesByLocationInRange(ctx, locationID, startDate, endDate)
	if err != nil {
		s.logError(ctx, "GetSchedulesByLocationInRange", "failed to list schedules by range", err)
		return nil, fmt.Errorf("failed to list schedules by range: %w", err)
	}
	return rows, nil
}

func (s *ScheduleService) GetScheduleByID(ctx context.Context, scheduleID uuid.UUID) (*domain.GetScheduleByIdResponse, error) {
	item, err := s.repository.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		s.logError(ctx, "GetScheduleByID", "failed to get schedule by id", err)
		return nil, fmt.Errorf("failed to get schedule by ID: %w", err)
	}
	return item, nil
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, scheduleID uuid.UUID, updaterEmployeeID uuid.UUID, req *domain.UpdateScheduleRequest) (*domain.UpdateScheduleResponse, error) {
	existingSchedule, err := s.repository.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		s.logError(ctx, "UpdateSchedule", "failed to fetch existing schedule", err)
		return nil, fmt.Errorf("failed to fetch existing schedule: %w", err)
	}

	isCustom := s.determineScheduleType(req, existingSchedule)
	var res *domain.UpdateScheduleResponse
	if isCustom {
		res, err = s.updateCustomSchedule(ctx, scheduleID, existingSchedule, req)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = s.updatePresetSchedule(ctx, scheduleID, existingSchedule, req)
		if err != nil {
			return nil, err
		}
	}

	s.sendNotificationForUpdatedSchedule(ctx, res.ID, updaterEmployeeID, res.EmployeeID, res.StartDatetime, res.EndDatetime, res.LocationName)
	return res, nil
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, scheduleID uuid.UUID) error {
	if err := s.repository.DeleteSchedule(ctx, scheduleID); err != nil {
		s.logError(ctx, "DeleteSchedule", "failed to delete schedule", err)
		return fmt.Errorf("failed to delete schedule: %w", err)
	}
	return nil
}

func (s *ScheduleService) validateCreateScheduleRequest(req *domain.CreateScheduleRequest) error {
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

func (s *ScheduleService) resolveRecurrence(req *domain.CreateScheduleRequest) (string, error) {
	if req.Recurrence == nil || *req.Recurrence == "" {
		return domain.CreateScheduleRecurrenceNone, nil
	}

	switch *req.Recurrence {
	case domain.CreateScheduleRecurrenceNone, domain.CreateScheduleRecurrenceEndOfWeek, domain.CreateScheduleRecurrenceEndOfMonth:
		return *req.Recurrence, nil
	default:
		return "", fmt.Errorf("recurrence must be one of: none, end_of_week, end_of_month")
	}
}

func (s *ScheduleService) buildCustomScheduleDates(baseDate time.Time, recurrence string) []time.Time {
	dayStart := time.Date(baseDate.Year(), baseDate.Month(), baseDate.Day(), 0, 0, 0, 0, baseDate.Location())
	endDate := dayStart

	switch recurrence {
	case domain.CreateScheduleRecurrenceEndOfWeek:
		daysUntilSunday := (7 - int(dayStart.Weekday())) % 7
		endDate = dayStart.AddDate(0, 0, daysUntilSunday)
	case domain.CreateScheduleRecurrenceEndOfMonth:
		endDate = time.Date(dayStart.Year(), dayStart.Month()+1, 0, 0, 0, 0, 0, dayStart.Location())
	}

	dates := make([]time.Time, 0)
	for date := dayStart; !date.After(endDate); date = date.AddDate(0, 0, 1) {
		dates = append(dates, date)
	}
	return dates
}

func (s *ScheduleService) validateCustomSchedule(req *domain.CreateScheduleRequest) error {
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

func (s *ScheduleService) validatePresetSchedule(req *domain.CreateScheduleRequest) error {
	if req.LocationShiftID == nil || req.ShiftDate == nil {
		return fmt.Errorf("location_shift_id and shift_date are required for preset shift schedules")
	}
	if req.StartDatetime != nil || req.EndDatetime != nil {
		return fmt.Errorf("start_datetime and end_datetime should not be provided for preset shift schedules")
	}
	return nil
}

func (s *ScheduleService) createCustomSchedule(ctx context.Context, creatorID, assigneeID, locationID uuid.UUID, startDatetime, endDatetime time.Time) (*domain.CreateScheduleResponse, error) {
	schedule, err := s.repository.CreateSchedule(ctx, domain.CreateScheduleParams{
		EmployeeID:             assigneeID,
		LocationID:             locationID,
		IsCustom:               true,
		LocationShiftID:        nil,
		ShiftNameSnapshot:      nil,
		ShiftStartTimeSnapshot: nil,
		ShiftEndTimeSnapshot:   nil,
		CreatedByEmployeeID:    creatorID,
		StartDatetime:          startDatetime,
		EndDatetime:            endDatetime,
	})
	if err != nil {
		s.logError(ctx, "createCustomSchedule", "failed to create custom schedule", err)
		return nil, fmt.Errorf("failed to create custom schedule: %w", err)
	}
	return schedule, nil
}

func (s *ScheduleService) getPresetScheduleContext(ctx context.Context, req *domain.CreateScheduleRequest) (*domain.ScheduleLocationShift, *time.Location, error) {
	if err := s.validatePresetSchedule(req); err != nil {
		s.logError(ctx, "getPresetScheduleContext", "preset schedule validation failed", err)
		return nil, nil, err
	}

	locationShift, err := s.repository.GetShiftByID(ctx, *req.LocationShiftID)
	if err != nil {
		s.logError(ctx, "getPresetScheduleContext", "failed to fetch location shift", err)
		return nil, nil, fmt.Errorf("failed to fetch location shift: %w", err)
	}
	if locationShift.LocationID != req.LocationID {
		return nil, nil, fmt.Errorf("location shift does not belong to the specified location")
	}

	location, err := s.repository.GetLocationByID(ctx, req.LocationID)
	if err != nil {
		s.logError(ctx, "getPresetScheduleContext", "failed to fetch location", err)
		return nil, nil, fmt.Errorf("failed to fetch location: %w", err)
	}
	locationTZ, err := time.LoadLocation(location.Timezone)
	if err != nil {
		s.logError(ctx, "getPresetScheduleContext", "invalid location timezone", err, zap.String("location_timezone", location.Timezone))
		return nil, nil, fmt.Errorf("invalid location timezone: %w", err)
	}

	return locationShift, locationTZ, nil
}

func (s *ScheduleService) createPresetScheduleForDate(ctx context.Context, creatorID, assigneeID, locationID uuid.UUID, locationShiftID *uuid.UUID, locationShift *domain.ScheduleLocationShift, shiftDate time.Time, locationTZ *time.Location) (*domain.CreateScheduleResponse, error) {
	shiftDate = time.Date(shiftDate.Year(), shiftDate.Month(), shiftDate.Day(), 0, 0, 0, 0, locationTZ)
	startHour, startMin, startSec, startNano := microsecondsToTimeComponents(locationShift.StartMicroseconds)
	endHour, endMin, endSec, endNano := microsecondsToTimeComponents(locationShift.EndMicroseconds)

	startDatetime := time.Date(shiftDate.Year(), shiftDate.Month(), shiftDate.Day(), startHour, startMin, startSec, startNano, locationTZ)
	endDatetime := time.Date(shiftDate.Year(), shiftDate.Month(), shiftDate.Day(), endHour, endMin, endSec, endNano, locationTZ)
	if locationShift.EndMicroseconds < locationShift.StartMicroseconds {
		endDatetime = endDatetime.AddDate(0, 0, 1)
	}

	shiftStart := locationShift.StartMicroseconds
	shiftEnd := locationShift.EndMicroseconds
	schedule, err := s.repository.CreateSchedule(ctx, domain.CreateScheduleParams{
		EmployeeID:             assigneeID,
		LocationID:             locationID,
		LocationShiftID:        locationShiftID,
		ShiftNameSnapshot:      &locationShift.ShiftName,
		ShiftStartTimeSnapshot: &shiftStart,
		ShiftEndTimeSnapshot:   &shiftEnd,
		IsCustom:               false,
		CreatedByEmployeeID:    creatorID,
		StartDatetime:          startDatetime,
		EndDatetime:            endDatetime,
	})
	if err != nil {
		s.logError(ctx, "createPresetSchedule", "failed to create preset schedule", err)
		return nil, fmt.Errorf("failed to create preset schedule: %w", err)
	}
	return schedule, nil
}

func (s *ScheduleService) sendNotificationForNewSchedule(ctx context.Context, scheduleID uuid.UUID, creatorID, recipientID uuid.UUID, startTime, endTime time.Time, locationName string) {
	if s.asynqClient == nil {
		return
	}
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
		s.logError(ctx, "sendNotificationForNewSchedule", "failed to enqueue new schedule notification", err, zap.String("schedule_id", scheduleID.String()))
	}
}

func (s *ScheduleService) determineScheduleType(req *domain.UpdateScheduleRequest, existingSchedule *domain.GetScheduleByIdResponse) bool {
	if req.IsCustom != nil {
		isCustom := *req.IsCustom
		if !isCustom {
			req.StartDatetime = nil
			req.EndDatetime = nil
		}
		return isCustom
	}
	if req.LocationShiftID != nil || req.ShiftDate != nil {
		req.StartDatetime = nil
		req.EndDatetime = nil
		return false
	}
	return existingSchedule.LocationShiftID == nil
}

func (s *ScheduleService) updateCustomSchedule(ctx context.Context, scheduleID uuid.UUID, existingSchedule *domain.GetScheduleByIdResponse, req *domain.UpdateScheduleRequest) (*domain.UpdateScheduleResponse, error) {
	if err := s.validateCustomScheduleUpdate(req); err != nil {
		s.logError(ctx, "updateCustomSchedule", "custom schedule validation failed", err)
		return nil, err
	}

	employeeID := existingSchedule.EmployeeID
	locationID := existingSchedule.LocationID
	startDatetime := existingSchedule.StartDatetime
	endDatetime := existingSchedule.EndDatetime

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
	if startDatetime.After(endDatetime) {
		return nil, fmt.Errorf("start_datetime must be before end_datetime")
	}

	schedule, err := s.repository.UpdateSchedule(ctx, scheduleID, domain.UpdateScheduleParams{
		EmployeeID:             employeeID,
		LocationID:             locationID,
		LocationShiftID:        nil,
		ShiftNameSnapshot:      nil,
		ShiftStartTimeSnapshot: nil,
		ShiftEndTimeSnapshot:   nil,
		IsCustom:               true,
		StartDatetime:          startDatetime,
		EndDatetime:            endDatetime,
	})
	if err != nil {
		s.logError(ctx, "updateCustomSchedule", "failed to update custom schedule", err)
		return nil, fmt.Errorf("failed to update custom schedule: %w", err)
	}
	return schedule, nil
}

func (s *ScheduleService) updatePresetSchedule(ctx context.Context, scheduleID uuid.UUID, existingSchedule *domain.GetScheduleByIdResponse, req *domain.UpdateScheduleRequest) (*domain.UpdateScheduleResponse, error) {
	employeeID := existingSchedule.EmployeeID
	locationID := existingSchedule.LocationID
	if req.EmployeeID != nil {
		employeeID = *req.EmployeeID
	}
	if req.LocationID != nil {
		locationID = *req.LocationID
	}

	var shiftIDToUse uuid.UUID
	if req.LocationShiftID != nil {
		shiftIDToUse = *req.LocationShiftID
	} else if existingSchedule.LocationShiftID != nil {
		shiftIDToUse = *existingSchedule.LocationShiftID
	} else {
		return nil, fmt.Errorf("location_shift_id is required for preset shift schedules")
	}

	shiftDateToUse := existingSchedule.StartDatetime.Format("2006-01-02")
	if req.ShiftDate != nil {
		shiftDateToUse = *req.ShiftDate
	}

	locationShift, err := s.repository.GetShiftByID(ctx, shiftIDToUse)
	if err != nil {
		s.logError(ctx, "updatePresetSchedule", "failed to fetch location shift", err)
		return nil, fmt.Errorf("invalid location_shift_id: %w", err)
	}
	if locationShift.LocationID != locationID {
		return nil, fmt.Errorf("location_shift_id does not belong to the specified location")
	}

	shiftDate, err := time.Parse("2006-01-02", shiftDateToUse)
	if err != nil {
		return nil, fmt.Errorf("invalid shift_date format, expected YYYY-MM-DD: %w", err)
	}

	location, err := s.repository.GetLocationByID(ctx, locationID)
	if err != nil {
		s.logError(ctx, "updatePresetSchedule", "failed to fetch location", err)
		return nil, fmt.Errorf("failed to fetch location: %w", err)
	}
	locationTZ, err := time.LoadLocation(location.Timezone)
	if err != nil {
		s.logError(ctx, "updatePresetSchedule", "invalid location timezone", err, zap.String("location_timezone", location.Timezone))
		return nil, fmt.Errorf("invalid location timezone: %w", err)
	}

	shiftDate = time.Date(shiftDate.Year(), shiftDate.Month(), shiftDate.Day(), 0, 0, 0, 0, locationTZ)
	startHour, startMin, startSec, startNano := microsecondsToTimeComponents(locationShift.StartMicroseconds)
	endHour, endMin, endSec, endNano := microsecondsToTimeComponents(locationShift.EndMicroseconds)
	startDatetime := time.Date(shiftDate.Year(), shiftDate.Month(), shiftDate.Day(), startHour, startMin, startSec, startNano, locationTZ)
	endDatetime := time.Date(shiftDate.Year(), shiftDate.Month(), shiftDate.Day(), endHour, endMin, endSec, endNano, locationTZ)
	if locationShift.EndMicroseconds < locationShift.StartMicroseconds {
		endDatetime = endDatetime.AddDate(0, 0, 1)
	}

	shiftStart := locationShift.StartMicroseconds
	shiftEnd := locationShift.EndMicroseconds
	schedule, err := s.repository.UpdateSchedule(ctx, scheduleID, domain.UpdateScheduleParams{
		EmployeeID:             employeeID,
		LocationID:             locationID,
		LocationShiftID:        &shiftIDToUse,
		ShiftNameSnapshot:      &locationShift.ShiftName,
		ShiftStartTimeSnapshot: &shiftStart,
		ShiftEndTimeSnapshot:   &shiftEnd,
		IsCustom:               false,
		StartDatetime:          startDatetime,
		EndDatetime:            endDatetime,
	})
	if err != nil {
		s.logError(ctx, "updatePresetSchedule", "failed to update preset schedule", err)
		return nil, fmt.Errorf("failed to update preset schedule: %w", err)
	}
	return schedule, nil
}

func (s *ScheduleService) validateCustomScheduleUpdate(req *domain.UpdateScheduleRequest) error {
	if req.LocationShiftID != nil || req.ShiftDate != nil {
		return fmt.Errorf("location_shift_id and shift_date should not be provided for custom schedules")
	}
	return nil
}

func (s *ScheduleService) sendNotificationForUpdatedSchedule(ctx context.Context, scheduleID uuid.UUID, updaterEmployeeID, recipientEmployeeID uuid.UUID, startTime, endTime time.Time, locationName string) {
	if s.asynqClient == nil {
		return
	}
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
		s.logError(ctx, "sendNotificationForUpdatedSchedule", "failed to enqueue notification task", err)
	}
}

func (s *ScheduleService) AutoGenerateSchedules(ctx context.Context, req *domain.AutoGenerateSchedulesRequest) (*domain.AutoGenerateSchedulesResponse, error) {
	if req.LocationID == uuid.Nil {
		return nil, fmt.Errorf("location_id is required")
	}
	if req.Week < 1 || req.Week > 53 {
		return nil, fmt.Errorf("invalid week")
	}
	if req.Year <= 0 {
		return nil, fmt.Errorf("invalid year")
	}
	if len(req.EmployeeIDs) == 0 {
		return nil, fmt.Errorf("employee_ids is required")
	}

	employees, err := s.repository.ListEmployeesWithContractHours(ctx, req.EmployeeIDs)
	if err != nil {
		s.logError(ctx, "AutoGenerateSchedules", "failed to fetch employee contract hours", err)
		return nil, err
	}
	if len(employees) == 0 {
		return nil, fmt.Errorf("no employees found with contract hours")
	}

	locationShifts, err := s.repository.GetShiftsByLocationID(ctx, req.LocationID)
	if err != nil {
		s.logError(ctx, "AutoGenerateSchedules", "failed to fetch location shifts", err)
		return nil, err
	}
	if len(locationShifts) == 0 {
		return nil, fmt.Errorf("no shifts configured for location")
	}

	location, err := s.repository.GetLocationByID(ctx, req.LocationID)
	if err != nil {
		s.logError(ctx, "AutoGenerateSchedules", "failed to fetch location", err)
		return nil, err
	}
	locationTZ, err := time.LoadLocation(location.Timezone)
	if err != nil {
		s.logError(ctx, "AutoGenerateSchedules", "invalid location timezone", err, zap.String("timezone", location.Timezone))
		return nil, fmt.Errorf("invalid location timezone: %w", err)
	}

	if err := s.ensureWeekEmpty(ctx, req.LocationID, req.Week, req.Year, locationTZ); err != nil {
		return nil, err
	}

	return s.generateSchedulesWithORTools(ctx, req.LocationID, location.Timezone, locationTZ, employees, locationShifts, req.Week, req.Year)
}

type autoGenEmployee struct {
	ID            uuid.UUID
	FirstName     string
	LastName      string
	TargetMinutes int64
}

type autoGenShift struct {
	ID              uuid.UUID
	Name            string
	StartMinutes    int
	EndMinutes      int
	DurationMinutes int64
}

func (s *ScheduleService) generateSchedulesWithORTools(ctx context.Context, locationID uuid.UUID, timezone string, locationTZ *time.Location, employees []domain.ScheduleEmployeeContractHours, locationShifts []domain.ScheduleLocationShift, week int32, year int32) (*domain.AutoGenerateSchedulesResponse, error) {
	const (
		minStaffPerShift = int64(1)
		maxStaffPerShift = int64(2)
		maxSolveSeconds  = 90.0
		minRestMinutes   = int64(8 * 60)
	)

	weekStart, err := isoWeekStartDate(int(year), int(week), locationTZ)
	if err != nil {
		return nil, err
	}
	weekStartStr := weekStart.Format("2006-01-02")

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
		return nil, fmt.Errorf("no employees with positive contract hours")
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
		return nil, fmt.Errorf("no valid shifts found for location")
	}

	dCount := 7
	planEmployees := make([]domain.SchedulePlanEmployee, 0, len(inputsEmp))
	for _, emp := range inputsEmp {
		planEmployees = append(planEmployees, domain.SchedulePlanEmployee{
			ID:            emp.ID,
			FirstName:     emp.FirstName,
			LastName:      emp.LastName,
			TargetMinutes: emp.TargetMinutes,
		})
	}

	shiftTemplates := make([]domain.ScheduleShiftTemplate, 0, len(inputsShifts))
	for _, sh := range inputsShifts {
		overnight := sh.EndMinutes < sh.StartMinutes
		shiftTemplates = append(shiftTemplates, domain.ScheduleShiftTemplate{
			ShiftID:         sh.ID,
			Name:            sh.Name,
			StartMinute:     int32(sh.StartMinutes),
			EndMinute:       int32(sh.EndMinutes),
			DurationMinutes: sh.DurationMinutes,
			Overnight:       overnight,
		})
	}

	emptySlots := make([]domain.SchedulePlanSlot, 0, dCount*len(inputsShifts))
	for dIdx := 0; dIdx < dCount; dIdx++ {
		dateStr := weekStart.AddDate(0, 0, dIdx).Format("2006-01-02")
		for _, sh := range inputsShifts {
			emptySlots = append(emptySlots, domain.SchedulePlanSlot{Date: dateStr, ShiftID: sh.ID, EmployeeIDs: []uuid.UUID{}})
		}
	}

	minEmployeesPerDay := int(minStaffPerShift) * len(inputsShifts)
	if len(inputsEmp) < minEmployeesPerDay {
		return &domain.AutoGenerateSchedulesResponse{
			Status:        "infeasible",
			PlanID:        uuid.New(),
			LocationID:    locationID,
			Timezone:      timezone,
			Week:          week,
			Year:          year,
			WeekStartDate: weekStartStr,
			Constraints: domain.SchedulePlanConstraints{
				MaxStaffPerShift: int32(maxStaffPerShift),
				AllowEmptyShift:  true,
			},
			Employees:      planEmployees,
			ShiftTemplates: shiftTemplates,
			Slots:          emptySlots,
			Summary:        []domain.ScheduleEmployeeSummary{},
			Warnings:       []domain.SchedulePlanWarning{{Code: "INFEASIBLE", Message: "Not enough employees to staff all shifts (1 shift/day per employee)."}},
		}, nil
	}

	model := cpmodel.NewCpModelBuilder()
	eCount := len(inputsEmp)
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
			model.AddLinearConstraint(expr, minStaffPerShift, maxStaffPerShift)
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
					rest := startAbs - endAbs
					if rest < minRestMinutes {
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
	params := &sppb.SatParameters{MaxTimeInSeconds: proto.Float64(maxSolveSeconds), NumSearchWorkers: proto.Int32(int32(runtime.NumCPU()))}
	res, err := cpmodel.SolveCpModelWithParameters(modelProto, params)
	if err != nil {
		return nil, err
	}

	status := res.GetStatus()
	statusStr := "infeasible"
	if status == cmpb.CpSolverStatus_OPTIMAL {
		statusStr = "optimal"
	} else if status == cmpb.CpSolverStatus_FEASIBLE {
		statusStr = "feasible"
	}
	if statusStr == "infeasible" {
		return &domain.AutoGenerateSchedulesResponse{
			Status:         "infeasible",
			PlanID:         uuid.New(),
			LocationID:     locationID,
			Timezone:       timezone,
			Week:           week,
			Year:           year,
			WeekStartDate:  weekStartStr,
			Constraints:    domain.SchedulePlanConstraints{MaxStaffPerShift: int32(maxStaffPerShift), AllowEmptyShift: true},
			Employees:      planEmployees,
			ShiftTemplates: shiftTemplates,
			Slots:          emptySlots,
			Summary:        []domain.ScheduleEmployeeSummary{},
			Warnings:       []domain.SchedulePlanWarning{{Code: "INFEASIBLE", Message: "Solver could not find a feasible schedule within the time limit."}},
		}, nil
	}

	slots := make([]domain.SchedulePlanSlot, 0, dCount*sCount)
	assignedMinutesByEmp := make(map[uuid.UUID]int64, eCount)
	shiftCountsByEmp := make(map[uuid.UUID]map[uuid.UUID]int, eCount)
	for _, emp := range inputsEmp {
		shiftCountsByEmp[emp.ID] = make(map[uuid.UUID]int)
	}
	for dIdx := 0; dIdx < dCount; dIdx++ {
		dateStr := weekStart.AddDate(0, 0, dIdx).Format("2006-01-02")
		for shIdx, sh := range inputsShifts {
			employeeIDs := make([]uuid.UUID, 0, int(maxStaffPerShift))
			for eIdx, emp := range inputsEmp {
				if cpmodel.SolutionBooleanValue(res, assign[eIdx][dIdx][shIdx]) {
					employeeIDs = append(employeeIDs, emp.ID)
					assignedMinutesByEmp[emp.ID] += sh.DurationMinutes
					shiftCountsByEmp[emp.ID][sh.ID]++
				}
			}
			slots = append(slots, domain.SchedulePlanSlot{Date: dateStr, ShiftID: sh.ID, EmployeeIDs: employeeIDs})
		}
	}

	summaries := make([]domain.ScheduleEmployeeSummary, 0, eCount)
	for _, emp := range inputsEmp {
		assigned := assignedMinutesByEmp[emp.ID]
		overtime := int64(0)
		if assigned > emp.TargetMinutes {
			overtime = assigned - emp.TargetMinutes
		}
		summaries = append(summaries, domain.ScheduleEmployeeSummary{
			EmployeeID:      emp.ID,
			TargetMinutes:   emp.TargetMinutes,
			AssignedMinutes: assigned,
			OvertimeMinutes: overtime,
			ShiftCounts:     shiftCountsByEmp[emp.ID],
		})
	}

	return &domain.AutoGenerateSchedulesResponse{
		Status:         statusStr,
		PlanID:         uuid.New(),
		LocationID:     locationID,
		Timezone:       timezone,
		Week:           week,
		Year:           year,
		WeekStartDate:  weekStartStr,
		Constraints:    domain.SchedulePlanConstraints{MaxStaffPerShift: int32(maxStaffPerShift), AllowEmptyShift: true},
		Employees:      planEmployees,
		ShiftTemplates: shiftTemplates,
		Slots:          slots,
		Summary:        summaries,
	}, nil
}

func isoWeekStartDate(year int, week int, loc *time.Location) (time.Time, error) {
	if week < 1 || week > 53 {
		return time.Time{}, fmt.Errorf("invalid ISO week: %d", week)
	}
	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, loc)
	isoDow := int(jan4.Weekday())
	if isoDow == 0 {
		isoDow = 7
	}
	week1Start := jan4.AddDate(0, 0, -(isoDow - 1))
	return week1Start.AddDate(0, 0, (week-1)*7), nil
}

func (s *ScheduleService) ensureWeekEmpty(ctx context.Context, locationID uuid.UUID, week int32, year int32, locationTZ *time.Location) error {
	weekStart, err := isoWeekStartDate(int(year), int(week), locationTZ)
	if err != nil {
		return err
	}
	weekEnd := weekStart.AddDate(0, 0, 6)
	rows, err := s.repository.GetSchedulesByLocationInRange(ctx, locationID, weekStart, weekEnd)
	if err != nil {
		return fmt.Errorf("failed to check existing schedules: %w", err)
	}
	if len(rows) > 0 {
		return domain.ErrWeekNotEmpty
	}
	return nil
}

func (s *ScheduleService) SaveGeneratedSchedules(ctx context.Context, creatorID uuid.UUID, req *domain.SaveGeneratedSchedulesRequest) error {
	if len(req.Slots) == 0 {
		return nil
	}
	if req.PlanID == uuid.Nil {
		return fmt.Errorf("plan_id is required")
	}
	if req.LocationID == uuid.Nil {
		return fmt.Errorf("location_id is required")
	}
	if req.Week < 1 || req.Week > 53 {
		return fmt.Errorf("invalid week")
	}
	if req.Year <= 0 {
		return fmt.Errorf("invalid year")
	}

	location, err := s.repository.GetLocationByID(ctx, req.LocationID)
	if err != nil {
		s.logError(ctx, "SaveGeneratedSchedules", "failed to fetch location", err, zap.String("location_id", req.LocationID.String()))
		return err
	}
	locationTZ, err := time.LoadLocation(location.Timezone)
	if err != nil {
		s.logError(ctx, "SaveGeneratedSchedules", "invalid location timezone", err, zap.String("timezone", location.Timezone))
		return fmt.Errorf("invalid location timezone: %w", err)
	}
	if err := s.ensureWeekEmpty(ctx, req.LocationID, req.Week, req.Year, locationTZ); err != nil {
		return err
	}

	weekStart, err := isoWeekStartDate(int(req.Year), int(req.Week), locationTZ)
	if err != nil {
		return err
	}
	weekEnd := weekStart.AddDate(0, 0, 6)
	_ = weekEnd

	locationShifts, err := s.repository.GetShiftsByLocationID(ctx, req.LocationID)
	if err != nil {
		return err
	}
	shiftIDSet := make(map[uuid.UUID]struct{}, len(locationShifts))
	shiftByID := make(map[uuid.UUID]domain.ScheduleLocationShift, len(locationShifts))
	for _, sh := range locationShifts {
		shiftIDSet[sh.ID] = struct{}{}
		shiftByID[sh.ID] = sh
	}

	maxStaffPerShift := 2
	seenSlot := make(map[string]struct{}, len(req.Slots))
	for _, slot := range req.Slots {
		if slot.Date == "" {
			return fmt.Errorf("slot.date is required")
		}
		if slot.ShiftID == uuid.Nil {
			return fmt.Errorf("slot.shift_id is required")
		}
		if _, ok := shiftIDSet[slot.ShiftID]; !ok {
			return fmt.Errorf("slot.shift_id does not belong to location")
		}
		slotKey := slot.Date + ":" + slot.ShiftID.String()
		if _, ok := seenSlot[slotKey]; ok {
			return fmt.Errorf("duplicate slot: %s", slotKey)
		}
		seenSlot[slotKey] = struct{}{}

		date, err := time.ParseInLocation("2006-01-02", slot.Date, locationTZ)
		if err != nil {
			return fmt.Errorf("invalid slot.date format, expected YYYY-MM-DD")
		}
		dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, locationTZ)
		if dayStart.Before(weekStart) || dayStart.After(weekEnd) {
			return fmt.Errorf("slot.date is outside requested ISO week")
		}

		if len(slot.EmployeeIDs) > maxStaffPerShift {
			return fmt.Errorf("too many employees for slot (max %d)", maxStaffPerShift)
		}
		empSeen := make(map[uuid.UUID]struct{}, len(slot.EmployeeIDs))
		for _, eid := range slot.EmployeeIDs {
			if eid == uuid.Nil {
				return fmt.Errorf("slot.employee_ids contains invalid uuid")
			}
			if _, ok := empSeen[eid]; ok {
				return fmt.Errorf("slot.employee_ids must not contain duplicates")
			}
			empSeen[eid] = struct{}{}
		}
	}

	assignByEmpByDate := make(map[uuid.UUID]map[string]uuid.UUID)
	for _, slot := range req.Slots {
		for _, eid := range slot.EmployeeIDs {
			m, ok := assignByEmpByDate[eid]
			if !ok {
				m = make(map[string]uuid.UUID)
				assignByEmpByDate[eid] = m
			}
			if _, exists := m[slot.Date]; exists {
				return fmt.Errorf("employee has more than one shift on %s", slot.Date)
			}
			m[slot.Date] = slot.ShiftID
		}
	}

	dateList := make([]string, 0, 7)
	for d := 0; d < 7; d++ {
		dateList = append(dateList, weekStart.AddDate(0, 0, d).Format("2006-01-02"))
	}
	for _, byDate := range assignByEmpByDate {
		for d := 0; d < 6; d++ {
			curDate := dateList[d]
			nextDate := dateList[d+1]
			curShiftID, okA := byDate[curDate]
			nextShiftID, okB := byDate[nextDate]
			if !okA || !okB {
				continue
			}
			curShift := shiftByID[curShiftID]
			nextShift := shiftByID[nextShiftID]
			curStartMin := curShift.StartMicroseconds / (60 * 1_000_000)
			curEndMin := curShift.EndMicroseconds / (60 * 1_000_000)
			curEndAbs := int64(d)*1440 + curEndMin
			if curEndMin < curStartMin {
				curEndAbs += 1440
			}
			nextStartMin := nextShift.StartMicroseconds / (60 * 1_000_000)
			nextStartAbs := int64(d+1)*1440 + nextStartMin
			if nextStartAbs-curEndAbs < int64(8*60) {
				return fmt.Errorf("minimum rest violation between %s and %s", curDate, nextDate)
			}
		}
	}

	for _, slot := range req.Slots {
		if len(slot.EmployeeIDs) == 0 {
			continue
		}
		shiftID := slot.ShiftID
		shiftDate := slot.Date
		if _, err := s.CreateSchedule(ctx, creatorID, &domain.CreateScheduleRequest{
			EmployeeIDs:     slot.EmployeeIDs,
			LocationID:      req.LocationID,
			IsCustom:        false,
			LocationShiftID: &shiftID,
			ShiftDate:       &shiftDate,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *ScheduleService) logError(ctx context.Context, operation, message string, err error, fields ...zap.Field) {
	if s.logger == nil {
		return
	}
	s.logger.LogError(ctx, "ScheduleService."+operation, message, err, fields...)
}

func microsecondsToTimeComponents(microseconds int64) (hour, min, sec, nano int) {
	totalSeconds := microseconds / 1_000_000
	hour = int(totalSeconds / 3600)
	min = int((totalSeconds % 3600) / 60)
	sec = int(totalSeconds % 60)
	nano = int((microseconds % 1_000_000) * 1000)
	return
}
