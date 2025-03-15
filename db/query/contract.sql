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
    financing_option
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
)
RETURNING *;


-- name: ListContracts :many
SELECT * FROM contract;


-- name: GetClientContract :one
SELECT * FROM contract
WHERE client_id = $1
limit 1;

-- name: GetSenderContracts :many
SELECT * FROM contract
WHERE sender_id = $1;

