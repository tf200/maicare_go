package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func TestClientGroupHPermissionPolicies(t *testing.T) {
	store := openIntegrationStore(t)
	ctx := context.Background()
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	fullActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	for _, permission := range []string{
		"CLIENT.STATUS.UPDATE", "CLIENT.INVOLVED_EMPLOYEE.VIEW",
		"APPOINTMENT_CARD.VIEW", "APPOINTMENT_CARD.UPDATE", "APPOINTMENT_CARD.DELETE",
	} {
		grantGroupHScopedPermission(t, ctx, tx, fullActor.userID, permission, "assigned")
	}
	contractActor := seedClientRLSActor(t, ctx, tx, "", false)
	grantGroupHScopedPermission(t, ctx, tx, contractActor.userID, "CONTRACT.VIEW", "assigned")
	outsideActor := seedClientRLSActor(t, ctx, tx, "", false)

	clientA := seedClientRLSClient(t, ctx, tx)
	clientB := seedClientRLSClient(t, ctx, tx)
	seedClientRLSAssignment(t, ctx, tx, clientA, fullActor.employeeID)
	seedClientRLSAssignment(t, ctx, tx, clientA, contractActor.employeeID)

	historyA := seedGroupHStatusHistory(t, ctx, tx, clientA)
	historyB := seedGroupHStatusHistory(t, ctx, tx, clientB)
	transferA := seedGroupHLocationTransfer(t, ctx, tx, clientA)
	transferB := seedGroupHLocationTransfer(t, ctx, tx, clientB)
	legacyAssignmentA := seedGroupHLegacyAssignment(t, ctx, tx, clientA, fullActor.employeeID)
	legacyAssignmentB := seedGroupHLegacyAssignment(t, ctx, tx, clientB, fullActor.employeeID)
	cardA := seedGroupHAppointmentCard(t, ctx, tx, clientA)
	cardB := seedGroupHAppointmentCard(t, ctx, tx, clientB)

	eventA := seedGroupHEvent(t, ctx, tx, fullActor.employeeID)
	attendeeAEmployee := seedGroupHEmployeeAttendee(t, ctx, tx, eventA, fullActor.employeeID)
	attendeeAClient := seedGroupHClientAttendee(t, ctx, tx, eventA, clientA)
	attendeeAEmail := seedGroupHEmailAttendee(t, ctx, tx, eventA)
	eventB := seedGroupHEvent(t, ctx, tx, outsideActor.employeeID)
	attendeeBClient := seedGroupHClientAttendee(t, ctx, tx, eventB, clientB)
	attendeeBEmail := seedGroupHEmailAttendee(t, ctx, tx, eventB)
	eventC := seedGroupHEvent(t, ctx, tx, contractActor.employeeID)
	seedGroupHEmployeeAttendee(t, ctx, tx, eventC, fullActor.employeeID)

	runtimeRole := "phase9_group_h_runtime_" + uuid.NewString()
	roleName := pgx.Identifier{runtimeRole}.Sanitize()
	if _, err := tx.Exec(ctx, fmt.Sprintf(`CREATE ROLE %s NOLOGIN NOSUPERUSER NOBYPASSRLS`, roleName)); err != nil {
		t.Fatalf("create runtime role: %v", err)
	}
	for _, table := range []string{
		"public.client_status_history", "public.client_location_transfer", "public.assignment",
		"public.calendar_event_attendees", "public.appointment_card", "public.calendar_events",
		"public.client_details", "public.employee_profile",
	} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON %s TO %s`, table, roleName)); err != nil {
			t.Fatalf("grant privileges on %s: %v", table, err)
		}
	}
	for _, signature := range []string{
		"public.get_current_user_id()", "public.get_current_employee_id()", "public.has_permission(text)",
		"public.get_permission_scope(text)", "public.is_assigned_to_client(uuid)", "public.can_access_client(uuid,text)",
		"public.can_participate_in_event(uuid)", "public.can_organize_event(uuid)",
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
		t.Fatal("Group H runtime role unexpectedly bypasses RLS")
	}
	assertGroupHTablesForced(t, ctx, tx)

	setAuthorizationActor(t, ctx, tx, fullActor.userID.String(), fullActor.employeeID.String())

	assertGroupHRowVisible(t, ctx, tx, "client_status_history", historyA, true)
	assertGroupHRowVisible(t, ctx, tx, "client_status_history", historyB, false)
	assertGroupHWriteApplied(t, ctx, tx, "group_h_history_update_immutable", false, func() (int64, error) {
		tag, err := tx.Exec(ctx, `UPDATE public.client_status_history SET reason = 'rewritten' WHERE id = $1`, historyA)
		return tag.RowsAffected(), err
	})
	assertGroupHWriteApplied(t, ctx, tx, "group_h_history_delete_immutable", false, func() (int64, error) {
		tag, err := tx.Exec(ctx, `DELETE FROM public.client_status_history WHERE id = $1`, historyA)
		return tag.RowsAffected(), err
	})
	assertRLSSavepointOperation(t, ctx, tx, "group_h_history_insert_owned", true, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.client_status_history (client_id, old_status, new_status, reason)
			VALUES ($1, NULL, 'in_care', 'group_h_test')`, clientA)
		return err
	})
	assertRLSSavepointOperation(t, ctx, tx, "group_h_history_insert_unowned", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.client_status_history (client_id, old_status, new_status, reason)
			VALUES ($1, NULL, 'in_care', 'group_h_test')`, clientB)
		return err
	})

	assertGroupHRowVisible(t, ctx, tx, "client_location_transfer", transferA, true)
	assertGroupHRowVisible(t, ctx, tx, "client_location_transfer", transferB, false)
	assertGroupHWriteApplied(t, ctx, tx, "group_h_transfer_approve_owned", true, func() (int64, error) {
		tag, err := tx.Exec(ctx, `UPDATE public.client_location_transfer
			SET status = 'approved', approved_rejected_at = NOW() WHERE id = $1`, transferA)
		return tag.RowsAffected(), err
	})
	assertGroupHWriteApplied(t, ctx, tx, "group_h_transfer_update_unowned", false, func() (int64, error) {
		tag, err := tx.Exec(ctx, `UPDATE public.client_location_transfer SET reason = 'moved' WHERE id = $1`, transferB)
		return tag.RowsAffected(), err
	})
	assertRLSSavepointOperation(t, ctx, tx, "group_h_transfer_insert_unowned", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.client_location_transfer (client_id, reason)
			VALUES ($1, 'group_h_test')`, clientB)
		return err
	})

	assertGroupHRowVisible(t, ctx, tx, "assignment", legacyAssignmentA, true)
	assertGroupHRowVisible(t, ctx, tx, "assignment", legacyAssignmentB, false)
	assertRLSSavepointOperation(t, ctx, tx, "group_h_legacy_assignment_insert_denied_without_all_scope", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.assignment (employee_id, client_id, start_datetime, end_datetime, status)
			VALUES ($1, $2, NOW(), NOW() + INTERVAL '1 day', 'Confirmed')`, fullActor.employeeID, clientB)
		return err
	})
	assertGroupHWriteApplied(t, ctx, tx, "group_h_legacy_assignment_delete_denied_without_all_scope", false, func() (int64, error) {
		tag, err := tx.Exec(ctx, `DELETE FROM public.assignment WHERE id = $1`, legacyAssignmentA)
		return tag.RowsAffected(), err
	})

	assertGroupHRowVisible(t, ctx, tx, "appointment_card", cardA, true)
	assertGroupHRowVisible(t, ctx, tx, "appointment_card", cardB, false)
	assertGroupHWriteApplied(t, ctx, tx, "group_h_card_update_owned", true, func() (int64, error) {
		tag, err := tx.Exec(ctx, `UPDATE public.appointment_card SET travel = ARRAY['card updated'] WHERE id = $1`, cardA)
		return tag.RowsAffected(), err
	})
	assertRLSSavepointOperation(t, ctx, tx, "group_h_card_upsert_unowned", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.appointment_card (client_id) VALUES ($1)`, clientB)
		return err
	})
	assertGroupHWriteApplied(t, ctx, tx, "group_h_card_update_unowned", false, func() (int64, error) {
		tag, err := tx.Exec(ctx, `UPDATE public.appointment_card SET travel = ARRAY['hijacked'] WHERE id = $1`, cardB)
		return tag.RowsAffected(), err
	})
	assertGroupHWriteApplied(t, ctx, tx, "group_h_card_delete_unowned", false, func() (int64, error) {
		tag, err := tx.Exec(ctx, `DELETE FROM public.appointment_card WHERE id = $1`, cardB)
		return tag.RowsAffected(), err
	})
	assertGroupHWriteApplied(t, ctx, tx, "group_h_card_delete_owned", true, func() (int64, error) {
		tag, err := tx.Exec(ctx, `DELETE FROM public.appointment_card WHERE id = $1`, cardA)
		return tag.RowsAffected(), err
	})

	assertGroupHRowVisible(t, ctx, tx, "calendar_event_attendees", attendeeAEmployee, true)
	assertGroupHRowVisible(t, ctx, tx, "calendar_event_attendees", attendeeAClient, true)
	assertGroupHRowVisible(t, ctx, tx, "calendar_event_attendees", attendeeAEmail, true)
	assertGroupHRowVisible(t, ctx, tx, "calendar_event_attendees", attendeeBClient, false)
	assertGroupHRowVisible(t, ctx, tx, "calendar_event_attendees", attendeeBEmail, false)

	setAuthorizationActor(t, ctx, tx, contractActor.userID.String(), contractActor.employeeID.String())
	assertGroupHRowVisible(t, ctx, tx, "calendar_event_attendees", attendeeAClient, true)
	assertGroupHRowVisible(t, ctx, tx, "calendar_event_attendees", attendeeAEmployee, false)
	assertGroupHRowVisible(t, ctx, tx, "calendar_event_attendees", attendeeAEmail, false)
	assertGroupHRowVisible(t, ctx, tx, "calendar_event_attendees", attendeeBClient, false)

	setAuthorizationActor(t, ctx, tx, fullActor.userID.String(), fullActor.employeeID.String())
	assertRLSSavepointOperation(t, ctx, tx, "group_h_attendee_insert_owned_event", true, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.calendar_event_attendees (event_id, email)
			VALUES ($1, $2)`, eventA, uuid.NewString()+"@test.invalid")
		return err
	})
	assertRLSSavepointOperation(t, ctx, tx, "group_h_attendee_insert_unreachable_client", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.calendar_event_attendees (event_id, client_id)
			VALUES ($1, $2)`, eventA, clientB)
		return err
	})
	assertRLSSavepointOperation(t, ctx, tx, "group_h_attendee_insert_other_event", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.calendar_event_attendees (event_id, email)
			VALUES ($1, $2)`, eventB, uuid.NewString()+"@test.invalid")
		return err
	})
	assertRLSSavepointOperation(t, ctx, tx, "group_h_attendee_insert_by_mere_participant", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.calendar_event_attendees (event_id, email)
			VALUES ($1, $2)`, eventC, uuid.NewString()+"@test.invalid")
		return err
	})
	assertGroupHWriteApplied(t, ctx, tx, "group_h_attendee_delete_owned_event", true, func() (int64, error) {
		tag, err := tx.Exec(ctx, `DELETE FROM public.calendar_event_attendees WHERE id = $1`, attendeeAEmail)
		return tag.RowsAffected(), err
	})
	assertGroupHWriteApplied(t, ctx, tx, "group_h_attendee_delete_other_event", false, func() (int64, error) {
		tag, err := tx.Exec(ctx, `DELETE FROM public.calendar_event_attendees WHERE id = $1`, attendeeBEmail)
		return tag.RowsAffected(), err
	})

	setAuthorizationActor(t, ctx, tx, outsideActor.userID.String(), outsideActor.employeeID.String())
	assertRLSSavepointOperation(t, ctx, tx, "group_h_attendee_insert_own_organized_event", true, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.calendar_event_attendees (event_id, email)
			VALUES ($1, $2)`, eventB, uuid.NewString()+"@test.invalid")
		return err
	})

	setAuthorizationActor(t, ctx, tx, "", "")
	assertGroupHRowVisible(t, ctx, tx, "client_status_history", historyA, false)
	assertGroupHRowVisible(t, ctx, tx, "client_location_transfer", transferA, false)
	assertGroupHRowVisible(t, ctx, tx, "assignment", legacyAssignmentA, false)
	assertGroupHRowVisible(t, ctx, tx, "calendar_event_attendees", attendeeAClient, false)
	assertGroupHRowVisible(t, ctx, tx, "appointment_card", cardB, false)
}

func assertGroupHWriteApplied(t *testing.T, ctx context.Context, tx pgx.Tx, name string, want bool, operation func() (int64, error)) {
	t.Helper()
	if _, err := tx.Exec(ctx, "SAVEPOINT "+name); err != nil {
		t.Fatalf("create %s savepoint: %v", name, err)
	}
	affected, err := operation()
	applied := err == nil && affected > 0
	if applied != want {
		t.Fatalf("%s applied = %t, want %t (affected=%d error: %v)", name, applied, want, affected, err)
	}
	if _, rollbackErr := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT "+name); rollbackErr != nil {
		t.Fatalf("rollback %s: %v", name, rollbackErr)
	}
}

func grantGroupHScopedPermission(t *testing.T, ctx context.Context, tx pgx.Tx, userID uuid.UUID, permission, scope string) {
	t.Helper()
	permissionID := ensureClientRLSPermission(t, ctx, tx, permission, true)
	if _, err := tx.Exec(ctx, `INSERT INTO public.role_permissions (role_id, permission_id, scope)
		SELECT role_id, $2, $3::public.permission_scope_enum FROM public.user_roles WHERE user_id = $1`,
		userID, permissionID, scope); err != nil {
		t.Fatalf("grant %s: %v", permission, err)
	}
}

func seedGroupHStatusHistory(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.client_status_history (client_id, old_status, new_status, reason)
		VALUES ($1, NULL, 'scheduled_in_care', 'group_h_seed') RETURNING id`, clientID).Scan(&id); err != nil {
		t.Fatalf("seed status history: %v", err)
	}
	return id
}

func seedGroupHLocationTransfer(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.client_location_transfer (client_id, request_date, reason)
		VALUES ($1, NOW(), 'group_h_seed') RETURNING id`, clientID).Scan(&id); err != nil {
		t.Fatalf("seed location transfer: %v", err)
	}
	return id
}

func seedGroupHLegacyAssignment(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, employeeID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.assignment (employee_id, client_id, start_datetime, end_datetime, status)
		VALUES ($1, $2, NOW(), NOW() + INTERVAL '30 days', 'Confirmed') RETURNING id`, employeeID, clientID).Scan(&id); err != nil {
		t.Fatalf("seed legacy assignment: %v", err)
	}
	return id
}

func seedGroupHAppointmentCard(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.appointment_card (client_id, general_information)
		VALUES ($1, ARRAY['group h seed']) RETURNING id`, clientID).Scan(&id); err != nil {
		t.Fatalf("seed appointment card: %v", err)
	}
	return id
}

func seedGroupHEvent(t *testing.T, ctx context.Context, tx pgx.Tx, organizerEmployeeID uuid.UUID) uuid.UUID {
	t.Helper()
	start := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.calendar_events (
		organizer_employee_id, created_by_employee_id, kind, title, start_at, end_at
	) VALUES ($1, $1, 'appointment', 'group h seed', $2, $3) RETURNING id`,
		organizerEmployeeID, start, start.Add(time.Hour)).Scan(&id); err != nil {
		t.Fatalf("seed calendar event: %v", err)
	}
	return id
}

func seedGroupHEmployeeAttendee(t *testing.T, ctx context.Context, tx pgx.Tx, eventID, employeeID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.calendar_event_attendees (event_id, employee_id)
		VALUES ($1, $2) RETURNING id`, eventID, employeeID).Scan(&id); err != nil {
		t.Fatalf("seed employee attendee: %v", err)
	}
	return id
}

func seedGroupHClientAttendee(t *testing.T, ctx context.Context, tx pgx.Tx, eventID, clientID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.calendar_event_attendees (event_id, client_id)
		VALUES ($1, $2) RETURNING id`, eventID, clientID).Scan(&id); err != nil {
		t.Fatalf("seed client attendee: %v", err)
	}
	return id
}

func seedGroupHEmailAttendee(t *testing.T, ctx context.Context, tx pgx.Tx, eventID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.calendar_event_attendees (event_id, email)
		VALUES ($1, $2) RETURNING id`, eventID, uuid.NewString()+"@test.invalid").Scan(&id); err != nil {
		t.Fatalf("seed email attendee: %v", err)
	}
	return id
}

func assertGroupHRowVisible(t *testing.T, ctx context.Context, tx pgx.Tx, table string, id uuid.UUID, want bool) {
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

func assertGroupHTablesForced(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	for _, table := range []string{
		"client_status_history", "client_location_transfer", "assignment",
		"calendar_event_attendees", "appointment_card",
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
