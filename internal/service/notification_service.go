package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"

	"maicare_go/internal/domain"
	"maicare_go/internal/ws"
)

type NotificationService struct {
	repo   domain.NotificationRepository
	hub    *ws.Hub
	logger domain.Logger
}

func NewNotificationService(repo domain.NotificationRepository, hub *ws.Hub, logger domain.Logger) domain.NotificationService {
	return &NotificationService{repo: repo, hub: hub, logger: logger}
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

func (s *NotificationService) CreateAndDeliver(ctx context.Context, payload domain.NotificationPayload) error {
	var firstError error

	dataBytes, err := json.Marshal(payload.Data)
	if err != nil {
		s.logger.LogError(ctx, "CreateAndDeliver", "Failed to marshal notification data", err)
		return fmt.Errorf("failed to marshal notification data: %w", err)
	}

	for _, recipientID := range payload.RecipientUserIDs {
		notif, dbErr := s.repo.CreateNotification(ctx, recipientID, payload.Type, dataBytes, payload.Message)
		if dbErr != nil {
			log.Printf("Error saving notification to DB for user %d: %v", recipientID, dbErr)
			if firstError == nil {
				firstError = fmt.Errorf("failed to save notification for user %d: %w", recipientID, dbErr)
			}
			continue
		}

		wsMsg := domain.WebSocketMessage{
			NotificationID:   notif.ID,
			NotificationType: notif.Type,
			Message:          notif.Message,
			IsRead:           notif.IsRead,
			Data:             payload.Data,
			CreatedAt:        notif.CreatedAt,
		}

		envelope := domain.WebSocketEnvelope[domain.WebSocketMessage]{
			Version: "1",
			ID:      uuid.NewString(),
			Type:    "notification.created",
			TS:      time.Now().UTC(),
			Data:    wsMsg,
		}

		wsPayload, err := json.Marshal(envelope)
		if err != nil {
			s.logger.LogError(ctx, "CreateAndDeliver", "Failed to marshal websocket payload", err)
			return fmt.Errorf("failed to marshal websocket payload: %w", err)
		}

		if s.hub != nil {
			s.hub.SendToUser(recipientID, wsPayload)
		}
	}

	return firstError
}
