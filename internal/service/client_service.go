package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

const maxFileSize = 100 << 20 // 100MB

type ClientService struct {
	repository domain.ClientRepository
	logger     domain.Logger
	audit      domain.AuditLogger
	taskQueue  domain.TaskQueue
	storage    domain.Storage
	reportGen  domain.AutoReportGenerator
	pdfService domain.PDFService
}

func NewClientService(repository domain.ClientRepository, taskQueue domain.TaskQueue, storage domain.Storage, reportGen domain.AutoReportGenerator, pdfService domain.PDFService, logger domain.Logger, audit domain.AuditLogger) domain.ClientService {
	return &ClientService{
		repository: repository,
		taskQueue:  taskQueue,
		storage:    storage,
		reportGen:  reportGen,
		pdfService: pdfService,
		logger:     logger,
		audit:      audit,
	}
}

func (s *ClientService) CreateClient(ctx context.Context, params domain.CreateClientParams) (*domain.Client, error) {
	client, err := s.repository.CreateClient(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.CreateClient", "failed to create client", err,
				zap.String("first_name", params.FirstName),
				zap.String("last_name", params.LastName),
			)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.CreateClient", "client created successfully",
			zap.String("client_id", client.ID.String()),
		)
	}

	return client, nil
}

func (s *ClientService) ListClients(ctx context.Context, params domain.ListClientsParams) (*domain.ClientPage, error) {
	page, err := s.repository.ListClients(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListClients", "failed to list clients", err)
		}
		return nil, err
	}

	if s.audit != nil && len(page.Items) > 0 {
		subjectIDs := make([]string, len(page.Items))
		for i, item := range page.Items {
			subjectIDs[i] = item.ID.String()
		}
		details := map[string]any{
			"limit":  params.Limit,
			"offset": params.Offset,
			"total":  page.TotalCount,
			"count":  len(page.Items),
		}
		if params.Search != nil && *params.Search != "" {
			details["search"] = *params.Search
		}
		if params.Status != nil && *params.Status != "" {
			details["status"] = *params.Status
		}
		if auditErr := s.audit.LogBatch(ctx, domain.AuditBatchEvent{
			EventGroupID: uuid.New().String(),
			EventType:    "record_access",
			Action:       "list",
			Result:       "success",
			SubjectType:  "client",
			SubjectIDs:   subjectIDs,
			AccessRule:   strPtr("CLIENT.VIEW"),
			Details:      details,
		}); auditErr != nil {
			s.logger.LogError(ctx, "ClientService.ListClients", "audit log failed", auditErr)
		}
	}

	return page, nil
}

func (s *ClientService) ListWaitingListClients(ctx context.Context, params domain.ListWaitingListClientsParams) (*domain.WaitingListClientPage, error) {
	if params.SortDays == "" {
		params.SortDays = "desc"
	}

	page, err := s.repository.ListWaitingListClients(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListWaitingListClients", "failed to list waiting list clients", err)
		}
		return nil, err
	}

	if s.audit != nil && len(page.Items) > 0 {
		subjectIDs := make([]string, len(page.Items))
		for i, item := range page.Items {
			subjectIDs[i] = item.ID.String()
		}
		details := map[string]any{
			"limit":     params.Limit,
			"offset":    params.Offset,
			"total":     page.TotalCount,
			"count":     len(page.Items),
			"sort_days": params.SortDays,
			"page_type": "waiting_list",
		}
		if params.Search != nil && *params.Search != "" {
			details["search"] = *params.Search
		}
		if auditErr := s.audit.LogBatch(ctx, domain.AuditBatchEvent{
			EventGroupID: uuid.New().String(),
			EventType:    "record_access",
			Action:       "list",
			Result:       "success",
			SubjectType:  "client",
			SubjectIDs:   subjectIDs,
			AccessRule:   strPtr("CLIENT.VIEW"),
			Details:      details,
		}); auditErr != nil {
			s.logger.LogError(ctx, "ClientService.ListWaitingListClients", "audit log failed", auditErr)
		}
	}

	return page, nil
}

func (s *ClientService) ListInCareClients(ctx context.Context, params domain.ListInCareClientsParams) (*domain.InCareClientPage, error) {
	if params.SortDaysInCare == "" {
		params.SortDaysInCare = "desc"
	}

	page, err := s.repository.ListInCareClients(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListInCareClients", "failed to list in-care clients", err)
		}
		return nil, err
	}

	if s.audit != nil && len(page.Items) > 0 {
		subjectIDs := make([]string, len(page.Items))
		for i, item := range page.Items {
			subjectIDs[i] = item.ID.String()
		}
		details := map[string]any{
			"limit":          params.Limit,
			"offset":         params.Offset,
			"total":          page.TotalCount,
			"count":          len(page.Items),
			"sort_days_care": params.SortDaysInCare,
			"page_type":      "in_care",
		}
		if auditErr := s.audit.LogBatch(ctx, domain.AuditBatchEvent{
			EventGroupID: uuid.New().String(),
			EventType:    "record_access",
			Action:       "list",
			Result:       "success",
			SubjectType:  "client",
			SubjectIDs:   subjectIDs,
			AccessRule:   strPtr("CLIENT.VIEW"),
			Details:      details,
		}); auditErr != nil {
			s.logger.LogError(ctx, "ClientService.ListInCareClients", "audit log failed", auditErr)
		}
	}

	return page, nil
}

func (s *ClientService) GetClientCounts(ctx context.Context) (*domain.ClientCounts, error) {
	counts, err := s.repository.GetClientCounts(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientCounts", "failed to get client counts", err)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.GetClientCounts", "client counts retrieved successfully")
	}

	if s.audit != nil {
		if auditErr := s.audit.Log(ctx, domain.AuditEvent{
			EventType:   "record_access",
			Action:      "read",
			Result:      "success",
			SubjectType: "client",
			SubjectID:   "summary",
			AccessRule:  strPtr("CLIENT.VIEW"),
			Details: map[string]any{
				"metric":              "client_counts",
				"total_clients":       counts.TotalClients,
				"clients_in_care":     counts.ClientsInCare,
				"clients_waiting":     counts.ClientsOnWaitingList,
				"clients_out_of_care": counts.ClientsOutOfCare,
			},
		}); auditErr != nil {
			s.logger.LogError(ctx, "ClientService.GetClientCounts", "audit log failed", auditErr)
		}
	}

	return counts, nil
}

func (s *ClientService) GetClientStatusCounts(ctx context.Context) (*domain.ClientStatusCounts, error) {
	counts, err := s.repository.GetClientStatusCounts(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientStatusCounts", "failed to get client status counts", err)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.GetClientStatusCounts", "client status counts retrieved successfully")
	}

	if s.audit != nil {
		if auditErr := s.audit.Log(ctx, domain.AuditEvent{
			EventType:   "record_access",
			Action:      "read",
			Result:      "success",
			SubjectType: "client",
			SubjectID:   "summary",
			AccessRule:  strPtr("CLIENT.VIEW"),
			Details: map[string]any{
				"metric":                       "client_status_counts",
				"in_or_scheduled_care":         counts.ClientsInOrScheduledInCare,
				"waiting_list":                 counts.ClientsOnWaitingList,
				"out_or_scheduled_out_of_care": counts.ClientsOutOrScheduledOutOfCare,
			},
		}); auditErr != nil {
			s.logger.LogError(ctx, "ClientService.GetClientStatusCounts", "audit log failed", auditErr)
		}
	}

	return counts, nil
}

func (s *ClientService) GetInCareStats(ctx context.Context) (*domain.InCareStats, error) {
	stats, err := s.repository.GetInCareStats(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetInCareStats", "failed to get in-care stats", err)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.GetInCareStats", "in-care stats retrieved successfully")
	}

	if s.audit != nil {
		if auditErr := s.audit.Log(ctx, domain.AuditEvent{
			EventType:   "record_access",
			Action:      "read",
			Result:      "success",
			SubjectType: "client",
			SubjectID:   "summary",
			AccessRule:  strPtr("CLIENT.VIEW"),
			Details: map[string]any{
				"metric":                "in_care_stats",
				"clients_in_care":       stats.ClientsInCare,
				"clients_scheduled":     stats.ClientsScheduledInCare,
				"contracts_ending_soon": stats.ContractsEndingSoon,
			},
		}); auditErr != nil {
			s.logger.LogError(ctx, "ClientService.GetInCareStats", "audit log failed", auditErr)
		}
	}

	return stats, nil
}

func (s *ClientService) GetWaitingListStats(ctx context.Context) (*domain.WaitingListStats, error) {
	stats, err := s.repository.GetWaitingListStats(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetWaitingListStats", "failed to get waiting list stats", err)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.GetWaitingListStats", "waiting list stats retrieved successfully")
	}

	if s.audit != nil {
		if auditErr := s.audit.Log(ctx, domain.AuditEvent{
			EventType:   "record_access",
			Action:      "read",
			Result:      "success",
			SubjectType: "client",
			SubjectID:   "summary",
			AccessRule:  strPtr("CLIENT.VIEW"),
			Details: map[string]any{
				"metric":            "waiting_list_stats",
				"total_clients":     stats.TotalClients,
				"total_crisis":      stats.TotalCrisis,
				"total_regular":     stats.TotalRegular,
				"avg_days_waitlist": stats.AvgDaysInWaitlist,
			},
		}); auditErr != nil {
			s.logger.LogError(ctx, "ClientService.GetWaitingListStats", "audit log failed", auditErr)
		}
	}

	return stats, nil
}

func (s *ClientService) GetClientByID(ctx context.Context, id uuid.UUID) (*domain.ClientPageDetail, error) {
	detail, err := s.repository.GetClientByID(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientByID", "failed to get client by id", err,
				zap.String("client_id", id.String()),
			)
		}
		return nil, err
	}

	if s.audit != nil {
		cid := id.String()
		if auditErr := s.audit.Log(ctx, domain.AuditEvent{
			EventType:   "record_access",
			Action:      "read",
			Result:      "success",
			SubjectType: "client",
			SubjectID:   cid,
			ClientID:    &cid,
			AccessRule:  strPtr("CLIENT.VIEW"),
		}); auditErr != nil {
			s.logger.LogError(ctx, "ClientService.GetClientByID", "audit log failed", auditErr)
		}
	}

	return detail, nil
}

func (s *ClientService) UpdateClient(ctx context.Context, id uuid.UUID, params domain.UpdateClientParams) (*domain.Client, error) {
	client, err := s.repository.UpdateClient(ctx, id, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.UpdateClient", "failed to update client", err,
				zap.String("client_id", id.String()),
			)
		}
		return nil, err
	}
	return client, nil
}

func (s *ClientService) GetClientAddresses(ctx context.Context, id uuid.UUID) ([]domain.ClientAddress, error) {
	addresses, err := s.repository.GetClientAddresses(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientAddresses", "failed to get client addresses", err,
				zap.String("client_id", id.String()),
			)
		}
		return nil, err
	}
	return addresses, nil
}

func (s *ClientService) UpdateClientStatus(ctx context.Context, clientID uuid.UUID, params domain.UpdateClientStatusParams) (*domain.UpdateClientStatusResult, error) {
	if params.IsScheduled {
		return s.handleScheduledStatusUpdate(ctx, clientID, params)
	}
	return s.handleNormalStatusUpdate(ctx, clientID, params)
}

func (s *ClientService) handleNormalStatusUpdate(ctx context.Context, clientID uuid.UUID, params domain.UpdateClientStatusParams) (*domain.UpdateClientStatusResult, error) {
	result, err := s.repository.UpdateClientStatus(ctx, clientID, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.UpdateClientStatus", "failed to update client status", err, zap.String("client_id", clientID.String()))
		}
		return nil, err
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.UpdateClientStatus", "client status updated successfully", zap.String("client_id", clientID.String()), zap.String("new_status", result.Status))
	}
	return result, nil
}

func (s *ClientService) handleScheduledStatusUpdate(ctx context.Context, clientID uuid.UUID, params domain.UpdateClientStatusParams) (*domain.UpdateClientStatusResult, error) {
	return &domain.UpdateClientStatusResult{ID: clientID}, nil
}

func (s *ClientService) PutClientInCare(ctx context.Context, clientID uuid.UUID, params domain.PutClientInCareParams) (*domain.PutClientInCareResult, error) {
	careStartDate, err := time.Parse("2006-01-02", params.CareStartDate)
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
		msg := "care_start_date is in the past; client was activated immediately"
		warning = &msg
	}

	result, err := s.repository.PutClientInCare(ctx, clientID, targetStatus, careStartDay, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.PutClientInCare", "failed to put client in care", err, zap.String("client_id", clientID.String()))
		}
		return nil, err
	}
	result.Warning = warning

	// TODO: enqueue coordinator notification via taskQueue (needs user ID from UpsertMainCoordinator)

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.PutClientInCare", "client moved into care lifecycle", zap.String("client_id", clientID.String()), zap.String("status", result.Status))
	}

	return result, nil
}

func (s *ClientService) PutClientOutOfCare(ctx context.Context, clientID uuid.UUID, params domain.PutClientOutOfCareParams) (*domain.PutClientOutOfCareResult, error) {
	dischargeDate, err := time.Parse("2006-01-02", params.DischargeDate)
	if err != nil {
		return nil, fmt.Errorf("discharge_date must be in YYYY-MM-DD format")
	}
	dischargeDay := time.Date(dischargeDate.Year(), dischargeDate.Month(), dischargeDate.Day(), 0, 0, 0, 0, time.UTC)

	var finalEvaluation *string
	if params.FinalEvaluation != nil {
		trimmed := strings.TrimSpace(*params.FinalEvaluation)
		if trimmed == "" {
			return nil, fmt.Errorf("final_evaluation cannot be empty")
		}
		finalEvaluation = &trimmed
	}

	dischargeReason, err := parseDischargeReason(params.DischargeReason)
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

	result, err := s.repository.PutClientOutOfCare(ctx, clientID, targetStatus, dischargeDay, dischargeReason, finalEvaluation, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.PutClientOutOfCare", "failed to put client out of care", err, zap.String("client_id", clientID.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.PutClientOutOfCare", "client discharge status updated", zap.String("client_id", clientID.String()), zap.String("status", result.Status))
	}

	return result, nil
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

func (s *ClientService) ListStatusHistory(ctx context.Context, clientID uuid.UUID) ([]domain.ClientStatusHistory, error) {
	histories, err := s.repository.ListStatusHistory(ctx, domain.ListStatusHistoryParams{
		ClientID: clientID,
		Limit:    10,
		Offset:   0,
	})
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListStatusHistory", "failed to list status history", err, zap.String("client_id", clientID.String()))
		}
		return nil, err
	}
	return histories, nil
}

func (s *ClientService) logGroupEAudit(ctx context.Context, action, subjectType string, subjectID, clientID uuid.UUID, permission string, count int) {
	if s.audit == nil {
		return
	}
	eventType := "record_change"
	if action == "read" || action == "list" {
		eventType = "record_access"
	}
	var details map[string]any
	if count >= 0 {
		details = map[string]any{"count": count}
	}
	clientIDString := clientID.String()
	if err := s.audit.Log(ctx, domain.AuditEvent{
		EventType:   eventType,
		Action:      action,
		Result:      "success",
		SubjectType: subjectType,
		SubjectID:   subjectID.String(),
		ClientID:    &clientIDString,
		AccessRule:  strPtr(permission),
		Details:     details,
	}); err != nil && s.logger != nil {
		s.logger.LogError(ctx, "ClientService.logGroupEAudit", "Group E audit log failed", err,
			zap.String("client_id", clientIDString),
			zap.String("subject_type", subjectType),
			zap.String("action", action),
		)
	}
}

func (s *ClientService) AddClientDocument(ctx context.Context, clientID uuid.UUID, params domain.AddClientDocumentParams) ([]domain.AddClientDocumentResult, error) {
	if len(params.Documents) == 0 {
		return nil, fmt.Errorf("at least one document is required")
	}

	seen := make(map[uuid.UUID]struct{}, len(params.Documents))
	attachmentIDs := make([]uuid.UUID, 0, len(params.Documents))
	for _, doc := range params.Documents {
		if _, ok := seen[doc.AttachmentID]; ok {
			return nil, fmt.Errorf("duplicate attachment_id %s in request", doc.AttachmentID)
		}
		seen[doc.AttachmentID] = struct{}{}
		attachmentIDs = append(attachmentIDs, doc.AttachmentID)
	}

	attachments, err := s.repository.GetAttachmentsByUUIDs(ctx, attachmentIDs)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.AddClientDocument", "failed to fetch attachments", err, zap.String("client_id", clientID.String()))
		}
		return nil, fmt.Errorf("failed to add client document")
	}

	if len(attachments) != len(params.Documents) {
		return nil, fmt.Errorf("one or more attachments were not found")
	}

	attachmentsByID := make(map[uuid.UUID]domain.AttachmentFile, len(attachments))
	objectKeys := make([]string, 0, len(attachments))
	for _, a := range attachments {
		attachmentsByID[a.UUID] = a
		objectKeys = append(objectKeys, a.File)
	}

	fileSizes, err := s.storage.GetFileInfos(ctx, objectKeys)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.AddClientDocument", "failed to verify uploaded files", err, zap.String("client_id", clientID.String()))
		}
		return nil, fmt.Errorf("failed to verify uploaded file")
	}

	for _, doc := range params.Documents {
		attachment, ok := attachmentsByID[doc.AttachmentID]
		if !ok {
			return nil, fmt.Errorf("attachment %s was not found", doc.AttachmentID)
		}

		size, ok := fileSizes[attachment.File]
		if !ok {
			return nil, fmt.Errorf("uploaded file is missing")
		}
		if size <= 0 {
			return nil, fmt.Errorf("uploaded file is empty")
		}
		if size > maxFileSize {
			return nil, fmt.Errorf("file size exceeds maximum limit of %dMB", maxFileSize>>20)
		}
	}

	result, err := s.repository.AddClientDocuments(ctx, clientID, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.AddClientDocument", "failed to add client documents", err, zap.String("client_id", clientID.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.AddClientDocument", "client documents added successfully",
			zap.String("client_id", clientID.String()),
			zap.Int("count", len(result)),
		)
	}
	s.logGroupEAudit(ctx, "upload", "client_document", clientID, clientID, domain.PermClientDocumentsUpload.String(), len(result))

	return result, nil
}

func (s *ClientService) ListClientDocuments(ctx context.Context, params domain.ListClientDocumentsParams) (*domain.ListClientDocumentsResult, error) {
	result, err := s.repository.ListClientDocuments(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListClientDocuments", "failed to list client documents", err, zap.String("client_id", params.ClientID.String()))
		}
		return nil, err
	}

	for i := range result.Documents {
		if result.Documents[i].File != "" {
			url, err := s.storage.GeneratePresignedURL(ctx, result.Documents[i].File, 15*time.Minute)
			if err == nil {
				result.Documents[i].FileURL = url
			}
		}
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.ListClientDocuments", "client documents listed successfully",
			zap.String("client_id", params.ClientID.String()),
			zap.Int("count", len(result.Documents)),
		)
	}
	s.logGroupEAudit(ctx, "list", "client_document", params.ClientID, params.ClientID, domain.PermClientDocumentsView.String(), len(result.Documents))

	return result, nil
}

func (s *ClientService) DeleteClientDocument(ctx context.Context, clientID uuid.UUID, documentID uuid.UUID) (*domain.DeleteClientDocumentResult, error) {
	result, err := s.repository.DeleteClientDocument(ctx, clientID, documentID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.DeleteClientDocument", "failed to delete client document", err,
				zap.String("client_id", clientID.String()),
				zap.String("document_id", documentID.String()),
			)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.DeleteClientDocument", "client document deleted successfully",
			zap.String("client_id", clientID.String()),
			zap.String("document_id", documentID.String()),
		)
	}
	s.logGroupEAudit(ctx, "delete", "client_document", documentID, clientID, domain.PermClientDocumentsDelete.String(), -1)

	return result, nil
}

func (s *ClientService) GetMissingClientDocuments(ctx context.Context, clientID uuid.UUID) ([]string, error) {
	missingDocs, err := s.repository.GetMissingClientDocuments(ctx, clientID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetMissingClientDocuments", "failed to get missing client documents", err, zap.String("client_id", clientID.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.GetMissingClientDocuments", "missing client documents retrieved successfully",
			zap.String("client_id", clientID.String()),
		)
	}
	s.logGroupEAudit(ctx, "list", "client_document_requirement", clientID, clientID, domain.PermClientDocumentsView.String(), len(missingDocs))

	return missingDocs, nil
}

func (s *ClientService) CreateClientGoal(ctx context.Context, clientID uuid.UUID, params domain.CreateClientGoalParams) (*domain.ClientGoal, error) {
	if strings.TrimSpace(params.Title) == "" {
		return nil, domain.ErrClientGoalTitleRequired
	}

	goal, err := s.repository.CreateClientGoal(ctx, clientID, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.CreateClientGoal", "failed to create client goal", err, zap.String("client_id", clientID.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.CreateClientGoal", "client goal created successfully",
			zap.String("client_id", clientID.String()),
			zap.String("goal_id", goal.ID.String()),
		)
	}
	s.logGroupEAudit(ctx, "create", "client_goal", goal.ID, clientID, domain.PermClientCarePlanCreate.String(), -1)

	return goal, nil
}

func (s *ClientService) UpdateClientGoal(ctx context.Context, clientID uuid.UUID, goalID uuid.UUID, params domain.UpdateClientGoalParams) (*domain.UpdateClientGoalResult, error) {
	if params.Title == nil && params.Description == nil && params.Priority == nil && params.TopicID == nil && params.SortOrder == nil {
		return nil, domain.ErrClientGoalEmptyPatch
	}

	result, err := s.repository.UpdateClientGoal(ctx, clientID, goalID, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.UpdateClientGoal", "failed to update client goal", err,
				zap.String("client_id", clientID.String()),
				zap.String("goal_id", goalID.String()),
			)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.UpdateClientGoal", "client goal updated successfully",
			zap.String("client_id", clientID.String()),
			zap.String("goal_id", goalID.String()),
		)
	}
	s.logGroupEAudit(ctx, "update", "client_goal", result.Goal.ID, clientID, domain.PermClientCarePlanUpdate.String(), -1)

	return result, nil
}

func (s *ClientService) GetClientGoalsForEvaluationPage(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID) (*domain.ClientGoalsForEvaluationPage, error) {
	result, err := s.repository.GetClientGoalsForEvaluationPage(ctx, clientID, employeeID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientGoalsForEvaluationPage", "failed to get goals for evaluation page", err, zap.String("client_id", clientID.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.GetClientGoalsForEvaluationPage", "goals for evaluation page retrieved successfully",
			zap.String("client_id", clientID.String()),
		)
	}
	s.logGroupEAudit(ctx, "list", "client_goal", clientID, clientID, domain.PermClientEvaluationView.String(), len(result.Goals))

	return result, nil
}

func (s *ClientService) CreateGoalEvaluation(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, params domain.CreateGoalEvaluationParams) (*domain.GoalEvaluation, error) {
	if err := validateGoalEvaluationItems(params.Items); err != nil {
		return nil, err
	}

	result, err := s.repository.CreateGoalEvaluation(ctx, clientID, employeeID, params)
	if err != nil {
		if result != nil {
			s.logGroupEAudit(ctx, "save", "client_goal_evaluation", result.ID, result.ClientID, domain.PermClientEvaluationCreate.String(), len(result.Items))
		}
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.CreateGoalEvaluation", "failed to create goal evaluation", err, zap.String("client_id", clientID.String()))
		}
		return result, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.CreateGoalEvaluation", "goal evaluation created successfully",
			zap.String("client_id", clientID.String()),
			zap.String("evaluation_id", result.ID.String()),
			zap.Bool("submitted", params.Submit),
		)
	}
	action := "save"
	if params.Submit {
		action = "submit"
	}
	s.logGroupEAudit(ctx, action, "client_goal_evaluation", result.ID, clientID, domain.PermClientEvaluationCreate.String(), len(result.Items))

	return result, nil
}

func (s *ClientService) UpdateGoalEvaluationDraft(ctx context.Context, evaluationID uuid.UUID, employeeID uuid.UUID, params domain.UpdateGoalEvaluationDraftParams) (*domain.GoalEvaluation, error) {
	if err := validateGoalEvaluationItems(params.Items); err != nil {
		return nil, err
	}

	result, err := s.repository.UpdateGoalEvaluationDraft(ctx, evaluationID, employeeID, params)
	if err != nil {
		return result, err
	}
	s.logGroupEAudit(ctx, "save", "client_goal_evaluation", result.ID, result.ClientID, domain.PermClientEvaluationCreate.String(), len(result.Items))
	return result, nil
}

func (s *ClientService) SubmitGoalEvaluationDraft(ctx context.Context, evaluationID uuid.UUID, employeeID uuid.UUID, params domain.SubmitGoalEvaluationDraftParams) (*domain.GoalEvaluation, error) {
	result, err := s.repository.SubmitGoalEvaluationDraft(ctx, evaluationID, employeeID, params)
	if err != nil {
		if result != nil {
			s.logGroupEAudit(ctx, "save", "client_goal_evaluation", result.ID, result.ClientID, domain.PermClientEvaluationCreate.String(), len(result.Items))
		}
		return result, err
	}
	s.logGroupEAudit(ctx, "submit", "client_goal_evaluation", result.ID, result.ClientID, domain.PermClientEvaluationCreate.String(), len(result.Items))
	return result, nil
}

func validateGoalEvaluationItems(items []domain.GoalEvaluationItemParams) error {
	itemsByGoal := make(map[uuid.UUID]struct{}, len(items))
	for _, item := range items {
		if _, exists := itemsByGoal[item.GoalID]; exists {
			return fmt.Errorf("%w: %s", domain.ErrGoalEvaluationDuplicateGoal, item.GoalID)
		}
		if err := validateProgress(item.Progress); err != nil {
			return err
		}
		itemsByGoal[item.GoalID] = struct{}{}
	}
	return nil
}

func (s *ClientService) GetGoalEvaluationBootstrap(ctx context.Context, clientID uuid.UUID) (*domain.GoalEvaluationBootstrap, error) {
	result, err := s.repository.GetGoalEvaluationBootstrap(ctx, clientID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetGoalEvaluationBootstrap", "failed to get goal evaluation bootstrap", err, zap.String("client_id", clientID.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.GetGoalEvaluationBootstrap", "goal evaluation bootstrap retrieved successfully",
			zap.String("client_id", clientID.String()),
		)
	}
	s.logGroupEAudit(ctx, "read", "client_goal_evaluation", clientID, clientID, domain.PermClientEvaluationView.String(), len(result.ActiveGoals))

	return result, nil
}

func (s *ClientService) ListClientSubmittedEvaluations(ctx context.Context, params domain.ListClientSubmittedEvaluationsParams) (*domain.ListClientSubmittedEvaluationsResult, error) {
	result, err := s.repository.ListClientSubmittedEvaluations(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListClientSubmittedEvaluations", "failed to list submitted evaluations", err, zap.String("client_id", params.ClientID.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.ListClientSubmittedEvaluations", "submitted evaluations listed successfully",
			zap.String("client_id", params.ClientID.String()),
			zap.Int("count", len(result.Items)),
		)
	}
	s.logGroupEAudit(ctx, "list", "client_goal_evaluation", params.ClientID, params.ClientID, domain.PermClientEvaluationView.String(), len(result.Items))

	return result, nil
}

func (s *ClientService) ListGoalEvaluationHistory(ctx context.Context, params domain.ListGoalEvaluationHistoryParams) (*domain.ListGoalEvaluationHistoryResult, error) {
	result, err := s.repository.ListGoalEvaluationHistory(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListGoalEvaluationHistory", "failed to list goal evaluation history", err,
				zap.String("client_id", params.ClientID.String()),
				zap.String("goal_id", params.GoalID.String()),
			)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.ListGoalEvaluationHistory", "goal evaluation history listed successfully",
			zap.String("client_id", params.ClientID.String()),
			zap.String("goal_id", params.GoalID.String()),
			zap.Int("count", len(result.Items)),
		)
	}
	s.logGroupEAudit(ctx, "list", "client_goal_evaluation", params.GoalID, params.ClientID, domain.PermClientEvaluationView.String(), len(result.Items))

	return result, nil
}

func validateProgress(value string) error {
	progress := strings.TrimSpace(value)
	if progress == "" {
		return nil
	}
	allowed := map[string]struct{}{
		"no_progress":      {},
		"regression":       {},
		"limited_progress": {},
		"good_progress":    {},
		"achieved":         {},
		"blocked":          {},
	}
	if _, ok := allowed[progress]; !ok {
		return fmt.Errorf("%w: %s", domain.ErrGoalEvaluationInvalidProgress, value)
	}
	return nil
}

func (s *ClientService) RequestLocationTransfer(ctx context.Context, clientID uuid.UUID, params domain.CreateLocationTransferParams) error {
	err := s.repository.CreateLocationTransfer(ctx, clientID, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.RequestLocationTransfer", "failed to create location transfer request", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return fmt.Errorf("failed to create location transfer request")
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.RequestLocationTransfer", "location transfer request created successfully",
			zap.String("client_id", clientID.String()),
		)
	}

	return nil
}

func (s *ClientService) ApproveLocationTransfer(ctx context.Context, employeeID uuid.UUID, params domain.ApproveLocationTransferParams) error {
	status := strings.TrimSpace(strings.ToLower(params.Status))
	if status != "approved" && status != "rejected" {
		return fmt.Errorf("invalid status: must be approved or rejected")
	}

	err := s.repository.ApproveLocationTransfer(ctx, employeeID, domain.ApproveLocationTransferParams{
		TransferID: params.TransferID,
		Status:     status,
	})
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ApproveLocationTransfer", "failed to approve location transfer request", err,
				zap.String("transfer_id", params.TransferID.String()),
				zap.String("status", status),
			)
		}
		return fmt.Errorf("failed to approve location transfer request")
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.ApproveLocationTransfer", "location transfer request processed successfully",
			zap.String("transfer_id", params.TransferID.String()),
			zap.String("status", status),
		)
	}

	return nil
}

func (s *ClientService) ListLocationTransferRequests(ctx context.Context, params domain.ListLocationTransferParams) (*domain.ListLocationTransferResult, error) {
	result, err := s.repository.ListLocationTransferRequests(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListLocationTransferRequests", "failed to list location transfer requests", err)
		}
		return nil, fmt.Errorf("failed to list location transfer requests")
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.ListLocationTransferRequests", "location transfer requests listed successfully",
			zap.Int("count", len(result.Items)),
		)
	}
	return result, nil
}

func (s *ClientService) GetGoalEvaluation(ctx context.Context, evaluationID uuid.UUID) (*domain.GoalEvaluation, error) {
	result, err := s.repository.GetGoalEvaluation(ctx, evaluationID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetGoalEvaluation", "failed to get goal evaluation", err,
				zap.String("evaluation_id", evaluationID.String()),
			)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.GetGoalEvaluation", "goal evaluation fetched successfully",
			zap.String("evaluation_id", evaluationID.String()),
		)
	}
	s.logGroupEAudit(ctx, "read", "client_goal_evaluation", result.ID, result.ClientID, domain.PermClientEvaluationView.String(), len(result.Items))

	return result, nil
}

func (s *ClientService) ListUpcomingEvaluations(ctx context.Context, params domain.ListUpcomingEvaluationsParams) (*domain.ListUpcomingEvaluationsResult, error) {
	result, err := s.repository.ListUpcomingEvaluations(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListUpcomingEvaluations", "failed to list upcoming evaluations", err,
				zap.String("employee_id", params.EmployeeID.String()),
			)
		}
		return nil, fmt.Errorf("failed to list upcoming evaluations")
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.ListUpcomingEvaluations", "upcoming evaluations listed successfully",
			zap.String("employee_id", params.EmployeeID.String()),
			zap.Int("count", len(result.Items)),
		)
	}
	for _, item := range result.Items {
		s.logGroupEAudit(ctx, "read", "client_goal_evaluation_schedule", item.ClientID, item.ClientID, domain.PermClientEvaluationView.String(), -1)
	}

	return result, nil
}

func (s *ClientService) ListRecentSubmittedEvaluations(ctx context.Context, params domain.ListRecentSubmittedEvaluationsParams) (*domain.ListRecentSubmittedEvaluationsResult, error) {
	result, err := s.repository.ListRecentSubmittedEvaluations(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListRecentSubmittedEvaluations", "failed to list recent submitted evaluations", err,
				zap.String("employee_id", params.EmployeeID.String()),
			)
		}
		return nil, fmt.Errorf("failed to list recent submitted evaluations")
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.ListRecentSubmittedEvaluations", "recent submitted evaluations listed successfully",
			zap.String("employee_id", params.EmployeeID.String()),
			zap.Int("count", len(result.Items)),
		)
	}
	for _, item := range result.Items {
		s.logGroupEAudit(ctx, "read", "client_goal_evaluation", item.EvaluationID, item.ClientID, domain.PermClientEvaluationView.String(), -1)
	}

	return result, nil
}

func (s *ClientService) ListRecentDraftEvaluations(ctx context.Context, params domain.ListRecentDraftEvaluationsParams) (*domain.ListRecentDraftEvaluationsResult, error) {
	result, err := s.repository.ListRecentDraftEvaluations(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListRecentDraftEvaluations", "failed to list recent draft evaluations", err,
				zap.String("employee_id", params.EmployeeID.String()),
			)
		}
		return nil, fmt.Errorf("failed to list recent draft evaluations")
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ClientService.ListRecentDraftEvaluations", "recent draft evaluations listed successfully",
			zap.String("employee_id", params.EmployeeID.String()),
			zap.Int("count", len(result.Items)),
		)
	}
	for _, item := range result.Items {
		s.logGroupEAudit(ctx, "read", "client_goal_evaluation", item.EvaluationID, item.ClientID, domain.PermClientEvaluationView.String(), -1)
	}

	return result, nil
}

func (s *ClientService) GetEvaluationStats(ctx context.Context, employeeID uuid.UUID) (*domain.EvaluationStats, error) {
	result, err := s.repository.GetEvaluationStats(ctx, employeeID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetEvaluationStats", "failed to get evaluation stats", err,
				zap.String("employee_id", employeeID.String()),
			)
		}
		return nil, fmt.Errorf("failed to get evaluation stats")
	}
	return result, nil
}

// =====================
// Medical - Diagnoses
// =====================

func (s *ClientService) logMedicalAudit(ctx context.Context, action, subjectType string, subjectID, clientID uuid.UUID, permission string, count int) {
	if s.audit == nil {
		return
	}
	eventType := "record_change"
	if action == "read" || action == "list" {
		eventType = "record_access"
	}
	var details map[string]any
	if count >= 0 {
		details = map[string]any{"count": count}
	}
	clientIDString := clientID.String()
	if err := s.audit.Log(ctx, domain.AuditEvent{
		EventType:   eventType,
		Action:      action,
		Result:      "success",
		SubjectType: subjectType,
		SubjectID:   subjectID.String(),
		ClientID:    &clientIDString,
		AccessRule:  strPtr(permission),
		Details:     details,
	}); err != nil && s.logger != nil {
		s.logger.LogError(ctx, "ClientService.logMedicalAudit", "medical audit log failed", err,
			zap.String("client_id", clientIDString),
			zap.String("subject_type", subjectType),
			zap.String("action", action),
		)
	}
}

func (s *ClientService) CreateClientDiagnosis(ctx context.Context, clientID uuid.UUID, _ uuid.UUID, params domain.CreateClientDiagnosisParams) (*domain.ClientDiagnosis, error) {
	params.ClientID = clientID

	result, err := s.repository.CreateClientDiagnosis(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.CreateClientDiagnosis", "failed to create client diagnosis", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, err
	}

	s.logMedicalAudit(ctx, "create", "client_diagnosis", result.ID, clientID, "CLIENT.DIAGNOSIS.CREATE", -1)
	return result, nil
}

func (s *ClientService) ListClientDiagnoses(ctx context.Context, params domain.ListClientDiagnosesParams) (*domain.ListClientDiagnosesResult, error) {
	result, err := s.repository.ListClientDiagnoses(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListClientDiagnoses", "failed to list client diagnoses", err,
				zap.String("client_id", params.ClientID.String()),
			)
		}
		return nil, err
	}

	s.logMedicalAudit(ctx, "list", "client_diagnosis", params.ClientID, params.ClientID, "CLIENT.DIAGNOSIS.VIEW", len(result.Items))
	return result, nil
}

func (s *ClientService) GetClientDiagnosis(ctx context.Context, clientID, diagnosisID uuid.UUID) (*domain.ClientDiagnosis, error) {
	result, err := s.repository.GetClientDiagnosis(ctx, clientID, diagnosisID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientDiagnosis", "failed to get client diagnosis", err,
				zap.String("client_id", clientID.String()),
				zap.String("diagnosis_id", diagnosisID.String()),
			)
		}
		return nil, err
	}

	s.logMedicalAudit(ctx, "read", "client_diagnosis", diagnosisID, clientID, "CLIENT.DIAGNOSIS.VIEW", -1)
	return result, nil
}

func (s *ClientService) UpdateClientDiagnosis(ctx context.Context, clientID, diagnosisID, _ uuid.UUID, params domain.UpdateClientDiagnosisParams) (*domain.ClientDiagnosis, error) {
	params.ClientID = clientID
	params.ID = diagnosisID

	result, err := s.repository.UpdateClientDiagnosis(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.UpdateClientDiagnosis", "failed to update client diagnosis", err,
				zap.String("client_id", clientID.String()),
				zap.String("diagnosis_id", diagnosisID.String()),
			)
		}
		return nil, err
	}

	s.logMedicalAudit(ctx, "update", "client_diagnosis", diagnosisID, clientID, "CLIENT.DIAGNOSIS.UPDATE", -1)
	return result, nil
}

func (s *ClientService) DeleteClientDiagnosis(ctx context.Context, clientID, diagnosisID uuid.UUID) (*domain.DeleteClientDiagnosisResult, error) {
	result, err := s.repository.DeleteClientDiagnosis(ctx, clientID, diagnosisID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.DeleteClientDiagnosis", "failed to delete client diagnosis", err,
				zap.String("client_id", clientID.String()),
				zap.String("diagnosis_id", diagnosisID.String()),
			)
		}
		return nil, err
	}

	s.logMedicalAudit(ctx, "delete", "client_diagnosis", diagnosisID, clientID, "CLIENT.DIAGNOSIS.DELETE", -1)
	return result, nil
}

// =====================
// Medical - Medication Orders
// =====================

func (s *ClientService) CreateClientMedicationOrder(ctx context.Context, clientID uuid.UUID, _ uuid.UUID, params domain.CreateClientMedicationOrderParams) (*domain.ClientMedicationOrder, error) {
	params.ClientID = clientID

	result, err := s.repository.CreateClientMedicationOrder(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.CreateClientMedicationOrder", "failed to create client medication order", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, err
	}

	s.logMedicalAudit(ctx, "create", "client_medication_order", result.ID, clientID, "CLIENT.MEDICATION.CREATE", -1)
	return result, nil
}

func (s *ClientService) ListClientMedicationOrders(ctx context.Context, params domain.ListClientMedicationOrdersParams) (*domain.ListClientMedicationOrdersResult, error) {
	result, err := s.repository.ListClientMedicationOrders(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListClientMedicationOrders", "failed to list client medication orders", err,
				zap.String("client_id", params.ClientID.String()),
			)
		}
		return nil, err
	}

	s.logMedicalAudit(ctx, "list", "client_medication_order", params.ClientID, params.ClientID, "CLIENT.MEDICATION.VIEW", len(result.Items))
	return result, nil
}

func (s *ClientService) GetClientMedicationOrder(ctx context.Context, clientID, orderID uuid.UUID) (*domain.ClientMedicationOrder, error) {
	result, err := s.repository.GetClientMedicationOrder(ctx, clientID, orderID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientMedicationOrder", "failed to get client medication order", err,
				zap.String("client_id", clientID.String()),
				zap.String("order_id", orderID.String()),
			)
		}
		return nil, err
	}

	s.logMedicalAudit(ctx, "read", "client_medication_order", orderID, clientID, "CLIENT.MEDICATION.VIEW", -1)
	return result, nil
}

func (s *ClientService) UpdateClientMedicationOrder(ctx context.Context, clientID, orderID, _ uuid.UUID, params domain.UpdateClientMedicationOrderParams) (*domain.ClientMedicationOrder, error) {
	params.ClientID = clientID
	params.ID = orderID

	updated, err := s.repository.UpdateClientMedicationOrder(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.UpdateClientMedicationOrder", "failed to update client medication order", err,
				zap.String("client_id", clientID.String()),
				zap.String("order_id", orderID.String()),
			)
		}
		return nil, err
	}

	// Re-fetch enriched view for consistent response
	result, err := s.repository.GetClientMedicationOrder(ctx, clientID, updated.ID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.UpdateClientMedicationOrder", "failed to re-fetch updated medication order", err,
				zap.String("client_id", clientID.String()),
				zap.String("order_id", orderID.String()),
			)
		}
		return nil, err
	}

	s.logMedicalAudit(ctx, "update", "client_medication_order", orderID, clientID, "CLIENT.MEDICATION.UPDATE", -1)
	return result, nil
}

func (s *ClientService) DeleteClientMedicationOrder(ctx context.Context, clientID, orderID uuid.UUID) (*domain.DeleteClientMedicationOrderResult, error) {
	result, err := s.repository.DeleteClientMedicationOrder(ctx, clientID, orderID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.DeleteClientMedicationOrder", "failed to delete client medication order", err,
				zap.String("client_id", clientID.String()),
				zap.String("order_id", orderID.String()),
			)
		}
		return nil, err
	}

	s.logMedicalAudit(ctx, "delete", "client_medication_order", orderID, clientID, "CLIENT.MEDICATION.DELETE", -1)
	return result, nil
}

func (s *ClientService) GetClientMedicalOverview(ctx context.Context, clientID uuid.UUID) (*domain.ClientMedicalOverview, error) {
	result, err := s.repository.GetClientMedicalOverview(ctx, clientID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientMedicalOverview", "failed to get client medical overview", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, err
	}
	s.logMedicalAudit(ctx, "list", "client_diagnosis", clientID, clientID, "CLIENT.DIAGNOSIS.VIEW", len(result.Diagnoses))
	s.logMedicalAudit(ctx, "list", "client_medication_order", clientID, clientID, "CLIENT.MEDICATION.VIEW", len(result.MedicationOrders))
	return result, nil
}

// =====================
// Network - Sender
// =====================

func (s *ClientService) GetClientSender(ctx context.Context, clientID uuid.UUID) (*domain.Sender, error) {
	result, err := s.repository.GetClientSender(ctx, clientID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientSender", "failed to get client sender", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

// =====================
// Network - Emergency Contacts
// =====================

func (s *ClientService) CreateClientEmergencyContact(ctx context.Context, clientID uuid.UUID, params domain.CreateClientEmergencyContactParams) (*domain.ClientEmergencyContact, error) {
	params.ClientID = clientID
	result, err := s.repository.CreateClientEmergencyContact(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.CreateClientEmergencyContact", "failed to create client emergency contact", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) ListClientEmergencyContacts(ctx context.Context, params domain.ListClientEmergencyContactsParams) (*domain.ListClientEmergencyContactsResult, error) {
	result, err := s.repository.ListClientEmergencyContacts(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListClientEmergencyContacts", "failed to list client emergency contacts", err,
				zap.String("client_id", params.ClientID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) GetClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*domain.ClientEmergencyContact, error) {
	result, err := s.repository.GetClientEmergencyContact(ctx, contactID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientEmergencyContact", "failed to get client emergency contact", err,
				zap.String("contact_id", contactID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) UpdateClientEmergencyContact(ctx context.Context, contactID uuid.UUID, params domain.UpdateClientEmergencyContactParams) (*domain.ClientEmergencyContact, error) {
	params.ID = contactID
	result, err := s.repository.UpdateClientEmergencyContact(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.UpdateClientEmergencyContact", "failed to update client emergency contact", err,
				zap.String("contact_id", contactID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) DeleteClientEmergencyContact(ctx context.Context, contactID uuid.UUID) (*domain.DeleteClientEmergencyContactResult, error) {
	result, err := s.repository.DeleteClientEmergencyContact(ctx, contactID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.DeleteClientEmergencyContact", "failed to delete client emergency contact", err,
				zap.String("contact_id", contactID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

// =====================
// Network - Assigned Employees
// =====================

func (s *ClientService) CreateAssignedEmployee(ctx context.Context, clientID uuid.UUID, params domain.CreateAssignedEmployeeParams) (*domain.AssignedEmployee, error) {
	params.ClientID = clientID
	result, err := s.repository.CreateAssignedEmployee(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.CreateAssignedEmployee", "failed to create assigned employee", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, err
	}

	// Enqueue notification task (non-fatal on failure)
	err = s.taskQueue.EnqueueNotificationTask(ctx, domain.NotificationTaskPayload{
		RecipientUserIDs: []uuid.UUID{result.UserID},
		Type:             "new_client_assignment",
		Data: domain.NotificationTaskData{
			NewClientAssignment: &domain.NewClientAssignmentTaskData{
				ClientID:        result.ClientID,
				ClientFirstName: result.ClientFirstName,
				ClientLastName:  result.ClientLastName,
				ClientLocation:  result.ClientLocationName,
			},
		},
		CreatedAt: time.Now(),
		Message:   "You have been assigned to a new client",
	}, nil)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.CreateAssignedEmployee", "failed to enqueue notification task", err,
				zap.String("client_id", clientID.String()),
				zap.String("assignment_id", result.ID.String()),
			)
		}
	}

	return result, nil
}

func (s *ClientService) ListAssignedEmployees(ctx context.Context, params domain.ListAssignedEmployeesParams) (*domain.ListAssignedEmployeesResult, error) {
	result, err := s.repository.ListAssignedEmployees(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListAssignedEmployees", "failed to list assigned employees", err,
				zap.String("client_id", params.ClientID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) GetAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*domain.AssignedEmployee, error) {
	result, err := s.repository.GetAssignedEmployee(ctx, assignmentID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetAssignedEmployee", "failed to get assigned employee", err,
				zap.String("assignment_id", assignmentID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) UpdateAssignedEmployee(ctx context.Context, assignmentID uuid.UUID, params domain.UpdateAssignedEmployeeParams) (*domain.AssignedEmployee, error) {
	params.ID = assignmentID
	result, err := s.repository.UpdateAssignedEmployee(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.UpdateAssignedEmployee", "failed to update assigned employee", err,
				zap.String("assignment_id", assignmentID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) DeleteAssignedEmployee(ctx context.Context, assignmentID uuid.UUID) (*domain.DeleteAssignedEmployeeResult, error) {
	result, err := s.repository.DeleteAssignedEmployee(ctx, assignmentID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.DeleteAssignedEmployee", "failed to delete assigned employee", err,
				zap.String("assignment_id", assignmentID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

// =====================
// Network - Related Emails
// =====================

func (s *ClientService) GetClientRelatedEmails(ctx context.Context, clientID uuid.UUID) (*domain.ClientRelatedEmails, error) {
	result, err := s.repository.GetClientRelatedEmails(ctx, clientID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetClientRelatedEmails", "failed to get client related emails", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

// =====================
// Progress Reports
// =====================

func (s *ClientService) CreateProgressReport(ctx context.Context, clientID uuid.UUID, params domain.CreateProgressReportParams) (*domain.ProgressReport, error) {
	params.ClientID = clientID
	result, err := s.repository.CreateProgressReport(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.CreateProgressReport", "failed to create progress report", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) ListProgressReports(ctx context.Context, params domain.ListProgressReportsParams) (*domain.ListProgressReportsResult, error) {
	result, err := s.repository.ListProgressReports(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListProgressReports", "failed to list progress reports", err,
				zap.String("client_id", params.ClientID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) GetProgressReport(ctx context.Context, reportID uuid.UUID) (*domain.ProgressReport, error) {
	result, err := s.repository.GetProgressReport(ctx, reportID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetProgressReport", "failed to get progress report", err,
				zap.String("report_id", reportID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) UpdateProgressReport(ctx context.Context, reportID uuid.UUID, params domain.UpdateProgressReportParams) (*domain.ProgressReport, error) {
	params.ID = reportID
	result, err := s.repository.UpdateProgressReport(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.UpdateProgressReport", "failed to update progress report", err,
				zap.String("report_id", reportID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) DeleteProgressReport(ctx context.Context, reportID uuid.UUID) error {
	err := s.repository.DeleteProgressReport(ctx, reportID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.DeleteProgressReport", "failed to delete progress report", err,
				zap.String("report_id", reportID.String()),
			)
		}
		return err
	}
	return nil
}

func (s *ClientService) GenerateAutoReports(ctx context.Context, clientID uuid.UUID, startDate, endDate time.Time) (string, error) {
	reports, err := s.repository.GetProgressReportsByDateRange(ctx, domain.GetProgressReportsByDateRangeParams{
		ClientID:  clientID,
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GenerateAutoReports", "failed to get progress reports by date range", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return "", err
	}

	var builder strings.Builder
	for _, report := range reports {
		fmt.Fprintf(&builder, "Date: %s\nType: %s\nEmotional State: %s\nReport Text: %s\n\n", report.Date.GoString(), report.Type, report.EmotionalState, report.ReportText)
	}
	text := builder.String()

	result, err := s.reportGen.GenerateAutoReports(ctx, text)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GenerateAutoReports", "failed to generate auto reports", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return "", err
	}
	return result, nil
}

func (s *ClientService) ConfirmAiProgressReport(ctx context.Context, clientID uuid.UUID, reportText string, startDate, endDate time.Time) (*domain.AiGeneratedReport, error) {
	result, err := s.repository.CreateAiGeneratedReport(ctx, domain.CreateAiGeneratedReportParams{
		ClientID:   clientID,
		ReportText: reportText,
		StartDate:  startDate,
		EndDate:    endDate,
	})
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ConfirmAiProgressReport", "failed to create AI generated report", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

func (s *ClientService) ListAiGeneratedReports(ctx context.Context, params domain.ListAiGeneratedReportsParams) (*domain.ListAiGeneratedReportsResult, error) {
	result, err := s.repository.ListAiGeneratedReports(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.ListAiGeneratedReports", "failed to list AI generated reports", err,
				zap.String("client_id", params.ClientID.String()),
			)
		}
		return nil, err
	}
	return result, nil
}

// =====================
// Appointment Card
// =====================

func (s *ClientService) GetAppointmentCard(ctx context.Context, clientID uuid.UUID) (*domain.AppointmentCard, error) {
	card, err := s.repository.GetAppointmentCard(ctx, clientID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GetAppointmentCard", "failed to get appointment card", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, fmt.Errorf("failed to get appointment card")
	}
	if card == nil {
		if s.logger != nil {
			s.logger.LogInfo(ctx, "ClientService.GetAppointmentCard", "appointment card not found",
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, nil
	}
	return card, nil
}

func (s *ClientService) UpdateAppointmentCard(ctx context.Context, clientID uuid.UUID, params domain.UpdateAppointmentCardParams) (*domain.AppointmentCard, error) {
	card, err := s.repository.UpdateAppointmentCard(ctx, clientID, params)
	if err != nil {
		// If no existing row to update, create a new one (upsert)
		if errors.Is(err, pgx.ErrNoRows) {
			createdCard, createErr := s.repository.CreateAppointmentCard(ctx, clientID, params)
			if createErr != nil {
				if s.logger != nil {
					s.logger.LogError(ctx, "ClientService.UpdateAppointmentCard", "failed to create appointment card during upsert", createErr,
						zap.String("client_id", clientID.String()),
					)
				}
				return nil, fmt.Errorf("failed to update appointment card")
			}
			return createdCard, nil
		}
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.UpdateAppointmentCard", "failed to update appointment card", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, fmt.Errorf("failed to update appointment card")
	}
	return card, nil
}

func (s *ClientService) GenerateAppointmentCardDocument(ctx context.Context, clientID uuid.UUID) ([]byte, string, error) {
	card, err := s.repository.GetAppointmentCard(ctx, clientID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GenerateAppointmentCardDocument", "failed to get appointment card", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, "", fmt.Errorf("failed to retrieve appointment card")
	}
	if card == nil {
		if s.logger != nil {
			s.logger.LogWarn(ctx, "ClientService.GenerateAppointmentCardDocument", "appointment card not found",
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, "", fmt.Errorf("appointment card not found")
	}

	pdfData := domain.AppointmentCardPDF{
		ID:                     card.ID,
		ClientName:             card.ClientFirstName + " " + card.ClientLastName,
		Date:                   card.CreatedAt.Format("02-01-2006"),
		GeneralInformation:     card.GeneralInformation,
		ImportantContacts:      card.ImportantContacts,
		HouseholdInfo:          card.HouseholdInfo,
		OrganizationAgreements: card.OrganizationAgreements,
		YouthOfficerAgreements: card.YouthOfficerAgreements,
		TreatmentAgreements:    card.TreatmentAgreements,
		SmokingRules:           card.SmokingRules,
		Work:                   card.Work,
		SchoolInternship:       card.SchoolInternship,
		Travel:                 card.Travel,
		Leave:                  card.Leave,
	}

	pdfBytes, err := s.pdfService.GenerateAppointmentCardPDF(ctx, pdfData)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ClientService.GenerateAppointmentCardDocument", "failed to generate appointment card PDF", err,
				zap.String("client_id", clientID.String()),
			)
		}
		return nil, "", fmt.Errorf("failed to generate appointment card PDF")
	}

	fileName := fmt.Sprintf("appointment_card_%s.pdf", card.ID.String())
	return pdfBytes, fileName, nil
}

// strPtr returns a pointer to the given string.
// Useful for audit event fields that take *string.
func strPtr(s string) *string {
	return &s
}
