/*
 *  RBAC – Role & Permission Management
 *  This file contains all sqlc queries for the simple role-based
 *  access-control system.  They are grouped in the following order:
 *    1. Core roles
 *    2. Core permissions
 *    3. Role <-> Permission mapping
 *    4. User <-> Role mapping
 *    5. User <-> Permission mapping
 *    6. Read queries
 *    7. Helper utilities
 *
 *  sqlc will generate Go types and funcs named after the `-- name:` tags.
 */

/* ---------- 1. ROLES ---------- */

-- name: CreateRole :one
/* Insert a new role and return the created row. */
INSERT INTO roles (name)
VALUES ($1)
RETURNING *;

-- name: ListRoles :many
/* Returns every role ordered by id with count of permissions. */
SELECT 
    r.*,
        COALESCE(COUNT(rp.permission_id), 0)::BIGINT AS permission_count
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
GROUP BY r.id, r.name
ORDER BY r.id;

/* ---------- 2. PERMISSIONS ---------- */

-- name: ListAllPermissions :many
/* Returns every permission ordered by id. */
SELECT *
FROM permissions
ORDER BY id;

/* ---------- 3. ROLE-PERMISSION MAPPING ---------- */

-- name: ListAllRolePermissions :many
/* Returns all permissions attached to a single role. */
SELECT p.id   AS permission_id,
       p.name AS permission_name,
       p.resource
FROM role_permissions rp
JOIN permissions p ON p.id = rp.permission_id
WHERE rp.role_id = $1
ORDER BY p.id;


-- name: AddPermissionsToRole :exec
/* Bulk-insert permission IDs into a role (idempotent). */
INSERT INTO role_permissions (role_id, permission_id)
SELECT sqlc.arg('role_id'), unnest(sqlc.arg('permission_ids')::int[])
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- name: RemovePermissionsFromRole :exec
/* Removes *all* permissions from the given role. */
DELETE FROM role_permissions
WHERE role_id = $1;

/* ---------- 4. USER-ROLE MAPPING ---------- */

-- name: ListUserRoles :many
/* Returns every role granted to a user. */
SELECT r.id, r.name
FROM user_roles ur
JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = $1
ORDER BY r.id;

-- name: GrantRoleToUser :exec
/*
 * Assigns a role to a user (idempotent).
 * As a convenience it also copies all the role’s permissions
 * into the user_permissions table so permission checks stay cheap.
 */
WITH ins_role AS (
    INSERT INTO user_roles (user_id, role_id)
    VALUES ($1, $2)
    ON CONFLICT (user_id, role_id) DO NOTHING
)
INSERT INTO user_permissions (user_id, permission_id)
SELECT $1, rp.permission_id
FROM role_permissions rp
WHERE rp.role_id = $2
ON CONFLICT (user_id, permission_id) DO NOTHING;

/* ---------- 5. USER-PERMISSION MAPPING ---------- */

-- name: ListUserPermissions :many
/* Returns every permission granted to a user (direct or via roles). */
SELECT p.id   AS permission_id,
       p.name AS permission_name,
       p.resource
FROM user_permissions up
JOIN permissions p ON p.id = up.permission_id
WHERE up.user_id = $1
ORDER BY p.id;

-- name: GrantUserPermissions :exec
/* Bulk-insert permission IDs for a user (idempotent). */
INSERT INTO user_permissions (user_id, permission_id)
SELECT sqlc.arg('user_id'), unnest(sqlc.arg('permission_ids')::int[])
ON CONFLICT (user_id, permission_id) DO NOTHING;

-- name: DeleteUserPermissions :exec
/* Removes *all* permissions from the given user. */
DELETE FROM user_permissions
WHERE user_id = $1;

/* ---------- 6. CHECK UTILITIES ---------- */

-- name: CheckUserPermission :one
/* Returns true/false whether the user has the named permission. */
SELECT EXISTS (
    SELECT 1
    FROM user_permissions up
    JOIN permissions p ON p.id = up.permission_id
    WHERE up.user_id = $1
      AND p.name = $2
) AS has_permission;