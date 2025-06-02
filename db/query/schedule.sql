
-- name: CreateSchedule :one
INSERT INTO schedules (
    employee_id,
    location_id,
    color,
    start_datetime,
    end_datetime
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;


-- name: GetMonthlySchedulesByLocation :many
WITH target_month AS (
  SELECT make_date($1, $2, 1) AS start_date
),
dates AS (
  SELECT generate_series(
    (SELECT start_date FROM target_month),
    (SELECT start_date + INTERVAL '1 month - 1 day' FROM target_month),
    INTERVAL '1 day'
  )::date AS day
),
shift_days AS (
  SELECT 
    s.id AS shift_id,
    s.employee_id,
    s.location_id,
    s.start_datetime,
    s.end_datetime,
    s.color,
    d.day,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name
  FROM schedules s
  JOIN dates d 
    ON d.day BETWEEN DATE(s.start_datetime) AND DATE(s.end_datetime)
  JOIN employee_profile e 
    ON s.employee_id = e.id
  WHERE s.location_id = $3
)
SELECT *
FROM shift_days
ORDER BY day, start_datetime;



-- name: GetDailySchedulesByLocation :many
WITH target_day AS (
  SELECT make_date($1, $2, $3) AS day
),
shift_days AS (
  SELECT 
    s.id AS shift_id,
    s.employee_id,
    s.location_id,
    s.color,
    s.start_datetime,
    s.end_datetime,
    d.day,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name
  FROM schedules s
  JOIN target_day d 
    ON d.day BETWEEN DATE(s.start_datetime) AND DATE(s.end_datetime)
  JOIN employee_profile e 
    ON s.employee_id = e.id
  WHERE s.location_id = $4
)
SELECT *
FROM shift_days
ORDER BY start_datetime;


-- name: GetScheduleById :one
SELECT s.*,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name,
    l.name AS location_name
FROM schedules s
JOIN employee_profile e ON s.employee_id = e.id
JOIN location l ON s.location_id = l.id
WHERE s.id = $1
LIMIT 1;

-- name: UpdateSchedule :one
UPDATE schedules
SET
    employee_id = $2,
    location_id = $3,
    start_datetime = $4,
    end_datetime = $5
WHERE id = $1
RETURNING *;

-- name: DeleteSchedule :exec
DELETE FROM schedules
WHERE id = $1;


