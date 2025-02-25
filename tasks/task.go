package tasks

import (
	"context"
	"fmt"
	"log"

	"github.com/goccy/go-json"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

type EmailDeliveryPayload struct {
	To           string `json:"to"`
	UserEmail    string `json:"user_email"`
	UserPassword string `json:"user_password"`
}

func (c *AsynqClient) EnqueueEmailDelivery(
	payload EmailDeliveryPayload,
	ctx context.Context,
	opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %v", err)
	}
	task := asynq.NewTask(TypeEmailDelivery, jsonPayload)
	info, err := c.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("client.EnqueueContext failed: %v", err)
	}
	log.Printf("task enqueued: id=%s queue=%s", info.ID, info.Queue)
	return nil

}
