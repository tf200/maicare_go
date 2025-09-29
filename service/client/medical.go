package clientp

import (
	"context"
	"database/sql"
	"errors"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
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

func (s *clientService) ListClientDiagnoses(ctx *gin.Context, req ListClientDiagnosesRequest, clientID int64) (*pagination.Response[ListClientDiagnosesResponse], error) {
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

	if len(diagnoses) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListClientDiagnoses", "No diagnoses found for client", zap.Int64("client_id", clientID))
		pag := pagination.NewResponse(ctx, req.Request, []ListClientDiagnosesResponse{}, 0)
		return &pag, nil
	}

	totalCount := diagnoses[0].TotalDiagnoses

	res := make([]ListClientDiagnosesResponse, len(diagnoses))
	diagnosisIDs := make([]int64, 0, len(diagnoses))
	diagIndexMap := make(map[int64]int, len(diagnoses))

	for i, d := range diagnoses {
		diagnosisIDs = append(diagnosisIDs, d.ID)
		diagIndexMap[d.ID] = i
		res[i] = ListClientDiagnosesResponse{
			ID:                  d.ID,
			Title:               d.Title,
			ClientID:            d.ClientID,
			DiagnosisCode:       d.DiagnosisCode,
			Description:         d.Description,
			Severity:            d.Severity,
			Status:              d.Status,
			DiagnosingClinician: d.DiagnosingClinician,
			Notes:               d.Notes,
			CreatedAt:           d.CreatedAt.Time,
			Medications:         []DiagnosisMedicationList{},
		}
	}

	// Fetch all related medications in a single database query.
	medications, err := s.Store.ListMedicationsByDiagnosisIDs(ctx, diagnosisIDs)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListClientDiagnoses", "Failed to list medications by diagnosis IDs", zap.Error(err), zap.Int64("client_id", clientID))
		return nil, err
	}

	for _, m := range medications {
		med := DiagnosisMedicationList{
			ID:               m.ID,
			Name:             m.Name,
			Dosage:           m.Dosage,
			StartDate:        m.StartDate.Time,
			EndDate:          m.EndDate.Time,
			Notes:            m.Notes,
			SelfAdministered: m.SelfAdministered,
			AdministeredByID: m.AdministeredByID,
			IsCritical:       m.IsCritical,
			CreatedAt:        m.CreatedAt.Time,
		}
		if index, ok := diagIndexMap[*m.DiagnosisID]; ok {
			res[index].Medications = append(res[index].Medications, med)
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, res, totalCount)
	return &pag, nil
}

func (s *clientService) GetClientDiagnosis(ctx context.Context, diagnosisID int64) (*GetClientDiagnosisResponse, error) {
	diagnosis, err := s.Store.GetClientDiagnosis(ctx, diagnosisID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientDiagnosis", "Failed to get client diagnosis", zap.Error(err), zap.Int64("diagnosis_id", diagnosisID))
		return nil, err
	}

	medications, err := s.Store.ListMedicationsByDiagnosisID(ctx, &diagnosisID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetClientDiagnosis", "Failed to list medications by diagnosis ID", zap.Error(err), zap.Int64("diagnosis_id", diagnosisID))
		return nil, err
	}

	var meds []DiagnosisMedicationList
	for _, m := range medications {
		meds = append(meds, DiagnosisMedicationList{
			ID:               m.ID,
			DiagnosisID:      m.DiagnosisID,
			Name:             m.Name,
			Dosage:           m.Dosage,
			StartDate:        m.StartDate.Time,
			EndDate:          m.EndDate.Time,
			Notes:            m.Notes,
			SelfAdministered: m.SelfAdministered,
			AdministeredByID: m.AdministeredByID,
			IsCritical:       m.IsCritical,
			UpdatedAt:        m.UpdatedAt.Time,
			CreatedAt:        m.CreatedAt.Time,
		})
	}

	res := &GetClientDiagnosisResponse{
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
		Medications:         meds,
	}
	return res, nil
}

func (s *clientService) UpdateClientDiagnosis(ctx context.Context, req UpdateClientDiagnosisRequest, diagnosisID int64) (*UpdateClientDiagnosisResponse, error) {
	arg := db.UpdateClientDiagnosisParams{
		ID:                  diagnosisID,
		Title:               req.Title,
		DiagnosisCode:       req.DiagnosisCode,
		Description:         req.Description,
		Severity:            req.Severity,
		Status:              req.Status,
		DiagnosingClinician: req.DiagnosingClinician,
		Notes:               req.Notes,
	}

	diagnosis, err := s.Store.UpdateClientDiagnosis(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateClientDiagnosis", "Failed to update client diagnosis", zap.Error(err), zap.Int64("diagnosis_id", diagnosisID))
		return nil, err
	}

	res := &UpdateClientDiagnosisResponse{
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

func (s *clientService) DeleteClientDiagnosis(ctx context.Context, diagnosisID int64) (*DeleteClientDiagnosisResponse, error) {
	diag, err := s.Store.DeleteClientDiagnosis(ctx, diagnosisID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteClientDiagnosis", "Failed to delete client diagnosis", zap.Error(err), zap.Int64("diagnosis_id", diagnosisID))
		return nil, err
	}
	res := &DeleteClientDiagnosisResponse{
		ID: diag.ID,
	}
	return res, nil
}


