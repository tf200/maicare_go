package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"maicare_go/config"
	db "maicare_go/db/sqlc"
	"maicare_go/docs"
	"maicare_go/internal/adapters"
	"maicare_go/internal/audit"
	"maicare_go/internal/domain"
	"maicare_go/internal/handler"
	"maicare_go/internal/middleware"
	"maicare_go/internal/repository"
	"maicare_go/internal/service"
	"maicare_go/internal/worker"
	"maicare_go/internal/ws"
	pkgasynq "maicare_go/pkg/asynq"
	pkgbucket "maicare_go/pkg/bucket"
	pkgemail "maicare_go/pkg/email"
	pkgjwt "maicare_go/pkg/jwt"
	pkglogger "maicare_go/pkg/logger"
	pkgpdf "maicare_go/pkg/pdf"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	cfg           config.Config
	dbPool        *pgxpool.Pool
	store         *db.Store
	logger        domain.Logger
	loggerSync    interface{ Sync() error }
	taskQueue     domain.TaskQueue
	hub           *ws.Hub
	ticketManager *ws.TicketManager
	httpServer    *http.Server
	workerServer  *worker.AsynqServer
	scheduler     *worker.Scheduler
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	if err := runMigrations(cfg.DbSource, cfg.MigrationsPath); err != nil {
		return nil, err
	}

	dbPool, err := newDBPool(ctx, cfg.DbSource)
	if err != nil {
		return nil, err
	}
	store := db.NewStore(dbPool)

	appLogger, err := pkglogger.Setup(cfg.Environment)
	if err != nil {
		return nil, fmt.Errorf("setup logger: %w", err)
	}

	jwtMaker, err := pkgjwt.New(cfg.AccessTokenSecretKey, cfg.RefreshTokenSecretKey, cfg.TwoFATokenSecretKey)
	if err != nil {
		return nil, fmt.Errorf("create jwt maker: %w", err)
	}
	tokenMaker := adapters.NewJWTTokenMakerAdapter(jwtMaker)

	if err := pingRedis(ctx, cfg); err != nil {
		return nil, err
	}
	asynqClient := pkgasynq.NewClient(cfg.RedisHost, "", cfg.RedisPassword, nil)
	taskQueue := adapters.NewTaskQueueAdapter(asynqClient)

	brevoConf := pkgemail.NewBrevoConf(cfg.BrevoSenderName, cfg.BrevoSenderEmail, cfg.BrevoApiKey)
	hub := ws.NewHub()
	ticketManager, err := ws.NewTicketManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("create websocket ticket manager: %w", err)
	}

	storage, bucketClient, err := newStorage(ctx, cfg)
	if err != nil {
		return nil, err
	}
	pdfService := pkgpdf.NewPdfService(bucketClient)
	incidentPDFGenerator := adapters.NewIncidentPDFGeneratorAdapter(pdfService)
	pdfSvc := adapters.NewPDFServiceAdapter(pdfService)
	aiService := adapters.NewAIServiceStub()

	notificationSvc, handlers := wireServicesAndHandlers(store, appLogger, tokenMaker, taskQueue, storage, incidentPDFGenerator, pdfSvc, aiService, hub, ticketManager, &cfg)

	router := newRouter(cfg, appLogger, tokenMaker, store, handlers)

	docs.SwaggerInfo.Title = "Maicare API"
	docs.SwaggerInfo.Description = "This is the Maicare server API documentation."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = cfg.Host
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	var workerBucket pkgbucket.ObjectStorageClient
	if bucketClient != nil {
		workerBucket = *bucketClient
	}

	return &App{
		cfg:           cfg,
		dbPool:        dbPool,
		store:         store,
		logger:        appLogger,
		loggerSync:    loggerSync(appLogger),
		taskQueue:     taskQueue,
		hub:           hub,
		ticketManager: ticketManager,
		httpServer: &http.Server{
			Addr:    cfg.ServerAddress,
			Handler: router,
		},
		workerServer: worker.NewAsynqServer(cfg.RedisHost, "", cfg.RedisPassword, store, nil, brevoConf, workerBucket, notificationSvc, pdfService),
		scheduler:    worker.NewScheduler(cfg.RedisHost, "", cfg.RedisPassword, nil),
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errChan := make(chan error, 3)

	go a.hub.Run()
	go func() {
		log.Println("Starting Asynq server...")
		if err := a.workerServer.Start(); err != nil {
			errChan <- fmt.Errorf("asynq server error: %w", err)
		}
	}()
	go func() {
		log.Println("Starting Asynq scheduler...")
		if err := a.scheduler.Start(); err != nil {
			errChan <- fmt.Errorf("asynq scheduler error: %w", err)
		}
	}()
	go func() {
		log.Printf("Starting HTTP server on %s", a.cfg.ServerAddress)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("http server error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	case sig := <-quit:
		log.Printf("Received signal: %v", sig)
	case <-ctx.Done():
		log.Printf("Context canceled: %v", ctx.Err())
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return a.Shutdown(shutdownCtx)
}

func (a *App) Shutdown(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() { defer wg.Done(); a.workerServer.Shutdown() }()
	go func() { defer wg.Done(); a.scheduler.Shutdown() }()
	go func() { defer wg.Done(); _ = a.httpServer.Shutdown(ctx) }()

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	a.hub.Shutdown()
	if a.ticketManager != nil {
		_ = a.ticketManager.Close()
	}
	if a.taskQueue != nil {
		_ = a.taskQueue.Close()
	}
	if a.dbPool != nil {
		a.dbPool.Close()
	}
	if a.loggerSync != nil {
		_ = a.loggerSync.Sync()
	}
	return nil
}

func runMigrations(dbSource, migrationsPath string) error {
	m, err := migrate.New(migrationsPath, dbSource)
	if err != nil {
		return fmt.Errorf("create migrate instance: %w", err)
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}
	return nil
}

func newDBPool(ctx context.Context, source string) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(source)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}
	poolConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error { return db.RegisterEnumTypes(ctx, conn) }
	poolConfig.MaxConns = 30
	poolConfig.MinConns = 5
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute
	poolConfig.HealthCheckPeriod = time.Minute
	poolConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool { return !conn.PgConn().IsClosed() }

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("connect database: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return pool, nil
}

func pingRedis(ctx context.Context, cfg config.Config) error {
	client := redis.NewClient(&redis.Options{Addr: cfg.RedisHost, Password: cfg.RedisPassword})
	defer func() { _ = client.Close() }()

	var lastErr error
	for attempt := 1; attempt <= 5; attempt++ {
		pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		_, lastErr = client.Ping(pingCtx).Result()
		cancel()
		if lastErr == nil {
			return nil
		}
		time.Sleep(time.Duration(attempt) * time.Second)
	}
	return fmt.Errorf("connect redis: %w", lastErr)
}

func newStorage(ctx context.Context, cfg config.Config) (domain.Storage, *pkgbucket.ObjectStorageClient, error) {
	if cfg.DisableBucket {
		return pkgbucket.NewNoop(), nil, nil
	}
	client, err := pkgbucket.New(ctx, pkgbucket.Config{Endpoint: cfg.B2Endpoint, KeyID: cfg.B2KeyID, Key: cfg.B2Key, Bucket: cfg.B2Bucket, Region: "eu-central-003", Secure: true})
	if err != nil {
		return nil, nil, fmt.Errorf("create object storage: %w", err)
	}
	return client, client, nil
}

func loggerSync(logger domain.Logger) interface{ Sync() error } {
	if syncer, ok := logger.(interface{ Sync() error }); ok {
		return syncer
	}
	return nil
}

type appHandlers struct {
	auth             *handler.AuthHandler
	attachment       *handler.AttachmentHandler
	client           *handler.ClientHandler
	contract         *handler.ContractHandler
	dashboard        *handler.DashboardHandler
	employee         *handler.EmployeeHandler
	handbook         *handler.HandbookHandler
	incident         *handler.IncidentHandler
	intakeForm       *handler.IntakeFormHandler
	lateArrival      *handler.LateArrivalHandler
	leave            *handler.LeaveHandler
	maturityMatrix   *handler.MaturityMatrixHandler
	notification     *handler.NotificationHandler
	organization     *handler.OrganizationHandler
	registrationForm *handler.RegistrationFormHandler
	role             *handler.RoleHandler
	schedule         *handler.ScheduleHandler
	event            *handler.EventHandler
	sender           *handler.SenderHandler
	settings         *handler.SettingsHandler
	shiftSwap        *handler.ShiftSwapHandler
	websocket        *handler.WebSocketHandler
	invoice          *handler.InvoiceHandler
}

func wireServicesAndHandlers(store *db.Store, logger domain.Logger, tokenMaker domain.TokenMaker, taskQueue domain.TaskQueue, storage domain.Storage, incidentPDFGenerator domain.IncidentPDFGenerator, pdfService domain.PDFService, aiService *adapters.AIServiceStub, hub *ws.Hub, ticketManager *ws.TicketManager, cfg *config.Config) (domain.NotificationService, appHandlers) {
	authRepo := repository.NewAuthRepository(store)
	attachmentRepo := repository.NewAttachmentRepository(store)
	clientRepo := repository.NewClientRepository(store)
	contractRepo := repository.NewContractRepository(store)
	dashboardRepo := repository.NewDashboardRepository(store)
	employeeRepo := repository.NewEmployeeRepository(store)
	handbookRepo := repository.NewHandbookRepository(store)
	incidentRepo := repository.NewIncidentRepository(store)
	intakeRepo := repository.NewIntakeFormRepository(store)
	lateArrivalRepo := repository.NewLateArrivalRepository(store)
	leaveRepo := repository.NewLeaveRepository(store)
	maturityRepo := repository.NewMaturityMatrixRepository(store)
	notificationRepo := repository.NewNotificationRepository(store)
	organizationRepo := repository.NewOrganizationRepository(store)
	registrationRepo := repository.NewRegistrationFormRepository(store)
	roleRepo := repository.NewRoleRepository(store)
	scheduleRepo := repository.NewScheduleRepository(store)
	senderRepo := repository.NewSenderRepository(store)
	departmentRepo := repository.NewDepartmentRepository(store)
	organizationProfileRepo := repository.NewOrganizationProfileRepository(store)

	notificationSvc := service.NewNotificationService(notificationRepo, hub, logger)
	scheduleSvc := service.NewScheduleService(scheduleRepo, notificationSvc, logger)
	eventSvc := service.NewEventService(store, taskQueue, logger)

	// NEN 7513 audit logger — shared across all services
	auditLogger := audit.New(store, logger)

	authSvc := service.NewAuthService(authRepo, tokenMaker, logger, auditLogger, cfg.AccessTokenDuration, cfg.RefreshTokenDuration, cfg.TwoFATokenDuration)
	attachmentSvc := service.NewAttachmentService(attachmentRepo, storage, logger)
	clientSvc := service.NewClientService(clientRepo, taskQueue, storage, aiService, pdfService, logger, auditLogger)
	contractSvc := service.NewContractService(contractRepo, storage, logger)
	dashboardSvc := service.NewDashboardService(dashboardRepo, logger, auditLogger)
	employeeSvc := service.NewEmployeeService(employeeRepo, logger)
	handbookSvc := service.NewHandbookService(handbookRepo, logger)
	incidentSvc := service.NewIncidentService(incidentRepo, incidentPDFGenerator, taskQueue, logger)
	intakeSvc := service.NewIntakeFormService(intakeRepo, logger, aiService, auditLogger)
	lateArrivalSvc := service.NewLateArrivalService(lateArrivalRepo, logger)
	leaveSvc := service.NewLeaveService(leaveRepo, logger)
	maturitySvc := service.NewMaturityMatrixService(maturityRepo)
	organizationSvc := service.NewOrganizationService(organizationRepo, logger)
	registrationSvc := service.NewRegistrationFormService(registrationRepo, logger, taskQueue, auditLogger)
	roleSvc := service.NewRoleService(roleRepo, logger)
	senderSvc := service.NewSenderService(senderRepo, logger)
	departmentSvc := service.NewDepartmentService(departmentRepo)
	organizationProfileSvc := service.NewOrganizationProfileService(organizationProfileRepo)

	invoiceSvc := service.NewInvoiceService(store, logger, storage, pdfService)

	return notificationSvc, appHandlers{
		auth:             handler.NewAuthHandler(authSvc),
		attachment:       handler.NewAttachmentHandler(attachmentSvc),
		client:           handler.NewClientHandler(clientSvc),
		contract:         handler.NewContractHandler(contractSvc),
		dashboard:        handler.NewDashboardHandler(dashboardSvc),
		employee:         handler.NewEmployeeHandler(employeeSvc),
		handbook:         handler.NewHandbookHandler(handbookSvc),
		incident:         handler.NewIncidentHandler(incidentSvc),
		intakeForm:       handler.NewIntakeFormHandler(intakeSvc),
		lateArrival:      handler.NewLateArrivalHandler(lateArrivalSvc),
		leave:            handler.NewLeaveHandler(leaveSvc),
		maturityMatrix:   handler.NewMaturityMatrixHandler(maturitySvc),
		notification:     handler.NewNotificationHandler(notificationSvc),
		organization:     handler.NewOrganizationHandler(organizationSvc),
		registrationForm: handler.NewRegistrationFormHandler(registrationSvc),
		role:             handler.NewRoleHandler(roleSvc),
		schedule:         handler.NewScheduleHandler(scheduleSvc),
		event:            handler.NewEventHandler(eventSvc),
		sender:           handler.NewSenderHandler(senderSvc),
		settings:         handler.NewSettingsHandler(departmentSvc, organizationProfileSvc),
		shiftSwap:        handler.NewShiftSwapHandler(scheduleSvc),
		websocket:        handler.NewWebSocketHandler(hub, ticketManager, cfg),
		invoice:          handler.NewInvoiceHandler(invoiceSvc),
	}
}

func newRouter(cfg config.Config, logger domain.Logger, tokenMaker domain.TokenMaker, store *db.Store, handlers appHandlers) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	router.Use(corsMiddleware())
	router.Use(middleware.NewRequestContextMiddleware(logger).Handle())
	router.Use(middleware.NewRequestLoggingMiddleware(logger, cfg.Environment).Handle())
	router.Use(gin.Recovery())

	auth := middleware.NewAuthMiddleware(tokenMaker, logger).Handle()
	rbac := middleware.NewRBACMiddleware(rolePermissionChecker{store: store}, logger)
	requirePermission := rbac.Require
	base := router.Group("/")

	handler.RegisterAuthRoutes(base, handlers.auth)
	handler.RegisterAttachmentRoutes(base, handlers.attachment, auth)
	handler.RegisterClientRoutes(base, handlers.client, auth, requirePermission)
	handler.RegisterEvaluationRoutes(base, handlers.client, auth, requirePermission)
	handler.RegisterContractRoutes(base, handlers.contract, auth, requirePermission)
	handler.RegisterDashboardRoutes(base, handlers.dashboard, auth, requirePermission)
	handler.RegisterEmployeeRoutes(base, handlers.employee, auth, requirePermission)
	handler.RegisterHandbookRoutes(base, handlers.handbook, auth, requirePermission)
	handler.RegisterIncidentRoutes(base, handlers.incident, auth, requirePermission)
	handler.RegisterIntakeFormRoutes(base, handlers.intakeForm, auth, requirePermission)
	handler.RegisterLateArrivalRoutes(base, handlers.lateArrival, auth, requirePermission)
	handler.RegisterLeaveRoutes(base, handlers.leave, auth, requirePermission)
	handler.RegisterMaturityMatrixRoutes(base, handlers.maturityMatrix, auth, requirePermission)
	handler.RegisterNotificationRoutes(base, handlers.notification, auth)
	handler.RegisterOrganizationRoutes(base, handlers.organization, auth, requirePermission)
	handler.RegisterRegistrationFormRoutes(base, handlers.registrationForm, auth, requirePermission)
	handler.RegisterRoleRoutes(base, handlers.role, auth, requirePermission)
	handler.RegisterScheduleRoutes(base, handlers.schedule, auth, requirePermission)
	handler.RegisterEventRoutes(base, handlers.event, auth, requirePermission)
	handler.RegisterSenderRoutes(base, handlers.sender, auth, requirePermission)
	handler.RegisterSettingsRoutes(base, handlers.settings, auth, requirePermission)
	handler.RegisterShiftSwapRoutes(base, handlers.shiftSwap, auth, requirePermission)
	handler.RegisterWebSocketRoutes(base, handlers.websocket, auth)
	handler.RegisterInvoiceRoutes(base, handlers.invoice, auth, requirePermission)

	return router
}

func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Credentials", "true")
		ctx.Header("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}
		ctx.Next()
	}
}

type rolePermissionChecker struct{ store *db.Store }

func (c rolePermissionChecker) HasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error) {
	perms, err := repository.NewRoleRepository(c.store).ListEffectiveUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}
	for _, perm := range perms {
		if perm.PermissionName == permission || perm.Resource == permission {
			return true, nil
		}
	}
	return false, nil
}

func (c rolePermissionChecker) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	roles, err := repository.NewRoleRepository(c.store).GetUserRoles(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(roles))
	for i, role := range roles {
		result[i] = role.Name
	}
	return result, nil
}
