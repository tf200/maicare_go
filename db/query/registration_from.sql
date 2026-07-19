-- name: CreateRegistrationForm :one
INSERT INTO registration_form (
    client_first_name,
    client_last_name,
    client_date_of_birth,
    client_bsn_number,
    client_gender,
    client_nationality,
    client_phone_number,
    client_email,
    client_street,
    client_house_number,
    client_postal_code,
    client_city,
    referrer_first_name,
    referrer_last_name,
    referrer_organization,
    referrer_job_title,
    referrer_phone_number,
    referrer_email,
    guardian1_first_name,
    guardian1_last_name,
    guardian1_relationship,
    guardian1_phone_number,
    guardian1_email,
    guardian2_first_name,
    guardian2_last_name,
    guardian2_relationship,
    guardian2_phone_number,
    guardian2_email,
    education_institution,
    education_mentor_name,
    education_mentor_phone,
    education_mentor_email,
    education_currently_enrolled,
    education_additional_notes,
    education_level,
    work_current_employer,
    work_employer_phone,
    work_employer_email,
    work_current_position,
    work_currently_employed,
    work_start_date,
    work_additional_notes,
    care_protected_living,
    care_assisted_independent_living,
    care_room_training_center,
    care_ambulatory_guidance,
    application_reason,
    client_goals,
    risk_aggressive_behavior,
    risk_suicidal_selfharm,
    risk_substance_abuse,
    risk_psychiatric_issues,
    risk_criminal_history,
    risk_flight_behavior,
    risk_weapon_possession,
    risk_sexual_behavior,
    risk_day_night_rhythm,
    risk_other,
    risk_other_description,
    risk_additional_notes,
    document_referral,
    document_education_report,
    document_psychiatric_report,
    document_diagnosis,
    document_safety_plan,
    document_id_copy,
    application_date,
    referrer_signature,
    client_house_number_addition
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19,
    $20, $21, $22, $23, $24, $25, $26, $27, $28,
    $29, $30, $31, $32, $33, $34, $35, $36, $37,
    $38, $39, $40, $41, $42, $43, $44, $45, $46,
    $47, $48, $49, $50, $51, $52, $53, $54, $55,
    $56, $57, $58, $59, $60, $61, $62, $63, $64, $65,
    $66, $67, $68, $69
) RETURNING *;




-- name: ListRegistrationForms :many
SELECT
    rf.*,
    iform.id AS intake_form_id
FROM registration_form rf
LEFT JOIN intake_forms iform ON rf.id = iform.registration_form_id
WHERE
    -- Form status filtering
    (sqlc.narg('status')::form_status_enum IS NULL OR rf.form_status = sqlc.narg('status')::form_status_enum)
    -- Risk filtering
    AND (sqlc.narg('risk_aggressive_behavior')::BOOLEAN IS NULL OR rf.risk_aggressive_behavior = sqlc.narg('risk_aggressive_behavior'))
    AND (sqlc.narg('risk_suicidal_selfharm')::BOOLEAN IS NULL OR rf.risk_suicidal_selfharm = sqlc.narg('risk_suicidal_selfharm'))
    AND (sqlc.narg('risk_substance_abuse')::BOOLEAN IS NULL OR rf.risk_substance_abuse = sqlc.narg('risk_substance_abuse'))
    AND (sqlc.narg('risk_psychiatric_issues')::BOOLEAN IS NULL OR rf.risk_psychiatric_issues = sqlc.narg('risk_psychiatric_issues'))
    AND (sqlc.narg('risk_criminal_history')::BOOLEAN IS NULL OR rf.risk_criminal_history = sqlc.narg('risk_criminal_history'))
    AND (sqlc.narg('risk_flight_behavior')::BOOLEAN IS NULL OR rf.risk_flight_behavior = sqlc.narg('risk_flight_behavior'))
    AND (sqlc.narg('risk_weapon_possession')::BOOLEAN IS NULL OR rf.risk_weapon_possession = sqlc.narg('risk_weapon_possession'))
    AND (sqlc.narg('risk_sexual_behavior')::BOOLEAN IS NULL OR rf.risk_sexual_behavior = sqlc.narg('risk_sexual_behavior'))
    AND (sqlc.narg('risk_day_night_rhythm')::BOOLEAN IS NULL OR rf.risk_day_night_rhythm = sqlc.narg('risk_day_night_rhythm'))
    AND (sqlc.narg('risk_other')::BOOLEAN IS NULL OR rf.risk_other = sqlc.narg('risk_other'))
ORDER BY rf.created_at DESC
LIMIT $1 OFFSET $2;


-- name: CountRegistrationForms :one
SELECT COUNT(*) FROM registration_form
WHERE
    -- Form status filtering
    (sqlc.narg('status')::form_status_enum IS NULL OR form_status = sqlc.narg('status')::form_status_enum)
    -- Risk filtering
    AND (sqlc.narg('risk_aggressive_behavior')::BOOLEAN IS NULL OR risk_aggressive_behavior = sqlc.narg('risk_aggressive_behavior'))
    AND (sqlc.narg('risk_suicidal_selfharm')::BOOLEAN IS NULL OR risk_suicidal_selfharm = sqlc.narg('risk_suicidal_selfharm'))
    AND (sqlc.narg('risk_substance_abuse')::BOOLEAN IS NULL OR risk_substance_abuse = sqlc.narg('risk_substance_abuse'))
    AND (sqlc.narg('risk_psychiatric_issues')::BOOLEAN IS NULL OR risk_psychiatric_issues = sqlc.narg('risk_psychiatric_issues'))
    AND (sqlc.narg('risk_criminal_history')::BOOLEAN IS NULL OR risk_criminal_history = sqlc.narg('risk_criminal_history'))
    AND (sqlc.narg('risk_flight_behavior')::BOOLEAN IS NULL OR risk_flight_behavior = sqlc.narg('risk_flight_behavior'))
    AND (sqlc.narg('risk_weapon_possession')::BOOLEAN IS NULL OR risk_weapon_possession = sqlc.narg('risk_weapon_possession'))
    AND (sqlc.narg('risk_sexual_behavior')::BOOLEAN IS NULL OR risk_sexual_behavior = sqlc.narg('risk_sexual_behavior'))
    AND (sqlc.narg('risk_day_night_rhythm')::BOOLEAN IS NULL OR risk_day_night_rhythm = sqlc.narg('risk_day_night_rhythm'));

-- name: GetRegistrationFormCounts :one
SELECT
    COUNT(*) AS total,
    COUNT(*) FILTER (WHERE form_status = 'pending') AS pending_review,
    COUNT(*) FILTER (WHERE form_status = 'processed') AS processed,
    COUNT(*) FILTER (WHERE
        (CASE WHEN risk_aggressive_behavior THEN 1 ELSE 0 END) +
        (CASE WHEN risk_suicidal_selfharm THEN 1 ELSE 0 END) +
        (CASE WHEN risk_substance_abuse THEN 1 ELSE 0 END) +
        (CASE WHEN risk_psychiatric_issues THEN 1 ELSE 0 END) +
        (CASE WHEN risk_criminal_history THEN 1 ELSE 0 END) +
        (CASE WHEN risk_flight_behavior THEN 1 ELSE 0 END) +
        (CASE WHEN risk_weapon_possession THEN 1 ELSE 0 END) +
        (CASE WHEN risk_sexual_behavior THEN 1 ELSE 0 END) +
        (CASE WHEN risk_day_night_rhythm THEN 1 ELSE 0 END) +
        (CASE WHEN risk_other THEN 1 ELSE 0 END) >= 3
    ) AS high_risk
FROM registration_form;





-- name: GetRegistrationForm :one
SELECT
    rf.*,
    ep.first_name AS processed_by_first_name,
    ep.last_name AS processed_by_last_name,
    iform.id AS intake_form_id
FROM registration_form rf
LEFT JOIN employee_profile ep ON rf.processed_by_employee_id = ep.id
LEFT JOIN intake_forms iform ON rf.id = iform.registration_form_id
WHERE rf.id = $1
LIMIT 1;

-- name: UpdateRegistrationForm :one
UPDATE registration_form
SET
    client_first_name = COALESCE(sqlc.narg('client_first_name'), client_first_name),
    client_last_name = COALESCE(sqlc.narg('client_last_name'), client_last_name),
    client_date_of_birth = COALESCE(sqlc.narg('client_date_of_birth'), client_date_of_birth),
    client_bsn_number = COALESCE(sqlc.narg('client_bsn_number'), client_bsn_number),
    client_gender = COALESCE(sqlc.narg('client_gender'), client_gender),
    client_nationality = COALESCE(sqlc.narg('client_nationality'), client_nationality),
    client_phone_number = COALESCE(sqlc.narg('client_phone_number'), client_phone_number),
    client_email = COALESCE(sqlc.narg('client_email'), client_email),
    client_street = COALESCE(sqlc.narg('client_street'), client_street),
    client_house_number = COALESCE(sqlc.narg('client_house_number'), client_house_number),
    client_house_number_addition = COALESCE(sqlc.narg('client_house_number_addition'), client_house_number_addition),
    client_postal_code = COALESCE(sqlc.narg('client_postal_code'), client_postal_code),
    client_city = COALESCE(sqlc.narg('client_city'), client_city),
    referrer_first_name = COALESCE(sqlc.narg('referrer_first_name'), referrer_first_name),
    referrer_last_name = COALESCE(sqlc.narg('referrer_last_name'), referrer_last_name),
    referrer_organization = COALESCE(sqlc.narg('referrer_organization'), referrer_organization),
    referrer_job_title = COALESCE(sqlc.narg('referrer_job_title'), referrer_job_title),
    referrer_phone_number = COALESCE(sqlc.narg('referrer_phone_number'), referrer_phone_number),
    referrer_email = COALESCE(sqlc.narg('referrer_email'), referrer_email),
    guardian1_first_name = COALESCE(sqlc.narg('guardian1_first_name'), guardian1_first_name),
    guardian1_last_name = COALESCE(sqlc.narg('guardian1_last_name'), guardian1_last_name),
    guardian1_relationship = COALESCE(sqlc.narg('guardian1_relationship'), guardian1_relationship),
    guardian1_phone_number = COALESCE(sqlc.narg('guardian1_phone_number'), guardian1_phone_number),
    guardian1_email = COALESCE(sqlc.narg('guardian1_email'), guardian1_email),
    guardian2_first_name = COALESCE(sqlc.narg('guardian2_first_name'), guardian2_first_name),
    guardian2_last_name = COALESCE(sqlc.narg('guardian2_last_name'), guardian2_last_name),
    guardian2_relationship = COALESCE(sqlc.narg('guardian2_relationship'), guardian2_relationship),
    guardian2_phone_number = COALESCE(sqlc.narg('guardian2_phone_number'), guardian2_phone_number),
    guardian2_email = COALESCE(sqlc.narg('guardian2_email'), guardian2_email),
    education_institution = COALESCE(sqlc.narg('education_institution'), education_institution),
    education_mentor_name = COALESCE(sqlc.narg('education_mentor_name'), education_mentor_name),
    education_mentor_phone = COALESCE(sqlc.narg('education_mentor_phone'), education_mentor_phone),
    education_mentor_email = COALESCE(sqlc.narg('education_mentor_email'), education_mentor_email),
    education_currently_enrolled = COALESCE(sqlc.narg('education_currently_enrolled'), education_currently_enrolled),
    education_additional_notes = COALESCE(sqlc.narg('education_additional_notes'), education_additional_notes),
    education_level = COALESCE(sqlc.narg('education_level'), education_level),
    work_current_employer = COALESCE(sqlc.narg('work_current_employer'), work_current_employer),
    work_employer_phone = COALESCE(sqlc.narg('work_employer_phone'), work_employer_phone),
    work_employer_email = COALESCE(sqlc.narg('work_employer_email'), work_employer_email),
    work_current_position = COALESCE(sqlc.narg('work_current_position'), work_current_position),
    work_currently_employed = COALESCE(sqlc.narg('work_currently_employed'), work_currently_employed),
    work_start_date = COALESCE(sqlc.narg('work_start_date'), work_start_date),
    work_additional_notes = COALESCE(sqlc.narg('work_additional_notes'), work_additional_notes),
    care_protected_living = COALESCE(sqlc.narg('care_protected_living'), care_protected_living),
    care_assisted_independent_living = COALESCE(sqlc.narg('care_assisted_independent_living'), care_assisted_independent_living),
    care_room_training_center = COALESCE(sqlc.narg('care_room_training_center'), care_room_training_center),
    care_ambulatory_guidance = COALESCE(sqlc.narg('care_ambulatory_guidance'), care_ambulatory_guidance),
    application_reason = COALESCE(sqlc.narg('application_reason'), application_reason),
    client_goals = COALESCE(sqlc.narg('client_goals'), client_goals),
    risk_aggressive_behavior = COALESCE(sqlc.narg('risk_aggressive_behavior'), risk_aggressive_behavior),
    risk_suicidal_selfharm = COALESCE(sqlc.narg('risk_suicidal_selfharm'), risk_suicidal_selfharm),
    risk_substance_abuse = COALESCE(sqlc.narg('risk_substance_abuse'), risk_substance_abuse),
    risk_psychiatric_issues = COALESCE(sqlc.narg('risk_psychiatric_issues'), risk_psychiatric_issues),
    risk_criminal_history = COALESCE(sqlc.narg('risk_criminal_history'), risk_criminal_history),
    risk_flight_behavior = COALESCE(sqlc.narg('risk_flight_behavior'), risk_flight_behavior),
    risk_weapon_possession = COALESCE(sqlc.narg('risk_weapon_possession'), risk_weapon_possession),
    risk_sexual_behavior = COALESCE(sqlc.narg('risk_sexual_behavior'), risk_sexual_behavior),
    risk_day_night_rhythm = COALESCE(sqlc.narg('risk_day_night_rhythm'), risk_day_night_rhythm),
    risk_other = COALESCE(sqlc.narg('risk_other'), risk_other),
    risk_other_description = COALESCE(sqlc.narg('risk_other_description'), risk_other_description),
    risk_additional_notes = COALESCE(sqlc.narg('risk_additional_notes'), risk_additional_notes),
    application_date = COALESCE(sqlc.narg('application_date'), application_date),
    referrer_signature = COALESCE(sqlc.narg('referrer_signature'), referrer_signature),
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;


-- name: DeleteRegistrationForm :exec
DELETE FROM registration_form
WHERE id = $1;

-- name: ReplaceRegistrationFormDocument :one
UPDATE registration_form
SET
    document_referral = CASE WHEN sqlc.arg('document_type') = 'document_referral' THEN sqlc.arg('file_id') ELSE document_referral END,
    document_education_report = CASE WHEN sqlc.arg('document_type') = 'document_education_report' THEN sqlc.arg('file_id') ELSE document_education_report END,
    document_action_plan = CASE WHEN sqlc.arg('document_type') = 'document_action_plan' THEN sqlc.arg('file_id') ELSE document_action_plan END,
    document_psychiatric_report = CASE WHEN sqlc.arg('document_type') = 'document_psychiatric_report' THEN sqlc.arg('file_id') ELSE document_psychiatric_report END,
    document_diagnosis = CASE WHEN sqlc.arg('document_type') = 'document_diagnosis' THEN sqlc.arg('file_id') ELSE document_diagnosis END,
    document_safety_plan = CASE WHEN sqlc.arg('document_type') = 'document_safety_plan' THEN sqlc.arg('file_id') ELSE document_safety_plan END,
    document_id_copy = CASE WHEN sqlc.arg('document_type') = 'document_id_copy' THEN sqlc.arg('file_id') ELSE document_id_copy END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = sqlc.arg('id')
RETURNING *;


-- name: UpdateRegistrationFormStatus :one
UPDATE registration_form
SET
    form_status = $2,
    processed_by_employee_id = $3,
    intake_appointment_location = COALESCE(sqlc.narg('intake_appointment_location'), intake_appointment_location),
    addmission_type = COALESCE(sqlc.narg('addmission_type'), addmission_type),
    intake_options = COALESCE(sqlc.narg('intake_options'), intake_options),
    intake_token = COALESCE(sqlc.narg('intake_token'), intake_token),
    rejection_reason = COALESCE(sqlc.narg('rejection_reason'), rejection_reason),
    processed_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;


-- name: GetRegistrationFormByToken :one
SELECT * FROM registration_form
WHERE intake_token = $1
LIMIT 1;

-- name: UpdateRegistrationFormIntakeDate :one
UPDATE registration_form
SET
    intake_appointment_datetime = $2,
    intake_token = NULL -- Invalidate token after use
WHERE id = $1
RETURNING *;
