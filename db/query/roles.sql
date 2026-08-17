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
INSERT INTO roles (name, description)
VALUES (sqlc.arg('name'), sqlc.narg('description'))
RETURNING *;

-- name: GetAdminRoleId :one
/* Returns the ID of the admin role. */
SELECT id
FROM roles
WHERE name = 'admin';

-- name: ListRoles :many
/* Returns every role ordered by id with count of permissions and employees. */
SELECT 
    r.id,
    r.name,
    r.description,
    COALESCE(COUNT(DISTINCT rp.permission_id), 0)::BIGINT AS permission_count,
    COALESCE(COUNT(DISTINCT ep.id), 0)::BIGINT AS employee_count
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
LEFT JOIN user_roles ur ON r.id = ur.role_id
LEFT JOIN employee_profile ep ON ep.user_id = ur.user_id
GROUP BY r.id, r.name, r.description
ORDER BY r.id;

/* ---------- 2. PERMISSIONS ---------- */

-- name: ListAllPermissions :many
/* Returns every permission ordered by id. */
SELECT *
FROM permissions
ORDER BY group_key, section_key, name;

/* ---------- 3. ROLE-PERMISSION MAPPING ---------- */

-- name: ListAllRolePermissions :many
/* Returns all permissions attached to a single role. */
SELECT p.id   AS permission_id,
       p.name AS permission_name,
       p.is_scoped,
       rp.scope
FROM role_permissions rp
JOIN permissions p ON p.id = rp.permission_id
WHERE rp.role_id = $1
ORDER BY p.id;


-- name: AddPermissionToRole :exec
/* Insert one permission grant while replacing a role's grants transactionally. */
INSERT INTO role_permissions (role_id, permission_id, scope)
VALUES (sqlc.arg('role_id'), sqlc.arg('permission_id'), sqlc.narg('scope'));

-- name: RemovePermissionsFromRole :exec
/* Removes *all* permissions from the given role. */
DELETE FROM role_permissions
WHERE role_id = $1;

/* ---------- 4. USER-ROLE MAPPING ---------- */

-- name: GetUserRoles :many
/* Returns every role granted to a user. */
SELECT r.id, r.name
FROM user_roles ur
JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = $1;

-- name: AssignRoleToUser :exec
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE SET role_id = $2;

/* ---------- 5. USER PERMISSIONS ---------- */

-- name: ListEffectiveUserPermissions :many
/* Returns permissions granted by the user's assigned role. */
SELECT p.id AS permission_id,
       p.name AS permission_name,
       p.is_scoped,
       rp.scope
FROM user_roles ur
JOIN role_permissions rp ON rp.role_id = ur.role_id
JOIN permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1
  AND ((p.is_scoped AND rp.scope IS NOT NULL)
       OR (NOT p.is_scoped AND rp.scope IS NULL))
ORDER BY p.id;

/* ---------- 6. CHECK UTILITIES ---------- */

-- name: GetEffectiveUserPermission :one
/* Returns the usable role grant for one named permission. */
SELECT p.id AS permission_id,
       p.name AS permission_name,
       p.is_scoped,
       rp.scope
FROM user_roles ur
JOIN role_permissions rp ON rp.role_id = ur.role_id
JOIN permissions p ON p.id = rp.permission_id
WHERE ur.user_id = $1
  AND p.name = $2
  AND ((p.is_scoped AND rp.scope IS NOT NULL)
       OR (NOT p.is_scoped AND rp.scope IS NULL))
LIMIT 1;
