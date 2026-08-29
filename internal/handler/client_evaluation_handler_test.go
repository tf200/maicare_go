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
	createErr          error
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

func (s *goalEvaluationServiceStub) CreateGoalEvaluation(_ context.Context, _ uuid.UUID, _ uuid.UUID, _ domain.CreateGoalEvaluationParams) (*domain.GoalEvaluation, error) {
	return nil, s.createErr
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

func TestGoalEvaluationMutationErrorCodes(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "not found", err: domain.ErrGoalEvaluationNotFound, wantStatus: http.StatusNotFound, wantCode: "EVALUATION_NOT_FOUND"},
		{name: "not owner", err: domain.ErrGoalEvaluationOwnedByOther, wantStatus: http.StatusConflict, wantCode: "EVALUATION_NOT_OWNER"},
		{name: "completed", err: domain.ErrGoalEvaluationNotDraft, wantStatus: http.StatusConflict, wantCode: "EVALUATION_ALREADY_COMPLETED"},
		{name: "historical", err: domain.ErrGoalEvaluationNotCurrentCycle, wantStatus: http.StatusConflict, wantCode: "EVALUATION_NOT_CURRENT_CYCLE"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(response)

			handleGoalEvaluationMutationError(ctx, test.err)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if !bytes.Contains(response.Body.Bytes(), []byte(`"code":"`+test.wantCode+`"`)) {
				t.Fatalf("body = %s, want code %s", response.Body.String(), test.wantCode)
			}
		})
	}
}

func TestCreateGoalEvaluationUsesStableConflictCode(t *testing.T) {
	gimMode := gin.Mode()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(gimMode) })

	service := &goalEvaluationServiceStub{createErr: domain.ErrGoalEvaluationOwnedByOther}
	router := gin.New()
	router.POST("/clients/:id/evaluations", evaluationActor(uuid.New()), NewClientHandler(service).CreateGoalEvaluation)
	request := httptest.NewRequest(http.MethodPost, "/clients/"+uuid.NewString()+"/evaluations", bytes.NewBufferString(`{"items":[]}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusConflict, response.Body.String())
	}
	if !bytes.Contains(response.Body.Bytes(), []byte(`"code":"EVALUATION_NOT_OWNER"`)) {
		t.Fatalf("body = %s, want stable ownership code", response.Body.String())
	}
}
