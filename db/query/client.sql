-- name: CreateClientDetails :one
INSERT INTO client_details (
    first_name,
    last_name,
    date_of_birth,
    "identity",
    "status",
    bsn,
    source,
    birthplace,
    email,
    phone_number,
    organisation,
    departement,
    gender,
    filenumber,
    profile_picture,
    infix,
    sender_id,
    location_id,
    identity_attachment_ids,
    departure_reason,
    departure_report,
    addresses,
    legal_measure
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 
    $17, $18, $19, $20, $21, $22, $23
) RETURNING *;


-- name: ListClientDetails :many
SELECT 
    *, 
    COUNT(*) OVER() AS total_count
FROM client_details
WHERE
    (status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
    (location_id = sqlc.narg('location_id') OR sqlc.narg('location_id') IS NULL) AND
    (sqlc.narg('search')::TEXT IS NULL OR 
        first_name ILIKE '%' || sqlc.narg('search') || '%' OR
        last_name ILIKE '%' || sqlc.narg('search') || '%' OR
        filenumber ILIKE '%' || sqlc.narg('search') || '%' OR
        email ILIKE '%' || sqlc.narg('search') || '%' OR
        phone_number ILIKE '%' || sqlc.narg('search') || '%')
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');


-- name: GetClientDetails :one
SELECT * FROM client_details
WHERE id = $1 LIMIT 1;




-- name: CreateClientAllergy :one
INSERT INTO client_allergy (
    client_id,
    allergy_type_id,
    severity,
    reaction,
    notes,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;


-- name: GetClientAllergy :one
SELECT * FROM client_allergy
WHERE id = $1 LIMIT 1;

-- name: UpdateClientAllergy :one
UPDATE client_allergy
SET
    allergy_type_id = COALESCE(sqlc.narg('allergy_type_id'), allergy_type_id),
    severity = COALESCE(sqlc.narg('severity'), severity),
    reaction = COALESCE(sqlc.narg('reaction'), reaction),
    notes = COALESCE(sqlc.narg('notes'), notes)
WHERE id = $1
RETURNING *;

-- name: ListClientAllergies :many
SELECT 
    al.*,
    ty.name AS allergy_type,
    (SELECT COUNT(*) FROM client_allergy WHERE al.client_id = $1) AS total_allergies
FROM client_allergy al
JOIN allergy_type ty ON al.allergy_type_id = ty.id
WHERE al.client_id = $1
LIMIT $2 OFFSET $3;
