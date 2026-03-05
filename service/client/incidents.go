package clientp

import (
	"context"
	"fmt"
	"time"

	"maicare_go/async/aclient"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/service/notification"
	"maicare_go/service/pdf"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func toInformedParties(values []string) []db.InformedPartyEnum {
	result := make([]db.InformedPartyEnum, 0, len(values))
	for _, value := range values {
		result = append(result, db.InformedPartyEnum(value))
	}
	return result
}

func informedPartiesToStrings(values []db.InformedPartyEnum) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, string(value))
	}
	return result
}

func toCauseCategories(values []string) []db.IncidentCauseCategoryEnum {
	result := make([]db.IncidentCauseCategoryEnum, 0, len(values))
	for _, value := range values {
		result = append(result, db.IncidentCauseCategoryEnum(value))
	}
	return result
}

func causeCategoriesToStrings(values []db.IncidentCauseCategoryEnum) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, string(value))
	}
	return result
}

func toFollowUpActions(values []string) []db.IncidentFollowUpActionEnum {
	result := make([]db.IncidentFollowUpActionEnum, 0, len(values))
	for _, value := range values {
		result = append(result, db.IncidentFollowUpActionEnum(value))
	}
	return result
}

func followUpActionsToStrings(values []db.IncidentFollowUpActionEnum) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, string(value))
	}
	return result
}

func nullIncidentTypeFromPtr(ptr *string) db.NullIncidentTypeEnum {
	if ptr == nil {
		return db.NullIncidentTypeEnum{Valid: false}
	}
	return db.NullIncidentTypeEnum{IncidentTypeEnum: db.IncidentTypeEnum(*ptr), Valid: true}
}

func (s *clientService) CreateIncident(ctx context.Context, req CreateIncidentRequest) (*CreateIncidentResponse, error) {
	arg := db.CreateIncidentParams{
		EmployeeID:              req.EmployeeID,
		LocationID:              req.LocationID,
		ReporterInvolvement:     db.IncidentReporterInvolvementEnum(req.ReporterInvolvement),
		InformedParties:         toInformedParties(req.InformedParties),
		OccurredAt:              pgtype.Timestamptz{Time: req.OccurredAt, Valid: true},
		IncidentType:            db.IncidentTypeEnum(req.IncidentType),
		SeverityOfIncident:      db.SeverityOfIncidentEnum(req.SeverityOfIncident),
		IncidentExplanation:     req.IncidentExplanation,
		RecurrenceRisk:          db.RecurrenceRiskEnum(req.RecurrenceRisk),
		IncidentPreventSteps:    req.IncidentPreventSteps,
		IncidentTakenMeasures:   req.IncidentTakenMeasures,
		CauseCategories:         toCauseCategories(req.CauseCategories),
		CauseExplanation:        req.CauseExplanation,
		PhysicalInjury:          db.PhysicalInjuryEnum(req.PhysicalInjury),
		PhysicalInjuryDesc:      req.PhysicalInjuryDesc,
		PsychologicalDamage:     db.PsychologicalDamageEnum(req.PsychologicalDamage),
		PsychologicalDamageDesc: req.PsychologicalDamageDesc,
		NeededConsultation:      db.NeededConsultationEnum(req.NeededConsultation),
		FollowUpActions:         toFollowUpActions(req.FollowUpActions),
		FollowUpNotes:           req.FollowUpNotes,
		IsEmployeeAbsent:        req.IsEmployeeAbsent,
		AdditionalDetails:       req.AdditionalDetails,
		ClientID:                req.ClientID,
		Emails:                  req.Emails,
	}

	var incident db.CreateIncidentRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		incident, err = q.CreateIncident(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateIncident", "Failed to create incident", zap.Error(err))
		return nil, err
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "CreateIncident", "Incident created successfully", zap.String("IncidentID", incident.ID.String()))

	receipients, err := s.Store.GetAllAdminUsers(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateIncident", "Failed to get admin users for incident notification", zap.Error(err))
	} else {
		var recipientUserIDs []uuid.UUID
		for _, user := range receipients {
			recipientUserIDs = append(recipientUserIDs, user.ID)
		}
		if len(recipientUserIDs) > 0 {
			notificationData := notification.NewIncidentReportData{
				ID:                 incident.ID,
				EmployeeID:         incident.EmployeeID,
				EmployeeFirstName:  util.DerefString(incident.EmployeeFirstName),
				EmployeeLastName:   util.DerefString(incident.EmployeeLastName),
				LocationID:         incident.LocationID,
				LocationName:       util.DerefString(incident.LocationName),
				ClientID:           incident.ClientID,
				ClientFirstName:    util.DerefString(incident.ClientFirstName),
				ClientLastName:     util.DerefString(incident.ClientLastName),
				SeverityOfIncident: string(incident.SeverityOfIncident),
			}
			message := fmt.Sprintf(
				"New incident reported for %s %s (%s)",
				notificationData.ClientFirstName,
				notificationData.ClientLastName,
				notificationData.SeverityOfIncident,
			)
			err = s.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
				RecipientUserIDs: recipientUserIDs,
				Type:             notification.TypeIncidentReport,
				Data: notification.NotificationData{
					NewIncidentReport: &notificationData,
				},
				CreatedAt: time.Now(),
				Message:   message,
			})
			if err != nil {
				s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateIncident", "Failed to enqueue incident notification task", zap.Error(err))
			}

		}
	}

	response := &CreateIncidentResponse{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     string(incident.ReporterInvolvement),
		InformedParties:         informedPartiesToStrings(incident.InformedParties),
		OccurredAt:              incident.OccurredAt.Time,
		IncidentType:            string(incident.IncidentType),
		SeverityOfIncident:      string(incident.SeverityOfIncident),
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          string(incident.RecurrenceRisk),
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		CauseCategories:         causeCategoriesToStrings(incident.CauseCategories),
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          string(incident.PhysicalInjury),
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     string(incident.PsychologicalDamage),
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      string(incident.NeededConsultation),
		FollowUpActions:         followUpActionsToStrings(incident.FollowUpActions),
		FollowUpNotes:           incident.FollowUpNotes,
		IsEmployeeAbsent:        incident.IsEmployeeAbsent,
		AdditionalDetails:       incident.AdditionalDetails,
		ClientID:                incident.ClientID,
		Emails:                  incident.Emails,
		UpdatedAt:               incident.UpdatedAt.Time,
		CreatedAt:               incident.CreatedAt.Time,
	}
	return response, nil
}

func (s *clientService) ListIncidents(ctx *gin.Context, req ListIncidentsRequest, clientID uuid.UUID) (*pagination.Response[ListIncidentsResponse], error) {
	params := req.GetParams()

	var incidents []db.ListIncidentsRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		incidents, err = q.ListIncidents(ctx, db.ListIncidentsParams{
			ClientID: clientID,
			Limit:    params.Limit,
			Offset:   params.Offset,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListIncidents", "Failed to list incidents", zap.Error(err))
		return nil, err
	}

	if len(incidents) == 0 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ListIncidents", "No incidents found for client", zap.String("ClientID", clientID.String()))
		pag := pagination.NewResponse(ctx, req.Request, []ListIncidentsResponse{}, 0)
		return &pag, nil
	}
	totalCount := incidents[0].TotalCount

	var incidentResponses []ListIncidentsResponse
	for _, incident := range incidents {
		incidentResponses = append(incidentResponses, ListIncidentsResponse{
			ID:                     incident.ID,
			OccurredAt:             incident.OccurredAt.Time,
			IncidentType:           string(incident.IncidentType),
			SeverityOfIncident:     string(incident.SeverityOfIncident),
			IsConfirmed:            incident.IsConfirmed,
			EmployeeFirstName:      incident.EmployeeFirstName,
			EmployeeLastName:       incident.EmployeeLastName,
			EmployeeProfilePicture: incident.EmployeeProfilePicture,
			LocationName:           incident.LocationName,
		})
	}

	paginatedResponse := pagination.NewResponse(ctx, req.Request, incidentResponses, totalCount)
	return &paginatedResponse, nil
}

func (s *clientService) GetIncident(ctx context.Context, incidentID uuid.UUID) (*GetIncidentResponse, error) {
	var incident db.GetIncidentRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		incident, err = q.GetIncident(ctx, incidentID)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetIncident", "Failed to get incident", zap.Error(err))
		return nil, err
	}

	response := &GetIncidentResponse{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       incident.EmployeeFirstName,
		EmployeeLastName:        incident.EmployeeLastName,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     string(incident.ReporterInvolvement),
		InformedParties:         informedPartiesToStrings(incident.InformedParties),
		OccurredAt:              incident.OccurredAt.Time,
		IncidentType:            string(incident.IncidentType),
		SeverityOfIncident:      string(incident.SeverityOfIncident),
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          string(incident.RecurrenceRisk),
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		CauseCategories:         causeCategoriesToStrings(incident.CauseCategories),
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          string(incident.PhysicalInjury),
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     string(incident.PsychologicalDamage),
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      string(incident.NeededConsultation),
		FollowUpActions:         followUpActionsToStrings(incident.FollowUpActions),
		FollowUpNotes:           incident.FollowUpNotes,
		IsEmployeeAbsent:        incident.IsEmployeeAbsent,
		AdditionalDetails:       incident.AdditionalDetails,
		ClientID:                incident.ClientID,
		UpdatedAt:               incident.UpdatedAt.Time,
		CreatedAt:               incident.CreatedAt.Time,
		IsConfirmed:             incident.IsConfirmed,
		LocationName:            incident.LocationName,
		Emails:                  incident.Emails,
	}
	return response, nil
}

func (s *clientService) UpdateIncident(ctx context.Context, req UpdateIncidentRequest, incidentID uuid.UUID) (*UpdateIncidentResponse, error) {
	arg := db.UpdateIncidentParams{
		ID:                      incidentID,
		EmployeeID:              req.EmployeeID,
		LocationID:              req.LocationID,
		ReporterInvolvement:     db.NullIncidentReporterInvolvementFromPtr(req.ReporterInvolvement),
		InformedParties:         toInformedParties(req.InformedParties),
		OccurredAt:              pgtype.Timestamptz{Time: req.OccurredAt, Valid: true},
		IncidentType:            nullIncidentTypeFromPtr(req.IncidentType),
		SeverityOfIncident:      db.NullSeverityOfIncidentFromPtr(req.SeverityOfIncident),
		IncidentExplanation:     req.IncidentExplanation,
		RecurrenceRisk:          db.NullRecurrenceRiskFromPtr(req.RecurrenceRisk),
		IncidentPreventSteps:    req.IncidentPreventSteps,
		IncidentTakenMeasures:   req.IncidentTakenMeasures,
		CauseCategories:         toCauseCategories(req.CauseCategories),
		CauseExplanation:        req.CauseExplanation,
		PhysicalInjury:          db.NullPhysicalInjuryFromPtr(req.PhysicalInjury),
		PhysicalInjuryDesc:      req.PhysicalInjuryDesc,
		PsychologicalDamage:     db.NullPsychologicalDamageFromPtr(req.PsychologicalDamage),
		PsychologicalDamageDesc: req.PsychologicalDamageDesc,
		NeededConsultation:      db.NullNeededConsultationFromPtr(req.NeededConsultation),
		FollowUpActions:         toFollowUpActions(req.FollowUpActions),
		FollowUpNotes:           req.FollowUpNotes,
		IsEmployeeAbsent:        req.IsEmployeeAbsent,
		AdditionalDetails:       req.AdditionalDetails,
		Emails:                  req.Emails,
	}
	var incident db.Incident
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		incident, err = q.UpdateIncident(ctx, arg)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateIncident", "Failed to update incident", zap.Error(err))
		return nil, err
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "UpdateIncident", "Incident updated successfully", zap.String("IncidentID", incident.ID.String()))

	response := &UpdateIncidentResponse{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     string(incident.ReporterInvolvement),
		InformedParties:         informedPartiesToStrings(incident.InformedParties),
		OccurredAt:              incident.OccurredAt.Time,
		IncidentType:            string(incident.IncidentType),
		SeverityOfIncident:      string(incident.SeverityOfIncident),
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          string(incident.RecurrenceRisk),
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		CauseCategories:         causeCategoriesToStrings(incident.CauseCategories),
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          string(incident.PhysicalInjury),
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     string(incident.PsychologicalDamage),
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      string(incident.NeededConsultation),
		FollowUpActions:         followUpActionsToStrings(incident.FollowUpActions),
		FollowUpNotes:           incident.FollowUpNotes,
		IsEmployeeAbsent:        incident.IsEmployeeAbsent,
		AdditionalDetails:       incident.AdditionalDetails,
		ClientID:                incident.ClientID,
		UpdatedAt:               incident.UpdatedAt.Time,
		CreatedAt:               incident.CreatedAt.Time,
		IsConfirmed:             incident.IsConfirmed,
		Emails:                  incident.Emails,
	}
	return response, nil
}

func (s *clientService) DeleteIncident(ctx context.Context, incidentID uuid.UUID) error {
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		return q.DeleteIncident(ctx, incidentID)
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteIncident", "Failed to delete incident", zap.Error(err))
		return err
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "DeleteIncident", "Incident deleted successfully", zap.String("IncidentID", incidentID.String()))
	return nil
}

func (s *clientService) GenerateIncidentFile(ctx context.Context, incidentID uuid.UUID) ([]byte, string, error) {
	var incident db.GetIncidentRow
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		incident, err = q.GetIncident(ctx, incidentID)
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GenerateIncidentFile", "Failed to get incident", zap.Error(err))
		return nil, "", err
	}

	incidentData := pdf.IncidentReportData{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       incident.EmployeeFirstName,
		EmployeeLastName:        incident.EmployeeLastName,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     string(incident.ReporterInvolvement),
		InformedParties:         informedPartiesToStrings(incident.InformedParties),
		OccurredAt:              incident.OccurredAt.Time,
		IncidentType:            string(incident.IncidentType),
		SeverityOfIncident:      string(incident.SeverityOfIncident),
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          string(incident.RecurrenceRisk),
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		CauseCategories:         causeCategoriesToStrings(incident.CauseCategories),
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          string(incident.PhysicalInjury),
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     string(incident.PsychologicalDamage),
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      string(incident.NeededConsultation),
		FollowUpActions:         followUpActionsToStrings(incident.FollowUpActions),
		FollowUpNotes:           incident.FollowUpNotes,
		IsEmployeeAbsent:        incident.IsEmployeeAbsent,
		AdditionalDetails:       incident.AdditionalDetails,
		ClientID:                incident.ClientID,
		ClientFirstName:         incident.ClientFirstName,
		ClientLastName:          incident.ClientLastName,
		LocationName:            incident.LocationName,
	}
	pdfBytes, err := s.PDFService.GenerateIncidentPDF(ctx, incidentData)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GenerateIncidentFile", "Failed to generate incident PDF", zap.Error(err))
		return nil, "", err
	}

	fileName := fmt.Sprintf("incident_report_%s.pdf", incident.ID.String())
	return pdfBytes, fileName, nil
}

func (s *clientService) ConfirmIncident(ctx context.Context, incidentID uuid.UUID, confirmedByUserID uuid.UUID) (*ConfirmIncidentResponse, error) {
	var rowsAffected int64
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		rowsAffected, err = q.ConfirmIncident(ctx, db.ConfirmIncidentParams{
			ID:          incidentID,
			ConfirmedBy: &confirmedByUserID,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ConfirmIncident", "Failed to confirm incident", zap.Error(err))
		return nil, err
	}

	if rowsAffected == 0 {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ConfirmIncident", "Incident already confirmed", zap.String("IncidentID", incidentID.String()))
		return &ConfirmIncidentResponse{ID: incidentID, FileUrl: nil}, nil
	}

	if enqueueErr := s.asynqClient.EnqueueIncidentConfirmedEmail(ctx, aclient.IncidentConfirmedEmailPayload{IncidentID: incidentID}); enqueueErr != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ConfirmIncident", "Failed to enqueue incident confirmation email task", zap.Error(enqueueErr), zap.String("incident_id", incidentID.String()))
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "ConfirmIncident", "Incident confirmed successfully", zap.String("IncidentID", incidentID.String()))
	return &ConfirmIncidentResponse{
		FileUrl: nil,
		ID:      incidentID,
	}, nil
}

func (s *clientService) ListAllIncidents(ctx *gin.Context, req *ListAllIncidentsRequest) (*pagination.Response[ListAllIncidentsResponse], error) {
	params := req.GetParams()
	arg := db.ListAllIncidentsParams{
		Limit:       params.Limit,
		Offset:      params.Offset,
		IsConfirmed: req.IsConfirmed,
		Search:      req.Search,
	}
	var incidents []db.ListAllIncidentsRow
	var count int64
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		var err error
		incidents, err = q.ListAllIncidents(ctx, arg)
		if err != nil {
			return err
		}
		count, err = q.CountAllIncidents(ctx, db.CountAllIncidentsParams{
			IsConfirmed: req.IsConfirmed,
			Search:      req.Search,
		})
		return err
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListAllIncidents", "Failed to list all incidents", zap.Error(err))
		return nil, err
	}

	response := []ListAllIncidentsResponse{}
	for _, incident := range incidents {
		response = append(response, ListAllIncidentsResponse{
			ID:                 incident.ID,
			OccurredAt:         incident.OccurredAt.Time,
			IncidentType:       string(incident.IncidentType),
			SeverityOfIncident: string(incident.SeverityOfIncident),
			IsConfirmed:        incident.IsConfirmed,
			ClientFirstName:    incident.ClientFirstName,
			ClientLastName:     incident.ClientLastName,
			ClientBSN:          incident.ClientBsn,
			EmployeeFirstName:  incident.EmployeeFirstName,
			EmployeeLastName:   incident.EmployeeLastName,
			LocationName:       incident.LocationName,
		})
	}

	paginatedResponse := pagination.NewResponse(ctx, req.Request, response, count)
	return &paginatedResponse, nil
}

func (s *clientService) GetIncidentCounts(ctx context.Context) (*GetIncidentCountsResponse, error) {
	counts, err := s.Store.GetIncidentCounts(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetIncidentCounts", "Failed to get incident counts", zap.Error(err))
		return nil, err
	}

	return &GetIncidentCountsResponse{
		SeriousFatalCount:        counts.SeriousFatalCount,
		PendingConfirmationCount: counts.PendingConfirmationCount,
		Past24hCount:             counts.Past24hCount,
	}, nil
}
