package employees

import (
	"time"

	"github.com/google/uuid"
)

// WorkingHourItem represents a single working hour item, which can be either a schedule or an appointment.
type WorkingHourItem struct {
	ID            uuid.UUID `json:"id"`
	Type          string    `json:"type"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	DurationHours float64   `json:"duration_hours"`
	Location      string    `json:"location"`
	LocationID    *int64    `json:"location_id,omitempty"`
	Description   *string   `json:"description,omitempty"`
	Status        *string   `json:"status"`
	Color         string    `json:"color"`
}

// Summary contains the summary of working hours for an employee in a given period.
type Summary struct {
	TotalHours       float64 `json:"total_hours"`
	AppointmentHours float64 `json:"appointment_hours"`
	ShiftHours       float64 `json:"shift_hours"`
	TotalDaysWorked  int     `json:"total_days_worked"`
	OverTime         float64 `json:"over_time"` // Optional field for overtime hours
}

// Period information
type Period struct {
	Year          int32  `json:"year"`
	Week          int32  `json:"Week"`
	MonthName     string `json:"month_name"`
	IsCurrentWeek bool   `json:"is_current_week"`
	DateRange     struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"date_range"`
}

// ListWorkingHoursRequest represents the request parameters for listing working hours.
type ListWorkingHoursRequest struct {
	Year int32 `form:"year" binding:"required"`
	Week int32 `form:"week" binding:"required"`
}

// ListWorkingHoursResponse represents the response structure for listing working hours.
type ListWorkingHoursResponse struct {
	EmployeeID   uuid.UUID         `json:"employee_id"`
	Period       Period            `json:"period"`
	Summary      Summary           `json:"summary"`
	WorkingHours []WorkingHourItem `json:"working_hours"`
}
