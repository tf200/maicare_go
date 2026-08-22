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

func TestClientGroupFPermissionScopePolicies(t *testing.T) {
	store := openIntegrationStore(t)
	ctx := context.Background()
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	assignedActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	allActor := seedClientRLSActor(t, ctx, tx, "all", false)
	clientOnlyActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	invoiceViewActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	paymentViewActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	invoiceCreateActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	invoiceUpdateActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	paymentCreateActor := seedClientRLSActor(t, ctx, tx, "assigned", false)

	groupFPermissions := []string{
		"CONTRACT.VIEW", "CONTRACT.CREATE", "CONTRACT.UPDATE", "CONTRACT.DELETE",
		"INVOICE.VIEW", "INVOICE.CREATE", "INVOICE.UPDATE", "INVOICE.DELETE",
		"INVOICE.PAYMENT.VIEW", "INVOICE.PAYMENT.CREATE", "INVOICE.PAYMENT.UPDATE", "INVOICE.PAYMENT.DELETE",
	}
	for _, permission := range groupFPermissions {
		grantClientRLSPermission(t, ctx, tx, assignedActor.userID, permission, "assigned")
		grantClientRLSPermission(t, ctx, tx, allActor.userID, permission, "all")
	}
	grantClientRLSPermission(t, ctx, tx, invoiceViewActor.userID, "INVOICE.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, paymentViewActor.userID, "INVOICE.PAYMENT.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, invoiceCreateActor.userID, "INVOICE.CREATE", "assigned")
	grantClientRLSPermission(t, ctx, tx, invoiceCreateActor.userID, "CONTRACT.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, invoiceUpdateActor.userID, "INVOICE.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, invoiceUpdateActor.userID, "INVOICE.UPDATE", "assigned")
	grantClientRLSPermission(t, ctx, tx, paymentCreateActor.userID, "INVOICE.PAYMENT.CREATE", "assigned")

	assignedClient := seedClientRLSClient(t, ctx, tx)
	unassignedClient := seedClientRLSClient(t, ctx, tx)
	for _, actor := range []clientRLSActor{
		assignedActor, clientOnlyActor, invoiceViewActor, paymentViewActor, invoiceCreateActor, invoiceUpdateActor, paymentCreateActor,
	} {
		seedClientRLSAssignment(t, ctx, tx, assignedClient, actor.employeeID)
	}

	senderID := seedGroupFSender(t, ctx, tx)
	assignedContract := seedGroupFContract(t, ctx, tx, assignedClient, senderID)
	unassignedContract := seedGroupFContract(t, ctx, tx, unassignedClient, senderID)
	assignedInvoice := seedGroupFInvoice(t, ctx, tx, assignedClient, senderID)
	unassignedInvoice := seedGroupFInvoice(t, ctx, tx, unassignedClient, senderID)
	assignedLine := seedGroupFInvoiceLine(t, ctx, tx, assignedInvoice, assignedClient, senderID, assignedContract)
	unassignedLine := seedGroupFInvoiceLine(t, ctx, tx, unassignedInvoice, unassignedClient, senderID, unassignedContract)
	assignedPayment := seedGroupFPayment(t, ctx, tx, assignedInvoice)
	unassignedPayment := seedGroupFPayment(t, ctx, tx, unassignedInvoice)
	assignedAgreement := seedGroupFContractChild(t, ctx, tx, "client_agreement", assignedContract)
	unassignedAgreement := seedGroupFContractChild(t, ctx, tx, "client_agreement", unassignedContract)
	if _, err := tx.Exec(ctx, `INSERT INTO public.framework_agreement (client_id, agreement_details)
		VALUES ($1, 'Group F framework')`, assignedClient); err != nil {
		t.Fatalf("seed Group F framework agreement: %v", err)
	}
	retainedAuditClient := seedClientRLSClient(t, ctx, tx)
	retainedContract := seedGroupFContract(t, ctx, tx, retainedAuditClient, senderID)
	retainedInvoice := seedGroupFInvoice(t, ctx, tx, retainedAuditClient, senderID)
	if _, err := tx.Exec(ctx, `DELETE FROM public.client_details WHERE id = $1`, retainedAuditClient); err != nil {
		t.Fatalf("delete client for audit retention check: %v", err)
	}
	for table, args := range map[string][]any{
		"contract_audit": {retainedContract},
		"invoice_audit":  {retainedInvoice},
	} {
		var retained bool
		idColumn := map[string]string{"contract_audit": "contract_id", "invoice_audit": "invoice_id"}[table]
		query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM public.%s WHERE %s = $1)`,
			pgx.Identifier{table}.Sanitize(), pgx.Identifier{idColumn}.Sanitize())
		if err := tx.QueryRow(ctx, query, args...).Scan(&retained); err != nil {
			t.Fatalf("check %s retention: %v", table, err)
		} else if !retained {
			t.Fatalf("%s was deleted with its client", table)
		}
	}

	runtimeRole := "phase9_group_f_runtime_" + uuid.NewString()
	roleName := pgx.Identifier{runtimeRole}.Sanitize()
	if _, err := tx.Exec(ctx, fmt.Sprintf(`CREATE ROLE %s NOLOGIN NOSUPERUSER NOBYPASSRLS`, roleName)); err != nil {
		t.Fatalf("create runtime role: %v", err)
	}
	for _, table := range []string{
		"public.employee_profile", "public.attachment_file",
		"public.contract_type", "public.contract", "public.contract_audit", "public.contract_reminder",
		"public.contract_working_hours", "public.client_agreement", "public.provision", "public.framework_agreement",
		"public.invoice", "public.invoice_audit", "public.invoice_line", "public.invoice_payment_history",
		"public.invoice_line_calendar_event", "public.billed_calendar_event", "public.invoice_run", "public.invoice_run_item",
	} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON %s TO %s`, table, roleName)); err != nil {
			t.Fatalf("grant privileges on %s: %v", table, err)
		}
	}
	for _, signature := range []string{
		"public.get_current_user_id()", "public.get_current_employee_id()", "public.has_permission(text)",
		"public.get_permission_scope(text)", "public.is_assigned_to_client(uuid)", "public.can_access_client(uuid,text)",
		"public.begin_invoice_payment_operation(uuid,text,uuid)", "public.can_access_invoice_payment_operation(uuid)",
		"public.can_access_invoice_payment_record(uuid,uuid,uuid)", "public.get_payment_operation_completed_sum(uuid)",
		"public.recalculate_invoice_payment_status(uuid)",
		"public.can_access_invoice_mutation(uuid,text)", "public.get_invoice_paid_total(uuid)",
		"public.allocate_invoice_sequence_for_date(timestamptz)", "public.attach_generated_invoice_pdf(uuid,uuid)",
		"public.can_manage_invoice_run(uuid)",
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
		t.Fatal("Group F runtime role unexpectedly bypasses RLS")
	}
	assertGroupFTablesForced(t, ctx, tx)

	setAuthorizationActor(t, ctx, tx, assignedActor.userID.String(), assignedActor.employeeID.String())
	for _, row := range []struct {
		table, idColumn string
		id              uuid.UUID
		want            bool
	}{
		{"contract", "id", assignedContract, true}, {"contract", "id", unassignedContract, false},
		{"invoice", "id", assignedInvoice, true}, {"invoice", "id", unassignedInvoice, false},
		{"invoice_line", "id", assignedLine, true}, {"invoice_line", "id", unassignedLine, false},
		{"invoice_payment_history", "id", assignedPayment, true}, {"invoice_payment_history", "id", unassignedPayment, false},
		{"client_agreement", "id", assignedAgreement, true}, {"client_agreement", "id", unassignedAgreement, false},
	} {
		assertGroupFRowVisible(t, ctx, tx, row.table, row.idColumn, row.id, row.want)
	}
	assertGroupFRowVisible(t, ctx, tx, "contract_audit", "contract_id", assignedContract, true)
	assertGroupFRowVisible(t, ctx, tx, "contract_audit", "contract_id", unassignedContract, false)
	assertGroupFRowVisible(t, ctx, tx, "invoice_audit", "invoice_id", assignedInvoice, true)
	assertGroupFRowVisible(t, ctx, tx, "invoice_audit", "invoice_id", unassignedInvoice, false)
	assertGroupFContractCreate(t, ctx, tx, assignedClient, senderID, true)
	assertGroupFContractCreate(t, ctx, tx, unassignedClient, senderID, false)
	assertGroupFInvoiceCreate(t, ctx, tx, assignedClient, senderID, true)
	assertGroupFInvoiceCreate(t, ctx, tx, unassignedClient, senderID, false)
	assertRLSSavepointOperation(t, ctx, tx, "group_f_cross_client_contract_line", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.invoice_line (
			invoice_id, client_id, sender_id, line_no, contract_id, service_type, description, unit
		) VALUES ($1, $2, $3, 99, $4, 'accommodation', 'forged', 'day')`,
			assignedInvoice, assignedClient, senderID, unassignedContract)
		return err
	})

	setAuthorizationActor(t, ctx, tx, allActor.userID.String(), allActor.employeeID.String())
	assertGroupFRowVisible(t, ctx, tx, "contract", "id", unassignedContract, true)
	assertGroupFRowVisible(t, ctx, tx, "invoice", "id", unassignedInvoice, true)
	assertGroupFRowVisible(t, ctx, tx, "invoice_payment_history", "id", unassignedPayment, true)
	assertRLSSavepointOperation(t, ctx, tx, "group_f_framework_owner", false, func() error {
		_, err := tx.Exec(ctx, `UPDATE public.framework_agreement SET client_id = $2 WHERE client_id = $1`, assignedClient, unassignedClient)
		return err
	})
	var invoiceRunID, invoiceRunItemID uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.invoice_run (
		billing_cycle, period_start, period_end, created_by
	) VALUES ('iso_4_week', NOW(), NOW() + INTERVAL '28 days', $1) RETURNING id`, uuid.New()).Scan(&invoiceRunID); err != nil {
		t.Fatalf("create actor-owned invoice run: %v", err)
	}
	if err := tx.QueryRow(ctx, `INSERT INTO public.invoice_run_item (run_id, client_id, sender_id)
		VALUES ($1, $2, $3) RETURNING id`, invoiceRunID, assignedClient, senderID).Scan(&invoiceRunItemID); err != nil {
		t.Fatalf("create actor-owned invoice run item: %v", err)
	}
	assertRLSSavepointOperation(t, ctx, tx, "group_f_run_item_owner", false, func() error {
		_, err := tx.Exec(ctx, `UPDATE public.invoice_run_item SET client_id = $2 WHERE id = $1`, invoiceRunItemID, unassignedClient)
		return err
	})

	setAuthorizationActor(t, ctx, tx, clientOnlyActor.userID.String(), clientOnlyActor.employeeID.String())
	assertGroupFRowVisible(t, ctx, tx, "contract", "id", assignedContract, false)
	assertGroupFRowVisible(t, ctx, tx, "invoice", "id", assignedInvoice, false)
	assertGroupFRowVisible(t, ctx, tx, "invoice_payment_history", "id", assignedPayment, false)

	setAuthorizationActor(t, ctx, tx, assignedActor.userID.String(), assignedActor.employeeID.String())
	assertRLSSavepointOperation(t, ctx, tx, "group_f_assigned_invoice_run", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.invoice_run (billing_cycle, period_start, period_end)
			VALUES ('iso_4_week', NOW(), NOW() + INTERVAL '28 days')`)
		return err
	})

	setAuthorizationActor(t, ctx, tx, invoiceViewActor.userID.String(), invoiceViewActor.employeeID.String())
	assertGroupFRowVisible(t, ctx, tx, "invoice", "id", assignedInvoice, true)
	assertGroupFRowVisible(t, ctx, tx, "invoice_payment_history", "id", assignedPayment, false)
	if paid, err := New(tx).GetInvoicePaidTotal(ctx, assignedInvoice); err != nil {
		t.Fatalf("get protected invoice paid total: %v", err)
	} else if paid != 25 {
		t.Fatalf("invoice paid total = %v, want 25", paid)
	}
	pdfAttachmentID := uuid.New()
	if _, err := tx.Exec(ctx, `INSERT INTO public.attachment_file (uuid, name, file, size, is_used, tag)
		VALUES ($1, 'invoice.pdf', 'group-f-invoice.pdf', 1, TRUE, 'invoice_pdf')`, pdfAttachmentID); err != nil {
		t.Fatalf("create generated invoice PDF attachment: %v", err)
	}
	if _, err := New(tx).InsertIncoicePdfUrl(ctx, InsertIncoicePdfUrlParams{InvoiceID: assignedInvoice, AttachmentID: pdfAttachmentID}); err != nil {
		t.Fatalf("attach generated invoice PDF with view permission: %v", err)
	}

	setAuthorizationActor(t, ctx, tx, paymentViewActor.userID.String(), paymentViewActor.employeeID.String())
	assertGroupFRowVisible(t, ctx, tx, "invoice", "id", assignedInvoice, false)
	assertGroupFRowVisible(t, ctx, tx, "invoice_payment_history", "id", assignedPayment, true)

	setAuthorizationActor(t, ctx, tx, invoiceCreateActor.userID.String(), invoiceCreateActor.employeeID.String())
	assertGroupFRowVisible(t, ctx, tx, "contract", "id", assignedContract, true)
	assertGroupFRowVisible(t, ctx, tx, "contract", "id", unassignedContract, false)
	createOnlyInvoice := createGroupFInvoiceAsRuntime(t, ctx, tx, assignedClient, senderID)
	if _, err := tx.Exec(ctx, `UPDATE public.invoice SET gross_total_amount = 75 WHERE id = $1`, createOnlyInvoice); err != nil {
		t.Fatalf("finalize create-only invoice: %v", err)
	}

	setAuthorizationActor(t, ctx, tx, invoiceUpdateActor.userID.String(), invoiceUpdateActor.employeeID.String())
	if _, err := New(tx).AllocateInvoiceSequenceForDate(ctx, pgtype.Timestamptz{Time: time.Now(), Valid: true}); err != nil {
		t.Fatalf("allocate credit invoice sequence with update permission: %v", err)
	}
	var creditInvoice uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.invoice (
		invoice_number, due_date, invoice_type, original_invoice_id, client_id, sender_id
	) VALUES ($1, CURRENT_DATE + 30, 'credit_note', $2, $3, $4) RETURNING id`,
		"CREDIT-"+uuid.NewString(), assignedInvoice, assignedClient, senderID).Scan(&creditInvoice); err != nil {
		t.Fatalf("create credit invoice with update permission: %v", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE public.invoice SET gross_total_amount = -25 WHERE id = $1`, creditInvoice); err != nil {
		t.Fatalf("finalize credit invoice with update permission: %v", err)
	}

	setAuthorizationActor(t, ctx, tx, paymentCreateActor.userID.String(), paymentCreateActor.employeeID.String())
	if err := New(tx).BeginInvoicePaymentOperation(ctx, BeginInvoicePaymentOperationParams{
		InvoiceID: assignedInvoice, PermissionName: "INVOICE.PAYMENT.CREATE",
	}); err != nil {
		t.Fatalf("begin payment creation: %v", err)
	}
	createdPayment, err := New(tx).CreatePayment(ctx, CreatePaymentParams{
		InvoiceID: assignedInvoice, PaymentMethod: PaymentMethodEnumCash, PaymentStatus: PaymentStatusEnumPending,
		Amount: 10, PaymentDate: pgtype.Date{Time: time.Now(), Valid: true}, RecordedBy: &allActor.employeeID,
	})
	if err != nil {
		t.Fatalf("create actor-owned payment: %v", err)
	}
	if createdPayment.ClientID != assignedClient {
		t.Fatalf("payment client = %s, want %s", createdPayment.ClientID, assignedClient)
	}
	if createdPayment.RecordedBy == nil || *createdPayment.RecordedBy != paymentCreateActor.employeeID {
		t.Fatalf("payment recorder = %v, want actor %s", createdPayment.RecordedBy, paymentCreateActor.employeeID)
	}
	assertRLSSavepointOperation(t, ctx, tx, "group_f_negative_payment", false, func() error {
		_, err := New(tx).CreatePayment(ctx, CreatePaymentParams{
			InvoiceID: assignedInvoice, PaymentMethod: PaymentMethodEnumCash, PaymentStatus: PaymentStatusEnumCompleted,
			Amount: -1, PaymentDate: pgtype.Date{Time: time.Now(), Valid: true},
		})
		return err
	})
	assertGroupFInvoiceUpdateFiltered(t, ctx, tx, assignedInvoice, `gross_total_amount = 999`)
	assertGroupFInvoiceUpdateFiltered(t, ctx, tx, assignedInvoice, `status = 'canceled'`)
	if status, err := New(tx).RecalculateInvoicePaymentStatus(ctx, assignedInvoice); err != nil {
		t.Fatalf("recalculate invoice payment status: %v", err)
	} else if status != InvoiceStatusEnumOutstanding {
		t.Fatalf("recalculated invoice status = %s, want outstanding", status)
	}

	setAuthorizationActor(t, ctx, tx, paymentViewActor.userID.String(), paymentViewActor.employeeID.String())
	if _, err := New(tx).GetPayment(ctx, GetPaymentParams{PaymentID: assignedPayment, InvoiceID: unassignedInvoice}); err != pgx.ErrNoRows {
		t.Fatalf("cross-invoice payment lookup error = %v, want pgx.ErrNoRows", err)
	}
}

func seedGroupFSender(t *testing.T, ctx context.Context, tx pgx.Tx) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.sender (types, name)
		VALUES ('local_authority', $1) RETURNING id`, "Group F "+uuid.NewString()).Scan(&id); err != nil {
		t.Fatalf("seed Group F sender: %v", err)
	}
	return id
}

func seedGroupFContract(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, senderID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.contract (
		status, start_date, end_date, price, price_time_unit, care_name, care_type, client_id, sender_id
	) VALUES ('approved', NOW(), NOW() + INTERVAL '30 days', 100, 'daily', 'Group F care', 'accommodation', $1, $2)
	RETURNING id`, clientID, senderID).Scan(&id); err != nil {
		t.Fatalf("seed Group F contract: %v", err)
	}
	return id
}

func seedGroupFInvoice(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, senderID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.invoice (
		invoice_number, due_date, status, client_id, sender_id, gross_total_amount
	) VALUES ($1, CURRENT_DATE + 30, 'outstanding', $2, $3, 100) RETURNING id`,
		"GF-"+uuid.NewString(), clientID, senderID).Scan(&id); err != nil {
		t.Fatalf("seed Group F invoice: %v", err)
	}
	return id
}

func seedGroupFInvoiceLine(t *testing.T, ctx context.Context, tx pgx.Tx, invoiceID, clientID, senderID, contractID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.invoice_line (
		invoice_id, client_id, sender_id, line_no, contract_id, service_type, description, quantity, unit, unit_price, net_amount, gross_amount
	) VALUES ($1, $2, $3, 1, $4, 'accommodation', 'Group F line', 1, 'day', 100, 100, 100)
	RETURNING id`, invoiceID, clientID, senderID, contractID).Scan(&id); err != nil {
		t.Fatalf("seed Group F invoice line: %v", err)
	}
	return id
}

func seedGroupFPayment(t *testing.T, ctx context.Context, tx pgx.Tx, invoiceID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.invoice_payment_history (
		invoice_id, payment_method, payment_status, amount
	) VALUES ($1, 'bank_transfer', 'completed', 25) RETURNING id`, invoiceID).Scan(&id); err != nil {
		t.Fatalf("seed Group F payment: %v", err)
	}
	return id
}

func seedGroupFContractChild(t *testing.T, ctx context.Context, tx pgx.Tx, table string, contractID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	query := fmt.Sprintf(`INSERT INTO public.%s (contract_id, agreement_details)
		VALUES ($1, 'Group F agreement') RETURNING id`, pgx.Identifier{table}.Sanitize())
	if err := tx.QueryRow(ctx, query, contractID).Scan(&id); err != nil {
		t.Fatalf("seed Group F contract child: %v", err)
	}
	return id
}

func assertGroupFTablesForced(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	for _, table := range []string{
		"contract_type", "contract", "contract_audit", "contract_reminder", "contract_working_hours",
		"client_agreement", "provision", "framework_agreement", "invoice", "invoice_audit", "invoice_line",
		"invoice_payment_history", "invoice_line_calendar_event", "billed_calendar_event", "invoice_run", "invoice_run_item",
	} {
		var enabled, forced bool
		if err := tx.QueryRow(ctx, `SELECT relrowsecurity, relforcerowsecurity FROM pg_catalog.pg_class
			WHERE oid = $1::regclass`, "public."+table).Scan(&enabled, &forced); err != nil {
			t.Fatalf("read %s RLS flags: %v", table, err)
		}
		if !enabled || !forced {
			t.Fatalf("%s RLS enabled=%t forced=%t, want true/true", table, enabled, forced)
		}
	}
}

func assertGroupFRowVisible(t *testing.T, ctx context.Context, tx pgx.Tx, table, idColumn string, id uuid.UUID, want bool) {
	t.Helper()
	query := fmt.Sprintf(`SELECT EXISTS (SELECT 1 FROM public.%s WHERE %s = $1)`,
		pgx.Identifier{table}.Sanitize(), pgx.Identifier{idColumn}.Sanitize())
	var got bool
	if err := tx.QueryRow(ctx, query, id).Scan(&got); err != nil {
		t.Fatalf("read %s visibility: %v", table, err)
	}
	if got != want {
		t.Fatalf("%s %s visibility = %t, want %t", table, id, got, want)
	}
}

func assertGroupFContractCreate(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, senderID uuid.UUID, want bool) {
	t.Helper()
	assertRLSSavepointOperation(t, ctx, tx, "group_f_contract_create", want, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.contract (
			start_date, end_date, price, price_time_unit, care_name, care_type, client_id, sender_id
		) VALUES (NOW(), NOW() + INTERVAL '30 days', 100, 'daily', 'Runtime care', 'accommodation', $1, $2)`, clientID, senderID)
		return err
	})
}

func assertGroupFInvoiceCreate(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, senderID uuid.UUID, want bool) {
	t.Helper()
	assertRLSSavepointOperation(t, ctx, tx, "group_f_invoice_create", want, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.invoice (invoice_number, due_date, client_id, sender_id)
			VALUES ($1, CURRENT_DATE + 30, $2, $3)`, "RUNTIME-"+uuid.NewString(), clientID, senderID)
		return err
	})
}

func createGroupFInvoiceAsRuntime(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, senderID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.invoice (invoice_number, due_date, client_id, sender_id)
		VALUES ($1, CURRENT_DATE + 30, $2, $3) RETURNING id`, "CREATE-ONLY-"+uuid.NewString(), clientID, senderID).Scan(&id); err != nil {
		t.Fatalf("create invoice as runtime actor: %v", err)
	}
	return id
}

func assertGroupFInvoiceUpdateFiltered(t *testing.T, ctx context.Context, tx pgx.Tx, invoiceID uuid.UUID, setClause string) {
	t.Helper()
	query := fmt.Sprintf(`UPDATE public.invoice SET %s WHERE id = $1`, setClause)
	tag, err := tx.Exec(ctx, query, invoiceID)
	if err != nil {
		t.Fatalf("attempt filtered invoice update: %v", err)
	}
	if tag.RowsAffected() != 0 {
		t.Fatalf("payment operation changed invoice through general UPDATE: %s", setClause)
	}
}
