package async

import (
	"context"
	"crypto/tls"
	"log"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/email"
	"maicare_go/hub"
	"maicare_go/notification"

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
	testB2client, err := bucket.NewB2Client(config)
	if err != nil {
		log.Fatalf("unable to create b2 client: %v", err)
	}

	testStore = db.NewStore(conn)
	testSmtp := email.NewSmtpConf(config.SmtpName, config.SmtpAddress, config.SmtpAuth, config.SmtpHost, config.SmtpPort)
	testasynqClient = NewAsynqClient(config.RedisHost, config.RedisUser, config.RedisPassword, &tls.Config{})
	hubInstance := hub.NewHub()
	testNotifService := notification.NewService(testStore, hubInstance)

	testWorker = NewAsynqServer(config.RedisHost, config.RedisUser, config.RedisPassword, testStore, &tls.Config{}, testSmtp, testB2client, testNotifService)

	os.Exit(m.Run())
}
