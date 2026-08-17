package middleware

import (
	"context"
	"net/netip"

	"maicare_go/internal/ctxkeys"
	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

const (
	RequestIDHeader = "X-Request-ID"
)

type contextKey string

const (
	requestIDContextKey    contextKey = "request_id"
	authPayloadKey         contextKey = "authorization_payload"
	actorRolesKey          contextKey = "actor_roles"
	effectivePermissionKey contextKey = "effective_permission"
	clientIPKey            contextKey = "client_ip"
	userAgentKey           contextKey = "user_agent"
	methodKey              contextKey = "method"
	routeKey               contextKey = "route"
	sessionIDKey           contextKey = "session_id"
)

// --- Request ID ---

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(requestIDContextKey).(string)
	return value, ok
}

// RequestIDFromContextSafe returns the request ID or empty string.
func RequestIDFromContextSafe(ctx context.Context) string {
	v, _ := ctx.Value(requestIDContextKey).(string)
	return v
}

// --- Auth Payload ---

func WithAuthPayload(ctx context.Context, payload *domain.TokenPayload) context.Context {
	return context.WithValue(ctx, authPayloadKey, payload)
}

func AuthPayloadFromContext(ctx context.Context) (*domain.TokenPayload, bool) {
	value, ok := ctx.Value(authPayloadKey).(*domain.TokenPayload)
	return value, ok
}

// --- Actor Roles ---

func WithActorRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, actorRolesKey, roles)
}

func ActorRolesFromContext(ctx context.Context) ([]string, bool) {
	value, ok := ctx.Value(actorRolesKey).([]string)
	return value, ok
}

// --- Effective Permission ---

func WithEffectivePermission(ctx context.Context, permission domain.UserPermission) context.Context {
	return context.WithValue(ctx, effectivePermissionKey, permission)
}

func EffectivePermissionFromContext(ctx context.Context) (domain.UserPermission, bool) {
	value, ok := ctx.Value(effectivePermissionKey).(domain.UserPermission)
	return value, ok
}

// --- Client IP ---

func WithClientIP(ctx context.Context, ip string) context.Context {
	parsed, err := netip.ParseAddr(ip)
	if err != nil {
		return ctx
	}
	return context.WithValue(ctx, clientIPKey, parsed)
}

func ClientIPFromContext(ctx context.Context) (netip.Addr, bool) {
	value, ok := ctx.Value(clientIPKey).(netip.Addr)
	return value, ok
}

// ClientIPFromContextSafe returns the parsed client IP or nil.
func ClientIPFromContextSafe(ctx context.Context) *netip.Addr {
	v, ok := ctx.Value(clientIPKey).(netip.Addr)
	if !ok {
		return nil
	}
	return &v
}

// --- User Agent ---

func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	return context.WithValue(ctx, userAgentKey, userAgent)
}

func UserAgentFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(userAgentKey).(string)
	return value, ok
}

// UserAgentFromContextSafe returns the user agent or empty string.
func UserAgentFromContextSafe(ctx context.Context) string {
	v, _ := ctx.Value(userAgentKey).(string)
	return v
}

// --- Method ---

func WithMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, methodKey, method)
}

func MethodFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(methodKey).(string)
	return value, ok
}

// MethodFromContextSafe returns the HTTP method or empty string.
func MethodFromContextSafe(ctx context.Context) string {
	v, _ := ctx.Value(methodKey).(string)
	return v
}

// --- Route ---

func WithRoute(ctx context.Context, route string) context.Context {
	return context.WithValue(ctx, routeKey, route)
}

func RouteFromContext(ctx context.Context) (string, bool) {
	value, ok := ctx.Value(routeKey).(string)
	return value, ok
}

// RouteFromContextSafe returns the route pattern or empty string.
func RouteFromContextSafe(ctx context.Context) string {
	v, _ := ctx.Value(routeKey).(string)
	return v
}

// --- Session ID ---

func WithSessionID(ctx context.Context, sessionID uuid.UUID) context.Context {
	return context.WithValue(ctx, sessionIDKey, sessionID)
}

func SessionIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value, ok := ctx.Value(sessionIDKey).(uuid.UUID)
	return value, ok
}

// --- Legacy helpers ---

func EmployeeIDFromContext(ctx context.Context) uuid.UUID {
	return ctxkeys.EmployeeIDFromContext(ctx)
}

func UserIDFromContext(ctx context.Context) uuid.UUID {
	payload, ok := AuthPayloadFromContext(ctx)
	if !ok || payload == nil {
		return uuid.Nil
	}
	return payload.UserID
}
