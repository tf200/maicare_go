package async

import (
	"context"
	"crypto/tls"
	"log"
	bucketmocks "maicare_go/bucket/mocks"
	db "maicare_go/db/sqlc"
	"maicare_go/email"
	"maicare_go/hub"
	"maicare_go/logger"
	"maicare_go/notification"
	"maicare_go/service"
	"maicare_go/token"

	"maicare_go/util"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/mock/gomock"
)

var testStore *db.Store
var testasynqClient AsynqClientInterface
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

	ctrl := gomock.NewController(&testing.T{})
	defer ctrl.Finish()

	testB2client := bucketmocks.NewMockObjectStorageInterface(ctrl)
	if err != nil {
		log.Fatalf("unable to create b2 client: %v", err)
	}

	testStore = db.NewStore(conn)
	testBrevo := email.NewBrevoConf(config.BrevoApiKey, config.BrevoSenderName, config.BrevoSenderEmail)
	testasynqClient = NewAsynqClient(config.RedisHost, "", config.RedisPassword, &tls.Config{})
	hubInstance := hub.NewHub()
	testNotifService := notification.NewService(testStore, hubInstance)

	ctrl = gomock.NewController(&testing.T{})
	defer ctrl.Finish()

	buisnessService := service.NewBusinessService(testStore, &token.JWTMaker{}, &logger.LoggerImpl{}, &config, testB2client)

	testWorker = NewAsynqServer(config.RedisHost, "", config.RedisPassword, testStore, &tls.Config{}, testBrevo, testB2client, testNotifService, buisnessService)

	os.Exit(m.Run())
}
