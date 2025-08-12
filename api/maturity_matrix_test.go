package api

import (
	"bytes"
	"context"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/token"
	"maicare_go/util"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

type RandomCarePlan struct {
	CarePlanID       int64 `json:"care_plan_id"`
	ObjectiveID      int64 `json:"objective_id"`
	ActionID         int64 `json:"action_id"`
	InterventionID   int64 `json:"intervention_id"`
	SuccessMetricID  int64 `json:"success_metric_id"`
	RiskID           int64 `json:"risk_id"`
	SupportNetworkID int64 `json:"support_network_id"`
	ResourceID       int64 `json:"resource_id"`
}

func createRandomCarePlan(t *testing.T, clientID int64) RandomCarePlan {
	clientAssessments, err := testStore.CreateClientMaturityMatrixAssessment(context.Background(), db.CreateClientMaturityMatrixAssessmentParams{
		ClientID:         clientID,
		MaturityMatrixID: 1,
		InitialLevel:     1,
		CurrentLevel:     1,
		TargetLevel:      2,
		StartDate:        pgtype.Date{Time: time.Now(), Valid: true},
		EndDate:          pgtype.Date{Time: time.Now().Add(time.Hour * 24 * 365), Valid: true},
	})
	require.NoError(t, err)
	mockLLmResp := CreateMockGrpcClient().Response
	mockLLmRespBytes, err := json.Marshal(mockLLmResp)
	require.NoError(t, err)
	require.NoError(t, err)
	carePlan, err := testStore.CreateCarePlan(context.Background(), db.CreateCarePlanParams{
		AssessmentID:          clientAssessments.ID,
		GeneratedByEmployeeID: nil,
		AssessmentSummary:     mockLLmResp.AssessmentSummary,
		RawLlmResponse:        mockLLmRespBytes,
		Status:                "draft",
	})
	require.NoError(t, err)
	var randomObjectiveID int64
	var randomActionID int64
	for _, objective := range mockLLmResp.CarePlanObjectives.ShortTermGoals {
		createdObj, err := testStore.CreateCarePlanObjective(context.Background(), db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlan.ID,
			Description: objective.Description,
			Timeframe:   "short_term",
			TargetDate:  pgtype.Date{Time: time.Now(), Valid: true}, // Use current time as target date
		})
		randomObjectiveID = createdObj.ID
		require.NoError(t, err)
		for i, action := range objective.SpecificActions {

			action, err := testStore.CreateCarePlanAction(context.Background(), db.CreateCarePlanActionParams{
				ObjectiveID:       createdObj.ID,
				ActionDescription: action,
				SortOrder:         int32(i + 1), // Use index + 1 as sort order
			})
			randomActionID = action.ID
			require.NoError(t, err)

		}
	}

	for _, objective := range mockLLmResp.CarePlanObjectives.MediumTermGoals {
		createdObj, err := testStore.CreateCarePlanObjective(context.Background(), db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlan.ID,
			Description: objective.Description,
			Timeframe:   "medium_term",
			TargetDate:  pgtype.Date{Time: time.Now().Add(time.Hour * 24 * 30), Valid: true}, // Use current time + 30 days as target date
		})
		require.NoError(t, err)
		for i, action := range objective.SpecificActions {

			_, err := testStore.CreateCarePlanAction(context.Background(), db.CreateCarePlanActionParams{
				ObjectiveID:       createdObj.ID,
				ActionDescription: action,
				SortOrder:         int32(i + 1), // Use index + 1 as sort order
			})
			require.NoError(t, err)

		}
	}

	for _, objective := range mockLLmResp.CarePlanObjectives.LongTermGoals {
		createdObj, err := testStore.CreateCarePlanObjective(context.Background(), db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlan.ID,
			Description: objective.Description,
			Timeframe:   "long_term",
			TargetDate:  pgtype.Date{Time: time.Now().Add(time.Hour * 24 * 365), Valid: true}, // Use current time + 1 year as target date
		})
		require.NoError(t, err)
		for i, action := range objective.SpecificActions {

			_, err := testStore.CreateCarePlanAction(context.Background(), db.CreateCarePlanActionParams{
				ObjectiveID:       createdObj.ID,
				ActionDescription: action,
				SortOrder:         int32(i + 1), // Use index + 1 as sort order
			})
			require.NoError(t, err)

		}
	}
	var interventionID int64
	for _, intervention := range mockLLmResp.Interventions.DailyActivities {
		intervention, err := testStore.CreateCarePlanIntervention(context.Background(), db.CreateCarePlanInterventionParams{
			CarePlanID:              carePlan.ID,
			Frequency:               "daily",
			InterventionDescription: intervention,
		})
		interventionID = intervention.ID
		require.NoError(t, err)
	}
	for _, intervention := range mockLLmResp.Interventions.WeeklyActivities {
		_, err := testStore.CreateCarePlanIntervention(context.Background(), db.CreateCarePlanInterventionParams{
			CarePlanID:              carePlan.ID,
			Frequency:               "weekly",
			InterventionDescription: intervention,
		})
		require.NoError(t, err)
	}
	var successMetricID int64
	for _, successMetric := range mockLLmResp.SuccessMetrics {
		successMetric, err := testStore.CreateCarePlanSuccessMetric(context.Background(), db.CreateCarePlanSuccessMetricParams{
			CarePlanID:        carePlan.ID,
			MetricName:        successMetric.Metric,
			TargetValue:       successMetric.Target,
			MeasurementMethod: successMetric.MeasurementMethod,
		})
		successMetricID = successMetric.ID
		require.NoError(t, err)
	}

	var randomRiskID int64
	for _, risk := range mockLLmResp.RiskFactors {
		risk, err := testStore.CreateCarePlanRisk(context.Background(), db.CreateCarePlanRiskParams{
			CarePlanID:         carePlan.ID,
			RiskDescription:    risk.Risk,
			MitigationStrategy: risk.Mitigation,
			RiskLevel:          nil, // TODO: Add risk level if needed
		})
		randomRiskID = risk.ID
		require.NoError(t, err)
	}
	var supportNetworkID int64
	for _, supportNetwork := range mockLLmResp.SupportNetwork {
		network, err := testStore.CreateCarePlanSupportNetwork(context.Background(), db.CreateCarePlanSupportNetworkParams{
			CarePlanID:                carePlan.ID,
			RoleTitle:                 supportNetwork.Role,
			ResponsibilityDescription: supportNetwork.Responsibility,
			ContactPerson:             nil, // TODO: Add contact person if needed
			ContactDetails:            nil, // TODO: Add contact details if needed
		})
		supportNetworkID = network.ID
		require.NoError(t, err)
	}
	var resourceID int64
	for _, resource := range mockLLmResp.ResourcesRequired {
		resource, err := testStore.CreateCarePlanResources(context.Background(), db.CreateCarePlanResourcesParams{
			CarePlanID:          carePlan.ID,
			ResourceDescription: resource,
			IsObtained:          false,
			ObtainedDate:        pgtype.Date{Time: time.Now(), Valid: false},
		})
		resourceID = resource.ID
		require.NoError(t, err)
	}

	rnd := RandomCarePlan{
		CarePlanID:       carePlan.ID,
		ObjectiveID:      randomObjectiveID,
		ActionID:         randomActionID,
		InterventionID:   interventionID,
		SuccessMetricID:  successMetricID,
		RiskID:           randomRiskID,
		SupportNetworkID: supportNetworkID,
		ResourceID:       resourceID,
	}

	return rnd
}

func TestCreateClientMaturityMatrixAssessmentApi(t *testing.T) {
	client := createRandomClientDetails(t)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				assessmentReq := CreateClientMaturityMatrixAssessmentRequest{
					MaturityMatrixID: 1,
					InitialLevel:     1,
					TargetLevel:      3,
				}
				data, err := json.Marshal(assessmentReq)
				require.NoError(t, err)

				url := fmt.Sprintf("/clients/%d/maturity_matrix_assessment", client.ID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil

			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log(recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var assessmentCard Response[CreateClientMaturityMatrixAssessmentResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &assessmentCard)
				require.NoError(t, err)
				require.NotEmpty(t, assessmentCard.Data)
				require.Equal(t, client.ID, assessmentCard.Data.ClientID)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request, err := tc.buildRequest()
			require.NoError(t, err)

			tc.setupAuth(t, request, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func TestGetCarePlanOverviewApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/care_plans/%d", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[GetCarePlanOverviewResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.CarePlanID, response.Data.ID)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}
func TestUpdateCarePlanOverviewApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				updateReq := UpdateCarePlanOverviewRequest{
					AssessmentSummary: util.StringPtr("Updated assessment summary"),
				}
				data, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/care_plans/%d", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateCarePlanOverviewResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.CarePlanID, response.Data.CarePlanID)
				require.Equal(t, "Updated assessment summary", response.Data.AssessmentSummary)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCarePlanApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/care_plans/%d", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateCarePlanObjectiveApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createReq := CreateCarePlanObjectiveRequest{
					TimeFrame:   "short_term",
					GoalTitle:   "New Objective",
					Description: "This is a new objective for the care plan.",
				}
				data, err := json.Marshal(createReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/care_plans/%d/objectives", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var response Response[CreateCarePlanObjectiveResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.CarePlanID, response.Data.CarePlanID)
				require.Equal(t, "New Objective", response.Data.GoalTitle)
				require.Equal(t, "This is a new objective for the care plan.", response.Data.Description)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetCarePlanObjectivesApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlanID := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/care_plans/%d/objectives", carePlanID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[GetCarePlanObjectivesResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)

				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.NotEmpty(t, response.Data.ShortTermGoals)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCarePlanObjectiveApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				updateReq := UpdateCarePlanObjectiveRequest{
					TimeFrame:   util.StringPtr("short_term"),
					GoalTitle:   util.StringPtr("Updated Objective"),
					Description: util.StringPtr("This is an updated objective for the care plan."),
					Status:      util.StringPtr("not_started"),
				}
				data, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/objectives/%d", carePlan.ObjectiveID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateCarePlanObjectiveResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.CarePlanID, response.Data.CarePlanId)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCarePlanObjectiveApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/objectives/%d", carePlan.ObjectiveID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateCarePlanActionsApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createReq := CreateCarePlanActionsRequest{
					ActionDescription: "New action for care plan objective",
				}
				data, err := json.Marshal(createReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/objectives/%d/actions", carePlan.ObjectiveID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var response Response[CreateCarePlanActionsResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.ObjectiveID, response.Data.ObjectiveID)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCarePlanActionsApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				updateReq := UpdateCarePlanActionsRequest{
					ActionDescription: util.StringPtr("Updated action description for care plan objective"),
				}
				data, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/actions/%d", carePlan.ActionID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateCarePlanActionsResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.ObjectiveID, response.Data.ObjectiveID)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCarePlanActionApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/actions/%d", carePlan.ActionID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateCarePlanInterventionApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createReq := CreateCarePlanInterventionRequest{
					Frequency:               "daily",
					InterventionDescription: "New daily intervention for care plan",
				}
				data, err := json.Marshal(createReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/care_plans/%d/interventions", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var response Response[CreateCarePlanInterventionResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.CarePlanID, response.Data.CarePlanID)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetCarePlanInterventionsApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/care_plans/%d/interventions", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[GetCarePlanInterventionsResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.NotEmpty(t, response.Data.DailyActivities)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCarePlanInterventionApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				updateReq := UpdateCarePlanInterventionRequest{
					Frequency:               util.StringPtr("weekly"),
					InterventionDescription: util.StringPtr("Updated intervention description for care plan"),
				}
				data, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/interventions/%d", carePlan.InterventionID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateCarePlanInterventionResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.CarePlanID, response.Data.CarePlanID)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCarePlanInterventionApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/interventions/%d", carePlan.InterventionID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateCarePlanSuccessMetricApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createReq := CreateCarePlanSuccessMetricsRequest{
					MetricName:        "Weight Loss",
					TargetValue:       "10",
					MeasurementMethod: "kg",
				}
				data, err := json.Marshal(createReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/care_plans/%d/success_metrics", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var response Response[CreateCarePlanSuccessMetricsResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.SuccessMetricID, response.Data.MetricID)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetCarePlanSuccessMetricsApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/care_plans/%d/success_metrics", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[[]GetCarePlanSuccessMetricsResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCarePlanSuccessMetricApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				updateReq := UpdateCarePlanSuccessMetricsRequest{
					MetricName:        util.StringPtr("Updated Weight Loss"),
					TargetValue:       util.StringPtr("15"),
					MeasurementMethod: util.StringPtr("kg"),
				}
				data, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/success_metrics/%d", carePlan.SuccessMetricID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateCarePlanSuccessMetricsResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.SuccessMetricID, response.Data.MetricID)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCarePlanSuccessMetricApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/success_metrics/%d", carePlan.SuccessMetricID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateCarePlanRisksApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createReq := CreateCarePlanRisksRequest{
					RiskDescription:    "High risk of falls",
					RiskLevel:          util.StringPtr("high"),
					MitigationStrategy: "Implement fall prevention measures",
				}
				data, err := json.Marshal(createReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/care_plans/%d/risks", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var response Response[CreateCarePlanRisksResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetCarePlanRisksApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/care_plans/%d/risks", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[[]GetCarePlanRisksResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCarePlanRiskApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				updateReq := UpdateCarePlanRisksRequest{
					RiskDescription:    util.StringPtr("Updated risk description for care plan"),
					RiskLevel:          util.StringPtr("medium"),
					MitigationStrategy: util.StringPtr("Implement updated mitigation strategy"),
				}
				data, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/risks/%d", carePlan.RiskID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateCarePlanRisksResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.RiskID, response.Data.RiskID)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCarePlanRiskApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/risks/%d", carePlan.RiskID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateCarePlanSupportNetworkApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createReq := CreateCarePlanSupportNetworkRequest{
					RoleTitle:                 "Caregiver",
					ResponsibilityDescription: "Assist with daily activities and provide emotional support.",
				}
				data, err := json.Marshal(createReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/care_plans/%d/support_network", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var response Response[CreateCarePlanSupportNetworkResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetCarePlanSupportNetworkApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/care_plans/%d/support_network", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[[]GetCarePlanSupportNetworkResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCarePlanSupportNetworkApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				updateReq := UpdateCarePlanSupportNetworkRequest{
					RoleTitle:                 util.StringPtr("Updated Caregiver"),
					ResponsibilityDescription: util.StringPtr("Updated responsibilities for caregiver."),
				}
				data, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/support_network/%d", carePlan.SupportNetworkID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateCarePlanSupportNetworkResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
				require.Equal(t, carePlan.SupportNetworkID, response.Data.SupportNetworkID)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCarePlanSupportNetworkApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/support_network/%d", carePlan.SupportNetworkID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestCreateCarePlanResourcesApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				createReq := CreateCarePlanResourcesRequest{

					ResourceDescription: "A comprehensive guide to nutrition for better health.",
					ResourceType:        util.StringPtr("Nutrition Guide"),
					IsObtained:          util.BoolPtr(true),
					ObtainedDate:        util.TimePtr(time.Now()),
				}
				data, err := json.Marshal(createReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/care_plans/%d/resources", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusCreated, recorder.Code)
				var response Response[CreateCarePlanResourcesResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetCarePlanResourcesApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/care_plans/%d/resources", carePlan.CarePlanID)
				req, err := http.NewRequest(http.MethodGet, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[[]GetCarePlanResourcesResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestUpdateCarePlanResourceApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				updateReq := UpdateCarePlanResourcesRequest{
					ResourceDescription: util.StringPtr("Updated resource description for care plan"),
					ResourceType:        util.StringPtr("Updated Resource Type"),
					IsObtained:          util.BoolPtr(false),
					ObtainedDate:        time.Now().AddDate(0, 0, 1),
				}
				data, err := json.Marshal(updateReq)
				require.NoError(t, err)
				url := fmt.Sprintf("/resources/%d", carePlan.ResourceID)
				req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
				var response Response[UpdateCarePlanResourcesResponse]
				err := json.Unmarshal(recorder.Body.Bytes(), &response)
				require.NoError(t, err)
				require.NotEmpty(t, response.Data)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)
			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}

func TestDeleteCarePlanResourceApi(t *testing.T) {
	client := createRandomClientDetails(t)
	carePlan := createRandomCarePlan(t, client.ID)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildRequest  func() (*http.Request, error)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, 1, time.Minute)
			},
			buildRequest: func() (*http.Request, error) {
				url := fmt.Sprintf("/resources/%d", carePlan.ResourceID)
				req, err := http.NewRequest(http.MethodDelete, url, nil)
				require.NoError(t, err)
				return req, nil
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				t.Log("Response Body:", recorder.Body.String())
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.buildRequest()
			require.NoError(t, err)

			recorder := httptest.NewRecorder()
			tc.setupAuth(t, req, testServer.tokenMaker)
			testServer.router.ServeHTTP(recorder, req)
			tc.checkResponse(recorder)
		})
	}
}
