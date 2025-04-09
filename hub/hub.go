package hub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type UserMessage struct {
	UserID  int64
	Message []byte
}

type Hub struct {
	// Registered clients. Maps userID to a set of client pointers.
	clients map[int64]map[*Client]bool // Keep unexported

	// Inbound messages from the clients (optional).
	// broadcast chan []byte

	// Register requests from the clients.
	register chan *Client // Keep unexported

	// Unregister requests from clients.
	unregister chan *Client // Keep unexported

	// Messages to be sent to a specific user.
	sendToUser chan *UserMessage // Keep unexported

	shutdown     chan struct{} // Channel to signal shutdown
	shutdownOnce sync.Once     // Ensures shutdown logic runs only once
}

// NewHub creates and returns a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]map[*Client]bool),
		register:   make(chan *Client), // Buffered or unbuffered? Unbuffered is fine.
		unregister: make(chan *Client),
		sendToUser: make(chan *UserMessage),
		shutdown:   make(chan struct{}), // Initialize the shutdown channel
	}
}

// Run starts the hub's processing loop. It should be run in a separate goroutine.
func (h *Hub) Run() {
	// ... (Run method remains the same as before) ...
	log.Println("Hub started running")
	for {
		select {
		case client := <-h.register:
			// (handle registration as before)
			userClients, ok := h.clients[client.userID]
			if !ok {
				userClients = make(map[*Client]bool)
				h.clients[client.userID] = userClients
			}
			userClients[client] = true
			log.Printf("Client registered via channel for user %d. Total connections for user: %d", client.userID, len(userClients))

		case client := <-h.unregister:
			// (handle unregistration as before)
			userClients, ok := h.clients[client.userID]
			if ok {
				if _, clientExists := userClients[client]; clientExists {
					close(client.send)
					delete(userClients, client)
					log.Printf("Client unregistered via channel for user %d. Remaining connections for user: %d", client.userID, len(userClients))
					if len(userClients) == 0 {
						delete(h.clients, client.userID)
						log.Printf("User %d has no more connections. Removed user entry.", client.userID)
					}
				}
			}

		case userMessage := <-h.sendToUser:
			// (handle sending as before)
			userClients, ok := h.clients[userMessage.UserID]
			if ok {
				// log.Printf("Sending message to %d connection(s) for user %d", len(userClients), userMessage.UserID) // Reduce log verbosity maybe
				activeClients := 0
				for client := range userClients {
					select {
					case client.send <- userMessage.Message:
						activeClients++
					default:
						log.Printf("Client send buffer full for user %d. Forcing unregister.", client.userID)
						close(client.send)
						delete(userClients, client)
						if len(userClients) == 0 {
							delete(h.clients, client.userID)
							log.Printf("User %d has no more connections after forced unregister. Removed user entry.", client.userID)
						}
					}
				}
				if activeClients == 0 && len(userClients) > 0 {
					// This case shouldn't happen often if cleanup is working, but good to notice
					log.Printf("Warning: No active clients could receive message for user %d, but %d clients were registered.", userMessage.UserID, len(userClients))
				} else if ok {
					// log.Printf("Message sent to %d active connections for user %d", activeClients, userMessage.UserID)
				}
			} // Don't log if user not found, could be too noisy

		case <-h.shutdown:
			log.Println("Hub shutting down...")
			// Cleanup: Close all active client connections
			for userID, userClients := range h.clients {
				log.Printf("Closing %d connections for user %d", len(userClients), userID)
				for client := range userClients {
					close(client.send)            // Close the send channel first
					_ = client.conn.WriteMessage( // Attempt to send close message
						websocket.CloseMessage,
						websocket.FormatCloseMessage(websocket.CloseGoingAway, "Server shutting down"),
					)
					_ = client.conn.Close() // Force close the underlying connection
				}
				// Optionally clear the map as we go (though loop exit handles this)
				delete(h.clients, userID)
			}
			// Ensure map is fully cleared
			h.clients = make(map[int64]map[*Client]bool)
			return // Exit the Run loop
			// --- End Shutdown Case ---
		}
	}
}

// Register handles registering a client with the hub.
// It sends the client to the internal register channel.
func (h *Hub) Register(client *Client) {
	// Send the client to the Run loop's channel for safe processing
	select {
	case h.register <- client:
		log.Printf("Client for user %d queued for registration", client.userID)
	default:
		// This might happen if the register channel buffer is full (if buffered)
		// or more likely if the Run loop is stalled/not running.
		log.Printf("CRITICAL: Hub register channel blocked. Cannot register client for user %d. Closing client.", client.userID)
		// Close the connection immediately if we can't even register it.
		_ = client.conn.Close()
	}
}

// Unregister handles unregistering a client. It's called by the client's readPump usually.
// We might not need this to be public if only the client calls it internally.
// Let's keep it internal for now, triggered by client.readPump sending to h.unregister.

// SendToUser sends a message to all active connections for a specific user ID.
func (h *Hub) SendToUser(userID int64, message []byte) {
	// ... (SendToUser method remains the same as before) ...
	msg := &UserMessage{
		UserID:  userID,
		Message: message,
	}
	select {
	case h.sendToUser <- msg:
	default:
		log.Printf("Warning: Hub's sendToUser channel is blocked or Hub is not running. Message for user %d dropped.", userID)
	}
}

func (h *Hub) Shutdown() {
	h.shutdownOnce.Do(func() {
		log.Println("Signaling Hub shutdown...")
		close(h.shutdown) // Close the channel to signal Run loop
	})
}
