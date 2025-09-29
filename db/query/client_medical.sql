-- name: CreateClientDiagnosis :one
INSERT INTO client_diagnosis (
    client_id,
    title,
    diagnosis_code,
    description,
    severity,
    status,
    diagnosing_clinician,
    notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: ListClientDiagnoses :many
SELECT 
    d.*,
    (SELECT COUNT(*) FROM client_diagnosis WHERE client_diagnosis.client_id = d.client_id) AS total_diagnoses
FROM client_diagnosis d
WHERE d.client_id = $1
LIMIT $2 OFFSET $3;

-- name: GetClientDiagnosis :one
SELECT * FROM client_diagnosis
WHERE id = $1
LIMIT 1; 



-- name: UpdateClientDiagnosis :one
UPDATE client_diagnosis
SET
    title = COALESCE(sqlc.narg('title'), title),
    diagnosis_code = COALESCE(sqlc.narg('diagnosis_code'), diagnosis_code),
    description = COALESCE(sqlc.narg('description'), description),
    severity = COALESCE(sqlc.narg('severity'), severity),
    status = COALESCE(sqlc.narg('status'), status),
    diagnosing_clinician = COALESCE(sqlc.narg('diagnosing_clinician'), diagnosing_clinician),
    notes = COALESCE(sqlc.narg('notes'), notes)
WHERE id = $1
RETURNING *;

-- name: DeleteClientDiagnosis :one
DELETE FROM client_diagnosis
WHERE id = $1  
RETURNING *;



-- name: CreateClientMedication :one
INSERT INTO client_medication (
    diagnosis_id,
    name,
    dosage,
    start_date,
    end_date,
    notes,
    self_administered,
    administered_by_id,
    is_critical
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;


-- name: UpdateClientMedication :one
UPDATE client_medication
SET
    name = COALESCE(sqlc.narg('name'), name),
    dosage = COALESCE(sqlc.narg('dosage'), dosage),
    start_date = COALESCE(sqlc.narg('start_date'), start_date),
    end_date = COALESCE(sqlc.narg('end_date'), end_date),
    notes = COALESCE(sqlc.narg('notes'), notes),
    self_administered = COALESCE(sqlc.narg('self_administered'), self_administered),
    administered_by_id = COALESCE(sqlc.narg('administered_by_id'), administered_by_id),
    is_critical = COALESCE(sqlc.narg('is_critical'), is_critical)
WHERE id = $1
RETURNING *;

-- name: DeleteClientMedication :exec
DELETE FROM client_medication
WHERE id = $1;

-- name: ListMedicationsByDiagnosisID :many
SELECT 
    m.*,
    (SELECT COUNT(*) FROM client_medication WHERE client_medication.diagnosis_id = $1) AS total_medications
FROM client_medication m
WHERE m.diagnosis_id = $1
ORDER BY m.id
LIMIT $2 OFFSET $3;

-- name: ListMedicationsByDiagnosisIDs :many
SELECT *
FROM client_medication
WHERE diagnosis_id = ANY($1::bigint[]);


-- name: GetMedication :one
SELECT m.*, e.first_name AS administered_by_first_name, e.last_name AS administered_by_last_name
FROM client_medication m
JOIN employee_profile e ON m.administered_by_id = e.id
WHERE m.id = $1 LIMIT 1;



-- -- name: UpdateClientMedication :one
-- UPDATE client_medication
-- SET
--     name = COALESCE(sqlc.narg('name'), name),
--     dosage = COALESCE(sqlc.narg('dosage'), dosage),
--     start_date = COALESCE(sqlc.narg('start_date'), start_date),
--     end_date = COALESCE(sqlc.narg('end_date'), end_date),
--     notes = COALESCE(sqlc.narg('notes'), notes),
--     self_administered = COALESCE(sqlc.narg('self_administered'), self_administered),
--     administered_by_id = COALESCE(sqlc.narg('administered_by_id'), administered_by_id),
--     is_critical = COALESCE(sqlc.narg('is_critical'), is_critical)
-- WHERE id = $1
-- RETURNING *;


-- -- name: DeleteClientMedication :one
-- DELETE FROM client_medication
-- WHERE id = $1
-- RETURNING *;


