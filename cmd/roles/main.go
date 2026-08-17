package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func main() {
	dbSourceFlag := flag.String("db", "", "database connection string (defaults to DB_SOURCE then local default)")
	flag.Parse()

	dbSource := strings.TrimSpace(*dbSourceFlag)
	if dbSource == "" {
		dbSource = strings.TrimSpace(os.Getenv("DB_SOURCE"))
	}
	if dbSource == "" {
		dbSource = "postgres://maicare:maicare@localhost:5432/maicare?sslmode=disable"
	}

	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		panic(err)
	}
	defer func() { _ = conn.Close(context.Background()) }()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	defer func() { _ = tx.Rollback(context.Background()) }()

	fmt.Println("Starting database RBAC sync...")

	// Process permissions from domain source of truth
	fmt.Println("\nProcessing permissions from domain definitions...")
	permissions := domain.AllPermissionDefinitions()

	for _, perm := range permissions {
		permName := string(perm.Key)
		var permissionID uuid.UUID
		err = tx.QueryRow(context.Background(), "SELECT id FROM permissions WHERE name=$1 LIMIT 1", permName).Scan(&permissionID)
		if err != nil && err != pgx.ErrNoRows {
			panic(err)
		}

		var desc *string
		if perm.Description != "" {
			desc = &perm.Description
		}

		if err == pgx.ErrNoRows {
			_, err = tx.Exec(
				context.Background(),
				`INSERT INTO permissions
					(name, group_key, section_key, display_name, description, is_scoped)
				 VALUES ($1, $2, $3, $4, $5, $6)`,
				permName,
				perm.GroupKey,
				perm.SectionKey,
				perm.DisplayName,
				desc,
				perm.IsScoped,
			)
			if err != nil {
				panic(err)
			}
			fmt.Println("Inserted permission:", permName)
		} else {
			_, err = tx.Exec(
				context.Background(),
				`UPDATE permissions
				 SET group_key = $2,
					     section_key = $3,
					     display_name = $4,
					     description = $5,
					     is_scoped = $6
					 WHERE id = $1`,
				permissionID,
				perm.GroupKey,
				perm.SectionKey,
				perm.DisplayName,
				desc,
				perm.IsScoped,
			)
			if err != nil {
				panic(err)
			}
			fmt.Println("Updated permission:", permName)
		}
	}

	// Process roles
	fmt.Println("\nProcessing system role seeds...")
	roles := domain.DefaultRoleSeeds()

	for _, role := range roles {
		var roleID uuid.UUID
		err = tx.QueryRow(context.Background(), "SELECT id FROM roles WHERE name=$1", role.Name).Scan(&roleID)
		if err != nil && err != pgx.ErrNoRows {
			panic(err)
		}

		var desc *string
		if role.Description != "" {
			desc = &role.Description
		}

		if err == pgx.ErrNoRows {
			roleID = uuid.New()
			_, err = tx.Exec(
				context.Background(),
				"INSERT INTO roles (id, name, description) VALUES ($1, $2, $3)",
				roleID,
				role.Name,
				desc,
			)
			if err != nil {
				panic(err)
			}
			fmt.Println("Inserted role:", role.Name)
		} else {
			_, err = tx.Exec(
				context.Background(),
				"UPDATE roles SET description = $2 WHERE id = $1",
				roleID,
				desc,
			)
			if err != nil {
				panic(err)
			}
			fmt.Println("Role already exists:", role.Name)
		}

		// Sync role_permissions
		for _, permKey := range role.Permissions {
			permName := string(permKey)
			var permID uuid.UUID
			var isScoped bool
			err = tx.QueryRow(context.Background(), "SELECT id, is_scoped FROM permissions WHERE name=$1", permName).Scan(&permID, &isScoped)
			if err != nil {
				if err == pgx.ErrNoRows {
					fmt.Printf("Warning: permission '%s' not found, skipping.\n", permName)
					continue
				}
				panic(err)
			}

			var scope *string
			if isScoped {
				if role.Scope == "" {
					panic(fmt.Sprintf("role %q has no default scope for scoped permission %q", role.Name, permName))
				}
				scopeValue := string(role.Scope)
				scope = &scopeValue
			}

			_, err = tx.Exec(
				context.Background(),
				`INSERT INTO role_permissions (role_id, permission_id, scope)
				 VALUES ($1, $2, $3)
				 ON CONFLICT (role_id, permission_id) DO UPDATE SET scope = EXCLUDED.scope`,
				roleID,
				permID,
				scope,
			)
			if err != nil {
				panic(err)
			}
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("\nDatabase RBAC sync completed successfully!")
}
