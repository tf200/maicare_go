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
	TypeNewIncidentReport      = "new_incident_report"
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
