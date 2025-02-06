-- name: CreateClientAllergy :one
INSERT INTO client_allergy (
    client_id,
    allergy_type_id,
    severity,
    reaction,
    notes,
    created_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;



-- name: ListClientAllergies :many
SELECT 
    al.*,
    ty.name AS allergy_type,
    (SELECT COUNT(*) FROM client_allergy WHERE client_allergy.client_id = al.client_id) AS total_allergies
FROM client_allergy al
JOIN allergy_type ty ON al.allergy_type_id = ty.id
WHERE al.client_id = $1
LIMIT $2 OFFSET $3;


-- name: GetClientAllergy :one
SELECT * FROM client_allergy
WHERE id = $1 LIMIT 1;

-- name: UpdateClientAllergy :one
UPDATE client_allergy
SET
    allergy_type_id = COALESCE(sqlc.narg('allergy_type_id'), allergy_type_id),
    severity = COALESCE(sqlc.narg('severity'), severity),
    reaction = COALESCE(sqlc.narg('reaction'), reaction),
    notes = COALESCE(sqlc.narg('notes'), notes)
WHERE id = $1
RETURNING *;


-- name: DeleteClientAllergy :one
DELETE FROM client_allergy
WHERE id = $1
RETURNING *;


-- name: CreateClientDiagnosis :one
INSERT INTO client_diagnosis (
    client_id,
    title,
    diagnosis_code,
    description,
    date_of_diagnosis,
    severity,
    status,
    diagnosing_clinician,
    notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
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
WHERE id = $1 LIMIT 1;


-- name: UpdateClientDiagnosis :one
UPDATE client_diagnosis
SET
    title = COALESCE(sqlc.narg('title'), title),
    diagnosis_code = COALESCE(sqlc.narg('diagnosis_code'), diagnosis_code),
    description = COALESCE(sqlc.narg('description'), description),
    date_of_diagnosis = COALESCE(sqlc.narg('date_of_diagnosis'), date_of_diagnosis),
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
