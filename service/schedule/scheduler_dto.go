package schedule

import "github.com/google/uuid"

// AutoGenerateSchedulesRequest is the request struct for the AutoGenerateSchedules
type AutoGenerateSchedulesRequest struct {
	LocationID  uuid.UUID   `json:"location_id"`
	Week        int32       `json:"week"` // e.g., "2024-W27"
	Year        int32       `json:"year"` // e.g., 2024
	EmployeeIDs []uuid.UUID `json:"employee_ids"`
}

// AutoGenerateSchedulesResponse is the response struct for the AutoGenerateSchedules
type AutoGenerateSchedulesResponse struct {
	Status   string            `json:"status"`
	Week     int32             `json:"week"`
	Year     int32             `json:"year"`
	Shifts   []ScheduledShift  `json:"shifts"`
	GridView GridView          `json:"grid_view"`
	Summary  []EmployeeSummary `json:"summary"`
}

// SaveGeneratedSchedulesRequest is the request struct for saving generated schedules
type SaveGeneratedSchedulesRequest struct {
	LocationID      uuid.UUID        `json:"location_id"`
	ScheduledShifts []ScheduledShift `json:"scheduled_shifts"`
}

type ScheduledShift struct {
	Date      string             `json:"date"`       // e.g., "2024-07-01"
	DayName   string             `json:"day_name"`   // e.g., "Monday"
	ShiftId   uuid.UUID          `json:"shift_id"`   // e.g., 1
	ShiftName string             `json:"shift_name"` // e.g., "Morning Shift"
	StartTime string             `json:"start_time"` // e.g., "09:00"
	EndTime   string             `json:"end_time"`   // e.g., "17:00"
	Hours     float64            `json:"hours"`      // e.g., 8.0
	Employees []AssignedEmployee `json:"employees"`
}

type AssignedEmployee struct {
	EmployeeID   uuid.UUID `json:"employee_id"`
	EmployeeName string    `json:"employee_name"`
}

type GridView struct {
	Days        []string             `json:"days"`          // e.g., ["2024-07-01", "2024-07-02", ...]
	Dates       []string             `json:"dates"`         // e.g., ["Mon", "Tue", ...]
	ShiftsByDay map[string][]GridDay `json:"shifts_by_day"` // date -> list of shifts
}

type GridDay struct {
	Date   string                 `json:"date"`   // e.g., "2024-07-01"
	Shifts map[string][]GridShift `json:"shifts"` // shift_name -> list of GridShift
}

type GridShift struct {
	Employees []string `json:"employees"` // list of employee names
	Hours     float64  `json:"hours"`     // e.g., 8.0
	Start     string   `json:"start"`     // e.g., "09:00"
	End       string   `json:"end"`       // e.g., "17:00"
}

type EmployeeSummary struct {
	ID          string         `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	TargetHours float64        `json:"target_hours"`
	ActualHours float64        `json:"actual_hours"`
	Deviation   float64        `json:"deviation"`
	Status      string         `json:"status"` // "perfect", "overtime", or "undertime"
	Shifts      map[string]int `json:"shifts"` // shift type -> count
}

type Schedule struct {
	Schedule map[string]map[string][]string `json:"schedule"` // date -> employee_id -> list of shift times
	Summary  []EmployeeScheduleSummary      `json:"summary"`
}

type EmployeeScheduleSummary struct {
	Name      string         `json:"name"`
	Target    int            `json:"target"`    // target hours
	Actual    int            `json:"actual"`    // actual hours worked
	Deviation int            `json:"deviation"` // difference from target
	Status    string         `json:"status"`    // "perfect", "overtime", or "undertime"
	Shifts    map[string]int `json:"shifts"`    // shift type -> count
}
