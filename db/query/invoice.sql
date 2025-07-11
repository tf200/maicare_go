-- name : CreateInvoice :one
INSERT INTO invoice (
    invoice_number,
    isuue_date,
    due_date,
    invoice_details,
    total_amount,
    extra_content,
    client_id
    ) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7
) RETURNING *;