
-- name: CreateIntakeForm :one
INSERT INTO intake_forms (
    registration_form_id,
    date_of_intake,
    care_type,
    intake_participants,
    family_situation,
    psychological_state,
    self_sufficiency,
    sender_id,
    assigned_location_id,
    risk_assessment,
    intake_conclusion,
    intake_conclusion_notes,
    evaluation_intervals_weeks,
    signature
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;



-- name: ListIntakeForms :many
SELECT
    i.id,
    i.registration_form_id,
    i.date_of_intake,
    i.care_type,
    i.assigned_location_id,
    i.intake_conclusion,
    r.client_first_name,
    r.client_last_name,
    r.client_bsn_number,
    l.street AS assigned_location_street,
    l.house_number AS assigned_location_house_number,
    l.house_number_addition AS assigned_location_house_number_addition,
    l.postal_code AS assigned_location_postal_code,
    l.city AS assigned_location_city,
    EXISTS (
        SELECT 1
        FROM intake_topic_assessments ita
        WHERE ita.intake_form_id = i.id
    ) AS goal_assessment_done,
    COUNT(*) OVER() AS total_count
FROM intake_forms i
JOIN registration_form r ON i.registration_form_id = r.id
LEFT JOIN location l ON i.assigned_location_id = l.id
WHERE
    (
        @search::text IS NULL
        OR @search::text = ''
        OR r.client_first_name ILIKE '%' || @search::text || '%'
        OR r.client_last_name ILIKE '%' || @search::text || '%'
    )
    AND (sqlc.narg('status')::intake_conclusion_enum IS NULL OR i.intake_conclusion = sqlc.narg('status')::intake_conclusion_enum)
ORDER BY
    CASE WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'asc' THEN i.created_at END ASC,
    CASE WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'desc' THEN i.created_at END DESC,
    CASE WHEN @sort_by IS NULL OR @sort_by = '' THEN i.id END DESC
LIMIT $1 OFFSET $2;



-- name: GetIntakeFormDetails :one
SELECT
    i.*,
    r.client_first_name,
    r.client_last_name,
    r.client_bsn_number,
    r.client_goals,
    s.name AS sender_name,
    l.name AS location_name,
    l.street AS location_street,
    l.house_number AS location_house_number,
    l.house_number_addition AS location_house_number_addition,
    l.postal_code AS location_postal_code,
    l.city AS location_city
FROM intake_forms i
JOIN registration_form r ON i.registration_form_id = r.id
LEFT JOIN sender s ON i.sender_id = s.id
LEFT JOIN location l ON i.assigned_location_id = l.id
WHERE i.id = $1
LIMIT 1;



-- name: GetIntakeForm :one
SELECT * FROM intake_forms
WHERE id = $1;


-- name: GetIntakeFormByRegistrationFormID :one
SELECT * FROM intake_forms
WHERE registration_form_id = $1;
