package notification

import (
	"github.com/google/uuid"
)

const (
	// Notification Types
	TypeNewAppointment      = "new_appointment"
	TypeAppointmentUpdate   = "appointment_update"
	TypeNewClientAssignment = "new_client_assigned"
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
