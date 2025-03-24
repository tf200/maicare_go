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
    start_date,
    end_date,
    reminder_period,
    tax,
    price,
    price_frequency,
    hours,
    hours_type,
    care_name,
    care_type,
    client_id,
    sender_id,
    attachment_ids,
    financing_act,
    financing_option,
    status
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
ORDER BY created DESC
LIMIT $2
OFFSET $3;



-- name: GetClientContract :one
SELECT * FROM contract
WHERE id = $1
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
        c.price_frequency,
        c.care_name,
        c.care_type,
        c.financing_act,
        c.financing_option,
        c.created,
        s.name AS sender_name,
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
    created DESC
LIMIT $1
OFFSET $2;