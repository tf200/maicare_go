package service

import (
	"context"
	"fmt"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SenderService struct {
	repo   domain.SenderRepository
	logger domain.Logger
}

func NewSenderService(repo domain.SenderRepository, logger domain.Logger) domain.SenderService {
	return &SenderService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SenderService) CreateSender(ctx context.Context, params domain.CreateSenderParams) (*domain.SenderEntity, error) {
	result, err := s.repo.Create(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "SenderService.CreateSender", "failed to create sender", err)
		}
		return nil, fmt.Errorf("failed to create sender: %w", err)
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "SenderService.CreateSender", "sender created successfully", zap.String("sender_id", result.ID.String()))
	}
	return result, nil
}

func (s *SenderService) ListSenders(ctx context.Context, params domain.ListSendersParams) (*domain.ListResult[domain.SenderEntity], error) {
	senders, totalCount, err := s.repo.List(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "SenderService.ListSenders", "failed to list senders", err)
		}
		return nil, fmt.Errorf("failed to list senders: %w", err)
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "SenderService.ListSenders", "senders listed successfully", zap.Int("count", len(senders)))
	}
	return &domain.ListResult[domain.SenderEntity]{
		Items:      senders,
		TotalCount: totalCount,
	}, nil
}

func (s *SenderService) GetSenderByID(ctx context.Context, id uuid.UUID) (*domain.SenderEntityWithTemplate, error) {
	sender, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "SenderService.GetSenderByID", "failed to get sender by ID", err, zap.String("sender_id", id.String()))
		}
		return nil, fmt.Errorf("failed to get sender: %w", err)
	}

	templateItems, err := s.repo.GetTemplateItemsBySourceTable(ctx, sender.InvoiceTemplate)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "SenderService.GetSenderByID", "failed to get template items", err, zap.String("sender_id", id.String()))
		}
		return nil, fmt.Errorf("failed to get template items: %w", err)
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "SenderService.GetSenderByID", "sender retrieved successfully", zap.String("sender_id", id.String()))
	}

	return &domain.SenderEntityWithTemplate{
		SenderEntity:         *sender,
		InvoiceTemplateItems: templateItems,
	}, nil
}

func (s *SenderService) UpdateSender(ctx context.Context, params domain.UpdateSenderParams) (*domain.SenderEntity, error) {
	result, err := s.repo.Update(ctx, params)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "SenderService.UpdateSender", "failed to update sender", err, zap.String("sender_id", params.ID.String()))
		}
		return nil, fmt.Errorf("failed to update sender: %w", err)
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "SenderService.UpdateSender", "sender updated successfully", zap.String("sender_id", params.ID.String()))
	}
	return result, nil
}

func (s *SenderService) DeleteSender(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "SenderService.DeleteSender", "failed to delete sender", err, zap.String("sender_id", id.String()))
		}
		return fmt.Errorf("failed to delete sender: %w", err)
	}
	if s.logger != nil {
		s.logger.LogInfo(ctx, "SenderService.DeleteSender", "sender deleted successfully", zap.String("sender_id", id.String()))
	}
	return nil
}

func (s *SenderService) CreateSenderInvoiceTemplate(ctx context.Context, senderID uuid.UUID, templateIDs []uuid.UUID) error {
	validIDs, err := s.repo.GetTemplateItemsByIDs(ctx, templateIDs)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "SenderService.CreateSenderInvoiceTemplate", "failed to validate template IDs", err, zap.String("sender_id", senderID.String()))
		}
		return fmt.Errorf("failed to get template items by IDs: %w", err)
	}

	err = s.repo.CreateInvoiceTemplate(ctx, senderID, validIDs)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "SenderService.CreateSenderInvoiceTemplate", "failed to create invoice template", err, zap.String("sender_id", senderID.String()))
		}
		return fmt.Errorf("failed to create sender invoice template: %w", err)
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "SenderService.CreateSenderInvoiceTemplate", "sender invoice template created successfully", zap.String("sender_id", senderID.String()))
	}
	return nil
}
