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

// AuthMiddleware handles authentication using Bearer tokens
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string // Variable to hold the token from either source

		// 1. Try to get the token from the Authorization header
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if authHeader != "" {
			// Header exists, process it (original logic)
			fields := strings.Fields(authHeader)
			if len(fields) < 2 { // Use < 2 for safety, although your original check was == 2
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrInvalidAuthFormat))
				return
			}

			authType := fields[0]
			if !strings.EqualFold(authType, authorizationTypeBearer) {
				// Use fmt.Errorf to include the specific unsupported type
				err := fmt.Errorf("unsupported authorization type: %s", authType)
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
				return
			}

			accessToken = fields[1] // Token found in header

		} else {
			// 2. If Authorization header is missing, try to get the token from the query parameter
			accessToken = ctx.Query(authorizationQueryKey) // Get token from query, might be "" if not present
		}

		// 3. Check if we ultimately found a token from either source
		if accessToken == "" {
			// Abort if token is missing from both header and query parameter
			// Use the new specific error message
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrMissingToken))
			return
		}

		// 4. Verify the token (this part is the same, regardless of where the token came from)
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			// Use the error from VerifyToken directly
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// 5. Store the payload in context and continue (same as before)
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

func RBACMiddleware(store *db.Store, requiredPermission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get auth payload from context (set by AuthMiddleware)
		payload, err := GetAuthPayload(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Check if role has required permission
		hasPermission, err := store.CheckRolePermission(ctx, db.CheckRolePermissionParams{
			RoleID: payload.RoleID,
			Name:   requiredPermission})
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
