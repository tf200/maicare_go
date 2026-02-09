-- name: CreateClientGoalsFromIntakeAssessments :many
WITH intake_rows AS (
    SELECT
        ita.id AS intake_assessment_id,
        ita.topic_id,
        t.topic_name,
        g.goal_item,
        g.goal_ord
    FROM intake_topic_assessments ita
    JOIN topics t ON t.id = ita.topic_id
    CROSS JOIN LATERAL jsonb_array_elements(COALESCE(ita.proposed_goals, '[]'::jsonb))
        WITH ORDINALITY AS g(goal_item, goal_ord)
    WHERE ita.intake_form_id = sqlc.arg(intake_form_id)::uuid
),
normalized AS (
    SELECT
        intake_assessment_id,
        topic_id,
        topic_name,
        goal_ord,
        CASE
            WHEN jsonb_typeof(goal_item) = 'object' THEN
                COALESCE(
                    NULLIF(BTRIM(goal_item->>'title'), ''),
                    NULLIF(BTRIM(goal_item->>'description'), '')
                )
            WHEN jsonb_typeof(goal_item) = 'string' THEN
                NULLIF(BTRIM(TRIM(BOTH '"' FROM goal_item::text)), '')
            ELSE NULL
        END AS title,
        CASE
            WHEN jsonb_typeof(goal_item) = 'object' THEN NULLIF(BTRIM(goal_item->>'description'), '')
            ELSE NULL
        END AS description,
        CASE
            WHEN jsonb_typeof(goal_item) = 'object'
                 AND LOWER(COALESCE(goal_item->>'priority', 'medium')) IN ('low', 'medium', 'high')
                THEN LOWER(goal_item->>'priority')
            ELSE 'medium'
        END AS priority
    FROM intake_rows
),
inserted AS (
    INSERT INTO client_goals (
        client_id,
        title,
        description,
        priority,
        status,
        topic_id,
        topic_name_snapshot,
        source,
        origin_intake_assessment_id,
        sort_order
    )
    SELECT
        sqlc.arg(client_id)::uuid,
        n.title,
        n.description,
        n.priority::client_goal_priority_enum,
        'active'::client_goal_status_enum,
        n.topic_id,
        n.topic_name,
        'intake'::client_goal_source_enum,
        n.intake_assessment_id,
        ROW_NUMBER() OVER (ORDER BY n.topic_name, n.intake_assessment_id, n.goal_ord)::int - 1
    FROM normalized n
    WHERE n.title IS NOT NULL
    RETURNING *
)
SELECT * FROM inserted;
