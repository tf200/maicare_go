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



-- name: AssignEmployee :one
INSERT INTO assigned_employee (
    client_id,
    employee_id,
    start_date,
    role
) VALUES (
    $1, $2, $3, $4
) RETURNING *;


-- name: ListAssignedEmployees :many
SELECT 
    ae.*,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name,
    COUNT(*) OVER() as total_count
FROM assigned_employee ae
JOIN employee_profile e ON ae.employee_id = e.id
WHERE ae.client_id = $1
ORDER BY ae.start_date DESC
LIMIT $2 OFFSET $3;


-- name: GetAssignedEmployee :one
SELECT 
    ae.*,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name
FROM assigned_employee ae
JOIN employee_profile e ON ae.employee_id = e.id
WHERE ae.id = $1 LIMIT 1;


-- name: UpdateAssignedEmployee :one
UPDATE assigned_employee
SET
    employee_id = COALESCE(sqlc.narg('employee_id'), employee_id),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    role = COALESCE(sqlc.narg('role'), role)
WHERE id = $1
RETURNING *;

