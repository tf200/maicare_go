-- name: ListMaturityMatrix :many
SELECT * FROM maturity_matrix;



-- name: CreateClientMaturityMatrixAssessment :one
WITH inserted AS (
    INSERT INTO client_maturity_matrix_assessment (
        client_id,
        maturity_matrix_id,
        start_date,
        end_date,
        initial_level,
        current_level
    ) VALUES (
        $1, $2, $3, $4, $5, $6
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