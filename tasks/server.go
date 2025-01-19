package tasks

import (
	"crypto/tls"
	db "maicare_go/db/sqlc"
	"time"

	"github.com/hibiken/asynq"
)

type AsynqServer struct {
	server *asynq.Server
	store  *db.Store
}

func NewAsynqServer(redisHost, redisUser, redisPassword string, store *db.Store, tls *tls.Config) *AsynqServer {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:         redisHost,
			Username:     redisUser,
			Password:     redisPassword,
			TLSConfig:    tls,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			RetryDelayFunc: func(n int, err error, t *asynq.Task) time.Duration {
				return time.Duration(n*n) * time.Second // Exponential backoff
			},
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
