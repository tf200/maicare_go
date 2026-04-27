package handler

import (
	"time"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/google/uuid"
)

type listNotificationsRequest struct {
	httpapi.PageRequest
}

type notificationResponse struct {
	NotificationID   uuid.UUID               `json:"notification_id"`
	NotificationType string                  `json:"type"`
	Message          string                  `json:"message"`
	IsRead           bool                    `json:"is_read"`
	Data             domain.NotificationData `json:"data"`
	CreatedAt        time.Time               `json:"created_at"`
}

func toNotificationResponse(n domain.Notification) notificationResponse {
	return notificationResponse{
		NotificationID:   n.ID,
		NotificationType: n.Type,
		Message:          n.Message,
		IsRead:           n.IsRead,
		Data:             n.Data,
		CreatedAt:        n.CreatedAt,
	}
}
