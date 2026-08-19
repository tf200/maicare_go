// Package ctxkeys provides context key definitions and accessors shared across
// internal packages. This is a low-level package with zero internal dependencies.
package ctxkeys

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

func (k ctxKey) String() string { return "ctxkeys." + string(k) }

const (
	employeeIDKey    ctxKey = "employee_id"
	actorIdentityKey ctxKey = "actor_identity"
)

type ActorIdentity struct {
	UserID     uuid.UUID `json:"user_id"`
	EmployeeID uuid.UUID `json:"employee_id"`
}

func (a ActorIdentity) IsValid() bool {
	return a.UserID != uuid.Nil && a.EmployeeID != uuid.Nil
}

func WithActorIdentity(ctx context.Context, actor ActorIdentity) context.Context {
	return context.WithValue(ctx, actorIdentityKey, actor)
}

func ActorIdentityFromContext(ctx context.Context) (ActorIdentity, bool) {
	actor, ok := ctx.Value(actorIdentityKey).(ActorIdentity)
	return actor, ok && actor.IsValid()
}

// WithEmployeeID stores the employee ID in the context.
func WithEmployeeID(ctx context.Context, employeeID uuid.UUID) context.Context {
	return context.WithValue(ctx, employeeIDKey, employeeID)
}

// EmployeeIDFromContext extracts the employee ID from the context.
// Returns uuid.Nil if not set.
func EmployeeIDFromContext(ctx context.Context) uuid.UUID {
	value, ok := ctx.Value(employeeIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return value
}
