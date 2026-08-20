-- name: CreateProgressReport :one
WITH new_report AS (
    SELECT public.begin_progress_report_creation($1) AS id
)
INSERT INTO progress_report (
        id,
        client_id,
        employee_id,
        title,
        date,
        report_text,
        type,
        emotional_state
    ) SELECT
        new_report.id, $1, $2, $3, $4, $5, $6, $7
    FROM new_report
    WHERE new_report.id IS NOT NULL
    RETURNING *;


-- name: ListProgressReports :many
SELECT 
    pr.*,
    COUNT(*) OVER() AS total_count,
    COALESCE(e.first_name, '') AS employee_first_name,
    COALESCE(e.last_name, '') AS employee_last_name,
    u.profile_picture AS employee_profile_picture
FROM progress_report pr
LEFT JOIN employee_profile e ON pr.employee_id = e.id
LEFT JOIN custom_User u ON e.user_id = u.id
WHERE pr.client_id = sqlc.arg('client_id')
  AND (
    sqlc.narg('type')::progress_report_type_enum IS NULL
    OR pr.type = sqlc.narg('type')::progress_report_type_enum
  )
ORDER BY pr.date DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');


-- name: GetProgressReport :one
SELECT 
    pr.*,
    COALESCE(e.first_name, '') AS employee_first_name,
    COALESCE(e.last_name, '') AS employee_last_name,
    u.profile_picture AS employee_profile_picture
FROM progress_report pr
LEFT JOIN employee_profile e ON pr.employee_id = e.id
LEFT JOIN custom_User u ON e.user_id = u.id
WHERE pr.id = $1 LIMIT 1;

-- name: UpdateProgressReport :one
UPDATE progress_report
SET
    employee_id = COALESCE(sqlc.narg('employee_id'), employee_id),
    title = COALESCE(sqlc.narg('title'), title),
    date = COALESCE(sqlc.narg('date'), date),
    report_text = COALESCE(sqlc.narg('report_text'), report_text),
    type = COALESCE(sqlc.narg('type'), type),
    emotional_state = COALESCE(sqlc.narg('emotional_state'), emotional_state)
WHERE id = $1
RETURNING *;

-- name: DeleteProgressReport :exec
DELETE FROM progress_report
WHERE id = $1;


-- name: GetProgressReportsByDateRange :many
SELECT *
FROM progress_report
WHERE client_id = @client_id
  AND date >= @start_date
  AND date <= @end_date
ORDER BY date ASC;



-- name: CreateAiGeneratedReport :one
WITH new_report AS (
    SELECT public.begin_ai_report_creation($1) AS id
)
INSERT INTO ai_generated_reports (
        id,
        client_id,
        report_text,
        start_date,
        end_date

    ) SELECT
        new_report.id, $1, $2, $3, $4
    FROM new_report
    WHERE new_report.id IS NOT NULL
    RETURNING *;


-- name: ListAiGeneratedReports :many
SELECT 
    agr.*,
    COUNT(*) OVER() AS total_count
FROM ai_generated_reports agr
WHERE agr.client_id = $1
ORDER BY agr.created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetAiGeneratedReport :one
SELECT 
    agr.*
FROM ai_generated_reports agr
WHERE agr.id = $1 LIMIT 1;



