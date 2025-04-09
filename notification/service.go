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
	store *db.Store
	wsHub *hub.Hub
}

// NewService creates a new notification service.
func NewService(store *db.Store, wsHub *hub.Hub) *Service {
	return &Service{
		store: store,
		wsHub: wsHub,
	}
}

type NotificationPayload struct {
	RecipientUserIDs []int64 `json:"recipient_user_ids"`
	Type             string  `json:"type"`
	Data             []byte
	CreatedAt        time.Time `json:"created_at"`
}

type WebSocketMessage struct {
	Type      string          `json:"type"`
	Data      json.RawMessage `json:"data"` // Use json.RawMessage to avoid double encoding if Data is already JSON
	CreatedAt time.Time       `json:"created_at"`
}

func (s *Service)CreateAndDeliver(ctx context.Context, payload []byte) error {
	var notificationPayload NotificationPayload
	if err := json.Unmarshal(payload, &notificationPayload); err != nil {
		log.Printf("Error unmarshalling payload: %v", err)
		// Don't retry if payload is invalid
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	log.Printf("Processing notification for %d recipients, Type: %s", len(notificationPayload.RecipientUserIDs), notificationPayload.Type)

	wsMsg := WebSocketMessage{
		Type:      notificationPayload.Type,
		Data:      notificationPayload.Data, // Assumes notificationPayload.Data is valid JSON bytes. Adjust if it's not.
		CreatedAt: notificationPayload.CreatedAt,
	}

	wsPayload, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("Error marshalling WebSocket message (Type: %s): %v", notificationPayload.Type, err)
		// If we can't marshal this, we can't send it via WS.
		// Depending on requirements, you might still want to proceed with DB saves,
		// or return an error here. Let's log and proceed with DB saves for now.
		return fmt.Errorf("failed to marshal websocket payload: %w", err) // Uncomment this if WS delivery is critical
	}
	// --- End Prepare WebSocket Message ---

	var firstError error // Keep track of the first error for potential return

	for _, recipientID := range notificationPayload.RecipientUserIDs {
		// 1. Save to Database
		_, dbErr := s.store.CreateNotification(ctx, db.CreateNotificationParams{
			UserID: recipientID,
			Type:   notificationPayload.Type,
			Data:   notificationPayload.Data,
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
