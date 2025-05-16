-- name: CreateRegistrationForm :one
INSERT INTO registration_form (
    client_first_name,
    client_last_name,
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
    care_protected_living,
    care_assisted_independent_living,
    care_room_training_center,
    care_ambulatory_guidance,
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
    referrer_signature
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19,
    $20, $21, $22, $23, $24, $25, $26, $27, $28,
    $29, $30, $31, $32, $33, $34, $35, $36, $37,
    $38, $39, $40, $41, $42, $43, $44, $45, $46,
    $47, $48, $49, $50, $51, $52, $53, $54, $55,
    $56, $57
) RETURNING *;




-- name: ListRegistrationForms :many
SELECT * FROM registration_form
WHERE 
    -- Form status filtering
    (sqlc.narg('status')::VARCHAR IS NULL OR form_status = sqlc.narg('status'))
    -- Risk filtering
    AND (sqlc.narg('risk_aggressive_behavior')::BOOLEAN IS NULL OR risk_aggressive_behavior = sqlc.narg('risk_aggressive_behavior'))
    AND (sqlc.narg('risk_suicidal_selfharm')::BOOLEAN IS NULL OR risk_suicidal_selfharm = sqlc.narg('risk_suicidal_selfharm'))
    AND (sqlc.narg('risk_substance_abuse')::BOOLEAN IS NULL OR risk_substance_abuse = sqlc.narg('risk_substance_abuse'))
    AND (sqlc.narg('risk_psychiatric_issues')::BOOLEAN IS NULL OR risk_psychiatric_issues = sqlc.narg('risk_psychiatric_issues'))
    AND (sqlc.narg('risk_criminal_history')::BOOLEAN IS NULL OR risk_criminal_history = sqlc.narg('risk_criminal_history'))
    AND (sqlc.narg('risk_flight_behavior')::BOOLEAN IS NULL OR risk_flight_behavior = sqlc.narg('risk_flight_behavior'))
    AND (sqlc.narg('risk_weapon_possession')::BOOLEAN IS NULL OR risk_weapon_possession = sqlc.narg('risk_weapon_possession'))
    AND (sqlc.narg('risk_sexual_behavior')::BOOLEAN IS NULL OR risk_sexual_behavior = sqlc.narg('risk_sexual_behavior'))
    AND (sqlc.narg('risk_day_night_rhythm')::BOOLEAN IS NULL OR risk_day_night_rhythm = sqlc.narg('risk_day_night_rhythm'))
    AND (sqlc.narg('risk_other')::BOOLEAN IS NULL OR risk_other = sqlc.narg('risk_other'))
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: CountRegistrationForms :one
SELECT COUNT(*) FROM registration_form
WHERE 
    -- Form status filtering
    (sqlc.narg('status')::VARCHAR IS NULL OR form_status = sqlc.narg('status'))
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





-- name: GetRegistrationForm :one
SELECT * FROM registration_form
WHERE id = $1
LIMIT 1;

-- name: UpdateRegistrationForm :one
UPDATE registration_form
SET
    client_first_name = COALESCE(sqlc.narg('client_first_name'), client_first_name),
    client_last_name = COALESCE(sqlc.narg('client_last_name'), client_last_name),
    client_bsn_number = COALESCE(sqlc.narg('client_bsn_number'), client_bsn_number),
    client_gender = COALESCE(sqlc.narg('client_gender'), client_gender),
    client_nationality = COALESCE(sqlc.narg('client_nationality'), client_nationality),
    client_phone_number = COALESCE(sqlc.narg('client_phone_number'), client_phone_number),
    client_email = COALESCE(sqlc.narg('client_email'), client_email),
    client_street = COALESCE(sqlc.narg('client_street'), client_street),
    client_house_number = COALESCE(sqlc.narg('client_house_number'), client_house_number),
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
    care_protected_living = COALESCE(sqlc.narg('care_protected_living'), care_protected_living),
    care_assisted_independent_living = COALESCE(sqlc.narg('care_assisted_independent_living'), care_assisted_independent_living),
    care_room_training_center = COALESCE(sqlc.narg('care_room_training_center'), care_room_training_center),
    care_ambulatory_guidance = COALESCE(sqlc.narg('care_ambulatory_guidance'), care_ambulatory_guidance),
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
    document_referral = COALESCE(sqlc.narg('document_referral'), document_referral),
    document_education_report = COALESCE(sqlc.narg('document_education_report'), document_education_report),
    document_psychiatric_report = COALESCE(sqlc.narg('document_psychiatric_report'), document_psychiatric_report),
    document_diagnosis = COALESCE(sqlc.narg('document_diagnosis'), document_diagnosis),
    document_safety_plan = COALESCE(sqlc.narg('document_safety_plan'), document_safety_plan),
    document_id_copy = COALESCE(sqlc.narg('document_id_copy'), document_id_copy),
    application_date = COALESCE(sqlc.narg('application_date'), application_date),
    referrer_signature = COALESCE(sqlc.narg('referrer_signature'), referrer_signature)
WHERE id = sqlc.arg('id')
RETURNING *;


-- name: DeleteRegistrationForm :exec
DELETE FROM registration_form
WHERE id = $1;


-- name: UpdateRegistrationFormStatus :exec
UPDATE registration_form
SET
    form_status = $2,
    updated_at = NOW(),
    processed_by_employee_id = $3
WHERE id = $1;