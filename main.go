package main

import (
	"context"
	"log"
	"maicare_go/api"
	"maicare_go/bucket"
	db "maicare_go/db/sqlc"
	"maicare_go/util"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DbSource)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer conn.Close()

	store := db.NewStore(conn)
	b2Client, err := bucket.NewB2Client(config)

	if err != nil {
		log.Fatalf("unable to create b2 client: %v", err)
	}

	server, err := api.NewServer(store, b2Client)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
