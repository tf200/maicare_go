package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
)

const (
	TypeEmailDelivery = "email:deliver"
	TypeImageResize   = "image:resize"
)

type EmailDeliveryPayload struct {
	UserID     int
	TemplateID string
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
