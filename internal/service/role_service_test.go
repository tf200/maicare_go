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
	userID       uuid.UUID
	roles        []domain.UserRole
	effective    []domain.UserPermission
}

func (s *rolePermissionRepoStub) ListAllPermissions(context.Context) ([]domain.SystemPermission, error) {
	return s.permissions, s.listErr
}

func (s *rolePermissionRepoStub) ReplaceRolePermissions(_ context.Context, _ uuid.UUID, permissions []domain.PermissionGrant) error {
	s.replaceCalls++
	s.replaced = permissions
	return s.replaceErr
}

func (s *rolePermissionRepoStub) GetUserIDByEmployeeID(context.Context, uuid.UUID) (uuid.UUID, error) {
	return s.userID, nil
}

func (s *rolePermissionRepoStub) GetUserRoles(context.Context, uuid.UUID) ([]domain.UserRole, error) {
	return s.roles, nil
}

func (s *rolePermissionRepoStub) ListEffectiveUserPermissions(context.Context, uuid.UUID) ([]domain.UserPermission, error) {
	return s.effective, nil
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

func TestListUserRolesAndPermissionsReturnsRolePermissions(t *testing.T) {
	userID := uuid.New()
	roleID := uuid.New()
	permissionID := uuid.New()
	assigned := domain.PermissionScopeAssigned
	repo := &rolePermissionRepoStub{
		userID: userID,
		roles:  []domain.UserRole{{ID: roleID, Name: "coordinator"}},
		effective: []domain.UserPermission{{
			PermissionID: permissionID, PermissionName: "CLIENT.VIEW", IsScoped: true, Scope: &assigned,
		}},
	}
	service := NewRoleService(repo, noopLogger{})

	result, err := service.ListUserRolesAndPermissions(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("ListUserRolesAndPermissions() error = %v", err)
	}
	if result.Role == nil || result.Role.ID != roleID {
		t.Fatalf("role = %#v, want role %s", result.Role, roleID)
	}
	if len(result.EffectivePermissions) != 1 || result.EffectivePermissions[0].ID != permissionID {
		t.Fatalf("effective permissions = %#v", result.EffectivePermissions)
	}
	permission := result.EffectivePermissions[0]
	if !permission.IsScoped || permission.Scope == nil || *permission.Scope != domain.PermissionScopeAssigned {
		t.Fatalf("effective permission scope = %#v", permission)
	}
}
