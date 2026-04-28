package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ─── Constants ─────────────────────────────────────────────

type EventKind string

const (
	EventKindAppointment EventKind = "appointment"
	EventKindReminder    EventKind = "reminder"
)

type MutationScope string

const (
	MutationScopeSingle MutationScope = "single"
	MutationScopeSeries MutationScope = "series"
	MutationScopeFuture MutationScope = "future"
)

// ─── Errors ────────────────────────────────────────────────

// ─── Request / Response DTOs ───────────────────────────────

type ReminderInput struct {
	MinutesBefore *int32     `json:"minutes_before"`
	RemindAt      *time.Time `json:"remind_at"`
}

type CreateEventRequest struct {
	Kind                EventKind       `json:"kind" binding:"required,oneof=appointment reminder"`
	Title               string          `json:"title"`
	Description         *string         `json:"description"`
	Location            *string         `json:"location"`
	Color               *string         `json:"color"`
	StartAt             time.Time       `json:"start_at" binding:"required"`
	EndAt               time.Time       `json:"end_at" binding:"required"`
	RRule               *string         `json:"rrule"`
	AttendeeEmployeeIDs []uuid.UUID     `json:"attendee_employee_ids"`
	AttendeeClientIDs   []uuid.UUID     `json:"attendee_client_ids"`
	Reminders           []ReminderInput `json:"reminders"`
}

type ListEventsRequest struct {
	StartAt    time.Time  `json:"start_at" binding:"required"`
	EndAt      time.Time  `json:"end_at" binding:"required"`
	EmployeeID *uuid.UUID `json:"employee_id"` // Optional: filter by specific employee (requires APPOINTMENT.VIEW_ALL permission)
}

type UpdateEventRequest struct {
	Scope               MutationScope    `json:"scope" binding:"required,oneof=single series future"`
	RecurrenceID        *time.Time       `json:"recurrence_id"`
	Title               *string          `json:"title"`
	Description         *string          `json:"description"`
	Location            *string          `json:"location"`
	Color               *string          `json:"color"`
	StartAt             *time.Time       `json:"start_at"`
	EndAt               *time.Time       `json:"end_at"`
	RRule               *string          `json:"rrule"`
	AttendeeEmployeeIDs *[]uuid.UUID     `json:"attendee_employee_ids"`
	AttendeeClientIDs   *[]uuid.UUID     `json:"attendee_client_ids"`
	Reminders           *[]ReminderInput `json:"reminders"`
}

type DeleteEventRequest struct {
	Scope        MutationScope `json:"scope" binding:"required,oneof=single series future"`
	RecurrenceID *time.Time    `json:"recurrence_id"`
}

type ReminderResponse struct {
	ID            uuid.UUID  `json:"id"`
	MinutesBefore *int32     `json:"minutes_before"`
	RemindAt      *time.Time `json:"remind_at"`
}

type EventResponse struct {
	ID                  uuid.UUID          `json:"id"`
	Kind                EventKind          `json:"kind"`
	Status              string             `json:"status"`
	WorkApprovalStatus  string             `json:"work_approval_status"`
	WorkApprovedBy      *uuid.UUID         `json:"work_approved_by"`
	WorkApprovedAt      *time.Time         `json:"work_approved_at"`
	WorkRejectedBy      *uuid.UUID         `json:"work_rejected_by"`
	WorkRejectedAt      *time.Time         `json:"work_rejected_at"`
	WorkRejectionReason *string            `json:"work_rejection_reason"`
	Title               string             `json:"title"`
	Description         *string            `json:"description"`
	Location            *string            `json:"location"`
	Color               *string            `json:"color"`
	OrganizerEmployeeID uuid.UUID          `json:"organizer_employee_id"`
	StartAt             time.Time          `json:"start_at"`
	EndAt               time.Time          `json:"end_at"`
	RRule               *string            `json:"rrule"`
	RecurringEventID    *uuid.UUID         `json:"recurring_event_id"`
	RecurrenceID        *time.Time         `json:"recurrence_id"`
	AttendeeEmployeeIDs []uuid.UUID        `json:"attendee_employee_ids"`
	AttendeeClientIDs   []uuid.UUID        `json:"attendee_client_ids"`
	Reminders           []ReminderResponse `json:"reminders"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
}

type EventOccurrenceResponse struct {
	ID                  uuid.UUID   `json:"id"`
	MasterEventID       *uuid.UUID  `json:"master_event_id"`
	Kind                EventKind   `json:"kind"`
	Title               string      `json:"title"`
	Description         *string     `json:"description"`
	Location            *string     `json:"location"`
	Color               *string     `json:"color"`
	StartAt             time.Time   `json:"start_at"`
	EndAt               time.Time   `json:"end_at"`
	WorkApprovalStatus  string      `json:"work_approval_status"`
	RecurrenceID        *time.Time  `json:"recurrence_id"`
	IsRecurringInstance bool        `json:"is_recurring_instance"`
	AttendeeEmployeeIDs []uuid.UUID `json:"attendee_employee_ids"`
	AttendeeClientIDs   []uuid.UUID `json:"attendee_client_ids"`
}

type SetEventWorkApprovalRequest struct {
	// If set for a recurring master event, approval is applied to this specific occurrence
	// by creating/updating an override row (recurring_event_id + recurrence_id).
	RecurrenceID *time.Time `json:"recurrence_id"`

	// One of: pending, approved, rejected
	Status string `json:"status" binding:"required,oneof=pending approved rejected"`

	// Required when rejecting
	RejectionReason *string `json:"rejection_reason"`
}

type ListWorkApprovalQueueRequest struct {
	StartAt time.Time `json:"start_at" binding:"required"`
	EndAt   time.Time `json:"end_at" binding:"required"`

	// Optional employee filter. Matches organizer or attendee employees.
	EmployeeIDs []uuid.UUID `json:"employee_ids"`

	// If true, only include occurrences that have ended (default true).
	OnlyEnded *bool `json:"only_ended"`

	// Pagination (applied after expansion/sort).
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type IDName struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type WorkApprovalQueueItem struct {
	// For recurring occurrences, this is the master event ID to call approval endpoints with.
	EventID uuid.UUID `json:"event_id"`
	// Non-nil for recurring occurrences.
	RecurrenceID *time.Time `json:"recurrence_id"`

	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`

	OrganizerEmployeeID uuid.UUID `json:"-"`
	OrganizerEmployee   IDName    `json:"organizer_employee"`

	AttendeeEmployeeIDs []uuid.UUID `json:"-"`
	AttendeeEmployees   []IDName    `json:"attendee_employees"`

	AttendeeClientIDs []uuid.UUID `json:"-"`
	AttendeeClients   []IDName    `json:"attendee_clients"`

	WorkApprovalStatus string `json:"work_approval_status"`
	IsConfirmed        bool   `json:"is_confirmed"`

	Title       string  `json:"title"`
	Description *string `json:"description"`
	Location    *string `json:"location"`

	CreatedAt time.Time `json:"created_at"`
}

type ListWorkApprovalQueueResponse struct {
	Items []WorkApprovalQueueItem `json:"items"`
	Total int32                   `json:"total"`
}

// ─── Service interface ─────────────────────────────────────

type EventService interface {
	CreateEvent(ctx context.Context, req *CreateEventRequest, employeeID uuid.UUID) (*EventResponse, error)
	ListEvents(ctx context.Context, req ListEventsRequest, employeeID uuid.UUID) ([]EventOccurrenceResponse, error)
	GetEvent(ctx context.Context, eventID uuid.UUID, employeeID uuid.UUID) (*EventResponse, error)
	UpdateEvent(ctx context.Context, eventID uuid.UUID, req *UpdateEventRequest, employeeID uuid.UUID) (*EventResponse, error)
	DeleteEvent(ctx context.Context, eventID uuid.UUID, req DeleteEventRequest, employeeID uuid.UUID) error

	SetEventWorkApproval(ctx context.Context, eventID uuid.UUID, req *SetEventWorkApprovalRequest, actorEmployeeID, actorUserID uuid.UUID) error
	ListWorkApprovalQueue(ctx context.Context, req *ListWorkApprovalQueueRequest) (*ListWorkApprovalQueueResponse, error)
}
