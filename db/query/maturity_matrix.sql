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

-- name: CreateObjective :one
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

-- name: CreateGoalAction :one
INSERT INTO care_plan_actions (
    objective_id,
    action_description,
    sort_order
) VALUES (
    $1, $2, $3
)
RETURNING *;


-- name: CreateIntervention :one
INSERT INTO care_plan_interventions (
    care_plan_id,
    frequency,
    intervention_description
) VALUES (
    $1, $2, $3
)
RETURNING *;


-- name: CreateSuccessMetric :one
INSERT INTO care_plan_metrics (
    care_plan_id,
    metric_name,
    target_value,
    measurement_method
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: CreateCarePlanRisks :one
INSERT INTO care_plan_risks (
    care_plan_id,
    risk_description,
    mitigation_strategy
) VALUES (
    $1, $2, $3
)
RETURNING *;


-- name: CreateCarePlanSupportNetwork :one
INSERT INTO care_plan_support_network (
    care_plan_id,
    role_title,
    responsibility_description
) VALUES (
    $1, $2, $3
)
RETURNING *;



-- name: CreateCarePlanResources :one
INSERT INTO care_plan_resources (
    care_plan_id,
    resource_description
) VALUES (
    $1, $2
)
RETURNING *;


-- -- name: GetCarePlan :one
-- SELECT 
--     cp.*,
--     cma.client_id,
--     cma.current_level,
--     cma.target_level,
--     mm.domain_name,
--     cd.first_name,
--     cd.last_name
-- FROM care_plans cp
-- JOIN client_maturity_matrix_assessment cma ON cp.assessment_id = cma.id
-- JOIN maturity_matrix mm ON cma.maturity_matrix_id = mm.id
-- JOIN client_details cd ON cma.client_id = cd.id
-- WHERE cp.id = $1;


-- -- name: GetObjectivesWithActions :many
-- SELECT 
--     o.id,
--     o.care_plan_id,
--     o.description,
--     o.timeframe,
--     o.created_at,
--     o.updated_at,
--     -- Add any other objective columns you need
--     json_agg(
--         json_build_object(
--             'id', a.id,
--             'description', a.action_description,
--             'is_completed', a.is_completed,
--             'completed_at', a.completed_at,
--             'sort_order', a.sort_order
--         ) ORDER BY a.sort_order
--     ) as actions
-- FROM care_plan_objectives o
-- LEFT JOIN care_plan_actions a ON o.id = a.objective_id
-- WHERE o.care_plan_id = $1
-- GROUP BY o.id, o.care_plan_id, o.description, o.timeframe, o.created_at, o.updated_at
-- ORDER BY 
--     CASE o.timeframe 
--         WHEN 'short_term' THEN 1 
--         WHEN 'medium_term' THEN 2 
--         WHEN 'long_term' THEN 3 
--     END;


--     -- name: GetInterventions :many
-- SELECT * FROM care_plan_interventions 
-- WHERE care_plan_id = $1 AND is_active = true 
-- ORDER BY 
--     CASE frequency 
--         WHEN 'daily' THEN 1 
--         WHEN 'weekly' THEN 2 
--         WHEN 'monthly' THEN 3 
--     END;

-- -- name: GetMetrics :many
-- SELECT * FROM care_plan_metrics 
-- WHERE care_plan_id = $1 
-- ORDER BY created_at;

-- -- name: GetRisks :many
-- SELECT * FROM care_plan_risks 
-- WHERE care_plan_id = $1 AND is_active = true 
-- ORDER BY 
--     CASE risk_level 
--         WHEN 'high' THEN 1 
--         WHEN 'medium' THEN 2 
--         WHEN 'low' THEN 3 
--     END;

-- -- name: GetSupportNetwork :many
-- SELECT * FROM care_plan_support_network 
-- WHERE care_plan_id = $1 AND is_active = true 
-- ORDER BY created_at;

-- -- name: GetResources :many
-- SELECT * FROM care_plan_resources 
-- WHERE care_plan_id = $1 
-- ORDER BY is_obtained, created_at;