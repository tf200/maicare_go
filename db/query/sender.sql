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
    email_adress,
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



