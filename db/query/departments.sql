-- name: CreateDepartment :one
INSERT INTO departments (name, description, department_head_employee_id)
VALUES (sqlc.arg('name'), sqlc.narg('description'), sqlc.narg('department_head_employee_id'))
RETURNING *;

-- name: ListDepartments :many
SELECT
    d.*,
    COUNT(ep.id)::BIGINT AS employee_count
FROM departments d
LEFT JOIN employee_profile ep ON ep.department_id = d.id
GROUP BY d.id
ORDER BY name ASC;

-- name: GetDepartment :one
SELECT *
FROM departments
WHERE id = $1
LIMIT 1;

-- name: GetDepartmentByName :one
SELECT *
FROM departments
WHERE name = $1
LIMIT 1;

-- name: UpdateDepartment :one
UPDATE departments
SET
    name = COALESCE(sqlc.narg('name'), name),
    description = COALESCE(sqlc.narg('description'), description),
    department_head_employee_id = COALESCE(sqlc.narg('department_head_employee_id'), department_head_employee_id),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteDepartment :exec
DELETE FROM departments
WHERE id = $1;
