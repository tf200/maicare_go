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
