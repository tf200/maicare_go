

-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    type,
    data

)  VALUES (
    $1,
    $2,
    $3
) RETURNING *;