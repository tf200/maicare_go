-- name: CreateAttachment :one
INSERT INTO attachment_file (
    "name",
    "file",
    "size",
    tag
) VALUES (
    $1,
    $2,
    $3,
    $4
) RETURNING *;



