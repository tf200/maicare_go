-- name: DischargeOverview :many
SELECT
    cd.id,
    cd.first_name,
    cd.last_name,
    cd.status AS current_status,
    NULL::client_status_enum AS scheduled_status,
    NULL::TEXT AS status_change_reason,
    NULL::TIMESTAMPTZ AS status_change_date,
    COALESCE(c.end_date, CURRENT_TIMESTAMP) AS contract_end_date,
    c.status AS contract_status,
    c.departure_reason,
    c.departure_report AS follow_up_plan,
    'status_change'::TEXT AS discharge_type
FROM client_details cd
LEFT JOIN LATERAL (
    SELECT *
    FROM contract c2
    WHERE c2.client_id = cd.id
    ORDER BY c2.end_date DESC
    LIMIT 1
) c ON TRUE
WHERE sqlc.arg('filter_type')::TEXT IN ('urgent', 'contract', 'status_change', 'all')
ORDER BY cd.created_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');


-- name: TotalDischargeCount :one
SELECT COUNT(*)::BIGINT
FROM client_details;


-- name: UrgentCasesCount :one
SELECT 0::BIGINT;


-- name: StatusChangeCount :one
SELECT 0::BIGINT;


-- name: ContractEndCount :one
SELECT COUNT(*)::BIGINT
FROM contract c
WHERE c.end_date IS NOT NULL
  AND c.end_date <= CURRENT_DATE + INTERVAL '3 months';


-- name: ListEmployeesByContractEndDate :many
SELECT
    ep.id,
    ep.user_id,
    ep.first_name,
    ep.last_name,
    ep.position,
    d.name AS department_name,
    ep.employee_number,
    ep.work_email_address,
    ep.contract_start_date,
    ep.contract_end_date,
    ep.contract_type
FROM employee_profile ep
LEFT JOIN departments d ON d.id = ep.department_id
WHERE contract_end_date IS NOT NULL
  AND ep.is_archived = FALSE
ORDER BY ep.contract_end_date ASC
LIMIT 10;


-- name: ListLatestPayments :many
SELECT
    i.id AS invoice_id,
    i.invoice_number,
    iph.payment_method,
    iph.payment_status,
    iph.amount,
    iph.payment_date,
    iph.updated_at
FROM invoice_payment_history iph
JOIN invoice i ON iph.invoice_id = i.id
ORDER BY iph.updated_at DESC
LIMIT 10;


-- name: ListUpcomingAppointments :many
SELECT
    ce.id,
    ce.start_at AS start_time,
    ce.end_at AS end_time,
    ce.location,
    ce.description
FROM
    calendar_events ce
LEFT JOIN
    calendar_event_attendees cea ON ce.id = cea.event_id
WHERE
    ce.kind = 'appointment'
    AND ce.status <> 'cancelled'
    AND ce.start_at >= CURRENT_TIMESTAMP
    AND (ce.organizer_employee_id = $1 OR cea.employee_id = $1)
ORDER BY
    ce.start_at ASC
LIMIT 10;
