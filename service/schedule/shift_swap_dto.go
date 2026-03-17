package schedule

import (
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
)

type CreateShiftSwapRequest struct {
	RecipientEmployeeID uuid.UUID  `json:"recipient_employee_id" binding:"required"`
	RequesterScheduleID uuid.UUID  `json:"requester_schedule_id" binding:"required"`
	RecipientScheduleID uuid.UUID  `json:"recipient_schedule_id" binding:"required"`
	ExpiresAt           *time.Time `json:"expires_at,omitempty"`
}

type CreateShiftSwapResponse struct {
	ID                  uuid.UUID  `json:"id"`
	RequesterEmployeeID uuid.UUID  `json:"requester_employee_id"`
	RecipientEmployeeID uuid.UUID  `json:"recipient_employee_id"`
	RequesterScheduleID uuid.UUID  `json:"requester_schedule_id"`
	RecipientScheduleID uuid.UUID  `json:"recipient_schedule_id"`
	Status              string     `json:"status"`
	RequestedAt         time.Time  `json:"requested_at"`
	ExpiresAt           *time.Time `json:"expires_at,omitempty"`
	Direction           string     `json:"direction,omitempty"`
}

type RespondShiftSwapRequest struct {
	Decision string  `json:"decision" binding:"required,oneof=accept reject"`
	Note     *string `json:"note,omitempty"`
}

type AdminDecisionShiftSwapRequest struct {
	Decision string  `json:"decision" binding:"required,oneof=approve reject"`
	Note     *string `json:"note,omitempty"`
}

type ListShiftSwapRequestsRequest struct {
	pagination.Request
	Status     *string    `form:"status" json:"status" binding:"omitempty,oneof=pending_recipient recipient_rejected pending_admin admin_rejected confirmed cancelled expired"`
	EmployeeID *uuid.UUID `form:"employee_id" json:"employee_id"`
}

type ShiftSwapScheduleSnapshot struct {
	ID            uuid.UUID `json:"id"`
	EmployeeID    uuid.UUID `json:"employee_id"`
	EmployeeName  string    `json:"employee_name"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
}

type ShiftSwapResponse struct {
	ID                    uuid.UUID                 `json:"id"`
	RequesterEmployeeID   uuid.UUID                 `json:"requester_employee_id"`
	RequesterEmployeeName string                    `json:"requester_employee_name"`
	RecipientEmployeeID   uuid.UUID                 `json:"recipient_employee_id"`
	RecipientEmployeeName string                    `json:"recipient_employee_name"`
	RequesterSchedule     ShiftSwapScheduleSnapshot `json:"requester_schedule"`
	RecipientSchedule     ShiftSwapScheduleSnapshot `json:"recipient_schedule"`
	Status                string                    `json:"status"`
	RequestedAt           time.Time                 `json:"requested_at"`
	RecipientRespondedAt  *time.Time                `json:"recipient_responded_at,omitempty"`
	AdminDecidedAt        *time.Time                `json:"admin_decided_at,omitempty"`
	RecipientResponseNote *string                   `json:"recipient_response_note,omitempty"`
	AdminDecisionNote     *string                   `json:"admin_decision_note,omitempty"`
	AdminEmployeeID       *uuid.UUID                `json:"admin_employee_id,omitempty"`
	AdminEmployeeName     *string                   `json:"admin_employee_name,omitempty"`
	ExpiresAt             *time.Time                `json:"expires_at,omitempty"`
	Direction             string                    `json:"direction,omitempty"`
}
