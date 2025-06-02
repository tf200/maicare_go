-- name: CreateAppointment :one
INSERT INTO scheduled_appointments (
    creator_employee_id, 
    start_time,           
    end_time,             
    location,
    color,             
    description       
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: CreateAppointmentTemplate :one
INSERT INTO appointment_templates (
    creator_employee_id, 
    start_time,           
    end_time,             
    location,             
    description,  
    color,        
    recurrence_type,
    recurrence_interval,
    recurrence_end_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;


-- name: BulkAddAppointmentParticipants :exec
INSERT INTO appointment_participants (appointment_id, employee_id)
SELECT
    $1, -- The single appointment_id
    unnest(sqlc.arg(employee_ids)::bigint[]); -- The array of employee_id

-- name: BulkAddAppointmentClients :exec
INSERT INTO appointment_clients (appointment_id, client_id)
SELECT
    $1, -- The single appointment_id
    unnest(sqlc.arg(client_ids)::bigint[]); -- The array of client_ids




-- name: GetAppointmentTemplate :one
SELECT * FROM appointment_templates
WHERE id = $1
LIMIT 1;



-- name: ListEmployeeAppointmentsInRange :many
-- Define the parameters for the query
-- employee_id: The ID of the employee whose appointments are being queried.
-- start_date: The beginning of the time range to search within.
-- end_date: The end of the time range to search within.
SELECT
    sa.id AS appointment_id,
    sa.start_time,
    sa.end_time,
    sa.location,
    sa.description,
    sa.color,
    sa.status,
    sa.is_confirmed,
    sa.creator_employee_id,
    sa.created_at,
    'CREATOR' AS involvement_type -- Indicate the employee created this appointment
FROM
    scheduled_appointments sa
WHERE
    sa.creator_employee_id = sqlc.arg(employee_id)
    -- Check for overlap: Appointment starts before the range ends AND appointment ends after the range starts
    AND sa.start_time < sqlc.arg(end_date)
    AND sa.end_time > sqlc.arg(start_date)

UNION -- Combine with participant appointments, removing duplicates

SELECT
    sa.id AS appointment_id,
    sa.start_time,
    sa.end_time,
    sa.location,
    sa.description,
    sa.color,
    sa.status,
    sa.is_confirmed,
    sa.creator_employee_id, -- Still show who created it
    sa.created_at,
    'PARTICIPANT' AS involvement_type -- Indicate the employee is a participant
FROM
    scheduled_appointments sa
JOIN
    appointment_participants ap ON sa.id = ap.appointment_id
WHERE
    ap.employee_id = sqlc.arg(employee_id)
    -- Check for overlap: Appointment starts before the range ends AND appointment ends after the range starts
    AND sa.start_time < sqlc.arg(end_date)
    AND sa.end_time > sqlc.arg(start_date)

-- Order the combined results by start time
ORDER BY
    start_time;




-- name: ListClientAppointmentsInRange :many
-- Define the parameters for the query
-- client_id: The ID of the client whose appointments are being queried.
-- start_date: The beginning of the time range to search within (inclusive).
-- end_date: The end of the time range to search within (exclusive).
SELECT
    sa.id AS appointment_id,
    sa.start_time,
    sa.end_time,
    sa.location,
    sa.description,
    sa.color,
    sa.status,
    sa.creator_employee_id, -- Include creator info if needed
    sa.created_at
    -- No 'involvement_type' needed as clients are always participants in this context
FROM
    scheduled_appointments sa
JOIN
    appointment_clients ac ON sa.id = ac.appointment_id
WHERE
    ac.client_id = sqlc.arg(client_id)
    -- Check for overlap: Appointment starts before the range ends AND appointment ends after the range starts
    AND sa.start_time < sqlc.arg(end_date)
    AND sa.end_time > sqlc.arg(start_date)

-- Order the results by start time
ORDER BY
    sa.start_time;




-- name: GetScheduledAppointmentByID :one
SELECT
    sa.id,
    sa.appointment_templates_id,
    sa.creator_employee_id,
    creator.first_name AS creator_first_name, -- Alias creator details
    creator.last_name AS creator_last_name,   -- Alias creator details
    sa.start_time,
    sa.end_time,
    sa.location,
    sa.description,
    sa.color,
    sa.status,
    sa.is_confirmed,
    sa.confirmed_by_employee_id,
    confirmer.first_name AS confirmer_first_name, -- Alias confirmer details
    confirmer.last_name AS confirmer_last_name,   -- Alias confirmer details
    sa.confirmed_at,
    sa.created_at,
    sa.updated_at
FROM
    scheduled_appointments sa
LEFT JOIN
    employee_profile creator ON sa.creator_employee_id = creator.id -- Join for creator
LEFT JOIN
    employee_profile confirmer ON sa.confirmed_by_employee_id = confirmer.id -- Join for confirmer
WHERE
    sa.id = $1 -- Filter by appointment ID
LIMIT 1;



-- name: GetAppointmentParticipants :many
SELECT
    ep.id AS employee_id,
    ep.first_name,
    ep.last_name
FROM
    appointment_participants ap
JOIN
    employee_profile ep ON ap.employee_id = ep.id
WHERE
    ap.appointment_id = $1 -- Filter by appointment ID
ORDER BY
    ep.last_name, ep.first_name; -- Optional ordering



-- name: GetAppointmentClients :many
SELECT
    cd.id AS client_id,
    cd.first_name,
    cd.last_name
FROM
    appointment_clients ac
JOIN
    client_details cd ON ac.client_id = cd.id
WHERE
    ac.appointment_id = $1 -- Filter by appointment ID
ORDER BY
    cd.last_name, cd.first_name; -- Optional ordering



-- name: UpdateAppointment :one
UPDATE scheduled_appointments
SET
    start_time = COALESCE($2, start_time),
    end_time = COALESCE($3, end_time),
    location = COALESCE ($4, location),
    description = COALESCE ($5, description),
    color = COALESCE ($6, color),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAppointment :exec
DELETE FROM scheduled_appointments
WHERE id = $1;


-- name: DeleteAppointmentParticipants :exec
DELETE FROM appointment_participants
WHERE appointment_id = $1;

-- name: DeleteAppointmentClients :exec
DELETE FROM appointment_clients
WHERE appointment_id = $1;











-- name: ConfirmAppointment :exec
UPDATE scheduled_appointments
SET
    status = 'CONFIRMED',
    is_confirmed = true,
    confirmed_by_employee_id = sqlc.arg(employee_id),
    confirmed_at = NOW()
WHERE id = $1;