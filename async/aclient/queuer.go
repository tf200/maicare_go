package aclient

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"maicare_go/service/notification"

	"github.com/hibiken/asynq"
)

const (
	// Queue names (optional, but good practice if using multiple)
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"

	// Task Type Names
	TypeEmailDelivery                = "email:deliver"
	TypeIncidentProcess              = "incident:process" // Renamed for clarity
	TypeIncidentConfirmedEmail       = "incident:confirmed_email"
	TypeNotificationSend             = "notification:send"     // Renamed for clarity
	TypeAcceptedRegistration         = "accepted:registration" // Renamed for clarity
	TypeProcessRegistrationFormEmail = "email:process_registration_form"
)

func (c *AsynqClient) EnqueueProcessRegistrationFormEmail(
	ctx context.Context,
	payload ProcessRegistrationFormEmailPayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("EnqueueProcessRegistrationFormEmail: json.Marshal failed: %w", err)
	}

	if len(opts) == 0 {
		opts = append(opts, asynq.Queue(QueueDefault), asynq.MaxRetry(5))
	}

	task := asynq.NewTask(TypeProcessRegistrationFormEmail, jsonPayload)
	info, err := c.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("EnqueueProcessRegistrationFormEmail: client.EnqueueContext failed: %w", err)
	}

	log.Printf("Process Registration Form Email task enqueued: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

func (c *AsynqClient) EnqueueEmailDelivery(
	payload EmailDeliveryPayload,
	ctx context.Context,
	opts ...asynq.Option,
) error {
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

func (c *AsynqClient) EnqueueIncident(
	payload IncidentPayload,
	ctx context.Context,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("json.Marshal failed: %v", err)
	}
	task := asynq.NewTask(TypeIncidentProcess, jsonPayload)
	info, err := c.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("client.EnqueueContext failed: %v", err)
	}
	log.Printf("task enqueued: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

func (c *AsynqClient) EnqueueIncidentConfirmedEmail(
	ctx context.Context,
	payload IncidentConfirmedEmailPayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("EnqueueIncidentConfirmedEmail: json.Marshal failed: %w", err)
	}

	if len(opts) == 0 {
		opts = append(opts, asynq.Queue(QueueDefault), asynq.MaxRetry(5))
	}

	task := asynq.NewTask(TypeIncidentConfirmedEmail, jsonPayload)
	info, err := c.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("EnqueueIncidentConfirmedEmail: client.EnqueueContext failed: %w", err)
	}
	log.Printf("Incident confirmed email task enqueued: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

func (c *AsynqClient) EnqueueNotificationTask(
	ctx context.Context,
	payload notification.NotificationPayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("EnqueueNotificationTask: json.Marshal failed: %w", err)
	}

	// Default options if none are provided
	if len(opts) == 0 {
		opts = append(opts, asynq.Queue(QueueDefault), asynq.MaxRetry(5))
	}

	task := asynq.NewTask(TypeNotificationSend, jsonPayload)
	info, err := c.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("EnqueueNotificationTask: client.EnqueueContext failed: %w", err)
	}

	log.Printf("Notification task enqueued: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

func (c *AsynqClient) EnqueueAcceptedRegistration(
	ctx context.Context,
	payload AcceptedRegistrationFormPayload,
	opts ...asynq.Option,
) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("EnqueueAcceptedRegistration: json.Marshal failed: %w", err)
	}

	// Default options if none are provided
	if len(opts) == 0 {
		opts = append(opts, asynq.Queue(QueueDefault), asynq.MaxRetry(5))
	}

	task := asynq.NewTask(TypeAcceptedRegistration, jsonPayload)
	info, err := c.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("EnqueueAcceptedRegistration: client.EnqueueContext failed: %w", err)
	}

	log.Printf("Accepted Registration task enqueued: id=%s queue=%s", info.ID, info.Queue)
	return nil
}
