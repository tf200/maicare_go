package appointment

import (
	"time"

	"github.com/google/uuid"
)

// GetAppointmentResponse represents the response payload for getting an appointment
type ParticipantsDetails struct {
	EmployeeID int64  `json:"employee_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}

// ClientsDetails represents the details of a client
type ClientsDetails struct {
	ClientID  int64  `json:"client_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// CreateAppointmentRequest represents the request payload for creating an appointment
type CreateAppointmentRequest struct {
	StartTime              time.Time `json:"start_time" binding:"required" example:"2023-10-01T10:00:00Z"`
	EndTime                time.Time `json:"end_time" binding:"required"`
	Location               *string   `json:"location"`
	Description            *string   `json:"description"`
	Color                  *string   `json:"color" example:"#FF5733"`
	RecurrenceType         string    `json:"recurrence_type" example:"NONE" enum:"NONE,DAILY,WEEKLY,MONTHLY"`
	RecurrenceInterval     *int32    `json:"recurrence_interval"`
	RecurrenceEndDate      time.Time `json:"recurrence_end_date" example:"2025-10-01T10:00:00Z"`
	ParticipantEmployeeIDs []int64   `json:"participant_employee_ids"`
	ClientIDs              []int64   `json:"client_ids"`
}

// CreateAppointmentResponse represents the response payload for creating an appointment
type CreateAppointmentResponse struct {
	ID                uuid.UUID `json:"id"`
	CreatorEmployeeID *int64    `json:"creator_employee_id"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	Color             *string   `json:"color"`
	Location          *string   `json:"location"`
	Description       *string   `json:"description"`
}

// AddParticipantToAppointmentRequest represents the request payload for adding participants to an appointment
type AddParticipantToAppointmentRequest struct {
	ParticipantEmployeeIDs []int64 `json:"participant_employee_ids" binding:"required"`
}

// AddClientToAppointmentRequest represents the request payload for adding clients to an appointment
type AddClientToAppointmentRequest struct {
	ClientIDs []int64 `json:"client_ids" binding:"required"`
}

// ListAppointmentsForEmployeeInRangeRequest represents the request payload for listing appointments for an employee in a date range
type ListAppointmentsForEmployeeInRangeRequest struct {
	StartDate time.Time `json:"start_date" binding:"required" example:"2025-04-27T00:00:00Z"`
	EndDate   time.Time `json:"end_date" binding:"required" example:"2025-04-30T23:59:59Z"`
}

// ListAppointmentsForEmployeeInRangeResponse represents the response payload for listing appointments for an employee in a date range
type ListAppointmentsForEmployeeInRangeResponse struct {
	ID                  uuid.UUID             `json:"id"`
	CreatorEmployeeID   *int64                `json:"creator_employee_id"`
	StartTime           time.Time             `json:"start_time"`
	EndTime             time.Time             `json:"end_time"`
	Location            *string               `json:"location"`
	Description         *string               `json:"description"`
	Color               *string               `json:"color"`
	Status              string                `json:"status"`
	IsConfirmed         bool                  `json:"is_confirmed"`
	CreatedAt           time.Time             `json:"created_at"`
	ParticipantsDetails []ParticipantsDetails `json:"participants_details"`
	ClientsDetails      []ClientsDetails      `json:"clients_details"`
}
