-- ////////////////////// Invoices //////////////////////

-- name: CreateInvoice :one
INSERT INTO invoice (
    invoice_number,
    invoice_sequence,
    due_date,
    issue_date,
    status,
    invoice_type,
    source,
    original_invoice_id,
    replaces_invoice_id,
    period_start,
    period_end,
    billing_cycle,
    billing_timezone,
    currency,
    extra_content,
    client_id,
    sender_id,
    warning_count,
    run_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19
) RETURNING *;


-- name: ListInvoices :many
WITH params AS (
    SELECT
        sqlc.arg('sort_by')::text AS sort_by,
        sqlc.arg('sort_dir')::text AS sort_dir
),
paid AS (
    SELECT
        iph.invoice_id,
        COALESCE(SUM(CASE WHEN iph.payment_status = 'completed' THEN iph.amount ELSE 0 END), 0)::NUMERIC(20,2) AS paid_total_amount
    FROM invoice_payment_history iph
    GROUP BY iph.invoice_id
),
base AS (
    SELECT
        i.id,
        i.invoice_number,
        i.issue_date,
        i.due_date,
        i.status,
        i.currency,
        i.gross_total_amount,
        i.client_id,
        i.sender_id,
        i.updated_at,

        COALESCE(NULLIF(i.bill_to_snapshot->>'name', ''), s.name) AS sender_name,
        COALESCE(NULLIF(i.client_snapshot->>'first_name', ''), cd.first_name) AS client_first_name,
        COALESCE(NULLIF(i.client_snapshot->>'last_name', ''), cd.last_name) AS client_last_name,
        COALESCE(NULLIF(i.client_snapshot->>'filenumber', ''), cd.filenumber) AS client_filenumber,

        COALESCE(p.paid_total_amount, 0)::NUMERIC(20,2) AS paid_total_amount,
        GREATEST(i.gross_total_amount - COALESCE(p.paid_total_amount, 0), 0)::NUMERIC(20,2) AS balance_due_amount,
        CASE
            WHEN i.gross_total_amount <= 0 THEN 0
            ELSE ROUND((COALESCE(p.paid_total_amount, 0) / i.gross_total_amount) * 100, 2)
        END::NUMERIC(10,2) AS payment_completion_prc,
        COALESCE((i.due_date < CURRENT_DATE AND i.status NOT IN ('paid', 'canceled')), false)::boolean AS is_overdue
    FROM invoice i
    JOIN client_details cd ON i.client_id = cd.id
    LEFT JOIN sender s ON i.sender_id = s.id
    LEFT JOIN paid p ON p.invoice_id = i.id
    WHERE
        (i.client_id = sqlc.narg('client_id') OR sqlc.narg('client_id') IS NULL)
        AND (i.sender_id = sqlc.narg('sender_id') OR sqlc.narg('sender_id') IS NULL)

        AND (i.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL)
        AND (sqlc.narg('statuses')::invoice_status_enum[] IS NULL OR i.status = ANY(sqlc.narg('statuses')::invoice_status_enum[]))

        AND (i.source = sqlc.narg('source') OR sqlc.narg('source') IS NULL)
        AND (i.invoice_type = sqlc.narg('invoice_type') OR sqlc.narg('invoice_type') IS NULL)
        AND (i.run_id = sqlc.narg('run_id') OR sqlc.narg('run_id') IS NULL)

        AND (i.issue_date >= sqlc.narg('start_date') OR sqlc.narg('start_date') IS NULL)
        AND (i.issue_date <= sqlc.narg('end_date') OR sqlc.narg('end_date') IS NULL)

        AND (i.period_start >= sqlc.narg('period_start') OR sqlc.narg('period_start') IS NULL)
        AND (i.period_end <= sqlc.narg('period_end') OR sqlc.narg('period_end') IS NULL)

        AND (i.warning_count >= sqlc.narg('min_warning_count') OR sqlc.narg('min_warning_count') IS NULL)

        AND (
            sqlc.narg('locked')::boolean IS NULL
            OR (sqlc.narg('locked') = true AND i.locked_at IS NOT NULL)
            OR (sqlc.narg('locked') = false AND i.locked_at IS NULL)
        )

        AND (
            sqlc.narg('q')::text IS NULL
            OR i.invoice_number ILIKE ('%' || sqlc.narg('q')::text || '%')
            OR COALESCE(NULLIF(i.bill_to_snapshot->>'name', ''), s.name) ILIKE ('%' || sqlc.narg('q')::text || '%')
            OR COALESCE(NULLIF(i.client_snapshot->>'filenumber', ''), cd.filenumber) ILIKE ('%' || sqlc.narg('q')::text || '%')
            OR (COALESCE(NULLIF(i.client_snapshot->>'first_name', ''), cd.first_name) || ' ' || COALESCE(NULLIF(i.client_snapshot->>'last_name', ''), cd.last_name)) ILIKE ('%' || sqlc.narg('q')::text || '%')
        )
),
paged AS (
    SELECT
        b.id,
        b.invoice_number,
        b.issue_date,
        b.due_date,
        b.status,
        b.currency,
        b.gross_total_amount,
        b.client_id,
        b.sender_id,
        b.updated_at,
        b.sender_name,
        b.client_first_name,
        b.client_last_name,
        b.client_filenumber,
        b.paid_total_amount,
        b.balance_due_amount,
        b.is_overdue,
        COUNT(*) OVER() AS total_count
    FROM base b
    CROSS JOIN params p
    ORDER BY
        CASE WHEN p.sort_by = 'updated_at' AND p.sort_dir = 'asc' THEN b.updated_at END ASC,
        CASE WHEN p.sort_by = 'updated_at' AND p.sort_dir = 'desc' THEN b.updated_at END DESC,

        CASE WHEN p.sort_by = 'issue_date' AND p.sort_dir = 'asc' THEN b.issue_date END ASC,
        CASE WHEN p.sort_by = 'issue_date' AND p.sort_dir = 'desc' THEN b.issue_date END DESC,

        CASE WHEN p.sort_by = 'due_date' AND p.sort_dir = 'asc' THEN b.due_date END ASC,
        CASE WHEN p.sort_by = 'due_date' AND p.sort_dir = 'desc' THEN b.due_date END DESC,

        CASE WHEN p.sort_by = 'invoice_number' AND p.sort_dir = 'asc' THEN b.invoice_number END ASC,
        CASE WHEN p.sort_by = 'invoice_number' AND p.sort_dir = 'desc' THEN b.invoice_number END DESC,

        CASE WHEN p.sort_by = 'gross_total_amount' AND p.sort_dir = 'asc' THEN b.gross_total_amount END ASC,
        CASE WHEN p.sort_by = 'gross_total_amount' AND p.sort_dir = 'desc' THEN b.gross_total_amount END DESC,

        CASE WHEN p.sort_by = 'balance_due_amount' AND p.sort_dir = 'asc' THEN b.balance_due_amount END ASC,
        CASE WHEN p.sort_by = 'balance_due_amount' AND p.sort_dir = 'desc' THEN b.balance_due_amount END DESC,

        b.updated_at DESC,
        b.id DESC
    LIMIT sqlc.arg('limit')
    OFFSET sqlc.arg('offset')
)
SELECT
    id,
    invoice_number,
    issue_date,
    due_date,
    status,
    currency,
    gross_total_amount,
    paid_total_amount,
    balance_due_amount,
    is_overdue,
    client_id,
    sender_id,
    sender_name,
    client_first_name,
    client_last_name,
    client_filenumber,
    total_count
FROM paged;


-- name: GetInvoice :one
SELECT
    i.*,
    s.name AS sender_name,
    s.contacts As sender_contacts,
    s.postal_code AS sender_postal_code,
    s.kvknumber AS sender_kvknumber,
    s.btwnumber AS sender_btwnumber,
    s.street AS sender_street,
    s.house_number AS sender_house_number,
    s.house_number_addition AS sender_house_number_addition,
    s.city AS sender_city,
    s.land AS sender_land,
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


-- name: GetMaxInvoiceSequenceForDate :one
SELECT COALESCE(MAX(invoice_sequence), 0)::BIGINT as max_sequence
FROM invoice
WHERE DATE(created_at) = DATE($1);


-- name: UpdateInvoice :one
UPDATE invoice
SET
    issue_date = COALESCE(sqlc.narg('issue_date'), issue_date),
    due_date = COALESCE(sqlc.narg('due_date'), due_date),
    period_start = COALESCE(sqlc.narg('period_start'), period_start),
    period_end = COALESCE(sqlc.narg('period_end'), period_end),
    billing_cycle = COALESCE(sqlc.narg('billing_cycle'), billing_cycle),
    billing_timezone = COALESCE(sqlc.narg('billing_timezone'), billing_timezone),
    source = COALESCE(sqlc.narg('source'), source),
    original_invoice_id = COALESCE(sqlc.narg('original_invoice_id'), original_invoice_id),
    replaces_invoice_id = COALESCE(sqlc.narg('replaces_invoice_id'), replaces_invoice_id),
    bill_to_snapshot = COALESCE(sqlc.narg('bill_to_snapshot'), bill_to_snapshot),
    client_snapshot = COALESCE(sqlc.narg('client_snapshot'), client_snapshot),
    details_snapshot = COALESCE(sqlc.narg('details_snapshot'), details_snapshot),
    net_total_amount = COALESCE(sqlc.narg('net_total_amount'), net_total_amount),
    vat_total_amount = COALESCE(sqlc.narg('vat_total_amount'), vat_total_amount),
    gross_total_amount = COALESCE(sqlc.narg('gross_total_amount'), gross_total_amount),
    currency = COALESCE(sqlc.narg('currency'), currency),
    extra_content = COALESCE(sqlc.narg('extra_content'), extra_content),
    status = COALESCE(sqlc.narg('status'), status),
    warning_count = COALESCE(sqlc.narg('warning_count'), warning_count),
    run_id = COALESCE(sqlc.narg('run_id'), run_id),
    locked_at = COALESCE(sqlc.narg('locked_at'), locked_at),
    calc_version = COALESCE(sqlc.narg('calc_version'), calc_version),
    calc_metadata = COALESCE(sqlc.narg('calc_metadata'), calc_metadata),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;


-- name: UpdateInvoiceStatus :one
UPDATE invoice
SET
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;


-- name: InsertIncoicePdfUrl :one
UPDATE invoice
SET
    pdf_attachment_id = $2
WHERE
    id = $1
    AND pdf_attachment_id IS NULL
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


-- ////////////////////// Invoice Lines //////////////////////

-- name: CreateInvoiceLine :one
INSERT INTO invoice_line (
    invoice_id,
    client_id,
    sender_id,
    line_no,
    line_type,
    contract_id,
    service_type,
    description,
    period_start,
    period_end,
    quantity,
    unit,
    unit_price,
    net_amount,
    vat_rate,
    vat_amount,
    gross_amount,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18
) RETURNING *;


-- name: ListInvoiceLinesByInvoice :many
SELECT *
FROM invoice_line
WHERE invoice_id = $1
ORDER BY line_no;


-- name: DeleteInvoiceLinesByInvoice :exec
DELETE FROM invoice_line
WHERE invoice_id = $1;


-- name: CountBilledCalendarEventsByInvoice :one
SELECT COUNT(*)::BIGINT AS count
FROM billed_calendar_event
WHERE invoice_id = $1;


-- name: CountInvoiceLineCalendarEventsByInvoice :one
SELECT COUNT(*)::BIGINT AS count
FROM invoice_line_calendar_event ilce
JOIN invoice_line il ON ilce.invoice_line_id = il.id
WHERE il.invoice_id = $1;


-- name: UpdateInvoiceLine :one
UPDATE invoice_line
SET
    line_type = COALESCE(sqlc.narg('line_type'), line_type),
    contract_id = COALESCE(sqlc.narg('contract_id'), contract_id),
    service_type = COALESCE(sqlc.narg('service_type'), service_type),
    description = COALESCE(sqlc.narg('description'), description),
    period_start = COALESCE(sqlc.narg('period_start'), period_start),
    period_end = COALESCE(sqlc.narg('period_end'), period_end),
    quantity = COALESCE(sqlc.narg('quantity'), quantity),
    unit = COALESCE(sqlc.narg('unit'), unit),
    unit_price = COALESCE(sqlc.narg('unit_price'), unit_price),
    net_amount = COALESCE(sqlc.narg('net_amount'), net_amount),
    vat_rate = COALESCE(sqlc.narg('vat_rate'), vat_rate),
    vat_amount = COALESCE(sqlc.narg('vat_amount'), vat_amount),
    gross_amount = COALESCE(sqlc.narg('gross_amount'), gross_amount),
    metadata = COALESCE(sqlc.narg('metadata'), metadata)
WHERE id = sqlc.arg('id')
RETURNING *;


-- ////////////////////// Line Sources (Appointments) //////////////////////

-- name: CreateInvoiceLineCalendarEvent :one
INSERT INTO invoice_line_calendar_event (
    invoice_line_id,
    calendar_event_id,
    client_id,
    start_at,
    end_at,
    minutes_billed,
    metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;


-- name: InsertBilledCalendarEvent :one
INSERT INTO billed_calendar_event (
    calendar_event_id,
    client_id,
    invoice_id,
    invoice_line_id
) VALUES (
    $1, $2, $3, $4
) RETURNING *;


-- name: VoidBilledCalendarEventsByInvoice :exec
UPDATE billed_calendar_event
SET voided_at = CURRENT_TIMESTAMP
WHERE invoice_id = $1
  AND voided_at IS NULL;


-- ////////////////////// Batch Runs //////////////////////

-- name: CreateInvoiceRun :one
INSERT INTO invoice_run (
    billing_cycle,
    period_start,
    period_end,
    timezone,
    dry_run,
    status,
    params,
    created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;


-- name: UpdateInvoiceRun :one
UPDATE invoice_run
SET
    status = COALESCE(sqlc.narg('status'), status),
    finished_at = COALESCE(sqlc.narg('finished_at'), finished_at)
WHERE id = $1
RETURNING *;


-- name: CreateInvoiceRunItem :one
INSERT INTO invoice_run_item (
    run_id,
    client_id,
    sender_id,
    status,
    invoice_id,
    error,
    warnings
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;


-- name: UpdateInvoiceRunItem :one
UPDATE invoice_run_item
SET
    status = COALESCE(sqlc.narg('status'), status),
    invoice_id = COALESCE(sqlc.narg('invoice_id'), invoice_id),
    error = COALESCE(sqlc.narg('error'), error),
    warnings = COALESCE(sqlc.narg('warnings'), warnings)
WHERE id = $1
RETURNING *;


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


-- name: GetCompletedPaymentSum :one
SELECT COALESCE(SUM(amount), 0)::DECIMAL as total_completed_amount
FROM invoice_payment_history
WHERE invoice_id = $1
  AND payment_status = 'completed';


-- name: GetPaymentWithInvoice :one
SELECT
    p.*,
    i.gross_total_amount as invoice_total_amount,
    i.status as invoice_status
FROM invoice_payment_history p
JOIN invoice i ON p.invoice_id = i.id
WHERE p.id = $1;


-- name: DeletePayment :one
DELETE FROM invoice_payment_history
WHERE id = $1
RETURNING *;


-- name: GetInvoiceSenderID :one
SELECT sender_id
FROM invoice
WHERE id = $1;
