package service

import (
	"context"
	"errors"
	"testing"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type rolePermissionRepoStub struct {
	domain.RoleRepository
	permissions  []domain.SystemPermission
	listErr      error
	replaceErr   error
	replaceCalls int
	replaced     []domain.PermissionGrant
}

func (s *rolePermissionRepoStub) ListAllPermissions(context.Context) ([]domain.SystemPermission, error) {
	return s.permissions, s.listErr
}

func (s *rolePermissionRepoStub) ReplaceRolePermissions(_ context.Context, _ uuid.UUID, permissions []domain.PermissionGrant) error {
	s.replaceCalls++
	s.replaced = permissions
	return s.replaceErr
}

type noopLogger struct{}

func (noopLogger) LogError(context.Context, string, string, error, ...zap.Field) {}
func (noopLogger) LogWarn(context.Context, string, string, ...zap.Field)         {}
func (noopLogger) LogInfo(context.Context, string, string, ...zap.Field)         {}

func TestReplaceRolePermissionsRejectsInvalidGrants(t *testing.T) {
	scopedID := uuid.New()
	unscopedID := uuid.New()
	assigned := domain.PermissionScopeAssigned
	invalid := domain.PermissionScope("organization")

	tests := []struct {
		name   string
		grants []domain.PermissionGrant
	}{
		{name: "unknown permission", grants: []domain.PermissionGrant{{PermissionID: uuid.New()}}},
		{name: "duplicate permission", grants: []domain.PermissionGrant{{PermissionID: unscopedID}, {PermissionID: unscopedID}}},
		{name: "missing scope", grants: []domain.PermissionGrant{{PermissionID: scopedID}}},
		{name: "invalid scope", grants: []domain.PermissionGrant{{PermissionID: scopedID, Scope: &invalid}}},
		{name: "scope on unscoped permission", grants: []domain.PermissionGrant{{PermissionID: unscopedID, Scope: &assigned}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &rolePermissionRepoStub{permissions: []domain.SystemPermission{
				{ID: scopedID, Name: "CLIENT.VIEW", IsScoped: true},
				{ID: unscopedID, Name: "CLIENT.CREATE"},
			}}
			service := NewRoleService(repo, noopLogger{})

			err := service.ReplaceRolePermissions(context.Background(), domain.ReplaceRolePermissionsParams{
				RoleID:      uuid.New(),
				Permissions: tt.grants,
			})
			if !errors.Is(err, domain.ErrInvalidRolePermissions) {
				t.Fatalf("expected ErrInvalidRolePermissions, got %v", err)
			}
			if repo.replaceCalls != 0 {
				t.Fatalf("expected no replacement, got %d calls", repo.replaceCalls)
			}
		})
	}
}

func TestReplaceRolePermissionsAcceptsMixedGrants(t *testing.T) {
	scopedID := uuid.New()
	unscopedID := uuid.New()
	all := domain.PermissionScopeAll
	grants := []domain.PermissionGrant{
		{PermissionID: scopedID, Scope: &all},
		{PermissionID: unscopedID},
	}
	repo := &rolePermissionRepoStub{permissions: []domain.SystemPermission{
		{ID: scopedID, Name: "CLIENT.VIEW", IsScoped: true},
		{ID: unscopedID, Name: "CLIENT.CREATE"},
	}}
	service := NewRoleService(repo, noopLogger{})

	if err := service.ReplaceRolePermissions(context.Background(), domain.ReplaceRolePermissionsParams{
		RoleID: uuid.New(), Permissions: grants,
	}); err != nil {
		t.Fatalf("ReplaceRolePermissions() error = %v", err)
	}
	if repo.replaceCalls != 1 {
		t.Fatalf("expected one atomic replacement call, got %d", repo.replaceCalls)
	}
	if len(repo.replaced) != len(grants) {
		t.Fatalf("replaced %d grants, want %d", len(repo.replaced), len(grants))
	}
}

func TestReplaceRolePermissionsReturnsRepositoryFailure(t *testing.T) {
	permissionID := uuid.New()
	replaceErr := errors.New("insert failed")
	repo := &rolePermissionRepoStub{
		permissions: []domain.SystemPermission{{ID: permissionID, Name: "CLIENT.CREATE"}},
		replaceErr:  replaceErr,
	}
	service := NewRoleService(repo, noopLogger{})

	err := service.ReplaceRolePermissions(context.Background(), domain.ReplaceRolePermissionsParams{
		RoleID: uuid.New(), Permissions: []domain.PermissionGrant{{PermissionID: permissionID}},
	})
	if !errors.Is(err, replaceErr) {
		t.Fatalf("expected repository error, got %v", err)
	}
	if repo.replaceCalls != 1 {
		t.Fatalf("expected one atomic replacement call, got %d", repo.replaceCalls)
	}
}
