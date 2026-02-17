package appointment

import (
	"time"

	"github.com/google/uuid"
)

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
	RecurrenceID        *time.Time  `json:"recurrence_id"`
	IsRecurringInstance bool        `json:"is_recurring_instance"`
	AttendeeEmployeeIDs []uuid.UUID `json:"attendee_employee_ids"`
	AttendeeClientIDs   []uuid.UUID `json:"attendee_client_ids"`
}
