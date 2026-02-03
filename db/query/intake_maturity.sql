-- name: CreateIntakeMaturityAssessment :one
INSERT INTO intake_maturity_assessments (
    intake_form_id,
    maturity_matrix_id,
    current_level,
    proposed_goals,
    notes
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;



-- name: GetIntakeMaturityAssessments :many
SELECT 
    ima.*,
    mm.topic_name
FROM intake_maturity_assessments ima
JOIN maturity_matrix mm ON ima.maturity_matrix_id = mm.id
WHERE ima.intake_form_id = $1
ORDER BY mm.topic_name;



-- name: GetIntakeMaturityAssessment :one
SELECT 
    ima.*,
    mm.topic_name
FROM intake_maturity_assessments ima
JOIN maturity_matrix mm ON ima.maturity_matrix_id = mm.id
WHERE ima.id = $1;



-- name: UpdateIntakeMaturityAssessment :one
UPDATE intake_maturity_assessments
SET
    current_level = COALESCE(sqlc.narg('current_level'), current_level),
    proposed_goals = COALESCE(sqlc.narg('proposed_goals'), proposed_goals),
    notes = COALESCE(sqlc.narg('notes'), notes)
WHERE id = sqlc.arg('id')
RETURNING *;



-- name: DeleteIntakeMaturityAssessment :exec
DELETE FROM intake_maturity_assessments
WHERE id = $1;



-- name: ListIntakeMaturityAssessmentsByIntake :many
SELECT 
    ima.*,
    mm.topic_name
FROM intake_maturity_assessments ima
JOIN maturity_matrix mm ON ima.maturity_matrix_id = mm.id
WHERE ima.intake_form_id = $1
ORDER BY mm.topic_name
LIMIT $2 OFFSET $3;


