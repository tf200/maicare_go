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
)

var (
	ErrMissingAuthHeader = errors.New("authorization header is not provided")
	ErrInvalidAuthFormat = errors.New("invalid authorization header format")
)

// AuthMiddleware handles authentication using Bearer tokens
func AuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get the Authorization header
		authHeader := ctx.GetHeader(authorizationHeaderKey)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrMissingAuthHeader))
			return
		}

		// Split the header into parts
		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrInvalidAuthFormat))
			return
		}

		// Validate the authorization type
		authType := fields[0]
		if !strings.EqualFold(authType, authorizationTypeBearer) {
			err := fmt.Errorf("unsupported authorization type: %s", authType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Verify the token
		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Store the payload in context and continue
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

func RBACMiddleware(store *db.Store) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, err := GetAuthPayload(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		// Get the role from the payload
		role := payload.RoleID

		// Check if the role is authorized to access the resource
		if !isAuthorized(role, ctx.Request.URL.Path, ctx.Request.Method, store, ctx) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, errorResponse(errors.New("forbidden")))
			return
		}

		// Continue if the role is authorized
		ctx.Next()
	}
}

func mapHTTPMethodToPermission(method string) string {
	switch method {
	case "POST":
		return "create"
	case "GET":
		return "read"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "read" // Default to read permission for unknown methods
	}
}

func isAuthorized(roleID int32, resourcePath string, method string, store *db.Store, ctx *gin.Context) bool {
	// Convert HTTP method to permission name
	permissionName := mapHTTPMethodToPermission(method)

	// Query to check if the role has the required permission for the resource
	hasPermission, err := store.CheckRolePermission(ctx, db.CheckRolePermissionParams{
		RoleID:   roleID,
		Resource: resourcePath,
		Name:     permissionName,
	})
	if err != nil {
		return false

	}

	return hasPermission
}
