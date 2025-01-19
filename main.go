package main

import (
	"context"
	"fmt"
	"log"
	"maicare_go/api"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/tasks"
	"maicare_go/util"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Add context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.New(ctx, config.DbSource)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	b2Client, err := bucket.NewB2Client(config)
	if err != nil {
		log.Fatalf("unable to create b2 client: %v", err)
	}

	// Initialize Asynq client
	asynqClient := tasks.NewAsynqClient(config.RedisHost, config.RedisUser, config.RedisPassword)

	// Initialize Asynq server
	asynqServer := tasks.NewAsynqServer(config.RedisHost, config.RedisUser, config.RedisPassword, store)

	// Create error channel to catch server errors
	errChan := make(chan error, 1)

	// Start the Asynq server in a goroutine
	go func() {
		if err := asynqServer.Start(); err != nil {
			errChan <- fmt.Errorf("asynq server error: %v", err)
		}
	}()

	// Start your main server
	server, err := api.NewServer(store, b2Client, asynqClient)
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
