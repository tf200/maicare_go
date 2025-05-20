
-- name: CreateSchedule :one
INSERT INTO schedules (
    employee_id,
    location_id,
    start_datetime,
    end_datetime
) VALUES (
    $1, $2, $3, $4
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
    d.day,
    e.first_name AS employee_first_name
  FROM schedules s
  JOIN dates d 
    ON d.day BETWEEN DATE(s.start_datetime) AND DATE(s.end_datetime)
  JOIN employee_profile e 
    ON s.employee_id = e.employee_id
  WHERE s.location_id = $3
)
SELECT *
FROM shift_days
ORDER BY day, start_datetime;