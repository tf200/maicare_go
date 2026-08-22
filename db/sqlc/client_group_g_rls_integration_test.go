package db

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestClientGroupGPermissionPoliciesAndPublicFlows(t *testing.T) {
	store := openIntegrationStore(t)
	ctx := context.Background()
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	fullActor := seedClientRLSActor(t, ctx, tx, "assigned", true)
	registrationActor := seedClientRLSActor(t, ctx, tx, "", false)
	intakeActor := seedClientRLSActor(t, ctx, tx, "", false)
	for _, permission := range []string{
		"REGISTRATION_FORM.VIEW", "REGISTRATION_FORM.UPDATE", "REGISTRATION_FORM.DELETE",
		"INTAKE_FORM.VIEW", "INTAKE_FORM.CREATE", "INTAKE_FORM.UPDATE", "INTAKE_FORM.DELETE",
	} {
		grantGroupGPermission(t, ctx, tx, fullActor.userID, permission)
	}
	grantGroupGPermission(t, ctx, tx, registrationActor.userID, "REGISTRATION_FORM.VIEW")
	grantGroupGPermission(t, ctx, tx, intakeActor.userID, "INTAKE_FORM.VIEW")

	preRegistration := seedGroupGRegistration(t, ctx, tx, nil, nil)
	preIntake := seedGroupGIntake(t, ctx, tx, preRegistration)
	preAssessment := seedGroupGAssessment(t, ctx, tx, preIntake)

	assignedRegistration := seedGroupGRegistration(t, ctx, tx, nil, nil)
	assignedIntake := seedGroupGIntake(t, ctx, tx, assignedRegistration)
	assignedAssessment := seedGroupGAssessment(t, ctx, tx, assignedIntake)
	assignedClient := seedGroupGPromotedClient(t, ctx, tx, assignedRegistration, assignedIntake)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, fullActor.employeeID)

	unassignedRegistration := seedGroupGRegistration(t, ctx, tx, nil, nil)
	unassignedIntake := seedGroupGIntake(t, ctx, tx, unassignedRegistration)
	unassignedAssessment := seedGroupGAssessment(t, ctx, tx, unassignedIntake)
	seedGroupGPromotedClient(t, ctx, tx, unassignedRegistration, unassignedIntake)

	validDate := time.Now().UTC().Add(48 * time.Hour).Truncate(time.Second)
	intakeToken := uuid.NewString()
	seedGroupGRegistration(t, ctx, tx, &intakeToken, &validDate)
	expiredToken := uuid.NewString()
	expiredRegistration := seedGroupGRegistration(t, ctx, tx, &expiredToken, &validDate)
	if _, err := tx.Exec(ctx, `UPDATE public.registration_form
		SET intake_token_expires_at = NOW() - INTERVAL '1 minute' WHERE id = $1`, expiredRegistration); err != nil {
		t.Fatalf("expire public intake token: %v", err)
	}

	uploadToken := uuid.NewString()
	uploadHash := sha256.Sum256([]byte(uploadToken))
	var uploadSession uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.registration_upload_sessions (token_hash, expires_at)
		VALUES ($1, NOW() + INTERVAL '1 hour') RETURNING id`, hex.EncodeToString(uploadHash[:])).Scan(&uploadSession); err != nil {
		t.Fatalf("seed public upload session: %v", err)
	}
	publicAttachment := uuid.New()
	unrelatedAttachment := uuid.New()
	for _, attachmentID := range []uuid.UUID{publicAttachment, unrelatedAttachment} {
		if _, err := tx.Exec(ctx, `INSERT INTO public.attachment_file (uuid, name, file, size, is_used)
			VALUES ($1, 'registration.pdf', $2, 1, TRUE)`, attachmentID, attachmentID.String()+".pdf"); err != nil {
			t.Fatalf("seed public attachment: %v", err)
		}
	}
	if _, err := tx.Exec(ctx, `UPDATE public.registration_upload_sessions
		SET attachment_ids = ARRAY[$2]::UUID[] WHERE id = $1`, uploadSession, publicAttachment); err != nil {
		t.Fatalf("bind public attachment: %v", err)
	}
	emptyUploadToken := uuid.NewString()
	emptyUploadHash := sha256.Sum256([]byte(emptyUploadToken))
	if _, err := tx.Exec(ctx, `INSERT INTO public.registration_upload_sessions (token_hash, expires_at)
		VALUES ($1, NOW() + INTERVAL '1 hour')`, hex.EncodeToString(emptyUploadHash[:])); err != nil {
		t.Fatalf("seed attachment-free upload session: %v", err)
	}

	runtimeRole := "phase9_group_g_runtime_" + uuid.NewString()
	roleName := pgx.Identifier{runtimeRole}.Sanitize()
	if _, err := tx.Exec(ctx, fmt.Sprintf(`CREATE ROLE %s NOLOGIN NOSUPERUSER NOBYPASSRLS`, roleName)); err != nil {
		t.Fatalf("create runtime role: %v", err)
	}
	for _, table := range []string{
		"public.registration_form", "public.intake_forms", "public.intake_topic_assessments",
		"public.registration_upload_sessions", "public.attachment_file", "public.client_details",
		"public.employee_profile", "public.topics",
	} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON %s TO %s`, table, roleName)); err != nil {
			t.Fatalf("grant privileges on %s: %v", table, err)
		}
	}
	for _, signature := range []string{
		"public.get_current_user_id()", "public.get_current_employee_id()", "public.has_permission(text)",
		"public.get_permission_scope(text)", "public.is_assigned_to_client(uuid)", "public.can_access_client(uuid,text)",
		"public.can_access_registration_form(uuid,text)", "public.can_access_intake_form(uuid,text)",
		"public.can_create_intake_form(uuid)", "public.can_access_intake_assessment(uuid,text)",
		"public.can_read_intake_promotion_source(uuid)", "public.begin_client_creation()",
		"public.can_read_created_client(uuid)",
		"public.begin_public_registration_submission(text,uuid[])",
		"public.authorize_public_registration_insert(uuid,uuid[])",
		"public.can_read_public_submitted_registration(uuid)",
		"public.consume_public_registration_submission(uuid)", "public.get_public_intake_options(text)",
		"public.select_public_intake_date(text,timestamptz)",
	} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT EXECUTE ON FUNCTION %s TO %s`, signature, roleName)); err != nil {
			t.Fatalf("grant helper %s: %v", signature, err)
		}
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`SET LOCAL ROLE %s`, roleName)); err != nil {
		t.Fatalf("set runtime role: %v", err)
	}
	if canBypass, err := New(tx).CurrentRoleBypassesRLS(ctx); err != nil {
		t.Fatalf("check runtime RLS bypass capability: %v", err)
	} else if canBypass {
		t.Fatal("Group G runtime role unexpectedly bypasses RLS")
	}
	assertGroupGTablesForced(t, ctx, tx)

	setAuthorizationActor(t, ctx, tx, fullActor.userID.String(), fullActor.employeeID.String())
	assertGroupGRowVisible(t, ctx, tx, "registration_form", preRegistration, true)
	assertGroupGRowVisible(t, ctx, tx, "intake_forms", preIntake, true)
	assertGroupGRowVisible(t, ctx, tx, "intake_topic_assessments", preAssessment, true)
	assertGroupGRowVisible(t, ctx, tx, "registration_form", assignedRegistration, true)
	assertGroupGRowVisible(t, ctx, tx, "intake_forms", assignedIntake, true)
	assertGroupGRowVisible(t, ctx, tx, "intake_topic_assessments", assignedAssessment, true)
	assertGroupGRowVisible(t, ctx, tx, "registration_form", unassignedRegistration, false)
	assertGroupGRowVisible(t, ctx, tx, "intake_forms", unassignedIntake, false)
	assertGroupGRowVisible(t, ctx, tx, "intake_topic_assessments", unassignedAssessment, false)
	assertRLSSavepointOperation(t, ctx, tx, "group_g_promoted_intake_immutable", false, func() error {
		_, err := tx.Exec(ctx, `UPDATE public.intake_forms SET family_situation = 'changed' WHERE id = $1`, assignedIntake)
		return err
	})
	assertRLSSavepointOperation(t, ctx, tx, "group_g_unowned_registration_attachment", false, func() error {
		_, err := tx.Exec(ctx, `UPDATE public.registration_form SET document_referral = $2 WHERE id = $1`, preRegistration, unrelatedAttachment)
		return err
	})

	assertRLSSavepointOperation(t, ctx, tx, "group_g_intake_provenance", false, func() error {
		_, err := tx.Exec(ctx, `UPDATE public.intake_topic_assessments
			SET intake_form_id = $2 WHERE id = $1`, preAssessment, assignedIntake)
		return err
	})

	setAuthorizationActor(t, ctx, tx, registrationActor.userID.String(), registrationActor.employeeID.String())
	assertGroupGRowVisible(t, ctx, tx, "registration_form", preRegistration, true)
	assertGroupGRowVisible(t, ctx, tx, "intake_forms", preIntake, false)

	setAuthorizationActor(t, ctx, tx, intakeActor.userID.String(), intakeActor.employeeID.String())
	assertGroupGRowVisible(t, ctx, tx, "registration_form", preRegistration, false)
	assertGroupGRowVisible(t, ctx, tx, "intake_forms", preIntake, true)
	assertGroupGRowVisible(t, ctx, tx, "intake_topic_assessments", preAssessment, true)

	setAuthorizationActor(t, ctx, tx, fullActor.userID.String(), fullActor.employeeID.String())
	var createdClientID uuid.UUID
	err = tx.QueryRow(ctx, `SELECT public.begin_client_creation()`).Scan(&createdClientID)
	if err != nil || createdClientID == uuid.Nil {
		t.Fatalf("begin intake promotion client creation: id=%s err=%v", createdClientID, err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO public.client_details (
		id, registration_form_id, intake_form_id, first_name, last_name, email,
		gender, filenumber, street, house_number, postal_code, city
	) VALUES ($1, $2, $3, 'Promoted', 'Client', $4, 'unknown', $5, 'Test', '1', '1000AA', 'Test')`,
		createdClientID, preRegistration, preIntake, uuid.NewString()+"@test.invalid", uuid.NewString()); err != nil {
		t.Fatalf("create promoted client: %v", err)
	}
	assertGroupGRowVisible(t, ctx, tx, "intake_topic_assessments", preAssessment, true)

	setAuthorizationActor(t, ctx, tx, "", "")
	emptySessionID, err := New(tx).BeginPublicRegistrationSubmission(ctx, BeginPublicRegistrationSubmissionParams{
		TokenHash: hex.EncodeToString(emptyUploadHash[:]), AttachmentIds: []uuid.UUID{},
	})
	if err != nil || emptySessionID == uuid.Nil {
		t.Fatalf("begin attachment-free public submission: id=%s err=%v", emptySessionID, err)
	}
	var emptyRegistration uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.registration_form (
		client_first_name, client_last_name, client_bsn_number, client_gender, client_nationality,
		client_phone_number, client_email, client_street, client_house_number, client_postal_code,
		client_city, referrer_first_name, referrer_last_name, referrer_organization,
		referrer_phone_number, referrer_email
	) VALUES (
		'No', 'Documents', 'empty-bsn', 'unknown', 'Dutch', '0612345678',
		'empty@test.invalid', 'Public street', '1', '1000AA', 'Amsterdam',
		'Public', 'Referrer', 'Referrer Org', '0612345678', 'empty-referrer@test.invalid'
	) RETURNING id`).Scan(&emptyRegistration); err != nil {
		t.Fatalf("submit attachment-free registration: %v", err)
	}
	if consumed, err := New(tx).ConsumePublicRegistrationSubmission(ctx, emptySessionID); err != nil || !consumed {
		t.Fatalf("consume attachment-free submission: consumed=%t err=%v", consumed, err)
	}
	assertGroupGRowVisible(t, ctx, tx, "registration_form", emptyRegistration, true)

	if duplicateSessionID, err := New(tx).BeginPublicRegistrationSubmission(ctx, BeginPublicRegistrationSubmissionParams{
		TokenHash: hex.EncodeToString(uploadHash[:]), AttachmentIds: []uuid.UUID{publicAttachment, publicAttachment},
	}); err == nil && duplicateSessionID != uuid.Nil {
		t.Fatal("duplicate public attachment unexpectedly authorized")
	}
	if sessionID, err := New(tx).BeginPublicRegistrationSubmission(ctx, BeginPublicRegistrationSubmissionParams{
		TokenHash: hex.EncodeToString(uploadHash[:]), AttachmentIds: []uuid.UUID{unrelatedAttachment},
	}); err == nil && sessionID != uuid.Nil {
		t.Fatal("unrelated attachment unexpectedly authorized for public submission")
	}
	sessionID, err := New(tx).BeginPublicRegistrationSubmission(ctx, BeginPublicRegistrationSubmissionParams{
		TokenHash: hex.EncodeToString(uploadHash[:]), AttachmentIds: []uuid.UUID{publicAttachment},
	})
	if err != nil {
		t.Fatalf("begin public registration submission: %v", err)
	}
	var submittedRegistration uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.registration_form (
		client_first_name, client_last_name, client_bsn_number, client_gender, client_nationality,
		client_phone_number, client_email, client_street, client_house_number, client_postal_code,
		client_city, referrer_first_name, referrer_last_name, referrer_organization,
		referrer_phone_number, referrer_email, document_referral
	) VALUES (
		'Public', 'Applicant', 'public-bsn', 'unknown', 'Dutch',
		'0612345678', 'public@test.invalid', 'Public street', '1', '1000AA',
		'Amsterdam', 'Public', 'Referrer', 'Referrer Org', '0612345678',
		'referrer@test.invalid', $1
	) RETURNING id`, publicAttachment).Scan(&submittedRegistration); err != nil {
		t.Fatalf("submit public registration: %v", err)
	}
	consumed, err := New(tx).ConsumePublicRegistrationSubmission(ctx, sessionID)
	if err != nil || !consumed {
		t.Fatalf("consume public registration submission: consumed=%t err=%v", consumed, err)
	}
	if reusedSessionID, err := New(tx).BeginPublicRegistrationSubmission(ctx, BeginPublicRegistrationSubmissionParams{
		TokenHash: hex.EncodeToString(uploadHash[:]), AttachmentIds: []uuid.UUID{publicAttachment},
	}); err == nil && reusedSessionID != uuid.Nil {
		t.Fatal("consumed upload session unexpectedly reusable")
	}
	assertGroupGRowVisible(t, ctx, tx, "registration_form", submittedRegistration, true)

	optionsJSON, err := New(tx).GetPublicIntakeOptions(ctx, intakeToken)
	if err != nil || len(optionsJSON) == 0 || string(optionsJSON) == "null" {
		t.Fatalf("get public intake options: options=%s err=%v", optionsJSON, err)
	}
	if expiredOptions, err := New(tx).GetPublicIntakeOptions(ctx, expiredToken); err != nil || (len(expiredOptions) > 0 && string(expiredOptions) != "null") {
		t.Fatalf("expired public intake token accepted: options=%s err=%v", expiredOptions, err)
	}
	selected, err := New(tx).SelectPublicIntakeDate(ctx, SelectPublicIntakeDateParams{
		Token: intakeToken, SelectedDate: pgtype.Timestamptz{Time: validDate.AddDate(0, 0, 1), Valid: true},
	})
	if err != nil || selected {
		t.Fatalf("arbitrary public intake date accepted: selected=%t err=%v", selected, err)
	}
	selected, err = New(tx).SelectPublicIntakeDate(ctx, SelectPublicIntakeDateParams{
		Token: intakeToken, SelectedDate: pgtype.Timestamptz{Time: validDate, Valid: true},
	})
	if err != nil || !selected {
		t.Fatalf("valid public intake date rejected: selected=%t err=%v", selected, err)
	}
	if selectedAgain, err := New(tx).SelectPublicIntakeDate(ctx, SelectPublicIntakeDateParams{
		Token: intakeToken, SelectedDate: pgtype.Timestamptz{Time: validDate, Valid: true},
	}); err != nil || selectedAgain {
		t.Fatalf("consumed intake token reusable: selected=%t err=%v", selectedAgain, err)
	}
}

func grantGroupGPermission(t *testing.T, ctx context.Context, tx pgx.Tx, userID uuid.UUID, permission string) {
	t.Helper()
	permissionID := ensureClientRLSPermission(t, ctx, tx, permission, false)
	if _, err := tx.Exec(ctx, `INSERT INTO public.role_permissions (role_id, permission_id, scope)
		SELECT role_id, $2, NULL FROM public.user_roles WHERE user_id = $1`, userID, permissionID); err != nil {
		t.Fatalf("grant %s: %v", permission, err)
	}
}

func seedGroupGRegistration(t *testing.T, ctx context.Context, tx pgx.Tx, intakeToken *string, optionDate *time.Time) uuid.UUID {
	t.Helper()
	options := []byte(`[]`)
	if optionDate != nil {
		options = []byte(fmt.Sprintf(`[%q]`, optionDate.Format(time.DateOnly)))
	}
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.registration_form (
		client_first_name, client_last_name, client_bsn_number, client_gender, client_nationality,
		client_phone_number, client_email, client_street, client_house_number, client_postal_code,
		client_city, referrer_first_name, referrer_last_name, referrer_organization,
		referrer_phone_number, referrer_email, form_status, intake_token, intake_token_expires_at, intake_options
	) VALUES (
		'Group', 'G', $1, 'unknown', 'Dutch', '0612345678', $2, 'Test', '1', '1000AA',
		'Test', 'Referrer', 'Group G', 'Test Org', '0612345678', $3, 'processed', $4,
		CASE WHEN $4::VARCHAR IS NULL THEN NULL ELSE NOW() + INTERVAL '7 days' END, $5
	) RETURNING id`, uuid.NewString(), uuid.NewString()+"@test.invalid", uuid.NewString()+"@test.invalid", intakeToken, options).Scan(&id); err != nil {
		t.Fatalf("seed Group G registration: %v", err)
	}
	return id
}

func seedGroupGIntake(t *testing.T, ctx context.Context, tx pgx.Tx, registrationID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.intake_forms (
		registration_form_id, care_type, self_sufficiency, intake_conclusion
	) VALUES ($1, 'protected_living', 3, 'suitable') RETURNING id`, registrationID).Scan(&id); err != nil {
		t.Fatalf("seed Group G intake: %v", err)
	}
	return id
}

func seedGroupGAssessment(t *testing.T, ctx context.Context, tx pgx.Tx, intakeID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.intake_topic_assessments (
		intake_form_id, topic_id, current_level
	) SELECT $1, id, 3 FROM public.topics LIMIT 1 RETURNING id`, intakeID).Scan(&id); err != nil {
		t.Fatalf("seed Group G assessment: %v", err)
	}
	return id
}

func seedGroupGPromotedClient(t *testing.T, ctx context.Context, tx pgx.Tx, registrationID, intakeID uuid.UUID) uuid.UUID {
	t.Helper()
	clientID := uuid.New()
	if _, err := tx.Exec(ctx, `INSERT INTO public.client_details (
		id, registration_form_id, intake_form_id, first_name, last_name, email,
		gender, filenumber, street, house_number, postal_code, city
	) VALUES ($1, $2, $3, 'Group', 'G', $4, 'unknown', $5, 'Test', '1', '1000AA', 'Test')`,
		clientID, registrationID, intakeID, uuid.NewString()+"@test.invalid", uuid.NewString()); err != nil {
		t.Fatalf("seed Group G promoted client: %v", err)
	}
	return clientID
}

func assertGroupGRowVisible(t *testing.T, ctx context.Context, tx pgx.Tx, table string, id uuid.UUID, want bool) {
	t.Helper()
	var got bool
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM public.%s WHERE id = $1)`, pgx.Identifier{table}.Sanitize())
	if err := tx.QueryRow(ctx, query, id).Scan(&got); err != nil {
		t.Fatalf("read %s visibility: %v", table, err)
	}
	if got != want {
		t.Fatalf("%s row %s visibility = %t, want %t", table, id, got, want)
	}
}

func assertGroupGTablesForced(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	for _, table := range []string{
		"registration_form", "intake_forms", "intake_topic_assessments",
		"collaboration_agreement", "risk_assessment", "consent_declaration",
		"youth_care_intake", "data_sharing_statement",
	} {
		var enabled, forced bool
		if err := tx.QueryRow(ctx, `SELECT relrowsecurity, relforcerowsecurity
			FROM pg_catalog.pg_class WHERE oid = ('public.' || $1)::regclass`, table).Scan(&enabled, &forced); err != nil {
			t.Fatalf("inspect %s RLS: %v", table, err)
		}
		if !enabled || !forced {
			t.Fatalf("%s RLS enabled=%t forced=%t", table, enabled, forced)
		}
	}
}
