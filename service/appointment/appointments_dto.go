package appointment

import (
	"time"

	"github.com/google/uuid"
)

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
