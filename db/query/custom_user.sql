-- name: CreateUser :one
INSERT INTO custom_user (
    password,
    email,
    is_active,
    profile_picture
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM custom_user
WHERE id = $1 LIMIT 1;


-- name: GetUserByEmail :one
SELECT * FROM custom_user
WHERE email= $1 LIMIT 1;