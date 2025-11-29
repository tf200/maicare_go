-- name: GetTemplateItemsByIds :many
SELECT id
FROM template_items
WHERE id = ANY($1::uuid[]);




-- name: GetTemplateItemsBySourceTable :many
SELECT *
FROM template_items
WHERE id = ANY($1::uuid[])
ORDER BY source_table;

-- name: GetAllTemplateItems :many
SELECT *
FROM template_items
ORDER BY source_table;