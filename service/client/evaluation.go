package clientp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *clientService) CreateGoalEvaluation(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, req CreateGoalEvaluationRequest) (*GoalEvaluationResponse, error) {
	evaluation, items, err := s.saveCurrentGoalEvaluationDraft(ctx, clientID, employeeID, req)
	if err != nil {
		return nil, err
	}

	var submitErrMessage *string
	if req.Submit {
		submittedEvaluation, submitErr := s.trySubmitGoalEvaluationDraft(ctx, evaluation.ID)
		if submitErr != nil {
			var blockedErr *goalEvaluationSubmitBlockedError
			if errors.As(submitErr, &blockedErr) {
				msg := blockedErr.Error()
				submitErrMessage = &msg
			} else {
				return nil, submitErr
			}
		} else {
			evaluation = submittedEvaluation
		}
	}

	response := mapGoalEvaluationToResponse(evaluation, items)
	response.SubmitError = submitErrMessage
	return response, nil
}

func (s *clientService) GetGoalEvaluation(ctx context.Context, evaluationID uuid.UUID) (*GoalEvaluationResponse, error) {
	evaluation, err := s.Store.GetGoalEvaluationByID(ctx, evaluationID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("goal evaluation not found")
		}
		return nil, fmt.Errorf("failed to get evaluation: %w", err)
	}

	items, err := s.Store.GetGoalEvaluationItems(ctx, evaluationID)
	if err != nil {
		return nil, fmt.Errorf("failed to load evaluation items: %w", err)
	}

	res := mapGoalEvaluationToResponse(db.ClientGoalEvaluation{
		ID:                      evaluation.ID,
		ClientID:                evaluation.ClientID,
		EvaluationDate:          evaluation.EvaluationDate,
		PeriodStart:             evaluation.PeriodStart,
		PeriodEnd:               evaluation.PeriodEnd,
		EvaluationIntervalWeeks: evaluation.EvaluationIntervalWeeks,
		Status:                  evaluation.Status,
		OverallNotes:            evaluation.OverallNotes,
		CreatedByEmployeeID:     evaluation.CreatedByEmployeeID,
		CreatedAt:               evaluation.CreatedAt,
		UpdatedAt:               evaluation.UpdatedAt,
	}, items)
	res.CreatorName = composeCreatorName(evaluation.CreatorFirstName, evaluation.CreatorLastName)
	return res, nil
}

type goalEvaluationSubmitBlockedError struct {
	message string
}

func (e *goalEvaluationSubmitBlockedError) Error() string {
	return e.message
}

func (s *clientService) saveCurrentGoalEvaluationDraft(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID, req CreateGoalEvaluationRequest) (db.ClientGoalEvaluation, []db.GetGoalEvaluationItemsRow, error) {
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

		if !client.NextEvaluationDate.Valid {
			return fmt.Errorf("client has no next evaluation date configured")
		}

		reqItemsByGoal, err := mapDraftItemsByGoalID(req.Items)
		if err != nil {
			return err
		}

		evaluationDate := client.NextEvaluationDate
		periodStart := client.LastEvaluationAnchorDate
		periodEnd := client.NextEvaluationDate
		intervalWeeks := client.EvaluationIntervalsWeeks
		if intervalWeeks <= 0 {
			intervalWeeks = 12
		}

		evaluation, err = q.GetDraftGoalEvaluationByClientAndDate(ctx, db.GetDraftGoalEvaluationByClientAndDateParams{
			ClientID:       clientID,
			EvaluationDate: evaluationDate,
		})
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("failed to get current draft evaluation: %w", err)
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
				var pgErr *pgconn.PgError
				if errors.As(err, &pgErr) && pgErr.Code == "23505" {
					evaluation, err = q.GetDraftGoalEvaluationByClientAndDate(ctx, db.GetDraftGoalEvaluationByClientAndDateParams{
						ClientID:       clientID,
						EvaluationDate: evaluationDate,
					})
					if err != nil {
						if errors.Is(err, pgx.ErrNoRows) {
							return fmt.Errorf("current evaluation already exists and is not editable")
						}
						return fmt.Errorf("failed to load concurrent draft evaluation: %w", err)
					}
				} else {
					return fmt.Errorf("failed to create evaluation header: %w", err)
				}
			}
		}

		evaluation, err = q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
			ID:           evaluation.ID,
			OverallNotes: req.OverallNotes,
			Status:       db.NullEvaluationStatusEnum{Valid: false},
		})
		if err != nil {
			return fmt.Errorf("failed to update evaluation header: %w", err)
		}

		activeGoalIDs := make(map[uuid.UUID]struct{}, len(activeGoals))
		for _, goal := range activeGoals {
			activeGoalIDs[goal.ID] = struct{}{}
		}

		existingItems, err := q.GetGoalEvaluationItems(ctx, evaluation.ID)
		if err != nil {
			return fmt.Errorf("failed to load evaluation items: %w", err)
		}

		existingByGoalID := make(map[uuid.UUID]struct{}, len(existingItems))
		for _, item := range existingItems {
			existingByGoalID[item.GoalID] = struct{}{}
		}

		for _, goal := range activeGoals {
			if _, exists := existingByGoalID[goal.ID]; exists {
				continue
			}

			progress := db.ClientGoalProgressEnumNoProgress
			var notes *string
			if reqItem, ok := reqItemsByGoal[goal.ID]; ok {
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
				return fmt.Errorf("failed to ensure evaluation item for goal %s: %w", goal.ID, err)
			}
		}

		for _, reqItem := range req.Items {
			if _, isActive := activeGoalIDs[reqItem.GoalID]; !isActive {
				return fmt.Errorf("goal %s is not an active goal for this client", reqItem.GoalID)
			}

			parsedProgress, parseErr := parseProgress(reqItem.Progress)
			if parseErr != nil {
				return parseErr
			}

			_, err = q.UpsertGoalEvaluationItem(ctx, db.UpsertGoalEvaluationItemParams{
				EvaluationID: evaluation.ID,
				GoalID:       reqItem.GoalID,
				Progress:     parsedProgress,
				Notes:        reqItem.Notes,
			})
			if err != nil {
				return fmt.Errorf("failed to save item for goal %s: %w", reqItem.GoalID, err)
			}
		}

		items, err = q.GetGoalEvaluationItems(ctx, evaluation.ID)
		if err != nil {
			return fmt.Errorf("failed to reload evaluation items: %w", err)
		}

		return nil
	})
	if err != nil {
		return db.ClientGoalEvaluation{}, nil, err
	}

	return evaluation, items, nil
}

func (s *clientService) trySubmitGoalEvaluationDraft(ctx context.Context, evaluationID uuid.UUID) (db.ClientGoalEvaluation, error) {
	var evaluation db.ClientGoalEvaluation

	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		if err := ensureAllGoalsEvaluated(ctx, q, evaluationID); err != nil {
			return &goalEvaluationSubmitBlockedError{message: err.Error()}
		}

		updatedEvaluation, err := q.UpdateGoalEvaluation(ctx, db.UpdateGoalEvaluationParams{
			ID:     evaluationID,
			Status: db.NullEvaluationStatusEnum{EvaluationStatusEnum: db.EvaluationStatusEnumCompleted, Valid: true},
		})
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "P0001" {
				return &goalEvaluationSubmitBlockedError{message: fmt.Sprintf("cannot submit evaluation yet: %s", pgErr.Message)}
			}
			return fmt.Errorf("failed to submit evaluation: %w", err)
		}

		evaluation = updatedEvaluation
		return nil
	})
	if err != nil {
		return db.ClientGoalEvaluation{}, err
	}

	return evaluation, nil
}

func mapDraftItemsByGoalID(items []SaveGoalEvaluationDraftItem) (map[uuid.UUID]SaveGoalEvaluationDraftItem, error) {
	itemsByGoal := make(map[uuid.UUID]SaveGoalEvaluationDraftItem, len(items))
	for _, item := range items {
		if _, exists := itemsByGoal[item.GoalID]; exists {
			return nil, fmt.Errorf("goal %s appears multiple times in request", item.GoalID)
		}
		itemsByGoal[item.GoalID] = item
	}
	return itemsByGoal, nil
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
			NextEvaluationDate: util.DatePtr(row.NextEvaluationDate),
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

func (s *clientService) GetClientGoalsForEvaluationPage(ctx context.Context, clientID uuid.UUID, employeeID uuid.UUID) (*GetClientGoalsForEvaluationPageResponse, error) {
	client, err := s.Store.GetClientDetails(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client details: %w", err)
	}

	response := &GetClientGoalsForEvaluationPageResponse{
		Goals: []ClientGoalForEvaluationPageResponse{},
	}

	if client.NextEvaluationDate.Valid {
		nextDate := client.NextEvaluationDate.Time
		response.NextEvaluationDate = &nextDate
	}

	coordinatorRows, err := s.Store.GetClientCoordinator(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get client coordinator: %w", err)
	}
	response.IsResponsibleEmployee = len(coordinatorRows) > 0 && coordinatorRows[0].EmployeeID == employeeID

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

	response.Goals = make([]ClientGoalForEvaluationPageResponse, 0, len(activeGoals))
	for _, goal := range activeGoals {
		goalResponse := ClientGoalForEvaluationPageResponse{
			ID:        goal.ID,
			TopicName: goal.TopicNameSnapshot,
			Title:     goal.Title,
			Priority:  string(goal.Priority),
		}

		if latestGoalProgress, ok := latestProgressByGoal[goal.ID]; ok {
			p := string(latestGoalProgress.Progress)
			goalResponse.LastEvaluationProgress = &p
		}

		response.Goals = append(response.Goals, goalResponse)
	}

	draftEval, err := s.Store.GetCurrentCycleDraftEvaluationByClientAndEmployee(ctx, db.GetCurrentCycleDraftEvaluationByClientAndEmployeeParams{
		ClientID:            clientID,
		CreatedByEmployeeID: &employeeID,
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get current cycle draft evaluation: %w", err)
		}
		return response, nil
	}

	response.MyDraftEvaluationID = &draftEval.ID
	return response, nil
}

func (s *clientService) ListClientSubmittedEvaluations(ctx *gin.Context, clientID uuid.UUID, req ListClientSubmittedEvaluationsRequest) (*pagination.Response[ListClientSubmittedEvaluationsResponse], error) {
	params := req.GetParams()

	rows, err := s.Store.ListSubmittedEvaluationsByClient(ctx, db.ListSubmittedEvaluationsByClientParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list submitted evaluations: %w", err)
	}

	if len(rows) == 0 {
		empty := pagination.NewResponse(ctx, req.Request, []ListClientSubmittedEvaluationsResponse{}, 0)
		return &empty, nil
	}

	items := make([]ListClientSubmittedEvaluationsResponse, 0, len(rows))
	for _, row := range rows {
		items = append(items, ListClientSubmittedEvaluationsResponse{
			EvaluationID:        row.ID,
			EvaluationDate:      row.EvaluationDate.Time,
			SubmittedAt:         row.SubmittedAt.Time,
			FilledGoalsCount:    row.FilledGoalsCount,
			TotalGoalsCount:     row.TotalGoalsCount,
			CreatedByEmployeeID: row.CreatedByEmployeeID,
			CreatorName:         composeCreatorName(row.CreatorFirstName, row.CreatorLastName),
		})
	}

	pag := pagination.NewResponse(ctx, req.Request, items, rows[0].TotalCount)
	return &pag, nil
}

func (s *clientService) ListGoalEvaluationHistory(ctx *gin.Context, clientID uuid.UUID, goalID uuid.UUID, req ListGoalEvaluationHistoryRequest) (*pagination.Response[ListGoalEvaluationHistoryResponse], error) {
	params := req.GetParams()

	rows, err := s.Store.ListGoalEvaluationHistoryByClientAndGoal(ctx, db.ListGoalEvaluationHistoryByClientAndGoalParams{
		ClientID: clientID,
		GoalID:   goalID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list goal evaluation history: %w", err)
	}

	if len(rows) == 0 {
		empty := pagination.NewResponse(ctx, req.Request, []ListGoalEvaluationHistoryResponse{}, 0)
		return &empty, nil
	}

	items := make([]ListGoalEvaluationHistoryResponse, 0, len(rows))
	for _, row := range rows {
		item := ListGoalEvaluationHistoryResponse{
			EvaluationID:        row.EvaluationID,
			EvaluationDate:      row.EvaluationDate.Time,
			SubmittedAt:         row.SubmittedAt.Time,
			Progress:            string(row.Progress),
			Notes:               row.Notes,
			CreatedByEmployeeID: row.CreatedByEmployeeID,
			CreatorName:         composeCreatorName(row.CreatorFirstName, row.CreatorLastName),
		}
		if row.PeriodStart.Valid {
			t := row.PeriodStart.Time
			item.PeriodStart = &t
		}
		if row.PeriodEnd.Valid {
			t := row.PeriodEnd.Time
			item.PeriodEnd = &t
		}
		items = append(items, item)
	}

	pag := pagination.NewResponse(ctx, req.Request, items, rows[0].TotalCount)
	return &pag, nil
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
