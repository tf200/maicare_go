package api

import (
	"fmt"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store      *db.Store
	router     *gin.Engine
	config     util.Config
	tokenMaker token.Maker
	b2Client   *bucket.B2Client
}

func NewServer(store *db.Store, b2Client *bucket.B2Client) (*Server, error) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		return nil, fmt.Errorf("cannot load env %v", err)
	}

	tokenMaker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create tokenmaker %v", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
		b2Client:   b2Client,
	}

	server.setupRoutes()
	return server, nil
}

func (server *Server) setupRoutes() {
	router := gin.Default()

	baseRouter := router.Group("/")

	// Setup routes from different modules
	server.setupAuthRoutes(baseRouter)
	server.setupEmployeeRoutes(baseRouter)
	server.setupLocationRoutes(baseRouter)
	server.setupAttachementRoutes(baseRouter)
	server.setupSenderRoutes(baseRouter)

	// Add more route setups as needed

	server.router = router
	
}

func (server *Server) Start() error {
	return server.router.Run(server.config.ServerAddress)
}
