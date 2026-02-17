-- name: CreateContractType :one
INSERT INTO contract_type (name)
VALUES
    ($1)
RETURNING *;

-- name: ListContractTypes :many
SELECT * FROM contract_type;

-- name: DeleteContractType :exec
DELETE FROM contract_type
WHERE id = $1;



-- name: CreateContract :one
INSERT INTO contract (
    type_id,
    status,
    start_date,
    end_date,
    reminder_period,
    VAT,
    price,
    price_time_unit,
    hours,
    hours_type,
    care_name,
    care_type,
    client_id,
    sender_id,
    attachment_ids,
    financing_act,
    financing_option
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
)
RETURNING *;

-- name: ListClientContracts :many
WITH client_contracts AS (
    SELECT
            c.start_date,
            c.end_date,
            (c.end_date::date - CURRENT_DATE)::int AS days_left,
            c.care_name,
            c.care_type,
            c.financing_act,
            c.financing_option,
            c.created_at
    FROM contract c
    WHERE client_id = $1
)
SELECT
    (SELECT COUNT(*) FROM client_contracts) AS total_count,
    *
FROM client_contracts
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;


-- name: UpdateContract :one
UPDATE contract
SET 
    type_id = $2,
    start_date = $3,
    end_date = $4,
    reminder_period = $5,
    VAT = $6,
    price = $7,
    price_time_unit = $8,
    hours = $9,
    hours_type = $10,
    care_name = $11,
    care_type = $12,
    sender_id = $13,
    attachment_ids = $14,
    financing_act = $15,
    financing_option = $16,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateContractStatus :one
UPDATE contract
SET
    status = @status::contract_status_enum,
    approved_at = CASE WHEN @status::contract_status_enum = 'approved' THEN NOW() ELSE approved_at END
WHERE id = @contract_id::uuid
RETURNING *;

-- name: GetClientContract :one
SELECT c.*,
        ct.name AS contract_type_name,
        cd.first_name AS client_first_name,
        cd.last_name AS client_last_name,
        cd.filenumber AS client_filenumber,
        cd.bsn AS client_bsn,
        s.name AS sender_name,
        s.types AS sender_type,
        s.street AS sender_street,
        s.house_number AS sender_house_number,
        s.house_number_addition AS sender_house_number_addition,
        s.postal_code AS sender_postal_code,
        s.city AS sender_city,
        s.land AS sender_land,
        s.kvknumber AS sender_kvknumber,
        s.btwnumber AS sender_btwnumber,
        s.phone_number AS sender_phone_number,
        s.client_number AS sender_client_number,
        s.email_address AS sender_email_address
FROM contract c
LEFT JOIN contract_type ct ON c.type_id = ct.id
JOIN client_details cd ON c.client_id = cd.id
JOIN sender s ON c.sender_id = s.id
WHERE c.id = $1
LIMIT 1;

-- name: GetSenderContracts :many
SELECT * FROM contract
WHERE sender_id = $1;

-- name: ListContracts :many
WITH filtered_contracts AS (
    SELECT
        c.id,
        c.client_id,
        c.status,
        c.approved_at,
        c.start_date,
        c.end_date,
        (c.end_date::date - CURRENT_DATE)::int AS days_left,
        c.price,
        c.price_time_unit,
        c.hours,
        c.hours_type,
        c.care_name,
        c.care_type,
        c.financing_act,
        c.financing_option,
        c.updated_at,
        s.name AS sender_name,
        c.sender_id AS sender_id,
        cd.filenumber AS client_filenumber,
        cd.first_name AS client_first_name,
        cd.last_name AS client_last_name
    FROM
        contract c
    JOIN
        sender s ON c.sender_id = s.id
    JOIN
        client_details cd ON c.client_id = cd.id
    WHERE
        (sqlc.narg(search)::varchar IS NULL OR 
            s.name ILIKE '%' || sqlc.narg(search) || '%' OR
            (cd.first_name || ' ' || cd.last_name) ILIKE '%' || sqlc.narg(search) || '%' OR
            (cd.last_name || ' ' || cd.first_name) ILIKE '%' || sqlc.narg(search) || '%' OR
            cd.first_name ILIKE '%' || sqlc.narg(search) || '%' OR
            cd.last_name ILIKE '%' || sqlc.narg(search) || '%')
    AND
        (sqlc.narg(status)::contract_status_enum[] IS NULL OR c.status = ANY(sqlc.narg(status)))
    AND
        (sqlc.narg(care_type)::care_type_enum[] IS NULL OR c.care_type = ANY(sqlc.narg(care_type)))
    AND
        (sqlc.narg(financing_act)::financing_act_enum[] IS NULL OR c.financing_act = ANY(sqlc.narg(financing_act)))
    AND
        (sqlc.narg(financing_option)::financing_option_enum IS NULL OR c.financing_option = sqlc.narg(financing_option))
    AND
        (sqlc.narg(end_date_from)::timestamptz IS NULL OR c.end_date >= sqlc.narg(end_date_from))
    AND
        (sqlc.narg(end_date_to)::timestamptz IS NULL OR c.end_date <= sqlc.narg(end_date_to))
)
SELECT
    (SELECT COUNT(*) FROM filtered_contracts) AS total_count,
    *
FROM
    filtered_contracts
ORDER BY
    updated_at DESC
LIMIT $1
OFFSET $2;




-- name: ListContractsTobeReminded :many
SELECT c.id,
       c.care_name,
       c.client_id,
       c.start_date,
       c.end_date,
       c.reminder_period,
       c.care_type,
       cd.id AS client_id,
         cd.first_name AS client_first_name,
         cd.last_name AS client_last_name,

       (c.end_date - INTERVAL '1 day' * c.reminder_period) AS reminder_date,

       COALESCE(MAX(cr.reminder_sent_at), '1970-01-01'::TIMESTAMPTZ)::TIMESTAMPTZ AS last_reminder_date

FROM contract c
JOIN client_details cd ON c.client_id = cd.id
LEFT JOIN contract_reminder cr ON c.id = cr.contract_id
    AND cr.reminder_sent_at IS NOT NULL

WHERE 
    c.status = 'approved'


    AND CURRENT_DATE >= (c.end_date - INTERVAL '1 day' * c.reminder_period)::date
    AND c.end_date > CURRENT_TIMESTAMP

GROUP BY c.id, c.care_name, c.client_id, c.end_date, c.reminder_period
HAVING 
    (MAX(cr.reminder_sent_at) IS NULL OR
      MAX(cr.reminder_sent_at) < CURRENT_TIMESTAMP - INTERVAL '7 days')
ORDER BY c.end_date ASC;
    



-- name: CreateContractReminder :one
INSERT INTO contract_reminder (
    contract_id,
    reminder_sent_at,
    reminder_type
) VALUES (
    $1, $2,
    CASE WHEN NOT EXISTS (
        SELECT 1 FROM contract_reminder
        WHERE contract_id = $1
        AND reminder_sent_at IS NOT NULL
    ) THEN 'initial'::contract_reminder_type_enum
    ELSE 'follow_up'::contract_reminder_type_enum END
)
RETURNING *;




-- name: GetBillablePeriodsForContract :many
WITH raw_status_changes AS (
  SELECT
    (new_values->>'status')::contract_status_enum AS status,
    changed_at AS effective_date
  FROM contract_audit
  WHERE contract_id = sqlc.arg(contract_id)
    AND new_values->>'status' IS NOT NULL
  
  UNION ALL
  
  SELECT
    status,
    updated_at
  FROM contract
  WHERE id = sqlc.arg(contract_id)
),

-- Deduplicate identical status within 1 second windows
status_history AS (
  SELECT
    status,
    MIN(effective_date) AS effective_date
  FROM raw_status_changes
  GROUP BY status, EXTRACT(EPOCH FROM effective_date)::INTEGER
  ORDER BY MIN(effective_date)
),

approved_periods AS (
  SELECT
    effective_date AS period_start,
    LEAD(effective_date, 1) OVER (ORDER BY effective_date) AS period_end
  FROM status_history
  WHERE status = 'approved'
)

SELECT
  GREATEST(period_start, sqlc.arg(invoice_start_date))::TIMESTAMPTZ AS billable_start,
  LEAST(COALESCE(period_end, sqlc.arg(invoice_end_date)), sqlc.arg(invoice_end_date))::TIMESTAMPTZ AS billable_end
FROM approved_periods
WHERE period_start < sqlc.arg(invoice_end_date)
  AND COALESCE(period_end, 'infinity'::TIMESTAMPTZ) > sqlc.arg(invoice_start_date);




-- name: GetContractAudit :many
SELECT ca.*,
         e.first_name AS changed_by_first_name,
         e.last_name AS changed_by_last_name
FROM contract_audit ca
LEFT JOIN employee_profile e ON ca.changed_by = e.id
WHERE ca.contract_id = $1
ORDER BY ca.changed_at DESC;
