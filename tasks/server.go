package tasks

import (
	db "maicare_go/db/sqlc"

	"github.com/hibiken/asynq"
)

type AsynqServer struct {
	server *asynq.Server
	store  *db.Store
}

func NewAsynqServer(redisAddr string, store *db.Store) *AsynqServer {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)
	return &AsynqServer{server: srv, store: store}
}

func (a *AsynqServer) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeEmailDelivery, a.ProcessEmailTask)
	return a.server.Start(mux)
}

func (a *AsynqServer) Shutdown() {
	a.server.Shutdown()
}
