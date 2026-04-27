package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ContractService struct {
	repository domain.ContractRepository
	storage    domain.Storage
	logger     domain.Logger
}

func NewContractService(repository domain.ContractRepository, storage domain.Storage, logger domain.Logger) domain.ContractService {
	return &ContractService{
		repository: repository,
		storage:    storage,
		logger:     logger,
	}
}

// ==================== ContractType ====================

func (s *ContractService) CreateContractType(ctx context.Context, params domain.CreateContractTypeParams) (*domain.CreateContractTypeResult, error) {
	ct, err := s.repository.CreateContractType(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.CreateContractType", "failed to create contract type", err, zap.String("name", params.Name))
		}
		return nil, err
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "ContractService.CreateContractType", "contract type created successfully", zap.String("contract_type_id", ct.ID.String()))
	}
	return &domain.CreateContractTypeResult{ID: ct.ID, Name: ct.Name}, nil
}

func (s *ContractService) ListContractTypes(ctx context.Context) ([]domain.CreateContractTypeResult, error) {
	types, err := s.repository.ListContractTypes(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.ListContractTypes", "failed to list contract types", err)
		}
		return nil, err
	}
	result := make([]domain.CreateContractTypeResult, len(types))
	for i, t := range types {
		result[i] = domain.CreateContractTypeResult{ID: t.ID, Name: t.Name}
	}
	return result, nil
}

func (s *ContractService) DeleteContractType(ctx context.Context, contractTypeID uuid.UUID) (*domain.CreateContractTypeResult, error) {
	err := s.repository.DeleteContractType(ctx, contractTypeID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.DeleteContractType", "failed to delete contract type", err, zap.String("contract_type_id", contractTypeID.String()))
		}
		return nil, err
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "ContractService.DeleteContractType", "contract type deleted successfully", zap.String("contract_type_id", contractTypeID.String()))
	}
	return &domain.CreateContractTypeResult{ID: contractTypeID}, nil
}

// ==================== Contract ====================

func (s *ContractService) CreateContract(ctx context.Context, params domain.CreateContractParams) (*domain.CreateContractResult, error) {
	careName := strings.TrimSpace(params.CareName)
	if careName == "" {
		return nil, fmt.Errorf("care_name is required")
	}
	params.CareName = careName

	if params.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than 0")
	}

	if params.EndDate.Before(params.StartDate) || params.EndDate.Equal(params.StartDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	if params.ReminderPeriod != nil && *params.ReminderPeriod < 0 {
		return nil, fmt.Errorf("reminder_period must be greater than or equal to 0")
	}

	if err := validateContractCarePricing(params.CareType, params.PriceTimeUnit, params.Hours, params.HoursType); err != nil {
		return nil, err
	}

	attachmentIDs := normalizeAttachmentIDs(params.AttachmentIds)
	if err := s.validateAttachmentIds(ctx, attachmentIDs); err != nil {
		return nil, err
	}
	params.AttachmentIds = attachmentIDs

	if params.Vat != nil && (*params.Vat < 0 || *params.Vat > 100) {
		return nil, fmt.Errorf("VAT must be between 0 and 100")
	}

	contract, err := s.repository.CreateContract(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.CreateContract", "failed to create contract", err, zap.String("client_id", params.ClientID.String()))
		}
		return nil, err
	}

	attachments := s.fetchAttachmentDetails(ctx, contract.AttachmentIds)

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ContractService.CreateContract", "contract created successfully", zap.String("contract_id", contract.ID.String()), zap.String("client_id", params.ClientID.String()))
	}

	return &domain.CreateContractResult{
		Contract:    *contract,
		Attachments: attachments,
	}, nil
}

func (s *ContractService) GetContractByID(ctx context.Context, contractID uuid.UUID) (*domain.ContractDetail, error) {
	contract, err := s.repository.GetContractByID(ctx, contractID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.GetContractByID", "failed to get contract", err, zap.String("contract_id", contractID.String()))
		}
		return nil, err
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "ContractService.GetContractByID", "contract retrieved successfully", zap.String("contract_id", contractID.String()))
	}
	return contract, nil
}

func (s *ContractService) UpdateContract(ctx context.Context, params domain.UpdateContractParams, employeeID uuid.UUID) (*domain.UpdateContractResult, error) {
	existing, err := s.repository.GetContractByID(ctx, params.ContractID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.UpdateContract", "failed to get existing contract", err, zap.String("contract_id", params.ContractID.String()))
		}
		return nil, err
	}

	startDate := existing.StartDate
	if params.StartDate != nil {
		startDate = *params.StartDate
	}
	endDate := existing.EndDate
	if params.EndDate != nil {
		endDate = *params.EndDate
	}
	if endDate.Before(startDate) || endDate.Equal(startDate) {
		return nil, fmt.Errorf("end_date must be after start_date")
	}

	reminderPeriod := existing.ReminderPeriod
	if params.ReminderPeriod != nil {
		reminderPeriod = *params.ReminderPeriod
	}
	if reminderPeriod < 0 {
		return nil, fmt.Errorf("reminder_period must be greater than or equal to 0")
	}

	vat := existing.Vat
	if params.Vat != nil {
		vat = params.Vat
	}
	if vat != nil && (*vat < 0 || *vat > 100) {
		return nil, fmt.Errorf("VAT must be between 0 and 100")
	}

	price := existing.Price
	if params.Price != nil {
		price = *params.Price
	}
	if price <= 0 {
		return nil, fmt.Errorf("price must be greater than 0")
	}

	priceTimeUnit := existing.PriceTimeUnit
	if params.PriceTimeUnit != nil {
		priceTimeUnit = *params.PriceTimeUnit
	}

	hours := existing.Hours
	if params.Hours != nil {
		hours = params.Hours
	}

	hoursType := existing.HoursType
	if params.HoursType != nil {
		hoursType = params.HoursType
	}

	careName := existing.CareName
	if params.CareName != nil {
		careName = strings.TrimSpace(*params.CareName)
	}
	if careName == "" {
		return nil, fmt.Errorf("care_name is required")
	}

	careType := existing.CareType
	if err := validateContractCarePricing(careType, priceTimeUnit, hours, hoursType); err != nil {
		return nil, err
	}

	attachmentIDs := existing.AttachmentIds
	if params.AttachmentIds != nil {
		attachmentIDs = normalizeAttachmentIDs(params.AttachmentIds)
	}
	if err := s.validateAttachmentIds(ctx, attachmentIDs); err != nil {
		return nil, err
	}
	params.AttachmentIds = attachmentIDs

	contract, err := s.repository.UpdateContract(ctx, params, employeeID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.UpdateContract", "failed to update contract", err, zap.String("contract_id", params.ContractID.String()))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ContractService.UpdateContract", "contract updated successfully", zap.String("contract_id", params.ContractID.String()))
	}

	return &domain.UpdateContractResult{Contract: *contract}, nil
}

func (s *ContractService) UpdateContractStatus(ctx context.Context, params domain.UpdateContractStatusParams, employeeID uuid.UUID) (*domain.UpdateContractStatusResult, error) {
	existing, err := s.repository.GetContractByID(ctx, params.ContractID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.UpdateContractStatus", "failed to get contract", err, zap.String("contract_id", params.ContractID.String()))
		}
		return nil, err
	}

	if params.Status == "approved" && existing.EndDate.Before(time.Now()) {
		err := fmt.Errorf("cannot approve contract that has already ended")
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.UpdateContractStatus", "cannot approve ended contract", err, zap.String("contract_id", params.ContractID.String()))
		}
		return nil, err
	}

	contract, err := s.repository.UpdateContractStatus(ctx, params, employeeID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.UpdateContractStatus", "failed to update contract status", err, zap.String("contract_id", params.ContractID.String()), zap.String("status", params.Status))
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "ContractService.UpdateContractStatus", "contract status updated successfully", zap.String("contract_id", params.ContractID.String()), zap.String("status", params.Status))
	}

	return &domain.UpdateContractStatusResult{ID: contract.ID, Status: string(contract.Status)}, nil
}

func (s *ContractService) ListClientContracts(ctx context.Context, params domain.ListClientContractsParams) (*domain.ListClientContractsResult, error) {
	contracts, totalCount, err := s.repository.ListClientContracts(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.ListClientContracts", "failed to list client contracts", err, zap.String("client_id", params.ClientID.String()))
		}
		return nil, err
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "ContractService.ListClientContracts", "client contracts listed successfully", zap.String("client_id", params.ClientID.String()))
	}
	return &domain.ListClientContractsResult{Contracts: contracts, TotalCount: totalCount}, nil
}

func (s *ContractService) ListContracts(ctx context.Context, params domain.ListContractsParams) (*domain.ListContractsResult, error) {
	contracts, totalCount, err := s.repository.ListContracts(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.ListContracts", "failed to list contracts", err)
		}
		return nil, err
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "ContractService.ListContracts", "contracts listed successfully")
	}
	return &domain.ListContractsResult{Contracts: contracts, TotalCount: totalCount}, nil
}

func (s *ContractService) GetContractAuditLog(ctx context.Context, contractID uuid.UUID) ([]domain.ContractAuditLog, error) {
	logs, err := s.repository.GetContractAuditLog(ctx, contractID)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.GetContractAuditLog", "failed to get contract audit logs", err, zap.String("contract_id", contractID.String()))
		}
		return nil, err
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "ContractService.GetContractAuditLog", "contract audit logs retrieved successfully", zap.String("contract_id", contractID.String()))
	}
	return logs, nil
}

// ==================== Helpers ====================

func normalizeAttachmentIDs(attachmentIDs []uuid.UUID) []uuid.UUID {
	if len(attachmentIDs) == 0 {
		return []uuid.UUID{}
	}
	normalized := make([]uuid.UUID, 0, len(attachmentIDs))
	seen := make(map[uuid.UUID]struct{}, len(attachmentIDs))
	for _, id := range attachmentIDs {
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		normalized = append(normalized, id)
	}
	return normalized
}

func (s *ContractService) validateAttachmentIds(ctx context.Context, attachmentIds []uuid.UUID) error {
	if len(attachmentIds) == 0 {
		return nil
	}
	attachments, err := s.repository.GetAttachmentFiles(ctx, attachmentIds)
	if err != nil {
		return fmt.Errorf("failed to validate attachment IDs: %w", err)
	}
	found := make(map[uuid.UUID]bool, len(attachments))
	for _, a := range attachments {
		found[a.UUID] = true
	}
	var missing []string
	for _, id := range attachmentIds {
		if !found[id] {
			missing = append(missing, id.String())
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("attachment IDs not found: %v", missing)
	}
	for _, a := range attachments {
		if !a.IsUsed {
			return fmt.Errorf("attachment %s has not been confirmed (upload not completed)", a.UUID.String())
		}
	}
	return nil
}

func (s *ContractService) fetchAttachmentDetails(ctx context.Context, attachmentIds []uuid.UUID) []domain.ContractAttachment {
	if len(attachmentIds) == 0 {
		return []domain.ContractAttachment{}
	}
	attachments, err := s.repository.GetAttachmentFiles(ctx, attachmentIds)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "ContractService.fetchAttachmentDetails", "failed to fetch attachments", err)
		}
		return []domain.ContractAttachment{}
	}
	details := make([]domain.ContractAttachment, 0, len(attachments))
	for _, a := range attachments {
		detail := domain.ContractAttachment{
			ID:   a.UUID,
			Name: a.Name,
			Size: int64(a.Size),
		}
		url, err := s.storage.GeneratePresignedURL(ctx, a.File, 15*time.Minute)
		if err == nil {
			detail.DownloadURL = url
		}
		details = append(details, detail)
	}
	return details
}

func validateContractCarePricing(careType string, priceTimeUnit string, hours *float64, hoursType *string) error {
	switch careType {
	case string(db.CareTypeEnumAmbulante):
		if priceTimeUnit != string(db.PriceTimeUnitEnumMinute) && priceTimeUnit != string(db.PriceTimeUnitEnumHourly) {
			return fmt.Errorf("ambulante contracts require price_time_unit to be minute or hourly")
		}
		if hours == nil || *hours <= 0 {
			return fmt.Errorf("ambulante contracts require hours to be greater than 0")
		}
		if hoursType == nil {
			return fmt.Errorf("ambulante contracts require hours_type")
		}
	case string(db.CareTypeEnumAccommodation):
		if priceTimeUnit != string(db.PriceTimeUnitEnumDaily) && priceTimeUnit != string(db.PriceTimeUnitEnumWeekly) {
			return fmt.Errorf("accommodation contracts require price_time_unit to be daily or weekly")
		}
		if hours != nil || hoursType != nil {
			return fmt.Errorf("accommodation contracts require hours and hours_type to be null")
		}
	default:
		return fmt.Errorf("invalid care_type")
	}
	return nil
}
