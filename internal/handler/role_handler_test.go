package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"maicare_go/internal/domain"

	"github.com/gin-gonic/gin"
)

type roleServiceStub struct {
	domain.RoleService
	replaceErr error
}

func (s roleServiceStub) ReplaceRolePermissions(context.Context, domain.ReplaceRolePermissionsParams) error {
	return s.replaceErr
}

func TestReplaceRolePermissionsReturnsBadRequestForInvalidGrant(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewRoleHandler(roleServiceStub{replaceErr: domain.ErrInvalidRolePermissions})
	router := gin.New()
	router.POST("/roles/:role_id/permissions", handler.ReplaceRolePermissions)

	request := httptest.NewRequest(
		http.MethodPost,
		"/roles/00000000-0000-0000-0000-000000000001/permissions",
		bytes.NewBufferString(`{"permissions":[{"permission_id":"00000000-0000-0000-0000-000000000002","scope":"assigned"}]}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

func TestReplaceRolePermissionsReturnsInternalServerErrorForUnexpectedFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewRoleHandler(roleServiceStub{replaceErr: errors.New("database unavailable")})
	router := gin.New()
	router.POST("/roles/:role_id/permissions", handler.ReplaceRolePermissions)

	request := httptest.NewRequest(
		http.MethodPost,
		"/roles/00000000-0000-0000-0000-000000000001/permissions",
		bytes.NewBufferString(`{"permissions":[{"permission_id":"00000000-0000-0000-0000-000000000002","scope":null}]}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}

func TestRoleRoutesDoNotExposeDirectUserPermissionMutation(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	RegisterRoleRoutes(
		&router.RouterGroup,
		NewRoleHandler(roleServiceStub{}),
		func(ctx *gin.Context) { ctx.Next() },
		func(string) gin.HandlerFunc { return func(ctx *gin.Context) { ctx.Next() } },
	)

	request := httptest.NewRequest(
		http.MethodPost,
		"/employees/00000000-0000-0000-0000-000000000001/permissions",
		nil,
	)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
	}
}
