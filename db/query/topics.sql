-- name: ListCarePlanTopics :many
SELECT
    topic.id,
    topic.topic_name,
    topic.level_description
FROM topics topic;



-- name: GetTopicLevel :one
SELECT
    t.topic_name,
    -- Explicitly cast to text to force string type in Go
    (element.value ->> 'name')::text as level_name,
    (element.value ->> 'description')::text as level_description
FROM
    topics t,
    jsonb_array_elements(t.level_description) element
WHERE
    t.id = $1
    AND (element.value ->> 'level')::int = sqlc.arg('level')::int;
