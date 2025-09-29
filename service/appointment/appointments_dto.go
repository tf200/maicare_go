package appointment

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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

// ListAppointmentsForClientRequest represents the request payload for listing appointments for a client in a date range
type ListAppointmentsForClientRequest struct {
	StartDate time.Time `json:"start_date" binding:"required" example:"2025-04-27T00:00:00Z"`
	EndDate   time.Time `json:"end_date" binding:"required" example:"2025-04-30T23:59:59Z"`
}

// ListAppointmentsForClientResponse represents the response payload for listing appointments for a client in a date range
type ListAppointmentsForClientResponse struct {
	ID                    uuid.UUID             `json:"id"`
	CreatorEmployeeID     *int64                `json:"creator_employee_id"`
	StartTime             time.Time             `json:"start_time"`
	EndTime               time.Time             `json:"end_time"`
	Location              *string               `json:"location"`
	Description           *string               `json:"description"`
	Color                 *string               `json:"color"`
	Status                string                `json:"status"`
	RecurrenceType        *string               `json:"recurrence_type"`
	RecurrenceInterval    *int32                `json:"recurrence_interval"`
	RecurrenceEndDate     pgtype.Date           `json:"recurrence_end_date"`
	ConfirmedByEmployeeID *int32                `json:"confirmed_by_employee_id"`
	ConfirmedAt           time.Time             `json:"confirmed_at"`
	CreatedAt             time.Time             `json:"created_at"`
	ParticipantsDetails   []ParticipantsDetails `json:"participants_details"`
	ClientsDetails        []ClientsDetails      `json:"clients_details"`
}

// GetAppointmentResponse represents the response payload for getting an appointment
type GetAppointmentResponse struct {
	ID                     uuid.UUID             `json:"id"`
	AppointmentTemplatesID *int64                `json:"appointment_templates_id"`
	CreatorEmployeeID      *int64                `json:"creator_employee_id"`
	CreatorFirstName       *string               `json:"creator_first_name"`
	CreatorLastName        *string               `json:"creator_last_name"`
	StartTime              time.Time             `json:"start_time"`
	EndTime                time.Time             `json:"end_time"`
	Location               *string               `json:"location"`
	Description            *string               `json:"description"`
	Color                  *string               `json:"color"`
	Status                 string                `json:"status"`
	IsConfirmed            bool                  `json:"is_confirmed"`
	ConfirmedByEmployeeID  *int64                `json:"confirmed_by_employee_id"`
	ConfirmerFirstName     *string               `json:"confirmer_first_name"`
	ConfirmerLastName      *string               `json:"confirmer_last_name"`
	ConfirmedAt            time.Time             `json:"confirmed_at"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
	ParticipantsDetails    []ParticipantsDetails `json:"participants_details"`
	ClientsDetails         []ClientsDetails      `json:"clients_details"`
}


// UpdateAppointmentRequest represents the request payload for updating an appointment
type UpdateAppointmentRequest struct {
	StartTime              time.Time `json:"start_time" binding:"required" example:"2023-10-01T10:00:00Z"`
	EndTime                time.Time `json:"end_time" binding:"required"`
	Location               *string   `json:"location"`
	Color                  *string   `json:"color" example:"#FF5733"`
	Description            *string   `json:"description"`
	ClientIDs              *[]int64  `json:"client_ids"`
	ParticipantEmployeeIDs *[]int64  `json:"participant_employee_ids"`
}

// UpdateAppointmentResponse represents the response payload for updating an appointment
type UpdateAppointmentResponse struct {
	ID                     uuid.UUID        `json:"id"`
	AppointmentTemplatesID *uuid.UUID       `json:"appointment_templates_id"`
	CreatorEmployeeID      *int64           `json:"creator_employee_id"`
	StartTime              pgtype.Timestamp `json:"start_time"`
	EndTime                pgtype.Timestamp `json:"end_time"`
	Location               *string          `json:"location"`
	Description            *string          `json:"description"`
	Color                  *string          `json:"color"`
	Status                 string           `json:"status"`
	IsConfirmed            bool             `json:"is_confirmed"`
	ConfirmedByEmployeeID  *int64           `json:"confirmed_by_employee_id"`
	ConfirmedAt            pgtype.Timestamp `json:"confirmed_at"`
	CreatedAt              pgtype.Timestamp `json:"created_at"`
	UpdatedAt              pgtype.Timestamp `json:"updated_at"`
}

