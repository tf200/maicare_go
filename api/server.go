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

// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @description Add 'Bearer ' prefix before your JWT token for authentication

// @Security Bearer
import (
	"context"
	"fmt"
	"log"
	"net/http"

	"maicare_go/ai"
	"maicare_go/async"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/docs"
	"maicare_go/hub"
	"maicare_go/token"
	"maicare_go/util"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	store       *db.Store
	router      *gin.Engine
	config      util.Config
	tokenMaker  token.Maker
	b2Client    *bucket.B2Client
	asynqClient *async.AsynqClient
	httpServer  *http.Server
	aiHandler   *ai.AiHandler
	hub         *hub.Hub
	logger      *zap.Logger
}

func NewServer(store *db.Store, b2Client *bucket.B2Client, asyqClient *async.AsynqClient, apiKey string, hubInstance *hub.Hub) (*Server, error) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		return nil, fmt.Errorf("cannot load env %v", err)
	}

	tokenMaker, err := token.NewJWTMaker(config.AccessTokenSecretKey, config.RefreshTokenSecretKey, config.TwoFATokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create tokenmaker %v", err)
	}

	aiHandler := ai.NewAiHandler(apiKey)

	logger, err := setupLogger(config.Environment)
	if err != nil {
		return nil, fmt.Errorf("cannot create logger %v", err)
	}

	server := &Server{
		store:       store,
		config:      config,
		tokenMaker:  tokenMaker,
		b2Client:    b2Client,
		asynqClient: asyqClient,
		aiHandler:   aiHandler,
		hub:         hubInstance,
		logger:      logger,
	}

	// Initialize swagger docs
	docs.SwaggerInfo.Title = "Maicare API"
	docs.SwaggerInfo.Description = "This is the Maicare server API documentation."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = server.config.Host // This will use your configured server address
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	server.setupRoutes()
	return server, nil
}

func (server *Server) setupRoutes() {
	gin.SetMode(func() string {
		if server.config.Environment == "production" {
			return gin.ReleaseMode
		}
		return gin.DebugMode
	}())
	router := gin.New()

	corsConf := cors.DefaultConfig()
	corsConf.AllowOrigins = []string{"*"}
	corsConf.AllowCredentials = true
	corsConf.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConf.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}

	router.Use(cors.New(corsConf))
	router.Use(server.requestLogger())
	router.Use(gin.Recovery())
	// Add swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	baseRouter := router.Group("/")

	// Setup routes from different modules
	server.setupTestRoutes(baseRouter)
	server.setupAuthRoutes(baseRouter)
	server.setupRolesRoutes(baseRouter)
	server.setupEmployeeRoutes(baseRouter)
	server.setupLocationRoutes(baseRouter)
	server.setupAttachementRoutes(baseRouter)
	server.setupSenderRoutes(baseRouter)
	server.setupClientRoutes(baseRouter)
	server.setupClientMedicalRoutes(baseRouter)
	server.setupClientNetworkRoutes(baseRouter)
	server.setupClientIncidentRoutes(baseRouter)
	server.setupAiRoutes(baseRouter)
	server.setupProgressReportsRoutes(baseRouter)
	server.setupAppointmentCardRoutes(baseRouter)
	server.setupMaturityMatrixRoutes(baseRouter)
	server.setupIntakeFormRoutes(baseRouter)
	server.setupContractRoutes(baseRouter)
	server.setupECRRoutes(baseRouter)
	server.setupAppointmentRoutes(baseRouter)
	server.setupIncidentsAllRoutes(baseRouter)
	server.setupRegistrationFormRoutes(baseRouter)
	server.setupScheduleRoutes(baseRouter)
	server.setupShiftsRoutes(baseRouter)
	server.setupWorkingHours(baseRouter)
	server.setupInvoiceRoutes(baseRouter)
	// Add more route setups as needed

	server.setupWebsocketRoutes(baseRouter)

	server.router = router
}

func (server *Server) Start() error {
	// --- Start the Hub's processing loop ---
	// Start this in a goroutine BEFORE starting the blocking HTTP server listener.
	go server.hub.Run() // <-- START HUB HERE
	log.Println("WebSocket Hub processing loop started")
	// ---

	// Create http.Server (if not already done in NewServer, though your code does it here)
	server.httpServer = &http.Server{
		Addr:    server.config.ServerAddress,
		Handler: server.router,
		// ... other http.Server configurations ...
	}

	// Start the HTTP server listener (this is blocking)
	log.Printf("Starting HTTP server on %s", server.config.ServerAddress)
	if err := server.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		// Log the error from ListenAndServe failing
		log.Printf("Failed to start HTTP server: %v", err)
		// Consider if you need to signal the hub to shut down here,
		// although currently hub.Run() doesn't have a shutdown mechanism.
		return fmt.Errorf("failed to start server: %v", err)
	}

	log.Println("HTTP server stopped gracefully or encountered an error.")
	return nil
}

func (server *Server) Shutdown(ctx context.Context) error {
	var httpErr error
	shutdownComplete := make(chan struct{})

	log.Println("Initiating server shutdown...")

	// Signal the WebSocket Hub to shut down first
	// This happens quickly; the actual cleanup runs in the Hub's goroutine.
	server.hub.Shutdown()
	log.Println("Hub shutdown signaled.")

	// Shutdown HTTP server concurrently
	go func() {
		log.Println("Shutting down HTTP server...")
		if err := server.httpServer.Shutdown(ctx); err != nil {
			httpErr = fmt.Errorf("http server shutdown failed: %v", err)
			log.Printf("HTTP server shutdown error: %v", err)
		} else {
			log.Println("HTTP server shut down successfully.")
		}
		close(shutdownComplete)
	}()

	// Wait for HTTP shutdown to complete or context deadline
	select {
	case <-ctx.Done():
		log.Printf("Shutdown context timed out/canceled: %v", ctx.Err())
		// Return the context error, potentially masking the httpErr if it also occurred
		return ctx.Err()
	case <-shutdownComplete:
		log.Println("HTTP server shutdown process finished.")
		// Return any error encountered during HTTP shutdown
		return httpErr
	}
}

func setupLogger(environment string) (*zap.Logger, error) {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()

		config.DisableCaller = true
		config.DisableStacktrace = true

		config.OutputPaths = []string{
			"stdout",
			"/var/log/maicare/app.log",
		}
	} else {
		config = zap.NewDevelopmentConfig()
		config.OutputPaths = []string{
			"stdout",
		}
	}

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %v", err)
	}

	if environment == "production" {
		fileWritter := &lumberjack.Logger{
			Filename:   "/var/log/maicare/app.log",
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28,   // days
			Compress:   true, // compress log files
		}

		core := zapcore.NewTee(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(config.EncoderConfig),
				zapcore.AddSync(fileWritter),
				zap.InfoLevel,
			),
			logger.Core(),
		)
		logger = zap.New(core)
	}

	return logger, nil

}
