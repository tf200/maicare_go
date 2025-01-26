-- name: Createallergy :one
INSERT INTO allergy_type (
    name
) VALUES (
    $1
) RETURNING *;



-- name: ListAllergies :many
SELECT * FROM allergy_type 
WHERE 
    CASE 
        WHEN sqlc.narg('search')::text IS NULL THEN true
        ELSE name ILIKE concat('%', sqlc.narg('search')::text, '%')
    END
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;