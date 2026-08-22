-- name: CreateIncident :one
WITH new_incident AS (
    SELECT public.begin_incident_creation($23) AS id
), inserted_incident AS (
    INSERT INTO incident (
        id,
        employee_id,
        location_id,
        reporter_involvement,
        informed_parties,
        occurred_at,
        incident_type,
        severity_of_incident,
        incident_explanation,
        recurrence_risk,
        incident_prevent_steps,
        incident_taken_measures,
        cause_categories,
        cause_explanation,
        physical_injury,
        physical_injury_desc,
        psychological_damage,
        psychological_damage_desc,
        needed_consultation,
        follow_up_actions,
        follow_up_notes,
        is_employee_absent,
        additional_details,
        client_id,
        emails
    ) SELECT
        new_incident.id,
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
        $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
        $21, $22, $23, $24
    FROM new_incident
    WHERE new_incident.id IS NOT NULL
    RETURNING *
)
SELECT 
    i.*,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name,
    c.first_name AS client_first_name,
    c.last_name AS client_last_name,
    l.name AS location_name
FROM inserted_incident i
LEFT JOIN employee_profile e ON i.employee_id = e.id
LEFT JOIN client_details c ON i.client_id = c.id
LEFT JOIN location l ON i.location_id = l.id;


-- name: ListIncidents :many
SELECT 
    i.id,
    i.occurred_at,
    i.incident_type,
    i.severity_of_incident,
    i.is_confirmed,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name,
    u.profile_picture AS employee_profile_picture,
    l.name AS location_name,
    COUNT(*) OVER() AS total_count
FROM incident i
JOIN employee_profile e ON i.employee_id = e.id
JOIN custom_user u ON e.user_id = u.id
JOIN location l ON i.location_id = l.id
WHERE i.client_id = $1
  AND public.can_access_client(i.client_id, 'CLIENT.VIEW')
ORDER BY i.occurred_at DESC
LIMIT $2 OFFSET $3;

-- name: CountClientIncidents :one
SELECT COUNT(*)::BIGINT
FROM incident
WHERE client_id = $1
  AND public.can_access_client(client_id, 'CLIENT.VIEW');


-- name: GetIncident :one
SELECT 
    i.*,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name,
    l.name AS location_name,
    c.first_name AS client_first_name,
    c.last_name AS client_last_name
FROM incident i
JOIN employee_profile e ON i.employee_id = e.id
JOIN client_details c ON i.client_id = c.id
JOIN location l ON i.location_id = l.id
WHERE i.id = $1 LIMIT 1;

-- name: UpdateIncident :one
UPDATE incident
SET
    employee_id = COALESCE(sqlc.narg('employee_id'), employee_id),
    location_id = COALESCE(sqlc.narg('location_id'), location_id),
    reporter_involvement = COALESCE(sqlc.narg('reporter_involvement'), reporter_involvement),
    informed_parties = COALESCE(sqlc.narg('informed_parties'), informed_parties),
    occurred_at = COALESCE(sqlc.narg('occurred_at'), occurred_at),
    incident_type = COALESCE(sqlc.narg('incident_type'), incident_type),
    severity_of_incident = COALESCE(sqlc.narg('severity_of_incident'), severity_of_incident),
    incident_explanation = COALESCE(sqlc.narg('incident_explanation'), incident_explanation),
    recurrence_risk = COALESCE(sqlc.narg('recurrence_risk'), recurrence_risk),
    incident_prevent_steps = COALESCE(sqlc.narg('incident_prevent_steps'), incident_prevent_steps),
    incident_taken_measures = COALESCE(sqlc.narg('incident_taken_measures'), incident_taken_measures),
    cause_categories = COALESCE(sqlc.narg('cause_categories'), cause_categories),
    cause_explanation = COALESCE(sqlc.narg('cause_explanation'), cause_explanation),
    physical_injury = COALESCE(sqlc.narg('physical_injury'), physical_injury),
    physical_injury_desc = COALESCE(sqlc.narg('physical_injury_desc'), physical_injury_desc),
    psychological_damage = COALESCE(sqlc.narg('psychological_damage'), psychological_damage),
    psychological_damage_desc = COALESCE(sqlc.narg('psychological_damage_desc'), psychological_damage_desc),
    needed_consultation = COALESCE(sqlc.narg('needed_consultation'), needed_consultation),
    follow_up_actions = COALESCE(sqlc.narg('follow_up_actions'), follow_up_actions),
    follow_up_notes = COALESCE(sqlc.narg('follow_up_notes'), follow_up_notes),
    is_employee_absent = COALESCE(sqlc.narg('is_employee_absent'), is_employee_absent),
    additional_details = COALESCE(sqlc.narg('additional_details'), additional_details),
    emails = COALESCE(sqlc.narg('emails'), emails)
WHERE id = $1
RETURNING *;


-- name: DeleteIncident :one
DELETE FROM incident
WHERE id = $1
RETURNING client_id;


-- name: ConfirmIncident :one
SELECT public.confirm_incident($1)::BIGINT AS affected;


-- name: MarkIncidentConfirmationEmailSent :one
SELECT public.mark_incident_confirmation_email_sent($1, $2)::BIGINT AS affected;

-- name: ClaimIncidentConfirmationEmail :one
SELECT COALESCE(public.claim_incident_confirmation_email($1), '00000000-0000-0000-0000-000000000000'::UUID)::UUID AS claim_token;

-- name: ReleaseIncidentConfirmationEmail :one
SELECT public.release_incident_confirmation_email($1, $2)::BIGINT AS affected;

-- name: SeedIncident :one
INSERT INTO incident (
    employee_id,
    location_id,
    reporter_involvement,
    informed_parties,
    occurred_at,
    incident_type,
    severity_of_incident,
    incident_explanation,
    recurrence_risk,
    incident_prevent_steps,
    incident_taken_measures,
    cause_categories,
    cause_explanation,
    physical_injury,
    physical_injury_desc,
    psychological_damage,
    psychological_damage_desc,
    needed_consultation,
    follow_up_actions,
    follow_up_notes,
    is_employee_absent,
    additional_details,
    client_id,
    emails
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24
)
RETURNING id;


-- name: UpdateIncidentFileUrl :one
UPDATE incident
SET file_url = $2
WHERE id = $1
RETURNING file_url;
