-- name: CreateAppointment :one
INSERT INTO appointments (
    creator_employee_id, 
    start_time,           
    end_time,             
    location,             
    description,          
    status,               
    recurrence_type,      
    recurrence_interval,  
    recurrence_end_date   
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;


-- name: AddAppointmentParticipant :exec
INSERT INTO appointment_participants (
    appointment_id,
    employee_id     
) VALUES (
    $1, $2
);

-- name: AddAppointmentClient :exec
INSERT INTO appointment_clients (
    appointment_id,
    client_id       
) VALUES (
    $1, $2
);



-- name: ListAppointmentsForEmployeeInRange :many
-- Parameters:
-- @employee_id: BIGINT
-- @start_date: TIMESTAMP - The beginning of the desired time range (inclusive)
-- @end_date: TIMESTAMP - The end of the desired time range (inclusive)
WITH EmployeeAppointments AS (
    -- CTE to get all relevant appointment IDs (both created and participated)
    SELECT id AS appointment_id
    FROM appointments a
    WHERE a.creator_employee_id = @employee_id -- Use named parameter

    UNION -- Use UNION to automatically handle duplicates

    SELECT appointment_id
    FROM appointment_participants
    WHERE employee_id = @employee_id -- Use named parameter
),

RecurringOccurrences AS (
    -- CTE to generate potential future occurrences for recurring appointments
    SELECT
        a.id AS original_appointment_id,
        a.creator_employee_id,
        -- Calculate the start time of the specific occurrence
        ts.occurrence_start_time::timestamp AS start_time,
        -- Calculate the end time of the specific occurrence by adding the original duration
        (ts.occurrence_start_time + (a.end_time - a.start_time))::timestamp AS end_time,
        a.location,
        a.description,
        a.status,
        a.recurrence_type,
        a.recurrence_interval,
        a.recurrence_end_date,
        a.confirmed_by_employee_id,
        a.confirmed_at,
        a.created_at,
        a.updated_at,
        TRUE AS is_recurring_occurrence -- Add a flag to indicate this is a generated occurrence
    FROM
        appointments a
    INNER JOIN EmployeeAppointments ea ON a.id = ea.appointment_id -- Only consider appointments involving the employee
    -- Use generate_series to create timestamps based on recurrence rules
    CROSS JOIN LATERAL generate_series(
        -- Start generating from the appointment's original start time
        a.start_time,
        -- Stop generating at the recurrence end date OR the query's end date, whichever is EARLIER
        LEAST(COALESCE(a.recurrence_end_date::timestamp, 'infinity'::timestamp), @end_date::timestamp), -- Use named parameter
        -- Calculate the interval step based on recurrence type and interval
        CASE a.recurrence_type
            WHEN 'DAILY' THEN (COALESCE(a.recurrence_interval, 1) || ' day')::interval
            WHEN 'WEEKLY' THEN (COALESCE(a.recurrence_interval, 1) || ' week')::interval
            WHEN 'MONTHLY' THEN (COALESCE(a.recurrence_interval, 1) || ' month')::interval
            -- Default to a very large interval if type is NONE or unexpected
            ELSE '1000 years'::interval
        END
    ) AS ts(occurrence_start_time)
    WHERE
        a.recurrence_type != 'NONE' -- Only process recurring appointments
        AND a.recurrence_type IS NOT NULL -- Safety check
        -- Optimization: Ensure the base appointment's start is before the query window ends
        AND a.start_time <= @end_date::timestamp -- Use named parameter
        -- Optimization: Ensure the recurrence doesn't end before the query window starts
        AND COALESCE(a.recurrence_end_date::timestamp, 'infinity'::timestamp) >= @start_date::timestamp -- Use named parameter
)

-- Final SELECT combining non-recurring and calculated recurring appointments
SELECT
    a.id,
    a.creator_employee_id,
    a.start_time,
    a.end_time,
    a.location,
    a.description,
    a.status,
    a.recurrence_type,
    a.recurrence_interval,
    a.recurrence_end_date,
    a.confirmed_by_employee_id,
    a.confirmed_at,
    a.created_at,
    a.updated_at,
    FALSE AS is_recurring_occurrence -- Flag for non-recurring
FROM
    appointments a
INNER JOIN EmployeeAppointments ea ON a.id = ea.appointment_id
WHERE
    a.recurrence_type = 'NONE' -- Select only non-recurring appointments
    -- Standard overlap check: (StartA <= EndB) AND (EndA >= StartB)
    AND (a.start_time <= @end_date::timestamp) -- Use named parameter
    AND (a.end_time >= @start_date::timestamp) -- Use named parameter

UNION ALL -- Combine with recurring occurrences, keeping all rows

SELECT
    ro.original_appointment_id AS id, -- Use the original ID for consistency
    ro.creator_employee_id,
    ro.start_time, -- Calculated start time
    ro.end_time,   -- Calculated end time
    ro.location,
    ro.description,
    ro.status,
    ro.recurrence_type,
    ro.recurrence_interval,
    ro.recurrence_end_date,
    ro.confirmed_by_employee_id,
    ro.confirmed_at,
    ro.created_at,
    ro.updated_at,
    ro.is_recurring_occurrence -- Flag indicating it's a calculated occurrence
FROM
    RecurringOccurrences ro
WHERE
    -- Filter the generated occurrences to only those that OVERLAP the requested time frame
    -- Standard overlap check: (StartA <= EndB) AND (EndA >= StartB)
    (ro.start_time <= @end_date::timestamp)   -- Use named parameter
    AND (ro.end_time >= @start_date::timestamp) -- <<<<< CORRECTED to use @start_date

-- Optional: Order the final results for consistent output
ORDER BY start_time;



-- name: ListAppointmentsForClientInRange :many
-- Parameters:
-- @client_id: BIGINT       - The ID of the client to find appointments for.
-- @start_date: TIMESTAMP  - The beginning of the desired time range (inclusive).
-- @end_date: TIMESTAMP    - The end of the desired time range (inclusive).
WITH ClientAppointments AS (
    -- CTE to get all relevant appointment IDs for the specified client
    SELECT appointment_id
    FROM appointment_clients
    WHERE client_id = @client_id -- Use named parameter

),

RecurringOccurrences AS (
    -- CTE to generate potential future occurrences for recurring appointments
    SELECT
        a.id AS original_appointment_id,
        a.creator_employee_id,
        -- Calculate the start time of the specific occurrence
        ts.occurrence_start_time::timestamp AS start_time,
        -- Calculate the end time of the specific occurrence by adding the original duration
        (ts.occurrence_start_time + (a.end_time - a.start_time))::timestamp AS end_time,
        a.location,
        a.description,
        a.status,
        a.recurrence_type,       -- Uses the recurrence_type column from the corrected table definition
        a.recurrence_interval,
        a.recurrence_end_date,
        a.confirmed_by_employee_id,
        a.confirmed_at,
        a.created_at,
        a.updated_at,
        TRUE AS is_recurring_occurrence -- Add a flag to indicate this is a generated occurrence
    FROM
        appointments a
    INNER JOIN ClientAppointments ca ON a.id = ca.appointment_id -- Only consider appointments involving the client
    -- Use generate_series to create timestamps based on recurrence rules
    CROSS JOIN LATERAL generate_series(
        -- Start generating from the appointment's original start time
        a.start_time,
        -- Stop generating at the recurrence end date OR the query's end date, whichever is EARLIER
        LEAST(COALESCE(a.recurrence_end_date::timestamp, 'infinity'::timestamp), @end_date::timestamp), -- Use named parameter
        -- Calculate the interval step based on recurrence type and interval
        CASE a.recurrence_type -- Uses the recurrence_type column
            WHEN 'DAILY' THEN (COALESCE(a.recurrence_interval, 1) || ' day')::interval
            WHEN 'WEEKLY' THEN (COALESCE(a.recurrence_interval, 1) || ' week')::interval
            WHEN 'MONTHLY' THEN (COALESCE(a.recurrence_interval, 1) || ' month')::interval
            -- Default to a very large interval if type is NONE or unexpected
            ELSE '1000 years'::interval
        END
    ) AS ts(occurrence_start_time)
    WHERE
        a.recurrence_type != 'NONE' -- Only process recurring appointments
        AND a.recurrence_type IS NOT NULL -- Safety check (though DEFAULT 'NONE' helps)
        -- Optimization: Ensure the base appointment's start is before the query window ends
        AND a.start_time <= @end_date::timestamp -- Use named parameter
        -- Optimization: Ensure the recurrence doesn't end before the query window starts
        AND COALESCE(a.recurrence_end_date::timestamp, 'infinity'::timestamp) >= @start_date::timestamp -- Use named parameter
)

-- Final SELECT combining non-recurring and calculated recurring appointments
SELECT
    a.id,
    a.creator_employee_id,
    a.start_time,
    a.end_time,
    a.location,
    a.description,
    a.status,
    a.recurrence_type,       -- Select the recurrence_type
    a.recurrence_interval,
    a.recurrence_end_date,
    a.confirmed_by_employee_id,
    a.confirmed_at,
    a.created_at,
    a.updated_at,
    FALSE AS is_recurring_occurrence -- Flag for non-recurring
FROM
    appointments a
INNER JOIN ClientAppointments ca ON a.id = ca.appointment_id -- Join based on client participation
WHERE
    (a.recurrence_type = 'NONE' OR a.recurrence_type IS NULL) -- Select only non-recurring appointments (IS NULL check for safety)
    -- Standard overlap check: (StartA <= EndB) AND (EndA >= StartB)
    AND (a.start_time <= @end_date::timestamp) -- Use named parameter
    AND (a.end_time >= @start_date::timestamp) -- Use named parameter

UNION ALL -- Combine with recurring occurrences, keeping all rows

SELECT
    ro.original_appointment_id AS id, -- Use the original ID for consistency
    ro.creator_employee_id,
    ro.start_time, -- Calculated start time
    ro.end_time,   -- Calculated end time
    ro.location,
    ro.description,
    ro.status,
    ro.recurrence_type,
    ro.recurrence_interval,
    ro.recurrence_end_date,
    ro.confirmed_by_employee_id,
    ro.confirmed_at,
    ro.created_at,
    ro.updated_at,
    ro.is_recurring_occurrence -- Flag indicating it's a calculated occurrence
FROM
    RecurringOccurrences ro
WHERE
    -- Filter the generated occurrences to only those that OVERLAP the requested time frame
    -- Standard overlap check: (StartA <= EndB) AND (EndA >= StartB)
    (ro.start_time <= @end_date::timestamp)   -- Use named parameter
    AND (ro.end_time >= @start_date::timestamp) -- Use named parameter

-- Optional: Order the final results for consistent output
ORDER BY start_time;