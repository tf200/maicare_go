package domain

import (
	"strings"
	"testing"
)

func TestClientPermissionsAreScopedExceptCreate(t *testing.T) {
	for _, definition := range AllPermissionDefinitions() {
		wantScoped := strings.HasPrefix(definition.Key.String(), "CLIENT.") && definition.Key != PermClientCreate
		if definition.IsScoped != wantScoped {
			t.Errorf("permission %s: IsScoped = %t, want %t", definition.Key, definition.IsScoped, wantScoped)
		}
	}
}

func TestClientEvaluationPermissionsAreRegistered(t *testing.T) {
	want := map[PermissionKey]bool{
		PermClientEvaluationCreate: false,
		PermClientEvaluationView:   false,
	}

	for _, key := range AllPermissionKeys {
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}

	for key, found := range want {
		if !found {
			t.Errorf("permission %s is not registered", key)
		}
	}
}

func TestDefaultRolesIncludeEvaluationPermissionsWithExpectedScopes(t *testing.T) {
	want := map[string]PermissionScope{
		"admin":       PermissionScopeAll,
		"coordinator": PermissionScopeAssigned,
	}

	for _, role := range DefaultRoleSeeds() {
		wantScope, ok := want[role.Name]
		if !ok {
			continue
		}
		if role.Scope != wantScope {
			t.Errorf("role %s: scope = %s, want %s", role.Name, role.Scope, wantScope)
		}

		permissions := make(map[PermissionKey]bool, len(role.Permissions))
		for _, permission := range role.Permissions {
			permissions[permission] = true
		}
		for _, permission := range []PermissionKey{PermClientEvaluationView, PermClientEvaluationCreate} {
			if !permissions[permission] {
				t.Errorf("role %s does not include %s", role.Name, permission)
			}
		}
		delete(want, role.Name)
	}

	for role := range want {
		t.Errorf("default role %s was not found", role)
	}
}
