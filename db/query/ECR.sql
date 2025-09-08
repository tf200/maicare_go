-- name: ListEmployeesByContractEndDate :many
SELECT 
    id,
    user_id,
    first_name,
    last_name,
    position,
    department,
    employee_number,
    employment_number,
    email,
    contract_start_date,
    contract_end_date,
    contract_type
FROM employee_profile 
WHERE contract_end_date IS NOT NULL
  AND is_archived = FALSE
ORDER BY contract_end_date ASC
LIMIT 10;





-- name: ListLatestPayments :many
SELECT
    i.id as invoice_id,
    i.invoice_number,
    iph.payment_method,
    iph.payment_status,
    iph.amount,
    iph.payment_date,
    iph.updated_at
FROM
    invoice_payment_history iph
JOIN
    invoice i ON iph.invoice_id = i.id
ORDER BY
    iph.updated_at DESC
LIMIT 10;


-- name: ListUpcomingAppointments :many
SELECT
    t1.id,
    t1.start_time, 
    t1.end_time, 
    t1.location, 
    t1.description
FROM 
    scheduled_appointments AS t1
LEFT JOIN 
    appointment_participants AS t2 
ON 
    t1.id = t2.appointment_id
WHERE 
    t1.creator_employee_id = $1 
    OR t2.employee_id = $1
ORDER BY 
    t1.start_time ASC
LIMIT 10;





-- name: DischargeOverview :many
WITH client_discharges AS (
    -- Get clients with scheduled status change to "Out Of Care"
    SELECT 
        cd.id, 
        cd.first_name,
        cd.last_name,
        cd.status AS current_status, 
        ssc.new_status AS scheduled_status, 
        ssc.reason AS status_change_reason, 
        ssc.scheduled_date AS status_change_date, 
        c.end_date::date AS contract_end_date, 
        c.status AS contract_status, 
        c.departure_reason, 
        c.departure_report AS follow_up_plan,
        'scheduled_status' AS discharge_type
    FROM client_details cd
    JOIN scheduled_status_changes ssc ON cd.id = ssc.client_id
    LEFT JOIN contract c ON cd.id = c.client_id AND c.status = 'approved'
    WHERE cd.status = 'In Care' 
      AND ssc.new_status = 'Out Of Care'
      AND ssc.scheduled_date <= CURRENT_DATE + INTERVAL '3 months'
    
    UNION ALL
    
    -- Get clients with contracts ending in the next 3 months
    SELECT 
        cd.id, 
        cd.first_name,
        cd.last_name,
        cd.status AS current_status, 
        NULL AS scheduled_status, 
        NULL AS status_change_reason, 
        NULL AS status_change_date, 
        c.end_date::date AS contract_end_date, 
        c.status AS contract_status, 
        c.departure_reason, 
        c.departure_report AS follow_up_plan,
        'contract_end' AS discharge_type
    FROM client_details cd
    JOIN contract c ON cd.id = c.client_id
    WHERE cd.status = 'In Care' 
      AND c.status = 'approved'
      AND c.end_date <= CURRENT_DATE + INTERVAL '3 months'
      -- Exclude clients who are already included in the scheduled status changes
      AND NOT EXISTS (
          SELECT 1 
          FROM scheduled_status_changes ssc 
          WHERE cd.id = ssc.client_id 
            AND ssc.new_status = 'Out Of Care'
            AND ssc.scheduled_date <= CURRENT_DATE + INTERVAL '3 months'
      )
)
SELECT * FROM client_discharges 
WHERE 
    -- Filter based on parameter filter_type:
    -- 'all' or NULL = Show all (default)
    -- 'status_change' = Show only status changes within 3 months
    -- 'contract' = Show only contract endings within 3 months
    -- 'urgent' = Show both status changes and contract endings within 1 month
    (@filter_type::text IS NULL OR @filter_type::text = 'all') OR 
    (@filter_type::text = 'status_change' AND discharge_type = 'scheduled_status' AND status_change_date <= CURRENT_DATE + INTERVAL '3 months') OR
    (@filter_type::text = 'contract' AND discharge_type = 'contract_end' AND contract_end_date <= CURRENT_DATE + INTERVAL '3 months') OR
    (@filter_type::text = 'urgent' AND (
        (discharge_type = 'scheduled_status' AND status_change_date <= CURRENT_DATE + INTERVAL '1 month') OR
        (discharge_type = 'contract_end' AND contract_end_date <= CURRENT_DATE + INTERVAL '1 month')
    ))
ORDER BY 
    CASE WHEN discharge_type = 'scheduled_status' THEN status_change_date ELSE contract_end_date END ASC
LIMIT $1 OFFSET $2;


-- name: TotalDischargeCount :one
WITH client_discharges AS (
    -- Get clients with scheduled status change to "Out Of Care"
    SELECT cd.id, 'scheduled_status' AS discharge_type
    FROM client_details cd
    JOIN scheduled_status_changes ssc ON cd.id = ssc.client_id
    WHERE cd.status = 'In Care' 
      AND ssc.new_status = 'Out Of Care'
      AND ssc.scheduled_date <= CURRENT_DATE + INTERVAL '3 months'
    
    UNION ALL
    
    -- Get clients with contracts ending in the next 3 months
    SELECT cd.id, 'contract_end' AS discharge_type
    FROM client_details cd
    JOIN contract c ON cd.id = c.client_id
    WHERE cd.status = 'In Care' 
      AND c.status = 'approved'
      AND c.end_date <= CURRENT_DATE + INTERVAL '3 months'
      -- Exclude clients who are already included in the scheduled status changes
      AND NOT EXISTS (
          SELECT 1 
          FROM scheduled_status_changes ssc 
          WHERE cd.id = ssc.client_id 
            AND ssc.new_status = 'Out Of Care'
            AND ssc.scheduled_date <= CURRENT_DATE + INTERVAL '3 months'
      )
)
SELECT COUNT(*) as total_discharges
FROM client_discharges;


-- name: UrgentCasesCount :one
WITH client_discharges AS (
    -- Get clients with scheduled status change to "Out Of Care"
    SELECT 
        cd.id,
        ssc.scheduled_date AS relevant_date
    FROM client_details cd
    JOIN scheduled_status_changes ssc ON cd.id = ssc.client_id
    WHERE cd.status = 'In Care' 
      AND ssc.new_status = 'Out Of Care'
      AND ssc.scheduled_date <= CURRENT_DATE + INTERVAL '30 days'
    
    UNION ALL
    
    -- Get clients with contracts ending in the next 30 days
    SELECT 
        cd.id,
        c.end_date AS relevant_date
    FROM client_details cd
    JOIN contract c ON cd.id = c.client_id
    WHERE cd.status = 'In Care' 
      AND c.status = 'approved'
      AND c.end_date <= CURRENT_DATE + INTERVAL '30 days'
      -- Exclude clients who are already included in the scheduled status changes
      AND NOT EXISTS (
          SELECT 1 
          FROM scheduled_status_changes ssc 
          WHERE cd.id = ssc.client_id 
            AND ssc.new_status = 'Out Of Care'
            AND ssc.scheduled_date <= CURRENT_DATE + INTERVAL '30 days'
      )
)
SELECT COUNT(*) as urgent_count
FROM client_discharges;



-- name: StatusChangeCount :one
SELECT COUNT(*) as status_changes_count
FROM client_details cd
JOIN scheduled_status_changes ssc ON cd.id = ssc.client_id
WHERE cd.status = 'In Care' 
  AND ssc.new_status = 'Out Of Care'
  AND ssc.scheduled_date <= CURRENT_DATE + INTERVAL '3 months';



-- name: ContractEndCount :one
SELECT COUNT(*) as contract_end_count
FROM client_details cd
JOIN contract c ON cd.id = c.client_id
WHERE cd.status = 'In Care' 
  AND c.status = 'approved'
  AND c.end_date <= CURRENT_DATE + INTERVAL '3 months'
  -- Exclude clients who are already included in the scheduled status changes
  AND NOT EXISTS (
      SELECT 1 
      FROM scheduled_status_changes ssc 
      WHERE cd.id = ssc.client_id 
        AND ssc.new_status = 'Out Of Care'
        AND ssc.scheduled_date <= CURRENT_DATE + INTERVAL '3 months'
  );




-- name: TotalActiveClients :one
SELECT COUNT(id) AS total_active_clients 
FROM client_details
WHERE status = 'In Care';

-- name: ClientsOnWaitlist :one
SELECT COUNT(id) AS total_clients_on_waitlist
FROM client_details
WHERE status = 'On Waitlist';

-- name: RecentIncidents :one
SELECT COUNT(id) AS total_recent_incidents
FROM incident
WHERE created_at >= CURRENT_DATE - INTERVAL '48 hours';


