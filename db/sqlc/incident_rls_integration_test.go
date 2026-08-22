package db

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestIncidentPermissionScopePolicies(t *testing.T) {
	store := openIntegrationStore(t)
	ctx := context.Background()
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	assignedActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	allActor := seedClientRLSActor(t, ctx, tx, "all", false)
	confirmActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	deniedActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	for _, permission := range []string{
		"CLIENT.INCIDENT.VIEW", "CLIENT.INCIDENT.CREATE", "CLIENT.INCIDENT.UPDATE", "CLIENT.INCIDENT.DELETE", "CLIENT.INCIDENT.CONFIRM",
	} {
		grantClientRLSPermission(t, ctx, tx, assignedActor.userID, permission, "assigned")
		grantClientRLSPermission(t, ctx, tx, allActor.userID, permission, "all")
	}
	grantClientRLSPermission(t, ctx, tx, confirmActor.userID, "CLIENT.INCIDENT.CONFIRM", "assigned")
	grantClientRLSPermission(t, ctx, tx, confirmActor.userID, "CLIENT.INCIDENT.VIEW", "assigned")

	assignedClient := seedClientRLSClient(t, ctx, tx)
	unassignedClient := seedClientRLSClient(t, ctx, tx)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, assignedActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, assignedClient, confirmActor.employeeID)
	locationID := seedIncidentLocation(t, ctx, tx)
	assignedIncident := seedIncidentFixture(t, ctx, tx, assignedClient, assignedActor.employeeID, locationID)
	unassignedIncident := seedIncidentFixture(t, ctx, tx, unassignedClient, assignedActor.employeeID, locationID)

	runtimeRole := "phase9_group_d_runtime_" + uuid.NewString()
	if _, err := tx.Exec(ctx, fmt.Sprintf(`CREATE ROLE %s NOLOGIN NOSUPERUSER NOBYPASSRLS`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("create runtime role: %v", err)
	}
	for _, table := range []string{"public.incident", "public.client_details", "public.employee_profile", "public.custom_user", "public.location"} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT ON %s TO %s`, table, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
			t.Fatalf("grant select on %s: %v", table, err)
		}
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT INSERT, UPDATE, DELETE ON public.incident TO %s`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("grant incident mutations: %v", err)
	}
	for _, signature := range []string{
		"public.get_current_user_id()",
		"public.get_current_employee_id()",
		"public.has_permission(text)",
		"public.get_permission_scope(text)",
		"public.is_assigned_to_client(uuid)",
		"public.can_access_client(uuid,text)",
		"public.can_read_created_client(uuid)",
		"public.begin_incident_creation(uuid)",
		"public.can_read_created_incident(uuid)",
		"public.confirm_incident(uuid)",
		"public.claim_incident_confirmation_email(uuid)",
		"public.release_incident_confirmation_email(uuid,uuid)",
		"public.mark_incident_confirmation_email_sent(uuid,uuid)",
	} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT EXECUTE ON FUNCTION %s TO %s`, signature, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
			t.Fatalf("grant helper %s: %v", signature, err)
		}
	}
	if _, err := tx.Exec(ctx, fmt.Sprintf(`SET LOCAL ROLE %s`, pgx.Identifier{runtimeRole}.Sanitize())); err != nil {
		t.Fatalf("set runtime role: %v", err)
	}
	assertIncidentTableForced(t, ctx, tx)

	setAuthorizationActor(t, ctx, tx, assignedActor.userID.String(), assignedActor.employeeID.String())
	assertIncidentVisible(t, ctx, tx, assignedIncident, true)
	assertIncidentVisible(t, ctx, tx, unassignedIncident, false)
	if _, err := New(tx).GetIncident(ctx, assignedIncident); err != nil {
		t.Fatalf("get assigned incident for file flow: %v", err)
	}
	assertIncidentUpdate(t, ctx, tx, assignedIncident, uuid.New(), true)
	assertRLSSavepointOperation(t, ctx, tx, "direct_incident_confirmation", false, func() error {
		_, err := tx.Exec(ctx, `UPDATE public.incident SET is_confirmed = TRUE, confirmed_at = NOW() WHERE id = $1`, assignedIncident)
		return err
	})
	assertIncidentUpdate(t, ctx, tx, unassignedIncident, uuid.New(), false)
	createdIncident := createIncidentWithQuery(t, ctx, tx, assignedClient, locationID, assignedActor.employeeID)
	forgedStateIncident := createIncidentWithForgedConfirmationState(t, ctx, tx, assignedClient, locationID, assignedActor.userID)
	createIncidentAsRuntime(t, ctx, tx, unassignedClient, locationID, uuid.New(), false)
	assertIncidentDelete(t, ctx, tx, createdIncident, true)
	assertIncidentDelete(t, ctx, tx, forgedStateIncident, true)
	if _, err := New(tx).DeleteIncident(ctx, createdIncident); err == nil {
		t.Fatal("second incident delete unexpectedly succeeded")
	} else if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("second incident delete error = %v, want pgx.ErrNoRows", err)
	}
	assertIncidentCounts(t, ctx, tx, 1)

	setAuthorizationActor(t, ctx, tx, allActor.userID.String(), allActor.employeeID.String())
	assertIncidentVisible(t, ctx, tx, assignedIncident, true)
	assertIncidentVisible(t, ctx, tx, unassignedIncident, true)
	assertIncidentCounts(t, ctx, tx, 2)

	setAuthorizationActor(t, ctx, tx, confirmActor.userID.String(), confirmActor.employeeID.String())
	assertIncidentVisible(t, ctx, tx, assignedIncident, true)
	if affected, err := New(tx).ConfirmIncident(ctx, assignedIncident); err != nil {
		t.Fatalf("confirm assigned incident: %v", err)
	} else if affected != 1 {
		t.Fatalf("confirmed rows = %d, want 1", affected)
	}
	if incident, err := New(tx).GetIncident(ctx, assignedIncident); err != nil {
		t.Fatalf("load confirmed incident for delegated worker: %v", err)
	} else if incident.ClientID != assignedClient {
		t.Fatalf("confirmed incident client = %s, want %s", incident.ClientID, assignedClient)
	}
	claimToken, err := New(tx).ClaimIncidentConfirmationEmail(ctx, assignedIncident)
	if err != nil {
		t.Fatalf("claim incident confirmation email: %v", err)
	} else if claimToken == uuid.Nil {
		t.Fatal("incident confirmation email claim returned no token")
	}
	if secondToken, err := New(tx).ClaimIncidentConfirmationEmail(ctx, assignedIncident); err != nil {
		t.Fatalf("reclaim incident confirmation email: %v", err)
	} else if secondToken != uuid.Nil {
		t.Fatalf("reclaimed confirmation email token = %s, want nil UUID", secondToken)
	}
	if affected, err := New(tx).MarkIncidentConfirmationEmailSent(ctx, MarkIncidentConfirmationEmailSentParams{IncidentID: assignedIncident, ClaimToken: uuid.New()}); err != nil {
		t.Fatalf("mark incident confirmation email with wrong token: %v", err)
	} else if affected != 0 {
		t.Fatalf("wrong-token marked rows = %d, want 0", affected)
	}
	if affected, err := New(tx).MarkIncidentConfirmationEmailSent(ctx, MarkIncidentConfirmationEmailSentParams{IncidentID: assignedIncident, ClaimToken: claimToken}); err != nil {
		t.Fatalf("mark incident confirmation email sent: %v", err)
	} else if affected != 1 {
		t.Fatalf("marked confirmation email rows = %d, want 1", affected)
	}
	if affected, err := New(tx).ConfirmIncident(ctx, unassignedIncident); err != nil {
		t.Fatalf("confirm unassigned incident: %v", err)
	} else if affected != 0 {
		t.Fatalf("unassigned confirmed rows = %d, want 0", affected)
	}
	assertIncidentUpdate(t, ctx, tx, assignedIncident, confirmActor.employeeID, false)

	setAuthorizationActor(t, ctx, tx, deniedActor.userID.String(), deniedActor.employeeID.String())
	assertIncidentVisible(t, ctx, tx, assignedIncident, false)
	assertIncidentCounts(t, ctx, tx, 0)
}

func createIncidentWithQuery(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, locationID, employeeID uuid.UUID) uuid.UUID {
	t.Helper()
	incident, err := New(tx).CreateIncident(ctx, CreateIncidentParams{
		EmployeeID:          uuid.New(),
		LocationID:          locationID,
		ReporterInvolvement: IncidentReporterInvolvementEnumWitness,
		InformedParties:     []InformedPartyEnum{},
		OccurredAt:          pgtype.Timestamptz{Time: time.Now(), Valid: true},
		IncidentType:        IncidentTypeEnumAccident,
		SeverityOfIncident:  SeverityOfIncidentEnumLessSerious,
		RecurrenceRisk:      RecurrenceRiskEnumVeryLow,
		CauseCategories:     []IncidentCauseCategoryEnum{},
		PhysicalInjury:      PhysicalInjuryEnumNoInjuries,
		PsychologicalDamage: PsychologicalDamageEnumNo,
		NeededConsultation:  NeededConsultationEnumNo,
		FollowUpActions:     []IncidentFollowUpActionEnum{},
		ClientID:            clientID,
		Emails:              []string{},
	})
	if err != nil {
		t.Fatalf("create incident with production query: %v", err)
	}
	if incident.EmployeeID != employeeID {
		t.Fatalf("production incident reporter = %s, want actor %s", incident.EmployeeID, employeeID)
	}
	return incident.ID
}

func createIncidentWithForgedConfirmationState(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, locationID, userID uuid.UUID) uuid.UUID {
	t.Helper()
	var incidentID uuid.UUID
	var isConfirmed bool
	var confirmedAt, claimedAt, sentAt *time.Time
	var confirmedBy, claimToken *uuid.UUID
	err := tx.QueryRow(ctx, `WITH new_incident AS (
		SELECT public.begin_incident_creation($1) AS id
	) INSERT INTO public.incident (
		id, employee_id, location_id, client_id, reporter_involvement, occurred_at, incident_type,
		severity_of_incident, recurrence_risk, physical_injury, psychological_damage, needed_consultation,
		is_confirmed, confirmed_at, confirmed_by, confirmation_email_claim_token,
		confirmation_email_claimed_at, confirmation_email_sent_at
	) SELECT id, $2, $3, $1, 'witness', NOW(), 'accident', 'less_serious', 'very_low',
		'no_injuries', 'no', 'no', TRUE, NOW(), $4, $5, NOW(), NOW()
	FROM new_incident WHERE id IS NOT NULL
	RETURNING id, is_confirmed, confirmed_at, confirmed_by, confirmation_email_claim_token,
		confirmation_email_claimed_at, confirmation_email_sent_at`,
		clientID, uuid.New(), locationID, userID, uuid.New()).Scan(
		&incidentID, &isConfirmed, &confirmedAt, &confirmedBy, &claimToken, &claimedAt, &sentAt)
	if err != nil {
		t.Fatalf("create incident with forged confirmation state: %v", err)
	}
	if isConfirmed || confirmedAt != nil || confirmedBy != nil || claimToken != nil || claimedAt != nil || sentAt != nil {
		t.Fatalf("incident accepted forged confirmation state: confirmed=%t at=%v by=%v token=%v claimed=%v sent=%v",
			isConfirmed, confirmedAt, confirmedBy, claimToken, claimedAt, sentAt)
	}
	return incidentID
}

func seedIncidentLocation(t *testing.T, ctx context.Context, tx pgx.Tx) uuid.UUID {
	t.Helper()
	organisationID := uuid.New()
	locationID := uuid.New()
	if _, err := tx.Exec(ctx, `INSERT INTO public.organisations (id, name, street, house_number, postal_code, city)
		VALUES ($1, 'Phase Nine', 'Test', '1', '1000AA', 'Test')`, organisationID); err != nil {
		t.Fatalf("seed incident organisation: %v", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO public.location (id, organisation_id, name, street, house_number, postal_code, city)
		VALUES ($1, $2, 'Phase Nine', 'Test', '1', '1000AA', 'Test')`, locationID, organisationID); err != nil {
		t.Fatalf("seed incident location: %v", err)
	}
	return locationID
}

func seedIncidentFixture(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, employeeID, locationID uuid.UUID) uuid.UUID {
	t.Helper()
	incidentID := uuid.New()
	_, err := tx.Exec(ctx, `INSERT INTO public.incident (
		id, employee_id, location_id, client_id, reporter_involvement, occurred_at, incident_type,
		severity_of_incident, recurrence_risk, physical_injury, psychological_damage, needed_consultation
	) VALUES ($1, $2, $3, $4, 'witness', NOW(), 'accident', 'less_serious', 'very_low', 'no_injuries', 'no', 'no')`,
		incidentID, employeeID, locationID, clientID)
	if err != nil {
		t.Fatalf("seed incident: %v", err)
	}
	return incidentID
}

func createIncidentAsRuntime(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, locationID, forgedEmployeeID uuid.UUID, want bool) uuid.UUID {
	t.Helper()
	var incidentID, employeeID uuid.UUID
	err := tx.QueryRow(ctx, `WITH new_incident AS (
		SELECT public.begin_incident_creation($1) AS id
	) INSERT INTO public.incident (
		id, employee_id, location_id, client_id, reporter_involvement, occurred_at, incident_type,
		severity_of_incident, recurrence_risk, physical_injury, psychological_damage, needed_consultation
	) SELECT id, $2, $3, $1, 'witness', NOW(), 'accident', 'less_serious', 'very_low', 'no_injuries', 'no', 'no'
	FROM new_incident WHERE id IS NOT NULL RETURNING id, employee_id`, clientID, forgedEmployeeID, locationID).Scan(&incidentID, &employeeID)
	if !want {
		if err == nil {
			t.Fatal("unauthorized incident creation unexpectedly succeeded")
		}
		return uuid.Nil
	}
	if err != nil {
		t.Fatalf("create incident as runtime: %v", err)
	}
	var actorEmployeeID uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT current_setting('myapp.current_employee_id')::uuid`).Scan(&actorEmployeeID); err != nil {
		t.Fatalf("read actor employee: %v", err)
	}
	if employeeID != actorEmployeeID {
		t.Fatalf("incident reporter = %s, want actor %s", employeeID, actorEmployeeID)
	}
	return incidentID
}

func assertIncidentTableForced(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	var enabled, forced bool
	if err := tx.QueryRow(ctx, `SELECT relrowsecurity, relforcerowsecurity FROM pg_class WHERE oid = 'public.incident'::regclass`).Scan(&enabled, &forced); err != nil {
		t.Fatalf("read incident RLS flags: %v", err)
	}
	if !enabled || !forced {
		t.Fatalf("incident RLS enabled=%t forced=%t, want true/true", enabled, forced)
	}
}

func assertIncidentVisible(t *testing.T, ctx context.Context, tx pgx.Tx, incidentID uuid.UUID, want bool) {
	t.Helper()
	var got bool
	if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM public.incident WHERE id = $1)`, incidentID).Scan(&got); err != nil {
		t.Fatalf("check incident visibility: %v", err)
	}
	if got != want {
		t.Fatalf("incident %s visible = %t, want %t", incidentID, got, want)
	}
}

func assertIncidentUpdate(t *testing.T, ctx context.Context, tx pgx.Tx, incidentID, forgedEmployeeID uuid.UUID, want bool) {
	t.Helper()
	tag, err := tx.Exec(ctx, `UPDATE public.incident SET additional_details = 'updated', employee_id = $2 WHERE id = $1`, incidentID, forgedEmployeeID)
	if err != nil {
		t.Fatalf("update incident: %v", err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("incident %s updated = %t, want %t", incidentID, got, want)
	}
	if want {
		var employeeID uuid.UUID
		if err := tx.QueryRow(ctx, `SELECT employee_id FROM public.incident WHERE id = $1`, incidentID).Scan(&employeeID); err != nil {
			t.Fatalf("read incident reporter: %v", err)
		}
		if employeeID == forgedEmployeeID {
			t.Fatalf("incident accepted forged reporter %s", forgedEmployeeID)
		}
	}
}

func assertIncidentDelete(t *testing.T, ctx context.Context, tx pgx.Tx, incidentID uuid.UUID, want bool) {
	t.Helper()
	tag, err := tx.Exec(ctx, `DELETE FROM public.incident WHERE id = $1`, incidentID)
	if err != nil {
		t.Fatalf("delete incident: %v", err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("incident %s deleted = %t, want %t", incidentID, got, want)
	}
}

func assertIncidentCounts(t *testing.T, ctx context.Context, tx pgx.Tx, want int64) {
	t.Helper()
	counts, err := New(tx).GetIncidentCounts(ctx)
	if err != nil {
		t.Fatalf("get incident counts: %v", err)
	}
	if counts.Past24hCount != want {
		t.Fatalf("past 24h incident count = %d, want %d", counts.Past24hCount, want)
	}
}
