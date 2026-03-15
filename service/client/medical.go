package clientp

import (
	"context"
	"github.com/goccy/go-json"
	"errors"

	"maicare_go/infra"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"

	db "maicare_go/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// =====================
// Diagnoses
// =====================

func (s *clientService) CreateClientDiagnosis(ctx context.Context, req CreateClientDiagnosisRequest, clientID uuid.UUID) (*ClientDiagnosisResponse, error) {
	employeeID := infra.GetEmployeeID(ctx)
	var createdBy *uuid.UUID
	if employeeID != uuid.Nil {
		createdBy = &employeeID
	}

	status := db.DiagnosisStatusEnum("confirmed")
	if req.Status != nil {
		status = db.DiagnosisStatusEnum(*req.Status)
	}

	severity := db.DiagnosisSeverityEnum("unknown")
	if req.Severity != nil {
		severity = db.DiagnosisSeverityEnum(*req.Severity)
	}

	arg := db.CreateClientDiagnosisParams{
		ClientID:            clientID,
		CodeSystem:          req.CodeSystem,
		Code:                req.Code,
		Title:               req.Title,
		Description:         req.Description,
		Status:              status,
		Severity:            severity,
		DiagnosedOn:         util.NullableDate(req.DiagnosedOn),
		ResolvedOn:          util.NullableDate(req.ResolvedOn),
		DiagnosingClinician: req.DiagnosingClinician,
		Notes:               req.Notes,
		CreatedByEmployeeID: createdBy,
		UpdatedByEmployeeID: createdBy,
	}

	var diagnosis db.ClientDiagnosis
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		diagnosis, err = q.CreateClientDiagnosis(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientDiagnosis", "Failed to create client diagnosis", zap.Error(err), zap.String("client_id", clientID.String()))
		return nil, err
	}

	res := &ClientDiagnosisResponse{
		ID:                  diagnosis.ID,
		ClientID:            diagnosis.ClientID,
		CodeSystem:          diagnosis.CodeSystem,
		Code:                diagnosis.Code,
		Title:               diagnosis.Title,
		Description:         diagnosis.Description,
		Status:              string(diagnosis.Status),
		Severity:            string(diagnosis.Severity),
		DiagnosedOn:         util.DatePtr(diagnosis.DiagnosedOn),
		ResolvedOn:          util.DatePtr(diagnosis.ResolvedOn),
		DiagnosingClinician: diagnosis.DiagnosingClinician,
		Notes:               diagnosis.Notes,
		CreatedAt:           diagnosis.CreatedAt.Time,
		UpdatedAt:           diagnosis.UpdatedAt.Time,
	}
	return res, nil
}

func (s *clientService) ListClientDiagnoses(ctx *gin.Context, req ListClientDiagnosesRequest, clientID uuid.UUID) (*pagination.Response[ClientDiagnosisResponse], error) {
	params := req.GetParams()
	arg := db.ListClientDiagnosesParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	}

	var rows []db.ListClientDiagnosesRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		rows, err = q.ListClientDiagnoses(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientDiagnoses", "Failed to list client diagnoses", zap.Error(err), zap.String("client_id", clientID.String()))
		return nil, err
	}

	if len(rows) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ClientDiagnosisResponse{}, 0)
		return &pag, nil
	}

	totalCount := rows[0].TotalCount
	res := make([]ClientDiagnosisResponse, 0, len(rows))
	for _, d := range rows {
		res = append(res, ClientDiagnosisResponse{
			ID:                  d.ID,
			ClientID:            d.ClientID,
			CodeSystem:          d.CodeSystem,
			Code:                d.Code,
			Title:               d.Title,
			Description:         d.Description,
			Status:              string(d.Status),
			Severity:            string(d.Severity),
			DiagnosedOn:         util.DatePtr(d.DiagnosedOn),
			ResolvedOn:          util.DatePtr(d.ResolvedOn),
			DiagnosingClinician: d.DiagnosingClinician,
			Notes:               d.Notes,
			CreatedAt:           d.CreatedAt.Time,
			UpdatedAt:           d.UpdatedAt.Time,
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, res, totalCount)
	return &pag, nil
}

func (s *clientService) GetClientDiagnosis(ctx context.Context, clientID uuid.UUID, diagnosisID uuid.UUID) (*ClientDiagnosisResponse, error) {
	arg := db.GetClientDiagnosisParams{ClientID: clientID, ID: diagnosisID}
	var diagnosis db.ClientDiagnosis
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		diagnosis, err = q.GetClientDiagnosis(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientDiagnosis", "Failed to get client diagnosis", zap.Error(err), zap.String("client_id", clientID.String()), zap.String("diagnosis_id", diagnosisID.String()))
		return nil, err
	}

	res := &ClientDiagnosisResponse{
		ID:                  diagnosis.ID,
		ClientID:            diagnosis.ClientID,
		CodeSystem:          diagnosis.CodeSystem,
		Code:                diagnosis.Code,
		Title:               diagnosis.Title,
		Description:         diagnosis.Description,
		Status:              string(diagnosis.Status),
		Severity:            string(diagnosis.Severity),
		DiagnosedOn:         util.DatePtr(diagnosis.DiagnosedOn),
		ResolvedOn:          util.DatePtr(diagnosis.ResolvedOn),
		DiagnosingClinician: diagnosis.DiagnosingClinician,
		Notes:               diagnosis.Notes,
		CreatedAt:           diagnosis.CreatedAt.Time,
		UpdatedAt:           diagnosis.UpdatedAt.Time,
	}
	return res, nil
}

func (s *clientService) UpdateClientDiagnosis(ctx context.Context, req UpdateClientDiagnosisRequest, clientID uuid.UUID, diagnosisID uuid.UUID) (*ClientDiagnosisResponse, error) {
	employeeID := infra.GetEmployeeID(ctx)
	var updatedBy *uuid.UUID
	if employeeID != uuid.Nil {
		updatedBy = &employeeID
	}

	status := db.NullDiagnosisStatusEnum{Valid: false}
	if req.Status != nil {
		status = db.NullDiagnosisStatusEnum{DiagnosisStatusEnum: db.DiagnosisStatusEnum(*req.Status), Valid: true}
	}
	severity := db.NullDiagnosisSeverityEnum{Valid: false}
	if req.Severity != nil {
		severity = db.NullDiagnosisSeverityEnum{DiagnosisSeverityEnum: db.DiagnosisSeverityEnum(*req.Severity), Valid: true}
	}

	arg := db.UpdateClientDiagnosisParams{
		CodeSystem:          req.CodeSystem,
		Code:                req.Code,
		Title:               req.Title,
		Description:         req.Description,
		Status:              status,
		Severity:            severity,
		DiagnosedOn:         util.NullableDate(req.DiagnosedOn),
		ResolvedOn:          util.NullableDate(req.ResolvedOn),
		DiagnosingClinician: req.DiagnosingClinician,
		Notes:               req.Notes,
		UpdatedByEmployeeID: updatedBy,
		ClientID:            clientID,
		ID:                  diagnosisID,
	}

	var diagnosis db.ClientDiagnosis
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		diagnosis, err = q.UpdateClientDiagnosis(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientDiagnosis", "Failed to update client diagnosis", zap.Error(err), zap.String("client_id", clientID.String()), zap.String("diagnosis_id", diagnosisID.String()))
		return nil, err
	}

	res := &ClientDiagnosisResponse{
		ID:                  diagnosis.ID,
		ClientID:            diagnosis.ClientID,
		CodeSystem:          diagnosis.CodeSystem,
		Code:                diagnosis.Code,
		Title:               diagnosis.Title,
		Description:         diagnosis.Description,
		Status:              string(diagnosis.Status),
		Severity:            string(diagnosis.Severity),
		DiagnosedOn:         util.DatePtr(diagnosis.DiagnosedOn),
		ResolvedOn:          util.DatePtr(diagnosis.ResolvedOn),
		DiagnosingClinician: diagnosis.DiagnosingClinician,
		Notes:               diagnosis.Notes,
		CreatedAt:           diagnosis.CreatedAt.Time,
		UpdatedAt:           diagnosis.UpdatedAt.Time,
	}
	return res, nil
}

func (s *clientService) DeleteClientDiagnosis(ctx context.Context, clientID uuid.UUID, diagnosisID uuid.UUID) (*DeleteClientDiagnosisResponse, error) {
	var diagnosis db.ClientDiagnosis
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		diagnosis, err = q.DeleteClientDiagnosis(ctx, db.DeleteClientDiagnosisParams{ClientID: clientID, ID: diagnosisID})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteClientDiagnosis", "Failed to delete client diagnosis", zap.Error(err), zap.String("client_id", clientID.String()), zap.String("diagnosis_id", diagnosisID.String()))
		return nil, err
	}
	return &DeleteClientDiagnosisResponse{ID: diagnosis.ID}, nil
}

// =====================
// Medication Orders
// =====================

func (s *clientService) CreateClientMedicationOrder(ctx context.Context, req CreateClientMedicationOrderRequest, clientID uuid.UUID) (*ClientMedicationOrderResponse, error) {
	employeeID := infra.GetEmployeeID(ctx)
	var createdBy *uuid.UUID
	if employeeID != uuid.Nil {
		createdBy = &employeeID
	}

	status := db.MedicationOrderStatusEnum("active")
	if req.Status != nil {
		status = db.MedicationOrderStatusEnum(*req.Status)
	}
	adminMode := db.MedicationAdminModeEnum("self")
	if req.AdminMode != nil {
		adminMode = db.MedicationAdminModeEnum(*req.AdminMode)
	}

	arg := db.CreateClientMedicationOrderParams{
		ClientID:              clientID,
		DiagnosisID:           req.DiagnosisID,
		MedicationName:        req.MedicationName,
		DosageText:            req.DosageText,
		DoseAmount:            req.DoseAmount,
		DoseUnit:              req.DoseUnit,
		Route:                 req.Route,
		FrequencyText:         req.FrequencyText,
		Schedule:              util.EnsureJSONArray(req.Schedule),
		IsPrn:                 req.IsPrn,
		PrnIndication:         req.PrnIndication,
		MaxDosesPer24h:        req.MaxDosesPer24h,
		StartDate:             pgtype.Date{Time: req.StartDate, Valid: true},
		EndDate:               util.NullableDate(req.EndDate),
		Status:                status,
		AdminMode:             adminMode,
		ResponsibleEmployeeID: req.ResponsibleEmployeeID,
		IsCritical:            req.IsCritical,
		Notes:                 req.Notes,
		SourceAttachmentUuid:  req.SourceAttachmentUUID,
		CreatedByEmployeeID:   createdBy,
		UpdatedByEmployeeID:   createdBy,
	}

	var order db.ClientMedicationOrder
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		order, err = q.CreateClientMedicationOrder(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientMedicationOrder", "Failed to create medication order", zap.Error(err), zap.String("client_id", clientID.String()))
		return nil, err
	}

	res := &ClientMedicationOrderResponse{
		ID:                    order.ID,
		ClientID:              order.ClientID,
		DiagnosisID:           order.DiagnosisID,
		MedicationName:        order.MedicationName,
		DosageText:            order.DosageText,
		DoseAmount:            order.DoseAmount,
		DoseUnit:              order.DoseUnit,
		Route:                 order.Route,
		FrequencyText:         order.FrequencyText,
		Schedule:              json.RawMessage(order.Schedule),
		IsPrn:                 order.IsPrn,
		PrnIndication:         order.PrnIndication,
		MaxDosesPer24h:        order.MaxDosesPer24h,
		StartDate:             order.StartDate.Time,
		EndDate:               util.DatePtr(order.EndDate),
		Status:                string(order.Status),
		AdminMode:             string(order.AdminMode),
		ResponsibleEmployeeID: order.ResponsibleEmployeeID,
		IsCritical:            order.IsCritical,
		Notes:                 order.Notes,
		SourceAttachmentUUID:  order.SourceAttachmentUuid,
		CreatedAt:             order.CreatedAt.Time,
		UpdatedAt:             order.UpdatedAt.Time,
	}
	return res, nil
}

func (s *clientService) ListClientMedicationOrders(ctx *gin.Context, req ListClientMedicationOrdersRequest, clientID uuid.UUID) (*pagination.Response[ClientMedicationOrderResponse], error) {
	params := req.GetParams()

	status := db.NullMedicationOrderStatusEnum{Valid: false}
	if req.Status != nil {
		status = db.NullMedicationOrderStatusEnum{MedicationOrderStatusEnum: db.MedicationOrderStatusEnum(*req.Status), Valid: true}
	}
	adminMode := db.NullMedicationAdminModeEnum{Valid: false}
	if req.AdminMode != nil {
		adminMode = db.NullMedicationAdminModeEnum{MedicationAdminModeEnum: db.MedicationAdminModeEnum(*req.AdminMode), Valid: true}
	}

	arg := db.ListClientMedicationOrdersParams{
		ClientID:    clientID,
		Status:      status,
		AdminMode:   adminMode,
		DiagnosisID: req.DiagnosisID,
		Search:      req.Search,
		Offset:      params.Offset,
		Limit:       params.Limit,
	}

	var rows []db.ListClientMedicationOrdersRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		rows, err = q.ListClientMedicationOrders(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListClientMedicationOrders", "Failed to list medication orders", zap.Error(err), zap.String("client_id", clientID.String()))
		return nil, err
	}

	if len(rows) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ClientMedicationOrderResponse{}, 0)
		return &pag, nil
	}

	totalCount := rows[0].TotalCount
	res := make([]ClientMedicationOrderResponse, 0, len(rows))
	for _, o := range rows {
		res = append(res, ClientMedicationOrderResponse{
			ID:                           o.ID,
			ClientID:                     o.ClientID,
			DiagnosisID:                  o.DiagnosisID,
			MedicationName:               o.MedicationName,
			DosageText:                   o.DosageText,
			DoseAmount:                   o.DoseAmount,
			DoseUnit:                     o.DoseUnit,
			Route:                        o.Route,
			FrequencyText:                o.FrequencyText,
			Schedule:                     json.RawMessage(o.Schedule),
			IsPrn:                        o.IsPrn,
			PrnIndication:                o.PrnIndication,
			MaxDosesPer24h:               o.MaxDosesPer24h,
			StartDate:                    o.StartDate.Time,
			EndDate:                      util.DatePtr(o.EndDate),
			Status:                       string(o.Status),
			AdminMode:                    string(o.AdminMode),
			ResponsibleEmployeeID:        o.ResponsibleEmployeeID,
			ResponsibleEmployeeFirstName: o.ResponsibleEmployeeFirstName,
			ResponsibleEmployeeLastName:  o.ResponsibleEmployeeLastName,
			IsCritical:                   o.IsCritical,
			Notes:                        o.Notes,
			SourceAttachmentUUID:         o.SourceAttachmentUuid,
			DiagnosisTitle:               o.DiagnosisTitle,
			DiagnosisCodeSystem:          o.DiagnosisCodeSystem,
			DiagnosisCode:                o.DiagnosisCode,
			CreatedAt:                    o.CreatedAt.Time,
			UpdatedAt:                    o.UpdatedAt.Time,
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, res, totalCount)
	return &pag, nil
}

func (s *clientService) GetClientMedicationOrder(ctx context.Context, clientID uuid.UUID, orderID uuid.UUID) (*ClientMedicationOrderResponse, error) {
	var row db.GetClientMedicationOrderRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		row, err = q.GetClientMedicationOrder(ctx, db.GetClientMedicationOrderParams{ClientID: clientID, ID: orderID})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientMedicationOrder", "Failed to get medication order", zap.Error(err), zap.String("client_id", clientID.String()), zap.String("order_id", orderID.String()))
		return nil, err
	}

	res := &ClientMedicationOrderResponse{
		ID:                           row.ID,
		ClientID:                     row.ClientID,
		DiagnosisID:                  row.DiagnosisID,
		MedicationName:               row.MedicationName,
		DosageText:                   row.DosageText,
		DoseAmount:                   row.DoseAmount,
		DoseUnit:                     row.DoseUnit,
		Route:                        row.Route,
		FrequencyText:                row.FrequencyText,
		Schedule:                     json.RawMessage(row.Schedule),
		IsPrn:                        row.IsPrn,
		PrnIndication:                row.PrnIndication,
		MaxDosesPer24h:               row.MaxDosesPer24h,
		StartDate:                    row.StartDate.Time,
		EndDate:                      util.DatePtr(row.EndDate),
		Status:                       string(row.Status),
		AdminMode:                    string(row.AdminMode),
		ResponsibleEmployeeID:        row.ResponsibleEmployeeID,
		ResponsibleEmployeeFirstName: row.ResponsibleEmployeeFirstName,
		ResponsibleEmployeeLastName:  row.ResponsibleEmployeeLastName,
		IsCritical:                   row.IsCritical,
		Notes:                        row.Notes,
		SourceAttachmentUUID:         row.SourceAttachmentUuid,
		DiagnosisTitle:               row.DiagnosisTitle,
		DiagnosisCodeSystem:          row.DiagnosisCodeSystem,
		DiagnosisCode:                row.DiagnosisCode,
		CreatedAt:                    row.CreatedAt.Time,
		UpdatedAt:                    row.UpdatedAt.Time,
	}
	return res, nil
}

func (s *clientService) UpdateClientMedicationOrder(ctx context.Context, req UpdateClientMedicationOrderRequest, clientID uuid.UUID, orderID uuid.UUID) (*ClientMedicationOrderResponse, error) {
	employeeID := infra.GetEmployeeID(ctx)
	var updatedBy *uuid.UUID
	if employeeID != uuid.Nil {
		updatedBy = &employeeID
	}

	status := db.NullMedicationOrderStatusEnum{Valid: false}
	if req.Status != nil {
		status = db.NullMedicationOrderStatusEnum{MedicationOrderStatusEnum: db.MedicationOrderStatusEnum(*req.Status), Valid: true}
	}
	adminMode := db.NullMedicationAdminModeEnum{Valid: false}
	if req.AdminMode != nil {
		adminMode = db.NullMedicationAdminModeEnum{MedicationAdminModeEnum: db.MedicationAdminModeEnum(*req.AdminMode), Valid: true}
	}

	startDate := pgtype.Date{Valid: false}
	if req.StartDate != nil {
		startDate = pgtype.Date{Time: *req.StartDate, Valid: true}
	}
	endDate := pgtype.Date{Valid: false}
	if req.EndDate != nil {
		endDate = pgtype.Date{Time: *req.EndDate, Valid: true}
	}

	updateArg := db.UpdateClientMedicationOrderParams{
		DiagnosisID:           req.DiagnosisID,
		MedicationName:        req.MedicationName,
		DosageText:            req.DosageText,
		DoseAmount:            req.DoseAmount,
		DoseUnit:              req.DoseUnit,
		Route:                 req.Route,
		FrequencyText:         req.FrequencyText,
		Schedule:              req.Schedule,
		IsPrn:                 req.IsPrn,
		PrnIndication:         req.PrnIndication,
		MaxDosesPer24h:        req.MaxDosesPer24h,
		StartDate:             startDate,
		EndDate:               endDate,
		Status:                status,
		AdminMode:             adminMode,
		ResponsibleEmployeeID: req.ResponsibleEmployeeID,
		IsCritical:            req.IsCritical,
		Notes:                 req.Notes,
		SourceAttachmentUuid:  req.SourceAttachmentUUID,
		UpdatedByEmployeeID:   updatedBy,
		ClientID:              clientID,
		ID:                    orderID,
	}

	var order db.ClientMedicationOrder
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		order, err = q.UpdateClientMedicationOrder(ctx, updateArg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientMedicationOrder", "Failed to update medication order", zap.Error(err), zap.String("client_id", clientID.String()), zap.String("order_id", orderID.String()))
		return nil, err
	}

	// To keep response consistent (includes joins), re-fetch with GetClientMedicationOrder.
	row, err := s.GetClientMedicationOrder(ctx, clientID, order.ID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientMedicationOrder", "Updated medication order but failed to fetch enriched view", zap.Error(err), zap.String("client_id", clientID.String()), zap.String("order_id", orderID.String()))
		return nil, err
	}
	return row, nil
}

func (s *clientService) DeleteClientMedicationOrder(ctx context.Context, clientID uuid.UUID, orderID uuid.UUID) (*DeleteClientMedicationOrderResponse, error) {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteClientMedicationOrder(ctx, db.DeleteClientMedicationOrderParams{ClientID: clientID, ID: orderID})
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteClientMedicationOrder", "Failed to delete medication order", zap.Error(err), zap.String("client_id", clientID.String()), zap.String("order_id", orderID.String()))
		return nil, err
	}
	return &DeleteClientMedicationOrderResponse{ID: orderID}, nil
}

// =====================
// Overview
// =====================

func (s *clientService) GetClientMedicalOverview(ctx context.Context, clientID uuid.UUID) (*ClientMedicalOverviewResponse, error) {
	// Fetch diagnoses (first page large enough for now) + active medication orders.
	var diagnoses []db.ListClientDiagnosesRow
	var orders []db.ListClientMedicationOrdersRow

	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		diagnoses, err = q.ListClientDiagnoses(ctx, db.ListClientDiagnosesParams{ClientID: clientID, Limit: 500, Offset: 0})
		if err != nil {
			return err
		}

		orders, err = q.ListClientMedicationOrders(ctx, db.ListClientMedicationOrdersParams{
			ClientID:    clientID,
			Status:      db.NullMedicationOrderStatusEnum{MedicationOrderStatusEnum: db.MedicationOrderStatusEnum("active"), Valid: true},
			AdminMode:   db.NullMedicationAdminModeEnum{Valid: false},
			DiagnosisID: nil,
			Search:      nil,
			Offset:      0,
			Limit:       500,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetClientMedicalOverview", "Failed to get client medical overview", zap.Error(err), zap.String("client_id", clientID.String()))
		return nil, err
	}

	diagnosisRes := make([]ClientDiagnosisResponse, 0, len(diagnoses))
	for _, d := range diagnoses {
		diagnosisRes = append(diagnosisRes, ClientDiagnosisResponse{
			ID:                  d.ID,
			ClientID:            d.ClientID,
			CodeSystem:          d.CodeSystem,
			Code:                d.Code,
			Title:               d.Title,
			Description:         d.Description,
			Status:              string(d.Status),
			Severity:            string(d.Severity),
			DiagnosedOn:         util.DatePtr(d.DiagnosedOn),
			ResolvedOn:          util.DatePtr(d.ResolvedOn),
			DiagnosingClinician: d.DiagnosingClinician,
			Notes:               d.Notes,
			CreatedAt:           d.CreatedAt.Time,
			UpdatedAt:           d.UpdatedAt.Time,
		})
	}

	orderRes := make([]ClientMedicationOrderResponse, 0, len(orders))
	for _, o := range orders {
		orderRes = append(orderRes, ClientMedicationOrderResponse{
			ID:                           o.ID,
			ClientID:                     o.ClientID,
			DiagnosisID:                  o.DiagnosisID,
			MedicationName:               o.MedicationName,
			DosageText:                   o.DosageText,
			DoseAmount:                   o.DoseAmount,
			DoseUnit:                     o.DoseUnit,
			Route:                        o.Route,
			FrequencyText:                o.FrequencyText,
			Schedule:                     json.RawMessage(o.Schedule),
			IsPrn:                        o.IsPrn,
			PrnIndication:                o.PrnIndication,
			MaxDosesPer24h:               o.MaxDosesPer24h,
			StartDate:                    o.StartDate.Time,
			EndDate:                      util.DatePtr(o.EndDate),
			Status:                       string(o.Status),
			AdminMode:                    string(o.AdminMode),
			ResponsibleEmployeeID:        o.ResponsibleEmployeeID,
			ResponsibleEmployeeFirstName: o.ResponsibleEmployeeFirstName,
			ResponsibleEmployeeLastName:  o.ResponsibleEmployeeLastName,
			IsCritical:                   o.IsCritical,
			Notes:                        o.Notes,
			SourceAttachmentUUID:         o.SourceAttachmentUuid,
			DiagnosisTitle:               o.DiagnosisTitle,
			DiagnosisCodeSystem:          o.DiagnosisCodeSystem,
			DiagnosisCode:                o.DiagnosisCode,
			CreatedAt:                    o.CreatedAt.Time,
			UpdatedAt:                    o.UpdatedAt.Time,
		})
	}

	return &ClientMedicalOverviewResponse{Diagnoses: diagnosisRes, MedicationOrders: orderRes}, nil
}
