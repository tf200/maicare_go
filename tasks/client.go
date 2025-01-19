package tasks

import (
	"crypto/tls"
	"time"

	"github.com/hibiken/asynq"
)

type AsynqClient struct {
	client *asynq.Client
}

func NewAsynqClient(redisHost, redisUser, redisPassword string, tls *tls.Config) *AsynqClient {

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:         redisHost,
		Username:     redisUser,
		Password:     redisPassword,
		TLSConfig:    tls,
		DialTimeout:  5 * time.Second, // Connection timeout
		ReadTimeout:  3 * time.Second, // Read timeout
		WriteTimeout: 3 * time.Second,
		// Write timeout
		// Maximum number of retries
	})
	return &AsynqClient{client: client}
}
