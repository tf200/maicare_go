package db

import (
	"context"
	"testing"

	"maicare_go/internal/ctxkeys"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func TestAuthorizationHelpers(t *testing.T) {
	store := openIntegrationStore(t)
	ctx := context.Background()
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	actor, clientID, scopedPermission, unscopedPermission := seedAuthorizationHelperData(t, ctx, tx)

	assertAuthorizationBool(t, ctx, tx, false, `SELECT has_function_privilege('public', 'public.has_permission(text)', 'EXECUTE')`)
	assertAuthorizationBool(t, ctx, tx, true, `SELECT has_function_privilege(current_user, 'public.has_permission(text)', 'EXECUTE')`)
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.has_permission($1)`, scopedPermission)
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.can_access_client($1, $2)`, clientID, scopedPermission)
	assertAuthorizationNull(t, ctx, tx, `SELECT public.get_current_user_id()`)
	assertAuthorizationNull(t, ctx, tx, `SELECT public.get_current_employee_id()`)

	setAuthorizationActor(t, ctx, tx, "not-a-uuid", "also-not-a-uuid")
	assertAuthorizationNull(t, ctx, tx, `SELECT public.get_current_user_id()`)
	assertAuthorizationNull(t, ctx, tx, `SELECT public.get_current_employee_id()`)
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.can_access_client($1, $2)`, clientID, scopedPermission)

	setAuthorizationActor(t, ctx, tx, actor.UserID.String(), actor.EmployeeID.String())
	assertAuthorizationBool(t, ctx, tx, true, `SELECT public.has_permission($1)`, scopedPermission)
	assertAuthorizationText(t, ctx, tx, "assigned", `SELECT public.get_permission_scope($1)::TEXT`, scopedPermission)
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.has_permission('TEST.MISSING')`)
	assertAuthorizationNull(t, ctx, tx, `SELECT public.get_permission_scope('TEST.MISSING')`)
	assertAuthorizationBool(t, ctx, tx, true, `SELECT public.has_permission($1)`, unscopedPermission)
	assertAuthorizationNull(t, ctx, tx, `SELECT public.get_permission_scope($1)`, unscopedPermission)
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.can_access_client($1, $2)`, clientID, unscopedPermission)
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.can_access_client(NULL, $1)`, scopedPermission)

	if _, err := tx.Exec(ctx, `INSERT INTO public.assigned_employee (client_id, employee_id, start_date, role)
		VALUES ($1, $2, CURRENT_DATE + 1, 'support')`, clientID, actor.EmployeeID); err != nil {
		t.Fatalf("insert future assignment: %v", err)
	}
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.is_assigned_to_client($1)`, clientID)
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.can_access_client($1, $2)`, clientID, scopedPermission)

	if _, err := tx.Exec(ctx, `UPDATE public.assigned_employee SET start_date = CURRENT_DATE WHERE client_id = $1`, clientID); err != nil {
		t.Fatalf("activate assignment: %v", err)
	}
	assertAuthorizationBool(t, ctx, tx, true, `SELECT public.is_assigned_to_client($1)`, clientID)
	assertAuthorizationBool(t, ctx, tx, true, `SELECT public.can_access_client($1, $2)`, clientID, scopedPermission)

	if _, err := tx.Exec(ctx, `DELETE FROM public.assigned_employee WHERE client_id = $1`, clientID); err != nil {
		t.Fatalf("delete assignment: %v", err)
	}
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.can_access_client($1, $2)`, clientID, scopedPermission)

	if _, err := tx.Exec(ctx, `UPDATE public.role_permissions SET scope = 'all' WHERE permission_id = (
		SELECT id FROM public.permissions WHERE name = $1
	)`, scopedPermission); err != nil {
		t.Fatalf("broaden permission scope: %v", err)
	}
	assertAuthorizationText(t, ctx, tx, "all", `SELECT public.get_permission_scope($1)::TEXT`, scopedPermission)
	assertAuthorizationBool(t, ctx, tx, true, `SELECT public.can_access_client($1, $2)`, uuid.New(), scopedPermission)

	if _, err := tx.Exec(ctx, `UPDATE public.employee_profile SET is_archived = TRUE WHERE id = $1`, actor.EmployeeID); err != nil {
		t.Fatalf("archive employee: %v", err)
	}
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.can_access_client($1, $2)`, clientID, scopedPermission)
	if _, err := tx.Exec(ctx, `UPDATE public.employee_profile SET is_archived = FALSE, out_of_service = TRUE WHERE id = $1`, actor.EmployeeID); err != nil {
		t.Fatalf("mark employee out of service: %v", err)
	}
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.can_access_client($1, $2)`, clientID, scopedPermission)
	if _, err := tx.Exec(ctx, `UPDATE public.employee_profile SET out_of_service = FALSE WHERE id = $1`, actor.EmployeeID); err != nil {
		t.Fatalf("reactivate employee: %v", err)
	}

	setAuthorizationActor(t, ctx, tx, actor.UserID.String(), uuid.NewString())
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.can_access_client($1, $2)`, clientID, scopedPermission)
	setAuthorizationActor(t, ctx, tx, actor.UserID.String(), actor.EmployeeID.String())
	if _, err := tx.Exec(ctx, `UPDATE public.custom_user SET is_active = FALSE WHERE id = $1`, actor.UserID); err != nil {
		t.Fatalf("deactivate user: %v", err)
	}
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.has_permission($1)`, scopedPermission)
	assertAuthorizationNull(t, ctx, tx, `SELECT public.get_permission_scope($1)`, scopedPermission)
	assertAuthorizationBool(t, ctx, tx, false, `SELECT public.can_access_client($1, $2)`, clientID, scopedPermission)
}

func seedAuthorizationHelperData(t *testing.T, ctx context.Context, tx pgx.Tx) (ctxkeys.ActorIdentity, uuid.UUID, string, string) {
	t.Helper()
	actor := ctxkeys.ActorIdentity{UserID: uuid.New(), EmployeeID: uuid.New()}
	clientID := uuid.New()
	roleID := uuid.New()
	scopedPermissionID := uuid.New()
	unscopedPermissionID := uuid.New()
	scopedPermission := "TEST.SCOPED." + uuid.NewString()
	unscopedPermission := "TEST.UNSCOPED." + uuid.NewString()

	statements := []struct {
		query string
		args  []any
	}{
		{`INSERT INTO public.custom_user (id, password, email) VALUES ($1, 'disabled', $2)`, []any{actor.UserID, uuid.NewString() + "@test.invalid"}},
		{`INSERT INTO public.employee_profile (
			id, user_id, first_name, last_name, bsn, street, house_number, postal_code, city, gender
		) VALUES ($1, $2, 'Phase', 'Seven', $3, 'Test', '1', '1000AA', 'Test', 'unknown')`, []any{actor.EmployeeID, actor.UserID, uuid.NewString()}},
		{`INSERT INTO public.roles (id, name) VALUES ($1, $2)`, []any{roleID, "phase-seven-" + uuid.NewString()}},
		{`INSERT INTO public.permissions (id, name, is_scoped) VALUES ($1, $2, TRUE), ($3, $4, FALSE)`, []any{scopedPermissionID, scopedPermission, unscopedPermissionID, unscopedPermission}},
		{`INSERT INTO public.user_roles (user_id, role_id) VALUES ($1, $2)`, []any{actor.UserID, roleID}},
		{`INSERT INTO public.role_permissions (role_id, permission_id, scope) VALUES ($1, $2, 'assigned'), ($1, $3, NULL)`, []any{roleID, scopedPermissionID, unscopedPermissionID}},
		{`INSERT INTO public.client_details (
			id, first_name, last_name, email, gender, street, house_number, postal_code, city
		) VALUES ($1, 'Phase', 'Seven', $2, 'unknown', 'Test', '1', '1000AA', 'Test')`, []any{clientID, uuid.NewString() + "@test.invalid"}},
	}
	for _, statement := range statements {
		if _, err := tx.Exec(ctx, statement.query, statement.args...); err != nil {
			t.Fatalf("seed authorization helper data: %v", err)
		}
	}
	return actor, clientID, scopedPermission, unscopedPermission
}

func setAuthorizationActor(t *testing.T, ctx context.Context, tx pgx.Tx, userID, employeeID string) {
	t.Helper()
	if _, err := tx.Exec(ctx, `SELECT set_config('myapp.current_user_id', $1, true),
		set_config('myapp.current_employee_id', $2, true)`, userID, employeeID); err != nil {
		t.Fatalf("set actor identity: %v", err)
	}
}

func assertAuthorizationBool(t *testing.T, ctx context.Context, tx pgx.Tx, want bool, query string, args ...any) {
	t.Helper()
	var got bool
	if err := tx.QueryRow(ctx, query, args...).Scan(&got); err != nil {
		t.Fatalf("query %q: %v", query, err)
	}
	if got != want {
		t.Fatalf("query %q = %t, want %t", query, got, want)
	}
}

func assertAuthorizationText(t *testing.T, ctx context.Context, tx pgx.Tx, want, query string, args ...any) {
	t.Helper()
	var got string
	if err := tx.QueryRow(ctx, query, args...).Scan(&got); err != nil {
		t.Fatalf("query %q: %v", query, err)
	}
	if got != want {
		t.Fatalf("query %q = %q, want %q", query, got, want)
	}
}

func assertAuthorizationNull(t *testing.T, ctx context.Context, tx pgx.Tx, query string, args ...any) {
	t.Helper()
	var value any
	if err := tx.QueryRow(ctx, query, args...).Scan(&value); err != nil {
		t.Fatalf("query %q: %v", query, err)
	}
	if value != nil {
		t.Fatalf("query %q = %v, want NULL", query, value)
	}
}
