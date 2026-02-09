package infra

import (
	"context"
	"maicare_go/token"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Authentication related constants
const (
	AuthorizationHeaderKey  = "Authorization" // Changed to proper HTTP header case
	AuthorizationTypeBearer = "Bearer"        // Changed to proper case
	AuthorizationPayloadKey = "authorization_payload"
	ActorRoleKey            = "actor_role"

	AuthorizationQueryKey = "access_token" // You can change this query param name if needed (e.g., "token")
)

func GetEmployeeID(ctx context.Context) uuid.UUID {
	if ctx == nil {
		return uuid.Nil
	}

	defer func() {
		if recover() != nil {
			// Keep this helper fail-safe for non-request contexts (seeds, jobs, CLI commands).
		}
	}()

	if ginCtx, ok := ctx.(*gin.Context); ok {
		if ginCtx == nil {
			return uuid.Nil
		}
		payload, exists := ginCtx.Get(AuthorizationPayloadKey)
		if !exists {
			return uuid.Nil
		}
		p, ok := payload.(*token.Payload)
		if !ok {
			return uuid.Nil
		}
		employeeID := p.EmployeeID

		return employeeID
	}

	payload := ctx.Value(AuthorizationPayloadKey)
	if payload == nil {
		return uuid.Nil
	}
	p, ok := payload.(*token.Payload)
	if !ok {
		return uuid.Nil
	}
	employeeID := p.EmployeeID

	return employeeID
}
