package clientp

import (
	"context"
	"maicare_go/async/aclient"
	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/notification"
	"maicare_go/pagination"
	"maicare_go/pdf"
	"maicare_go/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

func (s *clientService) CreateIncident(ctx context.Context, req CreateIncidentRequest, clientID int64) (*CreateIncidentResponse, error) {

	arg := db.CreateIncidentParams{
		EmployeeID:              req.EmployeeID,
		LocationID:              req.LocationID,
		ReporterInvolvement:     req.ReporterInvolvement,
		InformWho:               req.InformWho,
		IncidentDate:            pgtype.Date{Time: req.IncidentDate, Valid: true},
		RuntimeIncident:         req.RuntimeIncident,
		IncidentType:            req.IncidentType,
		PassingAway:             req.PassingAway,
		SelfHarm:                req.SelfHarm,
		Violence:                req.Violence,
		FireWaterDamage:         req.FireWaterDamage,
		Accident:                req.Accident,
		ClientAbsence:           req.ClientAbsence,
		Medicines:               req.Medicines,
		Organization:            req.Organization,
		UseProhibitedSubstances: req.UseProhibitedSubstances,
		OtherNotifications:      req.OtherNotifications,
		SeverityOfIncident:      req.SeverityOfIncident,
		IncidentExplanation:     req.IncidentExplanation,
		RecurrenceRisk:          req.RecurrenceRisk,
		IncidentPreventSteps:    req.IncidentPreventSteps,
		IncidentTakenMeasures:   req.IncidentTakenMeasures,
		Technical:               req.Technical,
		Organizational:          req.Organizational,
		MeseWorker:              req.MeseWorker,
		ClientOptions:           req.ClientOptions,
		OtherCause:              req.OtherCause,
		CauseExplanation:        req.CauseExplanation,
		PhysicalInjury:          req.PhysicalInjury,
		PhysicalInjuryDesc:      req.PhysicalInjuryDesc,
		PsychologicalDamage:     req.PsychologicalDamage,
		PsychologicalDamageDesc: req.PsychologicalDamageDesc,
		NeededConsultation:      req.NeededConsultation,
		Succession:              req.Succession,
		SuccessionDesc:          req.SuccessionDesc,
		Other:                   req.Other,
		OtherDesc:               req.OtherDesc,
		AdditionalAppointments:  req.AdditionalAppointments,
		EmployeeAbsenteeism:     req.EmployeeAbsenteeism,
		ClientID:                clientID,
		Emails:                  req.Emails,
	}

	incident, err := s.Store.CreateIncident(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateIncident", "Failed to create incident", zap.Error(err))
		return nil, err
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "CreateIncident", "Incident created successfully", zap.Int64("IncidentID", incident.ID))

	err = s.AsynqClient.EnqueueIncident(aclient.IncidentPayload{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       "",
		EmployeeLastName:        "",
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               req.InformWho,
		IncidentDate:            incident.IncidentDate.Time,
		RuntimeIncident:         incident.RuntimeIncident,
		IncidentType:            incident.IncidentType,
		PassingAway:             incident.PassingAway,
		SelfHarm:                incident.SelfHarm,
		Violence:                incident.Violence,
		FireWaterDamage:         incident.FireWaterDamage,
		Accident:                incident.Accident,
		ClientAbsence:           incident.ClientAbsence,
		Medicines:               incident.Medicines,
		Organization:            incident.Organization,
		UseProhibitedSubstances: incident.UseProhibitedSubstances,
		OtherNotifications:      incident.OtherNotifications,
		SeverityOfIncident:      incident.SeverityOfIncident,
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          incident.RecurrenceRisk,
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		Technical:               req.Technical,
		Organizational:          req.Organizational,
		MeseWorker:              req.MeseWorker,
		ClientOptions:           req.ClientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              req.Succession,
		SuccessionDesc:          incident.SuccessionDesc,
		Other:                   incident.Other,
		OtherDesc:               incident.OtherDesc,
		AdditionalAppointments:  incident.AdditionalAppointments,
		EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
		ClientID:                incident.ClientID,
		LocationName:            "",
		To:                      incident.Emails,
	}, ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateIncident", "Failed to enqueue incident email task", zap.Error(err))
	}

	receipients, err := s.Store.GetAllAdminUsers(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateIncident", "Failed to get admin users for incident notification", zap.Error(err))

	} else {
		var recipientUserIDs []int64
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
				SeverityOfIncident: incident.SeverityOfIncident,
			}
			err = s.AsynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
				RecipientUserIDs: recipientUserIDs,
				Type:             notification.TypeNewClientAssignment,
				Data: notification.NotificationData{
					NewIncidentReport: &notificationData,
				},
				CreatedAt: time.Now(),
			})
			if err != nil {
				s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateIncident", "Failed to enqueue incident notification task", zap.Error(err))
			}

		}
	}

	response := &CreateIncidentResponse{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               req.InformWho,
		IncidentDate:            incident.IncidentDate.Time,
		RuntimeIncident:         incident.RuntimeIncident,
		IncidentType:            incident.IncidentType,
		PassingAway:             incident.PassingAway,
		SelfHarm:                incident.SelfHarm,
		Violence:                incident.Violence,
		FireWaterDamage:         incident.FireWaterDamage,
		Accident:                incident.Accident,
		ClientAbsence:           incident.ClientAbsence,
		Medicines:               incident.Medicines,
		Organization:            incident.Organization,
		UseProhibitedSubstances: incident.UseProhibitedSubstances,
		OtherNotifications:      incident.OtherNotifications,
		SeverityOfIncident:      incident.SeverityOfIncident,
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          incident.RecurrenceRisk,
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		Technical:               req.Technical,
		Organizational:          req.Organizational,
		MeseWorker:              req.MeseWorker,
		ClientOptions:           req.ClientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              req.Succession,
		SuccessionDesc:          incident.SuccessionDesc,
		Other:                   incident.Other,
		OtherDesc:               incident.OtherDesc,
		AdditionalAppointments:  incident.AdditionalAppointments,
		EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
		ClientID:                incident.ClientID,
		Emails:                  incident.Emails,
		SoftDelete:              incident.SoftDelete,
		UpdatedAt:               incident.UpdatedAt.Time,
		CreatedAt:               incident.CreatedAt.Time,
	}
	return response, nil
}

func (s *clientService) ListIncidents(ctx *gin.Context, req ListIncidentsRequest, clientID int64) (*pagination.Response[ListIncidentsResponse], error) {
	params := req.GetParams()

	incidents, err := s.Store.ListIncidents(ctx, db.ListIncidentsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListIncidents", "Failed to list incidents", zap.Error(err))
		return nil, err
	}

	if len(incidents) == 0 {
		s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ListIncidents", "No incidents found for client", zap.Int64("ClientID", clientID))
		return nil, nil
	}
	totalCount := incidents[0].TotalCount

	var incidentResponses []ListIncidentsResponse
	for _, incident := range incidents {
		incidentResponses = append(incidentResponses, ListIncidentsResponse{
			ID:                      incident.ID,
			EmployeeID:              incident.EmployeeID,
			EmployeeFirstName:       incident.EmployeeFirstName,
			EmployeeLastName:        incident.EmployeeLastName,
			LocationID:              incident.LocationID,
			ReporterInvolvement:     incident.ReporterInvolvement,
			InformWho:               incident.InformWho,
			IncidentDate:            incident.IncidentDate.Time,
			RuntimeIncident:         incident.RuntimeIncident,
			IncidentType:            incident.IncidentType,
			PassingAway:             incident.PassingAway,
			SelfHarm:                incident.SelfHarm,
			Violence:                incident.Violence,
			FireWaterDamage:         incident.FireWaterDamage,
			Accident:                incident.Accident,
			ClientAbsence:           incident.ClientAbsence,
			Medicines:               incident.Medicines,
			Organization:            incident.Organization,
			UseProhibitedSubstances: incident.UseProhibitedSubstances,
			OtherNotifications:      incident.OtherNotifications,
			SeverityOfIncident:      incident.SeverityOfIncident,
			IncidentExplanation:     incident.IncidentExplanation,
			RecurrenceRisk:          incident.RecurrenceRisk,
			IncidentPreventSteps:    incident.IncidentPreventSteps,
			IncidentTakenMeasures:   incident.IncidentTakenMeasures,
			Technical:               incident.Technical,
			Organizational:          incident.Organizational,
			MeseWorker:              incident.MeseWorker,
			ClientOptions:           incident.ClientOptions,
			OtherCause:              incident.OtherCause,
			CauseExplanation:        incident.CauseExplanation,
			PhysicalInjury:          incident.PhysicalInjury,
			PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
			PsychologicalDamage:     incident.PsychologicalDamage,
			PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
			NeededConsultation:      incident.NeededConsultation,
			Succession:              incident.Succession,
			SuccessionDesc:          incident.SuccessionDesc,
			Other:                   incident.Other,
			OtherDesc:               incident.OtherDesc,
			AdditionalAppointments:  incident.AdditionalAppointments,
			EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
			ClientID:                incident.ClientID,
			Emails:                  incident.Emails,
			SoftDelete:              incident.SoftDelete,
			UpdatedAt:               incident.UpdatedAt.Time,
			CreatedAt:               incident.CreatedAt.Time,
			IsConfirmed:             incident.IsConfirmed,
			EmployeeProfilePicture:  incident.EmployeeProfilePicture,
			LocationName:            incident.LocationName,
		})
	}

	paginatedResponse := pagination.NewResponse(ctx, req.Request, incidentResponses, totalCount)
	return &paginatedResponse, nil
}

func (s *clientService) GetIncident(ctx context.Context, incidentID int64) (*GetIncidentResponse, error) {
	incident, err := s.Store.GetIncident(ctx, incidentID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetIncident", "Failed to get incident", zap.Error(err))
		return nil, err
	}

	response := &GetIncidentResponse{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       incident.EmployeeFirstName,
		EmployeeLastName:        incident.EmployeeLastName,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               incident.InformWho,
		IncidentDate:            incident.IncidentDate.Time,
		RuntimeIncident:         incident.RuntimeIncident,
		IncidentType:            incident.IncidentType,
		PassingAway:             incident.PassingAway,
		SelfHarm:                incident.SelfHarm,
		Violence:                incident.Violence,
		FireWaterDamage:         incident.FireWaterDamage,
		Accident:                incident.Accident,
		ClientAbsence:           incident.ClientAbsence,
		Medicines:               incident.Medicines,
		Organization:            incident.Organization,
		UseProhibitedSubstances: incident.UseProhibitedSubstances,
		OtherNotifications:      incident.OtherNotifications,
		SeverityOfIncident:      incident.SeverityOfIncident,
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          incident.RecurrenceRisk,
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		Technical:               incident.Technical,
		Organizational:          incident.Organizational,
		MeseWorker:              incident.MeseWorker,
		ClientOptions:           incident.ClientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              incident.Succession,
		SuccessionDesc:          incident.SuccessionDesc,
		Other:                   incident.Other,
		OtherDesc:               incident.OtherDesc,
		AdditionalAppointments:  incident.AdditionalAppointments,
		EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
		ClientID:                incident.ClientID,
		SoftDelete:              incident.SoftDelete,
		UpdatedAt:               incident.UpdatedAt.Time,
		CreatedAt:               incident.CreatedAt.Time,
		IsConfirmed:             incident.IsConfirmed,
		LocationName:            incident.LocationName,
		Emails:                  incident.Emails,
	}
	return response, nil
}

func (s *clientService) UpdateIncident(ctx context.Context, req UpdateIncidentRequest, incidentID int64) (*UpdateIncidentResponse, error) {
	arg := db.UpdateIncidentParams{
		ID:                      incidentID,
		EmployeeID:              req.EmployeeID,
		LocationID:              req.LocationID,
		ReporterInvolvement:     req.ReporterInvolvement,
		IncidentDate:            pgtype.Date{Time: req.IncidentDate, Valid: true},
		RuntimeIncident:         req.RuntimeIncident,
		IncidentType:            req.IncidentType,
		PassingAway:             req.PassingAway,
		SelfHarm:                req.SelfHarm,
		Violence:                req.Violence,
		FireWaterDamage:         req.FireWaterDamage,
		Accident:                req.Accident,
		ClientAbsence:           req.ClientAbsence,
		Medicines:               req.Medicines,
		Organization:            req.Organization,
		UseProhibitedSubstances: req.UseProhibitedSubstances,
		OtherNotifications:      req.OtherNotifications,
		SeverityOfIncident:      req.SeverityOfIncident,
		IncidentExplanation:     req.IncidentExplanation,
		RecurrenceRisk:          req.RecurrenceRisk,
		IncidentPreventSteps:    req.IncidentPreventSteps,
		IncidentTakenMeasures:   req.IncidentTakenMeasures,
		OtherCause:              req.OtherCause,
		CauseExplanation:        req.CauseExplanation,
		PhysicalInjury:          req.PhysicalInjury,
		PhysicalInjuryDesc:      req.PhysicalInjuryDesc,
		PsychologicalDamage:     req.PsychologicalDamage,
		PsychologicalDamageDesc: req.PsychologicalDamageDesc,
		NeededConsultation:      req.NeededConsultation,
		SuccessionDesc:          req.SuccessionDesc,
		Other:                   req.Other,
		OtherDesc:               req.OtherDesc,
		AdditionalAppointments:  req.AdditionalAppointments,
		EmployeeAbsenteeism:     req.EmployeeAbsenteeism,
		Emails:                  req.Emails,
		InformWho:               req.InformWho,
		Succession:              req.Succession,
		Technical:               req.Technical,
		Organizational:          req.Organizational,
		MeseWorker:              req.MeseWorker,
		ClientOptions:           req.ClientOptions,
	}
	incident, err := s.Store.UpdateIncident(ctx, arg)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateIncident", "Failed to update incident", zap.Error(err))
		return nil, err
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "UpdateIncident", "Incident updated successfully", zap.Int64("IncidentID", incident.ID))
	err = s.AsynqClient.EnqueueIncident(aclient.IncidentPayload{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       "",
		EmployeeLastName:        "",
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               incident.InformWho,
		IncidentDate:            incident.IncidentDate.Time,
		RuntimeIncident:         incident.RuntimeIncident,
		IncidentType:            incident.IncidentType,
		PassingAway:             incident.PassingAway,
		SelfHarm:                incident.SelfHarm,
		Violence:                incident.Violence,
		FireWaterDamage:         incident.FireWaterDamage,
		Accident:                incident.Accident,
		ClientAbsence:           incident.ClientAbsence,
		Medicines:               incident.Medicines,
		Organization:            incident.Organization,
		UseProhibitedSubstances: incident.UseProhibitedSubstances,
		OtherNotifications:      incident.OtherNotifications,
		SeverityOfIncident:      incident.SeverityOfIncident,
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          incident.RecurrenceRisk,
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		Technical:               incident.Technical,
		Organizational:          incident.Organizational,
		MeseWorker:              incident.MeseWorker,
		ClientOptions:           incident.ClientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              incident.Succession,
		SuccessionDesc:          incident.SuccessionDesc,
		Other:                   incident.Other,
		OtherDesc:               incident.OtherDesc,
		AdditionalAppointments:  incident.AdditionalAppointments,
		EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
		ClientID:                incident.ClientID,
		LocationName:            "",
		To:                      incident.Emails,
	}, ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateIncident", "Failed to enqueue incident email task", zap.Error(err))
	}
	response := &UpdateIncidentResponse{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               incident.InformWho,
		IncidentDate:            incident.IncidentDate.Time,
		RuntimeIncident:         incident.RuntimeIncident,
		IncidentType:            incident.IncidentType,
		PassingAway:             incident.PassingAway,
		SelfHarm:                incident.SelfHarm,
		Violence:                incident.Violence,
		FireWaterDamage:         incident.FireWaterDamage,
		Accident:                incident.Accident,
		ClientAbsence:           incident.ClientAbsence,
		Medicines:               incident.Medicines,
		Organization:            incident.Organization,
		UseProhibitedSubstances: incident.UseProhibitedSubstances,
		OtherNotifications:      incident.OtherNotifications,
		SeverityOfIncident:      incident.SeverityOfIncident,
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          incident.RecurrenceRisk,
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		Technical:               incident.Technical,
		Organizational:          incident.Organizational,
		MeseWorker:              incident.MeseWorker,
		ClientOptions:           incident.ClientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              incident.Succession,
		SuccessionDesc:          incident.SuccessionDesc,
		Other:                   incident.Other,
		OtherDesc:               incident.OtherDesc,
		AdditionalAppointments:  incident.AdditionalAppointments,
		EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
		ClientID:                incident.ClientID,
		SoftDelete:              incident.SoftDelete,
		UpdatedAt:               incident.UpdatedAt.Time,
		CreatedAt:               incident.CreatedAt.Time,
		IsConfirmed:             incident.IsConfirmed,
		Emails:                  incident.Emails,
	}
	return response, nil
}

func (s *clientService) DeleteIncident(ctx context.Context, incidentID int64) error {
	err := s.Store.DeleteIncident(ctx, incidentID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteIncident", "Failed to delete incident", zap.Error(err))
		return err
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "DeleteIncident", "Incident deleted successfully", zap.Int64("IncidentID", incidentID))
	return nil
}

func (s *clientService) GenerateIncidentFile(ctx context.Context, incidentID int64) (*GenerateIncidentFileResponse, error) {

	incident, err := s.Store.GetIncident(ctx, incidentID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateIncidentFile", "Failed to get incident", zap.Error(err))
		return nil, err
	}

	incidentData := pdf.IncidentReportData{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       incident.EmployeeFirstName,
		EmployeeLastName:        incident.EmployeeLastName,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformWho:               incident.InformWho,
		IncidentDate:            incident.IncidentDate.Time,
		RuntimeIncident:         incident.RuntimeIncident,
		IncidentType:            incident.IncidentType,
		PassingAway:             incident.PassingAway,
		SelfHarm:                incident.SelfHarm,
		Violence:                incident.Violence,
		FireWaterDamage:         incident.FireWaterDamage,
		Accident:                incident.Accident,
		ClientAbsence:           incident.ClientAbsence,
		Medicines:               incident.Medicines,
		Organization:            incident.Organization,
		UseProhibitedSubstances: incident.UseProhibitedSubstances,
		OtherNotifications:      incident.OtherNotifications,
		SeverityOfIncident:      incident.SeverityOfIncident,
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          incident.RecurrenceRisk,
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		Technical:               incident.Technical,
		Organizational:          incident.Organizational,
		MeseWorker:              incident.MeseWorker,
		ClientOptions:           incident.ClientOptions,
		OtherCause:              incident.OtherCause,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		Succession:              incident.Succession,
		SuccessionDesc:          incident.SuccessionDesc,
		Other:                   incident.Other,
		OtherDesc:               incident.OtherDesc,
		AdditionalAppointments:  incident.AdditionalAppointments,
		EmployeeAbsenteeism:     incident.EmployeeAbsenteeism,
		ClientID:                incident.ClientID,
		ClientFirstName:         incident.ClientFirstName,
		ClientLastName:          incident.ClientLastName,
		LocationName:            incident.LocationName,
	}
	fileKey, err := pdf.GenerateAndUploadIncidentPDF(ctx, incidentData, s.B2Client)
	if err != nil && fileKey == "" {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateIncidentFile", "Failed to generate incident PDF", zap.Error(err))
		return nil, err
	}

	incidentWithUpdatedFileUrl, err := s.Store.UpdateIncidentFileUrl(ctx, db.UpdateIncidentFileUrlParams{
		ID:      incident.ID,
		FileUrl: &fileKey,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GenerateIncidentFile", "Failed to update incident file URL", zap.Error(err))
		return nil, err
	}
	response := &GenerateIncidentFileResponse{
		FileUrl: s.GenerateResponsePresignedURL(incidentWithUpdatedFileUrl, ctx),
		ID:      incident.ID,
	}
	return response, nil
}

func (s *clientService) ConfirmIncident(ctx context.Context, incidentID int64) (*ConfirmIncidentResponse, error) {
	incident, err := s.Store.ConfirmIncident(ctx, incidentID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ConfirmIncident", "Failed to confirm incident", zap.Error(err))
		return nil, err
	}

	s.Logger.LogBusinessEvent(logger.LogLevelInfo, "ConfirmIncident", "Incident confirmed successfully", zap.Int64("IncidentID", incidentID))
	return &ConfirmIncidentResponse{
		FileUrl: s.GenerateResponsePresignedURL(incident.FileUrl, ctx),
		ID:      incident.ID,
	}, nil
	// TODO: Send notification to the party responsivle for the client

}
