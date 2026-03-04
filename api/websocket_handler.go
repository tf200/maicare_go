package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"maicare_go/hub"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Configure the WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (server *Server) handleWebSocket(ctx *gin.Context) {
	localUpgrader := upgrader
	localUpgrader.CheckOrigin = server.checkWebSocketOrigin

	// --- 1. Authenticate via one-time WebSocket ticket ---
	ticketValue := strings.TrimSpace(ctx.Query(wsTicketQueryKey))
	authPayload, err := server.wsTicketManager.Consume(ctx, ticketValue)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("invalid websocket ticket")))
		return
	}

	userID := authPayload.UserId // Extract user ID
	log.Printf("Attempting WebSocket upgrade for authenticated user ID: %s", userID)

	// --- 2. Upgrade Connection ---
	conn, err := localUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// upgrader.Upgrade sends an HTTP error response itself if it fails
		log.Printf("Failed to upgrade WebSocket connection for user %s: %v", userID, err)
		// No need to write further HTTP response here
		return
	}
	log.Printf("WebSocket connection successfully upgraded for user ID: %s", userID)

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

func (server *Server) checkWebSocketOrigin(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return false
	}

	allowedOrigins := parseAllowedOrigins(server.config.WsAllowedOrigins)
	if len(allowedOrigins) == 0 {
		return server.config.Environment != "production"
	}

	for _, allowed := range allowedOrigins {
		if strings.EqualFold(allowed, origin) {
			return true
		}
	}

	return false
}

func parseAllowedOrigins(value string) []string {
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		origins = append(origins, trimmed)
	}
	return origins
}
