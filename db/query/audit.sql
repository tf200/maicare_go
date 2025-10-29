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