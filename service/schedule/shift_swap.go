package schedule

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/notification"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

const (
	shiftSwapDecisionAccept  = "accept"
	shiftSwapDecisionReject  = "reject"
	shiftSwapDecisionApprove = "approve"
)

func (s *scheduleService) CreateShiftSwapRequest(
	ctx context.Context,
	requesterEmployeeID uuid.UUID,
	req *CreateShiftSwapRequest,
) (*CreateShiftSwapResponse, error) {
	if req.RequesterScheduleID == uuid.Nil || req.RecipientScheduleID == uuid.Nil || req.RecipientEmployeeID == uuid.Nil {
		return nil, ErrShiftSwapInvalidRequest
	}
	if requesterEmployeeID == req.RecipientEmployeeID || req.RequesterScheduleID == req.RecipientScheduleID {
		return nil, ErrShiftSwapInvalidRequest
	}
	if req.ExpiresAt != nil && !req.ExpiresAt.After(time.Now().UTC()) {
		return nil, ErrShiftSwapInvalidRequest
	}

	requesterSchedule, err := s.Store.GetScheduleForSwapValidation(ctx, req.RequesterScheduleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to fetch requester schedule: %w", err)
	}
	recipientSchedule, err := s.Store.GetScheduleForSwapValidation(ctx, req.RecipientScheduleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScheduleNotFound
		}
		return nil, fmt.Errorf("failed to fetch recipient schedule: %w", err)
	}

	now := time.Now().UTC()
	if requesterSchedule.StartDatetime.Time.Before(now) || recipientSchedule.StartDatetime.Time.Before(now) {
		return nil, ErrShiftSwapInvalidRequest
	}
	if requesterSchedule.EmployeeID != requesterEmployeeID || recipientSchedule.EmployeeID != req.RecipientEmployeeID {
		return nil, ErrShiftSwapScheduleOwnership
	}

	_ = s.Store.ExpirePendingShiftSwapRequests(ctx)

	createArg := db.CreateShiftSwapRequestParams{
		RequesterEmployeeID: requesterEmployeeID,
		RecipientEmployeeID: req.RecipientEmployeeID,
		RequesterScheduleID: req.RequesterScheduleID,
		RecipientScheduleID: req.RecipientScheduleID,
		Status:              db.ShiftSwapStatusEnumPendingRecipient,
		ExpiresAt:           pgtype.Timestamptz{Valid: false},
	}
	if req.ExpiresAt != nil {
		createArg.ExpiresAt = pgtype.Timestamptz{Time: req.ExpiresAt.UTC(), Valid: true}
	}

	created, err := s.Store.CreateShiftSwapRequest(ctx, createArg)
	if err != nil {
		if isShiftSwapUniqueViolation(err) {
			return nil, ErrShiftSwapDuplicateActiveRequest
		}
		return nil, fmt.Errorf("failed to create shift swap request: %w", err)
	}

	s.notifyEmployees(
		ctx,
		[]uuid.UUID{created.RecipientEmployeeID},
		"You have a new shift swap request.",
	)

	resp := mapCreatedShiftSwapToResponse(created, requesterEmployeeID)
	return &resp, nil
}

func (s *scheduleService) RespondToShiftSwapRequest(
	ctx context.Context,
	recipientEmployeeID, swapID uuid.UUID,
	req *RespondShiftSwapRequest,
) (*ShiftSwapResponse, error) {
	if swapID == uuid.Nil || recipientEmployeeID == uuid.Nil {
		return nil, ErrShiftSwapInvalidRequest
	}

	decision := strings.ToLower(strings.TrimSpace(req.Decision))
	nextStatus := db.ShiftSwapStatusEnumRecipientRejected
	if decision == shiftSwapDecisionAccept {
		nextStatus = db.ShiftSwapStatusEnumPendingAdmin
	} else if decision != shiftSwapDecisionReject {
		return nil, ErrShiftSwapInvalidRequest
	}

	updated, err := s.Store.UpdateShiftSwapStatusAfterRecipientResponse(ctx, db.UpdateShiftSwapStatusAfterRecipientResponseParams{
		Status:                nextStatus,
		RecipientResponseNote: req.Note,
		ID:                    swapID,
		RecipientEmployeeID:   recipientEmployeeID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			existing, getErr := s.Store.GetShiftSwapRequestByID(ctx, swapID)
			if getErr != nil {
				return nil, ErrShiftSwapNotFound
			}
			if existing.Status == db.ShiftSwapStatusEnumExpired ||
				(existing.ExpiresAt.Valid && !existing.ExpiresAt.Time.After(time.Now().UTC())) {
				return nil, ErrShiftSwapExpired
			}
			return nil, ErrShiftSwapStateInvalid
		}
		return nil, fmt.Errorf("failed to update recipient decision: %w", err)
	}

	details, err := s.Store.GetShiftSwapRequestDetailsByID(ctx, updated.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load updated shift swap details: %w", err)
	}

	if nextStatus == db.ShiftSwapStatusEnumPendingAdmin {
		s.notifyAdmins(ctx, "A shift swap request is pending admin approval.")
	}
	s.notifyEmployees(
		ctx,
		[]uuid.UUID{updated.RequesterEmployeeID},
		fmt.Sprintf("Your shift swap request was %s.", strings.ReplaceAll(string(nextStatus), "_", " ")),
	)

	resp := mapShiftSwapRowToResponse(details, recipientEmployeeID)
	return &resp, nil
}

func (s *scheduleService) AdminDecisionShiftSwapRequest(
	ctx context.Context,
	adminEmployeeID, swapID uuid.UUID,
	req *AdminDecisionShiftSwapRequest,
) (*ShiftSwapResponse, error) {
	if adminEmployeeID == uuid.Nil || swapID == uuid.Nil {
		return nil, ErrShiftSwapInvalidRequest
	}

	decision := strings.ToLower(strings.TrimSpace(req.Decision))
	switch decision {
	case shiftSwapDecisionApprove:
		var confirmed db.ShiftSwapRequest
		err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
			swapRow, lockErr := q.LockShiftSwapRequestForAdminDecision(ctx, swapID)
			if lockErr != nil {
				if errors.Is(lockErr, pgx.ErrNoRows) {
					return ErrShiftSwapNotFound
				}
				return lockErr
			}
			if swapRow.Status == db.ShiftSwapStatusEnumExpired {
				return ErrShiftSwapExpired
			}
			if swapRow.Status != db.ShiftSwapStatusEnumPendingAdmin {
				return ErrShiftSwapStateInvalid
			}
			if swapRow.ExpiresAt.Valid && !swapRow.ExpiresAt.Time.After(time.Now().UTC()) {
				return ErrShiftSwapExpired
			}

			schedules, lockSchedErr := q.LockSchedulesByIDsForSwap(ctx, []uuid.UUID{swapRow.RequesterScheduleID, swapRow.RecipientScheduleID})
			if lockSchedErr != nil {
				return lockSchedErr
			}
			if len(schedules) != 2 {
				return ErrScheduleNotFound
			}

			scheduleByID := make(map[uuid.UUID]db.LockSchedulesByIDsForSwapRow, 2)
			for _, sched := range schedules {
				scheduleByID[sched.ID] = sched
			}
			requesterSchedule, okReq := scheduleByID[swapRow.RequesterScheduleID]
			recipientSchedule, okRec := scheduleByID[swapRow.RecipientScheduleID]
			if !okReq || !okRec {
				return ErrScheduleNotFound
			}

			if requesterSchedule.EmployeeID != swapRow.RequesterEmployeeID || recipientSchedule.EmployeeID != swapRow.RecipientEmployeeID {
				return ErrShiftSwapScheduleOwnership
			}

			now := time.Now().UTC()
			if requesterSchedule.StartDatetime.Time.Before(now) || recipientSchedule.StartDatetime.Time.Before(now) {
				return ErrShiftSwapInvalidRequest
			}

			excludeIDs := []uuid.UUID{requesterSchedule.ID, recipientSchedule.ID}
			requesterOverlapCount, overlapErr := q.CountScheduleOverlapsForEmployee(ctx, db.CountScheduleOverlapsForEmployeeParams{
				EmployeeID:          swapRow.RequesterEmployeeID,
				ExcludedScheduleIds: excludeIDs,
				ConflictStart:       recipientSchedule.StartDatetime,
				ConflictEnd:         recipientSchedule.EndDatetime,
			})
			if overlapErr != nil {
				return overlapErr
			}
			if requesterOverlapCount > 0 {
				return ErrShiftSwapConflict
			}

			recipientOverlapCount, overlapErr := q.CountScheduleOverlapsForEmployee(ctx, db.CountScheduleOverlapsForEmployeeParams{
				EmployeeID:          swapRow.RecipientEmployeeID,
				ExcludedScheduleIds: excludeIDs,
				ConflictStart:       requesterSchedule.StartDatetime,
				ConflictEnd:         requesterSchedule.EndDatetime,
			})
			if overlapErr != nil {
				return overlapErr
			}
			if recipientOverlapCount > 0 {
				return ErrShiftSwapConflict
			}

			if updErr := q.UpdateScheduleEmployeeAssignment(ctx, db.UpdateScheduleEmployeeAssignmentParams{
				ID:         requesterSchedule.ID,
				EmployeeID: swapRow.RecipientEmployeeID,
			}); updErr != nil {
				return updErr
			}
			if updErr := q.UpdateScheduleEmployeeAssignment(ctx, db.UpdateScheduleEmployeeAssignmentParams{
				ID:         recipientSchedule.ID,
				EmployeeID: swapRow.RequesterEmployeeID,
			}); updErr != nil {
				return updErr
			}

			var note *string
			if req.Note != nil && strings.TrimSpace(*req.Note) != "" {
				trimmed := strings.TrimSpace(*req.Note)
				note = &trimmed
			}
			var confirmErr error
			confirmed, confirmErr = q.MarkShiftSwapConfirmed(ctx, db.MarkShiftSwapConfirmedParams{
				ID:                swapID,
				AdminDecisionNote: note,
				AdminEmployeeID:   &adminEmployeeID,
			})
			if confirmErr != nil {
				if errors.Is(confirmErr, pgx.ErrNoRows) {
					return ErrShiftSwapStateInvalid
				}
				return confirmErr
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		details, err := s.Store.GetShiftSwapRequestDetailsByID(ctx, confirmed.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load confirmed shift swap details: %w", err)
		}

		s.notifyEmployees(ctx, []uuid.UUID{details.RequesterEmployeeID, details.RecipientEmployeeID}, "Your shift swap request was approved and schedules were swapped.")
		resp := mapShiftSwapRowToResponse(details, adminEmployeeID)
		return &resp, nil

	case shiftSwapDecisionReject:
		updated, err := s.Store.UpdateShiftSwapAdminDecision(ctx, db.UpdateShiftSwapAdminDecisionParams{
			Status:            db.ShiftSwapStatusEnumAdminRejected,
			AdminDecisionNote: req.Note,
			AdminEmployeeID:   &adminEmployeeID,
			ID:                swapID,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				existing, getErr := s.Store.GetShiftSwapRequestByID(ctx, swapID)
				if getErr != nil {
					if errors.Is(getErr, pgx.ErrNoRows) {
						return nil, ErrShiftSwapNotFound
					}
					return nil, fmt.Errorf("failed to load shift swap request: %w", getErr)
				}
				if existing.Status == db.ShiftSwapStatusEnumExpired {
					return nil, ErrShiftSwapExpired
				}
				if existing.ExpiresAt.Valid && !existing.ExpiresAt.Time.After(time.Now().UTC()) {
					return nil, ErrShiftSwapExpired
				}
				return nil, ErrShiftSwapStateInvalid
			}
			return nil, fmt.Errorf("failed to reject shift swap request: %w", err)
		}

		details, err := s.Store.GetShiftSwapRequestDetailsByID(ctx, updated.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load rejected shift swap details: %w", err)
		}
		s.notifyEmployees(ctx, []uuid.UUID{details.RequesterEmployeeID, details.RecipientEmployeeID}, "Your shift swap request was rejected by admin.")
		resp := mapShiftSwapRowToResponse(details, adminEmployeeID)
		return &resp, nil

	default:
		return nil, ErrShiftSwapInvalidRequest
	}
}

func (s *scheduleService) ListMyShiftSwapRequests(ctx context.Context, employeeID uuid.UUID) ([]ShiftSwapResponse, error) {
	rows, err := s.Store.ListMyShiftSwapRequests(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to list shift swap requests: %w", err)
	}

	resp := make([]ShiftSwapResponse, 0, len(rows))
	for _, row := range rows {
		resp = append(resp, mapShiftSwapListRowToResponse(row, employeeID))
	}
	return resp, nil
}

func (s *scheduleService) ListShiftSwapRequests(
	ctx *gin.Context,
	req *ListShiftSwapRequestsRequest,
) (*pagination.Response[ShiftSwapResponse], error) {
	if req == nil {
		return nil, ErrShiftSwapInvalidRequest
	}

	params := req.Request.GetParams()
	queryArg := db.ListShiftSwapRequestsPaginatedParams{
		Limit:      params.Limit,
		Offset:     params.Offset,
		EmployeeID: req.EmployeeID,
		Status:     db.NullShiftSwapStatusEnum{Valid: false},
	}

	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		status, ok := parseShiftSwapStatus(*req.Status)
		if !ok {
			return nil, ErrShiftSwapInvalidRequest
		}
		queryArg.Status = db.NullShiftSwapStatusEnum{
			ShiftSwapStatusEnum: status,
			Valid:               true,
		}
	}

	rows, err := s.Store.ListShiftSwapRequestsPaginated(ctx, queryArg)
	if err != nil {
		return nil, fmt.Errorf("failed to list shift swap requests: %w", err)
	}

	items := make([]ShiftSwapResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, mapShiftSwapPaginatedRowToResponse(row, uuid.Nil))
	}

	var totalCount int64
	if len(rows) > 0 {
		totalCount = rows[0].TotalCount
	}
	paginated := pagination.NewResponse(ctx, req.Request, items, totalCount)
	return &paginated, nil
}

func isShiftSwapUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" &&
			(strings.Contains(pgErr.ConstraintName, "uq_shift_swap_active_requester_schedule") ||
				strings.Contains(pgErr.ConstraintName, "uq_shift_swap_active_recipient_schedule") ||
				strings.Contains(pgErr.ConstraintName, "uq_shift_swap_active_schedule_any"))
	}
	return false
}

func mapShiftSwapRowToResponse(row db.GetShiftSwapRequestDetailsByIDRow, viewerEmployeeID uuid.UUID) ShiftSwapResponse {
	requesterName := strings.TrimSpace(row.RequesterFirstName + " " + row.RequesterLastName)
	recipientName := strings.TrimSpace(row.RecipientFirstName + " " + row.RecipientLastName)

	response := ShiftSwapResponse{
		ID:                    row.ID,
		RequesterEmployeeID:   row.RequesterEmployeeID,
		RequesterEmployeeName: requesterName,
		RecipientEmployeeID:   row.RecipientEmployeeID,
		RecipientEmployeeName: recipientName,
		RequesterSchedule: ShiftSwapScheduleSnapshot{
			ID:            row.RequesterScheduleID,
			EmployeeID:    row.RequesterEmployeeID,
			EmployeeName:  requesterName,
			StartDatetime: row.RequesterScheduleStartDatetime.Time,
			EndDatetime:   row.RequesterScheduleEndDatetime.Time,
		},
		RecipientSchedule: ShiftSwapScheduleSnapshot{
			ID:            row.RecipientScheduleID,
			EmployeeID:    row.RecipientEmployeeID,
			EmployeeName:  recipientName,
			StartDatetime: row.RecipientScheduleStartDatetime.Time,
			EndDatetime:   row.RecipientScheduleEndDatetime.Time,
		},
		Status:      string(row.Status),
		RequestedAt: row.RequestedAt.Time,
	}

	if row.RecipientRespondedAt.Valid {
		t := row.RecipientRespondedAt.Time
		response.RecipientRespondedAt = &t
	}
	if row.AdminDecidedAt.Valid {
		t := row.AdminDecidedAt.Time
		response.AdminDecidedAt = &t
	}
	response.RecipientResponseNote = row.RecipientResponseNote
	response.AdminDecisionNote = row.AdminDecisionNote
	response.AdminEmployeeID = row.AdminEmployeeID
	if row.AdminFirstName != nil && row.AdminLastName != nil {
		name := strings.TrimSpace(*row.AdminFirstName + " " + *row.AdminLastName)
		response.AdminEmployeeName = &name
	}
	if row.ExpiresAt.Valid {
		t := row.ExpiresAt.Time
		response.ExpiresAt = &t
	}

	if viewerEmployeeID == row.RequesterEmployeeID {
		response.Direction = "sent"
	} else if viewerEmployeeID == row.RecipientEmployeeID {
		response.Direction = "received"
	}

	return response
}

func mapCreatedShiftSwapToResponse(row db.ShiftSwapRequest, viewerEmployeeID uuid.UUID) CreateShiftSwapResponse {
	response := CreateShiftSwapResponse{
		ID:                  row.ID,
		RequesterEmployeeID: row.RequesterEmployeeID,
		RecipientEmployeeID: row.RecipientEmployeeID,
		RequesterScheduleID: row.RequesterScheduleID,
		RecipientScheduleID: row.RecipientScheduleID,
		Status:              string(row.Status),
		RequestedAt:         row.RequestedAt.Time,
	}

	if row.ExpiresAt.Valid {
		t := row.ExpiresAt.Time
		response.ExpiresAt = &t
	}

	if viewerEmployeeID == row.RequesterEmployeeID {
		response.Direction = "sent"
	} else if viewerEmployeeID == row.RecipientEmployeeID {
		response.Direction = "received"
	}

	return response
}

func mapShiftSwapListRowToResponse(row db.ListMyShiftSwapRequestsRow, viewerEmployeeID uuid.UUID) ShiftSwapResponse {
	requesterName := strings.TrimSpace(row.RequesterFirstName + " " + row.RequesterLastName)
	recipientName := strings.TrimSpace(row.RecipientFirstName + " " + row.RecipientLastName)

	response := ShiftSwapResponse{
		ID:                    row.ID,
		RequesterEmployeeID:   row.RequesterEmployeeID,
		RequesterEmployeeName: requesterName,
		RecipientEmployeeID:   row.RecipientEmployeeID,
		RecipientEmployeeName: recipientName,
		RequesterSchedule: ShiftSwapScheduleSnapshot{
			ID:            row.RequesterScheduleID,
			EmployeeID:    row.RequesterEmployeeID,
			EmployeeName:  requesterName,
			StartDatetime: row.RequesterScheduleStartDatetime.Time,
			EndDatetime:   row.RequesterScheduleEndDatetime.Time,
		},
		RecipientSchedule: ShiftSwapScheduleSnapshot{
			ID:            row.RecipientScheduleID,
			EmployeeID:    row.RecipientEmployeeID,
			EmployeeName:  recipientName,
			StartDatetime: row.RecipientScheduleStartDatetime.Time,
			EndDatetime:   row.RecipientScheduleEndDatetime.Time,
		},
		Status:      string(row.Status),
		RequestedAt: row.RequestedAt.Time,
	}

	if row.RecipientRespondedAt.Valid {
		t := row.RecipientRespondedAt.Time
		response.RecipientRespondedAt = &t
	}
	if row.AdminDecidedAt.Valid {
		t := row.AdminDecidedAt.Time
		response.AdminDecidedAt = &t
	}
	response.RecipientResponseNote = row.RecipientResponseNote
	response.AdminDecisionNote = row.AdminDecisionNote
	response.AdminEmployeeID = row.AdminEmployeeID
	if row.AdminFirstName != nil && row.AdminLastName != nil {
		name := strings.TrimSpace(*row.AdminFirstName + " " + *row.AdminLastName)
		response.AdminEmployeeName = &name
	}
	if row.ExpiresAt.Valid {
		t := row.ExpiresAt.Time
		response.ExpiresAt = &t
	}
	if viewerEmployeeID == row.RequesterEmployeeID {
		response.Direction = "sent"
	} else if viewerEmployeeID == row.RecipientEmployeeID {
		response.Direction = "received"
	}

	return response
}

func mapShiftSwapPaginatedRowToResponse(row db.ListShiftSwapRequestsPaginatedRow, viewerEmployeeID uuid.UUID) ShiftSwapResponse {
	requesterName := strings.TrimSpace(row.RequesterFirstName + " " + row.RequesterLastName)
	recipientName := strings.TrimSpace(row.RecipientFirstName + " " + row.RecipientLastName)

	response := ShiftSwapResponse{
		ID:                    row.ID,
		RequesterEmployeeID:   row.RequesterEmployeeID,
		RequesterEmployeeName: requesterName,
		RecipientEmployeeID:   row.RecipientEmployeeID,
		RecipientEmployeeName: recipientName,
		RequesterSchedule: ShiftSwapScheduleSnapshot{
			ID:            row.RequesterScheduleID,
			EmployeeID:    row.RequesterEmployeeID,
			EmployeeName:  requesterName,
			StartDatetime: row.RequesterScheduleStartDatetime.Time,
			EndDatetime:   row.RequesterScheduleEndDatetime.Time,
		},
		RecipientSchedule: ShiftSwapScheduleSnapshot{
			ID:            row.RecipientScheduleID,
			EmployeeID:    row.RecipientEmployeeID,
			EmployeeName:  recipientName,
			StartDatetime: row.RecipientScheduleStartDatetime.Time,
			EndDatetime:   row.RecipientScheduleEndDatetime.Time,
		},
		Status:      string(row.Status),
		RequestedAt: row.RequestedAt.Time,
	}

	if row.RecipientRespondedAt.Valid {
		t := row.RecipientRespondedAt.Time
		response.RecipientRespondedAt = &t
	}
	if row.AdminDecidedAt.Valid {
		t := row.AdminDecidedAt.Time
		response.AdminDecidedAt = &t
	}
	response.RecipientResponseNote = row.RecipientResponseNote
	response.AdminDecisionNote = row.AdminDecisionNote
	response.AdminEmployeeID = row.AdminEmployeeID
	if row.AdminFirstName != nil && row.AdminLastName != nil {
		name := strings.TrimSpace(*row.AdminFirstName + " " + *row.AdminLastName)
		response.AdminEmployeeName = &name
	}
	if row.ExpiresAt.Valid {
		t := row.ExpiresAt.Time
		response.ExpiresAt = &t
	}
	if viewerEmployeeID == row.RequesterEmployeeID {
		response.Direction = "sent"
	} else if viewerEmployeeID == row.RecipientEmployeeID {
		response.Direction = "received"
	}
	return response
}

func parseShiftSwapStatus(raw string) (db.ShiftSwapStatusEnum, bool) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case string(db.ShiftSwapStatusEnumPendingRecipient):
		return db.ShiftSwapStatusEnumPendingRecipient, true
	case string(db.ShiftSwapStatusEnumRecipientRejected):
		return db.ShiftSwapStatusEnumRecipientRejected, true
	case string(db.ShiftSwapStatusEnumPendingAdmin):
		return db.ShiftSwapStatusEnumPendingAdmin, true
	case string(db.ShiftSwapStatusEnumAdminRejected):
		return db.ShiftSwapStatusEnumAdminRejected, true
	case string(db.ShiftSwapStatusEnumConfirmed):
		return db.ShiftSwapStatusEnumConfirmed, true
	case string(db.ShiftSwapStatusEnumCancelled):
		return db.ShiftSwapStatusEnumCancelled, true
	case string(db.ShiftSwapStatusEnumExpired):
		return db.ShiftSwapStatusEnumExpired, true
	default:
		return "", false
	}
}

func (s *scheduleService) notifyEmployees(ctx context.Context, employeeIDs []uuid.UUID, message string) {
	if len(employeeIDs) == 0 {
		return
	}

	validEmployeeIDs := make([]uuid.UUID, 0, len(employeeIDs))
	for _, employeeID := range employeeIDs {
		if employeeID != uuid.Nil {
			validEmployeeIDs = append(validEmployeeIDs, employeeID)
		}
	}
	if len(validEmployeeIDs) == 0 {
		return
	}

	resolvedUserIDs, err := s.Store.ListUserIDsByEmployeeIDs(ctx, validEmployeeIDs)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "notifyEmployees", "failed to resolve user ids from employee ids", zap.Error(err))
		return
	}

	userIDs := make([]uuid.UUID, 0, len(resolvedUserIDs))
	seen := map[uuid.UUID]struct{}{}
	for _, userID := range resolvedUserIDs {
		if userID == uuid.Nil {
			continue
		}
		if _, exists := seen[userID]; exists {
			continue
		}
		seen[userID] = struct{}{}
		userIDs = append(userIDs, userID)
	}

	if len(userIDs) == 0 {
		return
	}
	if err := s.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
		RecipientUserIDs: userIDs,
		Type:             notification.TypeSystemReminder,
		Data:             notification.NotificationData{},
		CreatedAt:        time.Now().UTC(),
		Message:          message,
	}); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "notifyEmployees", "failed to enqueue system reminder notification", zap.Error(err))
	}
}

func (s *scheduleService) notifyAdmins(ctx context.Context, message string) {
	adminUsers, err := s.Store.GetAllAdminUsers(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "notifyAdmins", "failed to load admin users", zap.Error(err))
		return
	}

	userIDs := make([]uuid.UUID, 0, len(adminUsers))
	for _, admin := range adminUsers {
		if admin.ID != uuid.Nil {
			userIDs = append(userIDs, admin.ID)
		}
	}
	if len(userIDs) == 0 {
		return
	}
	if err := s.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
		RecipientUserIDs: userIDs,
		Type:             notification.TypeSystemReminder,
		Data:             notification.NotificationData{},
		CreatedAt:        time.Now().UTC(),
		Message:          message,
	}); err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelWarn, "notifyAdmins", "failed to enqueue admin system reminder notification", zap.Error(err))
	}
}
