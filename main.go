package main

import (
	"context"
	"fmt"
	"log"
	"maicare_go/api"
	"maicare_go/async/aclient"
	"maicare_go/async/processor"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/email"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/hub"
	"maicare_go/logger"
	"maicare_go/notification"
	"maicare_go/service"
	"maicare_go/token"
	"maicare_go/util"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Add context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	poolConfig, err := pgxpool.ParseConfig(config.DbSource)
	if err != nil {
		log.Fatalf("unable to parse database config: %v", err)
	}
	// Configure connection pool settings
	poolConfig.MaxConns = 30                      // Maximum connections
	poolConfig.MinConns = 5                       // Minimum connections to keep open
	poolConfig.MaxConnLifetime = time.Hour        // Close connections after 1 hour
	poolConfig.MaxConnIdleTime = time.Minute * 30 // Close idle connections after 30 min
	poolConfig.HealthCheckPeriod = time.Minute    // Health check every minute
	poolConfig.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
		// Return true if connection is healthy
		return !conn.PgConn().IsClosed()
	}
	// Create the pool with the config
	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}

	store := db.NewStore(conn)
	b2Client, err := bucket.NewObjectStorageClient(ctx, config)
	if err != nil {
		log.Fatalf("unable to create b2 client: %v", err)
	}

	var asynqClient aclient.AsynqClientInterface
	if !config.Remote {
		asynqClient = aclient.NewAsynqClient(config.RedisHost, "", config.RedisPassword, nil)
	} else {
		asynqClient = aclient.NewAsynqClient(config.RedisHost, "", config.RedisPassword, nil)
	}

	// Inirialize the SMTP Client for email deleviry
	brevoConf := email.NewBrevoConf(config.BrevoSenderName, config.BrevoSenderEmail, config.BrevoApiKey)

	// Initialize the ws Hub
	hubInstance := hub.NewHub()

	// Initialize the notification service
	notificationService := notification.NewService(store, hubInstance)

	// Initialize Asynq server
	var asynqServer *processor.AsynqServer

	grpcClient, err := grpclient.NewGrpcClient(config.GrpcUrl)
	if err != nil {
		log.Fatalf("Could not create gRPC client: %v", err)
	}

	// Create error channel to catch server errors
	errChan := make(chan error, 1)

	// move this to services
	tokenMaker, err := token.NewJWTMaker(config.AccessTokenSecretKey, config.RefreshTokenSecretKey, config.TwoFATokenSecretKey)
	if err != nil {
		log.Fatalf("cannot create tokenmaker: %v", err)
	}

	// move this to services
	logger, err := logger.SetupLogger(config.Environment)
	if err != nil {
		log.Fatalf("cannot setup logger: %v", err)
	}

	// Init the buisness service
	businessService := service.NewBusinessService(store, tokenMaker, logger, &config, b2Client)

	if !config.Remote {
		redisClient := redis.NewClient(&redis.Options{
			Addr:      config.RedisHost, // e.g., "frankfurt-keyvalue.render.com:6379"
			Username:  "",               // if applicable
			Password:  config.RedisPassword,
			TLSConfig: nil, // Only if using TLS (rediss://)
		})

		maxAttempts := 5
		delay := time.Second // start with 1 second delay

		var pingErr error
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			rctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, pingErr = redisClient.Ping(rctx).Result()
			if pingErr == nil {
				log.Println("✅ Redis connection successful!")
				break
			}

			log.Printf("⚠️ Redis ping failed (attempt %d/%d): %v", attempt, maxAttempts, pingErr)

			// Wait before retrying
			time.Sleep(delay)
			delay *= 2 // exponential backoff
		}

		if pingErr != nil {
			log.Fatalf("❌ Failed to connect to Redis after %d attempts: %v", maxAttempts, pingErr)
		}
		asynqServer = processor.NewAsynqServer(config.RedisHost, "", config.RedisPassword, store, nil, brevoConf, b2Client, notificationService, businessService)
	} else {
		redisClient := redis.NewClient(&redis.Options{
			Addr:      config.RedisHost, // e.g., "frankfurt-keyvalue.render.com:6379"
			Username:  "",               // if applicable
			Password:  config.RedisPassword,
			TLSConfig: nil, // Only if using TLS (rediss://)
		})

		maxAttempts := 5
		delay := time.Second // start with 1 second delay

		var pingErr error
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			rctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, pingErr = redisClient.Ping(rctx).Result()
			if pingErr == nil {
				log.Println("✅ Redis connection successful!")
				break
			}

			log.Printf("⚠️ Redis ping failed (attempt %d/%d): %v", attempt, maxAttempts, pingErr)

			// Wait before retrying
			time.Sleep(delay)
			delay *= 2 // exponential backoff
		}

		if pingErr != nil {
			log.Fatalf("❌ Failed to connect to Redis after %d attempts: %v", maxAttempts, pingErr)
		}
		asynqServer = processor.NewAsynqServer(config.RedisHost, "", config.RedisPassword, store, nil, brevoConf, b2Client, notificationService, businessService)
	}

	// Start the Asynq server in a goroutine
	go func() {
		log.Println("Starting Asynq server...")

		if err := asynqServer.Start(); err != nil {
			log.Printf("FATAL: Asynq server error: %v", err)
			errChan <- fmt.Errorf("asynq server error: %v", err)
		}
	}()

	log.Println("Asynq server started successfully in background")

	// Start your main server
	server, err := api.NewServer(store, b2Client, asynqClient,
		config.OpenRouterAPIKey, hubInstance, notificationService,
		grpcClient, tokenMaker, config, businessService)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	go func() {
		if err := server.Start(); err != nil {
			errChan <- fmt.Errorf("main server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Wait for either error or shutdown signal
	select {
	case err := <-errChan:
		log.Printf("Server error: %v", err)
	case sig := <-quit:
		log.Printf("Received signal: %v", sig)
	}

	log.Println("Shutting down servers...")

	// Add timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown both servers with timeout
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		asynqServer.Shutdown()
	}()

	go func() {
		defer wg.Done()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	// Wait for both servers to shutdown or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-shutdownCtx.Done():
		log.Println("Shutdown timed out")
	case <-done:
		log.Println("Servers shut down successfully")
	}
}
