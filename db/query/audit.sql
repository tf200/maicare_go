-- Audit table queries


-- name: CreateAuditRecord :exec
INSERT INTO audit (
    event_id,
    event_type,
    occured_at,
    actor_role,
    actor_id,
    subject_type,
    subject_id,
    access_reason,
    action,
    result,
    module,
    tenant_id,
    details,
    ip,
    user_agent,
    hash_prev,
    hash_self
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17
) ; 


-- name: GetLatestAuditHash :one
SELECT hash_self
FROM audit
ORDER BY occured_at DESC
LIMIT 1;



-- name: ListAuditRecords :many
SELECT
    event_id,
    event_type,
    occured_at,
    actor_role,
    actor_id,
    subject_type,
    subject_id,
    access_reason,
    action,
    result,
    module,
    tenant_id,
    details,
    ip,
    user_agent,
    hash_prev,
    hash_self
FROM audit
WHERE
    (sqlc.narg('subject_id')::UUID IS NULL OR subject_id = sqlc.narg('subject_id')::UUID)
    AND (sqlc.narg('actor_id')::UUID IS NULL OR actor_id = sqlc.narg('actor_id')::UUID)
    AND (sqlc.narg('start_time')::TIMESTAMPTZ IS NULL OR occured_at >= sqlc.narg('start_time')::TIMESTAMPTZ)
    AND (sqlc.narg('end_time')::TIMESTAMPTZ IS NULL OR occured_at <= sqlc.narg('end_time')::TIMESTAMPTZ)
    AND (sqlc.narg('tenant_id')::TEXT IS NULL OR tenant_id = sqlc.narg('tenant_id')::TEXT)
ORDER BY occured_at DESC
LIMIT $1 OFFSET $2;