-- name: CreateLeaveRequest :one
INSERT INTO leave_requests (
    employee_id,
    created_by_employee_id,
    leave_type,
    start_date,
    end_date,
    reason,
    requested_at
) VALUES (
    sqlc.arg(employee_id),
    sqlc.arg(created_by_employee_id),
    sqlc.arg(leave_type),
    sqlc.arg(start_date),
    sqlc.arg(end_date),
    sqlc.narg(reason),
    NOW()
)
RETURNING *;
