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



-- name: SetAttachmentAsUsed :one
UPDATE attachment_file
SET
    is_used = true
WHERE
    uuid = $1
RETURNING *;