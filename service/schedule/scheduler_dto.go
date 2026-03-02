package schedule

import "github.com/google/uuid"

// AutoGenerateSchedulesRequest is the request struct for AutoGenerateSchedules.
// Generation is allowed only for empty weeks (no schedules in DB for that location/week).
type AutoGenerateSchedulesRequest struct {
	LocationID  uuid.UUID   `json:"location_id"`
	Week        int32       `json:"week"` // ISO week number (1-53)
	Year        int32       `json:"year"` // e.g., 2026
	EmployeeIDs []uuid.UUID `json:"employee_ids"`
}

// SchedulePlanConstraints describes editor rules for the generated plan.
type SchedulePlanConstraints struct {
	MaxStaffPerShift int32 `json:"max_staff_per_shift"`
	AllowEmptyShift  bool  `json:"allow_empty_shift"`
}

type SchedulePlanEmployee struct {
	ID            uuid.UUID `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	TargetMinutes int64     `json:"target_minutes"`
}

type ScheduleShiftTemplate struct {
	ShiftID         uuid.UUID `json:"shift_id"`
	Name            string    `json:"name"`
	StartMinute     int32     `json:"start_minute"`     // minutes since 00:00 local time
	EndMinute       int32     `json:"end_minute"`       // minutes since 00:00 local time
	DurationMinutes int64     `json:"duration_minutes"` // computed duration (handles overnight)
	Overnight       bool      `json:"overnight"`
}

// SchedulePlanSlot is the canonical editable unit for the frontend.
// (date + shift_id) uniquely identifies a shift occurrence.
type SchedulePlanSlot struct {
	Date        string      `json:"date"` // "YYYY-MM-DD" in location timezone
	ShiftID     uuid.UUID   `json:"shift_id"`
	EmployeeIDs []uuid.UUID `json:"employee_ids"` // 0..MaxStaffPerShift
}

type ScheduleEmployeeSummary struct {
	EmployeeID      uuid.UUID         `json:"employee_id"`
	TargetMinutes   int64             `json:"target_minutes"`
	AssignedMinutes int64             `json:"assigned_minutes"`
	OvertimeMinutes int64             `json:"overtime_minutes"`
	ShiftCounts     map[uuid.UUID]int `json:"shift_counts"` // shift_id -> count
}

type SchedulePlanWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// AutoGenerateSchedulesResponse returns an editable draft plan.
// The frontend should render the grid from (week_start_date + shift_templates + slots).
type AutoGenerateSchedulesResponse struct {
	Status         string                    `json:"status"` // optimal, feasible, infeasible
	PlanID         uuid.UUID                 `json:"plan_id"`
	LocationID     uuid.UUID                 `json:"location_id"`
	Timezone       string                    `json:"timezone"`
	Week           int32                     `json:"week"`
	Year           int32                     `json:"year"`
	WeekStartDate  string                    `json:"week_start_date"` // "YYYY-MM-DD" in location timezone
	Constraints    SchedulePlanConstraints   `json:"constraints"`
	Employees      []SchedulePlanEmployee    `json:"employees"`
	ShiftTemplates []ScheduleShiftTemplate   `json:"shift_templates"`
	Slots          []SchedulePlanSlot        `json:"slots"`
	Summary        []ScheduleEmployeeSummary `json:"summary"`
	Warnings       []SchedulePlanWarning     `json:"warnings,omitempty"`
}

// SaveGeneratedSchedulesRequest is the request struct for saving a generated plan.
// Saving is allowed only for empty weeks (no schedules in DB for that location/week).
// The backend derives start/end datetimes from (date + shift template + location timezone).
type SaveGeneratedSchedulesRequest struct {
	PlanID     uuid.UUID          `json:"plan_id"`
	LocationID uuid.UUID          `json:"location_id"`
	Week       int32              `json:"week"`
	Year       int32              `json:"year"`
	Slots      []SchedulePlanSlot `json:"slots"`
}
