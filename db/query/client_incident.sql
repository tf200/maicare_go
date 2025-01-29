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
