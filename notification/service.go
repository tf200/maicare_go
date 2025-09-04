package notification

import (
	"context"

	"fmt"
	"log" // Use a structured logger in a real app
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/hub"

	"github.com/goccy/go-json"
	// "your_project_root/websocket" // Import when ready
)

// Service handles notification business logic.
type Service struct {
	store             *db.Store
	wsHub             *hub.Hub
	processorRegistry map[string]func([]byte) (any, error)
}

// NewService creates a new notification service.
func NewService(store *db.Store, wsHub *hub.Hub) *Service {
	service := &Service{
		store:             store,
		wsHub:             wsHub,
		processorRegistry: make(map[string]func([]byte) (any, error)),
	}

	// Register default processors
	service.Register(TypeNewAppointment, func(data []byte) (any, error) {
		var appData NewAppointmentData
		if err := json.Unmarshal(data, &appData); err != nil {
			return nil, fmt.Errorf("failed to unmarshal new appointment data: %w", err)
		}
		return appData, nil
	})
	return service
}

type WebSocketMessage struct {
	Type      string           `json:"type"`
	Data      NotificationData `json:"data"` // Use json.RawMessage to avoid double encoding if Data is already JSON
	CreatedAt time.Time        `json:"created_at"`
}

func (s *Service) CreateAndDeliver(ctx context.Context, payload NotificationPayload) error {

	wsMsg := WebSocketMessage{
		Type:      payload.Type,
		Data:      payload.Data,
		CreatedAt: payload.CreatedAt,
	}
	log.Printf("Preparing WebSocket message for type: %s", payload.Type)

	wsPayload, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("Error marshalling WebSocket message (Type: %s): %v", payload.Type, err)
		// If we can't marshal this, we can't send it via WS.
		// Depending on requirements, you might still want to proceed with DB saves,
		// or return an error here. Let's log and proceed with DB saves for now.
		return fmt.Errorf("failed to marshal websocket payload: %w", err) // Uncomment this if WS delivery is critical
	}
	// --- End Prepare WebSocket Message ---

	var firstError error // Keep track of the first error for potential return

	dataBytes, err := json.Marshal(payload.Data)
	if err != nil {
		log.Printf("Error marshalling notification data (Type: %s): %v", payload.Type, err)
		return fmt.Errorf("failed to marshal notification data: %w", err)
	}
	log.Printf("Notification data marshalled successfully for type: %s", payload.Type)

	for _, recipientID := range payload.RecipientUserIDs {
		log.Printf("Processing notification for recipient ID: %d", recipientID)
		// 1. Save to Database
		_, dbErr := s.store.CreateNotification(ctx, db.CreateNotificationParams{
			UserID:  recipientID,
			Type:    payload.Type,
			Data:    dataBytes,
			Message: "", // Use the original data bytes
			// You might want to store CreatedAt from the payload too,
			// ensure your DB schema/params support this if needed.
		})

		if dbErr != nil {
			log.Printf("Error saving notification to DB for user %d: %v", recipientID, dbErr)
			// Capture the first error encountered
			if firstError == nil {
				firstError = fmt.Errorf("failed to save notification for user %d: %w", recipientID, dbErr)
			}
			// Decide if you want to skip WS delivery on DB error. Usually yes.
			continue // Skip WS delivery for this user if DB save failed
		}

		log.Printf("Notification saved to DB for user %d.", recipientID)

		// 2. Deliver via WebSocket (if marshalling succeeded)
		if s.wsHub != nil { // Check if marshalling failed earlier and hub exists
			// The hub's SendToUser handles checking if the user is actually connected.
			// It iterates through all connections for that user ID.
			s.wsHub.SendToUser(recipientID, wsPayload)
			// Log the *attempt* to send. The hub logs success/failure per connection.
			log.Printf("Attempted WebSocket delivery to user %d.", recipientID)
		} else if s.wsHub == nil {
			log.Printf("WebSocket Hub is nil, skipping WS delivery for user %d.", recipientID)
		} else {
			// This means json.Marshal(wsMsg) failed earlier
			log.Printf("Skipping WebSocket delivery for user %d due to prior marshalling error.", recipientID)
		}
	}

	// Return the first error encountered during DB operations, or nil if all succeeded.
	// Asynq will handle retries based on this error return.
	return firstError
}

func (s *Service) Register(notificationType string, processor func([]byte) (interface{}, error)) {
	s.processorRegistry[notificationType] = processor
}

func (s *Service) Process(notificationType string, data []byte) (interface{}, error) {
	processor, exists := s.processorRegistry[notificationType]
	if !exists {
		return nil, fmt.Errorf("no processor registered for notification type: %s", notificationType)
	}

	return processor(data)
}
