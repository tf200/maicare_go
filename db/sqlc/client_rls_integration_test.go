package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func TestClientDetailsPermissionScopePolicies(t *testing.T) {
	store := openIntegrationStore(t)
	ctx := context.Background()
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	assignedActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	allActor := seedClientRLSActor(t, ctx, tx, "all", false)
	updateOnlyActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	creatorActor := seedClientRLSActor(t, ctx, tx, "", true)
	deniedActor := seedClientRLSActor(t, ctx, tx, "", false)
	assignedClient := seedClientRLSClient(t, ctx, tx)
	unassignedClient := seedClientRLSClient(t, ctx, tx)
	if _, err := tx.Exec(ctx, `INSERT INTO public.assigned_employee (client_id, employee_id, start_date, role)
		VALUES ($1, $2, CURRENT_DATE, 'support')`, assignedClient, assignedActor.employeeID); err != nil {
		t.Fatalf("seed assigned client: %v", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO public.assigned_employee (client_id, employee_id, start_date, role)
		VALUES ($1, $2, CURRENT_DATE, 'support')`, unassignedClient, updateOnlyActor.employeeID); err != nil {
		t.Fatalf("seed update-only assignment: %v", err)
	}
	removeClientRLSPermission(t, ctx, tx, updateOnlyActor.userID, "CLIENT.VIEW")
	removeClientRLSPermission(t, ctx, tx, updateOnlyActor.userID, "CLIENT.DELETE")

	runtimeRole := "phase8_runtime_" + uuid.NewString()
	if _, err := tx.Exec(ctx, fmt.Sprintf(`CREATE ROLE %s NOLOGIN NOSUPERUSER NOBYPASSRLS`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("create runtime role: %v", err)
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON public.client_details TO %s`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("grant client privileges: %v", err)
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT USAGE ON SEQUENCE public.client_filenumber_seq TO %s`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("grant client filenumber sequence: %v", err)
	}
	for _, signature := range []string{
		"public.get_current_user_id()",
		"public.get_current_employee_id()",
		"public.has_permission(text)",
		"public.get_permission_scope(text)",
		"public.is_assigned_to_client(uuid)",
		"public.can_access_client(uuid,text)",
		"public.begin_client_creation()",
		"public.can_read_created_client(uuid)",
	} {
		query := fmt.Sprintf(`GRANT EXECUTE ON FUNCTION %s TO %s`, signature, pgx.Identifier{runtimeRole}.Sanitize())
		if _, err := tx.Exec(ctx, query); err != nil {
			t.Fatalf("grant helper %s: %v", signature, err)
		}
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`SET LOCAL ROLE %s`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("set runtime role: %v", err)
	}

	clientIDs := []uuid.UUID{assignedClient, unassignedClient}
	setAuthorizationActor(t, ctx, tx, assignedActor.userID.String(), assignedActor.employeeID.String())
	assertClientRLSVisibleCount(t, ctx, tx, clientIDs, 1)
	assertClientRLSClientVisible(t, ctx, tx, assignedClient, true)
	assertClientRLSClientVisible(t, ctx, tx, unassignedClient, false)
	assertClientRLSUpdate(t, ctx, tx, assignedClient, true)
	assertClientRLSUpdate(t, ctx, tx, unassignedClient, false)

	setAuthorizationActor(t, ctx, tx, allActor.userID.String(), allActor.employeeID.String())
	assertClientRLSVisibleCount(t, ctx, tx, clientIDs, 2)
	assertClientRLSClientVisible(t, ctx, tx, assignedClient, true)
	assertClientRLSClientVisible(t, ctx, tx, unassignedClient, true)

	setAuthorizationActor(t, ctx, tx, updateOnlyActor.userID.String(), updateOnlyActor.employeeID.String())
	assertClientRLSClientVisible(t, ctx, tx, unassignedClient, false)
	assertClientRLSUpdate(t, ctx, tx, unassignedClient, false)

	setAuthorizationActor(t, ctx, tx, deniedActor.userID.String(), deniedActor.employeeID.String())
	assertClientRLSVisibleCount(t, ctx, tx, clientIDs, 0)
	assertClientRLSInsertDenied(t, ctx, tx)

	setAuthorizationActor(t, ctx, tx, creatorActor.userID.String(), creatorActor.employeeID.String())
	createdClient := insertClientAsRuntime(t, ctx, tx)
	assertClientRLSClientVisible(t, ctx, tx, createdClient, true)
	assertRuntimeCannotForgeClientCreation(t, ctx, tx, assignedClient)

	setAuthorizationActor(t, ctx, tx, assignedActor.userID.String(), assignedActor.employeeID.String())
	assertClientRLSDelete(t, ctx, tx, unassignedClient, false)
	assertClientRLSDelete(t, ctx, tx, assignedClient, true)

	setAuthorizationActor(t, ctx, tx, "", "")
	assertClientRLSVisibleCount(t, ctx, tx, clientIDs, 0)

	assertRuntimeCannotDisableClientRLS(t, ctx, tx)
}

func removeClientRLSPermission(t *testing.T, ctx context.Context, tx pgx.Tx, userID uuid.UUID, permission string) {
	t.Helper()
	if _, err := tx.Exec(ctx, `DELETE FROM public.role_permissions AS rp
		USING public.user_roles AS ur, public.permissions AS p
		WHERE ur.user_id = $1
		  AND rp.role_id = ur.role_id
		  AND p.id = rp.permission_id
		  AND p.name = $2`, userID, permission); err != nil {
		t.Fatalf("remove %s: %v", permission, err)
	}
}

type clientRLSActor struct {
	userID     uuid.UUID
	employeeID uuid.UUID
}

func seedClientRLSActor(t *testing.T, ctx context.Context, tx pgx.Tx, scope string, canCreate bool) clientRLSActor {
	t.Helper()
	actor := clientRLSActor{userID: uuid.New(), employeeID: uuid.New()}
	roleID := uuid.New()
	if _, err := tx.Exec(ctx, `INSERT INTO public.custom_user (id, password, email) VALUES ($1, 'disabled', $2)`,
		actor.userID, uuid.NewString()+"@test.invalid"); err != nil {
		t.Fatalf("seed RLS user: %v", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO public.employee_profile (
		id, user_id, first_name, last_name, bsn, street, house_number, postal_code, city, gender
	) VALUES ($1, $2, 'Phase', 'Eight', $3, 'Test', '1', '1000AA', 'Test', 'unknown')`,
		actor.employeeID, actor.userID, uuid.NewString()); err != nil {
		t.Fatalf("seed RLS employee: %v", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO public.roles (id, name) VALUES ($1, $2)`, roleID, "phase-eight-"+uuid.NewString()); err != nil {
		t.Fatalf("seed RLS role: %v", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO public.user_roles (user_id, role_id) VALUES ($1, $2)`, actor.userID, roleID); err != nil {
		t.Fatalf("assign RLS role: %v", err)
	}
	if scope != "" {
		for _, permission := range []string{"CLIENT.VIEW", "CLIENT.UPDATE", "CLIENT.DELETE"} {
			permissionID := ensureClientRLSPermission(t, ctx, tx, permission, true)
			if _, err := tx.Exec(ctx, `INSERT INTO public.role_permissions (role_id, permission_id, scope)
				VALUES ($1, $2, $3::public.permission_scope_enum)`, roleID, permissionID, scope); err != nil {
				t.Fatalf("grant %s: %v", permission, err)
			}
		}
	}
	if canCreate {
		permissionID := ensureClientRLSPermission(t, ctx, tx, "CLIENT.CREATE", false)
		if _, err := tx.Exec(ctx, `INSERT INTO public.role_permissions (role_id, permission_id, scope) VALUES ($1, $2, NULL)`, roleID, permissionID); err != nil {
			t.Fatalf("grant CLIENT.CREATE: %v", err)
		}
	}
	return actor
}

func ensureClientRLSPermission(t *testing.T, ctx context.Context, tx pgx.Tx, name string, scoped bool) uuid.UUID {
	t.Helper()
	var permissionID uuid.UUID
	err := tx.QueryRow(ctx, `INSERT INTO public.permissions (name, is_scoped)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET is_scoped = EXCLUDED.is_scoped
		RETURNING id`, name, scoped).Scan(&permissionID)
	if err != nil {
		t.Fatalf("ensure permission %s: %v", name, err)
	}
	return permissionID
}

func seedClientRLSClient(t *testing.T, ctx context.Context, tx pgx.Tx) uuid.UUID {
	t.Helper()
	clientID := uuid.New()
	if _, err := tx.Exec(ctx, `INSERT INTO public.client_details (
		id, first_name, last_name, email, gender, filenumber, street, house_number, postal_code, city
	) VALUES ($1, 'Phase', 'Eight', $2, 'unknown', $3, 'Test', '1', '1000AA', 'Test')`,
		clientID, uuid.NewString()+"@test.invalid", uuid.NewString()); err != nil {
		t.Fatalf("seed RLS client: %v", err)
	}
	return clientID
}

func assertClientRLSVisibleCount(t *testing.T, ctx context.Context, tx pgx.Tx, clientIDs []uuid.UUID, want int) {
	t.Helper()
	var got int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM public.client_details WHERE id = ANY($1)`, clientIDs).Scan(&got); err != nil {
		t.Fatalf("count visible clients: %v", err)
	}
	if got != want {
		t.Fatalf("visible client count = %d, want %d", got, want)
	}
}

func assertClientRLSClientVisible(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, want bool) {
	t.Helper()
	var got bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM public.client_details WHERE id = $1)`, clientID).Scan(&got); err != nil {
		t.Fatalf("read client visibility: %v", err)
	}
	if got != want {
		t.Fatalf("client %s visibility = %t, want %t", clientID, got, want)
	}
}

func assertClientRLSUpdate(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, want bool) {
	t.Helper()
	tag, err := tx.Exec(ctx, `UPDATE public.client_details SET city = 'Updated' WHERE id = $1`, clientID)
	if err != nil {
		t.Fatalf("update client: %v", err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("client %s updated = %t, want %t", clientID, got, want)
	}
}

func assertClientRLSDelete(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, want bool) {
	t.Helper()
	tag, err := tx.Exec(ctx, `DELETE FROM public.client_details WHERE id = $1`, clientID)
	if err != nil {
		t.Fatalf("delete client: %v", err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("client %s deleted = %t, want %t", clientID, got, want)
	}
}

func assertClientRLSInsertDenied(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	if _, err := tx.Exec(ctx, "SAVEPOINT denied_insert"); err != nil {
		t.Fatalf("create insert savepoint: %v", err)
	}
	_, insertErr := tx.Exec(ctx, `INSERT INTO public.client_details (
		id, first_name, last_name, email, gender, filenumber, street, house_number, postal_code, city
	) VALUES ($1, 'Denied', 'Client', $2, 'unknown', $3, 'Test', '1', '1000AA', 'Test')`,
		uuid.New(), uuid.NewString()+"@test.invalid", uuid.NewString())
	if insertErr == nil {
		t.Fatal("client insert succeeded without CLIENT.CREATE")
	}
	if _, err := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT denied_insert"); err != nil {
		t.Fatalf("rollback denied insert: %v", err)
	}
}

func insertClientAsRuntime(t *testing.T, ctx context.Context, tx pgx.Tx) uuid.UUID {
	t.Helper()
	client, err := New(tx).CreateClientDetails(ctx, CreateClientDetailsParams{
		FirstName:                "Created",
		LastName:                 "Client",
		Email:                    uuid.NewString() + "@test.invalid",
		Gender:                   GenderEnumUnknown,
		Street:                   "Test",
		HouseNumber:              "1",
		PostalCode:               "1000AA",
		City:                     "Test",
		EducationLevel:           EducationLevelEnumNone,
		EvaluationIntervalsWeeks: 0,
	})
	if err != nil {
		t.Fatalf("insert client with CLIENT.CREATE: %v", err)
	}
	return client.ID
}

func assertRuntimeCannotDisableClientRLS(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	if _, err := tx.Exec(ctx, "SAVEPOINT row_security_bypass"); err != nil {
		t.Fatalf("create RLS bypass savepoint: %v", err)
	}
	if _, err := tx.Exec(ctx, "SET LOCAL row_security = off"); err != nil {
		t.Fatalf("disable row_security setting: %v", err)
	}
	if _, err := tx.Exec(ctx, `SELECT count(*) FROM public.client_details`); err == nil {
		t.Fatal("runtime role bypassed client_details RLS")
	}
	if _, err := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT row_security_bypass"); err != nil {
		t.Fatalf("rollback RLS bypass attempt: %v", err)
	}
}

func assertRuntimeCannotForgeClientCreation(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID) {
	t.Helper()
	if _, err := tx.Exec(ctx, "SAVEPOINT forge_client_creation"); err != nil {
		t.Fatalf("create context-forgery savepoint: %v", err)
	}
	_, insertErr := tx.Exec(ctx, `INSERT INTO public.rls_client_creation_context (backend_pid, transaction_id, client_id)
		VALUES (pg_backend_pid(), pg_current_xact_id(), $1)`, clientID)
	if insertErr == nil {
		t.Fatal("runtime role wrote client creation context")
	}
	if _, err := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT forge_client_creation"); err != nil {
		t.Fatalf("rollback context-forgery attempt: %v", err)
	}
}
