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
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *leaveService) CreateLeaveRequest(
	ctx context.Context,
	employeeID uuid.UUID,
	req *CreateLeaveRequestRequest,
) (*CreateLeaveRequestResponse, error) {
	return s.createLeaveRequest(ctx, employeeID, employeeID, req)
}

func (s *leaveService) CreateLeaveRequestByAdmin(
	ctx context.Context,
	adminEmployeeID uuid.UUID,
	req *CreateLeaveRequestByAdminRequest,
) (*CreateLeaveRequestResponse, error) {
	if req == nil || req.EmployeeID == uuid.Nil || adminEmployeeID == uuid.Nil {
		return nil, fmt.Errorf("%w: leave request payload is required", ErrLeaveRequestInvalidRequest)
	}

	payload := &CreateLeaveRequestRequest{
		LeaveType: req.LeaveType,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Reason:    req.Reason,
	}
	return s.createLeaveRequest(ctx, req.EmployeeID, adminEmployeeID, payload)
}

func (s *leaveService) createLeaveRequest(
	ctx context.Context,
	employeeID uuid.UUID,
	createdByEmployeeID uuid.UUID,
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
	if err := ensureLeaveTypePolicyActive(ctx, s.Store, leaveType); err != nil {
		s.Logger.LogBusinessEvent(
			ctx,
			logger.LogLevelError,
			"CreateLeaveRequest",
			"Leave type is not active in policy",
			zap.Error(err),
			zap.String("EmployeeID", employeeID.String()),
			zap.String("LeaveType", string(leaveType)),
		)
		return nil, err
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
	policy, err := getActiveLeavePolicy(ctx, s.Store, leaveType)
	if err != nil {
		return nil, err
	}
	if err := validateDeductibleRange(policy, startDate, endDate); err != nil {
		return nil, err
	}

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

func (s *leaveService) DecideLeaveRequestByAdmin(
	ctx context.Context,
	adminEmployeeID, leaveRequestID uuid.UUID,
	req *DecideLeaveRequestRequest,
) (*DecideLeaveRequestResponse, error) {
	if req == nil || adminEmployeeID == uuid.Nil || leaveRequestID == uuid.Nil {
		return nil, ErrLeaveRequestInvalidRequest
	}

	decision := strings.TrimSpace(req.Decision)
	if decision != "approve" && decision != "reject" {
		return nil, ErrLeaveRequestInvalidRequest
	}

	var updated db.LeaveRequest
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		current, err := q.LockLeaveRequestByID(ctx, leaveRequestID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrLeaveRequestNotFound
			}
			return fmt.Errorf("failed to lock leave request: %w", err)
		}
		if current.Status != db.LeaveRequestStatusEnumPending {
			return ErrLeaveRequestStateInvalid
		}

		nextStatus := db.LeaveRequestStatusEnumRejected
		if decision == "approve" {
			nextStatus = db.LeaveRequestStatusEnumApproved

			policy, err := q.GetActiveLeavePolicyByType(ctx, current.LeaveType)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					return fmt.Errorf("%w: leave type is not enabled", ErrLeaveRequestInvalidRequest)
				}
				return fmt.Errorf("failed to load leave policy: %w", err)
			}

			if policy.DeductsBalance {
				start := dateOnlyUTC(current.StartDate.Time)
				end := dateOnlyUTC(current.EndDate.Time)
				if start.Year() != end.Year() {
					return fmt.Errorf("%w: leave date range must be within one year", ErrLeaveRequestInvalidRequest)
				}

				requestedDays := int32(end.Sub(start).Hours()/24) + 1
				if requestedDays <= 0 {
					return fmt.Errorf("%w: invalid leave duration", ErrLeaveRequestInvalidRequest)
				}

				year := int32(start.Year())
				if err := q.EnsureLeaveBalanceForYear(ctx, db.EnsureLeaveBalanceForYearParams{
					EmployeeID: current.EmployeeID,
					Year:       year,
				}); err != nil {
					return fmt.Errorf("failed to ensure leave balance row: %w", err)
				}

				balance, err := q.LockLeaveBalanceByEmployeeYear(ctx, db.LockLeaveBalanceByEmployeeYearParams{
					EmployeeID: current.EmployeeID,
					Year:       year,
				})
				if err != nil {
					if errors.Is(err, pgx.ErrNoRows) {
						return ErrLeaveBalanceInsufficient
					}
					return fmt.Errorf("failed to lock leave balance: %w", err)
				}

				extraRemaining := balance.ExtraTotalDays - balance.ExtraUsedDays
				legalRemaining := balance.LegalTotalDays - balance.LegalUsedDays
				totalRemaining := extraRemaining + legalRemaining
				if totalRemaining < requestedDays {
					return ErrLeaveBalanceInsufficient
				}

				extraToUse := minInt32(extraRemaining, requestedDays)
				legalToUse := requestedDays - extraToUse

				if _, err := q.ApplyLeaveBalanceDeduction(ctx, db.ApplyLeaveBalanceDeductionParams{
					ID:        balance.ID,
					ExtraDays: extraToUse,
					LegalDays: legalToUse,
				}); err != nil {
					return fmt.Errorf("failed to deduct leave balance: %w", err)
				}
			}
		}

		updated, err = q.UpdateLeaveRequestDecision(ctx, db.UpdateLeaveRequestDecisionParams{
			ID:                  leaveRequestID,
			Status:              nextStatus,
			DecisionNote:        req.DecisionNote,
			DecidedByEmployeeID: &adminEmployeeID,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrLeaveRequestNotFound
			}
			return fmt.Errorf("failed to update leave request decision: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	res := mapLeaveRequestToResponse(updated)
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

		updateArg, parsedLeaveType, finalStartDate, finalEndDate, hasUpdates, err := buildLeaveRequestUpdateParams(current, req.LeaveType, req.StartDate, req.EndDate, req.Reason)
		if err != nil {
			return err
		}
		if !hasUpdates {
			return ErrLeaveRequestInvalidRequest
		}
		if parsedLeaveType != nil {
			if err := ensureLeaveTypePolicyActive(ctx, q, *parsedLeaveType); err != nil {
				return err
			}
		}
		if finalEndDate.Before(finalStartDate) {
			return fmt.Errorf("%w: end date must be on or after start date", ErrLeaveRequestInvalidRequest)
		}
		if !finalStartDate.After(today) {
			return ErrLeaveRequestStateInvalid
		}
		effectiveType := current.LeaveType
		if parsedLeaveType != nil {
			effectiveType = *parsedLeaveType
		}
		policy, err := getActiveLeavePolicy(ctx, q, effectiveType)
		if err != nil {
			return err
		}
		if err := validateDeductibleRange(policy, finalStartDate, finalEndDate); err != nil {
			return err
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
			current.Status != db.LeaveRequestStatusEnumRejected {
			return ErrLeaveRequestStateInvalid
		}

		updateArg, parsedLeaveType, finalStartDate, finalEndDate, hasUpdates, err := buildLeaveRequestUpdateParams(current, req.LeaveType, req.StartDate, req.EndDate, req.Reason)
		if err != nil {
			return err
		}
		if !hasUpdates {
			return ErrLeaveRequestInvalidRequest
		}
		if parsedLeaveType != nil {
			if err := ensureLeaveTypePolicyActive(ctx, q, *parsedLeaveType); err != nil {
				return err
			}
		}
		if finalEndDate.Before(finalStartDate) {
			return fmt.Errorf("%w: end date must be on or after start date", ErrLeaveRequestInvalidRequest)
		}
		effectiveType := current.LeaveType
		if parsedLeaveType != nil {
			effectiveType = *parsedLeaveType
		}
		policy, err := getActiveLeavePolicy(ctx, q, effectiveType)
		if err != nil {
			return err
		}
		if err := validateDeductibleRange(policy, finalStartDate, finalEndDate); err != nil {
			return err
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
) (db.UpdateLeaveRequestEditableFieldsParams, *db.LeaveRequestTypeEnum, time.Time, time.Time, bool, error) {
	updateArg := db.UpdateLeaveRequestEditableFieldsParams{
		LeaveType: db.NullLeaveRequestTypeEnum{Valid: false},
		StartDate: pgtype.Date{Valid: false},
		EndDate:   pgtype.Date{Valid: false},
		Reason:    nil,
	}
	var parsedLeaveType *db.LeaveRequestTypeEnum
	finalStartDate := current.StartDate.Time
	finalEndDate := current.EndDate.Time
	hasUpdates := false

	if leaveTypeValue != nil {
		parsedType, err := parseLeaveRequestType(strings.TrimSpace(*leaveTypeValue))
		if err != nil {
			return updateArg, nil, finalStartDate, finalEndDate, false, fmt.Errorf("%w: %v", ErrLeaveRequestInvalidRequest, err)
		}
		updateArg.LeaveType = db.NullLeaveRequestTypeEnum{
			LeaveRequestTypeEnum: parsedType,
			Valid:                true,
		}
		parsedLeaveType = &parsedType
		hasUpdates = true
	}

	startDate, err := parseDateForUpdate("start_date", startDateValue)
	if err != nil {
		return updateArg, parsedLeaveType, finalStartDate, finalEndDate, false, err
	}
	if startDate != nil {
		finalStartDate = *startDate
		updateArg.StartDate = pgtype.Date{Time: *startDate, Valid: true}
		hasUpdates = true
	}

	endDate, err := parseDateForUpdate("end_date", endDateValue)
	if err != nil {
		return updateArg, parsedLeaveType, finalStartDate, finalEndDate, false, err
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

	return updateArg, parsedLeaveType, finalStartDate, finalEndDate, hasUpdates, nil
}

type leavePolicyReader interface {
	GetActiveLeavePolicyByType(ctx context.Context, leaveType db.LeaveRequestTypeEnum) (db.LeavePolicy, error)
}

func ensureLeaveTypePolicyActive(ctx context.Context, reader leavePolicyReader, leaveType db.LeaveRequestTypeEnum) error {
	_, err := getActiveLeavePolicy(ctx, reader, leaveType)
	return err
}

func getActiveLeavePolicy(ctx context.Context, reader leavePolicyReader, leaveType db.LeaveRequestTypeEnum) (db.LeavePolicy, error) {
	policy, err := reader.GetActiveLeavePolicyByType(ctx, leaveType)
	if err == nil {
		return policy, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return db.LeavePolicy{}, fmt.Errorf("%w: leave type is not enabled", ErrLeaveRequestInvalidRequest)
	}
	return db.LeavePolicy{}, fmt.Errorf("failed to load leave policy: %w", err)
}

func validateDeductibleRange(policy db.LeavePolicy, startDate, endDate time.Time) error {
	if !policy.DeductsBalance {
		return nil
	}
	if startDate.Year() != endDate.Year() {
		return fmt.Errorf("%w: leave date range must be within one year for deductible leave types", ErrLeaveRequestInvalidRequest)
	}
	return nil
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

func dateOnlyUTC(t time.Time) time.Time {
	u := t.UTC()
	return time.Date(u.Year(), u.Month(), u.Day(), 0, 0, 0, 0, time.UTC)
}

func minInt32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
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

func (s *leaveService) GetMyLeaveRequestStats(ctx context.Context, employeeID uuid.UUID) (*MyLeaveRequestStatsResponse, error) {
	stats, err := s.Store.GetMyLeaveRequestStats(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leave request stats: %w", err)
	}

	return &LeaveRequestStatsResponse{
		OpenRequests:     stats.OpenRequests,
		ApprovedRequests: stats.ApprovedRequests,
		RejectedRequests: stats.RejectedRequests,
		SicknessAbsence:  stats.SicknessAbsence,
	}, nil
}

func (s *leaveService) GetLeaveRequestStats(ctx context.Context) (*LeaveRequestStatsResponse, error) {
	stats, err := s.Store.GetLeaveRequestStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leave request stats: %w", err)
	}

	return &LeaveRequestStatsResponse{
		OpenRequests:     stats.OpenRequests,
		ApprovedRequests: stats.ApprovedRequests,
		RejectedRequests: stats.RejectedRequests,
		SicknessAbsence:  stats.SicknessAbsence,
	}, nil
}

func (s *leaveService) ListLeaveRequests(
	ctx *gin.Context,
	req *ListLeaveRequestsRequest,
) (*pagination.Response[LeaveRequestListItem], error) {
	if req == nil {
		return nil, ErrLeaveRequestInvalidRequest
	}

	params := req.Request.GetParams()
	var employeeSearch *string
	if req.EmployeeSearch != nil {
		trimmed := strings.TrimSpace(*req.EmployeeSearch)
		if trimmed != "" {
			employeeSearch = &trimmed
		}
	}

	queryArg := db.ListLeaveRequestsPaginatedParams{
		Status:         db.NullLeaveRequestStatusEnum{Valid: false},
		EmployeeSearch: employeeSearch,
		Limit:          params.Limit,
		Offset:         params.Offset,
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

func (s *leaveService) ListLeaveBalances(
	ctx *gin.Context,
	req *ListLeaveBalancesRequest,
) (*pagination.Response[LeaveBalanceListItem], error) {
	if req == nil {
		return nil, ErrLeaveRequestInvalidRequest
	}

	params := req.Request.GetParams()
	var employeeSearch *string
	if req.EmployeeSearch != nil {
		trimmed := strings.TrimSpace(*req.EmployeeSearch)
		if trimmed != "" {
			employeeSearch = &trimmed
		}
	}

	queryArg := db.ListLeaveBalancesPaginatedParams{
		EmployeeSearch: employeeSearch,
		Year:           nil,
		Limit:          params.Limit,
		Offset:         params.Offset,
	}
	if req.Year != nil {
		queryArg.Year = req.Year
	}

	rows, err := s.Store.ListLeaveBalancesPaginated(ctx, queryArg)
	if err != nil {
		return nil, fmt.Errorf("failed to list leave balances: %w", err)
	}

	items := make([]LeaveBalanceListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveBalanceListRow(
			row.ID,
			row.EmployeeID,
			strings.TrimSpace(row.EmployeeFirstName+" "+row.EmployeeLastName),
			row.Year,
			row.LegalTotalDays,
			row.ExtraTotalDays,
			row.LegalUsedDays,
			row.ExtraUsedDays,
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

func (s *leaveService) ListMyLeaveBalances(
	ctx *gin.Context,
	employeeID uuid.UUID,
	req *ListMyLeaveBalancesRequest,
) (*pagination.Response[LeaveBalanceListItem], error) {
	if req == nil {
		return nil, ErrLeaveRequestInvalidRequest
	}

	params := req.Request.GetParams()
	queryArg := db.ListMyLeaveBalancesPaginatedParams{
		EmployeeID: employeeID,
		Year:       nil,
		Limit:      params.Limit,
		Offset:     params.Offset,
	}
	if req.Year != nil {
		queryArg.Year = req.Year
	}

	rows, err := s.Store.ListMyLeaveBalancesPaginated(ctx, queryArg)
	if err != nil {
		return nil, fmt.Errorf("failed to list leave balances: %w", err)
	}

	items := make([]LeaveBalanceListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapLeaveBalanceListRow(
			row.ID,
			row.EmployeeID,
			strings.TrimSpace(row.EmployeeFirstName+" "+row.EmployeeLastName),
			row.Year,
			row.LegalTotalDays,
			row.ExtraTotalDays,
			row.LegalUsedDays,
			row.ExtraUsedDays,
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

func (s *leaveService) AdjustLeaveBalance(
	ctx context.Context,
	adminEmployeeID uuid.UUID,
	req *AdjustLeaveBalanceRequest,
) (*AdjustLeaveBalanceResponse, error) {
	if req == nil || adminEmployeeID == uuid.Nil || req.EmployeeID == uuid.Nil {
		return nil, ErrLeaveRequestInvalidRequest
	}
	if req.LegalDaysDelta == 0 && req.ExtraDaysDelta == 0 {
		return nil, fmt.Errorf("%w: at least one delta is required", ErrLeaveBalanceInvalidAdjust)
	}
	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return nil, fmt.Errorf("%w: reason is required", ErrLeaveBalanceInvalidAdjust)
	}

	var adjusted db.LeaveBalance
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		if err := q.EnsureLeaveBalanceForYear(ctx, db.EnsureLeaveBalanceForYearParams{
			EmployeeID: req.EmployeeID,
			Year:       req.Year,
		}); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23503" {
				return ErrLeaveRequestNotFound
			}
			return fmt.Errorf("failed to ensure leave balance row: %w", err)
		}

		current, err := q.LockLeaveBalanceByEmployeeYear(ctx, db.LockLeaveBalanceByEmployeeYearParams{
			EmployeeID: req.EmployeeID,
			Year:       req.Year,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return ErrLeaveRequestNotFound
			}
			return fmt.Errorf("failed to lock leave balance: %w", err)
		}

		nextLegalTotal := current.LegalTotalDays + req.LegalDaysDelta
		nextExtraTotal := current.ExtraTotalDays + req.ExtraDaysDelta
		if nextLegalTotal < 0 || nextExtraTotal < 0 {
			return fmt.Errorf("%w: totals cannot be negative", ErrLeaveBalanceInvalidAdjust)
		}
		if nextLegalTotal < current.LegalUsedDays || nextExtraTotal < current.ExtraUsedDays {
			return fmt.Errorf("%w: totals cannot be lower than already used days", ErrLeaveBalanceInvalidAdjust)
		}

		adjusted, err = q.ApplyLeaveBalanceTotalAdjustment(ctx, db.ApplyLeaveBalanceTotalAdjustmentParams{
			ID:             current.ID,
			LegalDaysDelta: req.LegalDaysDelta,
			ExtraDaysDelta: req.ExtraDaysDelta,
		})
		if err != nil {
			return fmt.Errorf("failed to apply leave balance adjustment: %w", err)
		}

		if _, err := q.CreateLeaveBalanceAdjustmentAudit(ctx, db.CreateLeaveBalanceAdjustmentAuditParams{
			LeaveBalanceID:       current.ID,
			EmployeeID:           req.EmployeeID,
			Year:                 req.Year,
			LegalDaysDelta:       req.LegalDaysDelta,
			ExtraDaysDelta:       req.ExtraDaysDelta,
			Reason:               reason,
			AdjustedByEmployeeID: adminEmployeeID,
			LegalTotalDaysBefore: current.LegalTotalDays,
			ExtraTotalDaysBefore: current.ExtraTotalDays,
			LegalTotalDaysAfter:  adjusted.LegalTotalDays,
			ExtraTotalDaysAfter:  adjusted.ExtraTotalDays,
		}); err != nil {
			return fmt.Errorf("failed to create leave balance adjustment audit: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &AdjustLeaveBalanceResponse{
		Balance: mapLeaveBalanceListRow(
			adjusted.ID,
			adjusted.EmployeeID,
			"",
			adjusted.Year,
			adjusted.LegalTotalDays,
			adjusted.ExtraTotalDays,
			adjusted.LegalUsedDays,
			adjusted.ExtraUsedDays,
			adjusted.CreatedAt,
			adjusted.UpdatedAt,
		),
	}, nil
}

func mapLeaveBalanceListRow(
	id uuid.UUID,
	employeeID uuid.UUID,
	employeeName string,
	year int32,
	legalTotalDays int32,
	extraTotalDays int32,
	legalUsedDays int32,
	extraUsedDays int32,
	createdAt pgtype.Timestamptz,
	updatedAt pgtype.Timestamptz,
) LeaveBalanceListItem {
	legalRemaining := legalTotalDays - legalUsedDays
	extraRemaining := extraTotalDays - extraUsedDays
	return LeaveBalanceListItem{
		ID:             id,
		EmployeeID:     employeeID,
		EmployeeName:   employeeName,
		Year:           year,
		LegalTotalDays: legalTotalDays,
		ExtraTotalDays: extraTotalDays,
		LegalUsedDays:  legalUsedDays,
		ExtraUsedDays:  extraUsedDays,
		LegalRemaining: legalRemaining,
		ExtraRemaining: extraRemaining,
		TotalRemaining: legalRemaining + extraRemaining,
		CreatedAt:      createdAt.Time,
		UpdatedAt:      updatedAt.Time,
	}
}
