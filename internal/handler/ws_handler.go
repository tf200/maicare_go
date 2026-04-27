package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"maicare_go/internal/httpapi"
	"maicare_go/internal/ws"
	"maicare_go/token"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

// Configure the WebSocket upgrader
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketHandler struct {
	hub           *ws.Hub
	ticketManager *ws.TicketManager
	config        *util.Config
}

func NewWebSocketHandler(hub *ws.Hub, ticketManager *ws.TicketManager, config *util.Config) *WebSocketHandler {
	return &WebSocketHandler{
		hub:           hub,
		ticketManager: ticketManager,
		config:        config,
	}
}

// HandleWebSocket upgrades an HTTP connection to a WebSocket connection using ticket-based auth.
// @Summary WebSocket connection
// @Description Upgrade to WebSocket using a one-time ticket
// @Tags websocket
// @Produce json
// @Param ticket query string true "WebSocket ticket"
// @Success 101 {object} httpapi.Envelope[any]
// @Failure 401 {object} httpapi.Envelope[struct{}]
// @Router /ws [get]
func (h *WebSocketHandler) HandleWebSocket(ctx *gin.Context) {
	localUpgrader := wsUpgrader
	localUpgrader.CheckOrigin = h.checkWebSocketOrigin

	// --- 1. Authenticate via one-time WebSocket ticket ---
	ticketValue := strings.TrimSpace(ctx.Query(ws.WsTicketQueryKey))
	authPayload, err := h.ticketManager.Consume(ctx, ticketValue)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpapi.Fail("invalid websocket ticket", err.Error()))
		return
	}

	userID := authPayload.UserId
	log.Printf("Attempting WebSocket upgrade for authenticated user ID: %s", userID)

	// --- 2. Upgrade Connection ---
	conn, err := localUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection for user %s: %v", userID, err)
		return
	}
	log.Printf("WebSocket connection successfully upgraded for user ID: %s", userID)

	// --- 3. Create and Register Client ---
	client := ws.NewClient(h.hub, userID, conn)
	h.hub.Register(client)

	// --- 4. Start Client Goroutines ---
	client.Start()
}

// CreateWebSocketTicket creates a one-time WebSocket authentication ticket.
// @Summary Create WebSocket ticket
// @Description Generate a one-time ticket for WebSocket authentication
// @Tags websocket
// @Accept json
// @Produce json
// @Success 200 {object} httpapi.Envelope[createWebSocketTicketResponse]
// @Failure 401,500 {object} httpapi.Envelope[struct{}]
// @Router /auth/ws-ticket [post]
func (h *WebSocketHandler) CreateWebSocketTicket(ctx *gin.Context) {
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized", err.Error()))
		return
	}

	ticketValue, expiresAt, err := h.ticketManager.Issue(ctx, payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create websocket ticket", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(createWebSocketTicketResponse{
		Ticket:    ticketValue,
		ExpiresAt: expiresAt.UTC().Format(time.RFC3339),
		WSPath:    "/ws",
	}, "websocket ticket issued"))
}

func (h *WebSocketHandler) checkWebSocketOrigin(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return false
	}

	allowedOrigins := parseAllowedOrigins(h.config.WsAllowedOrigins)
	if len(allowedOrigins) == 0 {
		return h.config.Environment != "production"
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

// ==================== Types ====================

type createWebSocketTicketResponse struct {
	Ticket    string `json:"ticket"`
	ExpiresAt string `json:"expires_at"`
	WSPath    string `json:"ws_path"`
}

// ==================== Helpers ====================

func GetAuthPayload(ctx *gin.Context) (*token.Payload, error) {
	payload, exists := ctx.Get("authorization_payload")
	if !exists {
		return nil, fmt.Errorf("authorization payload not found")
	}

	authPayload, ok := payload.(*token.Payload)
	if !ok {
		return nil, fmt.Errorf("invalid authorization payload type")
	}

	if authPayload.UserId == uuid.Nil {
		return nil, fmt.Errorf("invalid user ID in authorization payload")
	}

	return authPayload, nil
}
