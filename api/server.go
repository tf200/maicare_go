package api
// @title Maicare API
// @version 1.0
// @description This is the Maicare server API documentation.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email your-email@domain.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
import (
	"fmt"

	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/docs"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Initialize swagger docs
	docs.SwaggerInfo.Title = "Maicare API"
	docs.SwaggerInfo.Description = "This is the Maicare server API documentation."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.ServerAddress // This will use your configured server address
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	server.setupRoutes()
	return server, nil
}

func (server *Server) setupRoutes() {
	router := gin.Default()

	// Add swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
