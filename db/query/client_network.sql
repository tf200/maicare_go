-- name: CreateEmemrgencyContact :one
INSERT INTO client_emergency_contact (
    client_id,
    first_name,
    last_name,
    email,
    phone_number,
    address,
    relationship,
    relation_status,
    medical_reports,
    incidents_reports,
    goals_reports
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: ListEmergencyContacts :many
SELECT 
    ec.*,
    COUNT(*) OVER() as total_count
FROM client_emergency_contact ec
WHERE client_id = @client_id
    AND (
        LOWER(first_name) LIKE LOWER(CONCAT('%', @search::text, '%')) OR
        LOWER(last_name) LIKE LOWER(CONCAT('%', @search::text, '%'))
    )
ORDER BY id
LIMIT $1 OFFSET $2;


-- name: GetEmergencyContact :one
SELECT * FROM client_emergency_contact
WHERE id = $1 LIMIT 1;

-- name: UpdateEmergencyContact :one
UPDATE client_emergency_contact
SET
    first_name = COALESCE(sqlc.narg('first_name'), first_name),
    last_name = COALESCE(sqlc.narg('last_name'), last_name),
    email = COALESCE(sqlc.narg('email'), email),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    address = COALESCE(sqlc.narg('address'), address),
    relationship = COALESCE(sqlc.narg('relationship'), relationship),
    relation_status = COALESCE(sqlc.narg('relation_status'), relation_status),
    medical_reports = COALESCE(sqlc.narg('medical_reports'), medical_reports),
    incidents_reports = COALESCE(sqlc.narg('incidents_reports'), incidents_reports),
    goals_reports = COALESCE(sqlc.narg('goals_reports'), goals_reports)
WHERE id = $1
RETURNING *;


-- name: DeleteEmergencyContact :one
DELETE FROM client_emergency_contact
WHERE id = $1
RETURNING *;