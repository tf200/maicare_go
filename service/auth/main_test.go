package auth

// var testAuthService AuthService

// func TestMain(m *testing.M) {
// 	config, err := util.LoadConfig("../..")
// 	if err != nil {
// 		panic(err)
// 	}

// 	tokenMaker, err := token.NewJWTMaker(config.AccessTokenSecretKey,
// 		config.RefreshTokenSecretKey,
// 		config.TwoFATokenSecretKey,
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	conn, err := pgxpool.New(context.Background(), config.DbSource)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close()

// 	testStore := db.NewStore(conn)

// 	mockLogger, err := logger.SetupLogger("development")
// 	if err != nil {
// 		panic(err)
// 	}

// 	deps := deps.NewServiceDependencies(testStore, tokenMaker, mockLogger, &config)
// 	testAuthService = NewAuthService(deps)
// 	os.Exit(m.Run())
// }
