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




-- name: SetClientProfilePicture :one
UPDATE client_details
SET profile_picture = $2
WHERE id = $1
RETURNING *;


-- name: CreateClientDocument :one
INSERT INTO client_documents (
    client_id,
    attachment_uuid,
    label
) VALUES (
    $1, $2, $3
) RETURNING *;


-- name: ListClientDocuments :many
SELECT 
    cd.*,
    a.*,
    COUNT(*) OVER() AS total_count
FROM client_documents cd
JOIN attachment_file a ON cd.attachment_uuid = a.uuid
WHERE client_id = $1
LIMIT $2 OFFSET $3;


-- name: DeleteClientDocument :one
DELETE FROM client_documents
WHERE attachment_uuid = $1
RETURNING *;


-- name: GetMissingClientDocuments :many
WITH all_labels AS (
    SELECT unnest(ARRAY[
        'registration_form', 'intake_form', 'consent_form',
        'risk_assessment', 'self_reliance_matrix', 'force_inventory',
        'care_plan', 'signaling_plan', 'cooperation_agreement'
    ]) AS label
),
client_labels AS (
    SELECT label
    FROM client_documents
    WHERE client_id = $1
)
SELECT al.label::text AS missing_label
FROM all_labels al
LEFT JOIN client_labels cl ON al.label = cl.label
WHERE cl.label IS NULL;