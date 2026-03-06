-- name: CreateHandbookTemplateForDepartment :one
WITH next_version AS (
    SELECT COALESCE(MAX(version), 0) + 1 AS v
    FROM handbook_templates
    WHERE department_id = sqlc.arg('department_id')
), deactivated AS (
    UPDATE handbook_templates
    SET is_active = FALSE, updated_at = CURRENT_TIMESTAMP
    WHERE department_id = sqlc.arg('department_id') AND is_active = TRUE
)
INSERT INTO handbook_templates (
    department_id,
    title,
    description,
    version,
    is_active,
    created_by_employee_id
)
SELECT
    sqlc.arg('department_id'),
    sqlc.arg('title'),
    sqlc.narg('description'),
    next_version.v,
    TRUE,
    sqlc.narg('created_by_employee_id')
FROM next_version
RETURNING *;

-- name: ListHandbookTemplatesByDepartment :many
SELECT *
FROM handbook_templates
WHERE department_id = $1
ORDER BY version DESC;

-- name: GetActiveHandbookTemplateByDepartment :one
SELECT *
FROM handbook_templates
WHERE department_id = $1 AND is_active = TRUE
LIMIT 1;

-- name: GetHandbookTemplateByID :one
SELECT *
FROM handbook_templates
WHERE id = $1
LIMIT 1;

-- name: CreateHandbookStep :one
INSERT INTO handbook_steps (
    template_id,
    sort_order,
    kind,
    title,
    body,
    content,
    is_required
)
VALUES (
    sqlc.arg('template_id'),
    sqlc.arg('sort_order'),
    sqlc.arg('kind'),
    sqlc.arg('title'),
    sqlc.narg('body'),
    COALESCE(sqlc.narg('content'), '{}'::jsonb),
    COALESCE(sqlc.narg('is_required')::boolean, TRUE)
)
RETURNING *;

-- name: ListHandbookStepsByTemplate :many
SELECT *
FROM handbook_steps
WHERE template_id = $1
ORDER BY sort_order ASC;

-- name: CreateEmployeeHandbookFromTemplate :one
WITH hb AS (
    INSERT INTO employee_handbooks (
        employee_id,
        template_id,
        assigned_by_employee_id
    )
    VALUES (
        sqlc.arg('employee_id'),
        sqlc.arg('template_id'),
        sqlc.narg('assigned_by_employee_id')
    )
    RETURNING *
), progress AS (
    INSERT INTO employee_handbook_step_progress (employee_handbook_id, step_id)
    SELECT hb.id, hs.id
    FROM hb
    JOIN handbook_steps hs ON hs.template_id = hb.template_id
)
SELECT * FROM hb;

-- name: GetActiveEmployeeHandbookByEmployeeID :one
SELECT
    eh.*,
    ht.title AS template_title,
    ht.description AS template_description,
    ht.version AS template_version,
    d.id AS department_id,
    d.name AS department_name
FROM employee_handbooks eh
JOIN handbook_templates ht ON ht.id = eh.template_id
JOIN departments d ON d.id = ht.department_id
WHERE eh.employee_id = $1
  AND eh.status IN ('not_started', 'in_progress')
ORDER BY eh.assigned_at DESC
LIMIT 1;

-- name: ListEmployeeHandbookStepsByHandbookID :many
SELECT
    hs.id AS step_id,
    hs.sort_order,
    hs.kind,
    hs.title,
    hs.body,
    hs.content,
    hs.is_required,
    ehsp.status AS progress_status,
    ehsp.started_at AS progress_started_at,
    ehsp.completed_at AS progress_completed_at,
    ehsp.response AS progress_response
FROM employee_handbook_step_progress ehsp
JOIN handbook_steps hs ON hs.id = ehsp.step_id
WHERE ehsp.employee_handbook_id = $1
ORDER BY hs.sort_order ASC;

-- name: MarkEmployeeHandbookStarted :one
UPDATE employee_handbooks
SET
    status = 'in_progress',
    started_at = COALESCE(started_at, CURRENT_TIMESTAMP)
WHERE id = $1
RETURNING *;

-- name: WaiveActiveEmployeeHandbooksByEmployeeID :exec
UPDATE employee_handbooks
SET
    status = 'waived',
    completed_at = COALESCE(completed_at, CURRENT_TIMESTAMP)
WHERE employee_id = $1
  AND status IN ('not_started', 'in_progress');

-- name: CompleteEmployeeHandbookStep :one
UPDATE employee_handbook_step_progress
SET
    status = 'completed',
    started_at = COALESCE(started_at, CURRENT_TIMESTAMP),
    completed_at = CURRENT_TIMESTAMP,
    response = COALESCE(sqlc.narg('response'), response)
WHERE employee_handbook_id = sqlc.arg('employee_handbook_id')
  AND step_id = sqlc.arg('step_id')
RETURNING *;

-- name: CountRemainingRequiredHandbookSteps :one
SELECT COUNT(*)::INT
FROM employee_handbook_step_progress p
JOIN handbook_steps s ON s.id = p.step_id
WHERE p.employee_handbook_id = $1
  AND s.is_required = TRUE
  AND p.status <> 'completed';

-- name: MarkEmployeeHandbookCompleted :one
UPDATE employee_handbooks
SET
    status = 'completed',
    completed_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;
