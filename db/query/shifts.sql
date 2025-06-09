-- name: CreateShift :one
INSERT INTO location_shift (
    location_id,
    shift_name,
    start_time,
    end_time
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateShift :one
UPDATE location_shift
SET
    shift_name = COALESCE($2, shift_name),
    start_time = COALESCE($3, start_time),
    end_time = COALESCE($4, end_time),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetShiftByID :one
SELECT * FROM location_shift
WHERE id = $1
LIMIT 1;


-- name: DeleteShift :exec
DELETE FROM location_shift
WHERE id = $1;


-- name: GetShiftsByLocationID :many
SELECT * FROM location_shift
WHERE location_id = $1;




