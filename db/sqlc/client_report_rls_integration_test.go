package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestClientReportPermissionScopePolicies(t *testing.T) {
	store := openIntegrationStore(t)
	ctx := context.Background()
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	assignedActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	allActor := seedClientRLSActor(t, ctx, tx, "all", false)
	generateActor := seedClientRLSActor(t, ctx, tx, "", false)
	createActor := seedClientRLSActor(t, ctx, tx, "", false)
	confirmActor := seedClientRLSActor(t, ctx, tx, "", false)
	viewActor := seedClientRLSActor(t, ctx, tx, "", false)
	aiViewActor := seedClientRLSActor(t, ctx, tx, "", false)
	updateActor := seedClientRLSActor(t, ctx, tx, "", false)
	deleteActor := seedClientRLSActor(t, ctx, tx, "", false)
	mutationOnlyActor := seedClientRLSActor(t, ctx, tx, "", false)
	deniedActor := seedClientRLSActor(t, ctx, tx, "", false)
	for _, permission := range []string{
		"CLIENT.PROGRESS_REPORT.VIEW",
		"CLIENT.PROGRESS_REPORT.CREATE",
		"CLIENT.PROGRESS_REPORT.UPDATE",
		"CLIENT.PROGRESS_REPORT.DELETE",
		"CLIENT.AI_PROGRESS_REPORT.GENERATE",
		"CLIENT.AI_PROGRESS_REPORT.CONFIRM",
		"CLIENT.AI_PROGRESS_REPORT.VIEW",
	} {
		grantClientRLSPermission(t, ctx, tx, assignedActor.userID, permission, "assigned")
		grantClientRLSPermission(t, ctx, tx, allActor.userID, permission, "all")
	}
	grantClientRLSPermission(t, ctx, tx, generateActor.userID, "CLIENT.AI_PROGRESS_REPORT.GENERATE", "assigned")
	grantClientRLSPermission(t, ctx, tx, createActor.userID, "CLIENT.PROGRESS_REPORT.CREATE", "assigned")
	grantClientRLSPermission(t, ctx, tx, confirmActor.userID, "CLIENT.AI_PROGRESS_REPORT.CONFIRM", "assigned")
	grantClientRLSPermission(t, ctx, tx, viewActor.userID, "CLIENT.PROGRESS_REPORT.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, aiViewActor.userID, "CLIENT.AI_PROGRESS_REPORT.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, updateActor.userID, "CLIENT.PROGRESS_REPORT.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, updateActor.userID, "CLIENT.PROGRESS_REPORT.UPDATE", "assigned")
	grantClientRLSPermission(t, ctx, tx, deleteActor.userID, "CLIENT.PROGRESS_REPORT.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, deleteActor.userID, "CLIENT.PROGRESS_REPORT.DELETE", "assigned")
	grantClientRLSPermission(t, ctx, tx, mutationOnlyActor.userID, "CLIENT.PROGRESS_REPORT.UPDATE", "assigned")
	grantClientRLSPermission(t, ctx, tx, mutationOnlyActor.userID, "CLIENT.PROGRESS_REPORT.DELETE", "assigned")

	assignedClient := seedClientRLSClient(t, ctx, tx)
	unassignedClient := seedClientRLSClient(t, ctx, tx)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, assignedActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, generateActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, createActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, confirmActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, viewActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, aiViewActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, updateActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, deleteActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, mutationOnlyActor.employeeID)
	assignedReportID := seedProgressReport(t, ctx, tx, assignedClient, assignedActor.employeeID)
	unassignedReportID := seedProgressReport(t, ctx, tx, unassignedClient, assignedActor.employeeID)
	deleteTargetID := seedProgressReport(t, ctx, tx, assignedClient, assignedActor.employeeID)
	assignedAIReportID := seedAIReport(t, ctx, tx, assignedClient)
	unassignedAIReportID := seedAIReport(t, ctx, tx, unassignedClient)

	runtimeRole := "phase9_group_b_runtime_" + uuid.NewString()
	if _, err := tx.Exec(ctx, fmt.Sprintf(`CREATE ROLE %s NOLOGIN NOSUPERUSER NOBYPASSRLS`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("create runtime role: %v", err)
	}
	for _, table := range []string{"public.progress_report", "public.ai_generated_reports"} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON %s TO %s`, table, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
			t.Fatalf("grant %s privileges: %v", table, err)
		}
	}
	for _, signature := range []string{
		"public.get_current_user_id()",
		"public.get_current_employee_id()",
		"public.has_permission(text)",
		"public.get_permission_scope(text)",
		"public.is_assigned_to_client(uuid)",
		"public.can_access_client(uuid,text)",
		"public.begin_progress_report_creation(uuid)",
		"public.can_read_created_progress_report(uuid)",
		"public.begin_ai_report_creation(uuid)",
		"public.can_read_created_ai_report(uuid)",
	} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT EXECUTE ON FUNCTION %s TO %s`, signature, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
			t.Fatalf("grant helper %s: %v", signature, err)
		}
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`SET LOCAL ROLE %s`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("set runtime role: %v", err)
	}
	assertReportTablesForced(t, ctx, tx)

	setAuthorizationActor(t, ctx, tx, assignedActor.userID.String(), assignedActor.employeeID.String())
	assertReportRowVisible(t, ctx, tx, "progress_report", assignedReportID, true)
	assertReportRowVisible(t, ctx, tx, "progress_report", unassignedReportID, false)
	assertReportRowVisible(t, ctx, tx, "ai_generated_reports", assignedAIReportID, true)
	assertReportRowVisible(t, ctx, tx, "ai_generated_reports", unassignedAIReportID, false)
	createdProgressID := assertProgressReportInsert(t, ctx, tx, assignedClient, assignedActor.employeeID, true)
	assertProgressReportInsert(t, ctx, tx, unassignedClient, assignedActor.employeeID, false)
	assertProgressReportUpdate(t, ctx, tx, assignedReportID, true)
	assertProgressReportUpdate(t, ctx, tx, unassignedReportID, false)
	assertProgressReportMove(t, ctx, tx, assignedReportID, unassignedClient, false)
	assertProgressReportDelete(t, ctx, tx, createdProgressID, true)
	assertProgressReportDelete(t, ctx, tx, unassignedReportID, false)
	assertAIReportInsert(t, ctx, tx, assignedClient, true)
	assertAIReportInsert(t, ctx, tx, unassignedClient, false)

	setAuthorizationActor(t, ctx, tx, generateActor.userID.String(), generateActor.employeeID.String())
	assertReportRowVisible(t, ctx, tx, "progress_report", assignedReportID, true)
	assertReportRowVisible(t, ctx, tx, "progress_report", unassignedReportID, false)
	assertReportRowVisible(t, ctx, tx, "ai_generated_reports", assignedAIReportID, false)
	assertProgressDateRangeCount(t, ctx, tx, assignedClient, 2)
	assertAIReportInsert(t, ctx, tx, assignedClient, false)

	setAuthorizationActor(t, ctx, tx, createActor.userID.String(), createActor.employeeID.String())
	createdWithoutViewID := assertProgressReportInsert(t, ctx, tx, assignedClient, createActor.employeeID, true)
	assertReportRowVisible(t, ctx, tx, "progress_report", createdWithoutViewID, true)
	assertAIReportInsert(t, ctx, tx, assignedClient, false)

	setAuthorizationActor(t, ctx, tx, confirmActor.userID.String(), confirmActor.employeeID.String())
	confirmedWithoutViewID := assertAIReportInsert(t, ctx, tx, assignedClient, true)
	assertReportRowVisible(t, ctx, tx, "ai_generated_reports", confirmedWithoutViewID, true)
	assertProgressReportInsert(t, ctx, tx, assignedClient, confirmActor.employeeID, false)
	assertRuntimeCannotForgeReportCreation(t, ctx, tx, assignedReportID)

	setAuthorizationActor(t, ctx, tx, viewActor.userID.String(), viewActor.employeeID.String())
	assertReportRowVisible(t, ctx, tx, "progress_report", assignedReportID, true)
	assertProgressReportUpdate(t, ctx, tx, assignedReportID, false)
	assertProgressReportDelete(t, ctx, tx, deleteTargetID, false)

	setAuthorizationActor(t, ctx, tx, aiViewActor.userID.String(), aiViewActor.employeeID.String())
	assertReportRowVisible(t, ctx, tx, "ai_generated_reports", assignedAIReportID, true)
	assertReportRowVisible(t, ctx, tx, "ai_generated_reports", unassignedAIReportID, false)
	assertReportRowVisible(t, ctx, tx, "progress_report", assignedReportID, false)

	setAuthorizationActor(t, ctx, tx, updateActor.userID.String(), updateActor.employeeID.String())
	assertProgressReportUpdate(t, ctx, tx, assignedReportID, true)
	assertProgressReportDelete(t, ctx, tx, deleteTargetID, false)

	setAuthorizationActor(t, ctx, tx, deleteActor.userID.String(), deleteActor.employeeID.String())
	assertProgressReportUpdate(t, ctx, tx, assignedReportID, false)
	assertProgressReportDelete(t, ctx, tx, deleteTargetID, true)

	setAuthorizationActor(t, ctx, tx, mutationOnlyActor.userID.String(), mutationOnlyActor.employeeID.String())
	assertReportRowVisible(t, ctx, tx, "progress_report", assignedReportID, false)
	assertProgressReportUpdate(t, ctx, tx, assignedReportID, false)
	assertProgressReportDelete(t, ctx, tx, assignedReportID, false)

	setAuthorizationActor(t, ctx, tx, allActor.userID.String(), allActor.employeeID.String())
	assertReportRowVisible(t, ctx, tx, "progress_report", unassignedReportID, true)
	assertReportRowVisible(t, ctx, tx, "ai_generated_reports", unassignedAIReportID, true)
	allScopeReportID := assertProgressReportInsert(t, ctx, tx, unassignedClient, allActor.employeeID, true)
	assertProgressReportUpdate(t, ctx, tx, allScopeReportID, true)
	assertProgressReportDelete(t, ctx, tx, allScopeReportID, true)
	assertAIReportMutationDenied(t, ctx, tx, unassignedAIReportID)

	setAuthorizationActor(t, ctx, tx, deniedActor.userID.String(), deniedActor.employeeID.String())
	assertReportRowVisible(t, ctx, tx, "progress_report", assignedReportID, false)
	assertReportRowVisible(t, ctx, tx, "ai_generated_reports", assignedAIReportID, false)

	setAuthorizationActor(t, ctx, tx, "", "")
	assertReportRowVisible(t, ctx, tx, "progress_report", assignedReportID, false)
	assertReportRowVisible(t, ctx, tx, "ai_generated_reports", assignedAIReportID, false)
	assertReportRLSCannotBeDisabled(t, ctx, tx)
}

func seedProgressReport(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, employeeID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.progress_report (
		client_id, employee_id, date, report_text, type, emotional_state
	) VALUES ($1, $2, CURRENT_TIMESTAMP, 'Phase Nine', 'other', 'normal') RETURNING id`, clientID, employeeID).Scan(&id); err != nil {
		t.Fatalf("seed progress report: %v", err)
	}
	return id
}

func seedAIReport(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.ai_generated_reports (
		client_id, report_text, start_date, end_date
	) VALUES ($1, 'Phase Nine AI', CURRENT_DATE, CURRENT_DATE) RETURNING id`, clientID).Scan(&id); err != nil {
		t.Fatalf("seed AI report: %v", err)
	}
	return id
}

func assertReportRowVisible(t *testing.T, ctx context.Context, tx pgx.Tx, table string, id uuid.UUID, want bool) {
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

func assertProgressReportInsert(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, employeeID uuid.UUID, want bool) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	assertRLSSavepointOperation(t, ctx, tx, "progress_report_insert", want, func() error {
		report, err := New(tx).CreateProgressReport(ctx, CreateProgressReportParams{
			ClientID:       clientID,
			EmployeeID:     &employeeID,
			Date:           pgtype.Timestamptz{Time: time.Now(), Valid: true},
			ReportText:     "Created",
			Type:           ProgressReportTypeEnumOther,
			EmotionalState: EmotionalStateEnumNormal,
		})
		id = report.ID
		return err
	})
	return id
}

func assertProgressReportUpdate(t *testing.T, ctx context.Context, tx pgx.Tx, reportID uuid.UUID, want bool) {
	t.Helper()
	tag, err := tx.Exec(ctx, `UPDATE public.progress_report SET report_text = 'Updated' WHERE id = $1`, reportID)
	if err != nil {
		t.Fatalf("update progress report: %v", err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("progress report updated = %t, want %t", got, want)
	}
}

func assertProgressReportMove(t *testing.T, ctx context.Context, tx pgx.Tx, reportID, clientID uuid.UUID, want bool) {
	t.Helper()
	assertRLSSavepointOperation(t, ctx, tx, "progress_report_move", want, func() error {
		_, err := tx.Exec(ctx, `UPDATE public.progress_report SET client_id = $2 WHERE id = $1`, reportID, clientID)
		return err
	})
}

func assertProgressReportDelete(t *testing.T, ctx context.Context, tx pgx.Tx, reportID uuid.UUID, want bool) {
	t.Helper()
	tag, err := tx.Exec(ctx, `DELETE FROM public.progress_report WHERE id = $1`, reportID)
	if err != nil {
		t.Fatalf("delete progress report: %v", err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("progress report deleted = %t, want %t", got, want)
	}
}

func assertAIReportInsert(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, want bool) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	assertRLSSavepointOperation(t, ctx, tx, "ai_report_insert", want, func() error {
		report, err := New(tx).CreateAiGeneratedReport(ctx, CreateAiGeneratedReportParams{
			ClientID:   clientID,
			ReportText: "Confirmed",
			StartDate:  pgtype.Date{Time: time.Now(), Valid: true},
			EndDate:    pgtype.Date{Time: time.Now(), Valid: true},
		})
		id = report.ID
		return err
	})
	return id
}

func assertRuntimeCannotForgeReportCreation(t *testing.T, ctx context.Context, tx pgx.Tx, reportID uuid.UUID) {
	t.Helper()
	assertRLSSavepointOperation(t, ctx, tx, "forge_report_creation", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.rls_report_creation_context (
			backend_pid, transaction_id, report_kind, report_id, user_id, employee_id
		) VALUES (pg_backend_pid(), pg_current_xact_id(), 'progress', $1, $2, $3)`,
			reportID, uuid.New(), uuid.New())
		return err
	})
}

func assertProgressDateRangeCount(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, want int) {
	t.Helper()
	var count int
	if err := tx.QueryRow(ctx, `SELECT count(*) FROM public.progress_report
		WHERE client_id = $1 AND date BETWEEN $2 AND $3`, clientID, time.Now().Add(-time.Hour), time.Now().Add(time.Hour)).Scan(&count); err != nil {
		t.Fatalf("read progress report date range: %v", err)
	}
	if count != want {
		t.Fatalf("progress report date-range count = %d, want %d", count, want)
	}
}

func assertAIReportMutationDenied(t *testing.T, ctx context.Context, tx pgx.Tx, reportID uuid.UUID) {
	t.Helper()
	operations := []struct {
		name  string
		query string
	}{
		{name: "update", query: `UPDATE public.ai_generated_reports SET report_text = 'Updated' WHERE id = $1`},
		{name: "delete", query: `DELETE FROM public.ai_generated_reports WHERE id = $1`},
	}
	for _, operation := range operations {
		tag, err := tx.Exec(ctx, operation.query, reportID)
		if err != nil {
			t.Fatalf("%s AI report: %v", operation.name, err)
		}
		if tag.RowsAffected() != 0 {
			t.Fatalf("%s AI report succeeded without a policy", operation.name)
		}
	}
}

func assertReportRLSCannotBeDisabled(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	for _, table := range []string{"progress_report", "ai_generated_reports"} {
		if _, err := tx.Exec(ctx, "SAVEPOINT group_b_rls_bypass"); err != nil {
			t.Fatalf("create %s bypass savepoint: %v", table, err)
		}
		if _, err := tx.Exec(ctx, "SET LOCAL row_security = off"); err != nil {
			t.Fatalf("disable row_security setting: %v", err)
		}
		query := fmt.Sprintf(`SELECT count(*) FROM public.%s`, pgx.Identifier{table}.Sanitize())
		if _, err := tx.Exec(ctx, query); err == nil {
			t.Fatalf("runtime role bypassed %s RLS", table)
		}
		if _, err := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT group_b_rls_bypass"); err != nil {
			t.Fatalf("rollback %s bypass attempt: %v", table, err)
		}
	}
}

func assertReportTablesForced(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	for _, table := range []string{"progress_report", "ai_generated_reports"} {
		var forced bool
		if err := tx.QueryRow(ctx, `SELECT relforcerowsecurity FROM pg_catalog.pg_class
			WHERE oid = $1::regclass`, "public."+table).Scan(&forced); err != nil {
			t.Fatalf("read %s RLS configuration: %v", table, err)
		}
		if !forced {
			t.Fatalf("%s does not force row-level security", table)
		}
	}
}
