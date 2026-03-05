-- name: CreateClientDetails :one
INSERT INTO client_details (
    intake_form_id,
    registration_form_id,
    first_name,
    last_name,
    date_of_birth,
    "identity",
    bsn,
    bsn_verified_by,
    email,
    phone_number,
    gender,
    care_type,
    sender_id,
    location_id,
    street,
    house_number,
    house_number_addition,
    postal_code,
    city,
    education_currently_enrolled,
    education_institution,
    education_mentor_name,
    education_mentor_phone,
    education_mentor_email,
    education_additional_notes,
    education_level,
    nationality,
    work_currently_employed,
    work_current_employer,
    work_current_employer_phone,
    work_current_employer_email,
    work_current_position,
    work_start_date,
    work_additional_notes,
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
    evaluation_intervals_weeks
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
    $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
    $31, $32, $33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44,
    $45, $46, $47
) RETURNING *;


-- name: GetClientByIntakeFormID :one
SELECT * FROM client_details
WHERE intake_form_id = $1
LIMIT 1;

-- name: ListClientNamesByIDs :many
SELECT
    id,
    first_name,
    last_name
FROM client_details
WHERE id = ANY(sqlc.arg(client_ids)::uuid[]);


-- name: ListClientDetails :many
SELECT
    c.id,
    c.first_name,
    c.last_name,
    c.bsn,
    c.filenumber,
    l.name AS location_name,
    c.care_type,
    c.status,
    COALESCE(gc.goals_count, 0)::bigint AS goals_count,
    (
        CASE WHEN COALESCE(c.risk_aggressive_behavior, FALSE) THEN 1 ELSE 0 END +
        CASE WHEN COALESCE(c.risk_suicidal_selfharm, FALSE) THEN 1 ELSE 0 END +
        CASE WHEN COALESCE(c.risk_substance_abuse, FALSE) THEN 1 ELSE 0 END +
        CASE WHEN COALESCE(c.risk_psychiatric_issues, FALSE) THEN 1 ELSE 0 END +
        CASE WHEN COALESCE(c.risk_criminal_history, FALSE) THEN 1 ELSE 0 END +
        CASE WHEN COALESCE(c.risk_flight_behavior, FALSE) THEN 1 ELSE 0 END +
        CASE WHEN COALESCE(c.risk_weapon_possession, FALSE) THEN 1 ELSE 0 END +
        CASE WHEN COALESCE(c.risk_sexual_behavior, FALSE) THEN 1 ELSE 0 END +
        CASE WHEN COALESCE(c.risk_day_night_rhythm, FALSE) THEN 1 ELSE 0 END +
        CASE WHEN COALESCE(c.risk_other, FALSE) THEN 1 ELSE 0 END
    )::bigint AS risk_count,
    c.created_at,
    COUNT(*) OVER() AS total_count
FROM client_details c
LEFT JOIN location l ON c.location_id = l.id
LEFT JOIN (
    SELECT
        client_id,
        COUNT(*) FILTER (WHERE status = 'active')::bigint AS goals_count
    FROM client_goals
    GROUP BY client_id
) gc ON gc.client_id = c.id
WHERE
    (c.status = sqlc.narg('status') OR sqlc.narg('status') IS NULL) AND
    (c.location_id = sqlc.narg('location_id') OR sqlc.narg('location_id') IS NULL) AND
    (sqlc.narg('search')::TEXT IS NULL OR
        c.first_name ILIKE '%' || sqlc.narg('search') || '%' OR
        c.last_name ILIKE '%' || sqlc.narg('search') || '%' OR
        c.filenumber ILIKE '%' || sqlc.narg('search') || '%' OR
        c.bsn ILIKE '%' || sqlc.narg('search') || '%')
ORDER BY c.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: GetClientCounts :one
SELECT
    COUNT(*) AS total_clients,
    COUNT(*) FILTER (WHERE status = 'in_care') AS clients_in_care,
    COUNT(*) FILTER (WHERE status = 'on_waiting_list') AS clients_on_waiting_list,
    COUNT(*) FILTER (WHERE status = 'out_of_care') AS clients_out_of_care
FROM client_details;

-- name: GetClientStatusCounts :one
SELECT
    COUNT(*) FILTER (WHERE status IN ('in_care', 'scheduled_in_care')) AS clients_in_or_scheduled_in_care,
    COUNT(*) FILTER (WHERE status = 'on_waiting_list') AS clients_on_waiting_list,
    COUNT(*) FILTER (WHERE status IN ('out_of_care', 'scheduled_out_of_care')) AS clients_out_or_scheduled_out_of_care
FROM client_details;


-- name: ListWaitingListClients :many
SELECT
    c.id,
    c.first_name,
    c.last_name,
    c.care_type,
    c.bsn,
    s.name AS sender_name,
    (CURRENT_DATE - c.created_at::date)::int4 AS days_in_waitlist,
    rf.addmission_type AS admission_type,
    COUNT(*) OVER() AS total_count
FROM client_details c
LEFT JOIN sender s ON c.sender_id = s.id
LEFT JOIN registration_form rf ON c.registration_form_id = rf.id
WHERE
    c.status = 'on_waiting_list'
    AND (
        sqlc.narg('search')::text IS NULL
        OR sqlc.narg('search')::text = ''
        OR c.first_name ILIKE '%' || sqlc.narg('search') || '%'
        OR c.last_name ILIKE '%' || sqlc.narg('search') || '%'
        OR s.name ILIKE '%' || sqlc.narg('search') || '%'
    )
    AND (
        sqlc.narg('placement')::intake_care_type_enum IS NULL
        OR c.care_type = sqlc.narg('placement')::intake_care_type_enum
    )
ORDER BY
    CASE WHEN sqlc.arg('sort_days')::text = 'asc' THEN (CURRENT_DATE - c.created_at::date) END ASC,
    CASE WHEN sqlc.arg('sort_days')::text = 'desc' THEN (CURRENT_DATE - c.created_at::date) END DESC,
    c.created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');


-- name: ListInCareClients :many
SELECT
    c.id,
    c.bsn,
    c.first_name,
    c.last_name,
    NULLIF(CONCAT_WS(' ', ep.first_name, ep.last_name), '')::text AS coordinator_name,
    l.name AS location_name,
    c.status,
    c.care_start_date,
    CASE
        WHEN c.care_start_date IS NULL OR CURRENT_DATE < c.care_start_date THEN 0
        ELSE (CURRENT_DATE - c.care_start_date)::int4
    END AS days_in_care,
    EXISTS (
        SELECT 1
        FROM contract ct
        WHERE ct.client_id = c.id
          AND ct.status = 'approved'
          AND ct.start_date <= CURRENT_TIMESTAMP
          AND ct.end_date >= CURRENT_TIMESTAMP
    ) AS has_active_contract,
    COUNT(*) OVER() AS total_count
FROM client_details c
LEFT JOIN location l ON c.location_id = l.id
LEFT JOIN assigned_employee ae ON ae.client_id = c.id AND ae.role = 'coordinator'
LEFT JOIN employee_profile ep ON ep.id = ae.employee_id
WHERE
    c.status IN ('in_care', 'scheduled_in_care')
    AND (
        sqlc.narg('search')::text IS NULL
        OR sqlc.narg('search')::text = ''
        OR c.first_name ILIKE '%' || sqlc.narg('search') || '%'
        OR c.last_name ILIKE '%' || sqlc.narg('search') || '%'
    )
    AND (
        sqlc.narg('status')::client_status_enum[] IS NULL
        OR c.status = ANY(sqlc.narg('status')::client_status_enum[])
    )
ORDER BY
    CASE
        WHEN sqlc.arg('sort_days_in_care')::text = 'asc' THEN
            CASE
                WHEN c.care_start_date IS NULL OR CURRENT_DATE < c.care_start_date THEN 0
                ELSE (CURRENT_DATE - c.care_start_date)::int4
            END
    END ASC,
    CASE
        WHEN sqlc.arg('sort_days_in_care')::text = 'desc' THEN
            CASE
                WHEN c.care_start_date IS NULL OR CURRENT_DATE < c.care_start_date THEN 0
                ELSE (CURRENT_DATE - c.care_start_date)::int4
            END
    END DESC,
    c.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');


-- name: GetAllClientsIDs :many
SELECT id FROM client_details;




-- name: GetClientDetails :one
SELECT c.*,
       ep.first_name AS bsn_verified_by_first_name,
       ep.last_name AS bsn_verified_by_last_name,
       l.name AS location_name,
       s.name AS sender_name,
       s.email_address AS sender_email_address,
       s.phone_number AS sender_phone_number,
       i.self_sufficiency AS intake_self_sufficiency,
       i.intake_conclusion AS intake_conclusion,
       i.intake_conclusion_notes AS intake_conclusion_notes
FROM client_details c
LEFT JOIN employee_profile ep ON c.bsn_verified_by = ep.id
LEFT JOIN location l ON c.location_id = l.id
LEFT JOIN sender s ON c.sender_id = s.id
LEFT JOIN intake_forms i ON c.intake_form_id = i.id
WHERE c.id = $1 LIMIT 1;


-- name: ListActiveGoalSummariesByClientID :many
SELECT
    cg.title,
    cg.priority,
    COALESCE(cg.topic_name_snapshot, t.topic_name, '') AS topic_name
FROM client_goals cg
LEFT JOIN topics t ON t.id = cg.topic_id
WHERE cg.client_id = $1 AND cg.status = 'active'
ORDER BY cg.sort_order;


-- name: ListTopEmergencyContactsByClientID :many
SELECT
    cec.id,
    cec.first_name,
    cec.last_name,
    cec.relationship,
    cec.phone_number,
    cec.email,
    cec.relation_status
FROM client_emergency_contact cec
WHERE cec.client_id = $1
ORDER BY
    CASE
        WHEN cec.relation_status::text = 'Primary Relationship' THEN 1
        WHEN cec.relation_status::text = 'Secondary Relationship' THEN 2
        ELSE 3
    END,
    cec.created_at ASC
LIMIT 2;


-- name: ListExistingClientDocumentLabels :many
SELECT DISTINCT cd.label::text AS label
FROM client_documents cd
WHERE cd.client_id = $1
ORDER BY label;


-- name: GetClientPageCounts :one
SELECT
    (SELECT COUNT(*)::bigint FROM contract c WHERE c.client_id = $1) AS contracts_count,
    (SELECT COUNT(*)::bigint FROM incident i WHERE i.client_id = $1) AS incidents_count,
    (SELECT COUNT(*)::bigint FROM progress_report pr WHERE pr.client_id = $1) AS reports_count,
    (SELECT COUNT(*)::bigint FROM client_goal_evaluations e WHERE e.client_id = $1) AS evaluations_count,
    (SELECT COUNT(*)::bigint FROM client_documents d WHERE d.client_id = $1) AS documents_count,
    (
        SELECT COUNT(*)::bigint
        FROM calendar_event_attendees cea
        JOIN calendar_events ce ON ce.id = cea.event_id
        WHERE cea.client_id = $1
          AND ce.kind = 'appointment'
          AND ce.status <> 'cancelled'
    ) AS appointments_count;


-- name: GetClientCoordinator :many
SELECT
    ae.employee_id,
    ep.first_name,
    ep.last_name,
    ae.start_date
FROM assigned_employee ae
JOIN employee_profile ep ON ep.id = ae.employee_id
WHERE ae.client_id = $1
  AND ae.role = 'coordinator'
ORDER BY ae.start_date DESC, ae.created_at DESC
LIMIT 1;


-- name: GetClientLatestStatusHistory :one
SELECT
    (
        SELECT csh.reason
        FROM client_status_history csh
        WHERE csh.client_id = $1
        ORDER BY csh.changed_at DESC
        LIMIT 1
    ) AS last_change_reason,
    (
        SELECT csh.changed_at
        FROM client_status_history csh
        WHERE csh.client_id = $1
        ORDER BY csh.changed_at DESC
        LIMIT 1
    ) AS last_changed_at,
    (
        COALESCE((
            SELECT csh.new_status
            FROM client_status_history csh
            WHERE csh.client_id = $1
            ORDER BY csh.changed_at DESC
            LIMIT 1
        ), '')::text
    ) AS last_status;


-- name: ListClientActiveApprovedContracts :many
SELECT
    c.id,
    c.status::text AS status,
    c.start_date,
    c.end_date,
    c.financing_act::text AS financing_act,
    c.financing_option::text AS financing_option,
    c.care_type::text AS care_type,
    (DATE(c.end_date) - CURRENT_DATE)::int4 AS days_until_contract_end
FROM contract c
WHERE c.client_id = $1
  AND c.status = 'approved'
  AND c.start_date <= CURRENT_TIMESTAMP
  AND c.end_date >= CURRENT_TIMESTAMP
ORDER BY c.end_date ASC;



-- -- name: UpdateClientDetails :one
-- UPDATE client_details
-- SET
--     first_name = COALESCE (sqlc.narg('first_name'), first_name),
--     last_name = COALESCE (sqlc.narg('last_name'), last_name),
--     date_of_birth = COALESCE (sqlc.narg('date_of_birth'), date_of_birth),
--     "identity" = COALESCE (sqlc.narg('identity'), "identity"),
--     bsn = COALESCE (sqlc.narg('bsn'), bsn),
--     bsn_verified_by = COALESCE (sqlc.narg('bsn_verified_by'), bsn_verified_by),
--     email = COALESCE (sqlc.narg('email'), email),
--     phone_number = COALESCE (sqlc.narg('phone_number'), phone_number),
--     gender = COALESCE (sqlc.narg('gender'), gender),
--     filenumber = COALESCE (sqlc.narg('filenumber'), filenumber),
--     sender_id = COALESCE (sqlc.narg('sender_id'), sender_id),
--     location_id = COALESCE (sqlc.narg('location_id'), location_id),
--     departure_reason = COALESCE (sqlc.narg('departure_reason'), departure_reason),
--     departure_report = COALESCE (sqlc.narg('departure_report'), departure_report),
--     legal_measure = COALESCE (sqlc.narg('legal_measure'), legal_measure),
--     education_currently_enrolled = COALESCE (sqlc.narg('education_currently_enrolled'), education_currently_enrolled),
--     education_institution = COALESCE (sqlc.narg('education_institution'), education_institution),
--     education_mentor_name = COALESCE (sqlc.narg('education_mentor_name'), education_mentor_name),
--     education_mentor_phone = COALESCE (sqlc.narg('education_mentor_phone'), education_mentor_phone),
--     education_mentor_email = COALESCE (sqlc.narg('education_mentor_email'), education_mentor_email),
--     education_additional_notes = COALESCE (sqlc.narg('education_additional_notes'), education_additional_notes),
--     education_level = COALESCE (sqlc.narg('education_level'), education_level),
--     work_currently_employed = COALESCE (sqlc.narg('work_currently_employed'), work_currently_employed),
--     work_current_employer = COALESCE (sqlc.narg('work_current_employer'), work_current_employer),
--     work_current_employer_phone = COALESCE (sqlc.narg('work_current_employer_phone'), work_current_employer_phone),
--     work_current_employer_email = COALESCE (sqlc.narg('work_current_employer_email'), work_current_employer_email),
--     work_current_position = COALESCE (sqlc.narg('work_current_position'), work_current_position),
--     work_start_date = COALESCE (sqlc.narg('work_start_date'), work_start_date),
--     work_additional_notes = COALESCE (sqlc.narg('work_additional_notes'), work_additional_notes),
--     living_situation = COALESCE (sqlc.narg('living_situation'), living_situation),
--     living_situation_notes = COALESCE (sqlc.narg('living_situation_notes'), living_situation_notes),
--     nationality = COALESCE (sqlc.narg('nationality'), nationality)

-- WHERE id = $1
-- RETURNING *;

-- name: UpdateClientStatus :one
UPDATE client_details
SET status = $2
WHERE id = $1
RETURNING *;

-- name: PutClientInCare :one
UPDATE client_details
SET
    status = $2,
    placed_in_care_at = COALESCE(sqlc.narg('placed_in_care_at'), CURRENT_TIMESTAMP),
    care_start_date = $3
WHERE id = $1
RETURNING *;

-- name: PutClientOutOfCare :one
UPDATE client_details
SET
    status = $2,
    discharge_date = $3,
    discharge_reason = $4,
    final_evaluation = $5
WHERE id = $1
RETURNING *;

-- name: CountActiveGoalsByClientID :one
SELECT COUNT(*)
FROM client_goals
WHERE client_id = $1
  AND status = 'active';

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


-- name: ActivateDueScheduledInCareClients :many
WITH due_clients AS (
    SELECT id
    FROM client_details
    WHERE status = 'scheduled_in_care'
      AND care_start_date IS NOT NULL
      AND care_start_date <= CURRENT_DATE
),
updated AS (
    UPDATE client_details cd
    SET status = 'in_care'
    FROM due_clients dc
    WHERE cd.id = dc.id
    RETURNING cd.id
),
history AS (
    INSERT INTO client_status_history (
        client_id,
        old_status,
        new_status,
        reason
    )
    SELECT
        u.id,
        'scheduled_in_care',
        'in_care',
        'auto_transition_care_start_date_reached'
    FROM updated u
)
SELECT id FROM updated;

-- name: ActivateDueScheduledOutOfCareClients :many
WITH due_clients AS (
    SELECT id
    FROM client_details
    WHERE status = 'scheduled_out_of_care'
      AND discharge_date IS NOT NULL
      AND discharge_date <= CURRENT_DATE
      AND final_evaluation IS NOT NULL
),
updated AS (
    UPDATE client_details cd
    SET status = 'out_of_care'
    FROM due_clients dc
    WHERE cd.id = dc.id
    RETURNING cd.id
),
history AS (
    INSERT INTO client_status_history (
        client_id,
        old_status,
        new_status,
        reason
    )
    SELECT
        u.id,
        'scheduled_out_of_care',
        'out_of_care',
        'auto_transition_discharge_date_reached'
    FROM updated u
)
SELECT id FROM updated;

-- name: ListDueScheduledOutOfCareMissingFinalEvaluation :many
SELECT id
FROM client_details
WHERE status = 'scheduled_out_of_care'
  AND discharge_date IS NOT NULL
  AND discharge_date <= CURRENT_DATE
  AND final_evaluation IS NULL;



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
    SELECT label::text AS label
    FROM client_documents
    WHERE client_id = $1
)
SELECT al.label::text AS missing_label
FROM all_labels al
LEFT JOIN client_labels cl ON al.label = cl.label
WHERE cl.label IS NULL;



-- name: CreateClientLocationTransfer :exec
INSERT INTO client_location_transfer (
    client_id,
    from_location_id,
    to_location_id,
    request_date,
    new_mentor_id,
    reason
) VALUES (
    $1, $2, $3, $4, $5, $6
);


-- name: ApproveOrRejectClientLocationTransfer :exec
UPDATE client_location_transfer
SET
    status = $2,
    approved_rejected_at = NOW(),
    approved_rejected_by = $3
WHERE id = $1;



-- name: ListClientLocationTransfer :many
SELECT
    t.*,
    e.first_name AS mentor_first_name,
    e.last_name AS mentor_last_name,
    COUNT(*) OVER() AS total_count
FROM client_location_transfer t
LEFT JOIN employee_profile e ON t.new_mentor_id = e.id
ORDER BY request_date DESC
LIMIT $1 OFFSET $2;
