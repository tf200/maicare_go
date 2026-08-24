-- name: GetUserIDByEmployeeID :one
SELECT user_id FROM employee_profile
WHERE id = $1 LIMIT 1;

-- name: ListUserIDsByEmployeeIDs :many
SELECT user_id
FROM employee_profile
WHERE id = ANY($1::uuid[]);

-- name: SetEmployeeProfilePicture :one
UPDATE custom_user
SET profile_picture = $2
WHERE id = (
    SELECT user_id
    FROM employee_profile
    WHERE employee_profile.id = $1
)
RETURNING *;

-- name: SearchEmployeesByNameOrEmail :many
SELECT
    employee_profile.id,
    employee_profile.first_name,
    employee_profile.last_name,
    employee_profile.work_email_address
FROM employee_profile
JOIN custom_user cu ON cu.id = employee_profile.user_id
WHERE
    employee_profile.id <> sqlc.arg('system_employee_id') AND
    (employee_profile.first_name ILIKE '%' || @search || '%' OR
     employee_profile.last_name ILIKE '%' || @search || '%' OR
     cu.email ILIKE '%' || @search || '%')
LIMIT 10;

-- name: GetEmployeeCounts :one
SELECT
    COUNT(*) FILTER (WHERE is_subcontractor IS NOT TRUE) AS total_employees,
    COUNT(*) FILTER (WHERE is_subcontractor = TRUE) AS total_subcontractors,
    COUNT(*) FILTER (WHERE is_archived = TRUE) AS total_archived,
    COUNT(*) FILTER (WHERE out_of_service = TRUE) AS total_out_of_service
FROM
    employee_profile
WHERE id <> sqlc.arg('system_employee_id');

-- name: ListEmployeesWithContractHours :many
SELECT
    id,
    first_name,
    last_name,
    contract_hours
FROM employee_profile
WHERE id = ANY($1::uuid[])
AND contract_hours IS NOT NULL
AND contract_hours > 0;

-- name: ListEmployeeNamesByIDs :many
SELECT
    id,
    first_name,
    last_name
FROM employee_profile
WHERE id = ANY(sqlc.arg(employee_ids)::uuid[]);
