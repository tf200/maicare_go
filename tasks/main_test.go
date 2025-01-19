package tasks

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
var testasynqClient *AsynqClient
var testWorker *AsynqServer


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
	testasynqClient = NewAsynqClient(config.RedisHost, config.RedisUser, config.RedisPassword)

	testWorker = NewAsynqServer(config.RedisHost, config.RedisUser, config.RedisPassword, testStore)

	os.Exit(m.Run())
}
