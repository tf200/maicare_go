package schedule

import (
	"time"

	"github.com/google/uuid"
)

// CreateScheduleRequest represents the request body for creating a schedule.
type CreateScheduleRequest struct {
	EmployeeIDs []uuid.UUID `json:"employee_ids"`
	LocationID  uuid.UUID   `json:"location_id"`
	IsCustom    bool        `json:"is_custom" example:"true"`                   // true for custom schedule, false for preset shift
	Recurrence  *string     `json:"recurrence,omitempty" example:"end_of_week"` // none (default), end_of_week, end_of_month

	// For custom schedules (required when is_custom = true)
	StartDatetime *time.Time `json:"start_datetime,omitempty" example:"2023-10-01T09:00:00Z"`
	EndDatetime   *time.Time `json:"end_datetime,omitempty" example:"2023-10-01T17:00:00Z"`

	// For preset shift-based schedules (required when is_custom = false)
	LocationShiftID *uuid.UUID `json:"location_shift_id,omitempty" example:"1"`
	ShiftDate       *string    `json:"shift_date,omitempty" example:"2023-10-01"` // Date to apply the shift
}

// CreateScheduleResponse represents the response body after creating a schedule.
type CreateScheduleResponse struct {
	ID            uuid.UUID `json:"id"`
	EmployeeID    uuid.UUID `json:"employee_id"`
	LocationID    uuid.UUID `json:"location_id"`
	LocationName  string    `json:"location_name"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Additional info if created from preset shift
	LocationShiftID *uuid.UUID `json:"location_shift_id,omitempty"`
	ShiftName       *string    `json:"shift_name,omitempty"`
}

const (
	CreateScheduleRecurrenceNone       = "none"
	CreateScheduleRecurrenceEndOfWeek  = "end_of_week"
	CreateScheduleRecurrenceEndOfMonth = "end_of_month"
)

// GetSchedulesByLocationInRangeRequest retrieves schedules for a location within a date range.
type GetSchedulesByLocationInRangeRequest struct {
	StartDate string `form:"start_date" binding:"required" example:"2026-02-01"`
	EndDate   string `form:"end_date" binding:"required" example:"2026-02-29"`
}

// Shift represents a work shift for an employee.
type Shift struct {
	ScheduleID        uuid.UUID  `json:"schedule_id"`
	EmployeeID        uuid.UUID  `json:"employee_id"`
	EmployeeFirstName string     `json:"employee_first_name"`
	EmployeeLastName  string     `json:"employee_last_name"`
	StartTime         time.Time  `json:"start_time"`
	EndTime           time.Time  `json:"end_time"`
	LocationID        uuid.UUID  `json:"location_id"`
	ShiftName         *string    `json:"shift_name,omitempty"`
	LocationShiftID   *uuid.UUID `json:"location_shift_id,omitempty"` // Optional field for preset shift
	IsCustom          bool       `json:"is_custom"`                   // Indicates if this is a custom schedule
}

// GetSchedulesByLocationInRangeResponse represents schedules grouped by day within range.
type GetSchedulesByLocationInRangeResponse struct {
	Date   string  `json:"date"`
	Shifts []Shift `json:"shifts"`
}

// GetScheduleByIdResponse represents the response body for retrieving a schedule by ID.
type GetScheduleByIdResponse struct {
	ID                uuid.UUID  `json:"id"`
	EmployeeID        uuid.UUID  `json:"employee_id"`
	EmployeeFirstName string     `json:"employee_first_name"`
	EmployeeLastName  string     `json:"employee_last_name"`
	LocationID        uuid.UUID  `json:"location_id"`
	LocationName      string     `json:"location_name"`
	LocationShiftID   *uuid.UUID `json:"location_shift_id,omitempty"` // Optional field for preset shift
	LocationShiftName *string    `json:"shift_name,omitempty"`        // Optional field for shift name
	StartDatetime     time.Time  `json:"start_datetime"`
	EndDatetime       time.Time  `json:"end_datetime"`
	IsCustom          bool       `json:"is_custom"` // Indicates if this is a custom schedule
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

// UpdateScheduleRequest represents the request body for updating a schedule.
type UpdateScheduleRequest struct {
	EmployeeID *uuid.UUID `json:"employee_id,omitempty"`
	LocationID *uuid.UUID `json:"location_id,omitempty"`
	IsCustom   *bool      `json:"is_custom,omitempty" example:"true"` // true for custom schedule, false for preset shift

	// For custom schedules (required when is_custom = true)
	StartDatetime *time.Time `json:"start_datetime,omitempty" example:"2023-10-01T09:00:00Z"`
	EndDatetime   *time.Time `json:"end_datetime,omitempty" example:"2023-10-01T17:00:00Z"`

	// For preset shift-based schedules (required when is_custom = false)
	LocationShiftID *uuid.UUID `json:"location_shift_id,omitempty" example:"1"`
	ShiftDate       *string    `json:"shift_date,omitempty" example:"2023-10-01"` // Date to apply the shift

}

// UpdateScheduleResponse represents the response body after updating a schedule.
type UpdateScheduleResponse struct {
	ID            uuid.UUID `json:"id"`
	EmployeeID    uuid.UUID `json:"employee_id"`
	LocationID    uuid.UUID `json:"location_id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LocationName  string    `json:"location_name"`

	// Additional info if updated from preset shift
	LocationShiftID *uuid.UUID `json:"location_shift_id,omitempty"`
	ShiftName       *string    `json:"shift_name,omitempty"`
}
