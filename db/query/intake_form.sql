
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
    i.*,
    r.client_first_name,
    r.client_last_name,
    r.client_bsn_number,
    COUNT(*) OVER() AS total_count
FROM intake_forms i
JOIN registration_form r ON i.registration_form_id = r.id
WHERE
    @search::text IS NULL
    OR @search::text = ''
    OR r.client_first_name ILIKE '%' || @search::text || '%'
    OR r.client_last_name ILIKE '%' || @search::text || '%'
ORDER BY
    CASE WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'asc' THEN created_at END ASC,
    CASE WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'desc' THEN created_at END DESC,
    CASE WHEN @sort_by IS NULL OR @sort_by = '' THEN id END DESC
LIMIT $1 OFFSET $2;



-- name: GetIntakeForm :one
SELECT * FROM intake_forms
WHERE id = $1;


-- name: GetIntakeFormByRegistrationFormID :one
SELECT * FROM intake_forms
WHERE registration_form_id = $1;
