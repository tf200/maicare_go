-- name: GetTemplateItemsByIds :many
SELECT id
FROM template_items
WHERE id = ANY($1::bigint[]);




-- name: GetTemplateItemsBySourceTable :many
SELECT *
FROM template_items
WHERE id = ANY($1::bigint[])
ORDER BY source_table;

-- name: GetAllTemplateItems :many
SELECT *
FROM template_items
ORDER BY source_table;