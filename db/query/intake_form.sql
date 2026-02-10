
-- name: CreateIntakeForm :one
INSERT INTO intake_forms (
    registration_form_id,
    date_of_intake,
    care_type,
    intake_participants,
    family_situation,
    psychological_state,
    self_sufficiency,
    maturity_matrix_id,
    goals,
    risk_assessment,
    intake_conclusion,
    intake_conclusion_notes,
    signature,
    status
) VALUES (
    sqlc.arg(registration_form_id),
    sqlc.arg(date_of_intake),
    sqlc.arg(care_type),
    sqlc.arg(intake_participants)::text[]::intake_participants_enum[],
    sqlc.arg(family_situation),
    sqlc.arg(psychological_state),
    sqlc.arg(self_sufficiency),
    sqlc.arg(maturity_matrix_id),
    sqlc.arg(goals),
    sqlc.arg(risk_assessment),
    sqlc.arg(intake_conclusion),
    sqlc.arg(intake_conclusion_notes),
    sqlc.arg(signature),
    sqlc.arg(status)
)
ON CONFLICT (registration_form_id) DO UPDATE SET
    date_of_intake = EXCLUDED.date_of_intake,
    care_type = EXCLUDED.care_type,
    intake_participants = EXCLUDED.intake_participants,
    family_situation = EXCLUDED.family_situation,
    psychological_state = EXCLUDED.psychological_state,
    self_sufficiency = EXCLUDED.self_sufficiency,
    maturity_matrix_id = EXCLUDED.maturity_matrix_id,
    goals = EXCLUDED.goals,
    risk_assessment = EXCLUDED.risk_assessment,
    intake_conclusion = EXCLUDED.intake_conclusion,
    intake_conclusion_notes = EXCLUDED.intake_conclusion_notes,
    signature = EXCLUDED.signature,
    status = EXCLUDED.status,
    updated_at = NOW()
RETURNING
    id,
    registration_form_id,
    date_of_intake,
    care_type,
    intake_participants::text[] AS intake_participants,
    family_situation,
    psychological_state,
    self_sufficiency,
    maturity_matrix_id,
    goals,
    risk_assessment,
    intake_conclusion,
    intake_conclusion_notes,
    signature,
    created_at,
    updated_at,
    status,
    urgency_level;

-- name: UpsertIntakeFormDraft :one
INSERT INTO intake_forms (
    registration_form_id,
    date_of_intake,
    status
) VALUES (
    $1, $2, $3
)
ON CONFLICT (registration_form_id) DO UPDATE SET
    date_of_intake = EXCLUDED.date_of_intake,
    status = EXCLUDED.status,
    updated_at = NOW()
RETURNING
    id,
    registration_form_id,
    date_of_intake,
    care_type,
    intake_participants::text[] AS intake_participants,
    family_situation,
    psychological_state,
    self_sufficiency,
    maturity_matrix_id,
    goals,
    risk_assessment,
    intake_conclusion,
    intake_conclusion_notes,
    signature,
    created_at,
    updated_at,
    status,
    urgency_level;
    


-- name: ListIntakeForms :many
SELECT 
    i.id,
    i.registration_form_id,
    i.date_of_intake,
    i.care_type,
    i.intake_participants::text[] AS intake_participants,
    i.family_situation,
    i.psychological_state,
    i.self_sufficiency,
    i.maturity_matrix_id,
    i.goals,
    i.risk_assessment,
    i.intake_conclusion,
    i.intake_conclusion_notes,
    i.signature,
    i.created_at,
    i.updated_at,
    i.status,
    i.urgency_level,
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
SELECT
    id,
    registration_form_id,
    date_of_intake,
    care_type,
    intake_participants::text[] AS intake_participants,
    family_situation,
    psychological_state,
    self_sufficiency,
    maturity_matrix_id,
    goals,
    risk_assessment,
    intake_conclusion,
    intake_conclusion_notes,
    signature,
    created_at,
    updated_at,
    status,
    urgency_level
FROM intake_forms
WHERE id = $1;

-- name: UpdateIntakeFormOutcome :one
UPDATE intake_forms
SET
    status = $2,
    intake_conclusion = $3,
    urgency_level = $4,
    intake_conclusion_notes = COALESCE(sqlc.narg('intake_conclusion_notes'), intake_conclusion_notes),
    risk_assessment = COALESCE(sqlc.narg('risk_assessment'), risk_assessment),
    updated_at = NOW()
WHERE id = $1
RETURNING
    id,
    registration_form_id,
    date_of_intake,
    care_type,
    intake_participants::text[] AS intake_participants,
    family_situation,
    psychological_state,
    self_sufficiency,
    maturity_matrix_id,
    goals,
    risk_assessment,
    intake_conclusion,
    intake_conclusion_notes,
    signature,
    created_at,
    updated_at,
    status,
    urgency_level;
