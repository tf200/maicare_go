-- name: CreateClientGoalsFromIntakeAssessments :many
WITH intake_rows AS (
    SELECT
        ita.id AS intake_assessment_id,
        ita.topic_id,
        t.topic_name,
        g.goal_item,
        g.goal_ord
    FROM intake_topic_assessments ita
    JOIN topics t ON t.id = ita.topic_id
    CROSS JOIN LATERAL jsonb_array_elements(COALESCE(ita.proposed_goals, '[]'::jsonb))
        WITH ORDINALITY AS g(goal_item, goal_ord)
    WHERE ita.intake_form_id = sqlc.arg(intake_form_id)::uuid
),
normalized AS (
    SELECT
        intake_assessment_id,
        topic_id,
        topic_name,
        goal_ord,
        CASE
            WHEN jsonb_typeof(goal_item) = 'object' THEN
                COALESCE(
                    NULLIF(BTRIM(goal_item->>'title'), ''),
                    NULLIF(BTRIM(goal_item->>'description'), '')
                )
            WHEN jsonb_typeof(goal_item) = 'string' THEN
                NULLIF(BTRIM(TRIM(BOTH '"' FROM goal_item::text)), '')
            ELSE NULL
        END AS title,
        CASE
            WHEN jsonb_typeof(goal_item) = 'object' THEN NULLIF(BTRIM(goal_item->>'description'), '')
            ELSE NULL
        END AS description,
        CASE
            WHEN jsonb_typeof(goal_item) = 'object'
                 AND LOWER(COALESCE(goal_item->>'priority', 'medium')) IN ('low', 'medium', 'high')
                THEN LOWER(goal_item->>'priority')
            ELSE 'medium'
        END AS priority
    FROM intake_rows
),
inserted AS (
    INSERT INTO client_goals (
        client_id,
        title,
        description,
        priority,
        status,
        topic_id,
        topic_name_snapshot,
        source,
        origin_intake_assessment_id,
        sort_order
    )
    SELECT
        sqlc.arg(client_id)::uuid,
        n.title,
        n.description,
        n.priority::client_goal_priority_enum,
        'active'::client_goal_status_enum,
        n.topic_id,
        n.topic_name,
        'intake'::client_goal_source_enum,
        n.intake_assessment_id,
        ROW_NUMBER() OVER (ORDER BY n.topic_name, n.intake_assessment_id, n.goal_ord)::int - 1
    FROM normalized n
    WHERE n.title IS NOT NULL
    RETURNING *
)
SELECT * FROM inserted;

-- name: ListActiveGoalsByClientID :many
SELECT *
FROM client_goals
WHERE client_id = $1
  AND status = 'active'
ORDER BY sort_order;

-- name: GetNextActiveClientGoalSortOrder :one
SELECT COALESCE(MAX(sort_order), -1)::int + 1
FROM client_goals
WHERE client_id = $1
  AND status = 'active';

-- name: CreateManualClientGoal :one
INSERT INTO client_goals (
    client_id,
    title,
    description,
    priority,
    status,
    topic_id,
    topic_name_snapshot,
    source,
    sort_order
) VALUES (
    $1,
    $2,
    $3,
    $4,
    'active',
    $5,
    $6,
    'manual',
    $7
)
RETURNING *;

-- name: GetClientGoalByIDAndClientID :one
SELECT *
FROM client_goals
WHERE id = $1
  AND client_id = $2
LIMIT 1;

-- name: GoalHasEvaluationItems :one
SELECT public.goal_has_evaluation_history_for_update(
    sqlc.arg(goal_id)::uuid,
    sqlc.arg(client_id)::uuid
);

-- name: ClientHasDraftEvaluationForGoalUpdate :one
SELECT public.client_has_draft_evaluation_for_goal_update($1::uuid);

-- name: UpdateClientGoalByID :one
UPDATE client_goals
SET
    title = $3,
    description = $4,
    priority = $5,
    topic_id = $6,
    topic_name_snapshot = $7,
    sort_order = $8,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND client_id = $2
RETURNING *;

-- name: CancelClientGoalByID :one
UPDATE client_goals
SET
    status = 'cancelled',
    archived_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND client_id = $2
RETURNING *;

-- name: CreateReviewUpdatedClientGoal :one
INSERT INTO client_goals (
    client_id,
    title,
    description,
    priority,
    status,
    topic_id,
    topic_name_snapshot,
    source,
    sort_order
) VALUES (
    $1,
    $2,
    $3,
    $4,
    'active',
    $5,
    $6,
    'review_update',
    $7
)
RETURNING *;
