package async

import (
	"context"
	"fmt"
	"log"
	"maicare_go/notification"

	"github.com/goccy/go-json"

	"github.com/hibiken/asynq"
)

const (
	// Queue names (optional, but good practice if using multiple)
	QueueCritical = "critical"
	QueueDefault  = "default"
	QueueLow      = "low"

	// Task Type Names
	TypeEmailDelivery        = "email:deliver"
	TypeIncidentProcess      = "incident:process"      // Renamed for clarity
	TypeNotificationSend     = "notification:send"     // Renamed for clarity
	TypeAppointmentCreate    = "appointment:create"    // Renamed for clarity
	TypeAcceptedRegistration = "accepted:registration" // Renamed for clarity
)

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

func (c *AsynqClient) EnqueueIncident(
	payload IncidentPayload,
	ctx context.Context,
	opts ...asynq.Option) error {
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

func (c *AsynqClient) EnqueueNotificationTask(
	ctx context.Context,
	payload notification.NotificationPayload,
	opts ...asynq.Option) error {

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

func (c *AsynqClient) EnqueueAppointmentTask(
	ctx context.Context,
	payload AppointmentPayload,
	opts ...asynq.Option) error {

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("EnqueueAppointmentTask: json.Marshal failed: %w", err)
	}

	// Default options if none are provided
	if len(opts) == 0 {
		opts = append(opts, asynq.Queue(QueueDefault), asynq.MaxRetry(5))
	}

	task := asynq.NewTask(TypeAppointmentCreate, jsonPayload)
	info, err := c.client.EnqueueContext(ctx, task, opts...)
	if err != nil {
		return fmt.Errorf("EnqueueAppointmentTask: client.EnqueueContext failed: %w", err)
	}

	log.Printf("Appointment task enqueued: id=%s queue=%s", info.ID, info.Queue)
	return nil
}

func (c *AsynqClient) EnqueueAcceptedRegistration(
	ctx context.Context,
	payload AcceptedRegistrationFormPayload,
	opts ...asynq.Option) error {

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
