package hub

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 // Adjust as needed
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The authenticated user ID associated with this connection.
	userID int64

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// NewClient creates a new Client instance.
// This should be called by the HTTP handler after successful upgrade and authentication.
func NewClient(hub *Hub, userID int64, conn *websocket.Conn) *Client {
	return &Client{
		hub:    hub,
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 256), // Buffered channel
	}
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		log.Printf("Client %d disconnected (readPump exit)", c.userID)
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait)) // Use underscore to ignore error; SetReadDeadline failing often means conn is already closing
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// Read message from WebSocket
		// Note: We are not processing incoming messages in this basic example,
		// but simply keeping the connection alive and handling closure/errors.
		// If you need to receive messages *from* the client, you'd handle `messageType` and `message` here.
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message for user %d: %v", c.userID, err)
			} else {
				log.Printf("websocket closed for user %d: %v", c.userID, err) // Log expected closures too
			}
			break // Exit loop on error or closure
		}

		// Optional: Process message if needed (e.g., echo, handle commands)
		// For now, we can just log it or ignore it if client->server messages aren't needed
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Printf("Received message from user %d: %s", c.userID, message) // Example logging

		// Reset read deadline with every message read (optional, pong handler does this too)
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close() // Ensure connection is closed on exit
		log.Printf("Client %d disconnected (writePump exit)", c.userID)
		// No need to explicitly unregister here, readPump or Hub's cleanup handles it
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait)) // Set deadline for this write
			if !ok {
				// The hub closed the channel.
				log.Printf("Hub closed channel for user %d. Closing connection.", c.userID)
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error getting next writer for user %d: %v", c.userID, err)
				return // Exit if we can't get a writer
			}
			_, err = w.Write(message)
			if err != nil {
				log.Printf("Error writing message for user %d: %v", c.userID, err)
				// Don't return immediately, try closing the writer
			}

			// Add queued chat messages to the current websocket message.
			// This attempts to batch messages for efficiency.
			n := len(c.send)
			for i := 0; i < n; i++ {
				// Add a newline separator if you want distinct messages in one frame
				// _, _ = w.Write(newline)
				_, err = w.Write(<-c.send)
				if err != nil {
					log.Printf("Error writing batched message for user %d: %v", c.userID, err)
					// Don't return immediately
				}
			}

			if err := w.Close(); err != nil {
				log.Printf("Error closing writer for user %d: %v", c.userID, err)
				return // Exit if closing the writer fails
			}

		case <-ticker.C:
			// Send Ping message periodically
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error writing ping for user %d: %v", c.userID, err)
				return // Assume connection is dead if ping fails
			}
		}
	}
}

// Start starts the client's read and write pumps in separate goroutines.
// This should be called by the HTTP handler after the client is registered.
func (c *Client) Start() {
	go c.writePump()
	go c.readPump()
	log.Printf("Client pumps started for user %d", c.userID)
}
