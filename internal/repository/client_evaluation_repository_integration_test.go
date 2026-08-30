package repository

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/ctxkeys"
	"maicare_go/internal/domain"
	"maicare_go/internal/testdatabase"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var evaluationIntegrationPool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()
	databaseURL, cleanup, err := testdatabase.StartPostgres(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse integration database URL: %v\n", err)
		_ = cleanup()
		os.Exit(1)
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		return db.RegisterEnumTypes(ctx, conn)
	}
	evaluationIntegrationPool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect integration database: %v\n", err)
		_ = cleanup()
		os.Exit(1)
	}

	exitCode := m.Run()
	evaluationIntegrationPool.Close()
	if err := cleanup(); err != nil {
		fmt.Fprintf(os.Stderr, "terminate PostgreSQL test container: %v\n", err)
		exitCode = 1
	}
	os.Exit(exitCode)
}

func TestGoalEvaluationBootstrapSelectsOnlyCurrentCycleDraft(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}

	result, err := repository.GetGoalEvaluationBootstrap(fixture.ownerContext(), fixture.clientID)
	if err != nil {
		t.Fatalf("GetGoalEvaluationBootstrap() error = %v", err)
	}
	if result.ExistingDraft == nil {
		t.Fatal("GetGoalEvaluationBootstrap() returned no current-cycle draft")
	}
	if result.ExistingDraft.ID != fixture.currentEvaluationID {
		t.Fatalf("existing draft ID = %s, want current-cycle ID %s", result.ExistingDraft.ID, fixture.currentEvaluationID)
	}
	if !sameDate(result.ExistingDraft.EvaluationDate, fixture.currentDate) {
		t.Fatalf("existing draft date = %s, want %s", result.ExistingDraft.EvaluationDate, fixture.currentDate)
	}
}

func TestGetEvaluationStatsExecutesAsActorAndReturnsAsOf(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	if _, err := evaluationIntegrationPool.Exec(context.Background(), `UPDATE client_goal_evaluations SET status = 'completed', updated_at = CURRENT_TIMESTAMP WHERE id = $1`, fixture.historicalEvaluationID); err != nil {
		t.Fatalf("complete historical evaluation: %v", err)
	}
	if _, err := evaluationIntegrationPool.Exec(context.Background(), `UPDATE client_details SET last_evaluation_anchor_date = $2::date - 28, next_evaluation_date = $2::date WHERE id = $1`, fixture.clientID, fixture.currentDate); err != nil {
		t.Fatalf("restore current evaluation cycle: %v", err)
	}
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}

	result, err := repository.GetEvaluationStats(fixture.ownerContext(), fixture.owner.employeeID)
	if err != nil {
		t.Fatalf("GetEvaluationStats() error = %v", err)
	}
	if result.AsOf.IsZero() {
		t.Fatal("stats as_of is zero")
	}
	if result.AttentionRequired != 1 || result.InProgress != 1 || result.RecentlyFinalized != 1 {
		t.Fatalf("stats = %#v, want attention_required=1 in_progress=1 recently_finalized=1", result)
	}
}

func TestEvaluationBusinessDateUsesConfiguredTimezoneAcrossDST(t *testing.T) {
	ctx := context.Background()
	if _, err := evaluationIntegrationPool.Exec(ctx, `UPDATE app_organization_profile SET default_timezone = 'Europe/Amsterdam' WHERE singleton = TRUE`); err != nil {
		t.Fatalf("set evaluation timezone: %v", err)
	}

	tests := []struct {
		name string
		at   string
		want string
	}{
		{name: "before spring transition", at: "2026-03-28T23:30:00Z", want: "2026-03-29"},
		{name: "after spring transition local midnight", at: "2026-03-29T22:30:00Z", want: "2026-03-30"},
		{name: "before autumn transition", at: "2026-10-24T22:30:00Z", want: "2026-10-25"},
		{name: "after autumn transition local midnight", at: "2026-10-25T23:30:00Z", want: "2026-10-26"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got time.Time
			if err := evaluationIntegrationPool.QueryRow(ctx, `SELECT evaluation_business_date($1::timestamptz)`, tt.at).Scan(&got); err != nil {
				t.Fatalf("evaluation_business_date() error = %v", err)
			}
			if got.Format(time.DateOnly) != tt.want {
				t.Fatalf("evaluation_business_date(%s) = %s, want %s", tt.at, got.Format(time.DateOnly), tt.want)
			}
		})
	}
}

func TestGoalEvaluationSubmissionAllowsExactFourteenDayBoundary(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	var businessDate time.Time
	if err := evaluationIntegrationPool.QueryRow(context.Background(), `SELECT evaluation_business_date()`).Scan(&businessDate); err != nil {
		t.Fatalf("get evaluation business date: %v", err)
	}
	dueDate := businessDate.AddDate(0, 0, 14)
	if _, err := evaluationIntegrationPool.Exec(context.Background(), `UPDATE client_details SET next_evaluation_date = $2 WHERE id = $1`, fixture.clientID, dueDate); err != nil {
		t.Fatalf("move evaluation schedule: %v", err)
	}
	if _, err := evaluationIntegrationPool.Exec(context.Background(), `UPDATE client_goal_evaluations SET evaluation_date = $2 WHERE id = $1`, fixture.currentEvaluationID, dueDate); err != nil {
		t.Fatalf("move evaluation date: %v", err)
	}

	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}
	result, err := repository.SubmitGoalEvaluationDraft(fixture.ownerContext(), fixture.currentEvaluationID, fixture.owner.employeeID, submitParams(t, fixture.currentEvaluationID))
	if err != nil {
		t.Fatalf("SubmitGoalEvaluationDraft() error = %v", err)
	}
	if result.Status != "completed" {
		t.Fatalf("status = %s, want completed", result.Status)
	}
}

func TestGoalEvaluationBootstrapReturnsNoDraftWithoutScheduledCycle(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	if _, err := evaluationIntegrationPool.Exec(context.Background(), `UPDATE client_details SET next_evaluation_date = NULL WHERE id = $1`, fixture.clientID); err != nil {
		t.Fatalf("clear client evaluation date: %v", err)
	}
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}

	result, err := repository.GetGoalEvaluationBootstrap(fixture.ownerContext(), fixture.clientID)
	if err != nil {
		t.Fatalf("GetGoalEvaluationBootstrap() error = %v", err)
	}
	if result.NextEvaluationDate != nil || result.ExistingDraft != nil {
		t.Fatalf("bootstrap without schedule returned next_date=%v existing_draft=%v", result.NextEvaluationDate, result.ExistingDraft)
	}
}

func TestUpdateGoalEvaluationDraftMutatesExactID(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}
	notes := "updated current evaluation"
	itemNotes := "updated current goal"

	result, err := repository.UpdateGoalEvaluationDraft(
		fixture.ownerContext(),
		fixture.currentEvaluationID,
		fixture.owner.employeeID,
		domain.UpdateGoalEvaluationDraftParams{
			ExpectedUpdatedAt: submitParams(t, fixture.currentEvaluationID).ExpectedUpdatedAt,
			OverallNotes:      &notes,
			Items: []domain.GoalEvaluationItemParams{{
				GoalID: fixture.goalID, Progress: "achieved", Notes: &itemNotes,
			}},
		},
	)
	if err != nil {
		t.Fatalf("UpdateGoalEvaluationDraft() error = %v", err)
	}
	if result.ID != fixture.currentEvaluationID {
		t.Fatalf("updated evaluation ID = %s, want %s", result.ID, fixture.currentEvaluationID)
	}

	var currentNotes, historicalNotes *string
	if err := evaluationIntegrationPool.QueryRow(context.Background(), `SELECT overall_notes FROM client_goal_evaluations WHERE id = $1`, fixture.currentEvaluationID).Scan(&currentNotes); err != nil {
		t.Fatalf("read current evaluation: %v", err)
	}
	if err := evaluationIntegrationPool.QueryRow(context.Background(), `SELECT overall_notes FROM client_goal_evaluations WHERE id = $1`, fixture.historicalEvaluationID).Scan(&historicalNotes); err != nil {
		t.Fatalf("read historical evaluation: %v", err)
	}
	if currentNotes == nil || *currentNotes != notes {
		t.Fatalf("current evaluation notes = %v, want %q", currentNotes, notes)
	}
	if historicalNotes == nil || *historicalNotes != "historical notes" {
		t.Fatalf("historical evaluation notes = %v, want unchanged", historicalNotes)
	}
	var currentProgress, historicalProgress string
	if err := evaluationIntegrationPool.QueryRow(context.Background(), `SELECT progress::text FROM client_goal_evaluation_items WHERE evaluation_id = $1 AND goal_id = $2`, fixture.currentEvaluationID, fixture.goalID).Scan(&currentProgress); err != nil {
		t.Fatalf("read current evaluation item: %v", err)
	}
	if err := evaluationIntegrationPool.QueryRow(context.Background(), `SELECT progress::text FROM client_goal_evaluation_items WHERE evaluation_id = $1 AND goal_id = $2`, fixture.historicalEvaluationID, fixture.goalID).Scan(&historicalProgress); err != nil {
		t.Fatalf("read historical evaluation item: %v", err)
	}
	if currentProgress != "achieved" || historicalProgress != "good_progress" {
		t.Fatalf("evaluation progress current=%s historical=%s, want achieved/good_progress", currentProgress, historicalProgress)
	}
}

func TestHistoricalGoalEvaluationDraftIsReadOnly(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}
	notes := "must not be written"
	params := domain.UpdateGoalEvaluationDraftParams{
		OverallNotes: &notes,
		Items:        []domain.GoalEvaluationItemParams{{GoalID: fixture.goalID, Progress: "achieved"}},
	}

	if _, err := repository.UpdateGoalEvaluationDraft(fixture.ownerContext(), fixture.historicalEvaluationID, fixture.owner.employeeID, params); !errors.Is(err, domain.ErrGoalEvaluationNotCurrentCycle) {
		t.Fatalf("historical update error = %v, want ErrGoalEvaluationNotCurrentCycle", err)
	}
	if _, err := repository.SubmitGoalEvaluationDraft(fixture.ownerContext(), fixture.historicalEvaluationID, fixture.owner.employeeID, submitParams(t, fixture.historicalEvaluationID)); !errors.Is(err, domain.ErrGoalEvaluationNotCurrentCycle) {
		t.Fatalf("historical submit error = %v, want ErrGoalEvaluationNotCurrentCycle", err)
	}

	var notesAfter *string
	var status string
	if err := evaluationIntegrationPool.QueryRow(context.Background(), `SELECT overall_notes, status::text FROM client_goal_evaluations WHERE id = $1`, fixture.historicalEvaluationID).Scan(&notesAfter, &status); err != nil {
		t.Fatalf("read historical evaluation: %v", err)
	}
	if notesAfter == nil || *notesAfter != "historical notes" || status != "draft" {
		t.Fatalf("historical evaluation changed: notes=%v status=%s", notesAfter, status)
	}
}

func TestGoalEvaluationDraftRejectsNonOwner(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}
	notes := "must not be written"
	params := domain.UpdateGoalEvaluationDraftParams{OverallNotes: &notes}

	if _, err := repository.UpdateGoalEvaluationDraft(fixture.otherContext(), fixture.currentEvaluationID, fixture.other.employeeID, params); !errors.Is(err, domain.ErrGoalEvaluationOwnedByOther) {
		t.Fatalf("non-owner update error = %v, want ErrGoalEvaluationOwnedByOther", err)
	}
	if _, err := repository.SubmitGoalEvaluationDraft(fixture.otherContext(), fixture.currentEvaluationID, fixture.other.employeeID, submitParams(t, fixture.currentEvaluationID)); !errors.Is(err, domain.ErrGoalEvaluationOwnedByOther) {
		t.Fatalf("non-owner submit error = %v, want ErrGoalEvaluationOwnedByOther", err)
	}
}

func TestSubmitGoalEvaluationDraftAdvancesScheduleOnce(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}

	result, err := repository.SubmitGoalEvaluationDraft(fixture.ownerContext(), fixture.currentEvaluationID, fixture.owner.employeeID, submitParams(t, fixture.currentEvaluationID))
	if err != nil {
		t.Fatalf("SubmitGoalEvaluationDraft() error = %v", err)
	}
	if result.Status != "completed" {
		t.Fatalf("submit result status=%q, want completed", result.Status)
	}

	wantNextDate := fixture.currentDate.AddDate(0, 0, 28)
	assertClientSchedule(t, fixture.clientID, fixture.currentDate, wantNextDate)

	if _, err := repository.SubmitGoalEvaluationDraft(fixture.ownerContext(), fixture.currentEvaluationID, fixture.owner.employeeID, domain.SubmitGoalEvaluationDraftParams{ExpectedUpdatedAt: result.UpdatedAt}); !errors.Is(err, domain.ErrGoalEvaluationNotDraft) {
		t.Fatalf("second submit error = %v, want ErrGoalEvaluationNotDraft", err)
	}
	assertClientSchedule(t, fixture.clientID, fixture.currentDate, wantNextDate)
}

func TestCreateGoalEvaluationRejectsCompletedCurrentCycle(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	if _, err := evaluationIntegrationPool.Exec(context.Background(), `UPDATE client_goal_evaluations SET status = 'completed' WHERE id = $1`, fixture.currentEvaluationID); err != nil {
		t.Fatalf("complete current evaluation: %v", err)
	}
	if _, err := evaluationIntegrationPool.Exec(context.Background(), `UPDATE client_details SET next_evaluation_date = $2 WHERE id = $1`, fixture.clientID, fixture.currentDate); err != nil {
		t.Fatalf("restore current evaluation date: %v", err)
	}
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}

	_, err := repository.CreateGoalEvaluation(
		fixture.ownerContext(),
		fixture.clientID,
		fixture.owner.employeeID,
		domain.CreateGoalEvaluationParams{Items: []domain.GoalEvaluationItemParams{{GoalID: fixture.goalID, Progress: "achieved"}}},
	)
	if !errors.Is(err, domain.ErrGoalEvaluationNotDraft) {
		t.Fatalf("CreateGoalEvaluation() error = %v, want ErrGoalEvaluationNotDraft", err)
	}
}

func TestBlockedGoalEvaluationSubmissionReturnsSavedDraft(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	if _, err := evaluationIntegrationPool.Exec(context.Background(), `UPDATE client_goal_evaluation_items SET progress = 'no_progress', notes = 'saved before submit' WHERE evaluation_id = $1`, fixture.currentEvaluationID); err != nil {
		t.Fatalf("make evaluation incomplete: %v", err)
	}
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}

	result, err := repository.SubmitGoalEvaluationDraft(fixture.ownerContext(), fixture.currentEvaluationID, fixture.owner.employeeID, submitParams(t, fixture.currentEvaluationID))
	if !errors.Is(err, domain.ErrGoalEvaluationIncomplete) {
		t.Fatalf("SubmitGoalEvaluationDraft() error = %v, want ErrGoalEvaluationIncomplete", err)
	}
	if result == nil || result.Status != "draft" || len(result.Items) != 1 || result.Items[0].Notes == nil || *result.Items[0].Notes != "saved before submit" {
		t.Fatalf("blocked submission result = %#v, want persisted draft", result)
	}
	assertClientSchedule(t, fixture.clientID, fixture.currentDate.AddDate(0, 0, -28), fixture.currentDate)
}

func TestTooEarlyGoalEvaluationSubmissionReturnsSavedDraft(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	var businessDate time.Time
	if err := evaluationIntegrationPool.QueryRow(context.Background(), `SELECT evaluation_business_date()`).Scan(&businessDate); err != nil {
		t.Fatalf("get evaluation business date: %v", err)
	}
	futureDate := businessDate.AddDate(0, 0, 15)
	if _, err := evaluationIntegrationPool.Exec(context.Background(), `UPDATE client_details SET next_evaluation_date = $2 WHERE id = $1`, fixture.clientID, futureDate); err != nil {
		t.Fatalf("move evaluation schedule: %v", err)
	}
	if _, err := evaluationIntegrationPool.Exec(context.Background(), `UPDATE client_goal_evaluations SET evaluation_date = $2 WHERE id = $1`, fixture.currentEvaluationID, futureDate); err != nil {
		t.Fatalf("move evaluation date: %v", err)
	}
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}

	result, err := repository.SubmitGoalEvaluationDraft(fixture.ownerContext(), fixture.currentEvaluationID, fixture.owner.employeeID, submitParams(t, fixture.currentEvaluationID))
	if !errors.Is(err, domain.ErrGoalEvaluationTooEarly) {
		t.Fatalf("SubmitGoalEvaluationDraft() error = %v, want ErrGoalEvaluationTooEarly", err)
	}
	if result == nil || result.Status != "draft" {
		t.Fatalf("blocked submission result = %#v, want saved draft", result)
	}
}

func TestStaleGoalEvaluationUpdateReturnsCurrentServerDraft(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}
	staleRevision := submitParams(t, fixture.currentEvaluationID).ExpectedUpdatedAt
	winnerNotes := "winner notes"
	winner, err := repository.UpdateGoalEvaluationDraft(fixture.ownerContext(), fixture.currentEvaluationID, fixture.owner.employeeID, domain.UpdateGoalEvaluationDraftParams{
		ExpectedUpdatedAt: staleRevision,
		OverallNotes:      &winnerNotes,
		Items:             []domain.GoalEvaluationItemParams{{GoalID: fixture.goalID, Progress: "achieved"}},
	})
	if err != nil {
		t.Fatalf("winning update error = %v", err)
	}
	loserNotes := "stale overwrite"
	current, err := repository.UpdateGoalEvaluationDraft(fixture.ownerContext(), fixture.currentEvaluationID, fixture.owner.employeeID, domain.UpdateGoalEvaluationDraftParams{
		ExpectedUpdatedAt: staleRevision,
		OverallNotes:      &loserNotes,
		Items:             []domain.GoalEvaluationItemParams{{GoalID: fixture.goalID, Progress: "regression"}},
	})
	if !errors.Is(err, domain.ErrGoalEvaluationConflict) {
		t.Fatalf("stale update error = %v, want ErrGoalEvaluationConflict", err)
	}
	if current == nil || current.OverallNotes == nil || *current.OverallNotes != winnerNotes || current.Items[0].Progress != "achieved" || !current.UpdatedAt.Equal(winner.UpdatedAt) {
		t.Fatalf("conflict result = %#v, want current winning draft", current)
	}
}

func TestStaleGoalEvaluationSubmitDoesNotAdvanceSchedule(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}
	staleRevision := submitParams(t, fixture.currentEvaluationID).ExpectedUpdatedAt
	notes := "newer saved revision"
	current, err := repository.UpdateGoalEvaluationDraft(fixture.ownerContext(), fixture.currentEvaluationID, fixture.owner.employeeID, domain.UpdateGoalEvaluationDraftParams{
		ExpectedUpdatedAt: staleRevision,
		OverallNotes:      &notes,
		Items:             []domain.GoalEvaluationItemParams{{GoalID: fixture.goalID, Progress: "achieved"}},
	})
	if err != nil {
		t.Fatalf("update before stale submit error = %v", err)
	}

	conflict, err := repository.SubmitGoalEvaluationDraft(fixture.ownerContext(), fixture.currentEvaluationID, fixture.owner.employeeID, domain.SubmitGoalEvaluationDraftParams{ExpectedUpdatedAt: staleRevision})
	if !errors.Is(err, domain.ErrGoalEvaluationConflict) {
		t.Fatalf("stale submit error = %v, want ErrGoalEvaluationConflict", err)
	}
	if conflict == nil || conflict.Status != "draft" || !conflict.UpdatedAt.Equal(current.UpdatedAt) {
		t.Fatalf("stale submit result = %#v, want current draft", conflict)
	}
	assertClientSchedule(t, fixture.clientID, fixture.currentDate.AddDate(0, 0, -28), fixture.currentDate)
}

func TestConcurrentGoalEvaluationSubmissionAdvancesScheduleOnce(t *testing.T) {
	fixture := seedEvaluationLifecycleFixture(t)
	repository := &ClientRepository{store: db.NewStore(evaluationIntegrationPool)}
	params := submitParams(t, fixture.currentEvaluationID)
	type submitResult struct {
		evaluation *domain.GoalEvaluation
		err        error
	}
	results := make(chan submitResult, 2)
	var ready sync.WaitGroup
	ready.Add(2)
	start := make(chan struct{})
	for range 2 {
		go func() {
			ready.Done()
			<-start
			evaluation, err := repository.SubmitGoalEvaluationDraft(fixture.ownerContext(), fixture.currentEvaluationID, fixture.owner.employeeID, params)
			results <- submitResult{evaluation: evaluation, err: err}
		}()
	}
	ready.Wait()
	close(start)

	successes := 0
	rejections := 0
	for range 2 {
		result := <-results
		switch {
		case result.err == nil && result.evaluation != nil && result.evaluation.Status == "completed":
			successes++
		case errors.Is(result.err, domain.ErrGoalEvaluationConflict), errors.Is(result.err, domain.ErrGoalEvaluationNotDraft):
			rejections++
		default:
			t.Fatalf("unexpected concurrent submit result: evaluation=%#v error=%v", result.evaluation, result.err)
		}
	}
	if successes != 1 || rejections != 1 {
		t.Fatalf("concurrent submits: successes=%d rejections=%d, want 1/1", successes, rejections)
	}
	assertClientSchedule(t, fixture.clientID, fixture.currentDate, fixture.currentDate.AddDate(0, 0, 28))
}

type evaluationActor struct {
	userID     uuid.UUID
	employeeID uuid.UUID
}

type evaluationLifecycleFixture struct {
	clientID               uuid.UUID
	goalID                 uuid.UUID
	currentEvaluationID    uuid.UUID
	historicalEvaluationID uuid.UUID
	currentDate            time.Time
	owner                  evaluationActor
	other                  evaluationActor
}

func (f evaluationLifecycleFixture) ownerContext() context.Context {
	return actorContext(f.owner)
}

func (f evaluationLifecycleFixture) otherContext() context.Context {
	return actorContext(f.other)
}

func actorContext(actor evaluationActor) context.Context {
	return ctxkeys.WithActorIdentity(context.Background(), ctxkeys.ActorIdentity{
		UserID: actor.userID, EmployeeID: actor.employeeID,
	})
}

func seedEvaluationLifecycleFixture(t *testing.T) evaluationLifecycleFixture {
	t.Helper()
	ctx := context.Background()
	tx, err := evaluationIntegrationPool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin fixture transaction: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	currentDate := time.Now().UTC().Truncate(24 * time.Hour)
	fixture := evaluationLifecycleFixture{
		clientID:               uuid.New(),
		goalID:                 uuid.New(),
		currentEvaluationID:    uuid.New(),
		historicalEvaluationID: uuid.New(),
		currentDate:            currentDate,
		owner:                  seedEvaluationActor(t, ctx, tx),
		other:                  seedEvaluationActor(t, ctx, tx),
	}

	if _, err := tx.Exec(ctx, `INSERT INTO client_details (
		id, first_name, last_name, email, gender, filenumber, street, house_number, postal_code, city
	) VALUES ($1, 'Evaluation', 'Client', $2, 'unknown', $3, 'Test', '1', '1000AA', 'Test')`,
		fixture.clientID, uuid.NewString()+"@test.invalid", uuid.NewString()); err != nil {
		t.Fatalf("seed evaluation client: %v", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO client_goals (id, client_id, title, source, status)
		VALUES ($1, $2, 'Evaluation goal', 'manual', 'active')`, fixture.goalID, fixture.clientID); err != nil {
		t.Fatalf("seed evaluation goal: %v", err)
	}
	if _, err := tx.Exec(ctx, `UPDATE client_details SET
		status = 'in_care', care_start_date = $2::date - 28, placed_in_care_at = NOW() - INTERVAL '28 days',
		evaluation_intervals_weeks = 4, last_evaluation_anchor_date = $2::date - 28, next_evaluation_date = $2::date
		WHERE id = $1`, fixture.clientID, currentDate); err != nil {
		t.Fatalf("place evaluation client in care: %v", err)
	}
	for _, actor := range []evaluationActor{fixture.owner, fixture.other} {
		role := "support"
		if actor.employeeID == fixture.owner.employeeID {
			role = "coordinator"
		}
		if _, err := tx.Exec(ctx, `INSERT INTO assigned_employee (client_id, employee_id, start_date, role)
			VALUES ($1, $2, CURRENT_DATE, $3)`, fixture.clientID, actor.employeeID, role); err != nil {
			t.Fatalf("seed evaluation assignment: %v", err)
		}
	}
	grantEvaluationViewPermission(t, ctx, tx, fixture.owner)
	if _, err := tx.Exec(ctx, `INSERT INTO client_goal_evaluations (
		id, client_id, evaluation_date, evaluation_interval_weeks, status, overall_notes, created_by_employee_id, updated_at
	) VALUES
		($1, $3, $5::date, 4, 'draft', 'current notes', $4, NOW() - INTERVAL '1 day'),
		($2, $3, $5::date - 28, 4, 'draft', 'historical notes', $4, NOW())`,
		fixture.currentEvaluationID, fixture.historicalEvaluationID, fixture.clientID, fixture.owner.employeeID, currentDate); err != nil {
		t.Fatalf("seed evaluations: %v", err)
	}
	for _, evaluationID := range []uuid.UUID{fixture.currentEvaluationID, fixture.historicalEvaluationID} {
		if _, err := tx.Exec(ctx, `INSERT INTO client_goal_evaluation_items (client_id, evaluation_id, goal_id, progress, notes)
			VALUES ($1, $2, $3, 'good_progress', 'seed notes')`, fixture.clientID, evaluationID, fixture.goalID); err != nil {
			t.Fatalf("seed evaluation item: %v", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("commit evaluation fixture: %v", err)
	}
	return fixture
}

func grantEvaluationViewPermission(t *testing.T, ctx context.Context, tx pgx.Tx, actor evaluationActor) {
	t.Helper()
	roleID := uuid.New()
	if _, err := tx.Exec(ctx, `INSERT INTO roles (id, name) VALUES ($1, $2)`, roleID, uuid.NewString()); err != nil {
		t.Fatalf("seed evaluation role: %v", err)
	}
	for _, permission := range []string{"CLIENT.VIEW", "CLIENT.EVALUATION.VIEW"} {
		var permissionID uuid.UUID
		if err := tx.QueryRow(ctx, `INSERT INTO permissions (id, name, is_scoped) VALUES ($1, $2, true)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id`, uuid.New(), permission).Scan(&permissionID); err != nil {
			t.Fatalf("seed evaluation permission %s: %v", permission, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO role_permissions (role_id, permission_id, scope) VALUES ($1, $2, 'all')`, roleID, permissionID); err != nil {
			t.Fatalf("grant evaluation permission %s: %v", permission, err)
		}
	}
	if _, err := tx.Exec(ctx, `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`, actor.userID, roleID); err != nil {
		t.Fatalf("assign evaluation role: %v", err)
	}
}

func seedEvaluationActor(t *testing.T, ctx context.Context, tx pgx.Tx) evaluationActor {
	t.Helper()
	actor := evaluationActor{userID: uuid.New(), employeeID: uuid.New()}
	if _, err := tx.Exec(ctx, `INSERT INTO custom_user (id, password, email) VALUES ($1, 'disabled', $2)`, actor.userID, uuid.NewString()+"@test.invalid"); err != nil {
		t.Fatalf("seed evaluation user: %v", err)
	}
	if _, err := tx.Exec(ctx, `INSERT INTO employee_profile (
		id, user_id, first_name, last_name, bsn, street, house_number, postal_code, city, gender
	) VALUES ($1, $2, 'Evaluation', 'Employee', $3, 'Test', '1', '1000AA', 'Test', 'unknown')`,
		actor.employeeID, actor.userID, uuid.NewString()); err != nil {
		t.Fatalf("seed evaluation employee: %v", err)
	}
	return actor
}

func assertClientSchedule(t *testing.T, clientID uuid.UUID, wantAnchor, wantNext time.Time) {
	t.Helper()
	var anchor, next time.Time
	if err := evaluationIntegrationPool.QueryRow(context.Background(), `SELECT last_evaluation_anchor_date, next_evaluation_date FROM client_details WHERE id = $1`, clientID).Scan(&anchor, &next); err != nil {
		t.Fatalf("read client schedule: %v", err)
	}
	if !sameDate(anchor, wantAnchor) || !sameDate(next, wantNext) {
		t.Fatalf("client schedule anchor=%s next=%s, want anchor=%s next=%s", anchor, next, wantAnchor, wantNext)
	}
}

func submitParams(t *testing.T, evaluationID uuid.UUID) domain.SubmitGoalEvaluationDraftParams {
	t.Helper()
	var updatedAt time.Time
	if err := evaluationIntegrationPool.QueryRow(context.Background(), `SELECT updated_at FROM client_goal_evaluations WHERE id = $1`, evaluationID).Scan(&updatedAt); err != nil {
		t.Fatalf("read evaluation revision: %v", err)
	}
	return domain.SubmitGoalEvaluationDraftParams{ExpectedUpdatedAt: updatedAt}
}

func sameDate(left, right time.Time) bool {
	leftYear, leftMonth, leftDay := left.Date()
	rightYear, rightMonth, rightDay := right.Date()
	return leftYear == rightYear && leftMonth == rightMonth && leftDay == rightDay
}
