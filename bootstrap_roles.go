package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"

	db "maicare_go/db/sqlc"
	rolecfg "maicare_go/roles"
)

type rbacPermission = rolecfg.Permission
type rbacRole = rolecfg.Role
type rbacConfig = rolecfg.Config

//go:embed roles/rbac_config.yaml
var embeddedRBACConfig []byte

func ensureRolesAndPermissionsOnStartup(ctx context.Context, store *db.Store) error {
	var conf rbacConfig
	if err := yaml.Unmarshal(embeddedRBACConfig, &conf); err != nil {
		return fmt.Errorf("failed to parse embedded rbac config: %w", err)
	}

	for _, perm := range conf.Permissions {
		if err := ensurePermission(ctx, store, perm); err != nil {
			return err
		}
	}

	for _, role := range conf.Roles {
		if err := ensureRoleWithPermissions(ctx, store, role); err != nil {
			return err
		}
	}

	log.Printf("RBAC bootstrap ready (%d permissions, %d roles)", len(conf.Permissions), len(conf.Roles))
	return nil
}

func ensurePermission(ctx context.Context, store *db.Store, perm rbacPermission) error {
	methodJSON, err := json.Marshal(perm.Method)
	if err != nil {
		return fmt.Errorf("failed to marshal methods for permission %s: %w", perm.Name, err)
	}

	metadata := perm.Normalize()

	var permissionID uuid.UUID
	err = store.ConnPool.QueryRow(
		ctx,
		`SELECT id FROM permissions WHERE name = $1 ORDER BY id LIMIT 1`,
		perm.Name,
	).Scan(&permissionID)
	if err == nil {
		_, err = store.ConnPool.Exec(
			ctx,
			`UPDATE permissions
			 SET resource = $2,
			     method = $3,
			     group_key = $4,
			     section_key = $5,
			     display_name = $6,
			     description = $7,
			     sort_order = $8
			 WHERE id = $1`,
			permissionID,
			perm.Resource,
			string(methodJSON),
			metadata.GroupKey,
			metadata.SectionKey,
			metadata.DisplayName,
			metadata.Description,
			metadata.SortOrder,
		)
		if err != nil {
			return fmt.Errorf("failed to update permission %s: %w", perm.Name, err)
		}
		return nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("failed to check permission %s: %w", perm.Name, err)
	}

	_, err = store.ConnPool.Exec(
		ctx,
		`INSERT INTO permissions
			(name, resource, method, group_key, section_key, display_name, description, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		perm.Name,
		perm.Resource,
		string(methodJSON),
		metadata.GroupKey,
		metadata.SectionKey,
		metadata.DisplayName,
		metadata.Description,
		metadata.SortOrder,
	)
	if err != nil {
		return fmt.Errorf("failed to insert permission %s: %w", perm.Name, err)
	}

	return nil
}

func ensureRoleWithPermissions(ctx context.Context, store *db.Store, role rbacRole) error {
	var roleID uuid.UUID
	err := store.ConnPool.QueryRow(ctx, `SELECT id FROM roles WHERE name = $1 LIMIT 1`, role.Name).Scan(&roleID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to query role %s: %w", role.Name, err)
		}

		var description *string
		if role.Description != "" {
			description = &role.Description
		}

		createdRole, createErr := store.CreateRole(ctx, db.CreateRoleParams{
			Name:        role.Name,
			Description: description,
		})
		if createErr != nil {
			if isConflictError(createErr) {
				err = store.ConnPool.QueryRow(ctx, `SELECT id FROM roles WHERE name = $1 LIMIT 1`, role.Name).Scan(&roleID)
				if err != nil {
					return fmt.Errorf("role %s conflict detected but lookup failed: %w", role.Name, createErr)
				}
			} else {
				return fmt.Errorf("failed to create role %s: %w", role.Name, createErr)
			}
		} else {
			roleID = createdRole.ID
		}
	}

	if _, err = store.ConnPool.Exec(
		ctx,
		`UPDATE roles SET description = $2 WHERE id = $1`,
		roleID,
		nilIfEmpty(role.Description),
	); err != nil {
		return fmt.Errorf("failed to update role description for %s: %w", role.Name, err)
	}

	for _, permissionName := range role.Permissions {
		var permissionID uuid.UUID
		err = store.ConnPool.QueryRow(
			ctx,
			`SELECT id FROM permissions WHERE name = $1 ORDER BY id LIMIT 1`,
			permissionName,
		).Scan(&permissionID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Printf("RBAC bootstrap warning: permission %q not found for role %q", permissionName, role.Name)
				continue
			}
			return fmt.Errorf("failed to query permission %s for role %s: %w", permissionName, role.Name, err)
		}

		if _, err = store.ConnPool.Exec(
			ctx,
			`INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2)
			 ON CONFLICT (role_id, permission_id) DO NOTHING`,
			roleID,
			permissionID,
		); err != nil {
			return fmt.Errorf("failed to map permission %s to role %s: %w", permissionName, role.Name, err)
		}
	}

	return nil
}

func nilIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
