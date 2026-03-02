package appointment

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/service/notification"
	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/teambition/rrule-go"
	"go.uber.org/zap"
)

type eventRow struct {
	ID                  uuid.UUID
	OrganizerEmployeeID uuid.UUID
	CreatedByEmployeeID uuid.UUID
	Kind                string
	Status              string
	WorkApprovalStatus  string
	WorkApprovedBy      *uuid.UUID
	WorkApprovedAt      *time.Time
	WorkRejectedBy      *uuid.UUID
	WorkRejectedAt      *time.Time
	WorkRejectionReason *string
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
		s.Logger.LogWarn(ctx, "CreateEvent", "Invalid time range", zap.Time("start", req.StartAt), zap.Time("end", req.EndAt))
		return nil, fmt.Errorf("start_at must be before end_at")
	}
	if req.Kind != EventKindAppointment && req.Kind != EventKindReminder {
		s.Logger.LogWarn(ctx, "CreateEvent", "Invalid event kind", zap.String("kind", string(req.Kind)))
		return nil, fmt.Errorf("invalid event kind")
	}
	if req.RRule != nil {
		if err := validateRRule(*req.RRule, req.StartAt.UTC()); err != nil {
			s.Logger.LogWarn(ctx, "CreateEvent", "Invalid RRule", zap.String("rrule", *req.RRule), zap.Error(err))
			return nil, err
		}
	}
	if err := validateReminders(req.Reminders); err != nil {
		s.Logger.LogWarn(ctx, "CreateEvent", "Invalid reminders", zap.Error(err))
		return nil, err
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "CreateEvent", "Failed to begin transaction", err)
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
		s.Logger.LogError(ctx, "CreateEvent", "Failed to create event in DB", err)
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	// isNewEvent = true to skip redundant DELETE
	if err := upsertAttendees(ctx, qtx, eventModel.ID, req.AttendeeEmployeeIDs, req.AttendeeClientIDs, true); err != nil {
		s.Logger.LogError(ctx, "CreateEvent", "Failed to upsert attendees", err, zap.String("event_id", eventModel.ID.String()))
		return nil, err
	}
	reminders, err := upsertReminders(ctx, qtx, eventModel.ID, req.Reminders, true)
	if err != nil {
		s.Logger.LogError(ctx, "CreateEvent", "Failed to upsert reminders", err, zap.String("event_id", eventModel.ID.String()))
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogError(ctx, "CreateEvent", "Failed to commit transaction", err)
		return nil, fmt.Errorf("failed to commit tx: %w", err)
	}

	resp := &EventResponse{
		ID:                  eventModel.ID,
		Kind:                EventKind(eventModel.Kind),
		Status:              string(eventModel.Status),
		WorkApprovalStatus:  string(eventModel.WorkApprovalStatus),
		Title:               eventModel.Title,
		Description:         eventModel.Description,
		Location:            eventModel.Location,
		Color:               eventModel.Color,
		OrganizerEmployeeID: eventModel.OrganizerEmployeeID,
		StartAt:             eventModel.StartAt.Time.UTC(),
		EndAt:               eventModel.EndAt.Time.UTC(),
		RRule:               eventModel.Rrule,
		AttendeeEmployeeIDs: util.UniqueUUIDs(req.AttendeeEmployeeIDs),
		AttendeeClientIDs:   util.UniqueUUIDs(req.AttendeeClientIDs),
		Reminders:           reminders,
		CreatedAt:           eventModel.CreatedAt.Time.UTC(),
		UpdatedAt:           eventModel.UpdatedAt.Time.UTC(),
	}

	s.Logger.LogInfo(ctx, "CreateEvent", "Event created successfully", zap.String("event_id", eventModel.ID.String()))

	if err := s.enqueueReminderNotifications(ctx, *resp); err != nil {
		s.Logger.LogWarn(ctx, "CreateEvent", "Failed to enqueue reminder notifications", zap.Error(err), zap.String("event_id", eventModel.ID.String()))
	}

	return resp, nil
}

func (s *appointmentService) SetEventWorkApproval(
	ctx context.Context,
	eventID uuid.UUID,
	req *SetEventWorkApprovalRequest,
	actorEmployeeID, actorUserID uuid.UUID,
) error {
	if req.Status == "rejected" {
		if req.RejectionReason == nil || strings.TrimSpace(*req.RejectionReason) == "" {
			s.Logger.LogWarn(ctx, "SetEventWorkApproval", "Rejection reason missing", zap.String("event_id", eventID.String()))
			return fmt.Errorf("rejection_reason is required when status is rejected")
		}
	}

	event, err := s.Store.GetCalendarEventByID(ctx, eventID)
	if err != nil {
		s.Logger.LogError(ctx, "SetEventWorkApproval", "Event not found", err, zap.String("event_id", eventID.String()))
		return fmt.Errorf("event not found")
	}
	if event.Kind != db.CalendarEventKindEnumAppointment {
		s.Logger.LogWarn(ctx, "SetEventWorkApproval", "Work approval not supported for this kind", zap.String("kind", string(event.Kind)))
		return fmt.Errorf("work approval is only supported for appointment events")
	}
	if event.Status == db.CalendarEventStatusEnumCancelled {
		s.Logger.LogWarn(ctx, "SetEventWorkApproval", "Cannot approve/reject cancelled event", zap.String("event_id", eventID.String()))
		return fmt.Errorf("cannot approve/reject a cancelled event")
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "SetEventWorkApproval", "Failed to begin transaction", err)
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.New(tx)
	targetEventID := event.ID

	// For recurring series, approval is per occurrence: create/update an override row for (master_id, recurrence_id).
	if req.RecurrenceID != nil {
		// If this is already an override row, just update it.
		if event.RecurringEventID != nil {
			targetEventID = event.ID
		} else {
			if event.Rrule == nil {
				s.Logger.LogWarn(ctx, "SetEventWorkApproval", "Recurrence ID used with non-recurring event", zap.String("event_id", eventID.String()))
				return fmt.Errorf("recurrence_id can only be used with recurring events")
			}
			if !event.StartAt.Valid || !event.EndAt.Valid {
				return fmt.Errorf("invalid event time range")
			}
			duration := event.EndAt.Time.Sub(event.StartAt.Time)
			baseStart := req.RecurrenceID.UTC()
			baseEnd := baseStart.Add(duration)

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
				s.Logger.LogError(ctx, "SetEventWorkApproval", "Failed to upsert occurrence override", err)
				return fmt.Errorf("failed to upsert occurrence override: %w", err)
			}
			// Copy attendees from master so downstream queries (hours, invoicing) can "see" the occurrence row.
			attRows, err := qtx.ListAttendeesByEventIDs(ctx, []uuid.UUID{event.ID})
			if err != nil {
				s.Logger.LogError(ctx, "SetEventWorkApproval", "Failed to load master attendees", err)
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
				s.Logger.LogError(ctx, "SetEventWorkApproval", "Failed to copy attendees", err)
				return fmt.Errorf("failed to copy attendees to occurrence override: %w", err)
			}
			targetEventID = override.ID
		}
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
		s.Logger.LogError(ctx, "SetEventWorkApproval", "Failed to update work approval", err)
		return fmt.Errorf("failed to update work approval: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogError(ctx, "SetEventWorkApproval", "Failed to commit transaction", err)
		return fmt.Errorf("failed to commit tx: %w", err)
	}
	s.Logger.LogInfo(ctx, "SetEventWorkApproval", "Work approval updated successfully", zap.String("event_id", targetEventID.String()), zap.String("status", string(req.Status)))
	return nil
}

func (s *appointmentService) ListWorkApprovalQueue(
	ctx context.Context,
	req *ListWorkApprovalQueueRequest,
) (*ListWorkApprovalQueueResponse, error) {
	const maxWorkApprovalQueueRange = 90 * 24 * time.Hour

	if !req.StartAt.Before(req.EndAt) {
		s.Logger.LogWarn(ctx, "ListWorkApprovalQueue", "Invalid time range", zap.Time("start", req.StartAt), zap.Time("end", req.EndAt))
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
		s.Logger.LogWarn(ctx, "ListWorkApprovalQueue", "Time range too large",
			zap.Time("start", startAt),
			zap.Time("end", endAt),
			zap.Time("effective_end", effectiveEnd),
			zap.Duration("range", effectiveEnd.Sub(startAt)),
		)
		return nil, fmt.Errorf("time range too large; maximum is 90 days")
	}
	if !startAt.Before(effectiveEnd) {
		return &ListWorkApprovalQueueResponse{Items: []WorkApprovalQueueItem{}, Total: 0}, nil
	}

	oneOff, err := s.Store.ListWorkApprovalQueueOneOffAppointmentsStartingInRange(ctx, db.ListWorkApprovalQueueOneOffAppointmentsStartingInRangeParams{
		StartAt:     toPgTimestamptz(startAt),
		EndAt:       toPgTimestamptz(effectiveEnd),
		EmployeeIds: req.EmployeeIDs,
	})
	if err != nil {
		s.Logger.LogError(ctx, "ListWorkApprovalQueue", "Failed to list one-off appointments", err)
		return nil, fmt.Errorf("failed to list one-off appointments: %w", err)
	}

	masters, err := s.Store.ListWorkApprovalQueueRecurringMastersStartingBeforeEnd(ctx, db.ListWorkApprovalQueueRecurringMastersStartingBeforeEndParams{
		EndAt:       toPgTimestamptz(effectiveEnd),
		EmployeeIds: req.EmployeeIDs,
	})
	if err != nil {
		s.Logger.LogError(ctx, "ListWorkApprovalQueue", "Failed to list recurring masters", err)
		return nil, fmt.Errorf("failed to list recurring masters: %w", err)
	}

	masterRows := make([]eventRow, 0, len(masters))
	seriesIDs := make([]uuid.UUID, 0, len(masters))
	for _, m := range masters {
		masterRows = append(masterRows, toEventRow(m))
		seriesIDs = append(seriesIDs, m.ID)
	}

	exceptions := []eventRow{}
	if len(seriesIDs) > 0 {
		exRows, err := s.Store.ListSeriesExceptionsStartingInRange(ctx, db.ListSeriesExceptionsStartingInRangeParams{
			SeriesIds:   seriesIDs,
			StartAt:     toPgTimestamptz(startAt),
			EndAt:       toPgTimestamptz(effectiveEnd),
			EmployeeIds: req.EmployeeIDs,
		})
		if err != nil {
			s.Logger.LogError(ctx, "ListWorkApprovalQueue", "Failed to list series exceptions", err)
			return nil, fmt.Errorf("failed to list series exceptions: %w", err)
		}
		for _, ex := range exRows {
			exceptions = append(exceptions, toEventRow(ex))
		}
	}

	// Attendees for masters + exceptions (choose override attendees if present).
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
		s.Logger.LogError(ctx, "ListWorkApprovalQueue", "Failed to load attendees", err)
		return nil, fmt.Errorf("failed to load attendees: %w", err)
	}

	// Map exceptions by (series_id, recurrence_id_unix).
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

	// Optimization: Filter masters to avoid expanding RRULEs for events that don't match criteria
	filteredMasters := make([]eventRow, 0, len(masterRows))
	for _, m := range masterRows {
		if len(attendeeClients[m.ID]) == 0 {
			continue // Overrides with clients will be caught by the exceptions loop
		}
		if !matchesEmployeeFilter(m.OrganizerEmployeeID, attendeeEmployees[m.ID]) {
			continue // Overrides matching filter will be caught by the exceptions loop
		}
		filteredMasters = append(filteredMasters, m)
	}
	masterRows = filteredMasters

	items := make([]WorkApprovalQueueItem, 0, len(oneOff)+len(masterRows)*4)
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
		items = append(items, WorkApprovalQueueItem{
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

	// One-off items.
	for _, row := range oneOff {
		e := toEventRow(row)
		if onlyEnded && e.EndAt.After(now) {
			continue
		}
		clientIDs := attendeeClients[e.ID]
		if len(clientIDs) == 0 {
			continue
		}
		employeeIDs := attendeeEmployees[e.ID]
		if !matchesEmployeeFilter(e.OrganizerEmployeeID, employeeIDs) {
			continue
		}
		appendItem(e.ID, nil, e.StartAt, e.EndAt, e.OrganizerEmployeeID, employeeIDs, clientIDs, e.WorkApprovalStatus, e.Title, e.Description, e.Location, e.CreatedAt)
	}

	// Recurring occurrences: expand within [startAt, endAt) and use override row if present.
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
				// Override row has its own start/end and approval status.
				effectiveStart = effective.StartAt.UTC()
				effectiveOccEnd = effective.EndAt.UTC()
				effectiveStatus = effective.WorkApprovalStatus
			}

			// Start-within window semantics (matches invoicing selection).
			if effectiveStart.Before(startAt) || !effectiveStart.Before(endAt) {
				continue
			}
			if onlyEnded && effectiveOccEnd.After(now) {
				continue
			}

			// Choose attendees: override attendees if present, otherwise master attendees.
			overrideID := effective.ID
			clientIDs := chooseAttendees(attendeeClients, attendeeClients, overrideID, m.ID)
			if len(clientIDs) == 0 {
				continue
			}
			employeeIDs := chooseAttendees(attendeeEmployees, attendeeEmployees, overrideID, m.ID)
			if !matchesEmployeeFilter(m.OrganizerEmployeeID, employeeIDs) {
				continue
			}

			if _, ok := addedRecurring[m.ID]; !ok {
				addedRecurring[m.ID] = map[int64]bool{}
			}
			addedRecurring[m.ID][key] = true

			appendItem(m.ID, &occurrenceID, effectiveStart, effectiveOccEnd, m.OrganizerEmployeeID, employeeIDs, clientIDs, effectiveStatus, effective.Title, effective.Description, effective.Location, m.CreatedAt)
		}
	}

	// Include moved overrides whose recurrence_id falls outside the time window but effective start_at is inside.
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
		employeeIDs := chooseAttendees(attendeeEmployees, attendeeEmployees, ex.ID, *ex.RecurringEventID)
		if !matchesEmployeeFilter(ex.OrganizerEmployeeID, employeeIDs) {
			continue
		}

		occID := ex.RecurrenceID.UTC()
		appendItem(*ex.RecurringEventID, &occID, ex.StartAt.UTC(), ex.EndAt.UTC(), ex.OrganizerEmployeeID, employeeIDs, clientIDs, ex.WorkApprovalStatus, ex.Title, ex.Description, ex.Location, ex.CreatedAt)
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
		rows, err := s.Store.ListEmployeeNamesByIDs(ctx, employeeIDs)
		if err != nil {
			s.Logger.LogError(ctx, "ListWorkApprovalQueue", "Failed to list employee names", err)
			return nil, fmt.Errorf("failed to list employee names: %w", err)
		}
		for _, r := range rows {
			name := strings.TrimSpace(strings.TrimSpace(r.FirstName) + " " + strings.TrimSpace(r.LastName))
			employeeNames[r.ID] = name
		}
	}
	clientNames := map[uuid.UUID]string{}
	if len(clientIDs) > 0 {
		rows, err := s.Store.ListClientNamesByIDs(ctx, clientIDs)
		if err != nil {
			s.Logger.LogError(ctx, "ListWorkApprovalQueue", "Failed to list client names", err)
			return nil, fmt.Errorf("failed to list client names: %w", err)
		}
		for _, r := range rows {
			name := strings.TrimSpace(strings.TrimSpace(r.FirstName) + " " + strings.TrimSpace(r.LastName))
			clientNames[r.ID] = name
		}
	}

	for i := range paged {
		paged[i].OrganizerEmployee = IDName{
			ID:   paged[i].OrganizerEmployeeID,
			Name: employeeNames[paged[i].OrganizerEmployeeID],
		}

		paged[i].AttendeeEmployees = make([]IDName, 0, len(paged[i].AttendeeEmployeeIDs))
		for _, eid := range paged[i].AttendeeEmployeeIDs {
			paged[i].AttendeeEmployees = append(paged[i].AttendeeEmployees, IDName{ID: eid, Name: employeeNames[eid]})
		}

		paged[i].AttendeeClients = make([]IDName, 0, len(paged[i].AttendeeClientIDs))
		for _, cid := range paged[i].AttendeeClientIDs {
			paged[i].AttendeeClients = append(paged[i].AttendeeClients, IDName{ID: cid, Name: clientNames[cid]})
		}
	}

	return &ListWorkApprovalQueueResponse{
		Items: paged,
		Total: total,
	}, nil
}

func (s *appointmentService) ListEvents(ctx context.Context, req ListEventsRequest, employeeID uuid.UUID) ([]EventOccurrenceResponse, error) {
	if !req.StartAt.Before(req.EndAt) {
		s.Logger.LogWarn(ctx, "ListEvents", "Invalid time range", zap.Time("start", req.StartAt), zap.Time("end", req.EndAt))
		return nil, fmt.Errorf("start_at must be before end_at")
	}

	// Use target employee ID from request if provided, otherwise use the caller's employee ID
	targetEmployeeID := employeeID
	if req.EmployeeID != nil {
		targetEmployeeID = *req.EmployeeID
	}

	masters, err := s.listVisibleMasterEvents(ctx, targetEmployeeID, req.StartAt.UTC(), req.EndAt.UTC())
	if err != nil {
		s.Logger.LogError(ctx, "ListEvents", "Failed to list master events", err)
		return nil, err
	}

	attendeeEmployees, attendeeClients, err := s.loadAttendeesForEvents(ctx, collectEventIDs(masters))
	if err != nil {
		s.Logger.LogError(ctx, "ListEvents", "Failed to load attendees", err)
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
		s.Logger.LogError(ctx, "ListEvents", "Failed to load series exceptions", err)
		return nil, err
	}
	exAttendeeEmployees, exAttendeeClients, err := s.loadAttendeesForEvents(ctx, collectEventIDs(exceptions))
	if err != nil {
		s.Logger.LogError(ctx, "ListEvents", "Failed to load exception attendees", err)
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
					WorkApprovalStatus:  exForOcc.WorkApprovalStatus,
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

func (s *appointmentService) GetEvent(ctx context.Context, eventID uuid.UUID, employeeID uuid.UUID) (*EventResponse, error) {
	e, err := s.getVisibleEventByID(ctx, eventID, employeeID)
	if err != nil {
		s.Logger.LogError(ctx, "GetEvent", "Failed to fetch visible event", err, zap.String("event_id", eventID.String()))
		return nil, err
	}
	attendeeEmployees, attendeeClients, err := s.loadAttendeesForEvents(ctx, []uuid.UUID{eventID})
	if err != nil {
		s.Logger.LogError(ctx, "GetEvent", "Failed to load attendees", err, zap.String("event_id", eventID.String()))
		return nil, err
	}
	reminders, err := s.loadRemindersForEvent(ctx, eventID)
	if err != nil {
		s.Logger.LogError(ctx, "GetEvent", "Failed to load reminders", err, zap.String("event_id", eventID.String()))
		return nil, err
	}

	return &EventResponse{
		ID:                  e.ID,
		Kind:                EventKind(e.Kind),
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

func (s *appointmentService) UpdateEvent(ctx context.Context, eventID uuid.UUID, req *UpdateEventRequest, employeeID uuid.UUID) (*EventResponse, error) {
	e, err := s.getVisibleEventByID(ctx, eventID, employeeID)
	if err != nil {
		s.Logger.LogError(ctx, "UpdateEvent", "Failed to fetch visible event", err, zap.String("event_id", eventID.String()))
		return nil, err
	}

	if req.Scope == MutationScopeFuture {
		return s.updateEventFuture(ctx, e, req, employeeID)
	}

	if req.Scope == MutationScopeSingle && e.RRule != nil {
		if req.RecurrenceID == nil {
			s.Logger.LogWarn(ctx, "UpdateEvent", "Recurrence ID missing for single scope", zap.String("event_id", eventID.String()))
			return nil, fmt.Errorf("recurrence_id is required for single scope on recurring events")
		}
		return s.upsertSingleOccurrenceOverride(ctx, e, req, employeeID)
	}

	if req.StartAt != nil && req.EndAt != nil && !req.StartAt.Before(*req.EndAt) {
		s.Logger.LogWarn(ctx, "UpdateEvent", "Invalid time range", zap.Time("start", *req.StartAt), zap.Time("end", *req.EndAt))
		return nil, fmt.Errorf("start_at must be before end_at")
	}

	var attendeeEmployeeIDs []uuid.UUID
	var attendeeClientIDs []uuid.UUID
	var reminders []ReminderResponse

	if req.AttendeeEmployeeIDs == nil || req.AttendeeClientIDs == nil {
		empMap, clientMap, err := s.loadAttendeesForEvents(ctx, []uuid.UUID{eventID})
		if err != nil {
			s.Logger.LogError(ctx, "UpdateEvent", "Failed to load existing attendees", err, zap.String("event_id", eventID.String()))
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
			s.Logger.LogError(ctx, "UpdateEvent", "Failed to load existing reminders", err, zap.String("event_id", eventID.String()))
			return nil, err
		}
		reminders = rems
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "UpdateEvent", "Failed to begin transaction", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	if req.RRule != nil {
		if err := validateRRule(*req.RRule, e.StartAt); err != nil {
			s.Logger.LogWarn(ctx, "UpdateEvent", "Invalid RRule", zap.String("rrule", *req.RRule), zap.Error(err))
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
		s.Logger.LogError(ctx, "UpdateEvent", "Failed to update event in DB", err)
		return nil, err
	}

	if req.AttendeeEmployeeIDs != nil || req.AttendeeClientIDs != nil {
		if req.AttendeeEmployeeIDs != nil {
			attendeeEmployeeIDs = *req.AttendeeEmployeeIDs
		}
		if req.AttendeeClientIDs != nil {
			attendeeClientIDs = *req.AttendeeClientIDs
		}
		// isNewEvent = false because we are updating
		if err := upsertAttendees(ctx, qtx, eventID, attendeeEmployeeIDs, attendeeClientIDs, false); err != nil {
			s.Logger.LogError(ctx, "UpdateEvent", "Failed to upsert attendees", err)
			return nil, err
		}
		attendeeEmployeeIDs = util.UniqueUUIDs(attendeeEmployeeIDs)
		attendeeClientIDs = util.UniqueUUIDs(attendeeClientIDs)
	}

	if req.Reminders != nil {
		if err := validateReminders(*req.Reminders); err != nil {
			s.Logger.LogWarn(ctx, "UpdateEvent", "Invalid reminders", zap.Error(err))
			return nil, err
		}
		// isNewEvent = false because we are updating
		rems, err := upsertReminders(ctx, qtx, eventID, *req.Reminders, false)
		if err != nil {
			s.Logger.LogError(ctx, "UpdateEvent", "Failed to upsert reminders", err)
			return nil, err
		}
		reminders = rems
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogError(ctx, "UpdateEvent", "Failed to commit transaction", err)
		return nil, err
	}

	resp := &EventResponse{
		ID:                  updatedEvent.ID,
		Kind:                EventKind(updatedEvent.Kind),
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

	s.Logger.LogInfo(ctx, "UpdateEvent", "Event updated successfully", zap.String("event_id", eventID.String()))

	if err := s.enqueueReminderNotifications(ctx, *resp); err != nil {
		s.Logger.LogWarn(ctx, "UpdateEvent", "Failed to enqueue reminder notifications", zap.Error(err))
	}
	return resp, nil
}

func (s *appointmentService) DeleteEvent(ctx context.Context, eventID uuid.UUID, req DeleteEventRequest, employeeID uuid.UUID) error {
	e, err := s.getVisibleEventByID(ctx, eventID, employeeID)
	if err != nil {
		s.Logger.LogError(ctx, "DeleteEvent", "Failed to fetch visible event", err, zap.String("event_id", eventID.String()))
		return err
	}

	if req.Scope == MutationScopeFuture {
		if e.RRule == nil {
			return s.Store.CancelCalendarEvent(ctx, eventID)
		}
		if req.RecurrenceID == nil {
			s.Logger.LogWarn(ctx, "DeleteEvent", "Recurrence ID missing for future scope", zap.String("event_id", eventID.String()))
			return fmt.Errorf("recurrence_id is required for future scope")
		}
		until := req.RecurrenceID.UTC().Add(-time.Second)
		oldRule, err := withUntil(*e.RRule, e.StartAt.UTC(), until)
		if err != nil {
			s.Logger.LogError(ctx, "DeleteEvent", "Failed to update RRule for future scope", err, zap.String("event_id", eventID.String()))
			return err
		}
		err = s.Store.UpdateCalendarEventRRule(ctx, db.UpdateCalendarEventRRuleParams{ID: e.ID, Rrule: &oldRule})
		if err != nil {
			s.Logger.LogError(ctx, "DeleteEvent", "Failed to update RRule in DB", err, zap.String("event_id", eventID.String()))
			return err
		}
		s.Logger.LogInfo(ctx, "DeleteEvent", "Event series truncated for future scope", zap.String("event_id", eventID.String()))
		return nil
	}

	if req.Scope == MutationScopeSingle && e.RRule != nil {
		if req.RecurrenceID == nil {
			s.Logger.LogWarn(ctx, "DeleteEvent", "Recurrence ID missing for single scope", zap.String("event_id", eventID.String()))
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
		if err != nil {
			s.Logger.LogError(ctx, "DeleteEvent", "Failed to upsert cancelled override", err, zap.String("event_id", eventID.String()))
			return err
		}
		s.Logger.LogInfo(ctx, "DeleteEvent", "Single occurrence cancelled", zap.String("event_id", eventID.String()), zap.Time("recurrence_id", *req.RecurrenceID))
		return nil
	}

	err = s.Store.CancelCalendarEvent(ctx, eventID)
	if err != nil {
		s.Logger.LogError(ctx, "DeleteEvent", "Failed to cancel calendar event", err, zap.String("event_id", eventID.String()))
		return err
	}
	s.Logger.LogInfo(ctx, "DeleteEvent", "Event cancelled successfully", zap.String("event_id", eventID.String()))
	return nil
}

func (s *appointmentService) updateEventFuture(ctx context.Context, master eventRow, req *UpdateEventRequest, employeeID uuid.UUID) (*EventResponse, error) {
	if master.RRule == nil {
		s.Logger.LogWarn(ctx, "updateEventFuture", "Future scope requires recurring event", zap.String("event_id", master.ID.String()))
		return nil, fmt.Errorf("future scope requires recurring event")
	}
	if req.RecurrenceID == nil {
		s.Logger.LogWarn(ctx, "updateEventFuture", "Recurrence ID missing", zap.String("event_id", master.ID.String()))
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
		s.Logger.LogWarn(ctx, "updateEventFuture", "Invalid time range", zap.Time("start", futureStart), zap.Time("end", futureEnd))
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
			s.Logger.LogWarn(ctx, "updateEventFuture", "Invalid new RRule", zap.Error(err))
			return nil, err
		}
		newRule = *req.RRule
	}

	until := req.RecurrenceID.UTC().Add(-time.Second)
	oldRule, err := withUntil(*master.RRule, master.StartAt.UTC(), until)
	if err != nil {
		s.Logger.LogError(ctx, "updateEventFuture", "Failed to calculate old rule end", err)
		return nil, err
	}

	tx, err := s.Store.ConnPool.Begin(ctx)
	if err != nil {
		s.Logger.LogError(ctx, "updateEventFuture", "Failed to begin transaction", err)
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := db.New(tx)
	err = qtx.UpdateCalendarEventRRule(ctx, db.UpdateCalendarEventRRuleParams{ID: master.ID, Rrule: &oldRule})
	if err != nil {
		s.Logger.LogError(ctx, "updateEventFuture", "Failed to update old rule in DB", err)
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
		s.Logger.LogError(ctx, "updateEventFuture", "Failed to create new master event", err)
		return nil, err
	}

	attendeeEmployees, attendeeClients, err := s.loadAttendeesForEvents(ctx, []uuid.UUID{master.ID})
	if err != nil {
		s.Logger.LogError(ctx, "updateEventFuture", "Failed to load master attendees", err)
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
	// isNewEvent = true to skip redundant DELETE
	if err := upsertAttendees(ctx, qtx, newMaster.ID, employeeIDs, clientIDs, true); err != nil {
		s.Logger.LogError(ctx, "updateEventFuture", "Failed to upsert attendees", err)
		return nil, err
	}

	reminders, err := s.loadRemindersForEvent(ctx, master.ID)
	if err != nil {
		s.Logger.LogError(ctx, "updateEventFuture", "Failed to load master reminders", err)
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
	// isNewEvent = true to skip redundant DELETE
	insertedReminders, err := upsertReminders(ctx, qtx, newMaster.ID, reminderInputs, true)
	if err != nil {
		s.Logger.LogError(ctx, "updateEventFuture", "Failed to upsert reminders", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		s.Logger.LogError(ctx, "updateEventFuture", "Failed to commit transaction", err)
		return nil, err
	}

	resp := &EventResponse{
		ID:                  newMaster.ID,
		Kind:                EventKind(newMaster.Kind),
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
		AttendeeEmployeeIDs: util.UniqueUUIDs(employeeIDs),
		AttendeeClientIDs:   util.UniqueUUIDs(clientIDs),
		Reminders:           insertedReminders,
		CreatedAt:           newMaster.CreatedAt.Time.UTC(),
		UpdatedAt:           newMaster.UpdatedAt.Time.UTC(),
	}

	s.Logger.LogInfo(ctx, "updateEventFuture", "Future events split into new master successfully", zap.String("old_master_id", master.ID.String()), zap.String("new_master_id", newMaster.ID.String()))

	if err := s.enqueueReminderNotifications(ctx, *resp); err != nil {
		s.Logger.LogWarn(ctx, "updateEventFuture", "Failed to enqueue reminder notifications", zap.Error(err))
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
		s.Logger.LogWarn(ctx, "upsertSingleOccurrenceOverride", "Invalid time range", zap.Time("start", baseStart), zap.Time("end", baseEnd))
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
		s.Logger.LogError(ctx, "upsertSingleOccurrenceOverride", "Failed to upsert override in DB", err)
		return nil, err
	}

	var finalEmployeeIDs []uuid.UUID
	var finalClientIDs []uuid.UUID

	if req.AttendeeEmployeeIDs != nil || req.AttendeeClientIDs != nil {
		tx, err := s.Store.ConnPool.Begin(ctx)
		if err != nil {
			s.Logger.LogError(ctx, "upsertSingleOccurrenceOverride", "Failed to begin transaction", err)
			return nil, err
		}
		defer tx.Rollback(ctx)
		if req.AttendeeEmployeeIDs != nil {
			finalEmployeeIDs = *req.AttendeeEmployeeIDs
		} else {
			empMap, _, _ := s.loadAttendeesForEvents(ctx, []uuid.UUID{master.ID})
			finalEmployeeIDs = empMap[master.ID]
		}
		if req.AttendeeClientIDs != nil {
			finalClientIDs = *req.AttendeeClientIDs
		} else {
			_, cliMap, _ := s.loadAttendeesForEvents(ctx, []uuid.UUID{master.ID})
			finalClientIDs = cliMap[master.ID]
		}
		qtx := db.New(tx)
		// isNewEvent = true to skip redundant DELETE because UpsertCalendarEventOverride creates a new row if it didn't exist
		// NOTE: If it DID exist, we might need skipDelete=false.
		// But in this logic, we usually create overrides on the fly.
		// For safety, let's check if it was newly created. Actually, Upsert usually implies we might be overwriting.
		// Let's use false here to be safe unless we are sure.
		if err := upsertAttendees(ctx, qtx, override.ID, finalEmployeeIDs, finalClientIDs, false); err != nil {
			s.Logger.LogError(ctx, "upsertSingleOccurrenceOverride", "Failed to upsert attendees", err)
			return nil, err
		}
		if err := tx.Commit(ctx); err != nil {
			s.Logger.LogError(ctx, "upsertSingleOccurrenceOverride", "Failed to commit transaction", err)
			return nil, err
		}
		finalEmployeeIDs = util.UniqueUUIDs(finalEmployeeIDs)
		finalClientIDs = util.UniqueUUIDs(finalClientIDs)
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

	resp := &EventResponse{
		ID:                  override.ID,
		Kind:                EventKind(override.Kind),
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

	s.Logger.LogInfo(ctx, "upsertSingleOccurrenceOverride", "Occurrence override created successfully", zap.String("master_id", master.ID.String()), zap.String("override_id", override.ID.String()))
	return resp, nil
}

func (s *appointmentService) enqueueReminderNotifications(ctx context.Context, event EventResponse) error {
	now := time.Now().UTC()
	if len(event.Reminders) == 0 {
		return nil
	}

	// Resolve recipients once outside the loop (Fix N+1)
	allAttendeeIDs := append(event.AttendeeEmployeeIDs, event.OrganizerEmployeeID)
	recipients, err := s.resolveUserIDsForEmployees(ctx, util.UniqueUUIDs(allAttendeeIDs))
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
		err = s.asynqClient.EnqueueNotificationTask(ctx, notification.NotificationPayload{
			RecipientUserIDs: recipients,
			Type:             notification.TypeSystemReminder,
			Message:          message,
			Data:             notification.NotificationData{},
			CreatedAt:        now,
		}, asynq.ProcessAt(sendAt))
		if err != nil {
			s.Logger.LogWarn(ctx, "enqueueReminderNotifications", "Failed to enqueue task", zap.Error(err))
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
	return util.UniqueUUIDs(rows), nil
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

func upsertAttendees(ctx context.Context, q *db.Queries, eventID uuid.UUID, employeeIDs, clientIDs []uuid.UUID, skipDelete bool) error {
	if !skipDelete {
		if err := q.DeleteAttendeesByEventID(ctx, eventID); err != nil {
			return fmt.Errorf("failed to clear attendees: %w", err)
		}
	}

	uniqueEmployeeIDs := util.UniqueUUIDs(employeeIDs)
	if len(uniqueEmployeeIDs) > 0 {
		if err := q.AddEventEmployeeAttendeesBatch(ctx, db.AddEventEmployeeAttendeesBatchParams{
			EventID:     eventID,
			EmployeeIds: uniqueEmployeeIDs,
		}); err != nil {
			return fmt.Errorf("failed to insert employee attendees: %w", err)
		}
	}

	uniqueClientIDs := util.UniqueUUIDs(clientIDs)
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

func upsertReminders(ctx context.Context, q *db.Queries, eventID uuid.UUID, reminders []ReminderInput, skipDelete bool) ([]ReminderResponse, error) {
	if !skipDelete {
		if err := q.DeleteRemindersByEventID(ctx, eventID); err != nil {
			return nil, fmt.Errorf("failed to clear reminders: %w", err)
		}
	}
	out := make([]ReminderResponse, 0, len(reminders))
	for _, reminder := range reminders {
		remindAt := pgtype.Timestamptz{}
		if reminder.RemindAt != nil {
			remindAt = toPgTimestamptz(*reminder.RemindAt)
		}
		inserted, err := q.AddEventReminder(ctx, db.AddEventReminderParams{EventID: eventID, MinutesBefore: reminder.MinutesBefore, RemindAt: remindAt})
		if err != nil {
			return nil, fmt.Errorf("failed to insert reminder: %w", err)
		}

		var r ReminderResponse
		r.ID = inserted.ID
		r.MinutesBefore = inserted.MinutesBefore
		if inserted.RemindAt.Valid {
			t := inserted.RemindAt.Time.UTC()
			r.RemindAt = &t
		}
		out = append(out, r)
	}
	return out, nil
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
	if event.WorkApprovedAt.Valid {
		t := event.WorkApprovedAt.Time.UTC()
		out.WorkApprovedAt = &t
	}
	if event.WorkRejectedAt.Valid {
		t := event.WorkRejectedAt.Time.UTC()
		out.WorkRejectedAt = &t
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

func chooseAttendees(m map[uuid.UUID][]uuid.UUID, fallback map[uuid.UUID][]uuid.UUID, key uuid.UUID, fallbackKey uuid.UUID) []uuid.UUID {
	if values, ok := m[key]; ok && len(values) > 0 {
		return values
	}
	return fallback[fallbackKey]
}
