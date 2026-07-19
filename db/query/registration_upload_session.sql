-- name: CreateRegistrationUploadSession :one
INSERT INTO registration_upload_sessions (token_hash, expires_at)
VALUES ($1, $2)
RETURNING *;

-- name: GetActiveRegistrationUploadSession :one
SELECT * FROM registration_upload_sessions
WHERE token_hash = $1
  AND expires_at > CURRENT_TIMESTAMP
  AND submitted_at IS NULL
LIMIT 1;

-- name: AddRegistrationUploadAttachment :exec
UPDATE registration_upload_sessions
SET attachment_ids = array_append(attachment_ids, $2::uuid)
WHERE id = $1
  AND expires_at > CURRENT_TIMESTAMP
  AND submitted_at IS NULL;

-- name: RegistrationUploadSessionHasAttachments :one
SELECT EXISTS (
    SELECT 1 FROM registration_upload_sessions
    WHERE id = $1
      AND expires_at > CURRENT_TIMESTAMP
      AND submitted_at IS NULL
      AND attachment_ids @> $2::uuid[]
) AS valid;

-- name: ConsumeRegistrationUploadSession :exec
UPDATE registration_upload_sessions
SET submitted_at = CURRENT_TIMESTAMP
WHERE id = $1
  AND submitted_at IS NULL;
