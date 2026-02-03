-- name: CreateOrganisation :one
INSERT INTO organisations (
    name,
    street,
    house_number,
    house_number_addition,
    postal_code,
    city,
    phone_number,
    email,
    kvk_number,
    btw_number
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;


-- name: ListOrganisations :many
SELECT o.*,
         COUNT(l.id) AS location_count
FROM organisations o
LEFT JOIN location l ON o.id = l.organisation_id
GROUP BY o.id
ORDER BY o.name;


-- name: ListOrganisationsPaginated :many
SELECT o.*,
         COUNT(l.id) AS location_count,
         COUNT(*) OVER() AS total_count
FROM organisations o
LEFT JOIN location l ON o.id = l.organisation_id
WHERE ($3::text IS NULL OR o.name ILIKE '%' || $3 || '%')
GROUP BY o.id
ORDER BY o.name
LIMIT $1 OFFSET $2;


-- name: GetOrganisation :one
SELECT o.*,
       COUNT(l.id) AS location_count
FROM organisations o
LEFT JOIN location l ON o.id = l.organisation_id
WHERE o.id = $1
GROUP BY o.id;

-- name: GetOrganisationCounts :one
SELECT 
    o.id AS organisation_id,
    o.name AS organisation_name,
    COALESCE(COUNT(DISTINCT l.id), 0)::BIGINT AS location_count,
    COALESCE(COUNT(DISTINCT c.id), 0)::BIGINT AS client_count,
    COALESCE(COUNT(DISTINCT e.id), 0)::BIGINT AS employee_count
FROM 
    organisations o
    LEFT JOIN location l ON o.id = l.organisation_id
    LEFT JOIN client_details c ON l.id = c.location_id
    LEFT JOIN employee_profile e ON l.id = e.location_id
WHERE 
    o.id = $1
GROUP BY 
    o.id, o.name;


-- name: UpdateOrganisation :one
UPDATE organisations
SET
    name = COALESCE(sqlc.narg('name'), name),
    street = COALESCE(sqlc.narg('street'), street),
    house_number = COALESCE(sqlc.narg('house_number'), house_number),
    house_number_addition = COALESCE(sqlc.narg('house_number_addition'), house_number_addition),
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
DELETE FROM organisations
WHERE id = $1
RETURNING *;



-- name: CreateLocation :one
INSERT INTO location (
    organisation_id,
    name,
    street,
    house_number,
    house_number_addition,
    postal_code,
    city,
    capacity
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;



-- name: ListLocations :many
SELECT l.*,
    COUNT(c.id) AS client_count
FROM location l
LEFT JOIN client_details c ON l.id = c.location_id
WHERE organisation_id = $1
GROUP BY l.id;

-- name: ListLocationsPaginated :many
SELECT l.*,
    COUNT(c.id) AS client_count,
    COUNT(*) OVER() AS total_count
FROM location l
LEFT JOIN client_details c ON l.id = c.location_id
WHERE l.organisation_id = $1
  AND ($4::text IS NULL OR l.name ILIKE '%' || $4 || '%')
GROUP BY l.id
ORDER BY l.name
LIMIT $2 OFFSET $3;

-- name: GetLocation :one
SELECT * FROM location
WHERE id = $1;

-- name: UpdateLocation :one
UPDATE location
SET
    name = COALESCE(sqlc.narg('name'), name),
    street = COALESCE(sqlc.narg('street'), street),
    house_number = COALESCE(sqlc.narg('house_number'), house_number),
    house_number_addition = COALESCE(sqlc.narg('house_number_addition'), house_number_addition),
    postal_code = COALESCE(sqlc.narg('postal_code'), postal_code),
    city = COALESCE(sqlc.narg('city'), city),
    capacity = COALESCE(sqlc.narg('capacity'), capacity)
WHERE id = $1
RETURNING *;


-- name: DeleteLocation :one
DELETE FROM location
WHERE id = $1
RETURNING *;



-- name: ListAllLocations :many
SELECT l.*,
       COUNT(c.id) AS client_count
FROM location l
LEFT JOIN client_details c ON l.id = c.location_id
GROUP BY l.id
ORDER BY l.name;

-- name: ListAllLocationsPaginated :many
SELECT l.*,
       COUNT(c.id) AS client_count,
       COUNT(*) OVER() AS total_count
FROM location l
LEFT JOIN client_details c ON l.id = c.location_id
WHERE ($3::text IS NULL OR l.name ILIKE '%' || $3 || '%')
GROUP BY l.id
ORDER BY l.name
LIMIT $1 OFFSET $2;
