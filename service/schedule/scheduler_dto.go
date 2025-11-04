package schedule

import "github.com/google/uuid"

type AutoGenerateSchedulesRequest struct {
	LocationID  int64       `json:"location_id"`
	Week        int64       `json:"week"` // e.g., "2024-W27"
	Year        int64       `json:"year"` // e.g., 2024
	EmployeeIDs []uuid.UUID `json:"employee_ids"`
}

type AutoGenerateSchedulesResponse struct{}

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
