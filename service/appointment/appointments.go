package appointment

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/service/notification"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/teambition/rrule-go"
)

type eventRow struct {
	ID                  uuid.UUID
	OrganizerEmployeeID uuid.UUID
	CreatedByEmployeeID uuid.UUID
	Kind                string
	Status              string
	Title               string
	Description         *string
	Location            *string
	Color               *string
	StartAt             time.Time
	EndAt               time.Time
	Timezone            string
	RRule               *string
	RecurringEventID    *uuid.UUID
	RecurrenceID        *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

func (s *appointmentService) CreateEvent(ctx context.Context, req *CreateEventRequest, employeeID uuid.UUID) (*EventResponse, error) {
	if !req.StartAt.Before(req.EndAt) {
		return nil, fmt.Errorf("start_at must be before end_at")
	}
	if req.Kind != EventKindAppointment && req.Kind != EventKindReminder {
		return nil, fmt.Errorf("invalid event kind")
	}
	if req.RRule != nil {
		if err := validateRRule(*req.RRule, req.StartAt.UTC()); err != nil {
			return nil, err
		}
	}
	if err := validateReminders(req.Reminders); err != nil {
		return nil, err
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.New(tx)
	eventModel, err := qtx.CreateCalendarEvent(ctx, db.CreateCalendarEventParams{
		OrganizerEmployeeID: employeeID,
		CreatedByEmployeeID: employeeID,
		Kind:                db.CalendarEventKindEnum(req.Kind),
		Status:              db.CalendarEventStatusEnumConfirmed,
		Title:               strings.TrimSpace(req.Title),
		Description:         req.Description,
		Location:            req.Location,
		Color:               req.Color,
		StartAt:             toPgTimestamptz(req.StartAt),
		EndAt:               toPgTimestamptz(req.EndAt),
		Timezone:            "UTC",
		Rrule:               req.RRule,
		RecurringEventID:    nil,
		RecurrenceID:        pgtype.Timestamptz{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	if err := upsertAttendees(ctx, qtx, eventModel.ID, req.AttendeeEmployeeIDs, req.AttendeeClientIDs); err != nil {
		return nil, err
	}
	if err := upsertReminders(ctx, qtx, eventModel.ID, req.Reminders); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	resp, err := s.GetEvent(ctx, eventModel.ID, employeeID)
	if err != nil {
		return nil, err
	}

	if err := s.enqueueReminderNotifications(ctx, *resp); err != nil {
		// reminder scheduling failure should not fail creation
	}

	return resp, nil
}

func (s *appointmentService) ListEvents(ctx context.Context, req ListEventsRequest, employeeID uuid.UUID) ([]EventOccurrenceResponse, error) {
	if !req.StartAt.Before(req.EndAt) {
		return nil, fmt.Errorf("start_at must be before end_at")
	}

	// Use target employee ID from request if provided, otherwise use the caller's employee ID
	targetEmployeeID := employeeID
	if req.EmployeeID != nil {
		targetEmployeeID = *req.EmployeeID
	}

	masters, err := s.listVisibleMasterEvents(ctx, targetEmployeeID, req.StartAt.UTC(), req.EndAt.UTC())
	if err != nil {
		return nil, err
	}

	attendeeEmployees, attendeeClients, err := s.loadAttendeesForEvents(ctx, collectEventIDs(masters))
	if err != nil {
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
		return nil, err
	}
	exAttendeeEmployees, exAttendeeClients, err := s.loadAttendeesForEvents(ctx, collectEventIDs(exceptions))
	if err != nil {
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

	items := make([]EventOccurrenceResponse, 0)
	for _, e := range masters {
		if e.RRule == nil {
			items = append(items, EventOccurrenceResponse{
				ID:                  e.ID,
				Kind:                EventKind(e.Kind),
				Title:               e.Title,
				Description:         e.Description,
				Location:            e.Location,
				Color:               e.Color,
				StartAt:             e.StartAt.UTC(),
				EndAt:               e.EndAt.UTC(),
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
				items = append(items, EventOccurrenceResponse{
					ID:                  exForOcc.ID,
					MasterEventID:       exForOcc.RecurringEventID,
					Kind:                EventKind(exForOcc.Kind),
					Title:               exForOcc.Title,
					Description:         exForOcc.Description,
					Location:            exForOcc.Location,
					Color:               exForOcc.Color,
					StartAt:             exForOcc.StartAt.UTC(),
					EndAt:               exForOcc.EndAt.UTC(),
					RecurrenceID:        exForOcc.RecurrenceID,
					IsRecurringInstance: true,
					AttendeeEmployeeIDs: chooseAttendees(exAttendeeEmployees, attendeeEmployees, exForOcc.ID, e.ID),
					AttendeeClientIDs:   chooseAttendees(exAttendeeClients, attendeeClients, exForOcc.ID, e.ID),
				})
				continue
			}
			occTime := occ.UTC()
			items = append(items, EventOccurrenceResponse{
				ID:                  e.ID,
				MasterEventID:       &e.ID,
				Kind:                EventKind(e.Kind),
				Title:               e.Title,
				Description:         e.Description,
				Location:            e.Location,
				Color:               e.Color,
				StartAt:             occTime,
				EndAt:               occTime.Add(duration),
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

func (s *appointmentService) GetEvent(ctx context.Context, eventID uuid.UUID, employeeID uuid.UUID) (*EventResponse, error) {
	e, err := s.getVisibleEventByID(ctx, eventID, employeeID)
	if err != nil {
		return nil, err
	}
	attendeeEmployees, attendeeClients, err := s.loadAttendeesForEvents(ctx, []uuid.UUID{eventID})
	if err != nil {
		return nil, err
	}
	reminders, err := s.loadRemindersForEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return &EventResponse{
		ID:                  e.ID,
		Kind:                EventKind(e.Kind),
		Status:              e.Status,
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

func (s *appointmentService) UpdateEvent(ctx context.Context, eventID uuid.UUID, req *UpdateEventRequest, employeeID uuid.UUID) (*EventResponse, error) {
	e, err := s.getVisibleEventByID(ctx, eventID, employeeID)
	if err != nil {
		return nil, err
	}

	if req.Scope == MutationScopeFuture {
		return s.updateEventFuture(ctx, e, req, employeeID)
	}

	if req.Scope == MutationScopeSingle && e.RRule != nil {
		if req.RecurrenceID == nil {
			return nil, fmt.Errorf("recurrence_id is required for single scope on recurring events")
		}
		return s.upsertSingleOccurrenceOverride(ctx, e, req, employeeID)
	}

	if req.StartAt != nil && req.EndAt != nil && !req.StartAt.Before(*req.EndAt) {
		return nil, fmt.Errorf("start_at must be before end_at")
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if req.RRule != nil {
		if err := validateRRule(*req.RRule, e.StartAt.UTC()); err != nil {
			return nil, err
		}
	}

	qtx := db.New(tx)
	err = qtx.UpdateCalendarEvent(ctx, db.UpdateCalendarEventParams{
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
		return nil, err
	}

	if req.AttendeeEmployeeIDs != nil || req.AttendeeClientIDs != nil {
		employeeIDs := []uuid.UUID{}
		clientIDs := []uuid.UUID{}
		if req.AttendeeEmployeeIDs != nil {
			employeeIDs = *req.AttendeeEmployeeIDs
		}
		if req.AttendeeClientIDs != nil {
			clientIDs = *req.AttendeeClientIDs
		}
		if err := upsertAttendees(ctx, qtx, eventID, employeeIDs, clientIDs); err != nil {
			return nil, err
		}
	}

	if req.Reminders != nil {
		if err := validateReminders(*req.Reminders); err != nil {
			return nil, err
		}
		if err := upsertReminders(ctx, qtx, eventID, *req.Reminders); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	resp, err := s.GetEvent(ctx, eventID, employeeID)
	if err != nil {
		return nil, err
	}
	if err := s.enqueueReminderNotifications(ctx, *resp); err != nil {
	}
	return resp, nil
}

func (s *appointmentService) DeleteEvent(ctx context.Context, eventID uuid.UUID, req DeleteEventRequest, employeeID uuid.UUID) error {
	e, err := s.getVisibleEventByID(ctx, eventID, employeeID)
	if err != nil {
		return err
	}

	if req.Scope == MutationScopeFuture {
		if e.RRule == nil {
			return s.Store.CancelCalendarEvent(ctx, eventID)
		}
		if req.RecurrenceID == nil {
			return fmt.Errorf("recurrence_id is required for future scope")
		}
		until := req.RecurrenceID.UTC().Add(-time.Second)
		oldRule, err := withUntil(*e.RRule, e.StartAt.UTC(), until)
		if err != nil {
			return err
		}
		return s.Store.UpdateCalendarEventRRule(ctx, db.UpdateCalendarEventRRuleParams{ID: e.ID, Rrule: &oldRule})
	}

	if req.Scope == MutationScopeSingle && e.RRule != nil {
		if req.RecurrenceID == nil {
			return fmt.Errorf("recurrence_id is required for single scope")
		}
		_, err := s.Store.UpsertCalendarEventOverride(ctx, db.UpsertCalendarEventOverrideParams{
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
		return err
	}

	return s.Store.CancelCalendarEvent(ctx, eventID)
}

func (s *appointmentService) updateEventFuture(ctx context.Context, master eventRow, req *UpdateEventRequest, employeeID uuid.UUID) (*EventResponse, error) {
	if master.RRule == nil {
		return nil, fmt.Errorf("future scope requires recurring event")
	}
	if req.RecurrenceID == nil {
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
			return nil, err
		}
		newRule = *req.RRule
	}

	until := req.RecurrenceID.UTC().Add(-time.Second)
	oldRule, err := withUntil(*master.RRule, master.StartAt.UTC(), until)
	if err != nil {
		return nil, err
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := db.New(tx)
	err = qtx.UpdateCalendarEventRRule(ctx, db.UpdateCalendarEventRRuleParams{ID: master.ID, Rrule: &oldRule})
	if err != nil {
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
		return nil, err
	}

	attendeeEmployees, attendeeClients, err := s.loadAttendeesForEvents(ctx, []uuid.UUID{master.ID})
	if err != nil {
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
	if err := upsertAttendees(ctx, qtx, newMaster.ID, employeeIDs, clientIDs); err != nil {
		return nil, err
	}

	reminders, err := s.loadRemindersForEvent(ctx, master.ID)
	if err != nil {
		return nil, err
	}
	reminderInputs := make([]ReminderInput, 0, len(reminders))
	for _, reminder := range reminders {
		reminderInputs = append(reminderInputs, ReminderInput{MinutesBefore: reminder.MinutesBefore, RemindAt: reminder.RemindAt})
	}
	if req.Reminders != nil {
		if err := validateReminders(*req.Reminders); err != nil {
			return nil, err
		}
		reminderInputs = *req.Reminders
	}
	if err := upsertReminders(ctx, qtx, newMaster.ID, reminderInputs); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	resp, err := s.GetEvent(ctx, newMaster.ID, employeeID)
	if err != nil {
		return nil, err
	}
	if err := s.enqueueReminderNotifications(ctx, *resp); err != nil {
	}
	return resp, nil
}

func (s *appointmentService) upsertSingleOccurrenceOverride(ctx context.Context, master eventRow, req *UpdateEventRequest, employeeID uuid.UUID) (*EventResponse, error) {
	baseStart := req.RecurrenceID.UTC()
	baseEnd := baseStart.Add(master.EndAt.Sub(master.StartAt))
	if req.StartAt != nil {
		baseStart = req.StartAt.UTC()
	}
	if req.EndAt != nil {
		baseEnd = req.EndAt.UTC()
	}
	if !baseStart.Before(baseEnd) {
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

	override, err := s.Store.UpsertCalendarEventOverride(ctx, db.UpsertCalendarEventOverrideParams{
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
		return nil, err
	}

	if req.AttendeeEmployeeIDs != nil || req.AttendeeClientIDs != nil {
		tx, err := s.Store.ConnPool.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer tx.Rollback(ctx)
		employeeIDs := []uuid.UUID{}
		clientIDs := []uuid.UUID{}
		if req.AttendeeEmployeeIDs != nil {
			employeeIDs = *req.AttendeeEmployeeIDs
		}
		if req.AttendeeClientIDs != nil {
			clientIDs = *req.AttendeeClientIDs
		}
		qtx := db.New(tx)
		if err := upsertAttendees(ctx, qtx, override.ID, employeeIDs, clientIDs); err != nil {
			return nil, err
		}
		if err := tx.Commit(ctx); err != nil {
			return nil, err
		}
	}

	return s.GetEvent(ctx, override.ID, employeeID)
}

func (s *appointmentService) enqueueReminderNotifications(ctx context.Context, event EventResponse) error {
	now := time.Now().UTC()
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

		recipients, err := s.resolveUserIDsForEmployees(ctx, uniqueUUIDs(append(event.AttendeeEmployeeIDs, event.OrganizerEmployeeID)))
		if err != nil || len(recipients) == 0 {
			continue
		}
		message := fmt.Sprintf("Reminder: %s starts at %s", event.Title, event.StartAt.Format(time.RFC3339))
		err = s.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
			RecipientUserIDs: recipients,
			Type:             notification.TypeSystemReminder,
			Message:          message,
			Data:             notification.NotificationData{},
			CreatedAt:        now,
		}, asynq.ProcessAt(sendAt))
		if err != nil {
			continue
		}
	}
	return nil
}

func (s *appointmentService) resolveUserIDsForEmployees(ctx context.Context, employeeIDs []uuid.UUID) ([]uuid.UUID, error) {
	if len(employeeIDs) == 0 {
		return nil, nil
	}
	rows, err := s.Store.ListUserIDsByEmployeeIDs(ctx, employeeIDs)
	if err != nil {
		return nil, err
	}
	return uniqueUUIDs(rows), nil
}

func (s *appointmentService) getVisibleEventByID(ctx context.Context, eventID, employeeID uuid.UUID) (eventRow, error) {
	e, err := s.Store.GetVisibleEventByID(ctx, db.GetVisibleEventByIDParams{ID: eventID, EmployeeID: employeeID})
	if err != nil {
		if err == pgx.ErrNoRows {
			return eventRow{}, fmt.Errorf("event not found")
		}
		return eventRow{}, err
	}
	return toEventRow(e), nil
}

func (s *appointmentService) listVisibleMasterEvents(ctx context.Context, employeeID uuid.UUID, startAt, endAt time.Time) ([]eventRow, error) {
	rows, err := s.Store.ListVisibleMasterEvents(ctx, db.ListVisibleMasterEventsParams{
		EmployeeID: employeeID,
		StartAt:    toPgTimestamptz(startAt),
		EndAt:      toPgTimestamptz(endAt),
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

func (s *appointmentService) loadSeriesExceptions(ctx context.Context, seriesIDs []uuid.UUID) ([]eventRow, error) {
	if len(seriesIDs) == 0 {
		return []eventRow{}, nil
	}
	rows, err := s.Store.ListSeriesExceptions(ctx, seriesIDs)
	if err != nil {
		return nil, err
	}
	out := make([]eventRow, 0)
	for _, row := range rows {
		out = append(out, toEventRow(row))
	}
	return out, nil
}

func (s *appointmentService) loadAttendeesForEvents(ctx context.Context, eventIDs []uuid.UUID) (map[uuid.UUID][]uuid.UUID, map[uuid.UUID][]uuid.UUID, error) {
	employeeMap := map[uuid.UUID][]uuid.UUID{}
	clientMap := map[uuid.UUID][]uuid.UUID{}
	if len(eventIDs) == 0 {
		return employeeMap, clientMap, nil
	}

	rows, err := s.Store.ListAttendeesByEventIDs(ctx, eventIDs)
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

func (s *appointmentService) loadRemindersForEvent(ctx context.Context, eventID uuid.UUID) ([]ReminderResponse, error) {
	rows, err := s.Store.ListRemindersByEventID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	out := make([]ReminderResponse, 0)
	for _, row := range rows {
		var r ReminderResponse
		r.ID = row.ID
		r.MinutesBefore = row.MinutesBefore
		if row.RemindAt.Valid {
			t := row.RemindAt.Time.UTC()
			r.RemindAt = &t
		}
		out = append(out, r)
	}
	return out, nil
}

func upsertAttendees(ctx context.Context, q *db.Queries, eventID uuid.UUID, employeeIDs, clientIDs []uuid.UUID) error {
	if err := q.DeleteAttendeesByEventID(ctx, eventID); err != nil {
		return fmt.Errorf("failed to clear attendees: %w", err)
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

func upsertReminders(ctx context.Context, q *db.Queries, eventID uuid.UUID, reminders []ReminderInput) error {
	if err := q.DeleteRemindersByEventID(ctx, eventID); err != nil {
		return fmt.Errorf("failed to clear reminders: %w", err)
	}
	for _, reminder := range reminders {
		remindAt := pgtype.Timestamptz{}
		if reminder.RemindAt != nil {
			remindAt = toPgTimestamptz(*reminder.RemindAt)
		}
		if err := q.AddEventReminder(ctx, db.AddEventReminderParams{EventID: eventID, MinutesBefore: reminder.MinutesBefore, RemindAt: remindAt}); err != nil {
			return fmt.Errorf("failed to insert reminder: %w", err)
		}
	}
	return nil
}

func toEventRow(event db.CalendarEvent) eventRow {
	out := eventRow{
		ID:                  event.ID,
		OrganizerEmployeeID: event.OrganizerEmployeeID,
		CreatedByEmployeeID: event.CreatedByEmployeeID,
		Kind:                string(event.Kind),
		Status:              string(event.Status),
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
	if event.RecurrenceID.Valid {
		t := event.RecurrenceID.Time.UTC()
		out.RecurrenceID = &t
	}
	if event.CreatedAt.Valid {
		out.CreatedAt = event.CreatedAt.Time.UTC()
	}
	if event.UpdatedAt.Valid {
		out.UpdatedAt = event.UpdatedAt.Time.UTC()
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

func validateReminders(reminders []ReminderInput) error {
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

func uniqueUUIDs(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(ids))
	out := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func chooseAttendees(m map[uuid.UUID][]uuid.UUID, fallback map[uuid.UUID][]uuid.UUID, key uuid.UUID, fallbackKey uuid.UUID) []uuid.UUID {
	if values, ok := m[key]; ok && len(values) > 0 {
		return values
	}
	return fallback[fallbackKey]
}
