package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"

	rolecfg "maicare_go/roles"
)

type Permission = rolecfg.Permission
type Role = rolecfg.Role
type Config = rolecfg.Config

func main() {
	dbSourceFlag := flag.String("db", "", "database connection string (defaults to DB_SOURCE then local default)")
	flag.Parse()

	configFile, err := os.ReadFile("roles/rbac_config.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

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

	fmt.Println("Starting database sync...")

	// Process permissions
	fmt.Println("\nProcessing permissions...")
	for _, perm := range config.Permissions {
		methodJSON, err := json.Marshal(perm.Method)
		if err != nil {
			panic(err)
		}

		metadata := perm.Normalize()

		var permissionID uuid.UUID
		err = tx.QueryRow(context.Background(), "SELECT id FROM permissions WHERE name=$1 ORDER BY id LIMIT 1", perm.Name).Scan(&permissionID)
		if err != nil && err != pgx.ErrNoRows {
			panic(err)
		}

		if err == pgx.ErrNoRows {
			_, err = tx.Exec(
				context.Background(),
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
				panic(err)
			}
			fmt.Println("Inserted permission:", perm.Name)
		} else {
			_, err = tx.Exec(
				context.Background(),
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
				panic(err)
			}
			fmt.Println("Updated permission:", perm.Name)
		}
	}

	// Process roles
	fmt.Println("\nProcessing roles...")
	for _, role := range config.Roles {
		var roleID uuid.UUID
		err = tx.QueryRow(context.Background(), "SELECT id FROM roles WHERE name=$1", role.Name).Scan(&roleID)
		if err != nil && err != pgx.ErrNoRows {
			panic(err)
		}
		if err == pgx.ErrNoRows {
			roleID = uuid.New()
			_, err = tx.Exec(
				context.Background(),
				"INSERT INTO roles (id, name, description) VALUES ($1, $2, $3)",
				roleID,
				role.Name,
				nilIfEmpty(role.Description),
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
				nilIfEmpty(role.Description),
			)
			if err != nil {
				panic(err)
			}
			fmt.Println("Role already exists:", role.Name)
		}

		// Sync role_permissions
		for _, permName := range role.Permissions {
			var permID uuid.UUID
			err = tx.QueryRow(context.Background(), "SELECT id FROM permissions WHERE name=$1", permName).Scan(&permID)
			if err != nil {
				if err == pgx.ErrNoRows {
					fmt.Printf("Warning: permission '%s' not found, skipping.\n", permName)
					continue
				}
				panic(err)
			}

			_, err = tx.Exec(context.Background(), "INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT (role_id, permission_id) DO NOTHING", roleID, permID)
			if err != nil {
				panic(err)
			}
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println("\nDatabase sync completed successfully!")
}

func nilIfEmpty(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
