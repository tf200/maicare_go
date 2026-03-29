package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"maicare_go/internal/domain"
	pkgjwt "maicare_go/pkg/jwt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	ErrMissingToken      = errors.New("missing access token in authorization header")
	ErrInvalidAuthFormat = errors.New("invalid authorization header format")
)

type AuthMiddleware struct {
	tokenVerifier domain.TokenVerifier
	logger        domain.Logger
}

func NewAuthMiddleware(tokenVerifier domain.TokenVerifier, logger domain.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		tokenVerifier: tokenVerifier,
		logger:        logger,
	}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrMissingToken))
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(ErrInvalidAuthFormat))
			return
		}

		if !strings.EqualFold(fields[0], "Bearer") {
			err := fmt.Errorf("unsupported authorization type: %s", fields[0])
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		payload, err := m.tokenVerifier.VerifyToken(fields[1])
		if err != nil {
			if m.logger != nil {
				m.logger.LogWarn(ctx.Request.Context(), "AuthMiddleware", "token verification failed", zap.Error(err))
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if payload.TokenType != pkgjwt.AccessToken {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(pkgjwt.ErrInvalidToken))
			return
		}

		requestCtx := WithAuthPayload(ctx.Request.Context(), payload)
		ctx.Request = ctx.Request.WithContext(requestCtx)
		ctx.Set(string(authPayloadKey), payload)
		ctx.Set("user_id", payload.UserID.String())
		ctx.Set("employee_id", payload.EmployeeID.String())

		ctx.Next()
	}
}
