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

-- name: GetActorAttachmentById :one
SELECT * FROM attachment_file
WHERE uuid = $1
  AND uploaded_by_user_id IS NOT DISTINCT FROM public.get_current_user_id()
  AND public.can_access_actor_attachment(uuid)
LIMIT 1;

-- name: DeleteActorAttachment :one
DELETE FROM attachment_file
WHERE uuid = $1
  AND uploaded_by_user_id IS NOT DISTINCT FROM public.get_current_user_id()
  AND NOT public.attachment_file_is_referenced(uuid)
RETURNING *;


-- name: SetAttachmentAsUsedorUnused :one
UPDATE attachment_file
SET
    is_used = $2
WHERE
    uuid = $1
RETURNING *;

-- name: SetActorAttachmentAsUsedOrUnused :one
UPDATE attachment_file
SET is_used = $2
WHERE uuid = $1
  AND uploaded_by_user_id IS NOT DISTINCT FROM public.get_current_user_id()
  AND public.can_access_actor_attachment(uuid)
RETURNING *;

-- name: SetActorAttachmentsAsUsedByUUIDs :many
UPDATE attachment_file
SET is_used = $2
WHERE uuid = ANY($1::uuid[])
  AND uploaded_by_user_id IS NOT DISTINCT FROM public.get_current_user_id()
  AND public.can_access_actor_attachment(uuid)
RETURNING *;

-- name: GetAttachmentsByUUIDs :many
SELECT * FROM attachment_file
WHERE uuid = ANY($1::uuid[]);

-- name: GetActorAttachmentsByUUIDs :many
SELECT * FROM attachment_file
WHERE uuid = ANY($1::uuid[])
  AND uploaded_by_user_id IS NOT DISTINCT FROM public.get_current_user_id()
  AND public.can_access_actor_attachment(uuid);
