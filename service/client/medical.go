package clientp

import (
	"context"
	"database/sql"
	"errors"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateClientDiagnosis(ctx context.Context, req CreateClientDiagnosisRequest, clientID int64) (*CreateClientDiagnosisResponse, error) {
	arg := db.CreateClientDiagnosisParams{
		ClientID:            clientID,
		Title:               req.Title,
		DiagnosisCode:       req.DiagnosisCode,
		Description:         req.Description,
		Severity:            req.Severity,
		Status:              req.Status,
		DiagnosingClinician: req.DiagnosingClinician,
		Notes:               req.Notes,
	}
	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDiagnosis", "Failed to begin transaction", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, sql.ErrTxDone) {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDiagnosis", "Failed to rollback transaction", zap.Error(err))
		}
	}()
	qtx := s.Store.WithTx(tx)
	diagnosis, err := qtx.CreateClientDiagnosis(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDiagnosis", "Failed to create client diagnosis", zap.Error(err), zap.Int64("client_id", clientID))
		return nil, err
	}

	if len(req.Medications) > 0 {
		for _, med := range req.Medications {
			medArg := db.CreateClientMedicationParams{
				DiagnosisID:      &diagnosis.ID,
				Name:             med.Name,
				Dosage:           med.Dosage,
				StartDate:        pgtype.Date{Time: med.StartDate, Valid: true},
				EndDate:          pgtype.Date{Time: med.EndDate, Valid: true},
				Notes:            med.Notes,
				SelfAdministered: med.SelfAdministered,
				AdministeredByID: med.AdministeredByID,
				IsCritical:       med.IsCritical,
			}
			_, err := qtx.CreateClientMedication(ctx, medArg)
			if err != nil {
				s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDiagnosis", "Failed to create diagnosis medication", zap.Error(err), zap.Int64("client_id", clientID))
				return nil, err
			}
		}
	}
	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateClientDiagnosis", "Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	res := &CreateClientDiagnosisResponse{
		ID:                  diagnosis.ID,
		Title:               diagnosis.Title,
		ClientID:            diagnosis.ClientID,
		DiagnosisCode:       diagnosis.DiagnosisCode,
		Description:         diagnosis.Description,
		Severity:            diagnosis.Severity,
		Status:              diagnosis.Status,
		DiagnosingClinician: diagnosis.DiagnosingClinician,
		Notes:               diagnosis.Notes,
		CreatedAt:           diagnosis.CreatedAt.Time,
	}
	return res, nil
}


func (s *clientService) ListClientDiagnoses(ctx context.Context, req ListClientDiagnosesRequest, clientID int64) ([]DiagnosisMedicationList, error) {
	params := req.GetParams()

	arg := db.ListClientDiagnosesParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}
	diagnoses, err := s.Store.ListClientDiagnoses(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListClientDiagnoses", "Failed to list client diagnoses", zap.Error(err), zap.Int64("client_id", clientID))
		return nil, err
	}