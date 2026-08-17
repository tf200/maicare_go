package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"maicare_go/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type permissionCheckerStub struct {
	grant *domain.UserPermission
	roles []string
}

func (s *permissionCheckerStub) GetEffectivePermission(context.Context, uuid.UUID, string) (*domain.UserPermission, error) {
	return s.grant, nil
}

func (s *permissionCheckerStub) GetUserRoles(context.Context, uuid.UUID) ([]string, error) {
	return s.roles, nil
}

func TestRBACMiddlewareStoresEffectivePermissionScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	assigned := domain.PermissionScopeAssigned
	checker := &permissionCheckerStub{
		grant: &domain.UserPermission{
			PermissionID:   uuid.New(),
			PermissionName: "CLIENT.VIEW",
			IsScoped:       true,
			Scope:          &assigned,
		},
		roles: []string{"coordinator"},
	}
	router := gin.New()
	router.GET("/clients", NewRBACMiddleware(checker, nil).Require("CLIENT.VIEW"), func(ctx *gin.Context) {
		grant, ok := EffectivePermissionFromContext(ctx.Request.Context())
		if !ok || grant.Scope == nil || *grant.Scope != domain.PermissionScopeAssigned {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		ctx.Status(http.StatusNoContent)
	})

	response := serveAuthorizedRequest(router, uuid.New())
	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
}

func TestRBACMiddlewareUsesCurrentGrantForEveryRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	checker := &permissionCheckerStub{roles: []string{"coordinator"}}
	router := gin.New()
	router.GET("/clients", NewRBACMiddleware(checker, nil).Require("CLIENT.VIEW"), func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})
	userID := uuid.New()

	if response := serveAuthorizedRequest(router, userID); response.Code != http.StatusForbidden {
		t.Fatalf("status without grant = %d, want %d", response.Code, http.StatusForbidden)
	}

	all := domain.PermissionScopeAll
	checker.grant = &domain.UserPermission{
		PermissionID: uuid.New(), PermissionName: "CLIENT.VIEW", IsScoped: true, Scope: &all,
	}
	if response := serveAuthorizedRequest(router, userID); response.Code != http.StatusNoContent {
		t.Fatalf("status after grant change = %d, want %d", response.Code, http.StatusNoContent)
	}
}

func TestRBACMiddlewareRejectsScopedPermissionWithoutScope(t *testing.T) {
	gin.SetMode(gin.TestMode)
	checker := &permissionCheckerStub{
		grant: &domain.UserPermission{
			PermissionID: uuid.New(), PermissionName: "CLIENT.VIEW", IsScoped: true,
		},
	}
	router := gin.New()
	router.GET("/clients", NewRBACMiddleware(checker, nil).Require("CLIENT.VIEW"), func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	response := serveAuthorizedRequest(router, uuid.New())
	if response.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusForbidden)
	}
}

func serveAuthorizedRequest(router http.Handler, userID uuid.UUID) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/clients", nil)
	request = request.WithContext(WithAuthPayload(request.Context(), &domain.TokenPayload{UserID: userID}))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}
