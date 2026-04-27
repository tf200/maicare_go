package service

import (
	"context"
	"fmt"

	"maicare_go/internal/domain"
	"maicare_go/internal/repository"

	"github.com/google/uuid"
)

type NotificationService struct {
	repo   *repository.NotificationRepository
	logger domain.Logger
}

func NewNotificationService(repo *repository.NotificationRepository, logger domain.Logger) domain.NotificationService {
	return &NotificationService{repo: repo, logger: logger}
}

func (s *NotificationService) ListNotifications(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]domain.Notification, error) {
	return s.repo.ListNotifications(ctx, userID, limit, offset)
}

func (s *NotificationService) MarkNotificationAsRead(ctx context.Context, notificationID, userID uuid.UUID) (*domain.Notification, error) {
	notif, err := s.repo.MarkNotificationAsRead(ctx, notificationID)
	if err != nil {
		return nil, fmt.Errorf("failed to mark notification as read: %w", err)
	}

	if notif.UserID != userID {
		return nil, fmt.Errorf("notification does not belong to user")
	}

	return notif, nil
}
