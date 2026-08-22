package service

import (
	"context"
	"fmt"
	"time"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type incidentService struct {
	repo         domain.IncidentRepository
	pdfGenerator domain.IncidentPDFGenerator
	taskQueue    domain.TaskQueue
	logger       domain.Logger
	audit        domain.AuditLogger
}

func NewIncidentService(repo domain.IncidentRepository, pdfGenerator domain.IncidentPDFGenerator, taskQueue domain.TaskQueue, logger domain.Logger, audit domain.AuditLogger) domain.IncidentService {
	return &incidentService{
		repo:         repo,
		pdfGenerator: pdfGenerator,
		taskQueue:    taskQueue,
		logger:       logger,
		audit:        audit,
	}
}

func (s *incidentService) logIncidentAudit(ctx context.Context, action, subjectID string, clientID *uuid.UUID, permission string, count int) {
	if s.audit == nil {
		return
	}
	eventType := "record_change"
	if action == "read" || action == "list" {
		eventType = "record_access"
	} else if action == "export" {
		eventType = "export"
	}
	var clientIDString *string
	if clientID != nil {
		value := clientID.String()
		clientIDString = &value
	}
	var details map[string]any
	if count >= 0 {
		details = map[string]any{"count": count}
	}
	if err := s.audit.Log(ctx, domain.AuditEvent{
		EventType:   eventType,
		Action:      action,
		Result:      "success",
		SubjectType: "incident",
		SubjectID:   subjectID,
		ClientID:    clientIDString,
		AccessRule:  strPtr(permission),
		Details:     details,
	}); err != nil && s.logger != nil {
		s.logger.LogError(ctx, "IncidentService.logIncidentAudit", "incident audit log failed", err,
			zap.String("subject_id", subjectID), zap.String("action", action))
	}
}

func (s *incidentService) logIncidentListAudit(ctx context.Context, subjectIDs []string, clientID *uuid.UUID) {
	if s.audit == nil || len(subjectIDs) == 0 {
		return
	}
	var clientIDString *string
	if clientID != nil {
		value := clientID.String()
		clientIDString = &value
	}
	if err := s.audit.LogBatch(ctx, domain.AuditBatchEvent{
		EventGroupID: uuid.NewString(),
		EventType:    "record_access",
		Action:       "list",
		Result:       "success",
		SubjectType:  "incident",
		ClientID:     clientIDString,
		AccessRule:   strPtr("CLIENT.INCIDENT.VIEW"),
		SubjectIDs:   subjectIDs,
		Details:      map[string]any{"count": len(subjectIDs)},
	}); err != nil && s.logger != nil {
		s.logger.LogError(ctx, "IncidentService.logIncidentListAudit", "incident list audit log failed", err)
	}
}

func (s *incidentService) CreateIncident(ctx context.Context, params domain.CreateIncidentParams) (*domain.Incident, error) {
	incident, err := s.repo.CreateIncident(ctx, params)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.CreateIncident", "failed to create incident", err,
			zap.String("employee_id", params.EmployeeID.String()),
			zap.String("client_id", params.ClientID.String()),
		)
		return nil, err
	}

	s.logger.LogInfo(ctx, "IncidentService.CreateIncident", "incident created",
		zap.String("incident_id", incident.ID.String()),
	)
	s.logIncidentAudit(ctx, "create", incident.ID.String(), &incident.ClientID, "CLIENT.INCIDENT.CREATE", -1)

	adminUsers, err := s.repo.GetAllAdminUsers(ctx)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.CreateIncident", "failed to get admin users", err)
	}

	if len(adminUsers) > 0 {
		payload := domain.NotificationTaskPayload{
			RecipientUserIDs: adminUsers,
			Type:             "incident_report",
			Data: domain.NotificationTaskData{
				NewIncidentReport: &domain.NewIncidentReportTaskData{
					ID:                 incident.ID,
					EmployeeID:         incident.EmployeeID,
					EmployeeFirstName:  incident.EmployeeFirstName,
					EmployeeLastName:   incident.EmployeeLastName,
					LocationID:         incident.LocationID,
					LocationName:       incident.LocationName,
					ClientID:           incident.ClientID,
					ClientFirstName:    incident.ClientFirstName,
					ClientLastName:     incident.ClientLastName,
					SeverityOfIncident: incident.SeverityOfIncident,
				},
			},
			CreatedAt: time.Now(),
			Message:   fmt.Sprintf("New incident reported for %s %s (%s)", incident.ClientFirstName, incident.ClientLastName, incident.SeverityOfIncident),
		}

		if enqueueErr := s.taskQueue.EnqueueNotificationTask(ctx, payload, nil); enqueueErr != nil {
			s.logger.LogError(ctx, "IncidentService.CreateIncident", "failed to enqueue notification task", enqueueErr,
				zap.String("incident_id", incident.ID.String()),
			)
		}
	}

	return incident, nil
}

func (s *incidentService) ListIncidents(ctx context.Context, params domain.ListIncidentsParams) (*domain.ListIncidentsResult, error) {
	result, err := s.repo.ListIncidents(ctx, params)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.ListIncidents", "failed to list incidents", err,
			zap.String("client_id", params.ClientID.String()),
		)
		return nil, err
	}
	subjectIDs := make([]string, len(result.Items))
	for i, item := range result.Items {
		subjectIDs[i] = item.ID.String()
	}
	s.logIncidentListAudit(ctx, subjectIDs, &params.ClientID)
	return result, nil
}

func (s *incidentService) GetIncident(ctx context.Context, id uuid.UUID) (*domain.Incident, error) {
	incident, err := s.repo.GetIncident(ctx, id)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.GetIncident", "failed to get incident", err,
			zap.String("incident_id", id.String()),
		)
		return nil, err
	}
	s.logIncidentAudit(ctx, "read", incident.ID.String(), &incident.ClientID, "CLIENT.INCIDENT.VIEW", -1)
	return incident, nil
}

func (s *incidentService) UpdateIncident(ctx context.Context, params domain.UpdateIncidentParams) (*domain.Incident, error) {
	incident, err := s.repo.UpdateIncident(ctx, params)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.UpdateIncident", "failed to update incident", err,
			zap.String("incident_id", params.ID.String()),
		)
		return nil, err
	}

	s.logger.LogInfo(ctx, "IncidentService.UpdateIncident", "incident updated",
		zap.String("incident_id", incident.ID.String()),
	)
	s.logIncidentAudit(ctx, "update", incident.ID.String(), &incident.ClientID, "CLIENT.INCIDENT.UPDATE", -1)
	return incident, nil
}

func (s *incidentService) DeleteIncident(ctx context.Context, id uuid.UUID) error {
	clientID, err := s.repo.DeleteIncident(ctx, id)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.DeleteIncident", "failed to delete incident", err,
			zap.String("incident_id", id.String()),
		)
		return err
	}

	s.logger.LogInfo(ctx, "IncidentService.DeleteIncident", "incident deleted",
		zap.String("incident_id", id.String()),
	)
	s.logIncidentAudit(ctx, "delete", id.String(), &clientID, "CLIENT.INCIDENT.DELETE", -1)
	return nil
}

func (s *incidentService) GenerateIncidentFile(ctx context.Context, id uuid.UUID) ([]byte, string, error) {
	incident, err := s.repo.GetIncident(ctx, id)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.GenerateIncidentFile", "failed to get incident", err,
			zap.String("incident_id", id.String()),
		)
		return nil, "", err
	}

	data := domain.IncidentPDFData{
		ID:                      incident.ID,
		EmployeeID:              incident.EmployeeID,
		EmployeeFirstName:       incident.EmployeeFirstName,
		EmployeeLastName:        incident.EmployeeLastName,
		LocationID:              incident.LocationID,
		ReporterInvolvement:     incident.ReporterInvolvement,
		InformedParties:         incident.InformedParties,
		OccurredAt:              incident.OccurredAt,
		IncidentType:            incident.IncidentType,
		SeverityOfIncident:      incident.SeverityOfIncident,
		IncidentExplanation:     incident.IncidentExplanation,
		RecurrenceRisk:          incident.RecurrenceRisk,
		IncidentPreventSteps:    incident.IncidentPreventSteps,
		IncidentTakenMeasures:   incident.IncidentTakenMeasures,
		CauseCategories:         incident.CauseCategories,
		CauseExplanation:        incident.CauseExplanation,
		PhysicalInjury:          incident.PhysicalInjury,
		PhysicalInjuryDesc:      incident.PhysicalInjuryDesc,
		PsychologicalDamage:     incident.PsychologicalDamage,
		PsychologicalDamageDesc: incident.PsychologicalDamageDesc,
		NeededConsultation:      incident.NeededConsultation,
		FollowUpActions:         incident.FollowUpActions,
		FollowUpNotes:           incident.FollowUpNotes,
		IsEmployeeAbsent:        incident.IsEmployeeAbsent,
		AdditionalDetails:       incident.AdditionalDetails,
		ClientID:                incident.ClientID,
		ClientFirstName:         incident.ClientFirstName,
		ClientLastName:          incident.ClientLastName,
		LocationName:            incident.LocationName,
	}

	pdfBytes, err := s.pdfGenerator.GenerateIncidentPDF(ctx, data)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.GenerateIncidentFile", "failed to generate incident PDF", err,
			zap.String("incident_id", id.String()),
		)
		return nil, "", err
	}
	s.logIncidentAudit(ctx, "export", incident.ID.String(), &incident.ClientID, "CLIENT.INCIDENT.VIEW", -1)

	return pdfBytes, fmt.Sprintf("incident_report_%s.pdf", incident.ID.String()), nil
}

func (s *incidentService) ConfirmIncident(ctx context.Context, id uuid.UUID, confirmedByUserID uuid.UUID) (*domain.ConfirmIncidentResult, error) {
	incident, err := s.repo.GetIncident(ctx, id)
	if err != nil {
		return nil, err
	}
	rowsAffected, err := s.repo.ConfirmIncident(ctx, id)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.ConfirmIncident", "failed to confirm incident", err,
			zap.String("incident_id", id.String()),
			zap.String("confirmed_by", confirmedByUserID.String()),
		)
		return nil, err
	}

	if rowsAffected == 0 {
		incident, err = s.repo.GetIncident(ctx, id)
		if err != nil {
			return nil, err
		}
		if !incident.IsConfirmed {
			return nil, domain.ErrIncidentNotFound
		}
		s.logger.LogInfo(ctx, "IncidentService.ConfirmIncident", "incident already confirmed",
			zap.String("incident_id", id.String()),
		)
		if incident.ConfirmationEmailSentAt == nil {
			if enqueueErr := s.taskQueue.EnqueueIncidentConfirmedEmail(ctx, domain.IncidentConfirmedEmailTaskPayload{IncidentID: id}, nil); enqueueErr != nil {
				s.logger.LogError(ctx, "IncidentService.ConfirmIncident", "failed to re-enqueue incident confirmed email", enqueueErr,
					zap.String("incident_id", id.String()))
			}
		}
		return &domain.ConfirmIncidentResult{ID: id, FileUrl: nil}, nil
	}

	if enqueueErr := s.taskQueue.EnqueueIncidentConfirmedEmail(ctx, domain.IncidentConfirmedEmailTaskPayload{IncidentID: id}, nil); enqueueErr != nil {
		s.logger.LogError(ctx, "IncidentService.ConfirmIncident", "failed to enqueue incident confirmed email", enqueueErr,
			zap.String("incident_id", id.String()),
		)
	}

	s.logger.LogInfo(ctx, "IncidentService.ConfirmIncident", "incident confirmed",
		zap.String("incident_id", id.String()),
		zap.String("confirmed_by", confirmedByUserID.String()),
	)
	s.logIncidentAudit(ctx, "confirm", id.String(), &incident.ClientID, "CLIENT.INCIDENT.CONFIRM", -1)
	return &domain.ConfirmIncidentResult{ID: id, FileUrl: nil}, nil
}

func (s *incidentService) ListAllIncidents(ctx context.Context, params domain.ListAllIncidentsParams) (*domain.ListAllIncidentsResult, error) {
	result, err := s.repo.ListAllIncidents(ctx, params)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.ListAllIncidents", "failed to list all incidents", err)
		return nil, err
	}
	subjectIDsByClient := make(map[uuid.UUID][]string)
	for _, item := range result.Items {
		subjectIDsByClient[item.ClientID] = append(subjectIDsByClient[item.ClientID], item.ID.String())
	}
	for clientID, subjectIDs := range subjectIDsByClient {
		s.logIncidentListAudit(ctx, subjectIDs, &clientID)
	}
	return result, nil
}

func (s *incidentService) GetIncidentCounts(ctx context.Context) (*domain.IncidentCounts, error) {
	counts, err := s.repo.GetIncidentCounts(ctx)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.GetIncidentCounts", "failed to get incident counts", err)
		return nil, err
	}
	s.logIncidentAudit(ctx, "read", "incident_counts", nil, "CLIENT.INCIDENT.VIEW", -1)
	return counts, nil
}
