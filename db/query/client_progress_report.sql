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