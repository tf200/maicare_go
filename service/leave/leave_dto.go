package leave

import (
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
)

type CreateLeaveRequestRequest struct {
	LeaveType string  `json:"leave_type" binding:"required,oneof=vacation personal sick pregnancy late unpaid other"`
	StartDate string  `json:"start_date" binding:"required,datetime=2006-01-02"`
	EndDate   string  `json:"end_date" binding:"required,datetime=2006-01-02"`
	Reason    *string `json:"reason"`
}

type CreateLeaveRequestResponse struct {
	ID                  uuid.UUID  `json:"id"`
	EmployeeID          uuid.UUID  `json:"employee_id"`
	CreatedByEmployeeID *uuid.UUID `json:"created_by_employee_id,omitempty"`
	LeaveType           string     `json:"leave_type"`
	Status              string     `json:"status"`
	StartDate           time.Time  `json:"start_date"`
	EndDate             time.Time  `json:"end_date"`
	Reason              *string    `json:"reason,omitempty"`
	DecisionNote        *string    `json:"decision_note,omitempty"`
	DecidedByEmployeeID *uuid.UUID `json:"decided_by_employee_id,omitempty"`
	RequestedAt         time.Time  `json:"requested_at"`
	DecidedAt           *time.Time `json:"decided_at,omitempty"`
	CancelledAt         *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

type UpdateLeaveRequestRequest struct {
	LeaveType *string `json:"leave_type" binding:"omitempty,oneof=vacation personal sick pregnancy late unpaid other"`
	StartDate *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate   *string `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Reason    *string `json:"reason"`
}

type UpdateLeaveRequestAdminRequest struct {
	LeaveType       *string `json:"leave_type" binding:"omitempty,oneof=vacation personal sick pregnancy late unpaid other"`
	StartDate       *string `json:"start_date" binding:"omitempty,datetime=2006-01-02"`
	EndDate         *string `json:"end_date" binding:"omitempty,datetime=2006-01-02"`
	Reason          *string `json:"reason"`
	AdminUpdateNote string  `json:"admin_update_note" binding:"required"`
}

type UpdateLeaveRequestResponse = CreateLeaveRequestResponse

type ListMyLeaveRequestsRequest struct {
	pagination.Request
	Status *string `form:"status" json:"status" binding:"omitempty,oneof=pending approved rejected cancelled expired"`
}

type ListLeaveRequestsRequest struct {
	pagination.Request
	Status     *string    `form:"status" json:"status" binding:"omitempty,oneof=pending approved rejected cancelled expired"`
	EmployeeID *uuid.UUID `form:"employee_id" json:"employee_id"`
}

type LeaveRequestListItem struct {
	ID                  uuid.UUID  `json:"id"`
	EmployeeID          uuid.UUID  `json:"employee_id"`
	EmployeeName        string     `json:"employee_name"`
	CreatedByEmployeeID *uuid.UUID `json:"created_by_employee_id,omitempty"`
	LeaveType           string     `json:"leave_type"`
	Status              string     `json:"status"`
	StartDate           time.Time  `json:"start_date"`
	EndDate             time.Time  `json:"end_date"`
	Reason              *string    `json:"reason,omitempty"`
	DecisionNote        *string    `json:"decision_note,omitempty"`
	DecidedByEmployeeID *uuid.UUID `json:"decided_by_employee_id,omitempty"`
	RequestedAt         time.Time  `json:"requested_at"`
	DecidedAt           *time.Time `json:"decided_at,omitempty"`
	CancelledAt         *time.Time `json:"cancelled_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
