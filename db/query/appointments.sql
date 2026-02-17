-- name: ListEmployeeAppointmentsInRange :many
SELECT
    ce.id AS appointment_id,
    ce.start_at AS start_time,
    ce.end_at AS end_time,
    ce.location,
    ce.description,
    ce.color,
    ce.status,
    FALSE AS is_confirmed,
    ce.organizer_employee_id AS creator_employee_id,
    ce.created_at,
    CASE
        WHEN ce.organizer_employee_id = sqlc.arg(employee_id) THEN 'CREATOR'
        ELSE 'PARTICIPANT'
    END AS involvement_type
FROM
    calendar_events ce
LEFT JOIN
    calendar_event_attendees cea ON ce.id = cea.event_id
WHERE
    ce.kind = 'appointment'
    AND ce.status <> 'cancelled'
    AND (ce.organizer_employee_id = sqlc.arg(employee_id) OR cea.employee_id = sqlc.arg(employee_id))
    AND ce.start_at < sqlc.arg(end_date)
    AND ce.end_at > sqlc.arg(start_date)
ORDER BY
    ce.start_at;


-- name: CreateCalendarEvent :one
INSERT INTO calendar_events (
    organizer_employee_id,
    created_by_employee_id,
    kind,
    status,
    title,
    description,
    location,
    color,
    start_at,
    end_at,
    timezone,
    rrule,
    recurring_event_id,
    recurrence_id
) VALUES (
    sqlc.arg(organizer_employee_id),
    sqlc.arg(created_by_employee_id),
    sqlc.arg(kind),
    sqlc.arg(status),
    sqlc.arg(title),
    sqlc.narg(description),
    sqlc.narg(location),
    sqlc.narg(color),
    sqlc.arg(start_at),
    sqlc.arg(end_at),
    sqlc.arg(timezone),
    sqlc.narg(rrule),
    sqlc.narg(recurring_event_id),
    sqlc.narg(recurrence_id)
)
RETURNING *;


-- name: UpsertCalendarEventOverride :one
INSERT INTO calendar_events (
    organizer_employee_id,
    created_by_employee_id,
    kind,
    status,
    title,
    description,
    location,
    color,
    start_at,
    end_at,
    timezone,
    recurring_event_id,
    recurrence_id
) VALUES (
    sqlc.arg(organizer_employee_id),
    sqlc.arg(created_by_employee_id),
    sqlc.arg(kind),
    sqlc.arg(status),
    sqlc.arg(title),
    sqlc.narg(description),
    sqlc.narg(location),
    sqlc.narg(color),
    sqlc.arg(start_at),
    sqlc.arg(end_at),
    sqlc.arg(timezone),
    sqlc.arg(recurring_event_id),
    sqlc.arg(recurrence_id)
)
ON CONFLICT (recurring_event_id, recurrence_id)
DO UPDATE SET
    title = EXCLUDED.title,
    description = EXCLUDED.description,
    location = EXCLUDED.location,
    color = EXCLUDED.color,
    start_at = EXCLUDED.start_at,
    end_at = EXCLUDED.end_at,
    status = EXCLUDED.status,
    updated_at = now()
RETURNING *;


-- name: UpdateCalendarEvent :exec
UPDATE calendar_events
SET
    title = COALESCE(sqlc.narg(title), title),
    description = COALESCE(sqlc.narg(description), description),
    location = COALESCE(sqlc.narg(location), location),
    color = COALESCE(sqlc.narg(color), color),
    start_at = COALESCE(sqlc.narg(start_at), start_at),
    end_at = COALESCE(sqlc.narg(end_at), end_at),
    rrule = CASE
        WHEN sqlc.narg(rrule)::TEXT IS NOT NULL THEN sqlc.narg(rrule)
        ELSE rrule
    END,
    updated_at = now()
WHERE id = sqlc.arg(id);


-- name: UpdateCalendarEventRRule :exec
UPDATE calendar_events
SET rrule = sqlc.arg(rrule), updated_at = now()
WHERE id = sqlc.arg(id);


-- name: CancelCalendarEvent :exec
UPDATE calendar_events
SET status = 'cancelled', updated_at = now()
WHERE id = sqlc.arg(id);


-- name: GetVisibleEventByID :one
SELECT ce.*
FROM calendar_events ce
WHERE ce.id = sqlc.arg(id)
  AND (
      ce.organizer_employee_id = sqlc.arg(employee_id)
      OR EXISTS (
          SELECT 1
          FROM calendar_event_attendees cea
          WHERE cea.event_id = ce.id
            AND cea.employee_id = sqlc.arg(employee_id)
      )
  )
LIMIT 1;


-- name: ListVisibleMasterEvents :many
SELECT ce.*
FROM calendar_events ce
WHERE ce.recurring_event_id IS NULL
  AND ce.status <> 'cancelled'
  AND (
      ce.organizer_employee_id = sqlc.arg(employee_id)
      OR EXISTS (
          SELECT 1
          FROM calendar_event_attendees cea
          WHERE cea.event_id = ce.id
            AND cea.employee_id = sqlc.arg(employee_id)
      )
  )
  AND (
      (ce.rrule IS NULL AND ce.start_at < sqlc.arg(end_at) AND ce.end_at > sqlc.arg(start_at))
      OR
      (ce.rrule IS NOT NULL AND ce.start_at < sqlc.arg(end_at))
  );


-- name: ListSeriesExceptions :many
SELECT *
FROM calendar_events
WHERE recurring_event_id = ANY(sqlc.arg(series_ids)::uuid[]);


-- name: ListAttendeesByEventIDs :many
SELECT event_id, employee_id, client_id
FROM calendar_event_attendees
WHERE event_id = ANY(sqlc.arg(event_ids)::uuid[]);


-- name: DeleteAttendeesByEventID :exec
DELETE FROM calendar_event_attendees
WHERE event_id = sqlc.arg(event_id);


-- name: AddEventEmployeeAttendee :exec
INSERT INTO calendar_event_attendees (event_id, employee_id)
VALUES (sqlc.arg(event_id), sqlc.arg(employee_id));


-- name: AddEventEmployeeAttendeesBatch :exec
INSERT INTO calendar_event_attendees (event_id, employee_id)
SELECT sqlc.arg(event_id), unnest(sqlc.arg(employee_ids)::uuid[])
ON CONFLICT (event_id, employee_id) DO NOTHING;


-- name: AddEventClientAttendee :exec
INSERT INTO calendar_event_attendees (event_id, client_id)
VALUES (sqlc.arg(event_id), sqlc.arg(client_id));


-- name: AddEventClientAttendeesBatch :exec
INSERT INTO calendar_event_attendees (event_id, client_id)
SELECT sqlc.arg(event_id), unnest(sqlc.arg(client_ids)::uuid[])
ON CONFLICT (event_id, client_id) DO NOTHING;


-- name: ListRemindersByEventID :many
SELECT id, minutes_before, remind_at
FROM calendar_event_reminders
WHERE event_id = sqlc.arg(event_id)
ORDER BY created_at ASC;


-- name: DeleteRemindersByEventID :exec
DELETE FROM calendar_event_reminders
WHERE event_id = sqlc.arg(event_id);


-- name: AddEventReminder :exec
INSERT INTO calendar_event_reminders (event_id, channel, minutes_before, remind_at)
VALUES (sqlc.arg(event_id), 'in_app', sqlc.narg(minutes_before), sqlc.narg(remind_at));


-- name: ListClientAppointmentsInRange :many
SELECT
    ce.id AS appointment_id,
    ce.start_at AS start_time,
    ce.end_at AS end_time,
    ce.location,
    ce.description,
    ce.color,
    ce.status,
    ce.organizer_employee_id AS creator_employee_id,
    ce.created_at
FROM
    calendar_events ce
JOIN
    calendar_event_attendees cea ON ce.id = cea.event_id
WHERE
    ce.kind = 'appointment'
    AND ce.status <> 'cancelled'
    AND cea.client_id = sqlc.arg(client_id)
    AND ce.start_at < sqlc.arg(end_date)
    AND ce.end_at > sqlc.arg(start_date)
ORDER BY
    ce.start_at;


-- name: ListClientAppointmentsStartingInRange :many
SELECT
    ce.id AS appointment_id,
    ce.start_at AS start_time,
    ce.end_at AS end_time,
    ce.location,
    ce.description,
    ce.color,
    ce.status,
    ce.organizer_employee_id AS creator_employee_id,
    ce.created_at
FROM
    calendar_events ce
JOIN
    calendar_event_attendees cea ON ce.id = cea.event_id
WHERE
    ce.kind = 'appointment'
    AND ce.status <> 'cancelled'
    AND cea.client_id = sqlc.arg(client_id)
    AND ce.start_at >= sqlc.arg(start_date)
    AND ce.start_at < sqlc.arg(end_date)
ORDER BY
    ce.start_at;
