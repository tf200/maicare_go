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


