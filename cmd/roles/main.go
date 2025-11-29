package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

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

func loadEnvVariable(filename, key string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k == key {
			return v, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("key %s not found in %s", key, filename)
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

	dbSource, err := loadEnvVariable("app.env", "DB_SOURCE")
	if err != nil {
		panic(err)
	}

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
