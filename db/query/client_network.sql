
-- name: GetClientSender :one
SELECT s.* FROM sender s
JOIN client_details cd ON s.id = cd.sender_id
WHERE cd.id = $1
LIMIT 1;

-- name: AssignSender :one
UPDATE client_details
SET sender_id = $1
WHERE id = $2
RETURNING *;

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
WITH inserted_assignment AS (
    -- Perform the insert and return the inserted row
    INSERT INTO assigned_employee (
        client_id,
        employee_id,
        start_date,
        role
    ) VALUES (
        $1, $2, $3, $4
    )
    RETURNING * -- Return all columns from the inserted assigned_employee row
)
-- Select the columns from the inserted row AND join to get the user_id
SELECT
    ia.*,  -- Select all columns from the inserted_assignment CTE
    ep.user_id, -- Select the user_id from the employee_profile table
    cl.first_name AS client_first_name,
    cl.last_name AS client_last_name,
    l.name AS client_location_name
FROM
    inserted_assignment ia
JOIN
    employee_profile ep ON ia.employee_id = ep.id -- Join based on employee_id
JOIN
    client_details cl ON ia.client_id = cl.id
LEFT JOIN
    location l ON cl.location_id = l.id; -- Join to get the client location name




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


-- name: DeleteAssignedEmployee :one
DELETE FROM assigned_employee
WHERE id = $1
RETURNING *;



-- name: GetClientRelatedEmails :many
WITH employee_emails AS (
    SELECT ed.email AS employee_email
    FROM assigned_employee ae
    JOIN employee_profile ed ON ae.employee_id = ed.id
    WHERE ae.client_id = $1
),
emergency_contact_emails AS (
    SELECT cec.email AS contact_email
    FROM client_emergency_contact cec
    WHERE cec.client_id = $1
)
SELECT employee_email FROM employee_emails
UNION
SELECT contact_email FROM emergency_contact_emails;
