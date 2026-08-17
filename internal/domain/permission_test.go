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
