
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


-- name: CreateSeedIntakeForm :one
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
    $1,
    $2,
    $3,
    ARRAY['client','parents/guardians','care_coordinator']::intake_participants_enum[],
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13
) RETURNING
    id,
    registration_form_id,
    care_type,
    sender_id,
    assigned_location_id,
    evaluation_intervals_weeks;



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


-- name: GetIntakeFormTotals :one
SELECT
    COUNT(*) FILTER (
        WHERE i.intake_conclusion = 'further_investigation'::intake_conclusion_enum
    )::bigint AS further_investigation_total,
    COUNT(*) FILTER (
        WHERE NOT EXISTS (
            SELECT 1
            FROM intake_topic_assessments ita
            WHERE ita.intake_form_id = i.id
        )
    )::bigint AS without_goals_total
FROM intake_forms i;



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
    l.city AS location_city,
    EXISTS (
        SELECT 1
        FROM client_details cd
        WHERE cd.intake_form_id = i.id
    ) AS has_client
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


-- name: UpdateIntakeConclusion :one
UPDATE intake_forms
SET
    intake_conclusion = $2,
    intake_conclusion_notes = COALESCE(sqlc.narg('intake_conclusion_notes'), intake_conclusion_notes),
    updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: UpdateIntakeForm :one
UPDATE intake_forms
SET
    date_of_intake = COALESCE(sqlc.narg('date_of_intake'), date_of_intake),
    care_type = COALESCE(sqlc.narg('care_type'), care_type),
    intake_participants = COALESCE(sqlc.narg('intake_participants'), intake_participants),
    family_situation = CASE
        WHEN @clear_family_situation::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('family_situation'), family_situation)
    END,
    psychological_state = CASE
        WHEN @clear_psychological_state::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('psychological_state'), psychological_state)
    END,
    self_sufficiency = COALESCE(sqlc.narg('self_sufficiency'), self_sufficiency),
    sender_id = CASE
        WHEN @clear_sender_id::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('sender_id'), sender_id)
    END,
    assigned_location_id = CASE
        WHEN @clear_assigned_location_id::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('assigned_location_id'), assigned_location_id)
    END,
    risk_assessment = CASE
        WHEN @clear_risk_assessment::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('risk_assessment'), risk_assessment)
    END,
    evaluation_intervals_weeks = COALESCE(sqlc.narg('evaluation_intervals_weeks'), evaluation_intervals_weeks),
    signature = CASE
        WHEN @clear_signature::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('signature'), signature)
    END,
    updated_at = NOW()
WHERE id = @id
RETURNING *;
