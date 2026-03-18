package leave

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *leaveService) CreateLeaveRequest(
	ctx context.Context,
	employeeID uuid.UUID,
	req *CreateLeaveRequestRequest,
) (*CreateLeaveRequestResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("%w: leave request payload is required", ErrLeaveRequestInvalidRequest)
	}

	leaveType, err := parseLeaveRequestType(req.LeaveType)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"CreateLeaveRequest",
			"Invalid leave request type",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
			zap.String("LeaveType", req.LeaveType),
		)
		return nil, fmt.Errorf("%w: %v", ErrLeaveRequestInvalidRequest, err)
	}

	startDate, ok, err := util.ParseYYYYMMDD(req.StartDate)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"CreateLeaveRequest",
			"Failed to parse leave request start date",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
			zap.String("StartDate", req.StartDate),
		)
		return nil, fmt.Errorf("%w: %v", ErrLeaveRequestInvalidRequest, err)
	}
	if !ok {
		return nil, fmt.Errorf("%w: start_date is required", ErrLeaveRequestInvalidRequest)
	}

	endDate, ok, err := util.ParseYYYYMMDD(req.EndDate)
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"CreateLeaveRequest",
			"Failed to parse leave request end date",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
			zap.String("EndDate", req.EndDate),
		)
		return nil, fmt.Errorf("%w: %v", ErrLeaveRequestInvalidRequest, err)
	}
	if !ok {
		return nil, fmt.Errorf("%w: end_date is required", ErrLeaveRequestInvalidRequest)
	}
	if endDate.Before(startDate) {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"CreateLeaveRequest",
			"Invalid leave request date range",
			zap.String("EmployeeID", employeeID.String()),
			zap.String("StartDate", req.StartDate),
			zap.String("EndDate", req.EndDate),
		)
		return nil, fmt.Errorf("%w: end date must be on or after start date", ErrLeaveRequestInvalidRequest)
	}

	createdByEmployeeID := employeeID
	created, err := s.Store.CreateLeaveRequest(ctx, db.CreateLeaveRequestParams{
		EmployeeID:          employeeID,
		CreatedByEmployeeID: &createdByEmployeeID,
		LeaveType:           leaveType,
		StartDate:           pgtype.Date{Time: startDate, Valid: true},
		EndDate:             pgtype.Date{Time: endDate, Valid: true},
		Reason:              req.Reason,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"CreateLeaveRequest",
			"Failed to create leave request",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
		)
		return nil, fmt.Errorf("failed to create leave request: %w", err)
	}

	res := mapLeaveRequestToResponse(created)

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"CreateLeaveRequest",
		"Successfully created leave request",
		zap.String("EmployeeID", employeeID.String()),
		zap.String("LeaveRequestID", created.ID.String()),
		zap.String("LeaveType", string(created.LeaveType)),
	)

	return res, nil
}

func (s *leaveService) UpdateLeaveRequest(
	ctx context.Context,
	employeeID, leaveRequestID uuid.UUID,
	req *UpdateLeaveRequestRequest,
) (*UpdateLeaveRequestResponse, error) {
	if req == nil || employeeID == uuid.Nil || leaveRequestID == uuid.Nil {
		return nil, ErrLeaveRequestInvalidRequest
	}

	today := currentUTCDate()
	var updated db.LeaveRequest

	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		current, err := q.LockLeaveRequestByID(ctx, leaveRequestID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrLeaveRequestNotFound
			}
			return fmt.Errorf("failed to lock leave request: %w", err)
		}

		if current.EmployeeID != employeeID {
			return ErrLeaveRequestForbidden
		}
		if current.Status != db.LeaveRequestStatusEnumPending {
			return ErrLeaveRequestStateInvalid
		}
		if !current.StartDate.Time.After(today) {
			return ErrLeaveRequestStateInvalid
		}

		updateArg, finalStartDate, finalEndDate, hasUpdates, err := buildLeaveRequestUpdateParams(current, req.LeaveType, req.StartDate, req.EndDate, req.Reason)
		if err != nil {
			return err
		}
		if !hasUpdates {
			return ErrLeaveRequestInvalidRequest
		}
		if finalEndDate.Before(finalStartDate) {
			return fmt.Errorf("%w: end date must be on or after start date", ErrLeaveRequestInvalidRequest)
		}
		if !finalStartDate.After(today) {
			return ErrLeaveRequestStateInvalid
		}

		updateArg.ID = leaveRequestID
		updated, err = q.UpdateLeaveRequestEditableFields(ctx, updateArg)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrLeaveRequestNotFound
			}
			return fmt.Errorf("failed to update leave request: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	res := mapLeaveRequestToResponse(updated)
	return res, nil
}

func (s *leaveService) UpdateLeaveRequestByAdmin(
	ctx context.Context,
	adminEmployeeID, leaveRequestID uuid.UUID,
	req *UpdateLeaveRequestAdminRequest,
) (*UpdateLeaveRequestResponse, error) {
	if req == nil || adminEmployeeID == uuid.Nil || leaveRequestID == uuid.Nil {
		return nil, ErrLeaveRequestInvalidRequest
	}

	note := strings.TrimSpace(req.AdminUpdateNote)
	if note == "" {
		return nil, fmt.Errorf("%w: admin_update_note is required", ErrLeaveRequestInvalidRequest)
	}

	var updated db.LeaveRequest
	var changedFields []string

	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		current, err := q.LockLeaveRequestByID(ctx, leaveRequestID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrLeaveRequestNotFound
			}
			return fmt.Errorf("failed to lock leave request: %w", err)
		}

		if current.Status != db.LeaveRequestStatusEnumPending &&
			current.Status != db.LeaveRequestStatusEnumApproved &&
			current.Status != db.LeaveRequestStatusEnumRejected {
			return ErrLeaveRequestStateInvalid
		}

		updateArg, finalStartDate, finalEndDate, hasUpdates, err := buildLeaveRequestUpdateParams(current, req.LeaveType, req.StartDate, req.EndDate, req.Reason)
		if err != nil {
			return err
		}
		if !hasUpdates {
			return ErrLeaveRequestInvalidRequest
		}
		if finalEndDate.Before(finalStartDate) {
			return fmt.Errorf("%w: end date must be on or after start date", ErrLeaveRequestInvalidRequest)
		}

		updateArg.ID = leaveRequestID
		updated, err = q.UpdateLeaveRequestEditableFields(ctx, updateArg)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrLeaveRequestNotFound
			}
			return fmt.Errorf("failed to update leave request: %w", err)
		}

		changedFields = detectChangedLeaveFields(current, updated)
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.Logger.LogBusinessEvent(
		ctx,
		logger.LogLevelInfo,
		"UpdateLeaveRequestByAdmin",
		"Admin updated leave request",
		zap.String("LeaveRequestID", leaveRequestID.String()),
		zap.String("AdminEmployeeID", adminEmployeeID.String()),
		zap.String("AdminUpdateNote", note),
		zap.Strings("ChangedFields", changedFields),
	)

	res := mapLeaveRequestToResponse(updated)
	return res, nil
}

func parseLeaveRequestType(value string) (db.LeaveRequestTypeEnum, error) {
	switch db.LeaveRequestTypeEnum(value) {
	case db.LeaveRequestTypeEnumVacation,
		db.LeaveRequestTypeEnumPersonal,
		db.LeaveRequestTypeEnumSick,
		db.LeaveRequestTypeEnumPregnancy,
		db.LeaveRequestTypeEnumLate,
		db.LeaveRequestTypeEnumUnpaid,
		db.LeaveRequestTypeEnumOther:
		return db.LeaveRequestTypeEnum(value), nil
	default:
		return "", fmt.Errorf("%w: invalid leave type: %s", ErrLeaveRequestInvalidRequest, value)
	}
}

func mapLeaveRequestToResponse(row db.LeaveRequest) *CreateLeaveRequestResponse {
	return &CreateLeaveRequestResponse{
		ID:                  row.ID,
		EmployeeID:          row.EmployeeID,
		CreatedByEmployeeID: row.CreatedByEmployeeID,
		LeaveType:           string(row.LeaveType),
		Status:              string(row.Status),
		StartDate:           row.StartDate.Time,
		EndDate:             row.EndDate.Time,
		Reason:              row.Reason,
		DecisionNote:        row.DecisionNote,
		DecidedByEmployeeID: row.DecidedByEmployeeID,
		RequestedAt:         row.RequestedAt.Time,
		DecidedAt:           timestamptzPtr(row.DecidedAt),
		CancelledAt:         timestamptzPtr(row.CancelledAt),
		CreatedAt:           row.CreatedAt.Time,
		UpdatedAt:           row.UpdatedAt.Time,
	}
}

func currentUTCDate() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

func parseDateForUpdate(fieldName string, value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}
	parsed, ok, err := util.ParseYYYYMMDD(*value)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLeaveRequestInvalidRequest, err)
	}
	if !ok {
		return nil, fmt.Errorf("%w: %s is required", ErrLeaveRequestInvalidRequest, fieldName)
	}
	return &parsed, nil
}

func buildLeaveRequestUpdateParams(
	current db.LeaveRequest,
	leaveTypeValue *string,
	startDateValue *string,
	endDateValue *string,
	reasonValue *string,
) (db.UpdateLeaveRequestEditableFieldsParams, time.Time, time.Time, bool, error) {
	updateArg := db.UpdateLeaveRequestEditableFieldsParams{
		LeaveType: db.NullLeaveRequestTypeEnum{Valid: false},
		StartDate: pgtype.Date{Valid: false},
		EndDate:   pgtype.Date{Valid: false},
		Reason:    nil,
	}
	finalStartDate := current.StartDate.Time
	finalEndDate := current.EndDate.Time
	hasUpdates := false

	if leaveTypeValue != nil {
		parsedType, err := parseLeaveRequestType(strings.TrimSpace(*leaveTypeValue))
		if err != nil {
			return updateArg, finalStartDate, finalEndDate, false, fmt.Errorf("%w: %v", ErrLeaveRequestInvalidRequest, err)
		}
		updateArg.LeaveType = db.NullLeaveRequestTypeEnum{
			LeaveRequestTypeEnum: parsedType,
			Valid:                true,
		}
		hasUpdates = true
	}

	startDate, err := parseDateForUpdate("start_date", startDateValue)
	if err != nil {
		return updateArg, finalStartDate, finalEndDate, false, err
	}
	if startDate != nil {
		finalStartDate = *startDate
		updateArg.StartDate = pgtype.Date{Time: *startDate, Valid: true}
		hasUpdates = true
	}

	endDate, err := parseDateForUpdate("end_date", endDateValue)
	if err != nil {
		return updateArg, finalStartDate, finalEndDate, false, err
	}
	if endDate != nil {
		finalEndDate = *endDate
		updateArg.EndDate = pgtype.Date{Time: *endDate, Valid: true}
		hasUpdates = true
	}

	if reasonValue != nil {
		trimmed := strings.TrimSpace(*reasonValue)
		updateArg.Reason = &trimmed
		hasUpdates = true
	}

	return updateArg, finalStartDate, finalEndDate, hasUpdates, nil
}

func detectChangedLeaveFields(before, after db.LeaveRequest) []string {
	changed := make([]string, 0, 4)
	if before.LeaveType != after.LeaveType {
		changed = append(changed, "leave_type")
	}
	if !before.StartDate.Time.Equal(after.StartDate.Time) {
		changed = append(changed, "start_date")
	}
	if !before.EndDate.Time.Equal(after.EndDate.Time) {
		changed = append(changed, "end_date")
	}
	beforeReason := ""
	if before.Reason != nil {
		beforeReason = *before.Reason
	}
	afterReason := ""
	if after.Reason != nil {
		afterReason = *after.Reason
	}
	if beforeReason != afterReason {
		changed = append(changed, "reason")
	}
	return changed
}

func timestamptzPtr(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}

func (s *leaveService) ListMyLeaveRequests(
	ctx *gin.Context,
	employeeID uuid.UUID,
	req *ListMyLeaveRequestsRequest,
) (*pagination.Response[LeaveRequestListItem], error) {
	if req == nil {
		return nil, ErrLeaveRequestInvalidRequest
	}

	params := req.Request.GetParams()
	queryArg := db.ListMyLeaveRequestsPaginatedParams{
		EmployeeID: employeeID,
		Status:     db.NullLeaveRequestStatusEnum{Valid: false},
		Limit:      params.Limit,
		Offset:     params.Offset,
	}

	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		status, ok := parseLeaveRequestStatus(*req.Status)
		if !ok {
			return nil, ErrLeaveRequestInvalidRequest
		}
		queryArg.Status = db.NullLeaveRequestStatusEnum{
			LeaveRequestStatusEnum: status,
			Valid:                  true,
		}
	}

	rows, err := s.Store.ListMyLeaveRequestsPaginated(ctx, queryArg)
	if err != nil {
		return nil, fmt.Errorf("failed to list leave requests: %w", err)
	}

	items := make([]LeaveRequestListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveListRowToResponse(
			row.ID,
			row.EmployeeID,
			strings.TrimSpace(row.EmployeeFirstName+" "+row.EmployeeLastName),
			row.CreatedByEmployeeID,
			row.LeaveType,
			row.Status,
			row.StartDate,
			row.EndDate,
			row.Reason,
			row.DecisionNote,
			row.DecidedByEmployeeID,
			row.RequestedAt,
			row.DecidedAt,
			row.CancelledAt,
			row.CreatedAt,
			row.UpdatedAt,
		))
	}

	var totalCount int64
	if len(rows) > 0 {
		totalCount = rows[0].TotalCount
	}
	resp := pagination.NewResponse(ctx, req.Request, items, totalCount)
	return &resp, nil
}

func (s *leaveService) ListLeaveRequests(
	ctx *gin.Context,
	req *ListLeaveRequestsRequest,
) (*pagination.Response[LeaveRequestListItem], error) {
	if req == nil {
		return nil, ErrLeaveRequestInvalidRequest
	}

	params := req.Request.GetParams()
	queryArg := db.ListLeaveRequestsPaginatedParams{
		Status:     db.NullLeaveRequestStatusEnum{Valid: false},
		EmployeeID: req.EmployeeID,
		Limit:      params.Limit,
		Offset:     params.Offset,
	}

	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		status, ok := parseLeaveRequestStatus(*req.Status)
		if !ok {
			return nil, ErrLeaveRequestInvalidRequest
		}
		queryArg.Status = db.NullLeaveRequestStatusEnum{
			LeaveRequestStatusEnum: status,
			Valid:                  true,
		}
	}

	rows, err := s.Store.ListLeaveRequestsPaginated(ctx, queryArg)
	if err != nil {
		return nil, fmt.Errorf("failed to list leave requests: %w", err)
	}

	items := make([]LeaveRequestListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveListRowToResponse(
			row.ID,
			row.EmployeeID,
			strings.TrimSpace(row.EmployeeFirstName+" "+row.EmployeeLastName),
			row.CreatedByEmployeeID,
			row.LeaveType,
			row.Status,
			row.StartDate,
			row.EndDate,
			row.Reason,
			row.DecisionNote,
			row.DecidedByEmployeeID,
			row.RequestedAt,
			row.DecidedAt,
			row.CancelledAt,
			row.CreatedAt,
			row.UpdatedAt,
		))
	}

	var totalCount int64
	if len(rows) > 0 {
		totalCount = rows[0].TotalCount
	}
	resp := pagination.NewResponse(ctx, req.Request, items, totalCount)
	return &resp, nil
}

func parseLeaveRequestStatus(value string) (db.LeaveRequestStatusEnum, bool) {
	switch db.LeaveRequestStatusEnum(strings.TrimSpace(value)) {
	case db.LeaveRequestStatusEnumPending,
		db.LeaveRequestStatusEnumApproved,
		db.LeaveRequestStatusEnumRejected,
		db.LeaveRequestStatusEnumCancelled,
		db.LeaveRequestStatusEnumExpired:
		return db.LeaveRequestStatusEnum(strings.TrimSpace(value)), true
	default:
		return "", false
	}
}

func mapLeaveListRowToResponse(
	id uuid.UUID,
	employeeID uuid.UUID,
	employeeName string,
	createdByEmployeeID *uuid.UUID,
	leaveType db.LeaveRequestTypeEnum,
	status db.LeaveRequestStatusEnum,
	startDate pgtype.Date,
	endDate pgtype.Date,
	reason *string,
	decisionNote *string,
	decidedByEmployeeID *uuid.UUID,
	requestedAt pgtype.Timestamptz,
	decidedAt pgtype.Timestamptz,
	cancelledAt pgtype.Timestamptz,
	createdAt pgtype.Timestamptz,
	updatedAt pgtype.Timestamptz,
) LeaveRequestListItem {
	return LeaveRequestListItem{
		ID:                  id,
		EmployeeID:          employeeID,
		EmployeeName:        employeeName,
		CreatedByEmployeeID: createdByEmployeeID,
		LeaveType:           string(leaveType),
		Status:              string(status),
		StartDate:           startDate.Time,
		EndDate:             endDate.Time,
		Reason:              reason,
		DecisionNote:        decisionNote,
		DecidedByEmployeeID: decidedByEmployeeID,
		RequestedAt:         requestedAt.Time,
		DecidedAt:           timestamptzPtr(decidedAt),
		CancelledAt:         timestamptzPtr(cancelledAt),
		CreatedAt:           createdAt.Time,
		UpdatedAt:           updatedAt.Time,
	}
}
