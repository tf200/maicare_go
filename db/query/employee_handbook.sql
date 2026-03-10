-- name: CreateHandbookTemplateForDepartment :one
WITH next_version AS (
    SELECT COALESCE(MAX(version), 0) + 1 AS v
    FROM handbook_templates
    WHERE department_id = sqlc.arg('department_id')
)
INSERT INTO handbook_templates (
    department_id,
    title,
    description,
    version,
    status,
    created_by_employee_id
)
SELECT
    sqlc.arg('department_id'),
    sqlc.arg('title'),
    sqlc.narg('description'),
    next_version.v,
    'draft',
    sqlc.narg('created_by_employee_id')
FROM next_version
RETURNING *;

-- name: CloneHandbookTemplateToDraft :one
WITH source_template AS (
    SELECT id, department_id, title, description
    FROM handbook_templates
    WHERE handbook_templates.id = sqlc.arg('source_template_id')
    LIMIT 1
), next_version AS (
    SELECT COALESCE(MAX(version), 0) + 1 AS v
    FROM handbook_templates
    WHERE department_id = (SELECT department_id FROM source_template)
), cloned_template AS (
    INSERT INTO handbook_templates (
        department_id,
        title,
        description,
        version,
        status,
        created_by_employee_id
    )
    SELECT
        st.department_id,
        st.title,
        st.description,
        nv.v,
        'draft',
        sqlc.narg('created_by_employee_id')
    FROM source_template st
    CROSS JOIN next_version nv
    RETURNING *
), cloned_steps AS (
    INSERT INTO handbook_steps (
        template_id,
        sort_order,
        kind,
        title,
        body,
        content,
        is_required
    )
    SELECT
        ct.id,
        hs.sort_order,
        hs.kind,
        hs.title,
        hs.body,
        hs.content,
        hs.is_required
    FROM cloned_template ct
    JOIN source_template st ON TRUE
    JOIN handbook_steps hs ON hs.template_id = st.id
    RETURNING 1
)
SELECT * FROM cloned_template;

-- name: PublishHandbookTemplate :one
WITH target AS (
    SELECT id, department_id
    FROM handbook_templates
    WHERE handbook_templates.id = sqlc.arg('template_id')
      AND status = 'draft'
    LIMIT 1
), archived AS (
    UPDATE handbook_templates ht
    SET
        status = 'archived',
        archived_at = CURRENT_TIMESTAMP,
        updated_at = CURRENT_TIMESTAMP
    FROM target t
    WHERE ht.department_id = t.department_id
      AND ht.status = 'published'
)
UPDATE handbook_templates ht
SET
    status = 'published',
    published_by_employee_id = sqlc.narg('published_by_employee_id'),
    published_at = CURRENT_TIMESTAMP,
    archived_at = NULL,
    updated_at = CURRENT_TIMESTAMP
FROM target t
WHERE ht.id = t.id
RETURNING ht.*;

-- name: ListHandbookTemplatesByDepartment :many
SELECT *
FROM handbook_templates
WHERE department_id = $1
ORDER BY version DESC;

-- name: GetActiveHandbookTemplateByDepartment :one
SELECT *
FROM handbook_templates
WHERE department_id = $1 AND status = 'published'
LIMIT 1;

-- name: GetHandbookTemplateByID :one
SELECT *
FROM handbook_templates
WHERE id = $1
LIMIT 1;

-- name: UpdateHandbookTemplateMetadata :one
UPDATE handbook_templates
SET
    title = CASE
        WHEN sqlc.arg('set_title')::boolean THEN COALESCE(sqlc.narg('title'), title)
        ELSE title
    END,
    description = CASE
        WHEN sqlc.arg('set_description')::boolean THEN sqlc.narg('description')
        ELSE description
    END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('template_id')
  AND status = 'draft'
RETURNING *;

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

-- name: GetHandbookStepByID :one
SELECT *
FROM handbook_steps
WHERE id = $1
LIMIT 1;

-- name: UpdateHandbookStepByID :one
UPDATE handbook_steps
SET
    title = CASE
        WHEN sqlc.arg('set_title')::boolean THEN COALESCE(sqlc.narg('title'), title)
        ELSE title
    END,
    body = CASE
        WHEN sqlc.arg('set_body')::boolean THEN sqlc.narg('body')
        ELSE body
    END,
    content = CASE
        WHEN sqlc.arg('set_content')::boolean THEN COALESCE(sqlc.narg('content'), 'null'::jsonb)
        ELSE content
    END,
    is_required = CASE
        WHEN sqlc.arg('set_is_required')::boolean THEN COALESCE(sqlc.narg('is_required')::boolean, is_required)
        ELSE is_required
    END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('step_id')
RETURNING *;

-- name: UpdateHandbookStepSortOrder :exec
UPDATE handbook_steps
SET
    sort_order = sqlc.arg('sort_order'),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('step_id');

-- name: DeleteHandbookStepByID :exec
DELETE FROM handbook_steps
WHERE id = $1;

-- name: CreateEmployeeHandbookFromTemplate :one
WITH hb AS (
    INSERT INTO employee_handbooks (
        employee_id,
        template_id,
        template_version,
        assigned_by_employee_id
    )
    SELECT
        sqlc.arg('employee_id'),
        t.id,
        t.version,
        sqlc.narg('assigned_by_employee_id')
    FROM handbook_templates t
    WHERE t.id = sqlc.arg('template_id')
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

-- name: CreateEmployeeHandbookAssignmentHistory :one
INSERT INTO employee_handbook_assignment_history (
    employee_handbook_id,
    employee_id,
    template_id,
    template_version,
    event,
    actor_employee_id,
    metadata
)
VALUES (
    sqlc.narg('employee_handbook_id'),
    sqlc.arg('employee_id'),
    sqlc.arg('template_id'),
    sqlc.arg('template_version'),
    sqlc.arg('event'),
    sqlc.narg('actor_employee_id'),
    COALESCE(sqlc.narg('metadata'), '{}'::jsonb)
)
RETURNING *;

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

-- name: CountHandbookStepsByTemplateID :one
SELECT COUNT(*)::INT
FROM handbook_steps
WHERE template_id = $1;
