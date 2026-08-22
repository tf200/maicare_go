package db

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var groupEEvaluationDateOffset int32

func TestClientGroupEPermissionScopePolicies(t *testing.T) {
	store := openIntegrationStore(t)
	ctx := context.Background()
	tx, err := store.ConnPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(ctx) })

	assignedActor := seedClientRLSActor(t, ctx, tx, "", false)
	allActor := seedClientRLSActor(t, ctx, tx, "", false)
	generalViewActor := seedClientRLSActor(t, ctx, tx, "assigned", false)
	uploadActor := seedClientRLSActor(t, ctx, tx, "", false)
	deleteOnlyActor := seedClientRLSActor(t, ctx, tx, "", false)
	careViewActor := seedClientRLSActor(t, ctx, tx, "", false)
	careCreateActor := seedClientRLSActor(t, ctx, tx, "", false)
	careUpdateActor := seedClientRLSActor(t, ctx, tx, "", false)
	evaluationViewActor := seedClientRLSActor(t, ctx, tx, "", false)
	evaluationCreateActor := seedClientRLSActor(t, ctx, tx, "", false)
	deniedActor := seedClientRLSActor(t, ctx, tx, "", false)

	permissions := []string{
		"CLIENT.DOCUMENTS.VIEW", "CLIENT.DOCUMENTS.UPLOAD", "CLIENT.DOCUMENTS.DELETE",
		"CLIENT.CARE_PLAN.VIEW", "CLIENT.CARE_PLAN.CREATE", "CLIENT.CARE_PLAN.UPDATE", "CLIENT.CARE_PLAN.DELETE",
		"CLIENT.EVALUATION.VIEW", "CLIENT.EVALUATION.CREATE",
	}
	for _, permission := range permissions {
		grantClientRLSPermission(t, ctx, tx, assignedActor.userID, permission, "assigned")
		grantClientRLSPermission(t, ctx, tx, allActor.userID, permission, "all")
	}
	grantClientRLSPermission(t, ctx, tx, uploadActor.userID, "CLIENT.DOCUMENTS.UPLOAD", "assigned")
	grantClientRLSPermission(t, ctx, tx, deleteOnlyActor.userID, "CLIENT.DOCUMENTS.DELETE", "assigned")
	grantClientRLSPermission(t, ctx, tx, careViewActor.userID, "CLIENT.CARE_PLAN.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, careCreateActor.userID, "CLIENT.CARE_PLAN.CREATE", "assigned")
	grantClientRLSPermission(t, ctx, tx, careUpdateActor.userID, "CLIENT.CARE_PLAN.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, careUpdateActor.userID, "CLIENT.CARE_PLAN.UPDATE", "assigned")
	grantClientRLSPermission(t, ctx, tx, evaluationViewActor.userID, "CLIENT.EVALUATION.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, evaluationCreateActor.userID, "CLIENT.EVALUATION.VIEW", "assigned")
	grantClientRLSPermission(t, ctx, tx, evaluationCreateActor.userID, "CLIENT.EVALUATION.CREATE", "assigned")

	assignedClient := seedClientRLSClient(t, ctx, tx)
	unassignedClient := seedClientRLSClient(t, ctx, tx)
	for _, actor := range []clientRLSActor{
		assignedActor, generalViewActor, uploadActor, deleteOnlyActor, careViewActor,
		careCreateActor, careUpdateActor, evaluationViewActor, evaluationCreateActor,
	} {
		seedClientRLSAssignment(t, ctx, tx, assignedClient, actor.employeeID)
	}

	assignedAttachment := seedGroupEAttachment(t, ctx, tx, assignedActor.userID)
	unassignedAttachment := seedGroupEAttachment(t, ctx, tx, allActor.userID)
	assignedDocument := seedGroupEDocument(t, ctx, tx, assignedClient, assignedAttachment)
	unassignedDocument := seedGroupEDocument(t, ctx, tx, unassignedClient, unassignedAttachment)
	deleteTargetAttachment := seedGroupEAttachment(t, ctx, tx, allActor.userID)
	if _, err := tx.Exec(ctx, `UPDATE public.attachment_file SET is_used = TRUE WHERE uuid = $1`, deleteTargetAttachment); err != nil {
		t.Fatalf("mark delete target attachment used: %v", err)
	}
	deleteTarget := seedGroupEDocument(t, ctx, tx, assignedClient, deleteTargetAttachment)
	assignedGoal := seedGroupEGoal(t, ctx, tx, assignedClient)
	unassignedGoal := seedGroupEGoal(t, ctx, tx, unassignedClient)
	assignedEvaluation := seedGroupEEvaluation(t, ctx, tx, assignedClient, assignedActor.employeeID)
	unassignedEvaluation := seedGroupEEvaluation(t, ctx, tx, unassignedClient, assignedActor.employeeID)
	assignedItem := seedGroupEEvaluationItem(t, ctx, tx, assignedClient, assignedEvaluation, assignedGoal)
	unassignedItem := seedGroupEEvaluationItem(t, ctx, tx, unassignedClient, unassignedEvaluation, unassignedGoal)

	runtimeRole := "phase9_group_e_runtime_" + uuid.NewString()
	roleName := pgx.Identifier{runtimeRole}.Sanitize()
	if _, err := tx.Exec(ctx, fmt.Sprintf(`CREATE ROLE %s NOLOGIN NOSUPERUSER NOBYPASSRLS`, roleName)); err != nil {
		t.Fatalf("create runtime role: %v", err)
	}
	for _, table := range []string{
		"public.attachment_file", "public.client_documents", "public.client_goals",
		"public.client_goal_evaluations", "public.client_goal_evaluation_items",
	} {
		if _, err := tx.Exec(ctx, fmt.Sprintf(`GRANT SELECT, INSERT, UPDATE, DELETE ON %s TO %s`, table, roleName)); err != nil {
			t.Fatalf("grant privileges on %s: %v", table, err)
		}
	}
	for _, signature := range []string{
		"public.get_current_user_id()",
		"public.get_current_employee_id()",
		"public.has_permission(text)",
		"public.get_permission_scope(text)",
		"public.is_assigned_to_client(uuid)",
		"public.can_access_client(uuid,text)",
		"public.begin_client_document_creation(uuid)",
		"public.can_read_created_client_document(uuid)",
		"public.can_read_created_client(uuid)",
		"public.client_has_draft_evaluation_for_goal_update(uuid)",
		"public.goal_has_evaluation_history_for_update(uuid,uuid)",
		"public.can_mutate_goal_evaluation(uuid,uuid)",
		"public.attachment_file_is_referenced(uuid)",
		"public.can_access_actor_attachment(uuid)",
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
		t.Fatal("Group E runtime role unexpectedly bypasses RLS")
	}
	assertGroupETablesForced(t, ctx, tx)

	setAuthorizationActor(t, ctx, tx, assignedActor.userID.String(), assignedActor.employeeID.String())
	for _, row := range []struct {
		table string
		id    uuid.UUID
		want  bool
	}{
		{"client_documents", assignedDocument, true}, {"client_documents", unassignedDocument, false},
		{"client_goals", assignedGoal, true}, {"client_goals", unassignedGoal, false},
		{"client_goal_evaluations", assignedEvaluation, true}, {"client_goal_evaluations", unassignedEvaluation, false},
		{"client_goal_evaluation_items", assignedItem, true}, {"client_goal_evaluation_items", unassignedItem, false},
	} {
		assertGroupERowVisible(t, ctx, tx, row.table, row.id, row.want)
	}
	assertGroupEGoalCreate(t, ctx, tx, assignedClient, true)
	assertGroupEGoalCreate(t, ctx, tx, unassignedClient, false)
	assertGroupEGoalUpdateScope(t, ctx, tx, assignedGoal, unassignedClient, true)
	assertGroupEGoalUpdateScope(t, ctx, tx, unassignedGoal, assignedClient, false)
	assertGroupEEvaluationCreate(t, ctx, tx, assignedClient, assignedActor.employeeID, true)
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluations", assignedEvaluation, true)
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluations", unassignedEvaluation, false)

	setAuthorizationActor(t, ctx, tx, allActor.userID.String(), allActor.employeeID.String())
	assertGroupERowVisible(t, ctx, tx, "client_documents", unassignedDocument, true)
	assertGroupERowVisible(t, ctx, tx, "client_goals", unassignedGoal, true)
	assertGroupERowVisible(t, ctx, tx, "client_goal_evaluations", unassignedEvaluation, true)
	assertGroupERowVisible(t, ctx, tx, "client_goal_evaluation_items", unassignedItem, true)
	assertGroupEGoalCreate(t, ctx, tx, unassignedClient, true)
	assertGroupEGoalUpdateScope(t, ctx, tx, unassignedGoal, assignedClient, true)
	allActorEvaluation := assertGroupEEvaluationCreate(t, ctx, tx, unassignedClient, allActor.employeeID, true)
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluations", allActorEvaluation.ID, true)
	assertRLSSavepointOperation(t, ctx, tx, "group_e_incomplete_evaluation_submit", false, func() error {
		completed := EvaluationStatusEnumCompleted
		_, err := New(tx).UpdateGoalEvaluation(ctx, UpdateGoalEvaluationParams{ID: allActorEvaluation.ID, Status: &completed})
		return err
	})
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluations", unassignedEvaluation, false)
	assertGroupECrossClientItemRejected(t, ctx, tx, unassignedClient, allActorEvaluation.ID, assignedGoal)
	assertGroupEEvaluationDeleteDenied(t, ctx, tx, assignedEvaluation)

	setAuthorizationActor(t, ctx, tx, generalViewActor.userID.String(), generalViewActor.employeeID.String())
	assertGroupEAllRowsHidden(t, ctx, tx, assignedDocument, assignedGoal, assignedEvaluation, assignedItem)
	setAuthorizationActor(t, ctx, tx, deniedActor.userID.String(), deniedActor.employeeID.String())
	assertGroupEAllRowsHidden(t, ctx, tx, assignedDocument, assignedGoal, assignedEvaluation, assignedItem)
	setAuthorizationActor(t, ctx, tx, "", "")
	assertGroupEAllRowsHidden(t, ctx, tx, assignedDocument, assignedGoal, assignedEvaluation, assignedItem)

	setAuthorizationActor(t, ctx, tx, uploadActor.userID.String(), uploadActor.employeeID.String())
	assertGroupERowVisible(t, ctx, tx, "client_documents", assignedDocument, false)
	uploadedAttachment := createGroupEAttachment(t, ctx, tx, uuid.New())
	createdDocument, err := New(tx).CreateClientDocument(ctx, CreateClientDocumentParams{
		ClientID: assignedClient, AttachmentUuid: uploadedAttachment, Label: ClientDocumentLabelEnumOther,
	})
	if err != nil {
		t.Fatalf("upload document through production CreateClientDocument: %v", err)
	}
	assertGroupERowVisible(t, ctx, tx, "client_documents", createdDocument.ID, true)
	if _, err := New(tx).GetActorAttachmentById(ctx, uploadedAttachment); err != pgx.ErrNoRows {
		t.Fatalf("linked attachment without document view error = %v, want pgx.ErrNoRows", err)
	}
	if _, err := New(tx).DeleteActorAttachment(ctx, uploadedAttachment); err != pgx.ErrNoRows {
		t.Fatalf("delete linked attachment error = %v, want pgx.ErrNoRows", err)
	}
	if _, err := New(tx).GetAttachmentById(ctx, uploadedAttachment); err != nil {
		t.Fatalf("trusted internal attachment query was unexpectedly actor-filtered: %v", err)
	}
	assertRLSSavepointOperation(t, ctx, tx, "group_e_foreign_attachment", false, func() error {
		_, err := New(tx).CreateClientDocument(ctx, CreateClientDocumentParams{
			ClientID: assignedClient, AttachmentUuid: assignedAttachment, Label: ClientDocumentLabelEnumOther,
		})
		return err
	})

	setAuthorizationActor(t, ctx, tx, deleteOnlyActor.userID.String(), deleteOnlyActor.employeeID.String())
	if _, err := New(tx).DeleteClientDocument(ctx, DeleteClientDocumentParams{ID: deleteTarget, ClientID: assignedClient}); err == nil {
		t.Fatal("document delete succeeded without document visibility")
	} else if err != pgx.ErrNoRows {
		t.Fatalf("document delete without visibility error = %v, want pgx.ErrNoRows", err)
	}

	setAuthorizationActor(t, ctx, tx, careViewActor.userID.String(), careViewActor.employeeID.String())
	assertGroupERowVisible(t, ctx, tx, "client_goals", assignedGoal, true)
	assertGroupEGoalCreate(t, ctx, tx, assignedClient, false)
	assertGroupEGoalUpdateScope(t, ctx, tx, assignedGoal, unassignedClient, false)
	setAuthorizationActor(t, ctx, tx, careCreateActor.userID.String(), careCreateActor.employeeID.String())
	assertGroupERowVisible(t, ctx, tx, "client_goals", assignedGoal, false)
	assertGroupEGoalCreate(t, ctx, tx, assignedClient, true)
	assertGroupEGoalUpdateScope(t, ctx, tx, assignedGoal, unassignedClient, false)
	setAuthorizationActor(t, ctx, tx, careUpdateActor.userID.String(), careUpdateActor.employeeID.String())
	assertGroupEGoalUpdateScope(t, ctx, tx, assignedGoal, unassignedClient, true)
	assertGroupEGoalUpdateScope(t, ctx, tx, unassignedGoal, assignedClient, false)

	setAuthorizationActor(t, ctx, tx, evaluationViewActor.userID.String(), evaluationViewActor.employeeID.String())
	assertGroupERowVisible(t, ctx, tx, "client_goal_evaluations", assignedEvaluation, true)
	assertGroupERowVisible(t, ctx, tx, "client_goal_evaluation_items", assignedItem, true)
	assertGroupEEvaluationCreate(t, ctx, tx, assignedClient, evaluationViewActor.employeeID, false)
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluations", assignedEvaluation, false)
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluation_items", assignedItem, false)
	setAuthorizationActor(t, ctx, tx, evaluationCreateActor.userID.String(), evaluationCreateActor.employeeID.String())
	createdEvaluation := assertGroupEEvaluationCreate(t, ctx, tx, assignedClient, uuid.New(), true)
	if createdEvaluation.CreatedByEmployeeID == nil || *createdEvaluation.CreatedByEmployeeID != evaluationCreateActor.employeeID {
		t.Fatalf("evaluation creator = %v, want actor %s", createdEvaluation.CreatedByEmployeeID, evaluationCreateActor.employeeID)
	}
	assertGroupEEvaluationCreate(t, ctx, tx, unassignedClient, evaluationCreateActor.employeeID, false)
	activeGoals, err := New(tx).ListActiveGoalsByClientID(ctx, assignedClient)
	if err != nil {
		t.Fatalf("list active goals for actor-owned evaluation: %v", err)
	}
	var createdItem ClientGoalEvaluationItem
	for _, goal := range activeGoals {
		item, err := New(tx).UpsertGoalEvaluationItem(ctx, UpsertGoalEvaluationItemParams{
			ClientID:     assignedClient,
			EvaluationID: createdEvaluation.ID,
			GoalID:       goal.ID,
			Progress:     ClientGoalProgressEnumGoodProgress,
		})
		if err != nil {
			t.Fatalf("upsert item into actor-owned evaluation: %v", err)
		}
		if goal.ID == assignedGoal {
			createdItem = item
		}
	}
	if createdItem.ID == uuid.Nil {
		t.Fatal("assigned goal evaluation item was not created")
	}
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluations", createdEvaluation.ID, true)
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluation_items", createdItem.ID, true)
	assertRLSSavepointOperation(t, ctx, tx, "group_e_archive_draft", false, func() error {
		archived := EvaluationStatusEnumArchived
		_, err := New(tx).UpdateGoalEvaluation(ctx, UpdateGoalEvaluationParams{ID: createdEvaluation.ID, Status: &archived})
		return err
	})
	completed := EvaluationStatusEnumCompleted
	if _, err := New(tx).UpdateGoalEvaluation(ctx, UpdateGoalEvaluationParams{ID: createdEvaluation.ID, Status: &completed}); err != nil {
		t.Fatalf("complete actor-owned evaluation: %v", err)
	}
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluations", createdEvaluation.ID, false)
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluation_items", createdItem.ID, false)
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluations", assignedEvaluation, false)
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluation_items", assignedItem, false)
	assertGroupEUpdate(t, ctx, tx, "client_goal_evaluations", unassignedEvaluation, false)

	setAuthorizationActor(t, ctx, tx, assignedActor.userID.String(), assignedActor.employeeID.String())
	if _, err := New(tx).CancelClientGoalByID(ctx, CancelClientGoalByIDParams{ID: assignedGoal, ClientID: assignedClient}); err != nil {
		t.Fatalf("cancel evaluated goal during versioning: %v", err)
	}
	if _, err := New(tx).CreateReviewUpdatedClientGoal(ctx, CreateReviewUpdatedClientGoalParams{
		ClientID: assignedClient,
		Title:    "Versioned Group E goal",
		Priority: ClientGoalPriorityEnumMedium,
	}); err != nil {
		t.Fatalf("create replacement goal under care-plan update permission: %v", err)
	}
	assertGroupEOwnershipImmutable(t, ctx, tx, unassignedClient, assignedGoal, assignedEvaluation, assignedItem, unassignedGoal, unassignedEvaluation)
	assertGroupERowVisible(t, ctx, tx, "client_documents", deleteTarget, true)
	if _, err := New(tx).DeleteClientDocument(ctx, DeleteClientDocumentParams{ID: deleteTarget, ClientID: assignedClient}); err != nil {
		t.Fatalf("delete visible document with delete permission: %v", err)
	}
	var isUsed bool
	if err := tx.QueryRow(ctx, `SELECT is_used FROM public.attachment_file WHERE uuid = $1`, deleteTargetAttachment).Scan(&isUsed); err != nil {
		t.Fatalf("read released document attachment: %v", err)
	} else if isUsed {
		t.Fatal("deleted document attachment remained marked as used")
	}
}

func seedGroupEAttachment(t *testing.T, ctx context.Context, tx pgx.Tx, uploader uuid.UUID) uuid.UUID {
	t.Helper()
	id := uuid.New()
	if _, err := tx.Exec(ctx, `INSERT INTO public.attachment_file (uuid, uploaded_by_user_id, name, file, size)
		VALUES ($1, $2, 'Group E', 'group-e.pdf', 1)`, id, uploader); err != nil {
		t.Fatalf("seed Group E attachment: %v", err)
	}
	return id
}

func seedGroupEDocument(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, attachmentID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.client_documents (client_id, attachment_uuid, label)
		VALUES ($1, $2, 'other') RETURNING id`, clientID, attachmentID).Scan(&id); err != nil {
		t.Fatalf("seed Group E document: %v", err)
	}
	return id
}

func seedGroupEGoal(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.client_goals (client_id, title, source)
		VALUES ($1, $2, 'manual') RETURNING id`, clientID, "Goal "+uuid.NewString()).Scan(&id); err != nil {
		t.Fatalf("seed Group E goal: %v", err)
	}
	return id
}

func seedGroupEEvaluation(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, employeeID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.client_goal_evaluations (
		client_id, evaluation_date, evaluation_interval_weeks, status, created_by_employee_id
	) VALUES ($1, CURRENT_DATE, 4, 'draft', $2) RETURNING id`, clientID, employeeID).Scan(&id); err != nil {
		t.Fatalf("seed Group E evaluation: %v", err)
	}
	return id
}

func seedGroupEEvaluationItem(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, evaluationID, goalID uuid.UUID) uuid.UUID {
	t.Helper()
	var id uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.client_goal_evaluation_items (
		client_id, evaluation_id, goal_id, progress
	) VALUES ($1, $2, $3, 'limited_progress') RETURNING id`, clientID, evaluationID, goalID).Scan(&id); err != nil {
		t.Fatalf("seed Group E evaluation item: %v", err)
	}
	return id
}

func createGroupEAttachment(t *testing.T, ctx context.Context, tx pgx.Tx, forgedUploader uuid.UUID) uuid.UUID {
	t.Helper()
	var id, uploader uuid.UUID
	if err := tx.QueryRow(ctx, `INSERT INTO public.attachment_file (
		uuid, uploaded_by_user_id, name, file, size
	) VALUES ($1, $2, 'Runtime', 'runtime.pdf', 1) RETURNING uuid, uploaded_by_user_id`,
		uuid.New(), forgedUploader).Scan(&id, &uploader); err != nil {
		t.Fatalf("create runtime attachment: %v", err)
	}
	var actor uuid.UUID
	if err := tx.QueryRow(ctx, `SELECT current_setting('myapp.current_user_id')::uuid`).Scan(&actor); err != nil {
		t.Fatalf("read attachment actor: %v", err)
	}
	if uploader != actor {
		t.Fatalf("attachment uploader = %s, want actor %s", uploader, actor)
	}
	return id
}

func assertGroupETablesForced(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	for _, table := range []string{"client_documents", "client_goals", "client_goal_evaluations", "client_goal_evaluation_items"} {
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

func assertGroupERowVisible(t *testing.T, ctx context.Context, tx pgx.Tx, table string, id uuid.UUID, want bool) {
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

func assertGroupEAllRowsHidden(t *testing.T, ctx context.Context, tx pgx.Tx, documentID, goalID, evaluationID, itemID uuid.UUID) {
	t.Helper()
	assertGroupERowVisible(t, ctx, tx, "client_documents", documentID, false)
	assertGroupERowVisible(t, ctx, tx, "client_goals", goalID, false)
	assertGroupERowVisible(t, ctx, tx, "client_goal_evaluations", evaluationID, false)
	assertGroupERowVisible(t, ctx, tx, "client_goal_evaluation_items", itemID, false)
}

func assertGroupEGoalCreate(t *testing.T, ctx context.Context, tx pgx.Tx, clientID uuid.UUID, want bool) {
	t.Helper()
	assertRLSSavepointOperation(t, ctx, tx, "group_e_goal_create", want, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.client_goals (client_id, title, source)
			VALUES ($1, $2, 'manual')`, clientID, "Runtime "+uuid.NewString())
		return err
	})
}

func assertGroupEEvaluationCreate(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, forgedEmployeeID uuid.UUID, want bool) ClientGoalEvaluation {
	t.Helper()
	var evaluation ClientGoalEvaluation
	evaluationDate := time.Now().AddDate(0, 0, int(atomic.AddInt32(&groupEEvaluationDateOffset, 1)))
	assertRLSSavepointOperation(t, ctx, tx, "group_e_evaluation_create", want, func() error {
		var err error
		evaluation, err = New(tx).CreateGoalEvaluation(ctx, CreateGoalEvaluationParams{
			ClientID:                clientID,
			EvaluationDate:          pgtype.Date{Time: evaluationDate, Valid: true},
			EvaluationIntervalWeeks: 4,
			Status:                  EvaluationStatusEnumCompleted,
			CreatedByEmployeeID:     &forgedEmployeeID,
		})
		return err
	})
	if want && evaluation.Status != EvaluationStatusEnumDraft {
		t.Fatalf("runtime evaluation status = %s, want draft", evaluation.Status)
	}
	return evaluation
}

func assertGroupEUpdate(t *testing.T, ctx context.Context, tx pgx.Tx, table string, id uuid.UUID, want bool) {
	t.Helper()
	query := fmt.Sprintf(`UPDATE public.%s SET updated_at = clock_timestamp() WHERE id = $1`, pgx.Identifier{table}.Sanitize())
	tag, err := tx.Exec(ctx, query, id)
	if err != nil {
		t.Fatalf("update %s: %v", table, err)
	}
	if got := tag.RowsAffected() == 1; got != want {
		t.Fatalf("%s updated = %t, want %t", table, got, want)
	}
}

func assertGroupEGoalUpdateScope(t *testing.T, ctx context.Context, tx pgx.Tx, goalID, newClientID uuid.UUID, authorized bool) {
	t.Helper()
	if _, err := tx.Exec(ctx, "SAVEPOINT group_e_goal_update_scope"); err != nil {
		t.Fatalf("create goal update savepoint: %v", err)
	}
	tag, err := tx.Exec(ctx, `UPDATE public.client_goals SET client_id = $2 WHERE id = $1`, goalID, newClientID)
	if authorized {
		if err == nil {
			t.Fatal("authorized goal ownership change bypassed immutability trigger")
		}
		if _, rollbackErr := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT group_e_goal_update_scope"); rollbackErr != nil {
			t.Fatalf("rollback authorized goal update: %v", rollbackErr)
		}
		return
	}
	if err != nil {
		t.Fatalf("unauthorized goal update returned an error instead of filtering the row: %v", err)
	}
	if tag.RowsAffected() != 0 {
		t.Fatal("unauthorized goal update reached the ownership trigger")
	}
}

func assertGroupECrossClientItemRejected(t *testing.T, ctx context.Context, tx pgx.Tx, clientID, evaluationID, goalID uuid.UUID) {
	t.Helper()
	assertRLSSavepointOperation(t, ctx, tx, "group_e_cross_client_item", false, func() error {
		_, err := tx.Exec(ctx, `INSERT INTO public.client_goal_evaluation_items (
			client_id, evaluation_id, goal_id, progress
		) VALUES ($1, $2, $3, 'good_progress')`, clientID, evaluationID, goalID)
		return err
	})
}

func assertGroupEOwnershipImmutable(t *testing.T, ctx context.Context, tx pgx.Tx, unassignedClient, goalID, evaluationID, itemID, otherGoalID, otherEvaluationID uuid.UUID) {
	t.Helper()
	operations := []struct {
		name  string
		query string
		args  []any
	}{
		{"group_e_goal_owner", `UPDATE public.client_goals SET client_id = $2 WHERE id = $1`, []any{goalID, unassignedClient}},
		{"group_e_evaluation_owner", `UPDATE public.client_goal_evaluations SET client_id = $2 WHERE id = $1`, []any{evaluationID, unassignedClient}},
		{"group_e_evaluation_creator", `UPDATE public.client_goal_evaluations SET created_by_employee_id = $2 WHERE id = $1`, []any{evaluationID, uuid.New()}},
		{"group_e_item_client", `UPDATE public.client_goal_evaluation_items SET client_id = $2 WHERE id = $1`, []any{itemID, unassignedClient}},
		{"group_e_item_evaluation", `UPDATE public.client_goal_evaluation_items SET evaluation_id = $2 WHERE id = $1`, []any{itemID, otherEvaluationID}},
		{"group_e_item_goal", `UPDATE public.client_goal_evaluation_items SET goal_id = $2 WHERE id = $1`, []any{itemID, otherGoalID}},
	}
	for _, operation := range operations {
		operation := operation
		assertRLSSavepointOperation(t, ctx, tx, operation.name, false, func() error {
			_, err := tx.Exec(ctx, operation.query, operation.args...)
			return err
		})
	}
}

func assertGroupEEvaluationDeleteDenied(t *testing.T, ctx context.Context, tx pgx.Tx, evaluationID uuid.UUID) {
	t.Helper()
	tag, err := tx.Exec(ctx, `DELETE FROM public.client_goal_evaluations WHERE id = $1`, evaluationID)
	if err != nil {
		t.Fatalf("delete evaluation without policy: %v", err)
	}
	if tag.RowsAffected() != 0 {
		t.Fatal("evaluation delete succeeded without a delete policy")
	}
}
