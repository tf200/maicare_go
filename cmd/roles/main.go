package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"gopkg.in/yaml.v3"
)

type Permission struct {
	Name     string   `yaml:"name"`
	Resource string   `yaml:"resource"`
	Method   []string `yaml:"method"`
}

type Role struct {
	Name        string   `yaml:"name"`
	ID          int      `yaml:"id"` // ignored, using UUID
	Description string   `yaml:"description"`
	Permissions []string `yaml:"permissions"`
}

type Config struct {
	Permissions []Permission `yaml:"permissions"`
	Roles       []Role       `yaml:"roles"`
}

func main() {
	configFile, err := os.ReadFile("roles/rbac_config.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	dbSource := "postgres://maicare:maicare@localhost:5432/maicare?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), dbSource)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	tx, err := conn.Begin(context.Background())
	if err != nil {
		panic(err)
	}
	defer tx.Rollback(context.Background())

	fmt.Println("Starting database sync...")

	// Process permissions
	fmt.Println("\nProcessing permissions...")
	for _, perm := range config.Permissions {
		methodJSON, err := json.Marshal(perm.Method)
		if err != nil {
			panic(err)
		}

		var exists bool
		err = tx.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM permissions WHERE name=$1 AND resource=$2 AND method=$3)", perm.Name, perm.Resource, string(methodJSON)).Scan(&exists)
		if err != nil {
			panic(err)
		}

		if !exists {
			_, err = tx.Exec(context.Background(), "INSERT INTO permissions (name, resource, method) VALUES ($1, $2, $3)", perm.Name, perm.Resource, string(methodJSON))
			if err != nil {
				panic(err)
			}
			fmt.Println("Inserted permission:", perm.Name)
		} else {
			fmt.Println("Permission already exists:", perm.Name)
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
			_, err = tx.Exec(context.Background(), "INSERT INTO roles (id, name) VALUES ($1, $2)", roleID, role.Name)
			if err != nil {
				panic(err)
			}
			fmt.Println("Inserted role:", role.Name)
		} else {
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
