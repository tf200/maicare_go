

-- name: CreateNotification :one
INSERT INTO notifications (
    user_id,
    type,
    data,
    message
)  VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;



-- name: ListNotifications :many
SELECT * FROM notifications
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;



-- name: MarkNotificationAsRead :one
UPDATE notifications
SET is_read = TRUE
WHERE id = $1
RETURNING *;