package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func TestClientNetworkPermissionScopePolicies(t *testing.T) {
	store := openIntegrationStore(t)
	ctx := context.Background()
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	assignedActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	allActor := seedClientRLSActor(t, ctx, tx, "all", false)
	creatorActor := seedClientRLSActor(t, ctx, tx, "", true)
	recipientActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	deniedActor := seedClientRLSActor(t, ctx, tx, "", false)
	for _, permission := range []string{
		"CLIENT.EMERGENCY_CONTACT.VIEW",
		"CLIENT.EMERGENCY_CONTACT.CREATE",
		"CLIENT.EMERGENCY_CONTACT.UPDATE",
		"CLIENT.EMERGENCY_CONTACT.DELETE",
		"CLIENT.INVOLVED_EMPLOYEE.VIEW",
		"CLIENT.INVOLVED_EMPLOYEE.CREATE",
		"CLIENT.INVOLVED_EMPLOYEE.UPDATE",
		"CLIENT.INVOLVED_EMPLOYEE.DELETE",
	} {
		grantClientRLSPermission(t, ctx, tx, assignedActor.userID, permission, "assigned")
		grantClientRLSPermission(t, ctx, tx, allActor.userID, permission, "all")
	}
	grantClientRLSPermission(t, ctx, tx, assignedActor.userID, "CLIENT.INCIDENT.CONFIRM", "assigned")
	grantClientRLSPermission(t, ctx, tx, allActor.userID, "CLIENT.INCIDENT.CONFIRM", "all")
	grantClientRLSPermission(t, ctx, tx, recipientActor.userID, "CLIENT.INCIDENT.CONFIRM", "assigned")

	assignedClient := seedClientRLSClient(t, ctx, tx)
	unassignedClient := seedClientRLSClient(t, ctx, tx)
	assignedEmployeeID := seedClientRLSAssignment(t, ctx, tx, assignedClient, assignedActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, recipientActor.employeeID)
	otherEmployee := seedClientRLSActor(t, ctx, tx, "", false)
	otherAssignmentID := seedClientRLSAssignment(t, ctx, tx, unassignedClient, otherEmployee.employeeID)
	assignedContactID := seedClientRLSEmergencyContact(t, ctx, tx, assignedClient)
	unassignedContactID := seedClientRLSEmergencyContact(t, ctx, tx, unassignedClient)

	runtimeRole := "phase9_group_a_runtime_" + uuid.NewString()
	if _, err := tx.Exec(ctx, fmt.Sprintf(`CREATE ROLE %s NOLOGIN NOSUPERUSER NOBYPASSRLS`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("create runtime role: %v", err)
	}
	for _, table := range []string{"public.assigned_employee", "public.client_emergency_contact"} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON %s TO %s`, table, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
			t.Fatalf("grant %s privileges: %v", table, err)
		}
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT, INSERT ON public.client_details TO %s`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("grant client creation privileges: %v", err)
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
		"public.get_authorized_client_related_emails(uuid)",
		"public.get_authorized_incident_recipient_emails(uuid)",
	} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT EXECUTE ON FUNCTION %s TO %s`, signature, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
			t.Fatalf("grant helper %s: %v", signature, err)
		}
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`SET LOCAL ROLE %s`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("set runtime role: %v", err)
	}

	setAuthorizationActor(t, ctx, tx, assignedActor.userID.String(), assignedActor.employeeID.String())
	assertClientNetworkRowVisible(t, ctx, tx, "client_emergency_contact", assignedContactID, true)
	assertClientNetworkRowVisible(t, ctx, tx, "client_emergency_contact", unassignedContactID, false)
	assertClientNetworkRowVisible(t, ctx, tx, "assigned_employee", assignedEmployeeID, true)
	assertClientNetworkRowVisible(t, ctx, tx, "assigned_employee", otherAssignmentID, false)
	assertEmergencyContactInsert(t, ctx, tx, assignedClient, true)
	assertEmergencyContactInsert(t, ctx, tx, unassignedClient, false)
	assertEmergencyContactUpdate(t, ctx, tx, assignedContactID, true)
	assertEmergencyContactUpdate(t, ctx, tx, unassignedContactID, false)
	assertEmergencyContactMove(t, ctx, tx, assignedContactID, unassignedClient, false)
	assertEmergencyContactDelete(t, ctx, tx, unassignedContactID, false)
	assertAssignmentInsert(t, ctx, tx, unassignedClient, assignedActor.employeeID, false)
	assertAssignmentUpdate(t, ctx, tx, assignedEmployeeID, false)
	assertAssignmentDelete(t, ctx, tx, assignedEmployeeID, false)
	assertAuthorizedEmailCount(t, ctx, tx, assignedClient, 1)
	assertIncidentRecipientCount(t, ctx, tx, assignedClient, 1)
	assertAuthorizedEmailCount(t, ctx, tx, unassignedClient, 0)
	assertIncidentRecipientCount(t, ctx, tx, unassignedClient, 0)

	setAuthorizationActor(t, ctx, tx, recipientActor.userID.String(), recipientActor.employeeID.String())
	assertClientNetworkRowVisible(t, ctx, tx, "client_emergency_contact", assignedContactID, false)
	assertClientNetworkRowVisible(t, ctx, tx, "assigned_employee", assignedEmployeeID, false)
	assertAuthorizedEmailCount(t, ctx, tx, assignedClient, 1)
	assertIncidentRecipientCount(t, ctx, tx, assignedClient, 1)

	setAuthorizationActor(t, ctx, tx, creatorActor.userID.String(), creatorActor.employeeID.String())
	createdClientID := insertClientAsRuntime(t, ctx, tx)
	createdContactID := assertEmergencyContactInsert(t, ctx, tx, createdClientID, true)
	assertClientNetworkRowVisible(t, ctx, tx, "client_emergency_contact", createdContactID, true)

	setAuthorizationActor(t, ctx, tx, allActor.userID.String(), allActor.employeeID.String())
	assertClientNetworkRowVisible(t, ctx, tx, "client_emergency_contact", unassignedContactID, true)
	assertClientNetworkRowVisible(t, ctx, tx, "assigned_employee", otherAssignmentID, true)
	deletableContactID := assertEmergencyContactInsert(t, ctx, tx, unassignedClient, true)
	assertEmergencyContactDelete(t, ctx, tx, deletableContactID, true)
	managedAssignmentID := assertAssignmentInsert(t, ctx, tx, unassignedClient, allActor.employeeID, true)
	assertAssignmentUpdate(t, ctx, tx, managedAssignmentID, true)
	assertAssignmentDelete(t, ctx, tx, managedAssignmentID, true)

	setAuthorizationActor(t, ctx, tx, deniedActor.userID.String(), deniedActor.employeeID.String())
	assertClientNetworkRowVisible(t, ctx, tx, "client_emergency_contact", assignedContactID, false)
	assertClientNetworkRowVisible(t, ctx, tx, "assigned_employee", assignedEmployeeID, false)

	setAuthorizationActor(t, ctx, tx, "", "")
	assertClientNetworkRowVisible(t, ctx, tx, "client_emergency_contact", assignedContactID, false)
	assertClientNetworkRowVisible(t, ctx, tx, "assigned_employee", assignedEmployeeID, false)
	assertClientNetworkRLSCannotBeDisabled(t, ctx, tx)
}

func grantClientRLSPermission(t *testing.T, ctx context.Context, tx pgx.Tx, userID uuid.UUID, permission, scope string) {
	t.Helper()
	permissionID := ensureClientRLSPermission(t, ctx, tx, permission, true)
	if _, err := tx.Exec(ctx, `INSERT INTO public.role_permissions (role_id, permission_id, scope)
		SELECT role_id, $2, $3::public.permission_scope_enum
		FROM public.user_roles
		WHERE user_id = $1`, userID, permissionID, scope); err != nil {
		t.Fatalf("grant %s: %v", permission, err)
	}
}

func seedClientRLSAssignment(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, employeeID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.assigned_employee (client_id, employee_id, start_date, role)
		VALUES ($1, $2, CURRENT_DATE, 'support') RETURNING id`, clientID, employeeID).Scan(&id); err != nil {
		t.Fatalf("seed assignment: %v", err)
	}
	return id
}

func seedClientRLSEmergencyContact(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.client_emergency_contact (
		client_id, first_name, last_name, email, is_verified, incidents_reports
	) VALUES ($1, 'Phase', 'Nine', $2, TRUE, TRUE) RETURNING id`, clientID, uuid.NewString()+"@test.invalid").Scan(&id); err != nil {
		t.Fatalf("seed emergency contact: %v", err)
	}
	return id
}

func assertClientNetworkRowVisible(t *testing.T, ctx context.Context, tx pgx.Tx, table string, id uuid.UUID, want bool) {
	t.Helper()
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM public.%s WHERE id = $1)`, pgx.Identifier{table}.Sanitize())
	var got bool
	if err := tx.QueryRow(ctx, query, id).Scan(&got); err != nil {
		t.Fatalf("read %s visibility: %v", table, err)
	}
	if got != want {
		t.Fatalf("%s %s visibility = %t, want %t", table, id, got, want)
	}
}

func assertEmergencyContactInsert(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, want bool) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	assertRLSSavepointOperation(t, ctx, tx, "emergency_contact_insert", want, func() error {
		return tx.QueryRow(ctx, `INSERT INTO public.client_emergency_contact (client_id, first_name)
			VALUES ($1, 'Test') RETURNING id`, clientID).Scan(&id)
	})
	return id
}

func assertEmergencyContactDelete(t *testing.T, ctx context.Context, tx pgx.Tx, contactID uuid.UUID, want bool) {
	t.Helper()
	tag, err := tx.Exec(ctx, `DELETE FROM public.client_emergency_contact WHERE id = $1`, contactID)
	if err != nil {
		t.Fatalf("delete emergency contact: %v", err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("emergency contact deleted = %t, want %t", got, want)
	}
}

func assertAuthorizedEmailCount(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, want int) {
	t.Helper()
	emails, err := New(tx).GetClientRelatedEmails(ctx, clientID)
	if err != nil {
		t.Fatalf("get authorized client emails: %v", err)
	}
	if len(emails) != want {
		t.Fatalf("authorized client email count = %d, want %d", len(emails), want)
	}
}

func assertIncidentRecipientCount(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, want int) {
	t.Helper()
	emails, err := New(tx).ListIncidentReportRecipientEmails(ctx, clientID)
	if err != nil {
		t.Fatalf("get authorized incident recipients: %v", err)
	}
	if len(emails) != want {
		t.Fatalf("authorized incident recipient count = %d, want %d", len(emails), want)
	}
}

func assertEmergencyContactMove(t *testing.T, ctx context.Context, tx pgx.Tx, contactID, clientID uuid.UUID, want bool) {
	t.Helper()
	assertRLSSavepointOperation(t, ctx, tx, "emergency_contact_move", want, func() error {
		_, err := tx.Exec(ctx, `UPDATE public.client_emergency_contact SET client_id = $2 WHERE id = $1`, contactID, clientID)
		return err
	})
}

func assertEmergencyContactUpdate(t *testing.T, ctx context.Context, tx pgx.Tx, contactID uuid.UUID, want bool) {
	t.Helper()
	tag, err := tx.Exec(ctx, `UPDATE public.client_emergency_contact SET first_name = 'Updated' WHERE id = $1`, contactID)
	if err != nil {
		t.Fatalf("update emergency contact: %v", err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("emergency contact updated = %t, want %t", got, want)
	}
}

func assertAssignmentInsert(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, employeeID uuid.UUID, want bool) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	assertRLSSavepointOperation(t, ctx, tx, "assignment_insert", want, func() error {
		return tx.QueryRow(ctx, `INSERT INTO public.assigned_employee (client_id, employee_id, start_date, role)
			VALUES ($1, $2, CURRENT_DATE, 'support') RETURNING id`, clientID, employeeID).Scan(&id)
	})
	return id
}

func assertAssignmentUpdate(t *testing.T, ctx context.Context, tx pgx.Tx, assignmentID uuid.UUID, want bool) {
	t.Helper()
	tag, err := tx.Exec(ctx, `UPDATE public.assigned_employee SET role = 'updated' WHERE id = $1`, assignmentID)
	if err != nil {
		t.Fatalf("update assignment: %v", err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("assignment updated = %t, want %t", got, want)
	}
}

func assertAssignmentDelete(t *testing.T, ctx context.Context, tx pgx.Tx, assignmentID uuid.UUID, want bool) {
	t.Helper()
	tag, err := tx.Exec(ctx, `DELETE FROM public.assigned_employee WHERE id = $1`, assignmentID)
	if err != nil {
		t.Fatalf("delete assignment: %v", err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("assignment deleted = %t, want %t", got, want)
	}
}

func assertRLSSavepointOperation(t *testing.T, ctx context.Context, tx pgx.Tx, name string, want bool, operation func() error) {
	t.Helper()
	if _, err := tx.Exec(ctx, "SAVEPOINT "+name); err != nil {
		t.Fatalf("create %s savepoint: %v", name, err)
	}
	err := operation()
	if got := err == nil; got != want {
		t.Fatalf("%s succeeded = %t, want %t (error: %v)", name, got, want, err)
	}
	if err != nil {
		if _, rollbackErr := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT "+name); rollbackErr != nil {
			t.Fatalf("rollback %s: %v", name, rollbackErr)
		}
	}
}

func assertClientNetworkRLSCannotBeDisabled(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	for _, table := range []string{"assigned_employee", "client_emergency_contact"} {
		if _, err := tx.Exec(ctx, "SAVEPOINT group_a_rls_bypass"); err != nil {
			t.Fatalf("create %s bypass savepoint: %v", table, err)
		}
		if _, err := tx.Exec(ctx, "SET LOCAL row_security = off"); err != nil {
			t.Fatalf("disable row_security setting: %v", err)
		}
		query := fmt.Sprintf(`SELECT count(*) FROM public.%s`, pgx.Identifier{table}.Sanitize())
		if _, err := tx.Exec(ctx, query); err == nil {
			t.Fatalf("runtime role bypassed %s RLS", table)
		}
		if _, err := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT group_a_rls_bypass"); err != nil {
			t.Fatalf("rollback %s bypass attempt: %v", table, err)
		}
	}
}
