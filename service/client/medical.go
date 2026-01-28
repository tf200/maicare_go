package clientp

import (
	"context"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateClientDiagnosis(ctx context.Context, req CreateClientDiagnosisRequest, clientID uuid.UUID) (*CreateClientDiagnosisResponse, error) {
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
	var diagnosis db.ClientDiagnosis
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		diagnosis, err = q.CreateClientDiagnosis(ctx, arg)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientDiagnosis", "Failed to create client diagnosis", zap.Error(err), zap.String("client_id", clientID.String()))
			return err
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
				_, err := q.CreateClientMedication(ctx, medArg)
				if err != nil {
					s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientDiagnosis", "Failed to create diagnosis medication", zap.Error(err), zap.String("client_id", clientID.String()))
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
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

func (s *clientService) ListClientDiagnoses(ctx *gin.Context, req ListClientDiagnosesRequest, clientID uuid.UUID) (*pagination.Response[ListClientDiagnosesResponse], error) {
	params := req.GetParams()

	arg := db.ListClientDiagnosesParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}

	var diagnoses []db.ListClientDiagnosesRow
	var medications []db.ClientMedication
	var medicationsErr error
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		diagnoses, err = q.ListClientDiagnoses(ctx, arg)
		if err != nil {
			return err
		}
		if len(diagnoses) == 0 {
			return nil
		}
		diagnosisIDs := make([]uuid.UUID, 0, len(diagnoses))
		for _, d := range diagnoses {
			diagnosisIDs = append(diagnosisIDs, d.ID)
		}
		medications, err = q.ListMedicationsByDiagnosisIDs(ctx, diagnosisIDs)
		if err != nil {
			medicationsErr = err
		}
		return err
	})
	if err != nil {
		if medicationsErr != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientDiagnoses", "Failed to list medications by diagnosis IDs", zap.Error(err), zap.String("client_id", clientID.String()))
		} else {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientDiagnoses", "Failed to list client diagnoses", zap.Error(err), zap.String("client_id", clientID.String()))
		}
		return nil, err
	}

	if len(diagnoses) == 0 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListClientDiagnoses", "No diagnoses found for client", zap.String("client_id", clientID.String()))
		pag := pagination.NewResponse(ctx, req.Request, []ListClientDiagnosesResponse{}, 0)
		return &pag, nil
	}

	totalCount := diagnoses[0].TotalDiagnoses

	res := make([]ListClientDiagnosesResponse, len(diagnoses))
	diagnosisIDs := make([]uuid.UUID, 0, len(diagnoses))
	diagIndexMap := make(map[uuid.UUID]int, len(diagnoses))

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

func (s *clientService) GetClientDiagnosis(ctx context.Context, diagnosisID uuid.UUID) (*GetClientDiagnosisResponse, error) {
	var diagnosis db.ClientDiagnosis
	var medications []db.ListMedicationsByDiagnosisIDRow
	var medicationsErr error
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		diagnosis, err = q.GetClientDiagnosis(ctx, diagnosisID)
		if err != nil {
			return err
		}
		medications, err = q.ListMedicationsByDiagnosisID(ctx, db.ListMedicationsByDiagnosisIDParams{
			DiagnosisID: &diagnosisID,
			Limit:       100,
			Offset:      0,
		})
		if err != nil {
			medicationsErr = err
		}
		return err
	})
	if err != nil {
		if medicationsErr != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientDiagnosis", "Failed to list medications by diagnosis ID", zap.Error(err), zap.String("diagnosis_id", diagnosisID.String()))
		} else {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientDiagnosis", "Failed to get client diagnosis", zap.Error(err), zap.String("diagnosis_id", diagnosisID.String()))
		}
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

func (s *clientService) UpdateClientDiagnosis(ctx context.Context, req UpdateClientDiagnosisRequest, diagnosisID uuid.UUID) (*UpdateClientDiagnosisResponse, error) {
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

	var diagnosis db.ClientDiagnosis
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		diagnosis, err = q.UpdateClientDiagnosis(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientDiagnosis", "Failed to update client diagnosis", zap.Error(err), zap.String("diagnosis_id", diagnosisID.String()))
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

func (s *clientService) DeleteClientDiagnosis(ctx context.Context, diagnosisID uuid.UUID) (*DeleteClientDiagnosisResponse, error) {
	var diag db.ClientDiagnosis
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		diag, err = q.DeleteClientDiagnosis(ctx, diagnosisID)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteClientDiagnosis", "Failed to delete client diagnosis", zap.Error(err), zap.String("diagnosis_id", diagnosisID.String()))
		return nil, err
	}
	res := &DeleteClientDiagnosisResponse{
		ID: diag.ID,
	}
	return res, nil
}

func (s *clientService) CreateClientMedication(ctx context.Context, req CreateClientMedicationRequest, diagnosisID *uuid.UUID) (*CreateClientMedicationResponse, error) {
	arg := db.CreateClientMedicationParams{
		DiagnosisID:      diagnosisID,
		Name:             req.Name,
		Dosage:           req.Dosage,
		StartDate:        pgtype.Date{Time: req.StartDate, Valid: true},
		EndDate:          pgtype.Date{Time: req.EndDate, Valid: true},
		Notes:            req.Notes,
		SelfAdministered: req.SelfAdministered,
		AdministeredByID: req.AdministeredByID,
		IsCritical:       req.IsCritical,
	}
	var medication db.ClientMedication
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		medication, err = q.CreateClientMedication(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientMedication", "Failed to create client medication", zap.Error(err))
		return nil, err
	}
	res := &CreateClientMedicationResponse{
		ID:               medication.ID,
		DiagnosisID:      medication.DiagnosisID,
		Name:             medication.Name,
		Dosage:           medication.Dosage,
		StartDate:        medication.StartDate.Time,
		EndDate:          medication.EndDate.Time,
		Notes:            medication.Notes,
		SelfAdministered: medication.SelfAdministered,
		AdministeredByID: medication.AdministeredByID,
		IsCritical:       medication.IsCritical,
		UpdatedAt:        medication.UpdatedAt.Time,
		CreatedAt:        medication.CreatedAt.Time,
	}
	return res, nil
}

func (s *clientService) ListMedicationsByDiagnosisID(ctx *gin.Context, req ListClientMedicationsRequest, diagnosisID *uuid.UUID) (*pagination.Response[ListClientMedicationsResponse], error) {
	params := req.GetParams()
	arg := db.ListMedicationsByDiagnosisIDParams{
		DiagnosisID: diagnosisID,
		Limit:       params.Limit,
		Offset:      params.Offset,
	}
	var medications []db.ListMedicationsByDiagnosisIDRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		medications, err = q.ListMedicationsByDiagnosisID(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListMedicationsByDiagnosisID", "Failed to list medications by diagnosis ID", zap.Error(err), zap.String("diagnosis_id", diagnosisID.String()))
		return nil, err
	}
	totalCount := medications[0].TotalMedications
	var res []ListClientMedicationsResponse
	for _, m := range medications {
		res = append(res, ListClientMedicationsResponse{
			ID:               m.ID,
			Name:             m.Name,
			Dosage:           m.Dosage,
			StartDate:        m.StartDate.Time,
			EndDate:          m.EndDate.Time,
			Notes:            m.Notes,
			SelfAdministered: m.SelfAdministered,
			IsCritical:       m.IsCritical,
			DiagnosisID:      m.DiagnosisID,
			AdministeredByID: m.AdministeredByID,
			UpdatedAt:        m.UpdatedAt.Time,
			CreatedAt:        m.CreatedAt.Time,
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, res, totalCount)
	return &pag, nil
}

func (s *clientService) GetClientMedication(ctx context.Context, medicationID uuid.UUID) (*GetClientMedicationResponse, error) {
	var medication db.GetMedicationRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		medication, err = q.GetMedication(ctx, medicationID)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientMedication", "Failed to get client medication", zap.Error(err), zap.String("medication_id", medicationID.String()))
		return nil, err
	}
	res := &GetClientMedicationResponse{
		ID:                      medication.ID,
		Name:                    medication.Name,
		Dosage:                  medication.Dosage,
		StartDate:               medication.StartDate.Time,
		EndDate:                 medication.EndDate.Time,
		Notes:                   medication.Notes,
		SelfAdministered:        medication.SelfAdministered,
		DiagnosisID:             medication.DiagnosisID,
		AdministeredByID:        medication.AdministeredByID,
		IsCritical:              medication.IsCritical,
		UpdatedAt:               medication.UpdatedAt.Time,
		CreatedAt:               medication.CreatedAt.Time,
		AdministeredByFirstName: medication.AdministeredByFirstName,
		AdministeredByLastName:  medication.AdministeredByLastName,
	}
	return res, nil
}

func (s *clientService) UpdateClientMedication(ctx context.Context, req UpdateClientMedicationRequest, medicationID uuid.UUID) (*UpdateClientMedicationResponse, error) {
	arg := db.UpdateClientMedicationParams{
		ID:               medicationID,
		Name:             req.Name,
		Dosage:           req.Dosage,
		StartDate:        pgtype.Date{Time: req.StartDate, Valid: true},
		EndDate:          pgtype.Date{Time: req.EndDate, Valid: true},
		Notes:            req.Notes,
		SelfAdministered: req.SelfAdministered,
		AdministeredByID: req.AdministeredByID,
		IsCritical:       req.IsCritical,
	}
	var medication db.ClientMedication
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		medication, err = q.UpdateClientMedication(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientMedication", "Failed to update client medication", zap.Error(err), zap.String("medication_id", medicationID.String()))
		return nil, err
	}
	res := &UpdateClientMedicationResponse{
		ID:               medication.ID,
		Name:             medication.Name,
		Dosage:           medication.Dosage,
		StartDate:        medication.StartDate.Time,
		EndDate:          medication.EndDate.Time,
		Notes:            medication.Notes,
		SelfAdministered: medication.SelfAdministered,
		DiagnosisID:      medication.DiagnosisID,
		AdministeredByID: medication.AdministeredByID,
		IsCritical:       medication.IsCritical,
		UpdatedAt:        medication.UpdatedAt.Time,
		CreatedAt:        medication.CreatedAt.Time,
	}
	return res, nil
}

func (s *clientService) DeleteClientMedication(ctx context.Context, medicationID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteClientMedication(ctx, medicationID)
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteClientMedication", "Failed to delete client medication", zap.Error(err), zap.String("medication_id", medicationID.String()))
		return err
	}
	return nil
}
