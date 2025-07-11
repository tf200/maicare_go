package invoice

import (
	"context"
	"log"
	db "maicare_go/db/sqlc"

	"maicare_go/util"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore *db.Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalf("Could not load conf %v", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DbSource)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	testStore = db.NewStore(conn)

	os.Exit(m.Run())
}
