package async

import (
	"crypto/tls"
	"time"

	"github.com/hibiken/asynq"
)

// AsynqClient wraps the asynq.Client.
type AsynqClient struct {
	client *asynq.Client
}

// NewAsynqClient creates a new AsynqClient instance.
func NewAsynqClient(redisHost, redisUser, redisPassword string, tls *tls.Config) *AsynqClient {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr:         redisHost,
		Username:     redisUser,
		Password:     redisPassword,
		TLSConfig:    tls,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	return &AsynqClient{client: client}
}

// GetClient returns the underlying asynq.Client if direct access is needed.
func (c *AsynqClient) GetClient() *asynq.Client {
	return c.client
}

// Close closes the underlying client connection.
func (c *AsynqClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
