package tasks

import (
	"context"
	"fmt"
	"log"

	"github.com/goccy/go-json"

	"github.com/hibiken/asynq"
)

func (processor *AsynqServer) ProcessEmailTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	log.Printf("Sending Email to User: user_id=%d, template_id=%s", p.UserID, p.TemplateID)
	// Email delivery code ...
	return nil
}
