-- name: ListMaturityMatrix :many
SELECT * FROM maturity_matrix;




-- name: GetMaturityMatrix :one
SELECT * FROM maturity_matrix WHERE id = $1;


-- name: GetLevelDescription :one
SELECT 
    topic_name, 
    (jsonb_path_query_first(level_description, format('$[*] ? (@.level == %s).description', sqlc.arg('level')::text)::jsonpath))::text as level_description 
FROM maturity_matrix 
WHERE id = $1;



-- name: CreateClientMaturityMatrixAssessment :one
WITH inserted AS (
    INSERT INTO client_maturity_matrix_assessment (
        client_id,
        maturity_matrix_id,
        start_date,
        end_date,
        target_level,
        initial_level,
        current_level
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7
    )
    RETURNING *
)
SELECT 
    inserted.*,
    mm.topic_name AS topic_name
FROM inserted
JOIN maturity_matrix mm ON inserted.maturity_matrix_id = mm.id;

-- name: ListClientMaturityMatrixAssessments :many
SELECT
    cma.*,
    mm.topic_name AS topic_name,
    COUNT(*) OVER() AS total_count
FROM client_maturity_matrix_assessment cma
JOIN maturity_matrix mm ON cma.maturity_matrix_id = mm.id
WHERE cma.client_id = $1
ORDER BY cma.start_date DESC
LIMIT $2 OFFSET $3;

-- name: GetClientMaturityMatrixAssessment :one
SELECT
    cma.*,
    mm.topic_name AS topic_name
FROM client_maturity_matrix_assessment cma
JOIN maturity_matrix mm ON cma.maturity_matrix_id = mm.id
WHERE cma.id = $1;



-- name: CreateClientGoal :one
INSERT INTO client_goals (
    client_maturity_matrix_assessment_id,
    description,
    status,
    target_level,
    start_date,
    target_date,
    completion_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;


-- name: ListClientGoals :many
SELECT
    cg.*,
    COUNT(*) OVER() AS total_count
    FROM client_goals cg
WHERE cg.client_maturity_matrix_assessment_id = $1
ORDER BY cg.start_date DESC
LIMIT $2 OFFSET $3;


-- name: GetClientGoal :one
SELECT * FROM client_goals WHERE id = $1;


-- name: CreateGoalObjective :one
INSERT INTO goal_objectives (
    goal_id,
    objective_description,
    due_date,
    status,
    completion_date
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;


-- name: ListGoalObjectives :many
SELECT
    go.*,
    COUNT(*) OVER() AS total_count
FROM goal_objectives go
WHERE go.goal_id = $1
ORDER BY go.due_date DESC;







-- ==================== new code    ====================

-- name: CreateCarePlan :one
INSERT INTO care_plans (
    assessment_id,
    generated_by_employee_id,
    assessment_summary,
    raw_llm_response,
    status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;



-- name: GetCarePlanOverview :one
SELECT 
    cp.*,
    cma.client_id,
    cma.current_level,
    cma.target_level,
    mm.topic_name,
    cd.first_name,
    cd.last_name
FROM care_plans cp
JOIN client_maturity_matrix_assessment cma ON cp.assessment_id = cma.id
JOIN maturity_matrix mm ON cma.maturity_matrix_id = mm.id
JOIN client_details cd ON cma.client_id = cd.id
WHERE cp.id = $1;


-- name: UpdateCarePlanOverview :one
UPDATE care_plans
SET
    assessment_summary = COALESCE(sqlc.narg('assessment_summary'), assessment_summary)
WHERE id = $1
RETURNING *;


-- name: DeleteCarePlan :exec
DELETE FROM care_plans
WHERE id = $1;

-- ==================== care plan objectives and actions ====================


-- name: CreateCarePlanObjective :one
INSERT INTO care_plan_objectives (
    care_plan_id,
    timeframe,
    goal_title,
    description,
    target_date
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetCarePlanObjectivesWithActions :many
SELECT 
    -- Objective fields
    o.id as objective_id,
    o.goal_title as objective_title,
    o.description as objective_description,
    o.timeframe as objective_timeframe,
    o.status as objective_status,

    -- Action fields (will be NULL if no actions exist)
    a.id as action_id,
    a.action_description,
    a.is_completed,
    a.notes as action_notes,
    a.sort_order
FROM care_plan_objectives o
LEFT JOIN care_plan_actions a ON o.id = a.objective_id
WHERE o.care_plan_id = $1
ORDER BY 
    o.timeframe,
    o.id,
    a.sort_order;


-- name: UpdateCarePlanObjective :one
UPDATE care_plan_objectives
SET
    timeframe = COALESCE(sqlc.narg('timeframe'), timeframe),
    goal_title = COALESCE(sqlc.narg('goal_title'), goal_title),
    description = COALESCE(sqlc.narg('description'), description),
    status = COALESCE(sqlc.narg('status'), status), 
    updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: DeleteCarePlanObjective :exec
DELETE FROM care_plan_objectives
WHERE id = $1;

-- name: CreateCarePlanAction :one
INSERT INTO care_plan_actions (
    objective_id,
    action_description,
    sort_order
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetCarePlanActionsMaxSortOrder :one
SELECT COALESCE(MAX(sort_order), 0)::INT AS max_sort_order
FROM care_plan_actions
WHERE objective_id = $1;

-- name: UpdateCarePlanAction :one
UPDATE care_plan_actions
SET
    action_description = COALESCE(sqlc.narg('action_description'), action_description)
WHERE id = $1
RETURNING *;

-- name: DeleteCarePlanAction :exec
DELETE FROM care_plan_actions
WHERE id = $1;



-- ==================== care plan interventions ====================
-- name: CreateCarePlanIntervention :one
INSERT INTO care_plan_interventions (
    care_plan_id,
    frequency,
    intervention_description
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetCarePlanInterventions :many
SELECT * FROM care_plan_interventions 
WHERE care_plan_id = $1 AND is_active = true 
ORDER BY 
    CASE frequency 
        WHEN 'daily' THEN 1 
        WHEN 'weekly' THEN 2 
        WHEN 'monthly' THEN 3 
    END;

-- name: UpdateCarePlanIntervention :one
UPDATE care_plan_interventions
SET
    frequency = COALESCE(sqlc.narg('frequency'), frequency),
    intervention_description = COALESCE(sqlc.narg('intervention_description'), intervention_description),
    updated_at = NOW()
WHERE id = $1
RETURNING *;


-- name: DeleteCarePlanIntervention :exec
DELETE FROM care_plan_interventions
WHERE id = $1;



-- ==================== care plan success metrics ====================

-- name: CreateCarePlanSuccessMetric :one
INSERT INTO care_plan_metrics (
    care_plan_id,
    metric_name,
    target_value,
    measurement_method
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetCarePlanSuccessMetrics :many
SELECT * FROM care_plan_metrics 
WHERE care_plan_id = $1 
ORDER BY created_at;


-- name: UpdateCarePlanSuccessMetric :one
UPDATE care_plan_metrics
SET
    metric_name = COALESCE(sqlc.narg('metric_name'), metric_name),
    target_value = COALESCE(sqlc.narg('target_value'), target_value),
    measurement_method = COALESCE(sqlc.narg('measurement_method'), measurement_method),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCarePlanSuccessMetric :exec
DELETE FROM care_plan_metrics
WHERE id = $1;


-- ==================== care plan risks ====================

-- name: CreateCarePlanRisk :one
INSERT INTO care_plan_risks (
    care_plan_id,
    risk_description,
    mitigation_strategy,
    risk_level
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;


-- name: GetCarePlanRisks :many
SELECT * FROM care_plan_risks 
WHERE care_plan_id = $1 AND is_active = true 
ORDER BY 
    CASE risk_level 
        WHEN 'high' THEN 1 
        WHEN 'medium' THEN 2 
        WHEN 'low' THEN 3 
    END;

-- name: UpdateCarePlanRisk :one
UPDATE care_plan_risks
SET
    risk_description = COALESCE(sqlc.narg('risk_description'), risk_description),
    mitigation_strategy = COALESCE(sqlc.narg('mitigation_strategy'), mitigation_strategy),
    risk_level = COALESCE(sqlc.narg('risk_level'), risk_level),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCarePlanRisk :exec
DELETE FROM care_plan_risks
WHERE id = $1; 


-- ===================== care plan support network ====================

-- name: CreateCarePlanSupportNetwork :one
INSERT INTO care_plan_support_network (
    care_plan_id,
    role_title,
    responsibility_description,
    contact_person,
    contact_details
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetCarePlanSupportNetwork :many
SELECT * FROM care_plan_support_network 
WHERE care_plan_id = $1 AND is_active = true 
ORDER BY created_at;

-- name: UpdateCarePlanSupportNetwork :one
UPDATE care_plan_support_network
SET
    role_title = COALESCE(sqlc.narg('role_title'), role_title),
    responsibility_description = COALESCE(sqlc.narg('responsibility_description'), responsibility_description),
    contact_person = COALESCE(sqlc.narg('contact_person'), contact_person),
    contact_details = COALESCE(sqlc.narg('contact_details'), contact_details),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCarePlanSupportNetwork :exec
DELETE FROM care_plan_support_network
WHERE id = $1;




-- ==================== care plan resources ====================

-- name: CreateCarePlanResources :one
INSERT INTO care_plan_resources (
    care_plan_id,
    resource_description,
    is_obtained,
    obtained_date
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetCarePlanResources :many
SELECT * FROM care_plan_resources 
WHERE care_plan_id = $1 
ORDER BY is_obtained, created_at;

-- name: UpdateCarePlanResource :one
UPDATE care_plan_resources
SET
    resource_description = COALESCE(sqlc.narg('resource_description'), resource_description),
    is_obtained = COALESCE(sqlc.narg('is_obtained'), is_obtained),
    obtained_date = COALESCE(sqlc.narg('obtained_date'), obtained_date),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCarePlanResource :exec
DELETE FROM care_plan_resources
WHERE id = $1;