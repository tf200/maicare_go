package auth

import (
	"context"
	db "maicare_go/db/sqlc"
	"maicare_go/mocks"
	"maicare_go/service/deps"
	"maicare_go/token"
	"maicare_go/util"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/mock/gomock"
)

var testAuthService AuthService

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		panic(err)
	}

	tokenMaker, err := token.NewJWTMaker(config.AccessTokenSecretKey,
		config.RefreshTokenSecretKey,
		config.TwoFATokenSecretKey,
	)
	if err != nil {
		panic(err)
	}
	conn, err := pgxpool.New(context.Background(), config.DbSource)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	testStore := db.NewStore(conn)
	testMockCtrl := gomock.NewController(&testing.T{})
	defer testMockCtrl.Finish()

	mockLogger := mocks.NewMockLogger(testMockCtrl)

	deps := deps.NewServiceDependencies(testStore, tokenMaker, mockLogger, &config)
	testAuthService = NewAuthService(deps)
	os.Exit(m.Run())
}
