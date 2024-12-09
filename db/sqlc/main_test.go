package db

import (
	"context"
	"log"
	"maicare_go/util"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	var err error
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatalf("Unable to load Env Config: %v", err)
	}
	testDB, err = pgxpool.New(context.Background(), config.DbSource)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
