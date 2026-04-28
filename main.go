package main

import (
	"context"
	"log"

	"maicare_go/config"
	"maicare_go/internal/app"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load(".")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("cannot create app: %v", err)
	}

	if err := application.Run(ctx); err != nil {
		log.Fatalf("app stopped with error: %v", err)
	}
}
