package clientp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *clientService) CreateGoalEvaluation(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, req CreateGoalEvaluationRequest) (*GoalEvaluationResponse, error) {
	var evaluation db.ClientGoalEvaluation
	var items []db.GetGoalEvaluationItemsRow

	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		client, err := q.GetClientDetails(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to get client details: %w", err)
		}

		if client.Status != db.ClientStatusEnumInCare {
			return fmt.Errorf("evaluations can only be created for clients in care")
		}

		activeGoals, err := q.ListActiveGoalsByClientID(ctx, clientID)
		if err != nil {
			return fmt.Errorf("failed to list active goals: %w", err)
		}

		if len(activeGoals) == 0 {
			return fmt.Errorf("client must have at least one active goal to start an evaluation")
		}

		evaluationDate := nullableDate(req.EvaluationDate) // refactor
		if !evaluationDate.Valid {
			evaluationDate = pgtype.Date{Valid: true, Time: time.Now().UTC()}
		}
		if client.NextEvaluationDate.Valid {
			evaluationDate = client.NextEvaluationDate
		}

		intervalWeeks := client.EvaluationIntervalsWeeks
		if req.EvaluationIntervalWeeks != nil {
			intervalWeeks = *req.EvaluationIntervalWeeks
		}
		if intervalWeeks <= 0 {
			intervalWeeks = 12
		}

		periodStart := nullableDate(req.PeriodStart)
		periodEnd := nullableDate(req.PeriodEnd)

		if !periodStart.Valid {
			periodStart = pgtype.Date{Valid: false}
		}
		if client.LastEvaluationAnchorDate.Valid {
			periodStart = client.LastEvaluationAnchorDate
		}
		if !periodEnd.Valid {
			periodEnd = evaluationDate
		}

		evaluation, err = q.CreateGoalEvaluation(ctx, db.CreateGoalEvaluationParams{
			ClientID:                clientID,
			EvaluationDate:          evaluationDate,
			PeriodStart:             periodStart,
			PeriodEnd:               periodEnd,
			EvaluationIntervalWeeks: intervalWeeks,
			Status:                  db.EvaluationStatusEnumDraft,
			OverallNotes:            req.OverallNotes,
			CreatedByEmployeeID:     &employeeID,
		})
		if err != nil {
			return fmt.Errorf("failed to create evaluation header: %w", err)
		}

		for _, goal := range activeGoals {
			progress := db.ClientGoalProgressEnumNoProgress
			var notes *string

			if reqItem, ok := findDraftItemForGoal(req.Items, goal.ID); ok {
				parsedProgress, parseErr := parseProgress(reqItem.Progress)
				if parseErr != nil {
					return parseErr
				}
				progress = parsedProgress
				notes = reqItem.Notes
			}

			_, err = q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{
				EvaluationID: evaluation.ID,
				GoalID:       goal.ID,
				Progress:     progress,
				Notes:        notes,
			})
			if err != nil {
				return fmt.Errorf("failed to create evaluation item for goal %s: %w", goal.ID, err)
			}
		}

		if req.Submit {
			if err := ensureAllGoalsEvaluated(ctx, q, evaluation.ID); err != nil {
				return err
			}

			evaluation, err = q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
				ID:     evaluation.ID,
				Status: db.NullEvaluationStatusEnum{EvaluationStatusEnum: db.EvaluationStatusEnumCompleted, Valid: true},
			})
			if err != nil {
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == "P0001" {
					return fmt.Errorf("cannot submit evaluation yet: %s", pgErr.Message)
				}
				return fmt.Errorf("failed to submit evaluation: %w", err)
			}
		}

		items, err = q.GetGoalEvaluationItems(ctx, evaluation.ID)
		if err != nil {
			return fmt.Errorf("failed to load evaluation items: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapGoalEvaluationToResponse(evaluation, items), nil
}

func (s *clientService) SaveGoalEvaluationDraft(ctx context.Context, evaluationID uuid.UUID, req SaveGoalEvaluationDraftRequest) (*GoalEvaluationResponse, error) {
	var evaluation db.ClientGoalEvaluation
	var items []db.GetGoalEvaluationItemsRow

	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		existingEval, err := q.GetGoalEvaluation(ctx, evaluationID)
		if err != nil {
			return fmt.Errorf("evaluation not found: %w", err)
		}

		if existingEval.Status != db.EvaluationStatusEnumDraft {
			return fmt.Errorf("only draft evaluations can be saved")
		}

		client, err := q.GetClientDetails(ctx, existingEval.ClientID)
		if err != nil {
			return fmt.Errorf("failed to get evaluation client: %w", err)
		}

		if client.Status != db.ClientStatusEnumInCare {
			return fmt.Errorf("evaluations can only be saved for clients in care")
		}

		existingItems, err := q.GetGoalEvaluationItems(ctx, evaluationID)
		if err != nil {
			return fmt.Errorf("failed to load evaluation items: %w", err)
		}

		allowedGoalIDs := make(map[uuid.UUID]struct{}, len(existingItems))
		for _, item := range existingItems {
			allowedGoalIDs[item.GoalID] = struct{}{}
		}

		updateParams := db.UpdateGoalEvaluationParams{
			ID:                      evaluationID,
			EvaluationDate:          nullableDate(req.EvaluationDate),
			PeriodStart:             nullableDate(req.PeriodStart),
			PeriodEnd:               nullableDate(req.PeriodEnd),
			EvaluationIntervalWeeks: req.EvaluationIntervalWeeks,
			OverallNotes:            req.OverallNotes,
			Status:                  db.NullEvaluationStatusEnum{Valid: false},
		}

		evaluation, err = q.UpdateGoalEvaluation(ctx, updateParams)
		if err != nil {
			return fmt.Errorf("failed to update evaluation header: %w", err)
		}

		for _, item := range req.Items {
			if _, ok := allowedGoalIDs[item.GoalID]; !ok {
				return fmt.Errorf("goal %s is not part of this evaluation", item.GoalID)
			}

			parsedProgress, parseErr := parseProgress(item.Progress)
			if parseErr != nil {
				return parseErr
			}

			_, err = q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{
				EvaluationID: evaluationID,
				GoalID:       item.GoalID,
				Progress:     parsedProgress,
				Notes:        item.Notes,
			})
			if err != nil {
				return fmt.Errorf("failed to save item for goal %s: %w", item.GoalID, err)
			}
		}

		items, err = q.GetGoalEvaluationItems(ctx, evaluationID)
		if err != nil {
			return fmt.Errorf("failed to reload evaluation items: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapGoalEvaluationToResponse(evaluation, items), nil
}

func (s *clientService) SubmitGoalEvaluationDraft(ctx context.Context, evaluationID uuid.UUID) (*GoalEvaluationResponse, error) {
	var evaluation db.ClientGoalEvaluation
	var items []db.GetGoalEvaluationItemsRow

	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		existingEval, err := q.GetGoalEvaluation(ctx, evaluationID)
		if err != nil {
			return fmt.Errorf("evaluation not found: %w", err)
		}

		if existingEval.Status != db.EvaluationStatusEnumDraft {
			return fmt.Errorf("only draft evaluations can be submitted")
		}

		if err := ensureAllGoalsEvaluated(ctx, q, evaluationID); err != nil {
			return err
		}

		evaluation, err = q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
			ID:     evaluationID,
			Status: db.NullEvaluationStatusEnum{EvaluationStatusEnum: db.EvaluationStatusEnumCompleted, Valid: true},
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "P0001" {
				return fmt.Errorf("cannot submit evaluation yet: %s", pgErr.Message)
			}
			return fmt.Errorf("failed to submit evaluation: %w", err)
		}

		items, err = q.GetGoalEvaluationItems(ctx, evaluationID)
		if err != nil {
			return fmt.Errorf("failed to reload evaluation items: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return mapGoalEvaluationToResponse(evaluation, items), nil
}

func nullableDate(t *time.Time) pgtype.Date {
	if t == nil {
		return pgtype.Date{Valid: false}
	}
	return pgtype.Date{Time: *t, Valid: true}
}

func ensureAllGoalsEvaluated(ctx context.Context, q *db.Queries, evaluationID uuid.UUID) error {
	items, err := q.GetGoalEvaluationItems(ctx, evaluationID)
	if err != nil {
		return fmt.Errorf("failed to load evaluation items: %w", err)
	}
	if len(items) == 0 {
		return fmt.Errorf("evaluation has no goal items")
	}

	for _, item := range items {
		if item.Progress == db.ClientGoalProgressEnumNoProgress {
			return fmt.Errorf("all goals must be evaluated before submit (missing goal: %s)", item.GoalTitle)
		}
	}

	return nil
}

func findDraftItemForGoal(items []SaveGoalEvaluationDraftItem, goalID uuid.UUID) (SaveGoalEvaluationDraftItem, bool) {
	for _, item := range items {
		if item.GoalID == goalID {
			return item, true
		}
	}
	return SaveGoalEvaluationDraftItem{}, false
}

func parseProgress(value string) (db.ClientGoalProgressEnum, error) {
	progress := strings.TrimSpace(value)
	if progress == "" {
		return db.ClientGoalProgressEnumNoProgress, nil
	}

	allowed := map[string]db.ClientGoalProgressEnum{
		string(db.ClientGoalProgressEnumNoProgress):      db.ClientGoalProgressEnumNoProgress,
		string(db.ClientGoalProgressEnumRegression):      db.ClientGoalProgressEnumRegression,
		string(db.ClientGoalProgressEnumLimitedProgress): db.ClientGoalProgressEnumLimitedProgress,
		string(db.ClientGoalProgressEnumGoodProgress):    db.ClientGoalProgressEnumGoodProgress,
		string(db.ClientGoalProgressEnumAchieved):        db.ClientGoalProgressEnumAchieved,
		string(db.ClientGoalProgressEnumBlocked):         db.ClientGoalProgressEnumBlocked,
	}

	parsed, ok := allowed[progress]
	if !ok {
		return "", fmt.Errorf("invalid progress value: %s", value)
	}
	return parsed, nil
}

func mapGoalEvaluationToResponse(eval db.ClientGoalEvaluation, items []db.GetGoalEvaluationItemsRow) *GoalEvaluationResponse {
	res := &GoalEvaluationResponse{
		ID:                      eval.ID,
		ClientID:                eval.ClientID,
		EvaluationDate:          eval.EvaluationDate.Time,
		EvaluationIntervalWeeks: eval.EvaluationIntervalWeeks,
		Status:                  string(eval.Status),
		OverallNotes:            eval.OverallNotes,
		CreatedByEmployeeID:     eval.CreatedByEmployeeID,
		CreatedAt:               eval.CreatedAt.Time,
		UpdatedAt:               eval.UpdatedAt.Time,
	}

	if eval.PeriodStart.Valid {
		t := eval.PeriodStart.Time
		res.PeriodStart = &t
	}
	if eval.PeriodEnd.Valid {
		t := eval.PeriodEnd.Time
		res.PeriodEnd = &t
	}

	res.Items = make([]GoalEvaluationItemResponse, 0, len(items))
	for _, item := range items {
		res.Items = append(res.Items, GoalEvaluationItemResponse{
			ID:                item.ID,
			EvaluationID:      item.EvaluationID,
			GoalID:            item.GoalID,
			GoalTitle:         item.GoalTitle,
			GoalDescription:   item.GoalDescription,
			TopicNameSnapshot: item.TopicNameSnapshot,
			Progress:          string(item.Progress),
			Notes:             item.Notes,
			CreatedAt:         item.CreatedAt.Time,
			UpdatedAt:         item.UpdatedAt.Time,
		})
	}

	return res
}

func (s *clientService) ListUpcomingEvaluations(ctx *gin.Context, coordinatorID uuid.UUID, req ListUpcomingEvaluationsRequest) (*pagination.Response[ListUpcomingEvaluationsResponse], error) {
	params := req.GetParams()

	rows, err := s.Store.ListUpcomingEvaluationsForCoordinator(ctx, db.ListUpcomingEvaluationsForCoordinatorParams{
		EmployeeID: coordinatorID,
		Limit:      params.Limit,
		Offset:     params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list upcoming evaluations: %w", err)
	}

	if len(rows) == 0 {
		empty := pagination.NewResponse(ctx, req.Request, []ListUpcomingEvaluationsResponse{}, 0)
		return &empty, nil
	}

	items := make([]ListUpcomingEvaluationsResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, ListUpcomingEvaluationsResponse{
			ClientID:         row.ClientID,
			ClientFirstName:  row.ClientFirstName,
			ClientLastName:   row.ClientLastName,
			DueDate:          row.NextEvaluationDate.Time,
			DaysLeft:         row.DaysLeft,
			Priority:         row.Priority,
			HasDraft:         row.HasDraft,
			FilledGoalsCount: row.FilledGoalsCount,
			TotalGoalsCount:  row.TotalGoalsCount,
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, items, rows[0].TotalCount)
	return &pag, nil
}

func (s *clientService) ListRecentSubmittedEvaluations(ctx *gin.Context, employeeID uuid.UUID, req ListRecentSubmittedEvaluationsRequest) (*pagination.Response[ListRecentSubmittedEvaluationsResponse], error) {
	params := req.GetParams()

	rows, err := s.Store.ListRecentSubmittedEvaluationsByEmployee(ctx, db.ListRecentSubmittedEvaluationsByEmployeeParams{
		CreatedByEmployeeID: &employeeID,
		Limit:               params.Limit,
		Offset:              params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list recent submitted evaluations: %w", err)
	}

	if len(rows) == 0 {
		empty := pagination.NewResponse(ctx, req.Request, []ListRecentSubmittedEvaluationsResponse{}, 0)
		return &empty, nil
	}

	items := make([]ListRecentSubmittedEvaluationsResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, ListRecentSubmittedEvaluationsResponse{
			EvaluationID:       row.ID,
			ClientID:           row.ClientID,
			ClientFirstName:    row.ClientFirstName,
			ClientLastName:     row.ClientLastName,
			EvaluationDate:     row.EvaluationDate.Time,
			SubmittedAt:        row.SubmittedAt.Time,
			NextEvaluationDate: datePtr(row.NextEvaluationDate),
			FilledGoalsCount:   row.FilledGoalsCount,
			TotalGoalsCount:    row.TotalGoalsCount,
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, items, rows[0].TotalCount)
	return &pag, nil
}

func (s *clientService) ListRecentDraftEvaluations(ctx *gin.Context, employeeID uuid.UUID, req ListRecentDraftEvaluationsRequest) (*pagination.Response[ListRecentDraftEvaluationsResponse], error) {
	params := req.GetParams()

	rows, err := s.Store.ListRecentDraftEvaluationsByEmployee(ctx, db.ListRecentDraftEvaluationsByEmployeeParams{
		CreatedByEmployeeID: &employeeID,
		Limit:               params.Limit,
		Offset:              params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list recent draft evaluations: %w", err)
	}

	if len(rows) == 0 {
		empty := pagination.NewResponse(ctx, req.Request, []ListRecentDraftEvaluationsResponse{}, 0)
		return &empty, nil
	}

	items := make([]ListRecentDraftEvaluationsResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, ListRecentDraftEvaluationsResponse{
			EvaluationID:     row.ID,
			ClientID:         row.ClientID,
			ClientFirstName:  row.ClientFirstName,
			ClientLastName:   row.ClientLastName,
			DueDate:          row.EvaluationDate.Time,
			UpdatedAt:        row.UpdatedAt.Time,
			DaysLeft:         row.DaysLeft,
			Priority:         row.Priority,
			FilledGoalsCount: row.FilledGoalsCount,
			TotalGoalsCount:  row.TotalGoalsCount,
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, items, rows[0].TotalCount)
	return &pag, nil
}

func (s *clientService) GetGoalEvaluationBootstrap(ctx context.Context, clientID uuid.UUID) (*GoalEvaluationBootstrapResponse, error) {
	client, err := s.Store.GetClientDetails(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client details: %w", err)
	}

	response := &GoalEvaluationBootstrapResponse{
		ClientID:        client.ID,
		ClientFirstName: client.FirstName,
		ClientLastName:  client.LastName,
		ActiveGoals:     []GoalEvaluationBootstrapActiveGoal{},
	}

	if client.NextEvaluationDate.Valid {
		nextDate := client.NextEvaluationDate.Time
		response.NextEvaluationDate = &nextDate
		today := time.Now().UTC().Truncate(24 * time.Hour)
		due := nextDate.UTC().Truncate(24 * time.Hour)
		daysLeft := int32(due.Sub(today).Hours() / 24)
		response.DaysLeft = &daysLeft
		priority := "normal"
		if daysLeft <= 3 {
			priority = "critical"
		}
		response.Priority = &priority
	}

	latestDraft, err := s.Store.GetLatestDraftEvaluationByClient(ctx, clientID)
	if err == nil {
		response.ExistingDraft = &GoalEvaluationBootstrapDraft{
			ID:             latestDraft.ID,
			EvaluationDate: latestDraft.EvaluationDate.Time,
			UpdatedAt:      latestDraft.UpdatedAt.Time,
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to get latest draft evaluation: %w", err)
	}

	latestCompleted, err := s.Store.GetLatestCompletedEvaluationByClient(ctx, clientID)
	if err == nil {
		response.LastCompletedEvaluation = &GoalEvaluationBootstrapCompleted{
			ID:                  latestCompleted.ID,
			EvaluationDate:      latestCompleted.EvaluationDate.Time,
			SubmittedAt:         latestCompleted.SubmittedAt.Time,
			OverallNotes:        latestCompleted.OverallNotes,
			CreatedByEmployeeID: latestCompleted.CreatedByEmployeeID,
			CreatorName:         composeCreatorName(latestCompleted.CreatorFirstName, latestCompleted.CreatorLastName),
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed to get latest completed evaluation: %w", err)
	}

	activeGoals, err := s.Store.ListActiveGoalsByClientID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to list active goals: %w", err)
	}

	latestProgressRows, err := s.Store.ListLatestCompletedGoalProgressByClient(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to list latest completed goal progress: %w", err)
	}

	latestProgressByGoal := make(map[uuid.UUID]db.ListLatestCompletedGoalProgressByClientRow, len(latestProgressRows))
	for _, row := range latestProgressRows {
		latestProgressByGoal[row.GoalID] = row
	}

	response.ActiveGoals = make([]GoalEvaluationBootstrapActiveGoal, 0, len(activeGoals))
	for _, goal := range activeGoals {
		goalResponse := GoalEvaluationBootstrapActiveGoal{
			GoalID:            goal.ID,
			Title:             goal.Title,
			TopicNameSnapshot: goal.TopicNameSnapshot,
			Priority:          string(goal.Priority),
			SortOrder:         goal.SortOrder,
		}

		if latestGoalProgress, ok := latestProgressByGoal[goal.ID]; ok {
			p := string(latestGoalProgress.Progress)
			goalResponse.LastProgress = &p
			goalResponse.LastNotes = latestGoalProgress.Notes
		}

		response.ActiveGoals = append(response.ActiveGoals, goalResponse)
	}

	return response, nil
}

func datePtr(date pgtype.Date) *time.Time {
	if !date.Valid {
		return nil
	}
	t := date.Time
	return &t
}

func composeCreatorName(firstName, lastName *string) *string {
	if firstName == nil && lastName == nil {
		return nil
	}
	parts := make([]string, 0, 2)
	if firstName != nil && strings.TrimSpace(*firstName) != "" {
		parts = append(parts, strings.TrimSpace(*firstName))
	}
	if lastName != nil && strings.TrimSpace(*lastName) != "" {
		parts = append(parts, strings.TrimSpace(*lastName))
	}
	if len(parts) == 0 {
		return nil
	}
	name := strings.Join(parts, " ")
	return &name
}
