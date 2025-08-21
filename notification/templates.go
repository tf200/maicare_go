package notification

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	// Notification Types
	TypeNewAppointment         = "new_appointment"
	TypeAppointmentUpdate      = "appointment_update"
	TypeNewClientAssignment    = "new_client_assigned"
	TypeClientContractReminder = "client_contract_reminder"
	TypeNewIncidentReport      = "new_incident_report"
)

type NotificationPayload struct {
	RecipientUserIDs []int64          `json:"recipient_user_ids"`
	Type             string           `json:"type"`
	Data             NotificationData `json:"data"`
	CreatedAt        time.Time        `json:"created_at"`
	Message          string           `json:"message"`
}

type NotificationData struct {
	NewAppointment         *NewAppointmentData         `json:"new_appointment,omitempty"`
	NewClientAssignment    *NewClientAssignmentData    `json:"new_client_assignment,omitempty"`
	ClientContractReminder *ClientContractReminderData `json:"client_contract_reminder,omitempty"`
	NewIncidentReport      *NewIncidentReportData      `json:"new_incident_report,omitempty"`
}

// Notifications Data Templates

type NewAppointmentData struct {
	AppointmentID uuid.UUID `json:"appointment_id"`
	CreatedBy     string    `json:"created_by"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Location      string    `json:"location"`
}

func (a *NewAppointmentData) NewAppointmentMessage() string {
	message := fmt.Sprintf(
		"New appointment created by %s from %s to %s at %s",
		a.CreatedBy, a.StartTime.Format(time.RFC3339), a.EndTime.Format(time.RFC3339), a.Location,
	)
	return message
}

type NewClientAssignmentData struct {
	ClientID        int64   `json:"client_id"`
	ClientFirstName string  `json:"client_first_name"`
	ClientLastName  string  `json:"client_last_name"`
	ClientLocation  *string `json:"client_location"`
}

func (n *NewClientAssignmentData) NewClientAssignmentMessage() string {
	if n.ClientLocation != nil {
		return fmt.Sprintf("New client assigned: %s %s at %s", n.ClientFirstName, n.ClientLastName, *n.ClientLocation)
	}
	return fmt.Sprintf("New client assigned: %s %s", n.ClientFirstName, n.ClientLastName)
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

type NewIncidentReportData struct {
	ID                 int64  `json:"id"`
	EmployeeID         int64  `json:"employee_id"`
	EmployeeFirstName  string `json:"employee_first_name"`
	EmployeeLastName   string `json:"employee_last_name"`
	LocationID         int64  `json:"location_id"`
	LocationName       string `json:"location_name"`
	ClientID           int64  `json:"client_id"`
	ClientFirstName    string `json:"client_first_name"`
	ClientLastName     string `json:"client_last_name"`
	SeverityOfIncident string `json:"severity_of_incident"`
}
