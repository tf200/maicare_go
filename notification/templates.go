package notification

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Notification Types
	TypeNewAppointment         = "new_appointment"
	TypeAppointmentUpdate      = "appointment_update"
	TypeNewClientAssignment    = "new_client_assigned"
	TypeClientContractReminder = "client_contract_reminder"
)

// Notifications Data Templates

type NewAppointmentData struct {
	AppointmentID uuid.UUID `json:"appointment_id"`
	CreatedBy     string    `json:"created_by"`
	StartTime     string    `json:"start_time"`
	EndTime       string    `json:"end_time"`
	Location      string    `json:"location"`
}

type NewClientAssignmentData struct {
	ClientID        int64   `json:"client_id"`
	ClientFirstName string  `json:"client_first_name"`
	ClientLastName  string  `json:"client_last_name"`
	ClientLocation  *string `json:"client_location"`
}

type ClientContractReminderData struct {
	ClientID           int64      `json:"client_id"`
	ClientFirstName    string     `json:"client_first_name"`
	ClientLastName     string     `json:"client_last_name"`
	ContractID         int64      `json:"contract_id"`
	CareType           string     `json:"care_type"` // e.g., "ambulante", "accommodation"
	ContractStart      time.Time  `json:"contract_start"`
	ContractEnd        time.Time  `json:"contract_end"`
	ReminderType       string     `json:"reminder_type"` // e.g., "initial
	LastReminderSentAt *time.Time `json:"last_reminder_sent_at,omitempty"`
}
