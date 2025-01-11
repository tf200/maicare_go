package tasks

import "github.com/hibiken/asynq"

type AsynqClient struct {
	client *asynq.Client
}

func NewAsynqClient(redisAddr string) *AsynqClient {
	client := asynq.NewClient(
		asynq.RedisClientOpt{Addr: redisAddr},
	)
	return &AsynqClient{client: client}
}
