package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"maicare_go/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type goalEvaluationServiceStub struct {
	domain.ClientService
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

func (s *goalEvaluationServiceStub) UpdateGoalEvaluationDraft(_ context.Context, evaluationID uuid.UUID, employeeID uuid.UUID, params domain.UpdateGoalEvaluationDraftParams) (*domain.GoalEvaluation, error) {
	s.updateCalls++
	s.updateEvaluationID = evaluationID
	s.updateEmployeeID = employeeID
	s.updateParams = params
	return s.updateResult, s.updateErr
}

func (s *goalEvaluationServiceStub) SubmitGoalEvaluationDraft(_ context.Context, evaluationID uuid.UUID, employeeID uuid.UUID) (*domain.GoalEvaluation, error) {
	s.submitCalls++
	s.submitEvaluationID = evaluationID
	s.submitEmployeeID = employeeID
	return s.submitResult, s.submitErr
}

func evaluationActor(employeeID uuid.UUID) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("employee_id", employeeID.String())
		ctx.Next()
	}
}

func TestUpdateGoalEvaluationDraftUsesPathIDAndActor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	evaluationID := uuid.New()
	employeeID := uuid.New()
	goalID := uuid.New()
	service := &goalEvaluationServiceStub{updateResult: &domain.GoalEvaluation{ID: evaluationID, ClientID: uuid.New(), Status: "draft"}}
	router := gin.New()
	router.PATCH("/evaluations/:evaluation_id/draft", evaluationActor(employeeID), NewClientHandler(service).UpdateGoalEvaluationDraft)
	request := httptest.NewRequest(http.MethodPatch, "/evaluations/"+evaluationID.String()+"/draft", bytes.NewBufferString(`{"items":[{"goal_id":"`+goalID.String()+`","progress":"good_progress"}]}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if service.updateCalls != 1 || service.updateEvaluationID != evaluationID || service.updateEmployeeID != employeeID {
		t.Fatal("handler did not pass exact evaluation and employee IDs")
	}
	if len(service.updateParams.Items) != 1 || service.updateParams.Items[0].GoalID != goalID {
		t.Fatalf("update params = %#v", service.updateParams)
	}
}

func TestUpdateGoalEvaluationDraftMapsOwnershipConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &goalEvaluationServiceStub{updateErr: domain.ErrGoalEvaluationOwnedByOther}
	router := gin.New()
	router.PATCH("/evaluations/:evaluation_id/draft", evaluationActor(uuid.New()), NewClientHandler(service).UpdateGoalEvaluationDraft)
	request := httptest.NewRequest(http.MethodPatch, "/evaluations/"+uuid.NewString()+"/draft", bytes.NewBufferString(`{"items":[]}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusConflict)
	}
	if !bytes.Contains(response.Body.Bytes(), []byte(`"code":"EVALUATION_NOT_OWNER"`)) {
		t.Fatalf("body = %s, want stable ownership code", response.Body.String())
	}
}

func TestUpdateGoalEvaluationDraftMapsCurrentCycleConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &goalEvaluationServiceStub{updateErr: domain.ErrGoalEvaluationNotCurrentCycle}
	router := gin.New()
	router.PATCH("/evaluations/:evaluation_id/draft", evaluationActor(uuid.New()), NewClientHandler(service).UpdateGoalEvaluationDraft)
	request := httptest.NewRequest(http.MethodPatch, "/evaluations/"+uuid.NewString()+"/draft", bytes.NewBufferString(`{"items":[]}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusConflict)
	}
	if !bytes.Contains(response.Body.Bytes(), []byte(`"code":"EVALUATION_NOT_CURRENT_CYCLE"`)) {
		t.Fatalf("body = %s, want stable current-cycle code", response.Body.String())
	}
}

func TestSubmitGoalEvaluationDraftRequiresNoBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	evaluationID := uuid.New()
	employeeID := uuid.New()
	service := &goalEvaluationServiceStub{submitResult: &domain.GoalEvaluation{ID: evaluationID, ClientID: uuid.New(), Status: "completed"}}
	router := gin.New()
	router.POST("/evaluations/:evaluation_id/submit", evaluationActor(employeeID), NewClientHandler(service).SubmitGoalEvaluationDraft)
	request := httptest.NewRequest(http.MethodPost, "/evaluations/"+evaluationID.String()+"/submit", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if service.submitCalls != 1 || service.submitEvaluationID != evaluationID || service.submitEmployeeID != employeeID {
		t.Fatal("handler did not pass exact evaluation and employee IDs")
	}
}
