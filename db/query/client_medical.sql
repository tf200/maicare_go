-- ==========================================
-- Client Medical (v2)
-- Diagnoses + Medication Orders
-- ==========================================

-- name: CreateClientDiagnosis :one
INSERT INTO client_diagnosis (
    client_id,
    code_system,
    code,
    title,
    description,
    status,
    severity,
    diagnosed_on,
    resolved_on,
    diagnosing_clinician,
    notes,
    created_by_employee_id,
    updated_by_employee_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING
    id,
    client_id,
    code_system,
    code,
    title,
    description,
    status,
    severity,
    diagnosed_on,
    resolved_on,
    diagnosing_clinician,
    notes,
    created_by_employee_id,
    updated_by_employee_id,
    created_at,
    updated_at,
    archived_at;


-- name: ListClientDiagnoses :many
SELECT
    d.id,
    d.client_id,
    d.code_system,
    d.code,
    d.title,
    d.description,
    d.status,
    d.severity,
    d.diagnosed_on,
    d.resolved_on,
    d.diagnosing_clinician,
    d.notes,
    d.created_by_employee_id,
    d.updated_by_employee_id,
    d.created_at,
    d.updated_at,
    d.archived_at,
    COUNT(*) OVER() AS total_count
FROM client_diagnosis d
WHERE d.client_id = $1
  AND d.archived_at IS NULL
ORDER BY COALESCE(d.diagnosed_on, d.created_at::date) DESC, d.created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetClientDiagnosis :one
SELECT
    d.id,
    d.client_id,
    d.code_system,
    d.code,
    d.title,
    d.description,
    d.status,
    d.severity,
    d.diagnosed_on,
    d.resolved_on,
    d.diagnosing_clinician,
    d.notes,
    d.created_by_employee_id,
    d.updated_by_employee_id,
    d.created_at,
    d.updated_at,
    d.archived_at
FROM client_diagnosis d
WHERE d.client_id = $1
  AND d.id = $2
  AND d.archived_at IS NULL
LIMIT 1;


-- name: UpdateClientDiagnosis :one
UPDATE client_diagnosis
SET
    code_system = COALESCE(sqlc.narg('code_system'), code_system),
    code = COALESCE(sqlc.narg('code'), code),
    title = COALESCE(sqlc.narg('title'), title),
    description = COALESCE(sqlc.narg('description'), description),
    status = COALESCE(sqlc.narg('status'), status),
    severity = COALESCE(sqlc.narg('severity'), severity),
    diagnosed_on = COALESCE(sqlc.narg('diagnosed_on'), diagnosed_on),
    resolved_on = COALESCE(sqlc.narg('resolved_on'), resolved_on),
    diagnosing_clinician = COALESCE(sqlc.narg('diagnosing_clinician'), diagnosing_clinician),
    notes = COALESCE(sqlc.narg('notes'), notes),
    updated_by_employee_id = COALESCE(sqlc.narg('updated_by_employee_id'), updated_by_employee_id)
WHERE client_id = sqlc.arg('client_id')
  AND id = sqlc.arg('id')
  AND archived_at IS NULL
RETURNING
    id,
    client_id,
    code_system,
    code,
    title,
    description,
    status,
    severity,
    diagnosed_on,
    resolved_on,
    diagnosing_clinician,
    notes,
    created_by_employee_id,
    updated_by_employee_id,
    created_at,
    updated_at,
    archived_at;


-- name: DeleteClientDiagnosis :one
DELETE FROM client_diagnosis
WHERE client_id = $1
  AND id = $2
RETURNING
    id,
    client_id,
    code_system,
    code,
    title,
    description,
    status,
    severity,
    diagnosed_on,
    resolved_on,
    diagnosing_clinician,
    notes,
    created_by_employee_id,
    updated_by_employee_id,
    created_at,
    updated_at,
    archived_at;


-- name: CreateClientMedicationOrder :one
INSERT INTO client_medication_order (
    client_id,
    diagnosis_id,
    medication_name,
    dosage_text,
    dose_amount,
    dose_unit,
    route,
    frequency_text,
    schedule,
    is_prn,
    prn_indication,
    max_doses_per_24h,
    start_date,
    end_date,
    status,
    admin_mode,
    responsible_employee_id,
    is_critical,
    notes,
    source_attachment_uuid,
    created_by_employee_id,
    updated_by_employee_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9,
    $10, $11, $12, $13, $14, $15, $16, $17,
    $18, $19, $20, $21, $22
) RETURNING
    id,
    client_id,
    diagnosis_id,
    medication_name,
    dosage_text,
    dose_amount,
    dose_unit,
    route,
    frequency_text,
    schedule,
    is_prn,
    prn_indication,
    max_doses_per_24h,
    start_date,
    end_date,
    status,
    admin_mode,
    responsible_employee_id,
    is_critical,
    notes,
    source_attachment_uuid,
    created_by_employee_id,
    updated_by_employee_id,
    created_at,
    updated_at,
    archived_at;


-- name: ListClientMedicationOrders :many
SELECT
    mo.id,
    mo.client_id,
    mo.diagnosis_id,
    mo.medication_name,
    mo.dosage_text,
    mo.dose_amount,
    mo.dose_unit,
    mo.route,
    mo.frequency_text,
    mo.schedule,
    mo.is_prn,
    mo.prn_indication,
    mo.max_doses_per_24h,
    mo.start_date,
    mo.end_date,
    mo.status,
    mo.admin_mode,
    mo.responsible_employee_id,
    mo.is_critical,
    mo.notes,
    mo.source_attachment_uuid,
    mo.created_by_employee_id,
    mo.updated_by_employee_id,
    mo.created_at,
    mo.updated_at,
    mo.archived_at,
    e.first_name AS responsible_employee_first_name,
    e.last_name AS responsible_employee_last_name,
    d.title AS diagnosis_title,
    d.code_system AS diagnosis_code_system,
    d.code AS diagnosis_code,
    COUNT(*) OVER() AS total_count
FROM client_medication_order mo
LEFT JOIN employee_profile e ON e.id = mo.responsible_employee_id
LEFT JOIN client_diagnosis d ON d.id = mo.diagnosis_id AND d.client_id = mo.client_id
WHERE mo.client_id = sqlc.arg('client_id')
  AND mo.archived_at IS NULL
  AND (sqlc.narg('status')::medication_order_status_enum IS NULL OR mo.status = sqlc.narg('status')::medication_order_status_enum)
  AND (sqlc.narg('admin_mode')::medication_admin_mode_enum IS NULL OR mo.admin_mode = sqlc.narg('admin_mode')::medication_admin_mode_enum)
  AND (sqlc.narg('diagnosis_id')::uuid IS NULL OR mo.diagnosis_id = sqlc.narg('diagnosis_id')::uuid)
  AND (
    sqlc.narg('search')::text IS NULL
    OR mo.medication_name ILIKE '%' || sqlc.narg('search')::text || '%'
    OR mo.dosage_text ILIKE '%' || sqlc.narg('search')::text || '%'
  )
ORDER BY mo.start_date DESC, mo.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');


-- name: GetClientMedicationOrder :one
SELECT
    mo.id,
    mo.client_id,
    mo.diagnosis_id,
    mo.medication_name,
    mo.dosage_text,
    mo.dose_amount,
    mo.dose_unit,
    mo.route,
    mo.frequency_text,
    mo.schedule,
    mo.is_prn,
    mo.prn_indication,
    mo.max_doses_per_24h,
    mo.start_date,
    mo.end_date,
    mo.status,
    mo.admin_mode,
    mo.responsible_employee_id,
    mo.is_critical,
    mo.notes,
    mo.source_attachment_uuid,
    mo.created_by_employee_id,
    mo.updated_by_employee_id,
    mo.created_at,
    mo.updated_at,
    mo.archived_at,
    e.first_name AS responsible_employee_first_name,
    e.last_name AS responsible_employee_last_name,
    d.title AS diagnosis_title,
    d.code_system AS diagnosis_code_system,
    d.code AS diagnosis_code
FROM client_medication_order mo
LEFT JOIN employee_profile e ON e.id = mo.responsible_employee_id
LEFT JOIN client_diagnosis d ON d.id = mo.diagnosis_id AND d.client_id = mo.client_id
WHERE mo.client_id = $1
  AND mo.id = $2
  AND mo.archived_at IS NULL
LIMIT 1;


-- name: UpdateClientMedicationOrder :one
UPDATE client_medication_order
SET
    diagnosis_id = COALESCE(sqlc.narg('diagnosis_id'), diagnosis_id),
    medication_name = COALESCE(sqlc.narg('medication_name'), medication_name),
    dosage_text = COALESCE(sqlc.narg('dosage_text'), dosage_text),
    dose_amount = COALESCE(sqlc.narg('dose_amount'), dose_amount),
    dose_unit = COALESCE(sqlc.narg('dose_unit'), dose_unit),
    route = COALESCE(sqlc.narg('route'), route),
    frequency_text = COALESCE(sqlc.narg('frequency_text'), frequency_text),
    schedule = COALESCE(sqlc.narg('schedule'), schedule),
    is_prn = COALESCE(sqlc.narg('is_prn'), is_prn),
    prn_indication = COALESCE(sqlc.narg('prn_indication'), prn_indication),
    max_doses_per_24h = COALESCE(sqlc.narg('max_doses_per_24h'), max_doses_per_24h),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = COALESCE(sqlc.narg('end_date'), end_date),
    status = COALESCE(sqlc.narg('status'), status),
    admin_mode = COALESCE(sqlc.narg('admin_mode'), admin_mode),
    responsible_employee_id = COALESCE(sqlc.narg('responsible_employee_id'), responsible_employee_id),
    is_critical = COALESCE(sqlc.narg('is_critical'), is_critical),
    notes = COALESCE(sqlc.narg('notes'), notes),
    source_attachment_uuid = COALESCE(sqlc.narg('source_attachment_uuid'), source_attachment_uuid),
    updated_by_employee_id = COALESCE(sqlc.narg('updated_by_employee_id'), updated_by_employee_id)
WHERE client_id = sqlc.arg('client_id')
  AND id = sqlc.arg('id')
  AND archived_at IS NULL
RETURNING
    id,
    client_id,
    diagnosis_id,
    medication_name,
    dosage_text,
    dose_amount,
    dose_unit,
    route,
    frequency_text,
    schedule,
    is_prn,
    prn_indication,
    max_doses_per_24h,
    start_date,
    end_date,
    status,
    admin_mode,
    responsible_employee_id,
    is_critical,
    notes,
    source_attachment_uuid,
    created_by_employee_id,
    updated_by_employee_id,
    created_at,
    updated_at,
    archived_at;


-- name: DeleteClientMedicationOrder :exec
DELETE FROM client_medication_order
WHERE client_id = $1
  AND id = $2;
