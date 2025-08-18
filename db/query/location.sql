-- name: CreateOrganisation :one
INSERT INTO organisation (
    name,
    address,
    postal_code,
    city,
    phone_number,
    email,
    kvk_number,
    btw_number
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;


-- name: ListOrganisations :many
SELECT o.*,
         COUNT(l.id) AS location_count
FROM organisation o
LEFT JOIN location l ON o.id = l.organisation_id
GROUP BY o.id
ORDER BY o.name;


-- name: GetOrganisation :one
SELECT o.*,
       COUNT(l.id) AS location_count
FROM organisation o
LEFT JOIN location l ON o.id = l.organisation_id
WHERE o.id = $1;


-- name: UpdateOrganisation :one
UPDATE organisation
SET
    name = COALESCE(sqlc.narg('name'), name),
    address = COALESCE(sqlc.narg('address'), address),
    postal_code = COALESCE(sqlc.narg('postal_code'), postal_code),
    city = COALESCE(sqlc.narg('city'), city),
    phone_number = COALESCE(sqlc.narg('phone_number'), phone_number),
    email = COALESCE(sqlc.narg('email'), email),
    kvk_number = COALESCE(sqlc.narg('kvk_number'), kvk_number),
    btw_number = COALESCE(sqlc.narg('btw_number'), btw_number),
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: DeleteOrganisation :one
DELETE FROM organisation
WHERE id = $1
RETURNING *;



-- name: CreateLocation :one
INSERT INTO location (
    organisation_id,
    name,
    address,
    capacity
) VALUES (
    $1, $2, $3, $4
) RETURNING *;



-- name: ListLocations :many
SELECT * FROM location 
WHERE organisation_id = $1;

-- name: GetLocation :one
SELECT * FROM location
WHERE id = $1;

-- name: UpdateLocation :one
UPDATE location
SET
    name = COALESCE(sqlc.narg('name'), name),
    address = COALESCE(sqlc.narg('address'), address),
    capacity = COALESCE(sqlc.narg('capacity'), capacity)
WHERE id = $1
RETURNING *;


-- name: DeleteLocation :one
DELETE FROM location
WHERE id = $1
RETURNING *;



-- name: ListAllLocations :many
SELECT * FROM location
ORDER BY name;