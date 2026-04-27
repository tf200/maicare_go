package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// NotificationData mirrors the JSON structure stored in the data column
type NotificationData struct {
	NewAppointment          *NewAppointmentNotificationData          `json:"new_appointment,omitempty"`
	NewClientAssignment     *NewClientAssignmentNotificationData     `json:"new_client_assignment,omitempty"`
	ClientContractReminder  *ClientContractReminderNotificationData  `json:"client_contract_reminder,omitempty"`
	NewIncidentReport       *NewIncidentReportNotificationData       `json:"new_incident_report,omitempty"`
	NewScheduleNotification *NewScheduleNotificationNotificationData `json:"new_schedule_notification,omitempty"`
}

type NewAppointmentNotificationData struct {
	AppointmentID uuid.UUID `json:"appointment_id"`
	CreatedBy     string    `json:"created_by"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Location      string    `json:"location"`
}

type NewClientAssignmentNotificationData struct {
	ClientID        uuid.UUID `json:"client_id"`
	ClientFirstName string    `json:"client_first_name"`
	ClientLastName  string    `json:"client_last_name"`
	ClientLocation  *string   `json:"client_location"`
}

type ClientContractReminderNotificationData struct {
	ClientID           uuid.UUID  `json:"client_id"`
	ClientFirstName    string     `json:"client_first_name"`
	ClientLastName     string     `json:"client_last_name"`
	ContractID         uuid.UUID  `json:"contract_id"`
	CareType           string     `json:"care_type"`
	ContractStart      time.Time  `json:"contract_start"`
	ContractEnd        time.Time  `json:"contract_end"`
	ReminderType       string     `json:"reminder_type"`
	LastReminderSentAt *time.Time `json:"last_reminder_sent_at,omitempty"`
}

type NewIncidentReportNotificationData struct {
	ID                 uuid.UUID `json:"id"`
	EmployeeID         uuid.UUID `json:"employee_id"`
	EmployeeFirstName  string    `json:"employee_first_name"`
	EmployeeLastName   string    `json:"employee_last_name"`
	LocationID         uuid.UUID `json:"location_id"`
	LocationName       string    `json:"location_name"`
	ClientID           uuid.UUID `json:"client_id"`
	ClientFirstName    string    `json:"client_first_name"`
	ClientLastName     string    `json:"client_last_name"`
	SeverityOfIncident string    `json:"severity_of_incident"`
}

type NewScheduleNotificationNotificationData struct {
	ScheduleID uuid.UUID `json:"schedule_id"`
	CreatedBy  uuid.UUID `json:"created_by"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Location   string    `json:"location"`
}

type Notification struct {
	ID        uuid.UUID        `json:"id"`
	UserID    uuid.UUID        `json:"user_id"`
	Type      string           `json:"type"`
	Message   string           `json:"message"`
	IsRead    bool             `json:"is_read"`
	Data      NotificationData `json:"data"`
	CreatedAt time.Time        `json:"created_at"`
}

// Notification type constants
const (
	TypeNewAppointment          = "new_appointment"
	TypeAppointmentUpdate       = "appointment_update"
	TypeNewClientAssignment     = "new_client_assigned"
	TypeClientContractReminder  = "client_contract_reminder"
	TypeIncidentReport          = "incident_report"
	TypeNewScheduleNotification = "new_schedule_notification"
	TypeSystemReminder          = "system_reminder"
)

// NotificationPayload is the payload for creating and delivering notifications
type NotificationPayload struct {
	RecipientUserIDs []uuid.UUID      `json:"recipient_user_ids"`
	Type             string           `json:"type"`
	Data             NotificationData `json:"data"`
	CreatedAt        time.Time        `json:"created_at"`
	Message          string           `json:"message"`
}

// WebSocketMessage is the notification message sent over WebSocket
type WebSocketMessage struct {
	NotificationID   uuid.UUID        `json:"notification_id"`
	NotificationType string           `json:"type"`
	Message          string           `json:"message"`
	IsRead           bool             `json:"is_read"`
	Data             NotificationData `json:"data"`
	CreatedAt        time.Time        `json:"created_at"`
}

// WebSocketEnvelope is the envelope wrapper for WebSocket messages
type WebSocketEnvelope[T any] struct {
	Version string    `json:"v"`
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	TS      time.Time `json:"ts"`
	Data    T         `json:"data"`
}

func (n *NewScheduleNotificationNotificationData) NewScheduleMessage() string {
	return fmt.Sprintf(
		"New schedule created from %s to %s at %s",
		n.StartTime.Format(time.RFC3339), n.EndTime.Format(time.RFC3339), n.Location,
	)
}

func (n *NewScheduleNotificationNotificationData) UpdatedScheduleMessage() string {
	return fmt.Sprintf(
		"Schedule updated from %s to %s at %s",
		n.StartTime.Format(time.RFC3339), n.EndTime.Format(time.RFC3339), n.Location,
	)
}

// NotificationService defines the interface for notification operations
type NotificationService interface {
	ListNotifications(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID, userID uuid.UUID) (*Notification, error)
	CreateAndDeliver(ctx context.Context, payload NotificationPayload) error
}

type NotificationRepository interface {
	ListNotifications(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]Notification, error)
	MarkNotificationAsRead(ctx context.Context, notificationID uuid.UUID) (*Notification, error)
	CreateNotification(ctx context.Context, userID uuid.UUID, notifType string, data []byte, message string) (*Notification, error)
}
