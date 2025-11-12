package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"maicare_go/service/audit"
	"maicare_go/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Authentication related constants
const (
	authorizationHeaderKey  = "Authorization" // Changed to proper HTTP header case
	authorizationTypeBearer = "Bearer"        // Changed to proper case
	authorizationPayloadKey = "authorization_payload"
	actorRoleKey            = "actor_role"

	authorizationQueryKey = "access_token" // You can change this query param name if needed (e.g., "token")
)

var (
	ErrMissingAuthHeader = errors.New("authorization header is not provided")
	ErrInvalidAuthFormat = errors.New("invalid authorization header format")

	ErrMissingToken = errors.New("missing access token in header and query parameter") // New error for clarity
)

var ErrUnauthorizedRole = errors.New("role is not authorized to access this resource")

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
		roles, err := s.businessService.AuthService.GetUserRoles(ctx, payload.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		ctx.Set(actorRoleKey, roles)

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
		hasPermission, err := s.businessService.AuthService.HasPermission(ctx, payload.UserId, requiredPermission)
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

// api/middleware.go
func (s *Server) AuditMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// ===== PRE-REQUEST: Gather information before processing =====

		// Get authenticated user payload
		payload, err := GetAuthPayload(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Get subject ID from URL parameter
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		// Get request ID from context
		key, exists := ctx.Get(requestIDKey)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(errors.New("request ID not found in context")))
			return
		}
		requestID, err := uuid.Parse(key.(string))
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(errors.New("invalid request ID format")))
			return
		}

		// Capture timestamp at request start
		occurredAt := time.Now()

		// Determine subject type from path (first segment after /clients)
		subjectType := "client"

		// Build action from method and path
		action := fmt.Sprintf("%s %s", ctx.Request.Method, ctx.FullPath())

		// Get client IP address
		var clientIP *netip.Addr
		if ipStr := ctx.ClientIP(); ipStr != "" {
			if addr, err := netip.ParseAddr(ipStr); err == nil {
				clientIP = &addr
			}
		}

		// ===== PROCESS REQUEST =====
		ctx.Next()

		// ===== POST-REQUEST: Complete audit record after processing =====

		// Determine event type based on HTTP method and status
		statusCode := ctx.Writer.Status()
		eventType := determineEventType(ctx.Request.Method, statusCode, subjectType)

		// Set result based on status code
		result := ""
		if statusCode >= 200 && statusCode < 300 {
			result = "success"
		} else {
			result = fmt.Sprintf("failed_%d", statusCode)
		}

		// Calculate hash for this audit record
		actorRoles, exists := ctx.Get(actorRoleKey)
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(errors.New("actor roles not found in context")))
			return
		}

		// Convert actorRoles to []string
		rolesSlice, ok := actorRoles.([]string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(errors.New("invalid actor roles type")))
			return
		}

		// Build the audit record
		auditRecord := &audit.AuditRecord{
			EventID:      requestID,
			EventType:    eventType,
			OccuredAt:    occurredAt,
			ActorRole:    rolesSlice, // TODO: set actor role
			ActorID:      payload.UserId,
			SubjectType:  subjectType,
			SubjectID:    id,
			Action:       action,
			Result:       result,
			AccessReason: "placeholder", // TODO: set access reason
			Ip:           clientIP,
		}

		// Save the audit record asynchronously to avoid blocking response
		go func(record *audit.AuditRecord) {
			if err := s.businessService.AuditService.CreateAuditRecord(ctx, record); err != nil {
				fmt.Printf("failed to create audit record: %v\n", err)
			}
		}(auditRecord)
	}
}

// determineEventType creates a descriptive event type based on the action
func determineEventType(method string, statusCode int, subjectType string) string {
	if statusCode < 200 || statusCode >= 300 {
		return fmt.Sprintf("%s.failed", subjectType)
	}

	switch method {
	case http.MethodGet:
		return fmt.Sprintf("%s.viewed", subjectType)
	case http.MethodPost:
		return fmt.Sprintf("%s.created", subjectType)
	case http.MethodPut, http.MethodPatch:
		return fmt.Sprintf("%s.updated", subjectType)
	case http.MethodDelete:
		return fmt.Sprintf("%s.deleted", subjectType)
	default:
		return fmt.Sprintf("%s.accessed", subjectType)
	}
}
