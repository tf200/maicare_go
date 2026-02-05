-- name: CreateIntakeMaturityAssessment :one
INSERT INTO intake_topic_assessments (
    intake_form_id,
    topic_id,
    current_level,
    proposed_goals,
    notes
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;


-- name: CreateIntakeMaturityAssessmentsBatch :many
WITH inserted AS (
    INSERT INTO intake_topic_assessments (
        intake_form_id,
        topic_id,
        current_level,
        proposed_goals,
        notes
    )
    SELECT
        sqlc.arg(intake_form_id)::uuid,
        (item->>'topic_id')::uuid,
        (item->>'current_level')::int,
        COALESCE(item->'proposed_goals', '[]'::jsonb),
        item->>'notes'
    FROM jsonb_array_elements(sqlc.arg(items)::jsonb) AS item
    RETURNING id, intake_form_id, topic_id, current_level, proposed_goals, notes, created_at
)
SELECT
    inserted.*,
    t.topic_name
FROM inserted
JOIN topics t ON inserted.topic_id = t.id
ORDER BY t.topic_name;



-- name: GetIntakeMaturityAssessments :many
SELECT 
    ima.*,
    t.topic_name
FROM intake_topic_assessments ima
JOIN topics t ON ima.topic_id = t.id
WHERE ima.intake_form_id = $1
ORDER BY t.topic_name;



-- name: GetIntakeMaturityAssessment :one
SELECT 
    ima.*,
    t.topic_name
FROM intake_topic_assessments ima
JOIN topics t ON ima.topic_id = t.id
WHERE ima.id = $1;



-- name: UpdateIntakeMaturityAssessment :one
UPDATE intake_topic_assessments
SET
    current_level = COALESCE(sqlc.narg('current_level'), current_level),
    proposed_goals = COALESCE(sqlc.narg('proposed_goals'), proposed_goals),
    notes = COALESCE(sqlc.narg('notes'), notes)
WHERE id = sqlc.arg('id')
RETURNING *;



-- name: DeleteIntakeMaturityAssessment :exec
DELETE FROM intake_topic_assessments
WHERE id = $1;



-- name: ListIntakeMaturityAssessmentsByIntake :many
SELECT 
    ima.*,
    t.topic_name
FROM intake_topic_assessments ima
JOIN topics t ON ima.topic_id = t.id
WHERE ima.intake_form_id = $1
ORDER BY t.topic_name
LIMIT $2 OFFSET $3;
