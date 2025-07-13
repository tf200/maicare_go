-- name: CreateInvoice :one
INSERT INTO invoice (
    invoice_number,
    due_date,
    issue_date,
    invoice_details,
    total_amount,
    extra_content,
    client_id,
    sender_id,
    warning_count
    ) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;



-- name: ListInvoices :many
SELECT
    i.*,
    COUNT(*) OVER() AS total_count,
    s.name AS sender_name,
    cd.first_name AS client_first_name,
    cd.last_name AS client_last_name
FROM
    invoice i
JOIN
    client_details cd ON i.client_id = cd.id
LEFT JOIN
    sender s ON i.sender_id = s.id
WHERE
    -- Optional filter for client_id.
    -- The filter is applied only if the @client_id parameter is not NULL.
    (i.client_id = sqlc.narg('client_id') OR sqlc.narg('client_id') IS NULL)
    
    -- Optional filter for sender_id.
    -- The filter is applied only if the @sender_id parameter is not NULL.
    AND (i.sender_id = sqlc.narg('sender_id') OR sqlc.narg('sender_id') IS NULL)
    
    -- Optional filter for status.
    -- The filter is applied only if the @status parameter is not NULL.
    AND (i.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL)
    
    -- Optional filter for the start of the issue_date range.
    -- The filter is applied only if the @start_date parameter is not NULL.
    AND (i.issue_date >= sqlc.narg('start_date') OR sqlc.narg('start_date') IS NULL)
    
    -- Optional filter for the end of the issue_date range.
    -- The filter is applied only if the @end_date parameter is not NULL.
    AND (i.issue_date <= sqlc.narg('end_date') OR sqlc.narg('end_date') IS NULL)
ORDER BY
    i.updated_at DESC
LIMIT sqlc.arg('limit')
OFFSET sqlc.arg('offset');



-- name: GetInvoice :one
SELECT
    i.*,
    s.name AS sender_name,
    cd.first_name AS client_first_name,
    cd.last_name AS client_last_name
FROM
    invoice i
JOIN
    client_details cd ON i.client_id = cd.id
LEFT JOIN
    sender s ON i.sender_id = s.id
WHERE
    i.id = $1
LIMIT 1;


-- name: UpdateInvoice :one
UPDATE invoice
SET
    issue_date = COALESCE($2, issue_date),
    due_date = COALESCE($3, due_date),
    invoice_details = COALESCE($4, invoice_details),
    total_amount = COALESCE($5, total_amount),
    extra_content = COALESCE($6, extra_content),
    status = COALESCE($7, status),
    warning_count = COALESCE($8, warning_count)
WHERE
    id = $1
RETURNING *;


-- name: DeleteInvoice :exec
DELETE FROM invoice
WHERE id = $1;

    

