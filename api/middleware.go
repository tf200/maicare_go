package api

import (
	"errors"
	"fmt"
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
