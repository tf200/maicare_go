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
}

func NewIncidentService(repo domain.IncidentRepository, pdfGenerator domain.IncidentPDFGenerator, taskQueue domain.TaskQueue, logger domain.Logger) domain.IncidentService {
	return &incidentService{
		repo:         repo,
		pdfGenerator: pdfGenerator,
		taskQueue:    taskQueue,
		logger:       logger,
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
	return incident, nil
}

func (s *incidentService) DeleteIncident(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteIncident(ctx, id)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.DeleteIncident", "failed to delete incident", err,
			zap.String("incident_id", id.String()),
		)
		return err
	}

	s.logger.LogInfo(ctx, "IncidentService.DeleteIncident", "incident deleted",
		zap.String("incident_id", id.String()),
	)
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

	return pdfBytes, fmt.Sprintf("incident_report_%s.pdf", incident.ID.String()), nil
}

func (s *incidentService) ConfirmIncident(ctx context.Context, id uuid.UUID, confirmedByUserID uuid.UUID) (*domain.ConfirmIncidentResult, error) {
	rowsAffected, err := s.repo.ConfirmIncident(ctx, id, &confirmedByUserID)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.ConfirmIncident", "failed to confirm incident", err,
			zap.String("incident_id", id.String()),
			zap.String("confirmed_by", confirmedByUserID.String()),
		)
		return nil, err
	}

	if rowsAffected == 0 {
		s.logger.LogInfo(ctx, "IncidentService.ConfirmIncident", "incident already confirmed",
			zap.String("incident_id", id.String()),
		)
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
	return &domain.ConfirmIncidentResult{ID: id, FileUrl: nil}, nil
}

func (s *incidentService) ListAllIncidents(ctx context.Context, params domain.ListAllIncidentsParams) (*domain.ListAllIncidentsResult, error) {
	result, err := s.repo.ListAllIncidents(ctx, params)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.ListAllIncidents", "failed to list all incidents", err)
		return nil, err
	}
	return result, nil
}

func (s *incidentService) GetIncidentCounts(ctx context.Context) (*domain.IncidentCounts, error) {
	counts, err := s.repo.GetIncidentCounts(ctx)
	if err != nil {
		s.logger.LogError(ctx, "IncidentService.GetIncidentCounts", "failed to get incident counts", err)
		return nil, err
	}
	return counts, nil
}
