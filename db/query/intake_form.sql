
-- name: CreateIntakeForm :one
INSERT INTO intake_forms (
    first_name,
    last_name,
    date_of_birth,
    nationality,
    bsn,
    address,
    city,
    postal_code,
    phone_number,
    gender,
    email,
    id_type,
    id_number,
    referrer_name,
    referrer_organization,
    referrer_function,
    referrer_phone,
    referrer_email,
    signed_by,
    has_valid_indication,
    law_type,
    other_law_specification,
    main_provider_name,
    main_provider_contact,
    indication_start_date,
    indication_end_date,
    registration_reason,
    guidance_goals,
    registration_type,
    living_situation,
    other_living_situation,
    parental_authority,
    current_school,
    mentor_name,
    mentor_phone,
    mentor_email,
    previous_care,
    guardian_details,
    diagnoses,
    uses_medication,
    medication_details,
    addiction_issues,
    judicial_involvement,
    risk_aggression,
    risk_suicidality,
    risk_running_away,
    risk_self_harm,
    risk_weapon_possession,
    risk_drug_dealing,
    other_risks,
    sharing_permission,
    truth_declaration,
    client_signature,
    guardian_signature,
    referrer_signature,
    signature_date,
    attachement_ids
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
    $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
    $31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
    $41, $42, $43, $44, $45, $46, $47, $48, $49, $50,
    $51, $52, $53, $54, $55, $56, $57
) RETURNING *;




-- name: ListIntakeForms :many
SELECT *, COUNT(*) OVER() AS total_count FROM intake_forms
WHERE (
    LOWER(first_name) LIKE LOWER(CONCAT('%', COALESCE(@search::text, ''), '%')) OR
    LOWER(last_name) LIKE LOWER(CONCAT('%', COALESCE(@search::text, ''), '%'))
)
ORDER BY
    CASE 
        WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'asc' THEN created_at
    END ASC,
    CASE 
        WHEN @sort_by::text = 'created_at' AND @sort_order::text = 'desc' THEN created_at
    END DESC,
    CASE 
        WHEN @sort_by::text = 'urgency_score' AND @sort_order::text = 'asc' THEN urgency_score
    END ASC,
    CASE 
        WHEN @sort_by::text = 'urgency_score' AND @sort_order::text = 'desc' THEN urgency_score
    END DESC,
    CASE
        WHEN @sort_by::text IS NULL OR @sort_by::text = '' THEN id
    END DESC
LIMIT $1 OFFSET $2;


-- name: GetIntakeForm :one
SELECT * FROM intake_forms
WHERE id = $1;


-- name: AddUrgencyScore :one
UPDATE intake_forms
SET urgency_score = $2
WHERE id = $1
RETURNING *;