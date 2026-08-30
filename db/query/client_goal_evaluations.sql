-- name: CreateGoalEvaluation :one
INSERT INTO client_goal_evaluations (
    client_id,
    evaluation_date,
    period_start,
    period_end,
    evaluation_interval_weeks,
    status,
    overall_notes,
    created_by_employee_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetDraftGoalEvaluationByClientAndDate :one
SELECT *
FROM client_goal_evaluations
WHERE client_id = $1
  AND evaluation_date = $2
  AND status = 'draft'
ORDER BY updated_at DESC
LIMIT 1;

-- name: GetGoalEvaluationByClientAndDate :one
SELECT *
FROM client_goal_evaluations
WHERE client_id = $1
  AND evaluation_date = $2
ORDER BY updated_at DESC
LIMIT 1;

-- name: UpdateGoalEvaluation :one
UPDATE client_goal_evaluations
SET
    evaluation_date = COALESCE(sqlc.narg('evaluation_date'), evaluation_date),
    period_start = COALESCE(sqlc.narg('period_start'), period_start),
    period_end = COALESCE(sqlc.narg('period_end'), period_end),
    evaluation_interval_weeks = COALESCE(sqlc.narg('evaluation_interval_weeks'), evaluation_interval_weeks),
    status = COALESCE(sqlc.narg('status'), status),
    overall_notes = COALESCE(sqlc.narg('overall_notes'), overall_notes),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;

-- name: UpdateGoalEvaluationDraftCAS :one
UPDATE client_goal_evaluations
SET
    overall_notes = COALESCE(sqlc.narg('overall_notes'), overall_notes),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND status = 'draft'
  AND updated_at = sqlc.arg('expected_updated_at')
RETURNING *;

-- name: SubmitGoalEvaluationDraftCAS :one
UPDATE client_goal_evaluations
SET
    status = 'completed',
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND status = 'draft'
  AND updated_at = sqlc.arg('expected_updated_at')
RETURNING *;

-- name: UpsertGoalEvaluationItem :one
INSERT INTO client_goal_evaluation_items (
    client_id,
    evaluation_id,
    goal_id,
    progress,
    notes
) VALUES (
    $1, $2, $3, $4, $5
)
ON CONFLICT (evaluation_id, goal_id) DO UPDATE
SET
    progress = EXCLUDED.progress,
    notes = EXCLUDED.notes,
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetGoalEvaluationItems :many
SELECT 
    ei.*,
    g.title AS goal_title,
    g.description AS goal_description,
    g.topic_name_snapshot
FROM client_goal_evaluation_items ei
JOIN client_goals g ON ei.goal_id = g.id
WHERE ei.evaluation_id = $1
ORDER BY g.sort_order;

-- name: GetGoalEvaluationByID :one
SELECT
    e.*,
    ep.first_name AS creator_first_name,
    ep.last_name AS creator_last_name
FROM client_goal_evaluations e
LEFT JOIN employee_profile ep
    ON ep.id = e.created_by_employee_id
WHERE e.id = $1
LIMIT 1;

-- name: ListUpcomingEvaluationsForCoordinator :many
SELECT
    c.id AS client_id,
    c.first_name AS client_first_name,
    c.last_name AS client_last_name,
    c.next_evaluation_date,
    (c.next_evaluation_date - CURRENT_DATE)::int4 AS days_left,
    CASE
        WHEN (c.next_evaluation_date - CURRENT_DATE) <= 3 THEN 'critical'
        ELSE 'normal'
    END AS priority,
    COALESCE((d.id IS NOT NULL), false)::bool AS has_draft,
    COALESCE(di.filled_goals_count, 0)::int4 AS filled_goals_count,
    COALESCE(di.total_goals_count, 0)::int4 AS total_goals_count,
    COUNT(*) OVER() AS total_count
FROM assigned_employee ae
JOIN client_details c
    ON c.id = ae.client_id
LEFT JOIN LATERAL (
    SELECT e.id
    FROM client_goal_evaluations e
    WHERE e.client_id = c.id
      AND e.status = 'draft'
    ORDER BY e.updated_at DESC
    LIMIT 1
) d ON true
LEFT JOIN LATERAL (
    SELECT
        COUNT(*) FILTER (WHERE i.progress <> 'no_progress') AS filled_goals_count,
        COUNT(*) AS total_goals_count
    FROM client_goal_evaluation_items i
    WHERE i.evaluation_id = d.id
) di ON true
WHERE ae.employee_id = $1
  AND ae.role = 'coordinator'
  AND c.status = 'in_care'
  AND c.next_evaluation_date IS NOT NULL
ORDER BY c.next_evaluation_date ASC, c.first_name ASC, c.last_name ASC
LIMIT $2 OFFSET $3;

-- name: ListRecentSubmittedEvaluationsByEmployee :many
SELECT
    e.id,
    e.client_id,
    c.first_name AS client_first_name,
    c.last_name AS client_last_name,
    e.evaluation_date,
    e.updated_at AS submitted_at,
    c.next_evaluation_date,
    COALESCE(di.filled_goals_count, 0)::int4 AS filled_goals_count,
    COALESCE(di.total_goals_count, 0)::int4 AS total_goals_count,
    COUNT(*) OVER() AS total_count
FROM client_goal_evaluations e
JOIN client_details c
    ON c.id = e.client_id
LEFT JOIN LATERAL (
    SELECT
        COUNT(*) FILTER (WHERE i.progress <> 'no_progress') AS filled_goals_count,
        COUNT(*) AS total_goals_count
    FROM client_goal_evaluation_items i
    WHERE i.evaluation_id = e.id
) di ON true
WHERE e.created_by_employee_id = $1
  AND e.status = 'completed'
ORDER BY e.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListRecentDraftEvaluationsByEmployee :many
SELECT
    e.id,
    e.client_id,
    c.first_name AS client_first_name,
    c.last_name AS client_last_name,
    e.evaluation_date,
    e.updated_at,
    (e.evaluation_date - CURRENT_DATE)::int4 AS days_left,
    CASE
        WHEN (e.evaluation_date - CURRENT_DATE) <= 3 THEN 'critical'
        ELSE 'normal'
    END AS priority,
    COALESCE(di.filled_goals_count, 0)::int4 AS filled_goals_count,
    COALESCE(di.total_goals_count, 0)::int4 AS total_goals_count,
    COUNT(*) OVER() AS total_count
FROM client_goal_evaluations e
JOIN client_details c
    ON c.id = e.client_id
LEFT JOIN LATERAL (
    SELECT
        COUNT(*) FILTER (WHERE i.progress <> 'no_progress') AS filled_goals_count,
        COUNT(*) AS total_goals_count
    FROM client_goal_evaluation_items i
    WHERE i.evaluation_id = e.id
) di ON true
WHERE e.created_by_employee_id = $1
  AND e.status = 'draft'
ORDER BY e.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: GetEvaluationStatsByEmployee :one
SELECT
    (
        SELECT COUNT(DISTINCT c.id)
        FROM assigned_employee ae
        JOIN client_details c ON c.id = ae.client_id
        WHERE ae.employee_id = sqlc.arg(employee_id)::uuid
          AND ae.role = 'coordinator'
          AND c.status = 'in_care'
          AND c.next_evaluation_date IS NOT NULL
          AND c.next_evaluation_date <= CURRENT_DATE + 3
    )::int8 AS attention_required,
    (
        SELECT COUNT(DISTINCT e.id)
        FROM client_goal_evaluations e
        JOIN client_details c ON c.id = e.client_id
        WHERE e.created_by_employee_id = sqlc.arg(employee_id)::uuid
          AND e.status = 'draft'
          AND c.status = 'in_care'
          AND c.next_evaluation_date IS NOT NULL
          AND e.evaluation_date = c.next_evaluation_date
    )::int8 AS in_progress,
    (
        SELECT COUNT(DISTINCT e.id)
        FROM client_goal_evaluations e
        WHERE e.created_by_employee_id = sqlc.arg(employee_id)::uuid
          AND e.status = 'completed'
          AND e.updated_at >= CURRENT_TIMESTAMP - INTERVAL '30 days'
    )::int8 AS recently_finalized,
    CURRENT_TIMESTAMP::timestamptz AS as_of;

-- name: GetLatestDraftEvaluationByClient :one
SELECT
    e.id,
    e.client_id,
    e.evaluation_date,
    e.updated_at
FROM client_goal_evaluations e
WHERE e.client_id = $1
  AND e.status = 'draft'
ORDER BY e.updated_at DESC
LIMIT 1;

-- name: GetCurrentCycleDraftEvaluationByClientAndEmployee :one
SELECT
    e.*
FROM client_goal_evaluations e
JOIN client_details c
    ON c.id = e.client_id
WHERE e.client_id = $1
  AND e.created_by_employee_id = $2
  AND e.status = 'draft'
  AND c.next_evaluation_date IS NOT NULL
  AND e.evaluation_date = c.next_evaluation_date
ORDER BY e.updated_at DESC
LIMIT 1;

-- name: GetLatestCompletedEvaluationByClient :one
SELECT
    e.id,
    e.client_id,
    e.evaluation_date,
    e.overall_notes,
    e.created_by_employee_id,
    e.updated_at AS submitted_at,
    ep.first_name AS creator_first_name,
    ep.last_name AS creator_last_name
FROM client_goal_evaluations e
LEFT JOIN employee_profile ep ON ep.id = e.created_by_employee_id
WHERE e.client_id = $1
  AND e.status = 'completed'
ORDER BY e.evaluation_date DESC, e.updated_at DESC
LIMIT 1;

-- name: ListSubmittedEvaluationsByClient :many
SELECT
    e.id,
    e.client_id,
    e.evaluation_date,
    e.updated_at AS submitted_at,
    e.created_by_employee_id,
    ep.first_name AS creator_first_name,
    ep.last_name AS creator_last_name,
    COALESCE(di.filled_goals_count, 0)::int4 AS filled_goals_count,
    COALESCE(di.total_goals_count, 0)::int4 AS total_goals_count,
    COUNT(*) OVER() AS total_count
FROM client_goal_evaluations e
LEFT JOIN employee_profile ep
    ON ep.id = e.created_by_employee_id
LEFT JOIN LATERAL (
    SELECT
        COUNT(*) FILTER (WHERE i.progress <> 'no_progress') AS filled_goals_count,
        COUNT(*) AS total_goals_count
    FROM client_goal_evaluation_items i
    WHERE i.evaluation_id = e.id
) di ON true
WHERE e.client_id = $1
  AND e.status = 'completed'
ORDER BY e.evaluation_date DESC, e.updated_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLatestCompletedGoalProgressByClient :many
SELECT DISTINCT ON (i.goal_id)
    i.goal_id,
    i.progress,
    i.notes
FROM client_goal_evaluation_items i
JOIN client_goal_evaluations e ON e.id = i.evaluation_id
WHERE e.client_id = $1
  AND e.status = 'completed'
ORDER BY i.goal_id, e.evaluation_date DESC, e.updated_at DESC;

-- name: ListGoalEvaluationHistoryByClientAndGoal :many
SELECT
    e.id AS evaluation_id,
    e.evaluation_date,
    e.updated_at AS submitted_at,
    i.progress,
    i.notes,
    e.created_by_employee_id,
    ep.first_name AS creator_first_name,
    ep.last_name AS creator_last_name,
    e.period_start,
    e.period_end,
    COUNT(*) OVER() AS total_count
FROM client_goal_evaluation_items i
JOIN client_goal_evaluations e
    ON e.id = i.evaluation_id
JOIN client_goals g
    ON g.id = i.goal_id
LEFT JOIN employee_profile ep
    ON ep.id = e.created_by_employee_id
WHERE e.client_id = $1
  AND i.goal_id = $2
  AND g.client_id = $1
  AND e.status = 'completed'
ORDER BY e.evaluation_date DESC, e.updated_at DESC
LIMIT $3 OFFSET $4;
