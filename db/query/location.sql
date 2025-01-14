-- name: CreateLocation :one
INSERT INTO location (
    name,
    address,
    capacity
) VALUES (
    $1, $2, $3
) RETURNING *;



-- name: ListLocations :many
SELECT * FROM location;



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