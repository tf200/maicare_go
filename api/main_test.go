package api

import (
	"context"
	"log"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/util"
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore *db.Store
var testServer *Server
var testb2Client *bucket.B2Client

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
	testb2Client, err = bucket.NewB2Client(config)
	if err != nil {
		log.Fatal("cannot create b2 client:", err)
	}

	testServer, err = NewServer(testStore, testb2Client)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	testServer.router.GET("/auth", AuthMiddleware(testServer.tokenMaker), func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{})
	})

	os.Exit(m.Run())
}
