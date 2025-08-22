-- name: CheckUserPermission :one
SELECT EXISTS (
    SELECT 1
    FROM user_permissions up
    JOIN permissions p ON p.id = up.permission_id
    WHERE up.user_id = $1
      AND p.name = $2
) AS has_permission;




-- name: ListRoles :many
SELECT * FROM roles
ORDER BY id;


-- name: GrantRoleToUser :exec
WITH ins_role AS (
    INSERT INTO user_roles (user_id, role_id)
    VALUES ($1, $2)
    ON CONFLICT (user_id, role_id) DO NOTHING
)
INSERT INTO user_permissions (user_id, permission_id)
SELECT $1, rp.permission_id
FROM   role_permissions rp
WHERE  rp.role_id = $2
ON CONFLICT (user_id, permission_id) DO NOTHING;


-- name: ListAllPermissions :many
SELECT * FROM permissions
ORDER BY id;


-- name: ListAllRolePermissions :many
SELECT p.id AS permission_id,
       p.name AS permission_name,
       p.resource AS permission_resource
FROM role_permissions rp
JOIN permissions p ON rp.permission_id = p.id
WHERE rp.role_id = $1
ORDER BY p.id;

-- name: ListUserRoles :many
SELECT r.id, r.name
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.user_id = $1
ORDER BY r.id;


-- name: ListUserPermissions :many
SELECT p.id AS permission_id,
       p.name AS permission_name,
       p.resource AS permission_resource
FROM user_permissions up
JOIN permissions p ON up.permission_id = p.id
WHERE up.user_id = $1
ORDER BY p.id;