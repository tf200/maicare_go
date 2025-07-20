-- name: GetTemplateItemsByIds :many
-- GetTemplateItemsByIds retrieves a list of template item IDs that match the given input IDs.
SELECT id
FROM template_items
WHERE id IN (sqlc.slice('ids'));








-- name: GetInvoiceTemplateItems :many 
SELECT * 
FROM template_items
WHERE id IN (sqlc.slice('ids'))
GROUP BY source_table;
