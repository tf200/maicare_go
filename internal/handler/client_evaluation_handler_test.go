package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"maicare_go/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type goalEvaluationServiceStub struct {
	domain.ClientService
	statsEmployeeID    uuid.UUID
	statsResult        *domain.EvaluationStats
	statsErr           error
	statsCalls         int
	createErr          error
	updateEvaluationID uuid.UUID
	updateEmployeeID   uuid.UUID
	updateParams       domain.UpdateGoalEvaluationDraftParams
	updateResult       *domain.GoalEvaluation
	updateErr          error
	updateCalls        int
	submitEvaluationID uuid.UUID
	submitEmployeeID   uuid.UUID
	submitParams       domain.SubmitGoalEvaluationDraftParams
	submitResult       *domain.GoalEvaluation
	submitErr          error
	submitCalls        int
}

func (s *goalEvaluationServiceStub) GetEvaluationStats(_ context.Context, employeeID uuid.UUID) (*domain.EvaluationStats, error) {
	s.statsCalls++
	s.statsEmployeeID = employeeID
	return s.statsResult, s.statsErr
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

func (s *goalEvaluationServiceStub) SubmitGoalEvaluationDraft(_ context.Context, evaluationID uuid.UUID, employeeID uuid.UUID, params domain.SubmitGoalEvaluationDraftParams) (*domain.GoalEvaluation, error) {
	s.submitCalls++
	s.submitEvaluationID = evaluationID
	s.submitEmployeeID = employeeID
	s.submitParams = params
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
	request.Header.Set("If-Match", `"2026-08-29T12:34:56.123456Z"`)
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
	if service.updateParams.ExpectedUpdatedAt.Format(time.RFC3339Nano) != "2026-08-29T12:34:56.123456Z" {
		t.Fatalf("revision = %s, want exact If-Match value", service.updateParams.ExpectedUpdatedAt.Format(time.RFC3339Nano))
	}
}

func TestUpdateGoalEvaluationDraftMapsOwnershipConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &goalEvaluationServiceStub{updateErr: domain.ErrGoalEvaluationOwnedByOther}
	router := gin.New()
	router.PATCH("/evaluations/:evaluation_id/draft", evaluationActor(uuid.New()), NewClientHandler(service).UpdateGoalEvaluationDraft)
	request := httptest.NewRequest(http.MethodPatch, "/evaluations/"+uuid.NewString()+"/draft", bytes.NewBufferString(`{"items":[]}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("If-Match", `"2026-08-29T12:34:56Z"`)
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
	request.Header.Set("If-Match", `"2026-08-29T12:34:56Z"`)
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
	request.Header.Set("If-Match", `"2026-08-29T12:34:56.123456Z"`)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}
	if service.submitCalls != 1 || service.submitEvaluationID != evaluationID || service.submitEmployeeID != employeeID {
		t.Fatal("handler did not pass exact evaluation and employee IDs")
	}
	if service.submitParams.ExpectedUpdatedAt.Format(time.RFC3339Nano) != "2026-08-29T12:34:56.123456Z" {
		t.Fatalf("submit revision = %s", service.submitParams.ExpectedUpdatedAt.Format(time.RFC3339Nano))
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
		{name: "client not in care", err: domain.ErrGoalEvaluationClientNotInCare, wantStatus: http.StatusUnprocessableEntity, wantCode: "EVALUATION_CLIENT_NOT_IN_CARE"},
		{name: "no active goals", err: domain.ErrGoalEvaluationNoActiveGoals, wantStatus: http.StatusUnprocessableEntity, wantCode: "EVALUATION_NO_ACTIVE_GOALS"},
		{name: "no due date", err: domain.ErrGoalEvaluationNoDueDate, wantStatus: http.StatusUnprocessableEntity, wantCode: "EVALUATION_NO_DUE_DATE"},
		{name: "duplicate goal", err: domain.ErrGoalEvaluationDuplicateGoal, wantStatus: http.StatusUnprocessableEntity, wantCode: "EVALUATION_DUPLICATE_GOAL"},
		{name: "inactive goal", err: domain.ErrGoalEvaluationGoalNotActive, wantStatus: http.StatusUnprocessableEntity, wantCode: "EVALUATION_GOAL_NOT_ACTIVE"},
		{name: "invalid progress", err: domain.ErrGoalEvaluationInvalidProgress, wantStatus: http.StatusUnprocessableEntity, wantCode: "EVALUATION_INVALID_PROGRESS"},
		{name: "incomplete", err: domain.ErrGoalEvaluationIncomplete, wantStatus: http.StatusUnprocessableEntity, wantCode: "EVALUATION_INCOMPLETE"},
		{name: "too early", err: domain.ErrGoalEvaluationTooEarly, wantStatus: http.StatusUnprocessableEntity, wantCode: "EVALUATION_TOO_EARLY"},
		{name: "conflict", err: domain.ErrGoalEvaluationConflict, wantStatus: http.StatusConflict, wantCode: "EVALUATION_CONFLICT"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(response)

			handleGoalEvaluationMutationError(ctx, test.err, nil, false)

			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, test.wantStatus)
			}
			if !bytes.Contains(response.Body.Bytes(), []byte(`"code":"`+test.wantCode+`"`)) {
				t.Fatalf("body = %s, want code %s", response.Body.String(), test.wantCode)
			}
		})
	}
}

func TestSubmitGoalEvaluationDraftReturnsSavedDraftOnValidationFailure(t *testing.T) {
	evaluationID := uuid.New()
	result := &domain.GoalEvaluation{ID: evaluationID, ClientID: uuid.New(), Status: "draft"}
	service := &goalEvaluationServiceStub{submitResult: result, submitErr: domain.ErrGoalEvaluationIncomplete}
	router := gin.New()
	router.POST("/evaluations/:evaluation_id/submit", evaluationActor(uuid.New()), NewClientHandler(service).SubmitGoalEvaluationDraft)
	response := httptest.NewRecorder()

	request := httptest.NewRequest(http.MethodPost, "/evaluations/"+evaluationID.String()+"/submit", nil)
	request.Header.Set("If-Match", `"2026-08-29T12:34:56Z"`)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusUnprocessableEntity, response.Body.String())
	}
	for _, expected := range []string{`"code":"EVALUATION_INCOMPLETE"`, `"draft_saved":true`, `"id":"` + evaluationID.String() + `"`} {
		if !bytes.Contains(response.Body.Bytes(), []byte(expected)) {
			t.Fatalf("body = %s, want %s", response.Body.String(), expected)
		}
	}
}

func TestGoalEvaluationMutationRequiresRevision(t *testing.T) {
	service := &goalEvaluationServiceStub{}
	router := gin.New()
	router.POST("/evaluations/:evaluation_id/submit", evaluationActor(uuid.New()), NewClientHandler(service).SubmitGoalEvaluationDraft)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/evaluations/"+uuid.NewString()+"/submit", nil))

	if response.Code != http.StatusPreconditionRequired || service.submitCalls != 0 {
		t.Fatalf("status = %d, calls = %d; want 428 and no service call", response.Code, service.submitCalls)
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

func TestGetEvaluationStatsUsesAuthenticatedEmployee(t *testing.T) {
	gin.SetMode(gin.TestMode)
	employeeID := uuid.New()
	service := &goalEvaluationServiceStub{statsResult: &domain.EvaluationStats{AttentionRequired: 4, InProgress: 2, RecentlyFinalized: 1, AsOf: time.Date(2026, 8, 30, 12, 0, 0, 0, time.UTC)}}
	router := gin.New()
	router.GET("/evaluations/stats", evaluationActor(employeeID), NewClientHandler(service).GetEvaluationStats)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/evaluations/stats", nil))

	if response.Code != http.StatusOK || service.statsCalls != 1 || service.statsEmployeeID != employeeID {
		t.Fatalf("status=%d calls=%d employee=%s body=%s", response.Code, service.statsCalls, service.statsEmployeeID, response.Body.String())
	}
	for _, expected := range []string{`"attention_required":4`, `"in_progress":2`, `"recently_finalized":1`, `"as_of":"2026-08-30T12:00:00Z"`} {
		if !bytes.Contains(response.Body.Bytes(), []byte(expected)) {
			t.Fatalf("body=%s, want %s", response.Body.String(), expected)
		}
	}
}

func TestGoalEvaluationResponseUsesDateOnlyCalendarFields(t *testing.T) {
	date := time.Date(2026, time.March, 29, 0, 0, 0, 0, time.UTC)
	response := toGoalEvaluationResponse(domain.GoalEvaluation{
		EvaluationDate: date,
		PeriodStart:    &date,
		PeriodEnd:      &date,
		Items:          []domain.GoalEvaluationItem{},
	})

	body, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal response: %v", err)
	}
	for _, expected := range []string{
		`"evaluation_date":"2026-03-29"`,
		`"period_start":"2026-03-29"`,
		`"period_end":"2026-03-29"`,
	} {
		if !bytes.Contains(body, []byte(expected)) {
			t.Fatalf("body=%s, want %s", body, expected)
		}
	}
}
