package service

import (
	"context"
	"errors"
	"testing"

	"maicare_go/internal/domain"

	"github.com/google/uuid"
)

type goalEvaluationRepositoryStub struct {
	domain.ClientRepository
	updateEvaluationID uuid.UUID
	updateEmployeeID   uuid.UUID
	updateParams       domain.UpdateGoalEvaluationDraftParams
	updateResult       *domain.GoalEvaluation
	updateErr          error
	updateCalls        int
	submitEvaluationID uuid.UUID
	submitEmployeeID   uuid.UUID
	submitResult       *domain.GoalEvaluation
	submitErr          error
	submitCalls        int
}

func (s *goalEvaluationRepositoryStub) UpdateGoalEvaluationDraft(_ context.Context, evaluationID uuid.UUID, employeeID uuid.UUID, params domain.UpdateGoalEvaluationDraftParams) (*domain.GoalEvaluation, error) {
	s.updateCalls++
	s.updateEvaluationID = evaluationID
	s.updateEmployeeID = employeeID
	s.updateParams = params
	return s.updateResult, s.updateErr
}

func (s *goalEvaluationRepositoryStub) SubmitGoalEvaluationDraft(_ context.Context, evaluationID uuid.UUID, employeeID uuid.UUID) (*domain.GoalEvaluation, error) {
	s.submitCalls++
	s.submitEvaluationID = evaluationID
	s.submitEmployeeID = employeeID
	return s.submitResult, s.submitErr
}

func TestUpdateGoalEvaluationDraftDelegatesExactIDs(t *testing.T) {
	evaluationID := uuid.New()
	employeeID := uuid.New()
	goalID := uuid.New()
	want := &domain.GoalEvaluation{ID: evaluationID, ClientID: uuid.New()}
	repository := &goalEvaluationRepositoryStub{updateResult: want}
	service := NewClientService(repository, nil, nil, nil, nil, nil, nil)
	params := domain.UpdateGoalEvaluationDraftParams{Items: []domain.GoalEvaluationItemParams{{GoalID: goalID, Progress: "achieved"}}}

	got, err := service.UpdateGoalEvaluationDraft(context.Background(), evaluationID, employeeID, params)
	if err != nil {
		t.Fatalf("UpdateGoalEvaluationDraft() error = %v", err)
	}
	if got != want || repository.updateCalls != 1 || repository.updateEvaluationID != evaluationID || repository.updateEmployeeID != employeeID {
		t.Fatalf("UpdateGoalEvaluationDraft() did not delegate exact evaluation and employee IDs")
	}
	if len(repository.updateParams.Items) != 1 || repository.updateParams.Items[0].GoalID != goalID {
		t.Fatalf("UpdateGoalEvaluationDraft() params = %#v", repository.updateParams)
	}
}

func TestUpdateGoalEvaluationDraftRejectsDuplicateGoals(t *testing.T) {
	goalID := uuid.New()
	repository := &goalEvaluationRepositoryStub{}
	service := NewClientService(repository, nil, nil, nil, nil, nil, nil)
	params := domain.UpdateGoalEvaluationDraftParams{Items: []domain.GoalEvaluationItemParams{{GoalID: goalID}, {GoalID: goalID}}}

	_, err := service.UpdateGoalEvaluationDraft(context.Background(), uuid.New(), uuid.New(), params)
	if err == nil {
		t.Fatal("UpdateGoalEvaluationDraft() error = nil, want duplicate goal error")
	}
	if repository.updateCalls != 0 {
		t.Fatalf("repository update calls = %d, want 0", repository.updateCalls)
	}
}

func TestSubmitGoalEvaluationDraftDelegatesExactIDsAndError(t *testing.T) {
	evaluationID := uuid.New()
	employeeID := uuid.New()
	wantErr := errors.New("submit failed")
	repository := &goalEvaluationRepositoryStub{submitErr: wantErr}
	service := NewClientService(repository, nil, nil, nil, nil, nil, nil)

	_, err := service.SubmitGoalEvaluationDraft(context.Background(), evaluationID, employeeID)
	if !errors.Is(err, wantErr) {
		t.Fatalf("SubmitGoalEvaluationDraft() error = %v, want %v", err, wantErr)
	}
	if repository.submitCalls != 1 || repository.submitEvaluationID != evaluationID || repository.submitEmployeeID != employeeID {
		t.Fatalf("SubmitGoalEvaluationDraft() did not delegate exact evaluation and employee IDs")
	}
}
