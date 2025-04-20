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

-- name: GetAppointmentTemplatesForEmployee :many
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
    a.recurrence_end_date
FROM
    appointments a
LEFT JOIN appointment_participants ap ON a.id = ap.appointment_id
WHERE
    (a.creator_employee_id = $1 OR ap.employee_id = $1)
AND a.status = ANY($4::VARCHAR[])
AND a.start_time >= $3
AND (a.recurrence_end_date IS NULL OR a.recurrence_end_date >= $2)
GROUP BY a.id ;


-- name: GetParticipantsForAppointments :many
SELECT
    ap.appointment_id,
    ap.employee_id,
    e.first_name AS employee_first_name,
    e.last_name AS employee_last_name -- Join to get names for the frontend
    -- Include other employee details if needed by the frontend
FROM appointment_participants ap
JOIN employee_profile e ON ap.employee_id = e.employee_id -- Ensure 'employees' table and 'name' column exist
WHERE ap.appointment_id = ANY($1::INT[]);

-- name: GetClientsForAppointments :many
SELECT
    ac.appointment_id,
    ac.client_id,
    c.first_name AS client_first_name,
    c.last_name AS client_last_name 
FROM appointment_clients ac
JOIN client_details c ON ac.client_id = c.client_id -- Ensure 'clients' table and 'name' column exist
WHERE ac.appointment_id = ANY($1::INT[]);