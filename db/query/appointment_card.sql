-- name: CreateAppointmentCard :one
INSERT INTO appointment_card (
    client_id,
    general_information,
    important_contacts,
    household_info,
    organization_agreements,
    youth_officer_agreements,
    treatment_agreements,
    smoking_rules,
    work,
    school_internship,
    travel,
    leave
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;

-- name: GetAppointmentCard :zero_or_one
SELECT 
    ac.*,
    c.first_name,
    c.last_name
FROM appointment_card ac
JOIN client_details c ON ac.client_id = c.id
WHERE ac.client_id = $1
LIMIT 1;



-- name: UpdateAppointmentCard :one
UPDATE appointment_card
SET
    general_information = COALESCE(sqlc.narg('general_information'), general_information),
    important_contacts = COALESCE(sqlc.narg('important_contacts'), important_contacts),
    household_info = COALESCE(sqlc.narg('household_info'), household_info),
    organization_agreements = COALESCE(sqlc.narg('organization_agreements'), organization_agreements),
    youth_officer_agreements = COALESCE(sqlc.narg('youth_officer_agreements'), youth_officer_agreements),
    treatment_agreements = COALESCE(sqlc.narg('treatment_agreements'), treatment_agreements),
    smoking_rules = COALESCE(sqlc.narg('smoking_rules'), smoking_rules),
    work = COALESCE(sqlc.narg('work'), work),
    school_internship = COALESCE(sqlc.narg('school_internship'), school_internship),
    travel = COALESCE(sqlc.narg('travel'), travel),
    leave = COALESCE(sqlc.narg('leave'), leave)
WHERE client_id = $1
RETURNING *;

-- name: UpdateAppointmentCardUrl :one
UPDATE appointment_card
SET
    file_url = COALESCE(sqlc.narg('file_url'), file_url)
WHERE client_id = $1
RETURNING file_url;

