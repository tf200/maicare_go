-- Audit table queries

-- name: LockAuditHashChain :exec
SELECT pg_advisory_xact_lock(7513);

-- name: CreateAuditRecord :exec
INSERT INTO audit (
    event_id,
    event_group_id,
    occurred_at,
    event_type,
    action,
    result,
    actor_user_id,
    actor_employee_id,
    actor_roles,
    subject_type,
    subject_id,
    client_id,
    access_rule,
    access_reason,
    session_id,
    request_id,
    ip,
    user_agent,
    route,
    method,
    details,
    hash_prev,
    hash_self
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9,
    $10, $11, $12,
    $13, $14,
    $15, $16, $17, $18, $19, $20,
    $21, $22, $23
);

-- name: BulkCreateAuditRecords :copyfrom
INSERT INTO audit (
    event_id,
    event_group_id,
    occurred_at,
    event_type,
    action,
    result,
    actor_user_id,
    actor_employee_id,
    actor_roles,
    subject_type,
    subject_id,
    client_id,
    access_rule,
    access_reason,
    session_id,
    request_id,
    ip,
    user_agent,
    route,
    method,
    details,
    hash_prev,
    hash_self
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9,
    $10, $11, $12,
    $13, $14,
    $15, $16, $17, $18, $19, $20,
    $21, $22, $23
);

-- name: GetLatestAuditHash :one
SELECT hash_self
FROM audit
ORDER BY append_seq DESC
LIMIT 1;

-- name: ListAuditRecords :many
SELECT *
FROM audit
WHERE
    (sqlc.narg('subject_id')::TEXT IS NULL OR subject_id = sqlc.narg('subject_id')::TEXT)
    AND (sqlc.narg('actor_user_id')::UUID IS NULL OR actor_user_id = sqlc.narg('actor_user_id')::UUID)
    AND (sqlc.narg('start_time')::TIMESTAMPTZ IS NULL OR occurred_at >= sqlc.narg('start_time')::TIMESTAMPTZ)
    AND (sqlc.narg('end_time')::TIMESTAMPTZ IS NULL OR occurred_at <= sqlc.narg('end_time')::TIMESTAMPTZ)
    AND (sqlc.narg('client_id')::UUID IS NULL OR client_id = sqlc.narg('client_id')::UUID)
ORDER BY occurred_at DESC
LIMIT $1 OFFSET $2;

-- name: ListAuditRecordsByClientID :many
SELECT *
FROM audit
WHERE client_id = $1
ORDER BY occurred_at DESC
LIMIT $2 OFFSET $3;
