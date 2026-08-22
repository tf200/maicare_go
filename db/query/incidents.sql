-- name: ListAllIncidents :many
SELECT 
    i.id,
    i.client_id,
    i.occurred_at,
    i.incident_type,
    i.severity_of_incident,
    i.is_confirmed,
    c.first_name AS client_first_name,
    c.bsn AS client_bsn,
    c.last_name AS client_last_name,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name,
    l.name AS location_name
FROM 
    incident i
JOIN 
    client_details c ON i.client_id = c.id
JOIN 
    employee_profile e ON i.employee_id = e.id
JOIN
    location l ON i.location_id = l.id
WHERE 
    (
        sqlc.narg('is_confirmed')::boolean IS NULL
        OR i.is_confirmed = sqlc.narg('is_confirmed')::boolean
    )
    AND (
        sqlc.narg('search')::text IS NULL
        OR sqlc.narg('search')::text = ''
        OR c.first_name ILIKE '%' || sqlc.narg('search')::text || '%'
        OR c.last_name ILIKE '%' || sqlc.narg('search')::text || '%'
    )
ORDER BY 
    i.occurred_at DESC
LIMIT $1
OFFSET $2;


-- name: CountAllIncidents :one
SELECT COUNT(*) as total_count
FROM incident i
JOIN client_details c ON i.client_id = c.id
WHERE (
    sqlc.narg('is_confirmed')::boolean IS NULL
    OR i.is_confirmed = sqlc.narg('is_confirmed')::boolean
)
AND (
    sqlc.narg('search')::text IS NULL
    OR sqlc.narg('search')::text = ''
    OR c.first_name ILIKE '%' || sqlc.narg('search')::text || '%'
    OR c.last_name ILIKE '%' || sqlc.narg('search')::text || '%'
);

-- name: GetIncidentCounts :one
SELECT
    COUNT(*) FILTER (
        WHERE i.severity_of_incident IN ('serious', 'fatal')
    )::BIGINT AS serious_fatal_count,
    COUNT(*) FILTER (
        WHERE i.is_confirmed = FALSE
    )::BIGINT AS pending_confirmation_count,
    COUNT(*) FILTER (
        WHERE i.occurred_at >= NOW() - INTERVAL '24 hours'
    )::BIGINT AS past_24h_count
FROM incident i;
