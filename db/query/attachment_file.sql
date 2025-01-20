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

-- name: GetAttachmentById :one
SELECT * FROM attachment_file
WHERE uuid = $1 LIMIT 1;

-- name: DeleteAttachment :one
DELETE FROM attachment_file
WHERE uuid = $1
RETURNING *;


-- name: SetAttachmentAsUsed :one
UPDATE attachment_file
SET
    is_used = true
WHERE
    uuid = $1
RETURNING *;