package api

import (
	"context"
	"log"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/hub"
	"maicare_go/mocks"
	"maicare_go/notification"

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
var testb2Client *bucket.ObjectStorageClient
var testasynqClient *mocks.MockAsynqClientInterface
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

	testb2Client, err = bucket.NewObjectStorageClient(context.Background(), config)
	if err != nil {
		log.Fatal("cannot create b2 client:", err)
	}
	testasynqClient = mocks.NewMockAsynqClientInterface(testMockCtrl)

	hubInstance := hub.NewHub()

	testGrpcClient := CreateMockGrpcClient()
	testNotifService = notification.NewService(testStore, hubInstance)

	testServer, err = NewServer(testStore, testb2Client, testasynqClient, config.OpenRouterAPIKey, hubInstance, testNotifService, testGrpcClient)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	testServer.router.GET("/auth", testServer.AuthMiddleware(), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	os.Exit(m.Run())
}
