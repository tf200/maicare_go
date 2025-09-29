-- name: CreateClientDetails :one
INSERT INTO client_details (
    intake_form_id,
    first_name,
    last_name,
    date_of_birth,
    "identity",
    bsn,
    bsn_verified_by,
    source,
    birthplace,
    email,
    phone_number,
    organisation,
    departement,
    gender,
    filenumber,
    profile_picture,
    infix,
    sender_id,
    location_id,
    departure_reason,
    departure_report,
    addresses,
    legal_measure,
    education_currently_enrolled,
    education_institution,
    education_mentor_name,
    education_mentor_phone,
    education_mentor_email,
    education_additional_notes,
    education_level, 
    work_currently_employed,
    work_current_employer,
    work_current_employer_phone,
    work_current_employer_email,
    work_current_position,
    work_start_date,
    work_additional_notes, 
    living_situation,
    living_situation_notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, 
    $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36, $37, $38, $39
) RETURNING *;


-- name: ListClientDetails :many
SELECT 
    *,
    location.name AS location_name,
    COUNT(*) OVER() AS total_count
FROM client_details
LEFT JOIN location ON client_details.location_id = location.id
WHERE
    (status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
    (location_id = sqlc.narg('location_id') OR sqlc.narg('location_id') IS NULL) AND
    (sqlc.narg('search')::TEXT IS NULL OR 
        first_name ILIKE '%' || sqlc.narg('search') || '%' OR
        last_name ILIKE '%' || sqlc.narg('search') || '%' OR
        filenumber ILIKE '%' || sqlc.narg('search') || '%' OR
        email ILIKE '%' || sqlc.narg('search') || '%' OR
        phone_number ILIKE '%' || sqlc.narg('search') || '%')
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetClientCounts :one
SELECT 
    COUNT(*) AS total_clients,
    COUNT(*) FILTER (WHERE status = 'In Care') AS clients_in_care,
    COUNT(*) FILTER (WHERE status = 'On Waiting List') AS clients_on_waiting_list,
    COUNT(*) FILTER (WHERE status = 'Out Of Care') AS clients_out_of_care
FROM client_details;


-- name: GetAllClientsIDs :many
SELECT id FROM client_details;




-- name: GetClientDetails :one
SELECT c.*,
       ep.first_name AS bsn_verified_by_first_name,
       ep.last_name AS bsn_verified_by_last_name
FROM client_details c
LEFT JOIN employee_profile ep ON c.bsn_verified_by = ep.id
WHERE c.id = $1 LIMIT 1;

-- name: GetClientAddresses :one
SELECT addresses
FROM client_details
WHERE id = $1 LIMIT 1;


-- name: UpdateClientDetails :one
UPDATE client_details
SET 
    first_name = COALESCE (sqlc.narg('first_name'), first_name),
    last_name = COALESCE (sqlc.narg('last_name'), last_name),
    date_of_birth = COALESCE (sqlc.narg('date_of_birth'), date_of_birth),
    "identity" = COALESCE (sqlc.narg('identity'), "identity"),
    bsn = COALESCE (sqlc.narg('bsn'), bsn),
    bsn_verified_by = COALESCE (sqlc.narg('bsn_verified_by'), bsn_verified_by),
    source = COALESCE (sqlc.narg('source'), source),
    birthplace = COALESCE (sqlc.narg('birthplace'), birthplace),
    email = COALESCE (sqlc.narg('email'), email),
    phone_number = COALESCE (sqlc.narg('phone_number'), phone_number),
    organisation = COALESCE (sqlc.narg('organisation'), organisation),
    departement = COALESCE (sqlc.narg('departement'), departement),
    gender = COALESCE (sqlc.narg('gender'), gender),
    filenumber = COALESCE (sqlc.narg('filenumber'), filenumber),
    profile_picture = COALESCE (sqlc.narg('profile_picture'), profile_picture),
    infix = COALESCE (sqlc.narg('infix'), infix),
    sender_id = COALESCE (sqlc.narg('sender_id'), sender_id),
    location_id = COALESCE (sqlc.narg('location_id'), location_id),
    departure_reason = COALESCE (sqlc.narg('departure_reason'), departure_reason),
    departure_report = COALESCE (sqlc.narg('departure_report'), departure_report),
    legal_measure = COALESCE (sqlc.narg('legal_measure'), legal_measure),
    education_currently_enrolled = COALESCE (sqlc.narg('education_currently_enrolled'), education_currently_enrolled),
    education_institution = COALESCE (sqlc.narg('education_institution'), education_institution),
    education_mentor_name = COALESCE (sqlc.narg('education_mentor_name'), education_mentor_name),
    education_mentor_phone = COALESCE (sqlc.narg('education_mentor_phone'), education_mentor_phone),
    education_mentor_email = COALESCE (sqlc.narg('education_mentor_email'), education_mentor_email),
    education_additional_notes = COALESCE (sqlc.narg('education_additional_notes'), education_additional_notes),
    education_level = COALESCE (sqlc.narg('education_level'), education_level),
    work_currently_employed = COALESCE (sqlc.narg('work_currently_employed'), work_currently_employed),
    work_current_employer = COALESCE (sqlc.narg('work_current_employer'), work_current_employer),
    work_current_employer_phone = COALESCE (sqlc.narg('work_current_employer_phone'), work_current_employer_phone),
    work_current_employer_email = COALESCE (sqlc.narg('work_current_employer_email'), work_current_employer_email),
    work_current_position = COALESCE (sqlc.narg('work_current_position'), work_current_position),
    work_start_date = COALESCE (sqlc.narg('work_start_date'), work_start_date),
    work_additional_notes = COALESCE (sqlc.narg('work_additional_notes'), work_additional_notes),
    living_situation = COALESCE (sqlc.narg('living_situation'), living_situation),
    living_situation_notes = COALESCE (sqlc.narg('living_situation_notes'), living_situation_notes)

WHERE id = $1
RETURNING *;

-- name: UpdateClientStatus :one
UPDATE client_details
SET status = $2
WHERE id = $1
RETURNING *;

-- name: CreateClientStatusHistory :one
INSERT INTO client_status_history (
    client_id,
    old_status,
    new_status,
    reason
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: ListClientStatusHistory :many
SELECT * FROM client_status_history
WHERE client_id = $1
ORDER BY changed_at DESC
LIMIT $2 OFFSET $3;

-- name: CreateSchedueledClientStatusChange :one
INSERT INTO scheduled_status_changes (
    client_id,
    new_status,
    reason,
    scheduled_date
) VALUES (
    $1, $2, $3, $4
) RETURNING *;


-- name: SetClientProfilePicture :one
UPDATE client_details
SET profile_picture = $2
WHERE id = $1
RETURNING *;


-- name: CreateClientDocument :one
INSERT INTO client_documents (
    client_id,
    attachment_uuid,
    label
) VALUES (
    $1, $2, $3
) RETURNING *;


-- name: ListClientDocuments :many
SELECT 
    cd.*,
    a.*,
    COUNT(*) OVER() AS total_count
FROM client_documents cd
JOIN attachment_file a ON cd.attachment_uuid = a.uuid
WHERE client_id = $1
LIMIT $2 OFFSET $3;


-- name: DeleteClientDocument :one
DELETE FROM client_documents
WHERE attachment_uuid = $1
RETURNING *;


-- name: GetMissingClientDocuments :many
WITH all_labels AS (
    SELECT unnest(ARRAY[
        'registration_form', 'intake_form', 'consent_form',
        'risk_assessment', 'self_reliance_matrix', 'force_inventory',
        'care_plan', 'signaling_plan', 'cooperation_agreement'
    ]) AS label
),
client_labels AS (
    SELECT label
    FROM client_documents
    WHERE client_id = $1
)
SELECT al.label::text AS missing_label
FROM all_labels al
LEFT JOIN client_labels cl ON al.label = cl.label
WHERE cl.label IS NULL;
