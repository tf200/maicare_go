-- name: CheckRolePermission :one
SELECT EXISTS (
    SELECT 1
    FROM role_permissions rp
    JOIN permissions p ON p.id = rp.permission_id
    WHERE rp.role_id = $1
    AND p.resource = $2
    AND p.name = $3
) AS exists;