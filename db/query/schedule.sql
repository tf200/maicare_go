
-- name: CreateSchedule :one
INSERT INTO schedules (
    employee_id,
    location_id,
    location_shift_id,
    color,
    is_custom,
    start_datetime,
    end_datetime
) VALUES (
    $1, $2, $3, $4, $5 , $6, $7
)
RETURNING *;


-- name: GetMonthlySchedulesByLocation :many
WITH target_month AS (
  SELECT make_date($1, $2, 1) AS start_date
)
SELECT 
  s.id AS shift_id,
  s.employee_id,
  s.location_id,
  s.start_datetime,
  s.end_datetime,
  s.color,
  S.is_custom,
  ls.shift_name,
  ls.id AS location_shift_id,
  DATE(s.start_datetime) AS day,
  e.first_name AS employee_first_name,
  e.last_name AS employee_last_name
FROM schedules s
LEFT JOIN location_shift ls 
  ON s.location_shift_id = ls.id
JOIN employee_profile e 
  ON s.employee_id = e.id
WHERE s.location_id = $3
  AND DATE(s.start_datetime) >= (SELECT start_date FROM target_month)
  AND DATE(s.start_datetime) < (SELECT start_date + INTERVAL '1 month' FROM target_month)
ORDER BY DATE(s.start_datetime), s.start_datetime;



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
    s.is_custom,
    ls.id AS location_shift_id,
    ls.shift_name,
    d.day,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name
  FROM schedules s
  LEFT JOIN location_shift ls 
    ON s.location_shift_id = ls.id
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
    l.name AS location_name,
    ls.shift_name AS location_shift_name,
    ls.id AS location_shift_id
FROM schedules s
LEFT JOIN location_shift ls ON s.location_shift_id = ls.id
JOIN employee_profile e ON s.employee_id = e.id
JOIN location l ON s.location_id = l.id
WHERE s.id = $1
LIMIT 1;

-- name: UpdateSchedule :one
UPDATE schedules
SET
    employee_id = $2,
    location_id = $3,
    location_shift_id = $4,
    color = $5,
    start_datetime = $6,
    end_datetime = $7
WHERE id = $1
RETURNING *;

-- name: DeleteSchedule :exec
DELETE FROM schedules
WHERE id = $1;



-- name: GetEmployeeSchedules :many
SELECT 
    s.id,
    s.start_datetime,
    s.end_datetime,
    s.location_id,
    s.color,
    l.name as location_name,
    'shift'::text as type
FROM schedules s
JOIN location l ON s.location_id = l.id
WHERE s.employee_id = @employee_id
    AND s.start_datetime >= @period_start
    AND s.start_datetime < @period_end
ORDER BY s.start_datetime;