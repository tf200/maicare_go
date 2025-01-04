-- name: CheckRolePermission :one
SELECT EXISTS (
    SELECT 1
    FROM role_permissions rp
    JOIN permissions p ON p.id = rp.permission_id
    WHERE rp.role_id = $1
    AND p.name = $2
) AS has_permission;


-- name: GetPermissionsByRoleID :many
SELECT p.id, p.name, p.resource, p.method
FROM permissions p
JOIN role_permissions rp ON p.id = rp.permission_id
WHERE rp.role_id = $1;