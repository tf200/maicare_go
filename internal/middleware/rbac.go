package middleware

import (
	"context"
	"errors"
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var ErrUnauthorizedRole = errors.New("role is not authorized to access this resource")

type PermissionChecker interface {
	GetEffectivePermission(ctx context.Context, userID uuid.UUID, permission string) (*domain.UserPermission, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]string, error)
}

type RBACMiddleware struct {
	checker PermissionChecker
	logger  domain.Logger
}

func NewRBACMiddleware(checker PermissionChecker, logger domain.Logger) *RBACMiddleware {
	return &RBACMiddleware{
		checker: checker,
		logger:  logger,
	}
}

func (m *RBACMiddleware) Require(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, ok := AuthPayloadFromContext(ctx.Request.Context())
		if !ok || payload == nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, httpapi.Fail("authorization payload not found", ""))
			return
		}

		grant, err := m.checker.GetEffectivePermission(ctx.Request.Context(), payload.UserID, permission)
		if err != nil {
			if m.logger != nil {
				m.logger.LogError(ctx.Request.Context(), "RBACMiddleware", "permission check failed", err,
					zap.String("permission", permission),
					zap.String("user_id", payload.UserID.String()),
				)
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, httpapi.Fail("failed to check permissions", ""))
			return
		}

		if grant == nil || grant.PermissionName != permission || !grant.IsUsable() {
			if m.logger != nil {
				m.logger.LogWarn(ctx.Request.Context(), "RBACMiddleware", "permission denied",
					zap.String("permission", permission),
					zap.String("user_id", payload.UserID.String()),
				)
			}
			ctx.AbortWithStatusJSON(http.StatusForbidden, httpapi.Fail(ErrUnauthorizedRole.Error(), ""))
			return
		}

		roles, err := m.checker.GetUserRoles(ctx.Request.Context(), payload.UserID)
		if err != nil {
			if m.logger != nil {
				m.logger.LogError(ctx.Request.Context(), "RBACMiddleware", "failed to load user roles", err,
					zap.String("user_id", payload.UserID.String()),
				)
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, httpapi.Fail("failed to load user roles", ""))
			return
		}

		requestCtx := WithEffectivePermission(ctx.Request.Context(), *grant)
		requestCtx = WithActorRoles(requestCtx, roles)
		ctx.Request = ctx.Request.WithContext(requestCtx)
		ctx.Set(string(actorRolesKey), roles)

		ctx.Next()
	}
}
