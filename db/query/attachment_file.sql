-- name: CreateAttachment :one
INSERT INTO attachment_file (
    "uuid",
    "name",
    "file",
    "size",
    tag
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
) RETURNING *;

-- name: GetAttachmentById :one
SELECT * FROM attachment_file
WHERE uuid = $1 LIMIT 1;

-- name: DeleteAttachment :one
DELETE FROM attachment_file
WHERE uuid = $1
RETURNING *;


-- name: SetAttachmentAsUsedorUnused :one
UPDATE attachment_file
SET
    is_used = $2
WHERE
    uuid = $1
RETURNING *;