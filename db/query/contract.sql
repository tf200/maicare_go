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
    SELECT * FROM contract
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
    type_id = COALESCE(sqlc.narg('type_id'), type_id),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = COALESCE(sqlc.narg('end_date'), end_date),
    reminder_period = COALESCE(sqlc.narg('reminder_period'), reminder_period),
    VAT = COALESCE(sqlc.narg('VAT'), VAT),
    price = COALESCE(sqlc.narg('price'), price),
    price_time_unit = COALESCE(sqlc.narg('price_time_unit'), price_time_unit),
    hours = COALESCE(sqlc.narg('hours'), hours),
    hours_type = COALESCE(sqlc.narg('hours_type'), hours_type),
    care_name = COALESCE(sqlc.narg('care_name'), care_name),
    care_type = COALESCE(sqlc.narg('care_type'), care_type),
    sender_id = COALESCE(sqlc.narg('sender_id'), sender_id),
    attachment_ids = COALESCE(sqlc.narg('attachment_ids'), attachment_ids),
    financing_act = COALESCE(sqlc.narg('financing_act'), financing_act),
    financing_option = COALESCE(sqlc.narg('financing_option'), financing_option),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateContractStatus :one
UPDATE contract
SET
    status = @status::text,
    approved_at = CASE WHEN @status::text = 'approved' THEN NOW() ELSE approved_at END
WHERE id = @contract_id::BIGINT
RETURNING *;

-- name: GetClientContract :one
SELECT c.*,
        ct.name AS contract_type_name,
        cd.first_name AS client_first_name,
        cd.last_name AS client_last_name,
        s.name AS sender_name
FROM contract c
JOIN contract_type ct ON c.type_id = ct.id
JOIN client_details cd ON c.client_id = cd.id
LEFT JOIN sender s ON c.sender_id = s.id
WHERE c.id = $1
limit 1;

-- name: GetSenderContracts :many
SELECT * FROM contract
WHERE sender_id = $1;

-- name: ListContracts :many
WITH filtered_contracts AS (
    SELECT
        c.id,
        c.status,
        c.start_date,
        c.end_date,
        c.price,
        c.price_time_unit,
        c.care_name,
        c.care_type,
        c.financing_act,
        c.financing_option,
        c.created_at,
        s.name AS sender_name,
        cd.id AS client_id,
        cd.sender_id AS sender_id,
        cd.first_name AS client_first_name,
        cd.last_name AS client_last_name
    FROM
        contract c
    LEFT JOIN
        sender s ON c.sender_id = s.id
    JOIN
        client_details cd ON c.client_id = cd.id
    WHERE
        (sqlc.narg(search)::varchar IS NULL OR 
            s.name ILIKE '%' || sqlc.narg(search) || '%' OR
            cd.first_name ILIKE '%' || sqlc.narg(search) || '%' OR
            cd.last_name ILIKE '%' || sqlc.narg(search) || '%')
    AND
        (sqlc.narg(status)::varchar[] IS NULL OR c.status = ANY(sqlc.narg(status)))
    AND
        (sqlc.narg(care_type)::varchar[] IS NULL OR c.care_type = ANY(sqlc.narg(care_type)))
    AND
        (sqlc.narg(financing_act)::varchar[] IS NULL OR c.financing_act = ANY(sqlc.narg(financing_act)))
    AND
        (sqlc.narg(financing_option)::varchar[] IS NULL OR c.financing_option = ANY(sqlc.narg(financing_option)))
)
SELECT
    (SELECT COUNT(*) FROM filtered_contracts) AS total_count,
    *
FROM
    filtered_contracts
ORDER BY
    created_at DESC
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
    ) THEN 'initial' 
    ELSE 'follow_up' END
)
RETURNING *;



-- name: GetBillablePeriodsForContract :many
WITH status_history AS (
  SELECT
    new_values->>'status' AS status,
    changed_at AS effective_date
  FROM contract_audit
  WHERE contract_id = sqlc.arg(contract_id)
  
  UNION ALL
  
  SELECT
    status,
    updated_at
  FROM contract
  WHERE id = sqlc.arg(contract_id)
  ORDER BY effective_date
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