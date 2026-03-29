package clientp

import (
	"context"
	"errors"
	"fmt"
	"strings"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

var ErrClientGoalClientNotFound = errors.New("client for goal creation not found")
var ErrClientGoalGoalNotFound = errors.New("client goal not found")
var ErrClientGoalTopicNotFound = errors.New("topic for goal creation not found")
var ErrClientGoalDraftEvaluationExists = errors.New("client goal updates are blocked while a draft evaluation exists")
var ErrClientGoalEmptyPatch = errors.New("at least one goal field must be provided")
var ErrClientGoalTitleRequired = errors.New("title is required")

func (s *clientService) CreateClientGoal(ctx context.Context, clientID uuid.UUID, req CreateClientGoalRequest) (*CreateClientGoalResponse, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, ErrClientGoalTitleRequired
	}

	priority := db.ClientGoalPriorityEnumMedium
	if req.Priority != nil && strings.TrimSpace(*req.Priority) != "" {
		priority = db.ClientGoalPriorityEnum(strings.TrimSpace(*req.Priority))
	}

	var createdGoal db.ClientGoal
	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		if _, err := q.GetClientDetails(ctx, clientID); err != nil {
			return fmt.Errorf("%w: %w", ErrClientGoalClientNotFound, err)
		}

		topic, err := q.GetTopicByID(ctx, req.TopicID)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrClientGoalTopicNotFound, err)
		}

		sortOrder := int32(0)
		if req.SortOrder != nil {
			sortOrder = *req.SortOrder
		} else {
			nextSortOrder, err := q.GetNextActiveClientGoalSortOrder(ctx, clientID)
			if err != nil {
				return fmt.Errorf("failed to compute next goal sort order: %w", err)
			}
			sortOrder = nextSortOrder
		}

		createdGoal, err = q.CreateManualClientGoal(ctx, db.CreateManualClientGoalParams{
			ClientID:          clientID,
			Title:             title,
			Description:       req.Description,
			Priority:          priority,
			TopicID:           &topic.ID,
			TopicNameSnapshot: &topic.TopicName,
			SortOrder:         sortOrder,
		})
		if err != nil {
			return fmt.Errorf("failed to create client goal: %w", err)
		}

		return nil
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateClientGoal", "Failed to create client goal", zap.Error(err), zap.String("client_id", clientID.String()))
		return nil, err
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "CreateClientGoal", "Created client goal successfully", zap.String("client_id", clientID.String()), zap.String("goal_id", createdGoal.ID.String()))

	return mapClientGoal(createdGoal), nil
}

func (s *clientService) UpdateClientGoal(ctx context.Context, clientID uuid.UUID, goalID uuid.UUID, req UpdateClientGoalRequest) (*UpdateClientGoalResponse, error) {
	if req.Title == nil && req.Description == nil && req.Priority == nil && req.TopicID == nil && req.SortOrder == nil {
		return nil, ErrClientGoalEmptyPatch
	}

	var (
		activeGoal        db.ClientGoal
		replacementGoalID *uuid.UUID
		mutationType      = "updated"
	)

	err := s.Store.ExecTx(ctx, func(q *db.Queries) error {
		if _, err := q.GetClientDetails(ctx, clientID); err != nil {
			return fmt.Errorf("%w: %w", ErrClientGoalClientNotFound, err)
		}

		existingGoal, err := q.GetClientGoalByIDAndClientID(ctx, db.GetClientGoalByIDAndClientIDParams{
			ID:       goalID,
			ClientID: clientID,
		})
		if err != nil {
			return fmt.Errorf("%w: %w", ErrClientGoalGoalNotFound, err)
		}
		if existingGoal.Status != db.ClientGoalStatusEnumActive {
			return ErrClientGoalGoalNotFound
		}

		if _, err := q.GetLatestDraftEvaluationByClient(ctx, clientID); err == nil {
			return ErrClientGoalDraftEvaluationExists
		} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to check for draft evaluations: %w", err)
		}

		nextTitle := existingGoal.Title
		if req.Title != nil {
			trimmedTitle := strings.TrimSpace(*req.Title)
			if trimmedTitle == "" {
				return ErrClientGoalTitleRequired
			}
			nextTitle = trimmedTitle
		}

		nextDescription := existingGoal.Description
		if req.Description != nil {
			nextDescription = normalizeOptionalTrimmedString(req.Description)
		}

		nextPriority := existingGoal.Priority
		if req.Priority != nil && strings.TrimSpace(*req.Priority) != "" {
			nextPriority = db.ClientGoalPriorityEnum(strings.TrimSpace(*req.Priority))
		}

		nextTopicID := existingGoal.TopicID
		nextTopicNameSnapshot := existingGoal.TopicNameSnapshot
		if req.TopicID != nil {
			topic, err := q.GetTopicByID(ctx, *req.TopicID)
			if err != nil {
				return fmt.Errorf("%w: %w", ErrClientGoalTopicNotFound, err)
			}
			nextTopicID = &topic.ID
			nextTopicNameSnapshot = &topic.TopicName
		}

		nextSortOrder := existingGoal.SortOrder
		if req.SortOrder != nil {
			nextSortOrder = *req.SortOrder
		}

		contentChanged := nextTitle != existingGoal.Title ||
			!equalOptionalString(nextDescription, existingGoal.Description) ||
			!equalOptionalUUID(nextTopicID, existingGoal.TopicID)

		priorityChanged := nextPriority != existingGoal.Priority
		sortOrderChanged := nextSortOrder != existingGoal.SortOrder

		if !contentChanged && !priorityChanged && !sortOrderChanged {
			activeGoal = existingGoal
			return nil
		}

		hasHistory, err := q.GoalHasEvaluationItems(ctx, goalID)
		if err != nil {
			return fmt.Errorf("failed to check goal evaluation history: %w", err)
		}

		if hasHistory && contentChanged {
			if _, err := q.CancelClientGoalByID(ctx, db.CancelClientGoalByIDParams{
				ID:       goalID,
				ClientID: clientID,
			}); err != nil {
				return fmt.Errorf("failed to retire client goal: %w", err)
			}

			activeGoal, err = q.CreateReviewUpdatedClientGoal(ctx, db.CreateReviewUpdatedClientGoalParams{
				ClientID:          clientID,
				Title:             nextTitle,
				Description:       nextDescription,
				Priority:          nextPriority,
				TopicID:           nextTopicID,
				TopicNameSnapshot: nextTopicNameSnapshot,
				SortOrder:         nextSortOrder,
			})
			if err != nil {
				return fmt.Errorf("failed to create replacement client goal: %w", err)
			}

			mutationType = "versioned"
			replacementGoalID = &activeGoal.ID
			return nil
		}

		activeGoal, err = q.UpdateClientGoalByID(ctx, db.UpdateClientGoalByIDParams{
			ID:                goalID,
			ClientID:          clientID,
			Title:             nextTitle,
			Description:       nextDescription,
			Priority:          nextPriority,
			TopicID:           nextTopicID,
			TopicNameSnapshot: nextTopicNameSnapshot,
			SortOrder:         nextSortOrder,
		})
		if err != nil {
			return fmt.Errorf("failed to update client goal: %w", err)
		}

		return nil
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateClientGoal", "Failed to update client goal", zap.Error(err), zap.String("client_id", clientID.String()), zap.String("goal_id", goalID.String()))
		return nil, err
	}

	s.Logger.LogBusinessEvent(ctx, logger.LogLevelInfo, "UpdateClientGoal", "Updated client goal successfully", zap.String("client_id", clientID.String()), zap.String("goal_id", goalID.String()))

	return &UpdateClientGoalResponse{
		MutationType:      mutationType,
		GoalID:            goalID,
		ReplacementGoalID: replacementGoalID,
		Goal:              *mapClientGoal(activeGoal),
	}, nil
}

func mapClientGoal(goal db.ClientGoal) *CreateClientGoalResponse {
	return &CreateClientGoalResponse{
		ID:                goal.ID,
		ClientID:          goal.ClientID,
		Title:             goal.Title,
		Description:       goal.Description,
		Priority:          string(goal.Priority),
		Status:            string(goal.Status),
		TopicID:           goal.TopicID,
		TopicNameSnapshot: goal.TopicNameSnapshot,
		Source:            string(goal.Source),
		SortOrder:         goal.SortOrder,
		CreatedAt:         goal.CreatedAt.Time,
		UpdatedAt:         goal.UpdatedAt.Time,
	}
}

func normalizeOptionalTrimmedString(value *string) *string {
	if value == nil {
		return nil
	}

	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}

func equalOptionalString(a *string, b *string) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	default:
		return *a == *b
	}
}

func equalOptionalUUID(a *uuid.UUID, b *uuid.UUID) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil || b == nil:
		return false
	default:
		return *a == *b
	}
}
