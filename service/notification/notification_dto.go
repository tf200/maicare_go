package notification

import (
	"maicare_go/pagination"
	"time"

	"github.com/google/uuid"
)

// ListNotificationsRequest defines the request structure for listing notifications
type ListNotificationsRequest struct {
	pagination.Request
}

// ListNotificationsResponse defines the response structure for listing notifications
type ListNotificationsResponse struct {
	NotificationID   uuid.UUID        `json:"notification_id"`
	NotificationType string           `json:"type"`
	Message          string           `json:"message"`
	IsRead           bool             `json:"is_read"`
	Data             NotificationData `json:"data"`
	CreatedAT        time.Time        `json:"created_at"`
}

// MarkNotificationAsReadResponse represents the response for marking a notification as read
type MarkNotificationAsReadResponse struct {
	NotificationID   uuid.UUID `json:"notification_id"`
	NotificationType string    `json:"notification_type"`
	Message          string    `json:"message"`
	IsRead           bool      `json:"is_read"`
	CreatedAT        time.Time `json:"created_at"`
}
