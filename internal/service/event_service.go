package service

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
	"maicare_go/pkg/conv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	rrule "github.com/teambition/rrule-go"
	"go.uber.org/zap"
)

// ─── Helpers ───────────────────────────────────────────────

func uniqueUUIDs(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]bool, len(ids))
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
}

func toPgTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: t.UTC(), Valid: true}
}

func toNullablePgTimestamptz(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return toPgTimestamptz(*t)
}

func pgTimeToPtr(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	utc := t.Time.UTC()
	return &utc
}

// ─── eventRow is a flat internal representation ────────────

type eventRow struct {
	ID                  uuid.UUID
	OrganizerEmployeeID uuid.UUID
	CreatedByEmployeeID uuid.UUID
	Kind                string
	Status              string
	WorkApprovalStatus  string
	WorkApprovedBy      *uuid.UUID
	WorkRejectedBy      *uuid.UUID
	WorkRejectionReason *string
	Title               string
	Description         *string
	Location            *string
	Color               *string
	Timezone            string
	RRule               *string
	RecurringEventID    *uuid.UUID
	RecurrenceID        *time.Time
	StartAt             time.Time
	EndAt               time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	WorkApprovedAt      *time.Time
	WorkRejectedAt      *time.Time
}

func toEventRow(event db.CalendarEvent) eventRow {
	out := eventRow{
		ID:                  event.ID,
		OrganizerEmployeeID: event.OrganizerEmployeeID,
		CreatedByEmployeeID: event.CreatedByEmployeeID,
		Kind:                string(event.Kind),
		Status:              string(event.Status),
		WorkApprovalStatus:  string(event.WorkApprovalStatus),
		WorkApprovedBy:      event.WorkApprovedBy,
		WorkRejectedBy:      event.WorkRejectedBy,
		WorkRejectionReason: event.WorkRejectionReason,
		Title:               event.Title,
		Description:         event.Description,
		Location:            event.Location,
		Color:               event.Color,
		Timezone:            event.Timezone,
		RRule:               event.Rrule,
		RecurringEventID:    event.RecurringEventID,
	}
	if event.StartAt.Valid {
		out.StartAt = event.StartAt.Time.UTC()
	}
	if event.EndAt.Valid {
		out.EndAt = event.EndAt.Time.UTC()
	}
	out.RecurrenceID = pgTimeToPtr(event.RecurrenceID)
	out.CreatedAt = conv.TimeFromPgTimestamptz(event.CreatedAt).UTC()
	out.UpdatedAt = conv.TimeFromPgTimestamptz(event.UpdatedAt).UTC()
	out.WorkApprovedAt = pgTimeToPtr(event.WorkApprovedAt)
	out.WorkRejectedAt = pgTimeToPtr(event.WorkRejectedAt)
	return out
}

// ─── Service struct ────────────────────────────────────────

type EventService struct {
	store  *db.Store
	taskQ  domain.TaskQueue
	logger domain.Logger
}

func NewEventService(store *db.Store, taskQ domain.TaskQueue, logger domain.Logger) domain.EventService {
	return &EventService{
		store:  store,
		taskQ:  taskQ,
		logger: logger,
	}
}

// ─── CreateEvent ───────────────────────────────────────────

func (s *EventService) CreateEvent(ctx context.Context, req *domain.CreateEventRequest, employeeID uuid.UUID) (*domain.EventResponse, error) {
	if !req.StartAt.Before(req.EndAt) {
		s.logger.LogWarn(ctx, "CreateEvent", "Invalid time range", zap.Time("start", req.StartAt), zap.Time("end", req.EndAt))
		return nil, fmt.Errorf("start_at must be before end_at")
	}
	if req.RRule != nil {
		if err := validateRRule(*req.RRule, req.StartAt); err != nil {
			s.logger.LogWarn(ctx, "CreateEvent", "Invalid RRule", zap.String("rrule", *req.RRule), zap.Error(err))
			return nil, err
		}
	}
	if err := validateReminders(req.Reminders); err != nil {
		s.logger.LogWarn(ctx, "CreateEvent", "Invalid reminders", zap.Error(err))
		return nil, err
	}

	tx, err := s.store.BeginActorTx(ctx)
	if err != nil {
		s.logger.LogError(ctx, "CreateEvent", "Failed to begin transaction", err)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.New(tx)

	event, err := qtx.CreateCalendarEvent(ctx, db.CreateCalendarEventParams{
		OrganizerEmployeeID: employeeID,
		CreatedByEmployeeID: employeeID,
		Kind:                db.CalendarEventKindEnum(req.Kind),
		Status:              db.CalendarEventStatusEnumConfirmed,
		Title:               req.Title,
		Description:         req.Description,
		Location:            req.Location,
		Color:               req.Color,
		StartAt:             toPgTimestamptz(req.StartAt),
		EndAt:               toPgTimestamptz(req.EndAt),
		Timezone:            "UTC",
		Rrule:               req.RRule,
	})
	if err != nil {
		s.logger.LogError(ctx, "CreateEvent", "Failed to create event in DB", err)
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	if err := upsertAttendees(ctx, qtx, event.ID, req.AttendeeEmployeeIDs, req.AttendeeClientIDs, true); err != nil {
		s.logger.LogError(ctx, "CreateEvent", "Failed to upsert attendees", err)
		return nil, err
	}

	reminders, err := upsertReminders(ctx, qtx, event.ID, req.Reminders, true)
	if err != nil {
		s.logger.LogError(ctx, "CreateEvent", "Failed to upsert reminders", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.LogError(ctx, "CreateEvent", "Failed to commit transaction", err)
		return nil, fmt.Errorf("failed to commit: %w", err)
	}

	attendeeEmpIDs := req.AttendeeEmployeeIDs
	if attendeeEmpIDs == nil {
		attendeeEmpIDs = []uuid.UUID{}
	}
	attendeeCliIDs := req.AttendeeClientIDs
	if attendeeCliIDs == nil {
		attendeeCliIDs = []uuid.UUID{}
	}

	// Build response
	resp := &domain.EventResponse{
		ID:                  event.ID,
		Kind:                domain.EventKind(event.Kind),
		Status:              string(event.Status),
		WorkApprovalStatus:  string(event.WorkApprovalStatus),
		WorkApprovedBy:      event.WorkApprovedBy,
		WorkRejectedBy:      event.WorkRejectedBy,
		WorkRejectionReason: event.WorkRejectionReason,
		Title:               event.Title,
		Description:         event.Description,
		Location:            event.Location,
		Color:               event.Color,
		OrganizerEmployeeID: event.OrganizerEmployeeID,
		StartAt:             event.StartAt.Time.UTC(),
		EndAt:               event.EndAt.Time.UTC(),
		RRule:               event.Rrule,
		RecurringEventID:    event.RecurringEventID,
		AttendeeEmployeeIDs: attendeeEmpIDs,
		AttendeeClientIDs:   attendeeCliIDs,
		Reminders:           reminders,
		CreatedAt:           event.CreatedAt.Time.UTC(),
		UpdatedAt:           event.UpdatedAt.Time.UTC(),
	}

	s.logger.LogInfo(ctx, "CreateEvent", "Event created successfully", zap.String("event_id", event.ID.String()))

	if err := s.enqueueReminderNotifications(ctx, *resp); err != nil {
		s.logger.LogWarn(ctx, "CreateEvent", "Failed to enqueue reminder notifications", zap.Error(err))
	}
	return resp, nil
}

// ─── ListEvents ────────────────────────────────────────────

func (s *EventService) ListEvents(ctx context.Context, req domain.ListEventsRequest, employeeID uuid.UUID) ([]domain.EventOccurrenceResponse, error) {
	if !req.StartAt.Before(req.EndAt) {
		s.logger.LogWarn(ctx, "ListEvents", "Invalid time range", zap.Time("start", req.StartAt), zap.Time("end", req.EndAt))
		return nil, fmt.Errorf("start_at must be before end_at")
	}

	targetEmployeeID := employeeID
	if req.EmployeeID != nil {
		targetEmployeeID = *req.EmployeeID
	}

	masters, err := s.listVisibleMasterEvents(ctx, targetEmployeeID, req.StartAt.UTC(), req.EndAt.UTC())
	if err != nil {
		s.logger.LogError(ctx, "ListEvents", "Failed to list master events", err)
		return nil, err
	}

	attendeeEmployees, attendeeClients, err := s.loadAttendeesForEvents(ctx, collectEventIDs(masters))
	if err != nil {
		s.logger.LogError(ctx, "ListEvents", "Failed to load attendees", err)
		return nil, err
	}

	seriesIDs := make([]uuid.UUID, 0)
	for _, e := range masters {
		if e.RRule != nil {
			seriesIDs = append(seriesIDs, e.ID)
		}
	}
	exceptions, err := s.loadSeriesExceptions(ctx, seriesIDs)
	if err != nil {
		s.logger.LogError(ctx, "ListEvents", "Failed to load series exceptions", err)
		return nil, err
	}
	exAttendeeEmployees, exAttendeeClients, err := s.loadAttendeesForEvents(ctx, collectEventIDs(exceptions))
	if err != nil {
		s.logger.LogError(ctx, "ListEvents", "Failed to load exception attendees", err)
		return nil, err
	}

	exMap := map[uuid.UUID]map[int64]eventRow{}
	for _, ex := range exceptions {
		if ex.RecurringEventID == nil || ex.RecurrenceID == nil {
			continue
		}
		if _, ok := exMap[*ex.RecurringEventID]; !ok {
			exMap[*ex.RecurringEventID] = map[int64]eventRow{}
		}
		exMap[*ex.RecurringEventID][ex.RecurrenceID.UTC().Unix()] = ex
	}

	items := make([]domain.EventOccurrenceResponse, 0)
	for _, e := range masters {
		if e.RRule == nil {
			items = append(items, domain.EventOccurrenceResponse{
				ID:                  e.ID,
				Kind:                domain.EventKind(e.Kind),
				Title:               e.Title,
				Description:         e.Description,
				Location:            e.Location,
				Color:               e.Color,
				StartAt:             e.StartAt.UTC(),
				EndAt:               e.EndAt.UTC(),
				WorkApprovalStatus:  e.WorkApprovalStatus,
				IsRecurringInstance: false,
				AttendeeEmployeeIDs: attendeeEmployees[e.ID],
				AttendeeClientIDs:   attendeeClients[e.ID],
			})
			continue
		}

		r, err := buildRule(*e.RRule, e.StartAt.UTC())
		if err != nil {
			continue
		}
		duration := e.EndAt.Sub(e.StartAt)
		for _, occ := range r.Between(req.StartAt.UTC(), req.EndAt.UTC(), true) {
			key := occ.UTC().Unix()
			if exForOcc, ok := exMap[e.ID][key]; ok {
				if exForOcc.Status == "cancelled" {
					continue
				}
				items = append(items, domain.EventOccurrenceResponse{
					ID:                  exForOcc.ID,
					MasterEventID:       exForOcc.RecurringEventID,
					Kind:                domain.EventKind(exForOcc.Kind),
					Title:               exForOcc.Title,
					Description:         exForOcc.Description,
					Location:            exForOcc.Location,
					Color:               exForOcc.Color,
					StartAt:             exForOcc.StartAt.UTC(),
					EndAt:               exForOcc.EndAt.UTC(),
					WorkApprovalStatus:  exForOcc.WorkApprovalStatus,
					RecurrenceID:        exForOcc.RecurrenceID,
					IsRecurringInstance: true,
					AttendeeEmployeeIDs: chooseAttendees(exAttendeeEmployees, attendeeEmployees, exForOcc.ID, e.ID),
					AttendeeClientIDs:   chooseAttendees(exAttendeeClients, attendeeClients, exForOcc.ID, e.ID),
				})
				continue
			}
			occTime := occ.UTC()
			items = append(items, domain.EventOccurrenceResponse{
				ID:                  e.ID,
				MasterEventID:       &e.ID,
				Kind:                domain.EventKind(e.Kind),
				Title:               e.Title,
				Description:         e.Description,
				Location:            e.Location,
				Color:               e.Color,
				StartAt:             occTime,
				EndAt:               occTime.Add(duration),
				WorkApprovalStatus:  string(db.CalendarEventWorkApprovalStatusEnumPending),
				RecurrenceID:        &occTime,
				IsRecurringInstance: true,
				AttendeeEmployeeIDs: attendeeEmployees[e.ID],
				AttendeeClientIDs:   attendeeClients[e.ID],
			})
		}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].StartAt.Before(items[j].StartAt)
	})
	return items, nil
}

// ─── GetEvent ──────────────────────────────────────────────

func (s *EventService) GetEvent(ctx context.Context, eventID uuid.UUID, employeeID uuid.UUID) (*domain.EventResponse, error) {
	e, err := s.getVisibleEventByID(ctx, eventID, employeeID)
	if err != nil {
		s.logger.LogError(ctx, "GetEvent", "Failed to fetch visible event", err, zap.String("event_id", eventID.String()))
		return nil, err
	}
	attendeeEmployees, attendeeClients, err := s.loadAttendeesForEvents(ctx, []uuid.UUID{eventID})
	if err != nil {
		s.logger.LogError(ctx, "GetEvent", "Failed to load attendees", err, zap.String("event_id", eventID.String()))
		return nil, err
	}
	reminders, err := s.loadRemindersForEvent(ctx, eventID)
	if err != nil {
		s.logger.LogError(ctx, "GetEvent", "Failed to load reminders", err, zap.String("event_id", eventID.String()))
		return nil, err
	}

	return &domain.EventResponse{
		ID:                  e.ID,
		Kind:                domain.EventKind(e.Kind),
		Status:              e.Status,
		WorkApprovalStatus:  e.WorkApprovalStatus,
		WorkApprovedBy:      e.WorkApprovedBy,
		WorkApprovedAt:      e.WorkApprovedAt,
		WorkRejectedBy:      e.WorkRejectedBy,
		WorkRejectedAt:      e.WorkRejectedAt,
		WorkRejectionReason: e.WorkRejectionReason,
		Title:               e.Title,
		Description:         e.Description,
		Location:            e.Location,
		Color:               e.Color,
		OrganizerEmployeeID: e.OrganizerEmployeeID,
		StartAt:             e.StartAt.UTC(),
		EndAt:               e.EndAt.UTC(),
		RRule:               e.RRule,
		RecurringEventID:    e.RecurringEventID,
		RecurrenceID:        e.RecurrenceID,
		AttendeeEmployeeIDs: attendeeEmployees[eventID],
		AttendeeClientIDs:   attendeeClients[eventID],
		Reminders:           reminders,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
	}, nil
}

// ─── UpdateEvent ───────────────────────────────────────────

func (s *EventService) UpdateEvent(ctx context.Context, eventID uuid.UUID, req *domain.UpdateEventRequest, employeeID uuid.UUID) (*domain.EventResponse, error) {
	e, err := s.getVisibleEventByID(ctx, eventID, employeeID)
	if err != nil {
		s.logger.LogError(ctx, "UpdateEvent", "Failed to fetch visible event", err, zap.String("event_id", eventID.String()))
		return nil, err
	}

	if req.Scope == domain.MutationScopeFuture {
		return s.updateEventFuture(ctx, e, req, employeeID)
	}

	if req.Scope == domain.MutationScopeSingle && e.RRule != nil {
		if req.RecurrenceID == nil {
			s.logger.LogWarn(ctx, "UpdateEvent", "Recurrence ID missing for single scope", zap.String("event_id", eventID.String()))
			return nil, fmt.Errorf("recurrence_id is required for single scope on recurring events")
		}
		return s.upsertSingleOccurrenceOverride(ctx, e, req, employeeID)
	}

	if req.StartAt != nil && req.EndAt != nil && !req.StartAt.Before(*req.EndAt) {
		s.logger.LogWarn(ctx, "UpdateEvent", "Invalid time range", zap.Time("start", *req.StartAt), zap.Time("end", *req.EndAt))
		return nil, fmt.Errorf("start_at must be before end_at")
	}

	var attendeeEmployeeIDs []uuid.UUID
	var attendeeClientIDs []uuid.UUID
	var reminders []domain.ReminderResponse

	if req.AttendeeEmployeeIDs == nil || req.AttendeeClientIDs == nil {
		empMap, clientMap, err := s.loadAttendeesForEvents(ctx, []uuid.UUID{eventID})
		if err != nil {
			s.logger.LogError(ctx, "UpdateEvent", "Failed to load existing attendees", err, zap.String("event_id", eventID.String()))
			return nil, err
		}
		if req.AttendeeEmployeeIDs == nil {
			attendeeEmployeeIDs = empMap[eventID]
		}
		if req.AttendeeClientIDs == nil {
			attendeeClientIDs = clientMap[eventID]
		}
	}
	if req.Reminders == nil {
		rems, err := s.loadRemindersForEvent(ctx, eventID)
		if err != nil {
			s.logger.LogError(ctx, "UpdateEvent", "Failed to load existing reminders", err, zap.String("event_id", eventID.String()))
			return nil, err
		}
		reminders = rems
	}

	tx, err := s.store.BeginActorTx(ctx)
	if err != nil {
		s.logger.LogError(ctx, "UpdateEvent", "Failed to begin transaction", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	if req.RRule != nil {
		if err := validateRRule(*req.RRule, e.StartAt); err != nil {
			s.logger.LogWarn(ctx, "UpdateEvent", "Invalid RRule", zap.String("rrule", *req.RRule), zap.Error(err))
			return nil, err
		}
	}

	qtx := db.New(tx)
	updatedEvent, err := qtx.UpdateCalendarEvent(ctx, db.UpdateCalendarEventParams{
		ID:          eventID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		Color:       req.Color,
		StartAt:     toNullablePgTimestamptz(req.StartAt),
		EndAt:       toNullablePgTimestamptz(req.EndAt),
		Rrule:       req.RRule,
	})
	if err != nil {
		s.logger.LogError(ctx, "UpdateEvent", "Failed to update event in DB", err)
		return nil, err
	}

	if req.AttendeeEmployeeIDs != nil || req.AttendeeClientIDs != nil {
		if req.AttendeeEmployeeIDs != nil {
			attendeeEmployeeIDs = *req.AttendeeEmployeeIDs
		}
		if req.AttendeeClientIDs != nil {
			attendeeClientIDs = *req.AttendeeClientIDs
		}
		if err := upsertAttendees(ctx, qtx, eventID, attendeeEmployeeIDs, attendeeClientIDs, false); err != nil {
			s.logger.LogError(ctx, "UpdateEvent", "Failed to upsert attendees", err)
			return nil, err
		}
		attendeeEmployeeIDs = uniqueUUIDs(attendeeEmployeeIDs)
		attendeeClientIDs = uniqueUUIDs(attendeeClientIDs)
	}

	if req.Reminders != nil {
		if err := validateReminders(*req.Reminders); err != nil {
			s.logger.LogWarn(ctx, "UpdateEvent", "Invalid reminders", zap.Error(err))
			return nil, err
		}
		rems, err := upsertReminders(ctx, qtx, eventID, *req.Reminders, false)
		if err != nil {
			s.logger.LogError(ctx, "UpdateEvent", "Failed to upsert reminders", err)
			return nil, err
		}
		reminders = rems
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.LogError(ctx, "UpdateEvent", "Failed to commit transaction", err)
		return nil, err
	}

	resp := &domain.EventResponse{
		ID:                  updatedEvent.ID,
		Kind:                domain.EventKind(updatedEvent.Kind),
		Status:              string(updatedEvent.Status),
		WorkApprovalStatus:  string(updatedEvent.WorkApprovalStatus),
		WorkApprovedBy:      updatedEvent.WorkApprovedBy,
		WorkRejectedBy:      updatedEvent.WorkRejectedBy,
		WorkRejectionReason: updatedEvent.WorkRejectionReason,
		Title:               updatedEvent.Title,
		Description:         updatedEvent.Description,
		Location:            updatedEvent.Location,
		Color:               updatedEvent.Color,
		OrganizerEmployeeID: updatedEvent.OrganizerEmployeeID,
		StartAt:             updatedEvent.StartAt.Time.UTC(),
		EndAt:               updatedEvent.EndAt.Time.UTC(),
		RRule:               updatedEvent.Rrule,
		RecurringEventID:    updatedEvent.RecurringEventID,
		AttendeeEmployeeIDs: attendeeEmployeeIDs,
		AttendeeClientIDs:   attendeeClientIDs,
		Reminders:           reminders,
		CreatedAt:           updatedEvent.CreatedAt.Time.UTC(),
		UpdatedAt:           updatedEvent.UpdatedAt.Time.UTC(),
	}
	if updatedEvent.WorkApprovedAt.Valid {
		t := updatedEvent.WorkApprovedAt.Time.UTC()
		resp.WorkApprovedAt = &t
	}
	if updatedEvent.WorkRejectedAt.Valid {
		t := updatedEvent.WorkRejectedAt.Time.UTC()
		resp.WorkRejectedAt = &t
	}
	if updatedEvent.RecurrenceID.Valid {
		t := updatedEvent.RecurrenceID.Time.UTC()
		resp.RecurrenceID = &t
	}

	s.logger.LogInfo(ctx, "UpdateEvent", "Event updated successfully", zap.String("event_id", eventID.String()))

	if err := s.enqueueReminderNotifications(ctx, *resp); err != nil {
		s.logger.LogWarn(ctx, "UpdateEvent", "Failed to enqueue reminder notifications", zap.Error(err))
	}
	return resp, nil
}

// ─── DeleteEvent ───────────────────────────────────────────

func (s *EventService) DeleteEvent(ctx context.Context, eventID uuid.UUID, req domain.DeleteEventRequest, employeeID uuid.UUID) error {
	e, err := s.getVisibleEventByID(ctx, eventID, employeeID)
	if err != nil {
		s.logger.LogError(ctx, "DeleteEvent", "Failed to fetch visible event", err, zap.String("event_id", eventID.String()))
		return err
	}

	if req.Scope == domain.MutationScopeFuture {
		if e.RRule == nil {
			return s.store.CancelCalendarEvent(ctx, eventID)
		}
		if req.RecurrenceID == nil {
			s.logger.LogWarn(ctx, "DeleteEvent", "Recurrence ID missing for future scope", zap.String("event_id", eventID.String()))
			return fmt.Errorf("recurrence_id is required for future scope")
		}
		until := req.RecurrenceID.UTC().Add(-time.Second)
		oldRule, err := withUntil(*e.RRule, e.StartAt.UTC(), until)
		if err != nil {
			s.logger.LogError(ctx, "DeleteEvent", "Failed to update RRule for future scope", err, zap.String("event_id", eventID.String()))
			return err
		}
		err = s.store.UpdateCalendarEventRRule(ctx, db.UpdateCalendarEventRRuleParams{ID: e.ID, Rrule: &oldRule})
		if err != nil {
			s.logger.LogError(ctx, "DeleteEvent", "Failed to update RRule in DB", err, zap.String("event_id", eventID.String()))
			return err
		}
		s.logger.LogInfo(ctx, "DeleteEvent", "Event series truncated for future scope", zap.String("event_id", eventID.String()))
		return nil
	}

	if req.Scope == domain.MutationScopeSingle && e.RRule != nil {
		if req.RecurrenceID == nil {
			s.logger.LogWarn(ctx, "DeleteEvent", "Recurrence ID missing for single scope", zap.String("event_id", eventID.String()))
			return fmt.Errorf("recurrence_id is required for single scope")
		}
		_, err := s.store.UpsertCalendarEventOverride(ctx, db.UpsertCalendarEventOverrideParams{
			OrganizerEmployeeID: e.OrganizerEmployeeID,
			CreatedByEmployeeID: employeeID,
			Kind:                db.CalendarEventKindEnum(e.Kind),
			Status:              db.CalendarEventStatusEnumCancelled,
			Title:               e.Title,
			Description:         e.Description,
			Location:            e.Location,
			Color:               e.Color,
			StartAt:             toPgTimestamptz(*req.RecurrenceID),
			EndAt:               toPgTimestamptz(req.RecurrenceID.UTC().Add(e.EndAt.Sub(e.StartAt))),
			Timezone:            "UTC",
			RecurringEventID:    &e.ID,
			RecurrenceID:        toPgTimestamptz(*req.RecurrenceID),
		})
		if err != nil {
			s.logger.LogError(ctx, "DeleteEvent", "Failed to upsert cancelled override", err, zap.String("event_id", eventID.String()))
			return err
		}
		s.logger.LogInfo(ctx, "DeleteEvent", "Single occurrence cancelled", zap.String("event_id", eventID.String()), zap.Time("recurrence_id", *req.RecurrenceID))
		return nil
	}

	err = s.store.CancelCalendarEvent(ctx, eventID)
	if err != nil {
		s.logger.LogError(ctx, "DeleteEvent", "Failed to cancel calendar event", err, zap.String("event_id", eventID.String()))
		return err
	}
	s.logger.LogInfo(ctx, "DeleteEvent", "Event cancelled successfully", zap.String("event_id", eventID.String()))
	return nil
}

// ─── SetEventWorkApproval ──────────────────────────────────

func (s *EventService) SetEventWorkApproval(ctx context.Context, eventID uuid.UUID, req *domain.SetEventWorkApprovalRequest, actorEmployeeID, actorUserID uuid.UUID) error {
	event, err := s.store.GetCalendarEventByID(ctx, eventID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("event not found")
		}
		s.logger.LogError(ctx, "SetEventWorkApproval", "Failed to get event by ID", err, zap.String("event_id", eventID.String()))
		return fmt.Errorf("failed to get event: %w", err)
	}

	targetEventID := eventID

	tx, err := s.store.BeginActorTx(ctx)
	if err != nil {
		s.logger.LogError(ctx, "SetEventWorkApproval", "Failed to begin transaction", err)
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := db.New(tx)

	if req.RecurrenceID != nil {
		baseStart := req.RecurrenceID.UTC()
		baseEnd := baseStart.Add(event.EndAt.Time.Sub(event.StartAt.Time))

		override, err := qtx.UpsertCalendarEventOverride(ctx, db.UpsertCalendarEventOverrideParams{
			OrganizerEmployeeID: event.OrganizerEmployeeID,
			CreatedByEmployeeID: actorEmployeeID,
			Kind:                event.Kind,
			Status:              event.Status,
			Title:               event.Title,
			Description:         event.Description,
			Location:            event.Location,
			Color:               event.Color,
			StartAt:             toPgTimestamptz(baseStart),
			EndAt:               toPgTimestamptz(baseEnd),
			Timezone:            event.Timezone,
			RecurringEventID:    &event.ID,
			RecurrenceID:        toPgTimestamptz(*req.RecurrenceID),
		})
		if err != nil {
			s.logger.LogError(ctx, "SetEventWorkApproval", "Failed to upsert occurrence override", err)
			return fmt.Errorf("failed to upsert occurrence override: %w", err)
		}

		attRows, err := qtx.ListAttendeesByEventIDs(ctx, []uuid.UUID{event.ID})
		if err != nil {
			s.logger.LogError(ctx, "SetEventWorkApproval", "Failed to load master attendees", err)
			return fmt.Errorf("failed to load master attendees: %w", err)
		}
		employeeIDs := make([]uuid.UUID, 0)
		clientIDs := make([]uuid.UUID, 0)
		for _, row := range attRows {
			if row.EmployeeID != nil {
				employeeIDs = append(employeeIDs, *row.EmployeeID)
			}
			if row.ClientID != nil {
				clientIDs = append(clientIDs, *row.ClientID)
			}
		}
		if err := upsertAttendees(ctx, qtx, override.ID, employeeIDs, clientIDs, true); err != nil {
			s.logger.LogError(ctx, "SetEventWorkApproval", "Failed to copy attendees", err)
			return fmt.Errorf("failed to copy attendees to occurrence override: %w", err)
		}
		targetEventID = override.ID
	}

	var rejectionReason *string
	if req.Status == "rejected" {
		rejectionReason = req.RejectionReason
	}

	err = qtx.UpdateCalendarEventWorkApproval(ctx, db.UpdateCalendarEventWorkApprovalParams{
		EventID:            targetEventID,
		WorkApprovalStatus: db.CalendarEventWorkApprovalStatusEnum(req.Status),
		ActorUserID:        &actorUserID,
		RejectionReason:    rejectionReason,
	})
	if err != nil {
		s.logger.LogError(ctx, "SetEventWorkApproval", "Failed to update work approval", err)
		return fmt.Errorf("failed to update work approval: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.LogError(ctx, "SetEventWorkApproval", "Failed to commit transaction", err)
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	s.logger.LogInfo(ctx, "SetEventWorkApproval", "Work approval updated successfully", zap.String("event_id", targetEventID.String()), zap.String("status", string(req.Status)))
	return nil
}

// ─── ListWorkApprovalQueue ─────────────────────────────────

func (s *EventService) ListWorkApprovalQueue(ctx context.Context, req *domain.ListWorkApprovalQueueRequest) (*domain.ListWorkApprovalQueueResponse, error) {
	const maxWorkApprovalQueueRange = 90 * 24 * time.Hour

	if !req.StartAt.Before(req.EndAt) {
		s.logger.LogWarn(ctx, "ListWorkApprovalQueue", "Invalid time range", zap.Time("start", req.StartAt), zap.Time("end", req.EndAt))
		return nil, fmt.Errorf("start_at must be before end_at")
	}
	onlyEnded := true
	if req.OnlyEnded != nil {
		onlyEnded = *req.OnlyEnded
	}
	limit := req.Limit
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	offset := req.Offset
	if offset < 0 {
		offset = 0
	}
	employeeFilter := map[uuid.UUID]bool{}
	for _, id := range req.EmployeeIDs {
		employeeFilter[id] = true
	}
	hasEmployeeFilter := len(employeeFilter) > 0

	startAt := req.StartAt.UTC()
	endAt := req.EndAt.UTC()
	now := time.Now().UTC()

	effectiveEnd := endAt
	if onlyEnded && effectiveEnd.After(now) {
		effectiveEnd = now
	}
	if effectiveEnd.Sub(startAt) > maxWorkApprovalQueueRange {
		s.logger.LogWarn(ctx, "ListWorkApprovalQueue", "Time range too large",
			zap.Time("start", startAt),
			zap.Time("end", endAt),
			zap.Time("effective_end", effectiveEnd),
			zap.Duration("range", effectiveEnd.Sub(startAt)),
		)
		return nil, fmt.Errorf("time range too large; maximum is 90 days")
	}
	if !startAt.Before(effectiveEnd) {
		return &domain.ListWorkApprovalQueueResponse{Items: []domain.WorkApprovalQueueItem{}, Total: 0}, nil
	}

	var oneOff []db.CalendarEvent
	var masters []db.CalendarEvent
	var exRows []db.CalendarEvent
	err := s.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		oneOff, err = q.ListWorkApprovalQueueOneOffAppointmentsStartingInRange(ctx, db.ListWorkApprovalQueueOneOffAppointmentsStartingInRangeParams{
			StartAt:     toPgTimestamptz(startAt),
			EndAt:       toPgTimestamptz(effectiveEnd),
			EmployeeIds: req.EmployeeIDs,
		})
		if err != nil {
			return err
		}
		masters, err = q.ListWorkApprovalQueueRecurringMastersStartingBeforeEnd(ctx, db.ListWorkApprovalQueueRecurringMastersStartingBeforeEndParams{
			EndAt:       toPgTimestamptz(effectiveEnd),
			EmployeeIds: req.EmployeeIDs,
		})
		if err != nil {
			return err
		}
		seriesIDs := make([]uuid.UUID, 0, len(masters))
		for _, master := range masters {
			seriesIDs = append(seriesIDs, master.ID)
		}
		if len(seriesIDs) > 0 {
			exRows, err = q.ListSeriesExceptionsStartingInRange(ctx, db.ListSeriesExceptionsStartingInRangeParams{
				SeriesIds: seriesIDs, StartAt: toPgTimestamptz(startAt), EndAt: toPgTimestamptz(effectiveEnd), EmployeeIds: req.EmployeeIDs,
			})
		}
		return err
	})
	if err != nil {
		s.logger.LogError(ctx, "ListWorkApprovalQueue", "Failed to list work approval events", err)
		return nil, fmt.Errorf("failed to list work approval events: %w", err)
	}

	masterRows := make([]eventRow, 0, len(masters))
	seriesIDs := make([]uuid.UUID, 0, len(masters))
	for _, m := range masters {
		masterRows = append(masterRows, toEventRow(m))
		seriesIDs = append(seriesIDs, m.ID)
	}

	exceptions := []eventRow{}
	for _, ex := range exRows {
		exceptions = append(exceptions, toEventRow(ex))
	}

	eventIDs := make([]uuid.UUID, 0, len(oneOff)+len(seriesIDs)+len(exceptions))
	for _, e := range oneOff {
		eventIDs = append(eventIDs, e.ID)
	}
	eventIDs = append(eventIDs, seriesIDs...)
	for _, ex := range exceptions {
		eventIDs = append(eventIDs, ex.ID)
	}
	attendeeEmployees, attendeeClients, err := s.loadAttendeesForEvents(ctx, eventIDs)
	if err != nil {
		s.logger.LogError(ctx, "ListWorkApprovalQueue", "Failed to load attendees", err)
		return nil, fmt.Errorf("failed to load attendees: %w", err)
	}

	exMap := map[uuid.UUID]map[int64]eventRow{}
	for _, ex := range exceptions {
		if ex.RecurringEventID == nil || ex.RecurrenceID == nil {
			continue
		}
		if _, ok := exMap[*ex.RecurringEventID]; !ok {
			exMap[*ex.RecurringEventID] = map[int64]eventRow{}
		}
		exMap[*ex.RecurringEventID][ex.RecurrenceID.UTC().Unix()] = ex
	}

	matchesEmployeeFilter := func(organizerID uuid.UUID, attendeeEmployeeIDs []uuid.UUID) bool {
		if !hasEmployeeFilter {
			return true
		}
		if employeeFilter[organizerID] {
			return true
		}
		for _, attendeeID := range attendeeEmployeeIDs {
			if employeeFilter[attendeeID] {
				return true
			}
		}
		return false
	}

	filteredMasters := make([]eventRow, 0, len(masterRows))
	for _, m := range masterRows {
		if len(attendeeClients[m.ID]) == 0 {
			continue
		}
		if !matchesEmployeeFilter(m.OrganizerEmployeeID, attendeeEmployees[m.ID]) {
			continue
		}
		filteredMasters = append(filteredMasters, m)
	}
	masterRows = filteredMasters

	items := make([]domain.WorkApprovalQueueItem, 0, len(oneOff)+len(masterRows)*4)
	addedRecurring := map[uuid.UUID]map[int64]bool{}

	appendItem := func(
		eventID uuid.UUID,
		recurrenceID *time.Time,
		start time.Time,
		end time.Time,
		organizerEmployeeID uuid.UUID,
		attendeeEmployeeIDs []uuid.UUID,
		attendeeClientIDs []uuid.UUID,
		workApprovalStatus string,
		title string,
		description *string,
		location *string,
		createdAt time.Time,
	) {
		items = append(items, domain.WorkApprovalQueueItem{
			EventID:             eventID,
			RecurrenceID:        recurrenceID,
			StartAt:             start,
			EndAt:               end,
			OrganizerEmployeeID: organizerEmployeeID,
			AttendeeEmployeeIDs: attendeeEmployeeIDs,
			AttendeeClientIDs:   attendeeClientIDs,
			WorkApprovalStatus:  workApprovalStatus,
			IsConfirmed:         workApprovalStatus == string(db.CalendarEventWorkApprovalStatusEnumApproved),
			Title:               title,
			Description:         description,
			Location:            location,
			CreatedAt:           createdAt,
		})
	}

	for _, row := range oneOff {
		e := toEventRow(row)
		if onlyEnded && e.EndAt.After(now) {
			continue
		}
		clientIDs := attendeeClients[e.ID]
		if len(clientIDs) == 0 {
			continue
		}
		empIDs := attendeeEmployees[e.ID]
		if !matchesEmployeeFilter(e.OrganizerEmployeeID, empIDs) {
			continue
		}
		appendItem(e.ID, nil, e.StartAt, e.EndAt, e.OrganizerEmployeeID, empIDs, clientIDs, e.WorkApprovalStatus, e.Title, e.Description, e.Location, e.CreatedAt)
	}

	for _, m := range masterRows {
		if m.RRule == nil {
			continue
		}
		r, err := buildRule(*m.RRule, m.StartAt.UTC())
		if err != nil {
			continue
		}
		duration := m.EndAt.Sub(m.StartAt)

		rruleEnd := effectiveEnd
		next := r.Iterator()
		for {
			occ, ok := next()
			if !ok {
				break
			}
			occurrenceID := occ.UTC()
			if occurrenceID.Before(startAt) {
				continue
			}
			if occurrenceID.After(rruleEnd) {
				break
			}

			key := occurrenceID.Unix()

			effective := m
			isCancelled := false
			if exForOcc, ok := exMap[m.ID][key]; ok {
				if exForOcc.Status == "cancelled" {
					isCancelled = true
				} else {
					effective = exForOcc
				}
			}
			if isCancelled {
				continue
			}

			effectiveStart := occurrenceID
			effectiveOccEnd := occurrenceID.Add(duration)
			effectiveStatus := string(db.CalendarEventWorkApprovalStatusEnumPending)

			if effective.ID != m.ID {
				effectiveStart = effective.StartAt.UTC()
				effectiveOccEnd = effective.EndAt.UTC()
				effectiveStatus = effective.WorkApprovalStatus
			}

			if effectiveStart.Before(startAt) || !effectiveStart.Before(endAt) {
				continue
			}
			if onlyEnded && effectiveOccEnd.After(now) {
				continue
			}

			overrideID := effective.ID
			clientIDs := chooseAttendees(attendeeClients, attendeeClients, overrideID, m.ID)
			if len(clientIDs) == 0 {
				continue
			}
			empIDs := chooseAttendees(attendeeEmployees, attendeeEmployees, overrideID, m.ID)
			if !matchesEmployeeFilter(m.OrganizerEmployeeID, empIDs) {
				continue
			}

			if _, ok := addedRecurring[m.ID]; !ok {
				addedRecurring[m.ID] = map[int64]bool{}
			}
			addedRecurring[m.ID][key] = true

			appendItem(m.ID, &occurrenceID, effectiveStart, effectiveOccEnd, m.OrganizerEmployeeID, empIDs, clientIDs, effectiveStatus, effective.Title, effective.Description, effective.Location, m.CreatedAt)
		}
	}

	for _, ex := range exceptions {
		if ex.RecurringEventID == nil || ex.RecurrenceID == nil {
			continue
		}
		if ex.Status == "cancelled" {
			continue
		}
		if ex.StartAt.Before(startAt) || !ex.StartAt.Before(endAt) {
			continue
		}
		if onlyEnded && ex.EndAt.After(now) {
			continue
		}
		key := ex.RecurrenceID.UTC().Unix()
		if addedRecurring[*ex.RecurringEventID] != nil && addedRecurring[*ex.RecurringEventID][key] {
			continue
		}

		clientIDs := chooseAttendees(attendeeClients, attendeeClients, ex.ID, *ex.RecurringEventID)
		if len(clientIDs) == 0 {
			continue
		}
		empIDs := chooseAttendees(attendeeEmployees, attendeeEmployees, ex.ID, *ex.RecurringEventID)
		if !matchesEmployeeFilter(ex.OrganizerEmployeeID, empIDs) {
			continue
		}

		occID := ex.RecurrenceID.UTC()
		appendItem(*ex.RecurringEventID, &occID, ex.StartAt.UTC(), ex.EndAt.UTC(), ex.OrganizerEmployeeID, empIDs, clientIDs, ex.WorkApprovalStatus, ex.Title, ex.Description, ex.Location, ex.CreatedAt)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].StartAt.After(items[j].StartAt)
	})

	total := int32(len(items))
	start := int(offset)
	if start > len(items) {
		start = len(items)
	}
	end := start + int(limit)
	if end > len(items) {
		end = len(items)
	}
	paged := items[start:end]

	employeeIDSet := map[uuid.UUID]bool{}
	clientIDSet := map[uuid.UUID]bool{}
	for i := range paged {
		employeeIDSet[paged[i].OrganizerEmployeeID] = true
		for _, eid := range paged[i].AttendeeEmployeeIDs {
			employeeIDSet[eid] = true
		}
		for _, cid := range paged[i].AttendeeClientIDs {
			clientIDSet[cid] = true
		}
	}

	employeeIDs := make([]uuid.UUID, 0, len(employeeIDSet))
	for id := range employeeIDSet {
		employeeIDs = append(employeeIDs, id)
	}
	clientIDs := make([]uuid.UUID, 0, len(clientIDSet))
	for id := range clientIDSet {
		clientIDs = append(clientIDs, id)
	}

	employeeNames := map[uuid.UUID]string{}
	if len(employeeIDs) > 0 {
		rows, err := s.store.ListEmployeeNamesByIDs(ctx, employeeIDs)
		if err != nil {
			s.logger.LogError(ctx, "ListWorkApprovalQueue", "Failed to list employee names", err)
			return nil, fmt.Errorf("failed to list employee names: %w", err)
		}
		for _, r := range rows {
			name := strings.TrimSpace(strings.TrimSpace(r.FirstName) + " " + strings.TrimSpace(r.LastName))
			employeeNames[r.ID] = name
		}
	}
	clientNames := map[uuid.UUID]string{}
	if len(clientIDs) > 0 {
		rows, err := actorQuery(ctx, s.store, func(q *db.Queries) ([]db.ListClientNamesByIDsRow, error) {
			return q.ListClientNamesByIDs(ctx, clientIDs)
		})
		if err != nil {
			s.logger.LogError(ctx, "ListWorkApprovalQueue", "Failed to list client names", err)
			return nil, fmt.Errorf("failed to list client names: %w", err)
		}
		for _, r := range rows {
			name := strings.TrimSpace(strings.TrimSpace(r.FirstName) + " " + strings.TrimSpace(r.LastName))
			clientNames[r.ID] = name
		}
	}

	for i := range paged {
		paged[i].OrganizerEmployee = domain.IDName{
			ID:   paged[i].OrganizerEmployeeID,
			Name: employeeNames[paged[i].OrganizerEmployeeID],
		}

		paged[i].AttendeeEmployees = make([]domain.IDName, 0, len(paged[i].AttendeeEmployeeIDs))
		for _, eid := range paged[i].AttendeeEmployeeIDs {
			paged[i].AttendeeEmployees = append(paged[i].AttendeeEmployees, domain.IDName{ID: eid, Name: employeeNames[eid]})
		}

		paged[i].AttendeeClients = make([]domain.IDName, 0, len(paged[i].AttendeeClientIDs))
		for _, cid := range paged[i].AttendeeClientIDs {
			paged[i].AttendeeClients = append(paged[i].AttendeeClients, domain.IDName{ID: cid, Name: clientNames[cid]})
		}
	}

	return &domain.ListWorkApprovalQueueResponse{
		Items: paged,
		Total: total,
	}, nil
}

// ─── Private helpers ──────────────────────────────────────

func (s *EventService) getVisibleEventByID(ctx context.Context, eventID, employeeID uuid.UUID) (eventRow, error) {
	e, err := actorQuery(ctx, s.store, func(q *db.Queries) (db.CalendarEvent, error) {
		return q.GetVisibleEventByID(ctx, db.GetVisibleEventByIDParams{ID: eventID, EmployeeID: employeeID})
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return eventRow{}, fmt.Errorf("event not found")
		}
		return eventRow{}, err
	}
	return toEventRow(e), nil
}

func (s *EventService) listVisibleMasterEvents(ctx context.Context, employeeID uuid.UUID, startAt, endAt time.Time) ([]eventRow, error) {
	rows, err := actorQuery(ctx, s.store, func(q *db.Queries) ([]db.CalendarEvent, error) {
		return q.ListVisibleMasterEvents(ctx, db.ListVisibleMasterEventsParams{
			EmployeeID: employeeID,
			StartAt:    toPgTimestamptz(startAt),
			EndAt:      toPgTimestamptz(endAt),
		})
	})
	if err != nil {
		return nil, err
	}
	out := make([]eventRow, 0)
	for _, row := range rows {
		out = append(out, toEventRow(row))
	}
	return out, nil
}

func (s *EventService) loadSeriesExceptions(ctx context.Context, seriesIDs []uuid.UUID) ([]eventRow, error) {
	if len(seriesIDs) == 0 {
		return []eventRow{}, nil
	}
	rows, err := s.store.ListSeriesExceptions(ctx, seriesIDs)
	if err != nil {
		return nil, err
	}
	out := make([]eventRow, 0)
	for _, row := range rows {
		out = append(out, toEventRow(row))
	}
	return out, nil
}

func (s *EventService) loadAttendeesForEvents(ctx context.Context, eventIDs []uuid.UUID) (map[uuid.UUID][]uuid.UUID, map[uuid.UUID][]uuid.UUID, error) {
	var employeeMap, clientMap map[uuid.UUID][]uuid.UUID
	err := s.store.ExecActorTx(ctx, func(q *db.Queries) error {
		var err error
		employeeMap, clientMap, err = loadAttendeesWithQueries(ctx, q, eventIDs)
		return err
	})
	return employeeMap, clientMap, err
}

func loadAttendeesWithQueries(ctx context.Context, q *db.Queries, eventIDs []uuid.UUID) (map[uuid.UUID][]uuid.UUID, map[uuid.UUID][]uuid.UUID, error) {
	employeeMap := map[uuid.UUID][]uuid.UUID{}
	clientMap := map[uuid.UUID][]uuid.UUID{}
	if len(eventIDs) == 0 {
		return employeeMap, clientMap, nil
	}

	rows, err := q.ListAttendeesByEventIDs(ctx, eventIDs)
	if err != nil {
		return nil, nil, err
	}

	for _, row := range rows {
		if row.EmployeeID != nil {
			employeeMap[row.EventID] = append(employeeMap[row.EventID], *row.EmployeeID)
		}
		if row.ClientID != nil {
			clientMap[row.EventID] = append(clientMap[row.EventID], *row.ClientID)
		}
	}
	return employeeMap, clientMap, nil
}

func (s *EventService) loadRemindersForEvent(ctx context.Context, eventID uuid.UUID) ([]domain.ReminderResponse, error) {
	return loadRemindersWithQueries(ctx, s.store.Queries, eventID)
}

func loadRemindersWithQueries(ctx context.Context, q *db.Queries, eventID uuid.UUID) ([]domain.ReminderResponse, error) {
	rows, err := q.ListRemindersByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.ReminderResponse, 0)
	for _, row := range rows {
		var r domain.ReminderResponse
		r.ID = row.ID
		r.MinutesBefore = row.MinutesBefore
		r.RemindAt = pgTimeToPtr(row.RemindAt)
		out = append(out, r)
	}
	return out, nil
}

// ─── Future scope helper ───────────────────────────────────

func (s *EventService) updateEventFuture(ctx context.Context, master eventRow, req *domain.UpdateEventRequest, employeeID uuid.UUID) (*domain.EventResponse, error) {
	if master.RRule == nil {
		s.logger.LogWarn(ctx, "updateEventFuture", "Future scope requires recurring event", zap.String("event_id", master.ID.String()))
		return nil, fmt.Errorf("future scope requires recurring event")
	}
	if req.RecurrenceID == nil {
		s.logger.LogWarn(ctx, "updateEventFuture", "Recurrence ID missing", zap.String("event_id", master.ID.String()))
		return nil, fmt.Errorf("recurrence_id is required for future scope")
	}

	futureStart := req.RecurrenceID.UTC()
	futureEnd := futureStart.Add(master.EndAt.Sub(master.StartAt))
	if req.StartAt != nil {
		futureStart = req.StartAt.UTC()
	}
	if req.EndAt != nil {
		futureEnd = req.EndAt.UTC()
	}
	if !futureStart.Before(futureEnd) {
		s.logger.LogWarn(ctx, "updateEventFuture", "Invalid time range", zap.Time("start", futureStart), zap.Time("end", futureEnd))
		return nil, fmt.Errorf("start_at must be before end_at")
	}

	title := master.Title
	if req.Title != nil {
		title = *req.Title
	}
	description := master.Description
	if req.Description != nil {
		description = req.Description
	}
	location := master.Location
	if req.Location != nil {
		location = req.Location
	}
	color := master.Color
	if req.Color != nil {
		color = req.Color
	}
	newRule := *master.RRule
	if req.RRule != nil {
		if err := validateRRule(*req.RRule, futureStart); err != nil {
			s.logger.LogWarn(ctx, "updateEventFuture", "Invalid new RRule", zap.Error(err))
			return nil, err
		}
		newRule = *req.RRule
	}

	until := req.RecurrenceID.UTC().Add(-time.Second)
	oldRule, err := withUntil(*master.RRule, master.StartAt.UTC(), until)
	if err != nil {
		s.logger.LogError(ctx, "updateEventFuture", "Failed to calculate old rule end", err)
		return nil, err
	}

	tx, err := s.store.BeginActorTx(ctx)
	if err != nil {
		s.logger.LogError(ctx, "updateEventFuture", "Failed to begin transaction", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := db.New(tx)
	err = qtx.UpdateCalendarEventRRule(ctx, db.UpdateCalendarEventRRuleParams{ID: master.ID, Rrule: &oldRule})
	if err != nil {
		s.logger.LogError(ctx, "updateEventFuture", "Failed to update old rule in DB", err)
		return nil, err
	}

	newMaster, err := qtx.CreateCalendarEvent(ctx, db.CreateCalendarEventParams{
		OrganizerEmployeeID: master.OrganizerEmployeeID,
		CreatedByEmployeeID: employeeID,
		Kind:                db.CalendarEventKindEnum(master.Kind),
		Status:              db.CalendarEventStatusEnum(master.Status),
		Title:               title,
		Description:         description,
		Location:            location,
		Color:               color,
		StartAt:             toPgTimestamptz(futureStart),
		EndAt:               toPgTimestamptz(futureEnd),
		Timezone:            "UTC",
		Rrule:               &newRule,
		RecurringEventID:    nil,
		RecurrenceID:        pgtype.Timestamptz{},
	})
	if err != nil {
		s.logger.LogError(ctx, "updateEventFuture", "Failed to create new master event", err)
		return nil, err
	}

	attendeeEmployees, attendeeClients, err := loadAttendeesWithQueries(ctx, qtx, []uuid.UUID{master.ID})
	if err != nil {
		s.logger.LogError(ctx, "updateEventFuture", "Failed to load master attendees", err)
		return nil, err
	}
	employeeIDs := attendeeEmployees[master.ID]
	clientIDs := attendeeClients[master.ID]
	if req.AttendeeEmployeeIDs != nil {
		employeeIDs = *req.AttendeeEmployeeIDs
	}
	if req.AttendeeClientIDs != nil {
		clientIDs = *req.AttendeeClientIDs
	}
	if err := upsertAttendees(ctx, qtx, newMaster.ID, employeeIDs, clientIDs, true); err != nil {
		s.logger.LogError(ctx, "updateEventFuture", "Failed to upsert attendees", err)
		return nil, err
	}

	reminders, err := loadRemindersWithQueries(ctx, qtx, master.ID)
	if err != nil {
		s.logger.LogError(ctx, "updateEventFuture", "Failed to load master reminders", err)
		return nil, err
	}
	reminderInputs := make([]domain.ReminderInput, 0, len(reminders))
	for _, reminder := range reminders {
		reminderInputs = append(reminderInputs, domain.ReminderInput{MinutesBefore: reminder.MinutesBefore, RemindAt: reminder.RemindAt})
	}
	if req.Reminders != nil {
		if err := validateReminders(*req.Reminders); err != nil {
			return nil, err
		}
		reminderInputs = *req.Reminders
	}
	insertedReminders, err := upsertReminders(ctx, qtx, newMaster.ID, reminderInputs, true)
	if err != nil {
		s.logger.LogError(ctx, "updateEventFuture", "Failed to upsert reminders", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.logger.LogError(ctx, "updateEventFuture", "Failed to commit transaction", err)
		return nil, err
	}

	resp := &domain.EventResponse{
		ID:                  newMaster.ID,
		Kind:                domain.EventKind(newMaster.Kind),
		Status:              string(newMaster.Status),
		WorkApprovalStatus:  string(newMaster.WorkApprovalStatus),
		WorkApprovedBy:      newMaster.WorkApprovedBy,
		WorkRejectedBy:      newMaster.WorkRejectedBy,
		WorkRejectionReason: newMaster.WorkRejectionReason,
		Title:               newMaster.Title,
		Description:         newMaster.Description,
		Location:            newMaster.Location,
		Color:               newMaster.Color,
		OrganizerEmployeeID: newMaster.OrganizerEmployeeID,
		StartAt:             newMaster.StartAt.Time.UTC(),
		EndAt:               newMaster.EndAt.Time.UTC(),
		RRule:               newMaster.Rrule,
		RecurringEventID:    newMaster.RecurringEventID,
		AttendeeEmployeeIDs: uniqueUUIDs(employeeIDs),
		AttendeeClientIDs:   uniqueUUIDs(clientIDs),
		Reminders:           insertedReminders,
		CreatedAt:           newMaster.CreatedAt.Time.UTC(),
		UpdatedAt:           newMaster.UpdatedAt.Time.UTC(),
	}

	s.logger.LogInfo(ctx, "updateEventFuture", "Future events split into new master successfully", zap.String("old_master_id", master.ID.String()), zap.String("new_master_id", newMaster.ID.String()))

	if err := s.enqueueReminderNotifications(ctx, *resp); err != nil {
		s.logger.LogWarn(ctx, "updateEventFuture", "Failed to enqueue reminder notifications", zap.Error(err))
	}
	return resp, nil
}

// ─── Single occurrence override helper ─────────────────────

func (s *EventService) upsertSingleOccurrenceOverride(ctx context.Context, master eventRow, req *domain.UpdateEventRequest, employeeID uuid.UUID) (*domain.EventResponse, error) {
	baseStart := req.RecurrenceID.UTC()
	baseEnd := baseStart.Add(master.EndAt.Sub(master.StartAt))
	if req.StartAt != nil {
		baseStart = req.StartAt.UTC()
	}
	if req.EndAt != nil {
		baseEnd = req.EndAt.UTC()
	}
	if !baseStart.Before(baseEnd) {
		s.logger.LogWarn(ctx, "upsertSingleOccurrenceOverride", "Invalid time range", zap.Time("start", baseStart), zap.Time("end", baseEnd))
		return nil, fmt.Errorf("start_at must be before end_at")
	}

	title := master.Title
	if req.Title != nil {
		title = *req.Title
	}
	description := master.Description
	if req.Description != nil {
		description = req.Description
	}
	location := master.Location
	if req.Location != nil {
		location = req.Location
	}
	color := master.Color
	if req.Color != nil {
		color = req.Color
	}

	override, err := s.store.UpsertCalendarEventOverride(ctx, db.UpsertCalendarEventOverrideParams{
		OrganizerEmployeeID: master.OrganizerEmployeeID,
		CreatedByEmployeeID: employeeID,
		Kind:                db.CalendarEventKindEnum(master.Kind),
		Status:              db.CalendarEventStatusEnum(master.Status),
		Title:               title,
		Description:         description,
		Location:            location,
		Color:               color,
		StartAt:             toPgTimestamptz(baseStart),
		EndAt:               toPgTimestamptz(baseEnd),
		Timezone:            "UTC",
		RecurringEventID:    &master.ID,
		RecurrenceID:        toPgTimestamptz(*req.RecurrenceID),
	})
	if err != nil {
		s.logger.LogError(ctx, "upsertSingleOccurrenceOverride", "Failed to upsert override in DB", err)
		return nil, err
	}

	var finalEmployeeIDs []uuid.UUID
	var finalClientIDs []uuid.UUID

	if req.AttendeeEmployeeIDs != nil || req.AttendeeClientIDs != nil {
		tx, err := s.store.BeginActorTx(ctx)
		if err != nil {
			s.logger.LogError(ctx, "upsertSingleOccurrenceOverride", "Failed to begin transaction", err)
			return nil, err
		}
		defer tx.Rollback(ctx)
		qtx := db.New(tx)
		employeeMap, clientMap, err := loadAttendeesWithQueries(ctx, qtx, []uuid.UUID{master.ID})
		if err != nil {
			return nil, err
		}
		if req.AttendeeEmployeeIDs != nil {
			finalEmployeeIDs = *req.AttendeeEmployeeIDs
		} else {
			finalEmployeeIDs = employeeMap[master.ID]
		}
		if req.AttendeeClientIDs != nil {
			finalClientIDs = *req.AttendeeClientIDs
		} else {
			finalClientIDs = clientMap[master.ID]
		}
		if err := upsertAttendees(ctx, qtx, override.ID, finalEmployeeIDs, finalClientIDs, false); err != nil {
			s.logger.LogError(ctx, "upsertSingleOccurrenceOverride", "Failed to upsert attendees", err)
			return nil, err
		}
		if err := tx.Commit(ctx); err != nil {
			s.logger.LogError(ctx, "upsertSingleOccurrenceOverride", "Failed to commit transaction", err)
			return nil, err
		}
		finalEmployeeIDs = uniqueUUIDs(finalEmployeeIDs)
		finalClientIDs = uniqueUUIDs(finalClientIDs)
	} else {
		empMap, cliMap, _ := s.loadAttendeesForEvents(ctx, []uuid.UUID{override.ID})
		if len(empMap[override.ID]) == 0 && len(cliMap[override.ID]) == 0 {
			empMap, cliMap, _ = s.loadAttendeesForEvents(ctx, []uuid.UUID{master.ID})
			finalEmployeeIDs = empMap[master.ID]
			finalClientIDs = cliMap[master.ID]
		} else {
			finalEmployeeIDs = empMap[override.ID]
			finalClientIDs = cliMap[override.ID]
		}
	}

	reminders, _ := s.loadRemindersForEvent(ctx, override.ID)

	resp := &domain.EventResponse{
		ID:                  override.ID,
		Kind:                domain.EventKind(override.Kind),
		Status:              string(override.Status),
		WorkApprovalStatus:  string(override.WorkApprovalStatus),
		WorkApprovedBy:      override.WorkApprovedBy,
		WorkRejectedBy:      override.WorkRejectedBy,
		WorkRejectionReason: override.WorkRejectionReason,
		Title:               override.Title,
		Description:         override.Description,
		Location:            override.Location,
		Color:               override.Color,
		OrganizerEmployeeID: override.OrganizerEmployeeID,
		StartAt:             override.StartAt.Time.UTC(),
		EndAt:               override.EndAt.Time.UTC(),
		RRule:               override.Rrule,
		RecurringEventID:    override.RecurringEventID,
		AttendeeEmployeeIDs: finalEmployeeIDs,
		AttendeeClientIDs:   finalClientIDs,
		Reminders:           reminders,
		CreatedAt:           override.CreatedAt.Time.UTC(),
		UpdatedAt:           override.UpdatedAt.Time.UTC(),
	}

	if override.WorkApprovedAt.Valid {
		t := override.WorkApprovedAt.Time.UTC()
		resp.WorkApprovedAt = &t
	}
	if override.WorkRejectedAt.Valid {
		t := override.WorkRejectedAt.Time.UTC()
		resp.WorkRejectedAt = &t
	}
	if override.RecurrenceID.Valid {
		t := override.RecurrenceID.Time.UTC()
		resp.RecurrenceID = &t
	}

	s.logger.LogInfo(ctx, "upsertSingleOccurrenceOverride", "Occurrence override created successfully", zap.String("master_id", master.ID.String()), zap.String("override_id", override.ID.String()))
	return resp, nil
}

// ─── Reminder notifications ────────────────────────────────

func (s *EventService) enqueueReminderNotifications(ctx context.Context, event domain.EventResponse) error {
	now := time.Now().UTC()
	if len(event.Reminders) == 0 {
		return nil
	}

	allAttendeeIDs := append(event.AttendeeEmployeeIDs, event.OrganizerEmployeeID)
	recipients, err := s.resolveUserIDsForEmployees(ctx, uniqueUUIDs(allAttendeeIDs))
	if err != nil || len(recipients) == 0 {
		return err
	}

	for _, reminder := range event.Reminders {
		var sendAt time.Time
		if reminder.RemindAt != nil {
			sendAt = reminder.RemindAt.UTC()
		} else if reminder.MinutesBefore != nil {
			sendAt = event.StartAt.Add(-time.Duration(*reminder.MinutesBefore) * time.Minute)
		} else {
			continue
		}
		if sendAt.Before(now) {
			continue
		}

		message := fmt.Sprintf("Reminder: %s starts at %s", event.Title, event.StartAt.Format(time.RFC3339))
		err = s.taskQ.EnqueueNotificationTask(ctx, domain.NotificationTaskPayload{
			RecipientUserIDs: recipients,
			Type:             "system_reminder",
			Message:          message,
			Data:             domain.NotificationTaskData{},
			CreatedAt:        now,
		}, &domain.TaskEnqueueOptions{ProcessAt: &sendAt})
		if err != nil {
			s.logger.LogWarn(ctx, "enqueueReminderNotifications", "Failed to enqueue task", zap.Error(err))
			continue
		}
	}
	return nil
}

func (s *EventService) resolveUserIDsForEmployees(ctx context.Context, employeeIDs []uuid.UUID) ([]uuid.UUID, error) {
	if len(employeeIDs) == 0 {
		return nil, nil
	}
	rows, err := s.store.ListUserIDsByEmployeeIDs(ctx, employeeIDs)
	if err != nil {
		return nil, err
	}
	return uniqueUUIDs(rows), nil
}

// ─── Pure helpers (package level) ──────────────────────────

func upsertAttendees(ctx context.Context, q *db.Queries, eventID uuid.UUID, employeeIDs, clientIDs []uuid.UUID, skipDelete bool) error {
	if !skipDelete {
		if err := q.DeleteAttendeesByEventID(ctx, eventID); err != nil {
			return fmt.Errorf("failed to clear attendees: %w", err)
		}
	}

	uniqueEmployeeIDs := uniqueUUIDs(employeeIDs)
	if len(uniqueEmployeeIDs) > 0 {
		if err := q.AddEventEmployeeAttendeesBatch(ctx, db.AddEventEmployeeAttendeesBatchParams{
			EventID:     eventID,
			EmployeeIds: uniqueEmployeeIDs,
		}); err != nil {
			return fmt.Errorf("failed to insert employee attendees: %w", err)
		}
	}

	uniqueClientIDs := uniqueUUIDs(clientIDs)
	if len(uniqueClientIDs) > 0 {
		if err := q.AddEventClientAttendeesBatch(ctx, db.AddEventClientAttendeesBatchParams{
			EventID:   eventID,
			ClientIds: uniqueClientIDs,
		}); err != nil {
			return fmt.Errorf("failed to insert client attendees: %w", err)
		}
	}

	return nil
}

func upsertReminders(ctx context.Context, q *db.Queries, eventID uuid.UUID, reminders []domain.ReminderInput, skipDelete bool) ([]domain.ReminderResponse, error) {
	if !skipDelete {
		if err := q.DeleteRemindersByEventID(ctx, eventID); err != nil {
			return nil, fmt.Errorf("failed to clear reminders: %w", err)
		}
	}
	out := make([]domain.ReminderResponse, 0, len(reminders))
	for _, reminder := range reminders {
		remindAt := pgtype.Timestamptz{}
		if reminder.RemindAt != nil {
			remindAt = toPgTimestamptz(*reminder.RemindAt)
		}
		inserted, err := q.AddEventReminder(ctx, db.AddEventReminderParams{EventID: eventID, MinutesBefore: reminder.MinutesBefore, RemindAt: remindAt})
		if err != nil {
			return nil, fmt.Errorf("failed to insert reminder: %w", err)
		}

		var r domain.ReminderResponse
		r.ID = inserted.ID
		r.MinutesBefore = inserted.MinutesBefore
		r.RemindAt = pgTimeToPtr(inserted.RemindAt)
		out = append(out, r)
	}
	return out, nil
}

func buildRule(rr string, dtstart time.Time) (*rrule.RRule, error) {
	opt, err := rrule.StrToROption(rr)
	if err != nil {
		return nil, fmt.Errorf("invalid rrule: %w", err)
	}
	opt.Dtstart = dtstart.UTC()
	r, err := rrule.NewRRule(*opt)
	if err != nil {
		return nil, fmt.Errorf("invalid recurrence rule: %w", err)
	}
	return r, nil
}

func validateRRule(rule string, dtstart time.Time) error {
	_, err := buildRule(rule, dtstart)
	return err
}

func withUntil(rule string, dtstart time.Time, until time.Time) (string, error) {
	opt, err := rrule.StrToROption(rule)
	if err != nil {
		return "", fmt.Errorf("invalid rrule: %w", err)
	}
	opt.Dtstart = dtstart.UTC()
	opt.Until = until.UTC()
	opt.Count = 0
	r, err := rrule.NewRRule(*opt)
	if err != nil {
		return "", fmt.Errorf("invalid recurrence rule: %w", err)
	}
	return r.OrigOptions.RRuleString(), nil
}

func validateReminders(reminders []domain.ReminderInput) error {
	for _, reminder := range reminders {
		if (reminder.MinutesBefore == nil && reminder.RemindAt == nil) || (reminder.MinutesBefore != nil && reminder.RemindAt != nil) {
			return fmt.Errorf("each reminder must set exactly one of minutes_before or remind_at")
		}
		if reminder.MinutesBefore != nil && *reminder.MinutesBefore < 0 {
			return fmt.Errorf("minutes_before must be >= 0")
		}
	}
	return nil
}

func collectEventIDs(events []eventRow) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(events))
	for _, event := range events {
		out = append(out, event.ID)
	}
	return out
}

func chooseAttendees(m map[uuid.UUID][]uuid.UUID, fallback map[uuid.UUID][]uuid.UUID, key uuid.UUID, fallbackKey uuid.UUID) []uuid.UUID {
	if values, ok := m[key]; ok && len(values) > 0 {
		return values
	}
	return fallback[fallbackKey]
}
