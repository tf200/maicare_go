package api

import (
	"errors"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Authentication related constants
const (
	authorizationHeaderKey  = "Authorization" // Changed to proper HTTP header case
	authorizationTypeBearer = "Bearer"        // Changed to proper case
	authorizationPayloadKey = "authorization_payload"

	authorizationQueryKey = "access_token" // You can change this query param name if needed (e.g., "token")
)

var (
	ErrMissingAuthHeader = errors.New("authorization header is not provided")
	ErrInvalidAuthFormat = errors.New("invalid authorization header format")

	ErrMissingToken = errors.New("missing access token in header and query parameter") // New error for clarity
)
var (
	ErrUnauthorizedRole = errors.New("role is not authorized to access this resource")
)

type RoleID int32

const (
	RoleAdmin RoleID = 1
)

// Helper function to check if the request looks like a WebSocket upgrade request
func isWebSocketUpgrade(ctx *gin.Context) bool {
	// Standard headers for WebSocket upgrade requests (case-insensitive check recommended)
	// Note: The 'Connection' header value might contain multiple comma-separated values.
	return strings.ToLower(ctx.GetHeader("Upgrade")) == "websocket" &&
		strings.Contains(strings.ToLower(ctx.GetHeader("Connection")), "upgrade")
}

// AuthMiddleware restricts query param auth to WebSocket requests only.
func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string

		// 1. Try to get the token from the Authorization header (Primary Method)
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if authHeader != "" {
			fields := strings.Fields(authHeader)
			if len(fields) < 2 {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrInvalidAuthFormat))
				return
			}

			authType := fields[0]
			if !strings.EqualFold(authType, authorizationTypeBearer) {
				err := fmt.Errorf("unsupported authorization type: %s", authType)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			accessToken = fields[1] // Token found in header

		} else {
			// 2. If Authorization header is missing, check if it's a WebSocket upgrade request
			if isWebSocketUpgrade(ctx) {
				// ONLY if it's a WS upgrade request, try getting token from query parameter
				accessToken = ctx.Query(authorizationQueryKey)
				// If accessToken is still "" here, the next check will handle ErrMissingToken
			}
			// If header is missing AND it's NOT a WS upgrade request,
			// accessToken remains "" and the check below will trigger ErrMissingToken,
			// correctly enforcing header usage for non-WS requests.
		}

		// 3. Check if we ultimately found a token through an ALLOWED method
		if accessToken == "" {
			// Abort if token is missing (either header missing for non-WS,
			// or both header and query param missing for WS)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrMissingToken))
			return
		}

		// 4. Verify the token (this part is the same)
		payload, err := s.tokenMaker.VerifyToken(accessToken)
		if err != nil {
			// Handle specific token errors if needed (e.g., expired token)
			// Example: if errors.Is(err, token.ErrExpiredToken) { ... }
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err)) // Use the error from VerifyToken
			return
		}

		// 5. Store the payload in context and continue
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

// GetAuthPayload retrieves the authorization payload from the context
func GetAuthPayload(ctx *gin.Context) (*token.Payload, error) {
	payload, exists := ctx.Get(authorizationPayloadKey)
	if !exists {
		return nil, errors.New("authorization payload not found")
	}

	tokenPayload, ok := payload.(*token.Payload)
	if !ok {
		return nil, errors.New("invalid authorization payload type")
	}

	return tokenPayload, nil
}

func (s *Server) RBACMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get auth payload from context (set by AuthMiddleware)
		payload, err := GetAuthPayload(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Check if role has required permission
		hasPermission, err := s.store.CheckUserPermission(ctx, db.CheckUserPermissionParams{
			UserID: payload.UserId,
			Name:   requiredPermission,
		})
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		if !hasPermission {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(ErrUnauthorizedRole))
			return
		}

		ctx.Next()
	}
}
