package clientp

import (
	"context"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) PutClientOutOfCare(ctx context.Context, req PutClientOutOfCareRequest, clientID uuid.UUID) (*PutClientOutOfCareResponse, error) {
	dischargeDate, err := time.Parse("2006-01-02", req.DischargeDate)
	if err != nil {
		return nil, fmt.Errorf("discharge_date must be in YYYY-MM-DD format")
	}
	dischargeDay := time.Date(dischargeDate.Year(), dischargeDate.Month(), dischargeDate.Day(), 0, 0, 0, 0, time.UTC)

	var finalEvaluation *string
	if req.FinalEvaluation != nil {
		trimmed := strings.TrimSpace(*req.FinalEvaluation)
		if trimmed == "" {
			return nil, fmt.Errorf("final_evaluation cannot be empty")
		}
		finalEvaluation = &trimmed
	}

	dischargeReason, err := parseDischargeReason(req.DischargeReason)
	if err != nil {
		return nil, err
	}

	today := time.Now()
	todayDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	targetStatus := db.ClientStatusEnumOutOfCare
	if dischargeDay.After(todayDate) {
		targetStatus = db.ClientStatusEnumScheduledOutOfCare
	}
	if targetStatus == db.ClientStatusEnumOutOfCare && finalEvaluation == nil {
		return nil, fmt.Errorf("final_evaluation is required when discharge is effective today or in the past")
	}

	var updatedClient db.ClientDetail
	err = s.Store.ExecTx(ctx, func(q *db.Queries) error {
		existingClient, err := q.GetClientDetails(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to get client details: %w", err)
		}

		if existingClient.Status != db.ClientStatusEnumInCare && existingClient.Status != db.ClientStatusEnumScheduledOutOfCare {
			return fmt.Errorf("client must be in in_care or scheduled_out_of_care to be put out of care")
		}

		updatedClient, err = q.PutClientOutOfCare(ctx, db.PutClientOutOfCareParams{
			ID:            clientID,
			Status:        targetStatus,
			DischargeDate: pgtype.Date{Time: dischargeDay, Valid: true},
			DischargeReason: db.NullDischargeReasonEnum{
				DischargeReasonEnum: dischargeReason,
				Valid:               true,
			},
			FinalEvaluation: finalEvaluation,
		})
		if err != nil {
			return fmt.Errorf("failed to put client out of care: %w", err)
		}

		historyReason := fmt.Sprintf("client_put_out_of_care:%s", dischargeReason)
		if targetStatus == db.ClientStatusEnumScheduledOutOfCare {
			historyReason = fmt.Sprintf("client_scheduled_out_of_care:%s", dischargeReason)
		}
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
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "PutClientOutOfCare", "Failed to put client out of care", zap.String("client_id", clientID.String()), zap.Error(err))
		return nil, err
	}

	response := &PutClientOutOfCareResponse{
		ID:              updatedClient.ID,
		Status:          string(updatedClient.Status),
		DischargeDate:   updatedClient.DischargeDate.Time,
		DischargeReason: string(updatedClient.DischargeReason.DischargeReasonEnum),
		FinalEvaluation: updatedClient.FinalEvaluation,
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "PutClientOutOfCare", "Client discharge status updated", zap.String("client_id", clientID.String()), zap.String("status", response.Status))
	return response, nil
}

func parseDischargeReason(reason string) (db.DischargeReasonEnum, error) {
	normalized := strings.TrimSpace(reason)
	switch normalized {
	case string(db.DischargeReasonEnumTreatmentCompleted),
		string(db.DischargeReasonEnumTerminatedByMutualAgreement),
		string(db.DischargeReasonEnumTerminatedByClient),
		string(db.DischargeReasonEnumTerminatedByProvider),
		string(db.DischargeReasonEnumTerminatedDueToExternalFactors),
		string(db.DischargeReasonEnumOther):
		return db.DischargeReasonEnum(normalized), nil
	default:
		return "", fmt.Errorf("invalid discharge_reason")
	}
}
