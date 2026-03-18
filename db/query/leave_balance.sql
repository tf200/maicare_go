-- name: EnsureLeaveBalanceForYear :exec
INSERT INTO leave_balances (
    employee_id,
    year
) VALUES (
    sqlc.arg(employee_id),
    sqlc.arg('year')
)
ON CONFLICT (employee_id, year) DO NOTHING;

-- name: LockLeaveBalanceByEmployeeYear :one
SELECT *
FROM leave_balances
WHERE employee_id = sqlc.arg(employee_id)
  AND year = sqlc.arg('year')
FOR UPDATE;

-- name: ApplyLeaveBalanceDeduction :one
UPDATE leave_balances
SET
    extra_used_days = extra_used_days + sqlc.arg(extra_days),
    legal_used_days = legal_used_days + sqlc.arg(legal_days),
    updated_at = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;
