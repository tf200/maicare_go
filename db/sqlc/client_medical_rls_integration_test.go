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

func TestClientMedicalPermissionScopePolicies(t *testing.T) {
	store := openIntegrationStore(t)
	ctx := context.Background()
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	assignedActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	allActor := seedClientRLSActor(t, ctx, tx, "all", false)
	generalViewActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	diagnosisViewActor := seedClientRLSActor(t, ctx, tx, "", false)
	medicationViewActor := seedClientRLSActor(t, ctx, tx, "", false)
	diagnosisCreateActor := seedClientRLSActor(t, ctx, tx, "", false)
	medicationCreateActor := seedClientRLSActor(t, ctx, tx, "", false)
	mutationOnlyActor := seedClientRLSActor(t, ctx, tx, "", false)
	deniedActor := seedClientRLSActor(t, ctx, tx, "", false)

	for _, permission := range []string{
		"CLIENT.DIAGNOSIS.VIEW",
		"CLIENT.DIAGNOSIS.CREATE",
		"CLIENT.DIAGNOSIS.UPDATE",
		"CLIENT.DIAGNOSIS.DELETE",
		"CLIENT.MEDICATION.VIEW",
		"CLIENT.MEDICATION.CREATE",
		"CLIENT.MEDICATION.UPDATE",
		"CLIENT.MEDICATION.DELETE",
	} {
		grantClientRLSPermission(t, ctx, tx, assignedActor.userID, permission, "assigned")
		grantClientRLSPermission(t, ctx, tx, allActor.userID, permission, "all")
	}
	grantClientRLSPermission(t, ctx, tx, diagnosisViewActor.userID, "CLIENT.DIAGNOSIS.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, medicationViewActor.userID, "CLIENT.MEDICATION.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, diagnosisCreateActor.userID, "CLIENT.DIAGNOSIS.CREATE", "assigned")
	grantClientRLSPermission(t, ctx, tx, medicationCreateActor.userID, "CLIENT.MEDICATION.CREATE", "assigned")
	for _, permission := range []string{
		"CLIENT.DIAGNOSIS.UPDATE",
		"CLIENT.DIAGNOSIS.DELETE",
		"CLIENT.MEDICATION.UPDATE",
		"CLIENT.MEDICATION.DELETE",
	} {
		grantClientRLSPermission(t, ctx, tx, mutationOnlyActor.userID, permission, "assigned")
	}

	assignedClient := seedClientRLSClient(t, ctx, tx)
	unassignedClient := seedClientRLSClient(t, ctx, tx)
	for _, actor := range []clientRLSActor{
		assignedActor, generalViewActor, diagnosisViewActor, medicationViewActor,
		diagnosisCreateActor, medicationCreateActor, mutationOnlyActor,
	} {
		seedClientRLSAssignment(t, ctx, tx, assignedClient, actor.employeeID)
	}
	assignedDiagnosisID := seedMedicalDiagnosis(t, ctx, tx, assignedClient, assignedActor.employeeID)
	unassignedDiagnosisID := seedMedicalDiagnosis(t, ctx, tx, unassignedClient, assignedActor.employeeID)
	assignedMedicationID := seedMedicationOrder(t, ctx, tx, assignedClient, assignedDiagnosisID, assignedActor.employeeID)
	unassignedMedicationID := seedMedicationOrder(t, ctx, tx, unassignedClient, unassignedDiagnosisID, assignedActor.employeeID)

	runtimeRole := "phase9_group_c_runtime_" + uuid.NewString()
	if _, err := tx.Exec(ctx, fmt.Sprintf(`CREATE ROLE %s NOLOGIN NOSUPERUSER NOBYPASSRLS`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("create runtime role: %v", err)
	}
	for _, table := range []string{"public.client_diagnosis", "public.client_medication_order"} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON %s TO %s`, table, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
			t.Fatalf("grant %s privileges: %v", table, err)
		}
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT ON public.employee_profile TO %s`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("grant employee enrichment privilege: %v", err)
	}
	for _, signature := range []string{
		"public.get_current_user_id()",
		"public.get_current_employee_id()",
		"public.has_permission(text)",
		"public.get_permission_scope(text)",
		"public.is_assigned_to_client(uuid)",
		"public.can_access_client(uuid,text)",
		"public.begin_client_diagnosis_creation(uuid)",
		"public.can_read_created_client_diagnosis(uuid)",
		"public.begin_client_medication_creation(uuid)",
		"public.can_read_created_client_medication(uuid)",
	} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT EXECUTE ON FUNCTION %s TO %s`, signature, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
			t.Fatalf("grant helper %s: %v", signature, err)
		}
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`SET LOCAL ROLE %s`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("set runtime role: %v", err)
	}
	if canBypass, err := New(tx).CurrentRoleBypassesRLS(ctx); err != nil {
		t.Fatalf("check runtime RLS bypass capability: %v", err)
	} else if canBypass {
		t.Fatal("medical runtime role unexpectedly bypasses RLS")
	}
	assertMedicalTablesForced(t, ctx, tx)

	setAuthorizationActor(t, ctx, tx, assignedActor.userID.String(), assignedActor.employeeID.String())
	assertMedicalRowVisible(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, true)
	assertMedicalRowVisible(t, ctx, tx, "client_diagnosis", unassignedDiagnosisID, false)
	assertMedicalRowVisible(t, ctx, tx, "client_medication_order", assignedMedicationID, true)
	assertMedicalRowVisible(t, ctx, tx, "client_medication_order", unassignedMedicationID, false)
	assertMedicalUpdate(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, true)
	assertMedicalUpdate(t, ctx, tx, "client_diagnosis", unassignedDiagnosisID, false)
	assertMedicalUpdate(t, ctx, tx, "client_medication_order", assignedMedicationID, true)
	assertMedicalUpdate(t, ctx, tx, "client_medication_order", unassignedMedicationID, false)
	assertMedicalUpdaterAttribution(t, ctx, tx, assignedActor.employeeID, assignedClient, assignedDiagnosisID, assignedMedicationID)
	assertMedicalAttributionCannotBeForged(t, ctx, tx, assignedActor.employeeID, assignedDiagnosisID, assignedMedicationID)

	setAuthorizationActor(t, ctx, tx, generalViewActor.userID.String(), generalViewActor.employeeID.String())
	assertMedicalRowVisible(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, false)
	assertMedicalRowVisible(t, ctx, tx, "client_medication_order", assignedMedicationID, false)

	setAuthorizationActor(t, ctx, tx, diagnosisViewActor.userID.String(), diagnosisViewActor.employeeID.String())
	assertMedicalRowVisible(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, true)
	assertMedicalRowVisible(t, ctx, tx, "client_medication_order", assignedMedicationID, false)
	assertMedicalUpdate(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, false)

	setAuthorizationActor(t, ctx, tx, medicationViewActor.userID.String(), medicationViewActor.employeeID.String())
	assertMedicalRowVisible(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, false)
	assertMedicalRowVisible(t, ctx, tx, "client_medication_order", assignedMedicationID, true)
	assertMedicationDiagnosisEnrichmentHidden(t, ctx, tx, assignedClient, assignedMedicationID)
	assertMedicalUpdate(t, ctx, tx, "client_medication_order", assignedMedicationID, false)

	setAuthorizationActor(t, ctx, tx, diagnosisCreateActor.userID.String(), diagnosisCreateActor.employeeID.String())
	createdDiagnosis := createMedicalDiagnosis(t, ctx, tx, assignedClient, true)
	assertMedicalRowVisible(t, ctx, tx, "client_diagnosis", createdDiagnosis.ID, true)
	assertMedicalAttribution(t, diagnosisCreateActor.employeeID, createdDiagnosis.CreatedByEmployeeID, createdDiagnosis.UpdatedByEmployeeID)
	assertDiagnosisCreateAttributionCannotBeForged(t, ctx, tx, diagnosisCreateActor.employeeID, assignedClient)
	createMedicalDiagnosis(t, ctx, tx, unassignedClient, false)

	setAuthorizationActor(t, ctx, tx, medicationCreateActor.userID.String(), medicationCreateActor.employeeID.String())
	createdMedication := createMedicationAsRuntime(t, ctx, tx, assignedClient, &assignedDiagnosisID, true)
	assertMedicalRowVisible(t, ctx, tx, "client_medication_order", createdMedication.ID, true)
	assertMedicalAttribution(t, medicationCreateActor.employeeID, createdMedication.CreatedByEmployeeID, createdMedication.UpdatedByEmployeeID)
	assertMedicationCreateAttributionCannotBeForged(t, ctx, tx, medicationCreateActor.employeeID, assignedClient)
	createMedicationAsRuntime(t, ctx, tx, unassignedClient, &unassignedDiagnosisID, false)
	assertRuntimeCannotForgeMedicalCreation(t, ctx, tx, assignedDiagnosisID)

	setAuthorizationActor(t, ctx, tx, mutationOnlyActor.userID.String(), mutationOnlyActor.employeeID.String())
	assertMedicalRowVisible(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, false)
	assertMedicalRowVisible(t, ctx, tx, "client_medication_order", assignedMedicationID, false)
	assertMedicalUpdate(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, false)
	assertMedicalDelete(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, false)
	assertMedicalUpdate(t, ctx, tx, "client_medication_order", assignedMedicationID, false)
	assertMedicalDelete(t, ctx, tx, "client_medication_order", assignedMedicationID, false)

	setAuthorizationActor(t, ctx, tx, allActor.userID.String(), allActor.employeeID.String())
	assertMedicalRowVisible(t, ctx, tx, "client_diagnosis", unassignedDiagnosisID, true)
	assertMedicalRowVisible(t, ctx, tx, "client_medication_order", unassignedMedicationID, true)
	createMedicationAsRuntime(t, ctx, tx, assignedClient, &unassignedDiagnosisID, false)
	deletableDiagnosis := createMedicalDiagnosis(t, ctx, tx, unassignedClient, true)
	deletableMedication := createMedicationAsRuntime(t, ctx, tx, unassignedClient, &deletableDiagnosis.ID, true)
	assertLinkedDiagnosisDeleteRejected(t, ctx, tx, deletableDiagnosis.ID)
	assertMedicalDelete(t, ctx, tx, "client_medication_order", deletableMedication.ID, true)
	assertMedicalDelete(t, ctx, tx, "client_diagnosis", deletableDiagnosis.ID, true)

	setAuthorizationActor(t, ctx, tx, deniedActor.userID.String(), deniedActor.employeeID.String())
	assertMedicalRowVisible(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, false)
	assertMedicalRowVisible(t, ctx, tx, "client_medication_order", assignedMedicationID, false)

	setAuthorizationActor(t, ctx, tx, "", "")
	assertMedicalRowVisible(t, ctx, tx, "client_diagnosis", assignedDiagnosisID, false)
	assertMedicalRowVisible(t, ctx, tx, "client_medication_order", assignedMedicationID, false)
	assertMedicalRLSCannotBeDisabled(t, ctx, tx)
}

func seedMedicalDiagnosis(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, employeeID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.client_diagnosis (
		client_id, code_system, code, status, severity, created_by_employee_id, updated_by_employee_id
	) VALUES ($1, 'ICD-10', $2, 'confirmed', 'unknown', $3, $3) RETURNING id`,
		clientID, uuid.NewString(), employeeID).Scan(&id); err != nil {
		t.Fatalf("seed diagnosis: %v", err)
	}
	return id
}

func seedMedicationOrder(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, diagnosisID, employeeID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.client_medication_order (
		client_id, diagnosis_id, medication_name, dosage_text, start_date, created_by_employee_id, updated_by_employee_id
	) VALUES ($1, $2, 'Test', '10 mg', CURRENT_DATE, $3, $3) RETURNING id`,
		clientID, diagnosisID, employeeID).Scan(&id); err != nil {
		t.Fatalf("seed medication order: %v", err)
	}
	return id
}

func createMedicalDiagnosis(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, want bool) ClientDiagnosis {
	t.Helper()
	var diagnosis ClientDiagnosis
	assertRLSSavepointOperation(t, ctx, tx, "medical_diagnosis_create", want, func() error {
		var err error
		diagnosis, err = New(tx).CreateClientDiagnosis(ctx, CreateClientDiagnosisParams{
			ClientID:   clientID,
			CodeSystem: "ICD-10",
			Code:       uuid.NewString(),
			Status:     DiagnosisStatusEnumConfirmed,
			Severity:   DiagnosisSeverityEnumUnknown,
		})
		return err
	})
	return diagnosis
}

func createMedicationAsRuntime(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, diagnosisID *uuid.UUID, want bool) ClientMedicationOrder {
	t.Helper()
	var medication ClientMedicationOrder
	assertRLSSavepointOperation(t, ctx, tx, "medical_medication_create", want, func() error {
		var err error
		medication, err = New(tx).CreateClientMedicationOrder(ctx, CreateClientMedicationOrderParams{
			ClientID:       clientID,
			DiagnosisID:    diagnosisID,
			MedicationName: "Test",
			DosageText:     "10 mg",
			Schedule:       []byte("[]"),
			StartDate:      pgtype.Date{Time: time.Now(), Valid: true},
			Status:         MedicationOrderStatusEnumActive,
			AdminMode:      MedicationAdminModeEnumSelf,
		})
		return err
	})
	return medication
}

func assertMedicalRowVisible(t *testing.T, ctx context.Context, tx pgx.Tx, table string, id uuid.UUID, want bool) {
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

func assertMedicalUpdate(t *testing.T, ctx context.Context, tx pgx.Tx, table string, id uuid.UUID, want bool) {
	t.Helper()
	query := fmt.Sprintf(`UPDATE public.%s SET updated_at = CURRENT_TIMESTAMP WHERE id = $1`, pgx.Identifier{table}.Sanitize())
	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		t.Fatalf("update %s: %v", table, err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("%s updated = %t, want %t", table, got, want)
	}
}

func assertMedicalDelete(t *testing.T, ctx context.Context, tx pgx.Tx, table string, id uuid.UUID, want bool) {
	t.Helper()
	query := fmt.Sprintf(`DELETE FROM public.%s WHERE id = $1`, pgx.Identifier{table}.Sanitize())
	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		t.Fatalf("delete %s: %v", table, err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("%s deleted = %t, want %t", table, got, want)
	}
}

func assertLinkedDiagnosisDeleteRejected(t *testing.T, ctx context.Context, tx pgx.Tx, diagnosisID uuid.UUID) {
	t.Helper()
	assertRLSSavepointOperation(t, ctx, tx, "linked_diagnosis_delete", false, func() error {
		_, err := tx.Exec(ctx, `DELETE FROM public.client_diagnosis WHERE id = $1`, diagnosisID)
		return err
	})
}

func assertMedicalAttribution(t *testing.T, employeeID uuid.UUID, createdBy, updatedBy *uuid.UUID) {
	t.Helper()
	if createdBy == nil || *createdBy != employeeID {
		t.Fatalf("created_by_employee_id = %v, want %s", createdBy, employeeID)
	}
	if updatedBy == nil || *updatedBy != employeeID {
		t.Fatalf("updated_by_employee_id = %v, want %s", updatedBy, employeeID)
	}
}

func assertMedicalUpdaterAttribution(t *testing.T, ctx context.Context, tx pgx.Tx, employeeID, clientID, diagnosisID, medicationID uuid.UUID) {
	t.Helper()
	diagnosis, err := New(tx).UpdateClientDiagnosis(ctx, UpdateClientDiagnosisParams{ClientID: clientID, ID: diagnosisID})
	if err != nil {
		t.Fatalf("update diagnosis attribution: %v", err)
	}
	if diagnosis.UpdatedByEmployeeID == nil || *diagnosis.UpdatedByEmployeeID != employeeID {
		t.Fatalf("diagnosis updater = %v, want %s", diagnosis.UpdatedByEmployeeID, employeeID)
	}
	medication, err := New(tx).UpdateClientMedicationOrder(ctx, UpdateClientMedicationOrderParams{ClientID: clientID, ID: medicationID})
	if err != nil {
		t.Fatalf("update medication attribution: %v", err)
	}
	if medication.UpdatedByEmployeeID == nil || *medication.UpdatedByEmployeeID != employeeID {
		t.Fatalf("medication updater = %v, want %s", medication.UpdatedByEmployeeID, employeeID)
	}
}

func assertMedicalAttributionCannotBeForged(t *testing.T, ctx context.Context, tx pgx.Tx, employeeID, diagnosisID, medicationID uuid.UUID) {
	t.Helper()
	rows := []struct {
		table string
		id    uuid.UUID
	}{
		{table: "client_diagnosis", id: diagnosisID},
		{table: "client_medication_order", id: medicationID},
	}
	for _, row := range rows {
		var updatedBy *uuid.UUID
		query := fmt.Sprintf(`UPDATE public.%s SET updated_by_employee_id = $2 WHERE id = $1 RETURNING updated_by_employee_id`, pgx.Identifier{row.table}.Sanitize())
		if err := tx.QueryRow(ctx, query, row.id, uuid.New()).Scan(&updatedBy); err != nil {
			t.Fatalf("attempt forged %s updater: %v", row.table, err)
		}
		if updatedBy == nil || *updatedBy != employeeID {
			t.Fatalf("%s accepted forged updater %v, want %s", row.table, updatedBy, employeeID)
		}
	}
}

func assertDiagnosisCreateAttributionCannotBeForged(t *testing.T, ctx context.Context, tx pgx.Tx, employeeID, clientID uuid.UUID) {
	t.Helper()
	var createdBy, updatedBy *uuid.UUID
	if err := tx.QueryRow(ctx, `WITH new_row AS (
		SELECT public.begin_client_diagnosis_creation($1) AS id
	)
	INSERT INTO public.client_diagnosis (
		id, client_id, code_system, code, created_by_employee_id, updated_by_employee_id
	) SELECT id, $1, 'ICD-10', $2, $3, $3 FROM new_row WHERE id IS NOT NULL
	RETURNING created_by_employee_id, updated_by_employee_id`, clientID, uuid.NewString(), uuid.New()).Scan(&createdBy, &updatedBy); err != nil {
		t.Fatalf("attempt forged diagnosis creator: %v", err)
	}
	assertMedicalAttribution(t, employeeID, createdBy, updatedBy)
}

func assertMedicationCreateAttributionCannotBeForged(t *testing.T, ctx context.Context, tx pgx.Tx, employeeID, clientID uuid.UUID) {
	t.Helper()
	var createdBy, updatedBy *uuid.UUID
	if err := tx.QueryRow(ctx, `WITH new_row AS (
		SELECT public.begin_client_medication_creation($1) AS id
	)
	INSERT INTO public.client_medication_order (
		id, client_id, medication_name, dosage_text, start_date, created_by_employee_id, updated_by_employee_id
	) SELECT id, $1, 'Test', '10 mg', CURRENT_DATE, $2, $2 FROM new_row WHERE id IS NOT NULL
	RETURNING created_by_employee_id, updated_by_employee_id`, clientID, uuid.New()).Scan(&createdBy, &updatedBy); err != nil {
		t.Fatalf("attempt forged medication creator: %v", err)
	}
	assertMedicalAttribution(t, employeeID, createdBy, updatedBy)
}

func assertMedicationDiagnosisEnrichmentHidden(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, medicationID uuid.UUID) {
	t.Helper()
	row, err := New(tx).GetClientMedicationOrder(ctx, GetClientMedicationOrderParams{ClientID: clientID, ID: medicationID})
	if err != nil {
		t.Fatalf("get medication enrichment: %v", err)
	}
	if row.DiagnosisTitle != nil || row.DiagnosisCodeSystem != nil || row.DiagnosisCode != nil {
		t.Fatalf("medication query leaked diagnosis enrichment: %+v", row)
	}
}

func assertRuntimeCannotForgeMedicalCreation(t *testing.T, ctx context.Context, tx pgx.Tx, rowID uuid.UUID) {
	t.Helper()
	assertRLSSavepointOperation(t, ctx, tx, "forge_medical_creation", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.rls_report_creation_context (
			backend_pid, transaction_id, report_kind, report_id, user_id, employee_id
		) VALUES (pg_backend_pid(), pg_current_xact_id(), 'diagnosis', $1, $2, $3)`, rowID, uuid.New(), uuid.New())
		return err
	})
}

func assertMedicalTablesForced(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	for _, table := range []string{"client_diagnosis", "client_medication_order"} {
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

func assertMedicalRLSCannotBeDisabled(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	for _, table := range []string{"client_diagnosis", "client_medication_order"} {
		if _, err := tx.Exec(ctx, "SAVEPOINT group_c_rls_bypass"); err != nil {
			t.Fatalf("create %s bypass savepoint: %v", table, err)
		}
		if _, err := tx.Exec(ctx, "SET LOCAL row_security = off"); err != nil {
			t.Fatalf("disable row_security setting: %v", err)
		}
		query := fmt.Sprintf(`SELECT count(*) FROM public.%s`, pgx.Identifier{table}.Sanitize())
		if _, err := tx.Exec(ctx, query); err == nil {
			t.Fatalf("runtime role bypassed %s RLS", table)
		}
		if _, err := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT group_c_rls_bypass"); err != nil {
			t.Fatalf("rollback %s bypass attempt: %v", table, err)
		}
	}
}
