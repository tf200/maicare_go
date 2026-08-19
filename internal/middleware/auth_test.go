package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"maicare_go/internal/ctxkeys"
	"maicare_go/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type tokenVerifierStub struct {
	payload *domain.TokenPayload
}

func (s tokenVerifierStub) VerifyToken(string) (*domain.TokenPayload, error) {
	return s.payload, nil
}

func (tokenVerifierStub) CreateToken(uuid.UUID, uuid.UUID, time.Duration, domain.TokenType) (string, *domain.TokenPayload, error) {
	panic("not used")
}

func (tokenVerifierStub) CreateTokenWithSessionID(uuid.UUID, uuid.UUID, time.Duration, domain.TokenType, uuid.UUID) (string, *domain.TokenPayload, error) {
	panic("not used")
}

func TestAuthMiddlewareStoresActorIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	payload := &domain.TokenPayload{
		UserID: uuid.New(), EmployeeID: uuid.New(), TokenType: domain.AccessTokenType,
	}
	router := gin.New()
	router.GET("/", NewAuthMiddleware(tokenVerifierStub{payload: payload}, nil).Handle(), func(ctx *gin.Context) {
		actor, ok := ctxkeys.ActorIdentityFromContext(ctx.Request.Context())
		if !ok || actor.UserID != payload.UserID || actor.EmployeeID != payload.EmployeeID {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.Status(http.StatusNoContent)
	})

	response := serveBearerRequest(router)
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}

func TestAuthMiddlewareRejectsIncompleteActorIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	payload := &domain.TokenPayload{UserID: uuid.New(), TokenType: domain.AccessTokenType}
	router := gin.New()
	router.GET("/", NewAuthMiddleware(tokenVerifierStub{payload: payload}, nil).Handle(), func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	response := serveBearerRequest(router)
	if response.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func serveBearerRequest(router http.Handler) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Authorization", "Bearer test-token")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
