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
	if ctx, ok := ctx.(*gin.Context); ok {
		payload, exists := ctx.Get(AuthorizationPayloadKey)
		if !exists {
			return uuid.Nil
		}
		p, ok := payload.(*token.Payload)
		if !ok {
			return uuid.Nil
		}
		employeeID := p.EmployeeID

		return employeeID
	} else {
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
}
