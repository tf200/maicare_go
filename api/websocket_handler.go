package api

import (
	"fmt"
	"log"
	"maicare_go/hub"
	"maicare_go/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Configure the WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin determines whether a request from a different origin is allowed.
	// In production, you should restrict this to your frontend's domain(s).
	// Returning true allows all origins (useful for development).
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Implement proper origin check for production
		// Example:
		// origin := r.Header.Get("Origin")
		// allowedOrigins := []string{"http://localhost:3000", "https://yourfrontend.com"}
		// for _, allowed := range allowedOrigins {
		//     if origin == allowed {
		//         return true
		//     }
		// }
		// return false
		log.Printf("Upgrading WebSocket connection from origin: %s", r.Header.Get("Origin"))
		return true // Allow all origins for now
	},
}

func (server *Server) handleWebSocket(ctx *gin.Context) {
	// --- 1. Retrieve Authenticated User Payload ---
	authPayloadValue, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		// This should ideally not happen if AuthMiddleware is working correctly
		log.Printf("Error: %s not found in context after AuthMiddleware", authorizationPayloadKey)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("missing auth payload"))) // Define ErrAuthPayloadMissing if needed
		return
	}

	authPayload, ok := authPayloadValue.(*token.Payload) // Adjust type assertion if your payload type is different
	if !ok {
		log.Printf("Error: Could not assert type of %s to *token.Payload", authorizationPayloadKey)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("wrong auth Payload"))) // Define ErrAuthPayloadType if needed
		return
	}

	userID := authPayload.UserId // Extract user ID
	log.Printf("Attempting WebSocket upgrade for authenticated user ID: %d", userID)

	// --- 2. Upgrade Connection ---
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// upgrader.Upgrade sends an HTTP error response itself if it fails
		log.Printf("Failed to upgrade WebSocket connection for user %d: %v", userID, err)
		// No need to write further HTTP response here
		return
	}
	log.Printf("WebSocket connection successfully upgraded for user ID: %d", userID)

	// --- 3. Create and Register Client ---
	// Create a new client instance associated with the hub and user ID
	client := hub.NewClient(server.hub, userID, conn)

	// Register the client with the hub's register channel
	// This is done safely within the hub's Run() loop
	server.hub.Register(client)

	// --- 4. Start Client Goroutines ---
	// Start the read and write pumps for this client in separate goroutines.
	// These methods handle the lifecycle of the connection from now on.
	client.Start()

	// Note: From this point on, we don't use the gin context (ctx) to send responses.
	// The connection is now a WebSocket managed by the client's pumps.
}

