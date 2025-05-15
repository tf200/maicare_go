-- name: ListAllIncidents :many
SELECT 
    i.*,
    c.first_name AS client_first_name,
    c.last_name AS client_last_name,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name
FROM 
    incident i
JOIN 
    client_details c ON i.client_id = c.id
JOIN 
    employee_profile e ON i.employee_id = e.id
WHERE 
    i.soft_delete = false
    AND (
        sqlc.arg('is_confirmed')::boolean IS NULL 
        OR i.is_confirmed = sqlc.arg('is_confirmed')::boolean
    )
ORDER BY 
    i.incident_date DESC
LIMIT $1
OFFSET $2;


-- name: CountAllIncidents :one
SELECT COUNT(*) as total_count
FROM incident i
WHERE i.soft_delete = false
AND (
    sqlc.arg('is_confirmed')::boolean IS NULL 
    OR i.is_confirmed = sqlc.arg('is_confirmed')::boolean
);