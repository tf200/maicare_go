package late_arrival

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *lateArrivalService) CreateLateArrival(
	ctx context.Context,
	employeeID uuid.UUID,
	req *CreateLateArrivalRequest,
) (*CreateLateArrivalResponse, error) {
	return s.createLateArrival(ctx, employeeID, employeeID, req)
}

func (s *lateArrivalService) CreateLateArrivalByAdmin(
	ctx context.Context,
	adminEmployeeID uuid.UUID,
	req *CreateLateArrivalByAdminRequest,
) (*CreateLateArrivalResponse, error) {
	if req == nil || adminEmployeeID == uuid.Nil || req.EmployeeID == uuid.Nil {
		return nil, ErrLateArrivalInvalidRequest
	}

	payload := &CreateLateArrivalRequest{
		ArrivalDate: req.ArrivalDate,
		ArrivalTime: req.ArrivalTime,
		Reason:      req.Reason,
	}
	return s.createLateArrival(ctx, req.EmployeeID, adminEmployeeID, payload)
}

func (s *lateArrivalService) createLateArrival(
	ctx context.Context,
	employeeID uuid.UUID,
	createdByEmployeeID uuid.UUID,
	req *CreateLateArrivalRequest,
) (*CreateLateArrivalResponse, error) {
	if req == nil || employeeID == uuid.Nil || createdByEmployeeID == uuid.Nil {
		return nil, ErrLateArrivalInvalidRequest
	}

	arrivalDate, ok, err := util.ParseYYYYMMDD(req.ArrivalDate)
	if err != nil || !ok {
		return nil, ErrLateArrivalInvalidRequest
	}

	arrivalTimeValue, err := util.StringToPgTime(strings.TrimSpace(req.ArrivalTime))
	if err != nil {
		return nil, fmt.Errorf("%w: invalid arrival_time format", ErrLateArrivalInvalidRequest)
	}
	hour, minute, second, _ := util.MicrosecondsToTimeComponents(arrivalTimeValue.Microseconds)

	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return nil, fmt.Errorf("%w: reason is required", ErrLateArrivalInvalidRequest)
	}

	schedules, err := s.Store.ListAssignedSchedulesForEmployeeOnDate(ctx, db.ListAssignedSchedulesForEmployeeOnDateParams{
		EmployeeID: employeeID,
		ArrivalDate: pgtype.Date{
			Time:  arrivalDate,
			Valid: true,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve assigned shift: %w", err)
	}
	if len(schedules) == 0 {
		return nil, fmt.Errorf("%w: no assigned shift found for selected date", ErrLateArrivalInvalidRequest)
	}
	if len(schedules) > 1 {
		return nil, fmt.Errorf("%w: multiple assigned shifts found for selected date", ErrLateArrivalConflict)
	}

	resolved := schedules[0]
	locationTZ, err := time.LoadLocation(resolved.LocationTimezone)
	if err != nil {
		return nil, fmt.Errorf("failed to load location timezone: %w", err)
	}

	arrivalAtLocal := time.Date(
		arrivalDate.Year(),
		arrivalDate.Month(),
		arrivalDate.Day(),
		hour,
		minute,
		second,
		0,
		locationTZ,
	)
	shiftStartLocal := resolved.StartDatetime.Time.In(locationTZ)
	if !arrivalAtLocal.After(shiftStartLocal) {
		return nil, fmt.Errorf("%w: arrival_time must be after shift start", ErrLateArrivalInvalidRequest)
	}

	created, err := s.Store.CreateLateArrival(ctx, db.CreateLateArrivalParams{
		ScheduleID:          resolved.ScheduleID,
		EmployeeID:          employeeID,
		CreatedByEmployeeID: &createdByEmployeeID,
		ArrivalDate: pgtype.Date{
			Time:  arrivalDate,
			Valid: true,
		},
		ArrivalTime: arrivalTimeValue,
		Reason:      reason,
	})
	if err != nil {
		if isLateArrivalUniqueViolation(err) {
			return nil, fmt.Errorf("%w: late arrival already exists for this assigned shift", ErrLateArrivalConflict)
		}
		return nil, fmt.Errorf("failed to create late arrival: %w", err)
	}

	return mapCreateLateArrivalToResponse(created, resolved), nil
}

func (s *lateArrivalService) ListMyLateArrivals(
	ctx *gin.Context,
	employeeID uuid.UUID,
	req *ListMyLateArrivalsRequest,
) (*pagination.Response[LateArrivalListItem], error) {
	if req == nil {
		return nil, ErrLateArrivalInvalidRequest
	}

	dateFrom, dateTo, err := parseDateRange(req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}

	params := req.Request.GetParams()
	rows, err := s.Store.ListMyLateArrivalsPaginated(ctx, db.ListMyLateArrivalsPaginatedParams{
		EmployeeID: employeeID,
		DateFrom:   dateFrom,
		DateTo:     dateTo,
		Limit:      params.Limit,
		Offset:     params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list late arrivals: %w", err)
	}

	items := make([]LateArrivalListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLateArrivalListRowToResponse(
			row.ID,
			row.ScheduleID,
			row.EmployeeID,
			strings.TrimSpace(row.EmployeeFirstName+" "+row.EmployeeLastName),
			row.CreatedByEmployeeID,
			row.ArrivalDate,
			row.ArrivalTime,
			row.Reason,
			row.ShiftStartDatetime,
			row.ShiftEndDatetime,
			row.ShiftName,
			row.LocationName,
			row.CreatedAt,
			row.UpdatedAt,
		))
	}

	var totalCount int64
	if len(rows) > 0 {
		totalCount = rows[0].TotalCount
	}
	paginated := pagination.NewResponse(ctx, req.Request, items, totalCount)
	return &paginated, nil
}

func (s *lateArrivalService) ListLateArrivals(
	ctx *gin.Context,
	req *ListLateArrivalsRequest,
) (*pagination.Response[LateArrivalListItem], error) {
	if req == nil {
		return nil, ErrLateArrivalInvalidRequest
	}

	dateFrom, dateTo, err := parseDateRange(req.DateFrom, req.DateTo)
	if err != nil {
		return nil, err
	}

	var employeeSearch *string
	if req.EmployeeSearch != nil {
		search := strings.TrimSpace(*req.EmployeeSearch)
		if search != "" {
			employeeSearch = &search
		}
	}

	params := req.Request.GetParams()
	rows, err := s.Store.ListLateArrivalsPaginated(ctx, db.ListLateArrivalsPaginatedParams{
		EmployeeSearch: employeeSearch,
		DateFrom:       dateFrom,
		DateTo:         dateTo,
		Limit:          params.Limit,
		Offset:         params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list late arrivals: %w", err)
	}

	items := make([]LateArrivalListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLateArrivalListRowToResponse(
			row.ID,
			row.ScheduleID,
			row.EmployeeID,
			strings.TrimSpace(row.EmployeeFirstName+" "+row.EmployeeLastName),
			row.CreatedByEmployeeID,
			row.ArrivalDate,
			row.ArrivalTime,
			row.Reason,
			row.ShiftStartDatetime,
			row.ShiftEndDatetime,
			row.ShiftName,
			row.LocationName,
			row.CreatedAt,
			row.UpdatedAt,
		))
	}

	var totalCount int64
	if len(rows) > 0 {
		totalCount = rows[0].TotalCount
	}
	paginated := pagination.NewResponse(ctx, req.Request, items, totalCount)
	return &paginated, nil
}

func parseDateRange(dateFromValue, dateToValue *string) (pgtype.Date, pgtype.Date, error) {
	dateFrom := pgtype.Date{Valid: false}
	dateTo := pgtype.Date{Valid: false}

	if dateFromValue != nil && strings.TrimSpace(*dateFromValue) != "" {
		parsed, ok, err := util.ParseYYYYMMDD(strings.TrimSpace(*dateFromValue))
		if err != nil || !ok {
			return dateFrom, dateTo, ErrLateArrivalInvalidRequest
		}
		dateFrom = pgtype.Date{Time: parsed, Valid: true}
	}

	if dateToValue != nil && strings.TrimSpace(*dateToValue) != "" {
		parsed, ok, err := util.ParseYYYYMMDD(strings.TrimSpace(*dateToValue))
		if err != nil || !ok {
			return dateFrom, dateTo, ErrLateArrivalInvalidRequest
		}
		dateTo = pgtype.Date{Time: parsed, Valid: true}
	}

	if dateFrom.Valid && dateTo.Valid && dateTo.Time.Before(dateFrom.Time) {
		return dateFrom, dateTo, fmt.Errorf("%w: date_to must be on or after date_from", ErrLateArrivalInvalidRequest)
	}

	return dateFrom, dateTo, nil
}

func mapCreateLateArrivalToResponse(
	row db.LateArrival,
	scheduleRow db.ListAssignedSchedulesForEmployeeOnDateRow,
) *CreateLateArrivalResponse {
	return &CreateLateArrivalResponse{
		ID:                  row.ID,
		ScheduleID:          row.ScheduleID,
		EmployeeID:          row.EmployeeID,
		CreatedByEmployeeID: row.CreatedByEmployeeID,
		ArrivalDate:         row.ArrivalDate.Time,
		ArrivalTime:         util.PgTimeToString(row.ArrivalTime),
		Reason:              row.Reason,
		ShiftStartDatetime:  scheduleRow.StartDatetime.Time,
		ShiftEndDatetime:    scheduleRow.EndDatetime.Time,
		ShiftName:           scheduleRow.ShiftName,
		LocationName:        scheduleRow.LocationName,
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
	}
}

func mapLateArrivalListRowToResponse(
	id uuid.UUID,
	scheduleID uuid.UUID,
	employeeID uuid.UUID,
	employeeName string,
	createdByEmployeeID *uuid.UUID,
	arrivalDate pgtype.Date,
	arrivalTime pgtype.Time,
	reason string,
	shiftStartDatetime pgtype.Timestamptz,
	shiftEndDatetime pgtype.Timestamptz,
	shiftName string,
	locationName string,
	createdAt pgtype.Timestamptz,
	updatedAt pgtype.Timestamptz,
) LateArrivalListItem {
	return LateArrivalListItem{
		ID:                  id,
		ScheduleID:          scheduleID,
		EmployeeID:          employeeID,
		EmployeeName:        employeeName,
		CreatedByEmployeeID: createdByEmployeeID,
		ArrivalDate:         arrivalDate.Time,
		ArrivalTime:         util.PgTimeToString(arrivalTime),
		Reason:              reason,
		ShiftStartDatetime:  shiftStartDatetime.Time,
		ShiftEndDatetime:    shiftEndDatetime.Time,
		ShiftName:           shiftName,
		LocationName:        locationName,
		CreatedAt:           createdAt.Time,
		UpdatedAt:           updatedAt.Time,
	}
}

func isLateArrivalUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) &&
		pgErr.Code == "23505" &&
		strings.Contains(pgErr.ConstraintName, "late_arrivals_unique_schedule")
}
