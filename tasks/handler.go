package tasks

import (
	"context"
	"fmt"
	"log"
	"maicare_go/email"

	"github.com/goccy/go-json"

	"github.com/hibiken/asynq"
)

func (processor *AsynqServer) ProcessEmailTask(ctx context.Context, t *asynq.Task) error {
	var p EmailDeliveryPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("Failed to unmarshal email task payload: %v", err)
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if p.To == "" || p.UserEmail == "" || p.UserPassword == "" {
		return fmt.Errorf("invalid email payload: missing required fields: %w", asynq.SkipRetry)
	}

	err := processor.smtp.SendCredentials(ctx, []string{p.To}, email.Credentials{Email: p.UserEmail, Password: p.UserPassword})
	if err != nil {
		log.Printf("Failed to send email to %s: %v", p.To, err)
		return fmt.Errorf("failed to send email to %s: %v: %w", p.To, err, asynq.SkipRetry)
	}

	return nil
}
