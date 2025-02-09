-- name: CreateIncident :one
INSERT INTO incident (
    employee_id,
    location_id,
    reporter_involvement,
    inform_who,
    incident_date,
    runtime_incident,
    incident_type,
    passing_away,
    self_harm,
    violence,
    fire_water_damage,
    accident,
    client_absence,
    medicines,
    organization,
    use_prohibited_substances,
    other_notifications,
    severity_of_incident,
    incident_explanation,
    recurrence_risk,
    incident_prevent_steps,
    incident_taken_measures,
    technical,
    organizational,
    mese_worker,
    client_options,
    other_cause,
    cause_explanation,
    physical_injury,
    physical_injury_desc,
    psychological_damage,
    psychological_damage_desc,
    needed_consultation,
    succession,
    succession_desc,
    other,
    other_desc,
    additional_appointments,
    employee_absenteeism,
    client_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
    $31, $32, $33, $34, $35, $36, $37, $38, $39, $40
) RETURNING *;


-- name: ListIncidents :many
SELECT 
    i.*,
    COUNT(*) OVER() AS total_count,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name
FROM incident i
JOIN employee_profile e ON i.employee_id = e.id
WHERE i.client_id = $1
ORDER BY i.incident_date DESC
LIMIT $2 OFFSET $3;


-- name: GetIncident :one
SELECT 
    i.*,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name,
    l.name AS location_name
FROM incident i
JOIN employee_profile e ON i.employee_id = e.id
JOIN location l ON i.location_id = l.id
WHERE i.id = $1 LIMIT 1;


-- name: UpdateIncident :one
UPDATE incident
SET
    employee_id = COALESCE(sqlc.narg('employee_id'), employee_id),
    location_id = COALESCE(sqlc.narg('location_id'), location_id),
    reporter_involvement = COALESCE(sqlc.narg('reporter_involvement'), reporter_involvement),
    inform_who = COALESCE(sqlc.narg('inform_who'), inform_who),
    incident_date = COALESCE(sqlc.narg('incident_date'), incident_date),
    runtime_incident = COALESCE(sqlc.narg('runtime_incident'), runtime_incident),
    incident_type = COALESCE(sqlc.narg('incident_type'), incident_type),
    passing_away = COALESCE(sqlc.narg('passing_away'), passing_away),
    self_harm = COALESCE(sqlc.narg('self_harm'), self_harm),
    violence = COALESCE(sqlc.narg('violence'), violence),
    fire_water_damage = COALESCE(sqlc.narg('fire_water_damage'), fire_water_damage),
    accident = COALESCE(sqlc.narg('accident'), accident),
    client_absence = COALESCE(sqlc.narg('client_absence'), client_absence),
    medicines = COALESCE(sqlc.narg('medicines'), medicines),
    organization = COALESCE(sqlc.narg('organization'), organization),
    use_prohibited_substances = COALESCE(sqlc.narg('use_prohibited_substances'), use_prohibited_substances),
    other_notifications = COALESCE(sqlc.narg('other_notifications'), other_notifications),
    severity_of_incident = COALESCE(sqlc.narg('severity_of_incident'), severity_of_incident),
    incident_explanation = COALESCE(sqlc.narg('incident_explanation'), incident_explanation),
    recurrence_risk = COALESCE(sqlc.narg('recurrence_risk'), recurrence_risk),
    incident_prevent_steps = COALESCE(sqlc.narg('incident_prevent_steps'), incident_prevent_steps),
    incident_taken_measures = COALESCE(sqlc.narg('incident_taken_measures'), incident_taken_measures),
    technical = COALESCE(sqlc.narg('technical'), technical),
    organizational = COALESCE(sqlc.narg('organizational'), organizational),
    mese_worker = COALESCE(sqlc.narg('mese_worker'), mese_worker),
    client_options = COALESCE(sqlc.narg('client_options'), client_options),
    other_cause = COALESCE(sqlc.narg('other_cause'), other_cause),
    cause_explanation = COALESCE(sqlc.narg('cause_explanation'), cause_explanation),
    physical_injury = COALESCE(sqlc.narg('physical_injury'), physical_injury),
    physical_injury_desc = COALESCE(sqlc.narg('physical_injury_desc'), physical_injury_desc),
    psychological_damage = COALESCE(sqlc.narg('psychological_damage'), psychological_damage),
    psychological_damage_desc = COALESCE(sqlc.narg('psychological_damage_desc'), psychological_damage_desc),
    needed_consultation = COALESCE(sqlc.narg('needed_consultation'), needed_consultation),
    succession = COALESCE(sqlc.narg('succession'), succession),
    succession_desc = COALESCE(sqlc.narg('succession_desc'), succession_desc),
    other = COALESCE(sqlc.narg('other'), other),
    other_desc = COALESCE(sqlc.narg('other_desc'), other_desc),
    additional_appointments = COALESCE(sqlc.narg('additional_appointments'), additional_appointments),
    employee_absenteeism = COALESCE(sqlc.narg('employee_absenteeism'), employee_absenteeism)
WHERE id = $1
RETURNING *;


-- name: DeleteIncident :one
DELETE FROM incident
WHERE id = $1
RETURNING *;


-- name: ConfirmIncident :one
UPDATE incident
SET is_confirmed = true
WHERE id = $1
RETURNING id, is_confirmed, file_url;


-- name: UpdateIncidentFileUrl :one
UPDATE incident
SET file_url = $2
WHERE id = $1
RETURNING file_url;