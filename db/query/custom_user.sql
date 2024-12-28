-- name: CreateUser :one
INSERT INTO custom_user (
    password,
    username,
    first_name,
    last_name,
    email,
    is_superuser,
    is_staff,
    is_active,
    profile_picture,
    phone_number
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM custom_user
WHERE id = $1 LIMIT 1;



-- name: GetUserByUsername :one
SELECT * FROM custom_user
WHERE username= $1 LIMIT 1;



-- name: GetUserByEmail :one
SELECT * FROM custom_user
WHERE email= $1 LIMIT 1;