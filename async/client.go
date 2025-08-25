package async

import (
	"context"
	"crypto/tls"
	"maicare_go/notification"
	"time"

	"github.com/hibiken/asynq"
)

//go:generate mockgen -source=client.go -destination=../mocks/mock_asynq_client.go -package=mocks
type AsynqClientInterface interface {
	EnqueueEmailDelivery(
		payload EmailDeliveryPayload,
		ctx context.Context,
		opts ...asynq.Option) error
	EnqueueIncident(
		payload IncidentPayload,
		ctx context.Context,
		opts ...asynq.Option) error
	EnqueueNotificationTask(
		ctx context.Context,
		payload notification.NotificationPayload,
		opts ...asynq.Option) error
	EnqueueAppointmentTask(
		ctx context.Context,
		payload AppointmentPayload,
		opts ...asynq.Option) error
	EnqueueAcceptedRegistration(
		ctx context.Context,
		payload AcceptedRegistrationFormPayload,
		opts ...asynq.Option) error
	GetClient() *asynq.Client
	Close() error
}

// AsynqClient wraps the asynq.Client.
type AsynqClient struct {
	client *asynq.Client
}

// NewAsynqClient creates a new AsynqClient instance.
func NewAsynqClient(redisHost, redisUser, redisPassword string, tls *tls.Config) AsynqClientInterface {
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
