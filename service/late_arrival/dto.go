package late_arrival

import (
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
)

type CreateLateArrivalRequest struct {
	ArrivalDate string `json:"arrival_date" binding:"required,datetime=2006-01-02"`
	ArrivalTime string `json:"arrival_time" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
}

type CreateLateArrivalByAdminRequest struct {
	EmployeeID  uuid.UUID `json:"employee_id" binding:"required"`
	ArrivalDate string    `json:"arrival_date" binding:"required,datetime=2006-01-02"`
	ArrivalTime string    `json:"arrival_time" binding:"required"`
	Reason      string    `json:"reason" binding:"required"`
}

type CreateLateArrivalResponse struct {
	ID                  uuid.UUID  `json:"id"`
	ScheduleID          uuid.UUID  `json:"schedule_id"`
	EmployeeID          uuid.UUID  `json:"employee_id"`
	CreatedByEmployeeID *uuid.UUID `json:"created_by_employee_id,omitempty"`
	ArrivalDate         time.Time  `json:"arrival_date"`
	ArrivalTime         string     `json:"arrival_time"`
	Reason              string     `json:"reason"`
	ShiftStartDatetime  time.Time  `json:"shift_start_datetime"`
	ShiftEndDatetime    time.Time  `json:"shift_end_datetime"`
	ShiftName           string     `json:"shift_name"`
	LocationName        string     `json:"location_name"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type ListMyLateArrivalsRequest struct {
	pagination.Request
	DateFrom *string `form:"date_from" json:"date_from" binding:"omitempty,datetime=2006-01-02"`
	DateTo   *string `form:"date_to" json:"date_to" binding:"omitempty,datetime=2006-01-02"`
}

type ListLateArrivalsRequest struct {
	pagination.Request
	EmployeeSearch *string `form:"employee_search" json:"employee_search" binding:"omitempty,max=120"`
	DateFrom       *string `form:"date_from" json:"date_from" binding:"omitempty,datetime=2006-01-02"`
	DateTo         *string `form:"date_to" json:"date_to" binding:"omitempty,datetime=2006-01-02"`
}

type LateArrivalListItem struct {
	ID                  uuid.UUID  `json:"id"`
	ScheduleID          uuid.UUID  `json:"schedule_id"`
	EmployeeID          uuid.UUID  `json:"employee_id"`
	EmployeeName        string     `json:"employee_name"`
	CreatedByEmployeeID *uuid.UUID `json:"created_by_employee_id,omitempty"`
	ArrivalDate         time.Time  `json:"arrival_date"`
	ArrivalTime         string     `json:"arrival_time"`
	Reason              string     `json:"reason"`
	ShiftStartDatetime  time.Time  `json:"shift_start_datetime"`
	ShiftEndDatetime    time.Time  `json:"shift_end_datetime"`
	ShiftName           string     `json:"shift_name"`
	LocationName        string     `json:"location_name"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
