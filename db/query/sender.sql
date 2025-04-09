-- name: CreateSender :one
INSERT INTO sender (
    types,
    name,
    address,
    postal_code,
    place,
    land,
    kvknumber,
    btwnumber,
    phone_number,
    client_number,
    email_address,
    contacts
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
) RETURNING *;



-- name: GetSenderById :one
SELECT * FROM sender
WHERE id = $1 LIMIT 1;


-- name: ListSenders :many
SELECT * FROM sender
WHERE 
   (CASE WHEN sqlc.narg('include_archived')::boolean THEN true ELSE NOT is_archived END)
   AND (sqlc.narg('search')::TEXT IS NULL OR name ILIKE '%' || sqlc.narg('search') || '%')
ORDER BY name
LIMIT $1 OFFSET $2;


-- name: CountSenders :one
SELECT COUNT(*) 
FROM sender
WHERE (
  CASE 
    WHEN sqlc.narg('include_archived')::boolean IS NULL THEN NOT is_archived
    WHEN sqlc.narg('include_archived')::boolean THEN true 
    ELSE NOT is_archived
  END
);



-- name: UpdateSender :one
UPDATE sender
SET 
    name = COALESCE(sqlc.narg('name'), name),
    address = COALESCE(sqlc.narg('address'), address),
    postal_code = COALESCE(sqlc.narg('postal_code'), postal_code),
    place = COALESCE(sqlc.narg('place'), place),
    land = COALESCE(sqlc.narg('land'), land),
    kvknumber = COALESCE(sqlc.narg('kvknumber'), kvknumber),
    btwnumber = COALESCE(sqlc.narg('btwnumber'), btwnumber),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    client_number = COALESCE(sqlc.narg('client_number'), client_number),
    email_address = COALESCE(sqlc.narg('email_address'), email_address),
    contacts = COALESCE(sqlc.narg('contacts')::JSONB, contacts),
    updated_at = NOW(),
    is_archived = COALESCE(sqlc.narg('is_archived'), is_archived),
    types = COALESCE(sqlc.narg('types'), types)
WHERE 
    id = sqlc.arg('id')
RETURNING *;


