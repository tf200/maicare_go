package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/service/deps"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type NotificationService interface {
	CreateAndDeliver(ctx context.Context, payload NotificationPayload) error
	ListNotifications(ctx context.Context, req *ListNotificationsRequest, userID uuid.UUID) ([]ListNotificationsResponse, error)
	MarkNotificationAsRead(ctx context.Context, notificationID uuid.UUID, userID uuid.UUID) (*MarkNotificationAsReadResponse, error)
}

type notificationService struct {
	*deps.ServiceDependencies
}

func NewNotificationService(deps *deps.ServiceDependencies) NotificationService {
	return &notificationService{
		ServiceDependencies: deps,
	}
}

type WebSocketMessage struct {
	NotificationID   uuid.UUID        `json:"notification_id"`
	NotificationType string           `json:"type"`
	Message          string           `json:"message"`
	IsRead           bool             `json:"is_read"`
	Data             NotificationData `json:"data"`
	CreatedAt        time.Time        `json:"created_at"`
}

func (s *notificationService) CreateAndDeliver(ctx context.Context, payload NotificationPayload) error {
	var firstError error

	dataBytes, err := json.Marshal(payload.Data)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAndDeliver", "Failed to marshal notification data", zap.Error(err))
		return fmt.Errorf("failed to marshal notification data: %w", err)
	}

	for _, recipientID := range payload.RecipientUserIDs {
		log.Printf("Processing notification for recipient ID: %d", recipientID)
		// 1. Save to Database
		notif, dbErr := s.Store.CreateNotification(ctx, db.CreateNotificationParams{
			UserID:  recipientID,
			Type:    payload.Type,
			Data:    dataBytes,
			Message: payload.Message, // Use the original data bytes
			// You might want to store CreatedAt from the payload too,
			// ensure your DB schema/params support this if needed.
		})

		if dbErr != nil {
			log.Printf("Error saving notification to DB for user %d: %v", recipientID, dbErr)
			// Capture the first error encountered
			if firstError == nil {
				firstError = fmt.Errorf("failed to save notification for user %d: %w", recipientID, dbErr)
			}
			// Decide if you want to skip WS delivery on DB error. Usually yes.
			continue // Skip WS delivery for this user if DB save failed
		}

		log.Printf("Notification saved to DB for user %d.", recipientID)
		// Prepare WebSocket message
		wsMsg := WebSocketMessage{
			NotificationID:   notif.ID,
			NotificationType: notif.Type,
			Message:          notif.Message,
			IsRead:           notif.IsRead,
			Data:             payload.Data,
			CreatedAt:        notif.CreatedAt.Time,
		}

		wsPayload, err := json.Marshal(wsMsg)
		if err != nil {
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAndDeliver", fmt.Sprintf("Error marshalling WebSocket message (Type: %s): %v", payload.Type, err))
			// If we can't marshal this, we can't send it via WS.
			// Depending on requirements, you might still want to proceed with DB saves,
			// or return an error here. Let's log and proceed with DB saves for now.
			return fmt.Errorf("failed to marshal websocket payload: %w", err) // Uncomment this if WS delivery is critical
		}

		// 2. Deliver via WebSocket (if marshalling succeeded)
		if s.WsHub != nil { // Check if marshalling failed earlier and hub exists
			// The hub's SendToUser handles checking if the user is actually connected.
			// It iterates through all connections for that user ID.
			s.WsHub.SendToUser(recipientID, wsPayload)
			// Log the *attempt* to send. The hub logs success/failure per connection.
			s.Logger.LogBusinessEvent(logger.LogLevelInfo, "CreateAndDeliver", fmt.Sprintf("Attempted WebSocket delivery to user %d.", recipientID))
		} else if s.WsHub == nil {
			s.Logger.LogBusinessEvent(logger.LogLevelWarn, "CreateAndDeliver", fmt.Sprintf("WebSocket Hub is nil, skipping WS delivery for user %d.", recipientID))
		} else {
			// This means json.Marshal(wsMsg) failed earlier
			s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateAndDeliver", "Failed to marshal WebSocket message", zap.Error(err))
		}
	}

	// Return the first error encountered during DB operations, or nil if all succeeded.
	// Asynq will handle retries based on this error return.
	return firstError
}
