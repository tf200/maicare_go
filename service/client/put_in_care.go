package clientp

import (
	"context"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/service/notification"
	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) PutClientInCare(ctx context.Context, req PutClientInCareRequest, clientID uuid.UUID) (*PutClientInCareResponse, error) {
	careStartDate, err := time.Parse("2006-01-02", req.CareStartDate)
	if err != nil {
		return nil, fmt.Errorf("care_start_date must be in YYYY-MM-DD format")
	}

	today := time.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	careStartDay := time.Date(careStartDate.Year(), careStartDate.Month(), careStartDate.Day(), 0, 0, 0, 0, time.UTC)

	targetStatus := db.ClientStatusEnumInCare
	if careStartDay.After(todayDate) {
		targetStatus = db.ClientStatusEnumScheduledInCare
	}

	var warning *string
	if careStartDay.Before(todayDate) {
		warningMsg := "care_start_date is in the past; client was activated immediately"
		warning = &warningMsg
	}

	var coordinatorAssignment db.UpsertMainCoordinatorRow
	var updatedClient db.ClientDetail

	err = s.Store.ExecTx(ctx, func(q *db.Queries) error {
		existingClient, err := q.GetClientDetails(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to get client details: %w", err)
		}

		if existingClient.Status != db.ClientStatusEnumOnWaitingList {
			return fmt.Errorf("client must be in on_waiting_list to be put in care")
		}

		activeGoalsCount, err := q.CountActiveGoalsByClientID(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to validate active goals: %w", err)
		}
		if activeGoalsCount == 0 {
			return fmt.Errorf("client cannot be put in care without at least one active goal")
		}

		coordinatorAssignment, err = q.UpsertMainCoordinator(ctx, db.UpsertMainCoordinatorParams{
			ClientID:   clientID,
			EmployeeID: req.CoordinatorEmployeeID,
			StartDate:  pgtype.Date{Time: careStartDay, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to assign main coordinator: %w", err)
		}

		placedInCareAt := pgtype.Timestamptz{Valid: false}
		if req.PlacedInCareAt != nil {
			placedInCareAt = pgtype.Timestamptz{Time: *req.PlacedInCareAt, Valid: true}
		}

		updatedClient, err = q.PutClientInCare(ctx, db.PutClientInCareParams{
			ID:             clientID,
			Status:         targetStatus,
			CareStartDate:  pgtype.Date{Time: careStartDay, Valid: true},
			PlacedInCareAt: placedInCareAt,
		})
		if err != nil {
			return fmt.Errorf("failed to put client in care: %w", err)
		}

		historyReason := "client_put_in_care"
		if req.Reason != nil {
			trimmed := strings.TrimSpace(*req.Reason)
			if trimmed != "" {
				historyReason = trimmed
			}
		}

		_, err = q.CreateClientStatusHistory(ctx, db.CreateClientStatusHistoryParams{
			ClientID:  clientID,
			OldStatus: util.StringPtr(string(existingClient.Status)),
			NewStatus: string(updatedClient.Status),
			Reason:    &historyReason,
		})
		if err != nil {
			return fmt.Errorf("failed to create client status history: %w", err)
		}

		return nil
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PutClientInCare", "Failed to put client in care", zap.String("client_id", clientID.String()), zap.Error(err))
		return nil, err
	}

	notificationData := notification.NewClientAssignmentData{
		ClientID:        coordinatorAssignment.ClientID,
		ClientFirstName: coordinatorAssignment.ClientFirstName,
		ClientLastName:  coordinatorAssignment.ClientLastName,
		ClientLocation:  coordinatorAssignment.ClientLocationName,
	}
	err = s.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
		RecipientUserIDs: []uuid.UUID{coordinatorAssignment.UserID},
		Type:             notification.TypeNewClientAssignment,
		Data: notification.NotificationData{
			NewClientAssignment: &notificationData,
		},
		CreatedAt: time.Now(),
		Message:   notificationData.NewClientAssignmentMessage(),
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PutClientInCare", "Failed to enqueue coordinator notification", zap.String("client_id", clientID.String()), zap.Error(err))
	}

	response := &PutClientInCareResponse{
		ID:                  updatedClient.ID,
		Status:              string(updatedClient.Status),
		CareStartDate:       updatedClient.CareStartDate.Time,
		PlacedInCareAt:      updatedClient.PlacedInCareAt.Time,
		CoordinatorAssignID: coordinatorAssignment.ID,
		Warning:             warning,
	}
	if updatedClient.NextEvaluationDate.Valid {
		nextDate := updatedClient.NextEvaluationDate.Time
		response.NextEvaluationDate = &nextDate
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "PutClientInCare", "Client moved into care lifecycle", zap.String("client_id", clientID.String()), zap.String("status", response.Status))
	return response, nil
}
