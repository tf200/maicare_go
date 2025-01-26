-- name: Createallergy :one
INSERT INTO allergy_type (
    name
) VALUES (
    $1
) RETURNING *;