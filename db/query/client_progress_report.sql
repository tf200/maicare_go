-- name: CreateProgressReport :one
INSERT INTO progress_report (
        client_id,
        employee_id,
        title,
        date,
        report_text,
        type,
        emotional_state
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7
    ) RETURNING *;


-- name: ListProgressReports :many
SELECT 
    pr.*,
    COUNT(*) OVER() AS total_count,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name
FROM progress_report pr
JOIN employee_profile e ON pr.employee_id = e.id
WHERE pr.client_id = $1
ORDER BY pr.date DESC
LIMIT $2 OFFSET $3;


-- name: GetProgressReport :one
SELECT 
    pr.*,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name
FROM progress_report pr
JOIN employee_profile e ON pr.employee_id = e.id
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


-- name: GetProgressReportsByDateRange :many
SELECT *
FROM progress_report
WHERE client_id = @client_id
  AND date >= @start_date
  AND date <= @end_date
ORDER BY date ASC;



-- name: CreateAiGeneratedReport :one
INSERT INTO ai_generated_reports (
        client_id,
        report_text,
        start_date,
        end_date

    ) VALUES (
        $1, $2, $3, $4
    ) RETURNING *;


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






