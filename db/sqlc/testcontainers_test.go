package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"maicare_go/internal/testdatabase"
)

var integrationDatabaseURL string

func TestMain(m *testing.M) {
	ctx := context.Background()
	databaseURL, cleanup, err := testdatabase.StartPostgres(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	integrationDatabaseURL = databaseURL

	exitCode := m.Run()
	if err := cleanup(); err != nil {
		fmt.Fprintf(os.Stderr, "terminate PostgreSQL test container: %v\n", err)
		exitCode = 1
	}
	os.Exit(exitCode)
}
