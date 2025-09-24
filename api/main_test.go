package api

import (
	"context"
	"log"
	asyncmocks "maicare_go/async/mocks"
	bucketmocks "maicare_go/bucket/mocks"
	db "maicare_go/db/sqlc"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/hub"
	"maicare_go/logger"
	"maicare_go/notification"
	"maicare_go/service"
	"maicare_go/token"

	"maicare_go/util"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/mock/gomock"
)

var testStore *db.Store
var testServer *Server
var testb2Client *bucketmocks.MockObjectStorageInterface
var testasynqClient *asyncmocks.MockAsynqClientInterface
var testGrpcClient grpclient.GrpcClientInterface
var testNotifService *notification.Service
var testMockCtrl *gomock.Controller

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../")
	if err != nil {
		log.Fatalf("Could not load conf %v", err)
	}
	gin.SetMode(gin.TestMode)

	conn, err := pgxpool.New(context.Background(), config.DbSource)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	testStore = db.NewStore(conn)

	testMockCtrl = gomock.NewController(&testing.T{})
	defer testMockCtrl.Finish()

	testb2Client = bucketmocks.NewMockObjectStorageInterface(testMockCtrl)

	testasynqClient = asyncmocks.NewMockAsynqClientInterface(testMockCtrl)

	hubInstance := hub.NewHub()

	testGrpcClient := CreateMockGrpcClient()
	testNotifService = notification.NewService(testStore, hubInstance)

	tokenMaker, err := token.NewJWTMaker(config.AccessTokenSecretKey, config.RefreshTokenSecretKey, config.TwoFATokenSecretKey)
	if err != nil {
		log.Fatalf("cannot create tokenmaker: %v", err)
	}

	logger, err := logger.SetupLogger(config.Environment)
	if err != nil {
		log.Fatalf("cannot setup logger: %v", err)
	}

	businessService := service.NewBusinessService(testStore, tokenMaker, logger, &config, testb2Client)

	testServer, err = NewServer(testStore, testb2Client, testasynqClient, config.OpenRouterAPIKey,
		hubInstance, testNotifService, testGrpcClient,
		tokenMaker, config, businessService)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	testServer.router.GET("/auth", testServer.AuthMiddleware(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	os.Exit(m.Run())
}
