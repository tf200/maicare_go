package testdatabase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
)

// StartPostgres starts and migrates an isolated PostgreSQL container. An
// already-migrated TEST_DATABASE_URL can be supplied for local debugging.
func StartPostgres(ctx context.Context) (string, func() error, error) {
	if databaseURL := os.Getenv("TEST_DATABASE_URL"); databaseURL != "" {
		return databaseURL, func() error { return nil }, nil
	}

	container, err := postgrescontainer.Run(ctx,
		"postgres:17-alpine",
		postgrescontainer.WithDatabase("maicare_test"),
		postgrescontainer.WithUsername("postgres"),
		postgrescontainer.WithPassword("postgres"),
		postgrescontainer.BasicWaitStrategies(),
	)
	if err != nil {
		return "", nil, fmt.Errorf("start PostgreSQL test container: %w", err)
	}
	cleanup := func() error { return testcontainers.TerminateContainer(container) }

	databaseURL, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = cleanup()
		return "", nil, fmt.Errorf("get PostgreSQL test container connection string: %w", err)
	}
	if err := migrateDatabase(databaseURL); err != nil {
		_ = cleanup()
		return "", nil, err
	}

	return databaseURL, cleanup, nil
}

func migrateDatabase(databaseURL string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("resolve integration test source path")
	}
	migrationsPath := "file://" + filepath.ToSlash(filepath.Join(filepath.Dir(filename), "..", "..", "db", "migrations"))

	migrator, err := migrate.New(migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer func() { _, _ = migrator.Close() }()

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
