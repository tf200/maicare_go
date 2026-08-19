package ctxkeys

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestActorIdentityFromContext(t *testing.T) {
	actor := ActorIdentity{UserID: uuid.New(), EmployeeID: uuid.New()}
	ctx := WithActorIdentity(context.Background(), actor)

	result, ok := ActorIdentityFromContext(ctx)
	if !ok || result != actor {
		t.Fatalf("actor = %#v, ok = %t", result, ok)
	}
}

func TestActorIdentityFromContextRejectsIncompleteIdentity(t *testing.T) {
	ctx := WithActorIdentity(context.Background(), ActorIdentity{UserID: uuid.New()})

	if _, ok := ActorIdentityFromContext(ctx); ok {
		t.Fatal("expected incomplete actor identity to be rejected")
	}
}
