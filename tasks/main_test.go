package tasks

import (
	"context"
	"crypto/tls"
	"log"
	db "maicare_go/db/sqlc"
	"maicare_go/email"

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
	testSmtp := email.NewSmtpConf(config.SmtpName, config.SmtpAddress, config.SmtpAuth, config.SmtpHost, config.SmtpPort)
	testasynqClient = NewAsynqClient(config.RedisHost, config.RedisUser, config.RedisPassword, &tls.Config{})

	testWorker = NewAsynqServer(config.RedisHost, config.RedisUser, config.RedisPassword, testStore, &tls.Config{}, testSmtp)

	os.Exit(m.Run())
}
