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
    s.contacts As sender_contacts,
    s.postal_code AS sender_postal_code,
    s.kvknumber AS sender_kvknumber,
    s.btwnumber AS sender_btwnumber,
    s.address AS sender_address,
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
    issue_date = COALESCE(sqlc.narg('issue_date'), issue_date),
    due_date = COALESCE(sqlc.narg('due_date'), due_date),
    invoice_details = COALESCE(sqlc.narg('invoice_details'), invoice_details),
    total_amount = COALESCE(sqlc.narg('total_amount'), total_amount),
    extra_content = COALESCE(sqlc.narg('extra_content'), extra_content),
    status = COALESCE(sqlc.narg('status'), status),
    warning_count = COALESCE(sqlc.narg('warning_count'), warning_count)
WHERE id = $1
RETURNING *;


-- name: InsertIncoicePdfUrl :one 
UPDATE invoice 
SET 
    pdf_attachment_id = $2  
WHERE
    id = $1
RETURNING invoice.pdf_attachment_id;





-- name: DeleteInvoice :exec
DELETE FROM invoice
WHERE id = $1;



-- name: GetInvoiceAuditLogs :many
SELECT
    ia.*,
    e.first_name AS changed_by_first_name,
    e.last_name AS changed_by_last_name
FROM
    invoice_audit ia
LEFT JOIN
    employee_profile e ON ia.changed_by = e.id
WHERE
    ia.invoice_id = $1
ORDER BY
    ia.changed_at DESC;





-- ////////////////////// Payments //////////////////////

-- name: CreatePayment :one
INSERT INTO invoice_payment_history (
    invoice_id,
    payment_method,
    payment_status,
    amount,
    payment_date,
    payment_reference,
    notes,
    recorded_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;



-- name: GetTotalPaidAmountByInvoice :one
SELECT 
    COALESCE(SUM(amount), 0)::FLOAT AS total_paid
FROM invoice_payment_history 
WHERE invoice_id = $1 
  AND payment_status = 'completed';


-- name: ListPayments :many
SELECT
    iph.*,
    e.first_name AS recorded_by_first_name,
    e.last_name AS recorded_by_last_name
FROM
    invoice_payment_history iph
LEFT JOIN
    employee_profile e ON iph.recorded_by = e.id
WHERE
    iph.invoice_id = $1
ORDER BY
    iph.payment_date DESC;


-- name: GetPayment :one
SELECT
    iph.*,
    e.first_name AS recorded_by_first_name,
    e.last_name AS recorded_by_last_name
FROM
    invoice_payment_history iph
LEFT JOIN
    employee_profile e ON iph.recorded_by = e.id
WHERE
    iph.id = $1
LIMIT 1;


-- name: UpdatePayment :one
UPDATE invoice_payment_history 
SET 
    payment_method = COALESCE(sqlc.narg('payment_method'), payment_method),
    payment_status = COALESCE(sqlc.narg('payment_status'), payment_status),
    amount = COALESCE(sqlc.narg('amount'), amount),
    payment_date = COALESCE(sqlc.narg('payment_date'), payment_date),
    payment_reference = COALESCE(sqlc.narg('payment_reference'), payment_reference),
    notes = COALESCE(sqlc.narg('notes'), notes),
    recorded_by = COALESCE(sqlc.narg('recorded_by'), recorded_by),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;
 

-- name: GetPaymentWithInvoice :one
SELECT 
    p.*,
    i.total_amount as invoice_total_amount,
    i.status as invoice_status
FROM invoice_payment_history p
JOIN invoice i ON p.invoice_id = i.id
WHERE p.id = $1;


-- name: DeletePayment :one
DELETE FROM invoice_payment_history
WHERE id = $1
RETURNING *;



