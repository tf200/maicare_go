package schedule

import (
	"time"

	"github.com/google/uuid"
)

// CreateScheduleRequest represents the request body for creating a schedule.
type CreateScheduleRequest struct {
	EmployeeID uuid.UUID `json:"employee_id"`
	LocationID int64     `json:"location_id"`
	IsCustom   bool      `json:"is_custom" example:"true"`          // true for custom schedule, false for preset shift
	Color      *string   `json:"color,omitempty" example:"#FF5733"` // Optional color for the schedule

	// For custom schedules (required when is_custom = true)
	StartDatetime *time.Time `json:"start_datetime,omitempty" example:"2023-10-01T09:00:00Z"`
	EndDatetime   *time.Time `json:"end_datetime,omitempty" example:"2023-10-01T17:00:00Z"`

	// For preset shift-based schedules (required when is_custom = false)
	LocationShiftID *int64  `json:"location_shift_id,omitempty" example:"1"`
	ShiftDate       *string `json:"shift_date,omitempty" example:"2023-10-01"` // Date to apply the shift
}

// CreateScheduleResponse represents the response body after creating a schedule.
type CreateScheduleResponse struct {
	ID            uuid.UUID `json:"id"`
	EmployeeID    uuid.UUID `json:"employee_id"`
	LocationID    int64     `json:"location_id"`
	LocationName  string    `json:"location_name"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	Color         *string   `json:"color"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	// Additional info if created from preset shift
	LocationShiftID *int64  `json:"location_shift_id,omitempty"`
	ShiftName       *string `json:"shift_name,omitempty"`
}

// GetMonthlySchedulesByLocationApi retrieves the monthly schedules for a specific location.
type GetMonthlySchedulesByLocationRequest struct {
	Year  int32 `form:"year"`
	Month int32 `form:"month"`
}

// Shift represents a work shift for an employee.
type Shift struct {
	ShiftID           uuid.UUID `json:"shift_id"`
	EmployeeID        uuid.UUID `json:"employee_id"`
	EmployeeFirstName string    `json:"employee_first_name"`
	EmployeeLastName  string    `json:"employee_last_name"`
	StartTime         time.Time `json:"start_time"`
	EndTime           time.Time `json:"end_time"`
	LocationID        int64     `json:"location_id"`
	Color             *string   `json:"color"` // Optional field for color coding
	ShiftName         *string   `json:"shift_name,omitempty"`
	LocationShiftID   *int64    `json:"location_shift_id,omitempty"` // Optional field for preset shift
	IsCustom          bool      `json:"is_custom"`                   // Indicates if this is a custom schedule
}

// GetMonthlySchedulesByLocationResponse represents the response body for monthly schedules.
type GetMonthlySchedulesByLocationResponse struct {
	Date   string  `json:"date"`
	Shifts []Shift `json:"shifts"`
}

// GetDailySchedulesByLocationApi retrieves the daily schedules for a specific location.
type GetDailySchedulesByLocationRequest struct {
	Year  int32 `form:"year" binding:"required"`
	Month int32 `form:"month" binding:"required"`
	Day   int32 `form:"day" binding:"required"`
}

// GetDailySchedulesByLocationResponse represents the response body for daily schedules.
type GetDailySchedulesByLocationResponse struct {
	Date   string  `json:"date"`
	Shifts []Shift `json:"shifts"`
}

// GetScheduleByIdResponse represents the response body for retrieving a schedule by ID.
type GetScheduleByIdResponse struct {
	ID                uuid.UUID `json:"id"`
	EmployeeID        uuid.UUID `json:"employee_id"`
	EmployeeFirstName string    `json:"employee_first_name"`
	EmployeeLastName  string    `json:"employee_last_name"`
	LocationID        int64     `json:"location_id"`
	LocationName      string    `json:"location_name"`
	LocationShiftID   *int64    `json:"location_shift_id,omitempty"` // Optional field for preset shift
	LocationShiftName *string   `json:"shift_name,omitempty"`        // Optional field for shift name
	Color             *string   `json:"color"`                       // Optional field for color coding
	StartDatetime     time.Time `json:"start_datetime"`
	EndDatetime       time.Time `json:"end_datetime"`
	IsCustom          bool      `json:"is_custom"` // Indicates if this is a custom schedule
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// UpdateScheduleRequest represents the request body for updating a schedule.
type UpdateScheduleRequest struct {
	EmployeeID *uuid.UUID `json:"employee_id,omitempty"`
	LocationID *int64     `json:"location_id,omitempty"`
	IsCustom   *bool      `json:"is_custom,omitempty" example:"true"` // true for custom schedule, false for preset shift

	// For custom schedules (required when is_custom = true)
	StartDatetime *time.Time `json:"start_datetime,omitempty" example:"2023-10-01T09:00:00Z"`
	EndDatetime   *time.Time `json:"end_datetime,omitempty" example:"2023-10-01T17:00:00Z"`

	// For preset shift-based schedules (required when is_custom = false)
	LocationShiftID *int64  `json:"location_shift_id,omitempty" example:"1"`
	ShiftDate       *string `json:"shift_date,omitempty" example:"2023-10-01"` // Date to apply the shift

	Color *string `json:"color,omitempty" example:"#FF5733"`
}

// UpdateScheduleResponse represents the response body after updating a schedule.
type UpdateScheduleResponse struct {
	ID            uuid.UUID `json:"id"`
	EmployeeID    uuid.UUID `json:"employee_id"`
	LocationID    int64     `json:"location_id"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	Color         *string   `json:"color"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LocationName  string    `json:"location_name"`

	// Additional info if updated from preset shift
	LocationShiftID *int64  `json:"location_shift_id,omitempty"`
	ShiftName       *string `json:"shift_name,omitempty"`
}
