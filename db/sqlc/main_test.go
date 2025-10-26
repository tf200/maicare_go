package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func runMigrations(dbSource string, migrationsPath string) error {
	log.Println("Running database migrations...")
	m, err := migrate.New(migrationsPath, dbSource)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply")
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("✅ Migrations applied successfully")
	return nil
}

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Start test container
	pgContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("maicare_test"),
		postgres.WithUsername("maicare"),
		postgres.WithPassword("password"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	// Run migrations
	if err := runMigrations(connStr, "file://../migrations"); err != nil {
		log.Fatalf("could not run migrations: %v", err)
	}

	// Connect to database
	testDB, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("unable to connect to db: %v", err)
	}

	testQueries = New(testDB)

	// Run tests
	code := m.Run()

	// Cleanup
	testDB.Close()
	if err := pgContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	os.Exit(code)
}
