package api

import (
	"fmt"
	db "maicare_go/db/sqlc"
	grpclient "maicare_go/grpclient/proto"
	"maicare_go/pagination"
	"maicare_go/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"
)

// ListMaturityMatrixResponse represents a maturity matrix in the list
type ListMaturityMatrixResponse struct {
	ID        int64  `json:"id"`
	TopicName string `json:"topic_name"`
}

// @Summary List all maturity matrix
// @Description Get a list of all maturity matrix
// @Tags maturity_matrix
// @Produce json
// @Success 200 {object} Response[[]ListMaturityMatrixResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /maturity_matrix [get]
func (server *Server) ListMaturityMatrixApi(ctx *gin.Context) {
	maturityMatrix, err := server.store.ListMaturityMatrix(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	responseMaturityMatrix := make([]ListMaturityMatrixResponse, len(maturityMatrix))
	for i, matrix := range maturityMatrix {
		responseMaturityMatrix[i] = ListMaturityMatrixResponse{
			ID:        matrix.ID,
			TopicName: matrix.TopicName,
		}
	}

	res := SuccessResponse(responseMaturityMatrix, "Maturity matrix retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

type Level struct {
	Level       int32  `json:"level"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateClientMaturityMatrixAssessmentRequest represents a request to create a client maturity matrix assessment
type CreateClientMaturityMatrixAssessmentRequest struct {
	MaturityMatrixID int64 `json:"maturity_matrix_id"`
	InitialLevel     int32 `json:"initial_level"`
	TargetLevel      int32 `json:"target_level"`
}

// CreateClientMaturityMatrixAssessmentResponse represents a response for CreateClientMaturityMatrixAssessmentApi
type CreateClientMaturityMatrixAssessmentResponse struct {
	ClientID   int64 `json:"client_id"`
	CarePlanID int64 `json:"care_plan_id"`
}

// @Summary Create client maturity matrix assessment
// @Description Create a client maturity matrix assessment
// @Tags maturity_matrix
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body CreateClientMaturityMatrixAssessmentRequest true "Request body"
// @Success 201 {object} Response[CreateClientMaturityMatrixAssessmentResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/assessments [post]
func (server *Server) CreateClientMaturityMatrixAssessmentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req CreateClientMaturityMatrixAssessmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := server.store.WithTx(tx)

	employeeID, err := qtx.GetEmployeeIDByUserID(ctx, payload.UserId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateClientMaturityMatrixAssessmentParams{
		ClientID:         clientID,
		MaturityMatrixID: req.MaturityMatrixID,
		StartDate:        pgtype.Date{Time: time.Now(), Valid: true},
		EndDate:          pgtype.Date{Time: time.Now().Add(time.Hour * 24 * 365), Valid: true},
		InitialLevel:     req.InitialLevel,
		TargetLevel:      req.TargetLevel,
		CurrentLevel:     req.InitialLevel,
	}

	clientAssessments, err := qtx.CreateClientMaturityMatrixAssessment(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	topicDescription, err := qtx.GetMaturityMatrix(ctx, req.MaturityMatrixID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	var levelDescription []Level

	err = json.Unmarshal(topicDescription.LevelDescription, &levelDescription)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	clientDetails, err := qtx.GetClientDetails(ctx, clientID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return

	}

	generatedCarePlan, err := server.grpClient.GenerateCarePlan(ctx, &grpclient.PersonalizedCarePlanRequest{
		ClientData: &grpclient.ClientData{
			Age:              int32(time.Since(clientDetails.DateOfBirth.Time).Hours() / 24 / 365), // Calculate age from DateOfBirth
			LivingSituation:  *clientDetails.LivingSituation,
			EducationLevel:   *clientDetails.EducationLevel,
			DomainName:       clientAssessments.TopicName,
			CurrentLevel:     clientAssessments.InitialLevel,                       // Example current level, replace with actual data
			LevelDescription: levelDescription[req.MaturityMatrixID-1].Description, // Use the description from the level
		},
		DomainDefinitions: map[string]*grpclient.DomainLevels{
			topicDescription.TopicName: {
				Levels: map[int32]string{
					levelDescription[0].Level: levelDescription[0].Description,
					levelDescription[1].Level: levelDescription[1].Description,
					levelDescription[2].Level: levelDescription[2].Description,
					levelDescription[3].Level: levelDescription[3].Description,
					levelDescription[4].Level: levelDescription[4].Description,
				},
			},
		},
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// insert the care plan into the database
	rawllmResp, err := json.Marshal(generatedCarePlan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	carePlan, err := qtx.CreateCarePlan(ctx, db.CreateCarePlanParams{
		AssessmentID:          clientAssessments.ID,
		GeneratedByEmployeeID: &employeeID,
		AssessmentSummary:     generatedCarePlan.AssessmentSummary,
		RawLlmResponse:        rawllmResp,
		Status:                "draft",
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// insert Objectives into the database
	// starting with short term goals

	for _, objective := range generatedCarePlan.CarePlanObjectives.ShortTermGoals {
		createdObj, err := qtx.CreateCarePlanObjective(ctx, db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlan.ID,
			GoalTitle:   objective.GoalTitle,
			Description: objective.Description,
			Timeframe:   "short_term",
			TargetDate:  pgtype.Date{Time: time.Now(), Valid: true}, // Use current time as target date
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		for i, action := range objective.SpecificActions {

			_, err := qtx.CreateCarePlanAction(ctx, db.CreateCarePlanActionParams{
				ObjectiveID:       createdObj.ID,
				ActionDescription: action,
				SortOrder:         int32(i + 1), // Use index + 1 as sort order
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}

		}
	}

	for _, objective := range generatedCarePlan.CarePlanObjectives.MediumTermGoals {
		createdObj, err := qtx.CreateCarePlanObjective(ctx, db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlan.ID,
			GoalTitle:   objective.GoalTitle,
			Description: objective.Description,
			Timeframe:   "medium_term",
			TargetDate:  pgtype.Date{Time: time.Now(), Valid: true}, // Use current time as target date
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		for i, action := range objective.SpecificActions {
			_, err := qtx.CreateCarePlanAction(ctx, db.CreateCarePlanActionParams{
				ObjectiveID:       createdObj.ID,
				ActionDescription: action,
				SortOrder:         int32(i + 1), // Use index + 1 as sort order
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
	}

	for _, objective := range generatedCarePlan.CarePlanObjectives.LongTermGoals {
		createdObj, err := qtx.CreateCarePlanObjective(ctx, db.CreateCarePlanObjectiveParams{
			CarePlanID:  carePlan.ID,
			GoalTitle:   objective.GoalTitle,
			Description: objective.Description,
			Timeframe:   "long_term",
			TargetDate:  pgtype.Date{Time: time.Now(), Valid: true}, // Use current time as target date
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		for i, action := range objective.SpecificActions {
			_, err := qtx.CreateCarePlanAction(ctx, db.CreateCarePlanActionParams{
				ObjectiveID:       createdObj.ID,
				ActionDescription: action,
				SortOrder:         int32(i + 1), // Use index + 1 as sort order
			})
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, errorResponse(err))
				return
			}
		}
	}
	// insert interventions into the database

	for _, intervention := range generatedCarePlan.Interventions.DailyActivities {
		_, err := qtx.CreateCarePlanIntervention(ctx, db.CreateCarePlanInterventionParams{
			CarePlanID:              carePlan.ID,
			Frequency:               "daily",
			InterventionDescription: intervention,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}
	for _, intervention := range generatedCarePlan.Interventions.WeeklyActivities {
		_, err := qtx.CreateCarePlanIntervention(ctx, db.CreateCarePlanInterventionParams{
			CarePlanID:              carePlan.ID,
			Frequency:               "weekly",
			InterventionDescription: intervention,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	for _, intervention := range generatedCarePlan.Interventions.MonthlyActivities {
		_, err := qtx.CreateCarePlanIntervention(ctx, db.CreateCarePlanInterventionParams{
			CarePlanID:              carePlan.ID,
			Frequency:               "monthly",
			InterventionDescription: intervention,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	for _, successMetric := range generatedCarePlan.SuccessMetrics {
		_, err := qtx.CreateCarePlanSuccessMetric(ctx, db.CreateCarePlanSuccessMetricParams{
			CarePlanID:        carePlan.ID,
			MetricName:        successMetric.Metric,
			TargetValue:       successMetric.Target,
			MeasurementMethod: successMetric.MeasurementMethod,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	for _, risk := range generatedCarePlan.RiskFactors {
		_, err := qtx.CreateCarePlanRisk(ctx, db.CreateCarePlanRiskParams{
			CarePlanID:         carePlan.ID,
			RiskDescription:    risk.Risk,
			MitigationStrategy: risk.Mitigation,
			RiskLevel:          &risk.RiskLevel, // Use pointer to allow NULL values
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	for _, supportNetwork := range generatedCarePlan.SupportNetwork {
		_, err := qtx.CreateCarePlanSupportNetwork(ctx, db.CreateCarePlanSupportNetworkParams{
			CarePlanID:                carePlan.ID,
			RoleTitle:                 supportNetwork.Role,
			ResponsibilityDescription: supportNetwork.Responsibility,
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	for _, resource := range generatedCarePlan.ResourcesRequired {
		_, err := qtx.CreateCarePlanResources(ctx, db.CreateCarePlanResourcesParams{
			CarePlanID:          carePlan.ID,
			ResourceDescription: resource,
			IsObtained:          false,
			ObtainedDate:        pgtype.Date{Time: time.Now(), Valid: false},
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateClientMaturityMatrixAssessmentResponse{
		ClientID:   clientID,
		CarePlanID: carePlan.ID,
	}, "Client maturity matrix assessment created successfully")
	ctx.JSON(http.StatusCreated, res)
}

// ListClientMaturityMatrixAssessmentsRequest represents a request to list client maturity matrix assessments
type ListClientMaturityMatrixAssessmentsRequest struct {
	pagination.Request
}

// ListClientMaturityMatrixAssessmentsResponse represents a response for ListClientMaturityMatrixAssessmentsApi
type ListClientMaturityMatrixAssessmentsResponse struct {
	CarePlanID   int64       `json:"care_plan_id"`
	ClientID     int64       `json:"client_id"`
	StartDate    pgtype.Date `json:"start_date"`
	EndDate      pgtype.Date `json:"end_date"`
	InitialLevel int32       `json:"initial_level"`
	CurrentLevel int32       `json:"current_level"`
	IsActive     bool        `json:"is_active"`
	TopicName    string      `json:"topic_name"`
}

// @Summary List client maturity matrix assessments
// @Description Get a list of client maturity matrix assessments
// @Tags care_plan
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListClientMaturityMatrixAssessmentsResponse]]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/assessments [get]
func (server *Server) ListClientMaturityMatrixAssessmentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListClientMaturityMatrixAssessmentsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	clientAssessments, err := server.store.ListClientMaturityMatrixAssessments(ctx, db.ListClientMaturityMatrixAssessmentsParams{
		ClientID: clientID,
		Limit:    params.Limit,
		Offset:   params.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(clientAssessments) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListClientMaturityMatrixAssessmentsResponse{}, 0)
		res := SuccessResponse(pag, "No client maturity matrix assessments found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	responseClientAssessments := make([]ListClientMaturityMatrixAssessmentsResponse, len(clientAssessments))
	for i, assessment := range clientAssessments {
		responseClientAssessments[i] = ListClientMaturityMatrixAssessmentsResponse{
			CarePlanID:   util.DerefInt64(assessment.CarePlanID),
			TopicName:    assessment.TopicName,
			ClientID:     assessment.ClientID,
			StartDate:    assessment.StartDate,
			EndDate:      assessment.EndDate,
			InitialLevel: assessment.InitialLevel,
			CurrentLevel: assessment.CurrentLevel,
			IsActive:     assessment.IsActive,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, responseClientAssessments, clientAssessments[0].TotalCount)
	res := SuccessResponse(pag, "Client maturity matrix assessments retrieved successfully")

	ctx.JSON(http.StatusOK, res)

}

// ============================ CarePlan Overview ===========================

// care_plan represents the response for the GetCarePlanOverview API
type GetCarePlanOverviewResponse struct {
	ID                int64     `json:"id"`
	Domain            string    `json:"domain"`
	CurrentLevel      int32     `json:"current_level"`
	TargetLevel       int32     `json:"target_level"`
	Status            string    `json:"status"`
	GeneratedAt       time.Time `json:"generated_at"`
	AssessmentSummary string    `json:"assessment_summary"`
	RawLlmResponse    string    `json:"raw_llm_response"`
}

// GetCarePlanOverviewApi retrieves the care plan overview for a given assessment ID
// @Summary Get care plan overview
// @Description Get the care plan overview for a given assessment ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[GetCarePlanOverviewResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id} [get]
func (server *Server) GetCarePlanOverviewApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanOverviewApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	carePlan, err := server.store.GetCarePlanOverview(ctx, carePlanID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanOverviewApi", "Failed to get care plan overview", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan overview")))
		return
	}

	res := SuccessResponse(GetCarePlanOverviewResponse{
		ID:                carePlan.ID,
		Domain:            carePlan.TopicName,
		CurrentLevel:      carePlan.CurrentLevel,
		TargetLevel:       carePlan.TargetLevel,
		Status:            carePlan.Status,
		GeneratedAt:       carePlan.GeneratedAt.Time,
		AssessmentSummary: carePlan.AssessmentSummary,
		RawLlmResponse:    string(carePlan.RawLlmResponse),
	}, "Care plan overview retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanOverviewRequest represents the request body for updating the care plan overview
type UpdateCarePlanOverviewRequest struct {
	AssessmentSummary *string `json:"assessment_summary"`
}

// UpdateCarePlanOverviewResponse represents the response for the UpdateCarePlanOverview API
type UpdateCarePlanOverviewResponse struct {
	CarePlanID        int64  `json:"care_plan_id"`
	AssessmentID      int64  `json:"assessment_id"`
	AssessmentSummary string `json:"assessment_summary"`
}

// UpdateCarePlanOverviewApi updates the care plan overview for a given care plan ID
// @Summary Update care plan overview
// @Description Update the care plan overview for a given care plan ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body UpdateCarePlanOverviewRequest true "Request body"
// @Success 200 {object} Response[UpdateCarePlanOverviewResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id} [put]
func (server *Server) UpdateCarePlanOverviewApi(ctx *gin.Context) {
	CarePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanOverviewApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req UpdateCarePlanOverviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanOverviewApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	carePlan, err := server.store.UpdateCarePlanOverview(ctx, db.UpdateCarePlanOverviewParams{
		ID:                CarePlanID,
		AssessmentSummary: req.AssessmentSummary,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanOverviewApi", "Failed to update care plan overview", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan overview")))
		return
	}

	res := SuccessResponse(UpdateCarePlanOverviewResponse{
		CarePlanID:        carePlan.ID,
		AssessmentID:      carePlan.AssessmentID,
		AssessmentSummary: carePlan.AssessmentSummary,
	}, "Care plan overview updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteCarePlanApi deletes a care plan by its ID
// @Summary Delete care plan
// @Description Delete a care plan by its ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plan/{care_plan_id} [delete]
func (server *Server) DeleteCarePlanApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	err = server.store.DeleteCarePlan(ctx, carePlanID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanApi", "Failed to delete care plan", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// =========================== Care plan objectives and actions ===========================

// CreateCarePlanObjectiveRequest represents the request body for creating a care plan objective
type CreateCarePlanObjectiveRequest struct {
	TimeFrame   string `json:"timeframe" binding:"required,oneof=short_term medium_term long_term"`
	GoalTitle   string `json:"goal_title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// CreateCarePlanObjectiveResponse represents the response for the CreateCarePlanObjective API
type CreateCarePlanObjectiveResponse struct {
	ID              int64     `json:"id"`
	CarePlanID      int64     `json:"care_plan_id"`
	Timeframe       string    `json:"timeframe"`
	GoalTitle       string    `json:"goal_title"`
	Description     string    `json:"description"`
	TargetDate      time.Time `json:"target_date"`
	Status          string    `json:"status"`
	CompletionDate  time.Time `json:"completion_date"`
	CompletionNotes *string   `json:"completion_notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// CreateCarePlanObjectiveApi creates a new care plan objective
// @Summary Create care plan objective
// @Description Create a new care plan objective
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body CreateCarePlanObjectiveRequest true "Request body"
// @Success 201 {object} Response[CreateCarePlanObjectiveResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/objectives [post]
func (server *Server) CreateCarePlanObjectiveApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanObjectiveApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req CreateCarePlanObjectiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanObjectiveApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	objective, err := server.store.CreateCarePlanObjective(ctx, db.CreateCarePlanObjectiveParams{
		CarePlanID:  carePlanID,
		Description: req.Description,
		Timeframe:   req.TimeFrame,
		GoalTitle:   req.GoalTitle,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanObjectiveApi", "Failed to create care plan objective", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan objective")))
		return
	}

	res := SuccessResponse(CreateCarePlanObjectiveResponse{
		ID:              objective.ID,
		CarePlanID:      carePlanID,
		Timeframe:       objective.Timeframe,
		GoalTitle:       objective.GoalTitle,
		Description:     objective.Description,
		TargetDate:      objective.TargetDate.Time,
		Status:          objective.Status,
		CompletionDate:  objective.CompletionDate.Time,
		CompletionNotes: objective.CompletionNotes,
		CreatedAt:       objective.CreatedAt.Time,
		UpdatedAt:       objective.UpdatedAt.Time,
	}, "Care plan objective created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// CarePlanActions represents the actions in a care plan objective
type CarePlanActions struct {
	ActionID          int64  `json:"action_id"`
	SortOrder         int32  `json:"sort_order"`
	ActionDescription string `json:"action_description"`
	IsCompleted       bool   `json:"is_completed"`
	Notes             string `json:"notes"`
}

// CarePlanObjectives represents the objectives in a care plan
type CarePlanObjectives struct {
	ObjectiveID int64             `json:"objective_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	TimeFrame   string            `json:"timeframe"`
	Status      string            `json:"status"`
	Actions     []CarePlanActions `json:"actions"`
}

// GetCarePlanObjectivesResponse represents the response for the GetCarePlanObjectives API
type GetCarePlanObjectivesResponse struct {
	ShortTermGoals  []CarePlanObjectives `json:"short_term_goals"`
	MediumTermGoals []CarePlanObjectives `json:"medium_term_goals"`
	LongTermGoals   []CarePlanObjectives `json:"long_term_goals"`
}

// GetCarePlanObjectivesApi retrieves the care plan objectives for a given care plan ID
// @Summary Get care plan objectives
// @Description Get the care plan objectives for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[GetCarePlanObjectivesResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/objectives [get]
func (server *Server) GetCarePlanObjectivesApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanObjectivesApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	rows, err := server.store.GetCarePlanObjectivesWithActions(ctx, carePlanID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanObjectivesApi", "Failed to get care plan objectives with actions", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan objectives with actions")))
		return
	}

	objectiveMap := make(map[int64]*CarePlanObjectives)
	for _, row := range rows {
		objective, exists := objectiveMap[row.ObjectiveID]
		if !exists {
			objective = &CarePlanObjectives{
				ObjectiveID: row.ObjectiveID,
				Title:       row.ObjectiveTitle,
				Description: row.ObjectiveDescription,
				TimeFrame:   row.ObjectiveTimeframe,
				Status:      row.ObjectiveStatus,
				Actions:     []CarePlanActions{},
			}
			objectiveMap[row.ObjectiveID] = objective
		}
		if row.ActionID != nil {
			action := CarePlanActions{
				ActionID:          *row.ActionID,
				SortOrder:         util.DerefInt32(row.SortOrder),
				ActionDescription: util.DerefString(row.ActionDescription),
				IsCompleted:       util.DerefBool(row.IsCompleted),
				Notes:             util.DerefString(row.ActionNotes),
			}
			objective.Actions = append(objective.Actions, action)
		}
	}

	response := &GetCarePlanObjectivesResponse{
		ShortTermGoals:  make([]CarePlanObjectives, 0),
		MediumTermGoals: make([]CarePlanObjectives, 0),
		LongTermGoals:   make([]CarePlanObjectives, 0),
	}

	for _, objective := range objectiveMap {
		switch objective.TimeFrame {
		case "short_term":
			response.ShortTermGoals = append(response.ShortTermGoals, *objective)
		case "medium_term":
			response.MediumTermGoals = append(response.MediumTermGoals, *objective)
		case "long_term":
			response.LongTermGoals = append(response.LongTermGoals, *objective)
		}
	}

	res := SuccessResponse(response, "Care plan objectives retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanObjectiveRequest represents the request body for updating a care plan objective
type UpdateCarePlanObjectiveRequest struct {
	TimeFrame   *string `json:"timeframe" binding:"oneof=short_term medium_term long_term"`
	GoalTitle   *string `json:"goal_title"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"oneof=not_started in_progress completed discontinued"`
}

// UpdateCarePlanObjectiveResponse represents the response for the UpdateCarePlanObjective API
type UpdateCarePlanObjectiveResponse struct {
	ObjectiveID int64 `json:"goal_id"`
	CarePlanId  int64 `json:"care_plan_id"`
}

// @Summary Update care plan objective
// @Description Update a care plan objective by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param objective_id path int true "Objective ID"
// @Param request body UpdateCarePlanObjectiveRequest true "Request body"
// @Success 200 {object} Response[UpdateCarePlanObjectiveResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /objectives/{objective_id} [put]
func (server *Server) UpdateCarePlanObjectiveApi(ctx *gin.Context) {
	objectiveId, err := strconv.ParseInt(ctx.Param("objective_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanObjectiveApi", "Invalid objective ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid objective ID")))
		return
	}

	var req UpdateCarePlanObjectiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanObjectiveApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	objective, err := server.store.UpdateCarePlanObjective(ctx, db.UpdateCarePlanObjectiveParams{
		ID:          objectiveId,
		Timeframe:   req.TimeFrame,
		GoalTitle:   req.GoalTitle,
		Description: req.Description,
		Status:      req.Status,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanObjectiveApi", "Failed to update care plan objective", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan objective")))
		return
	}

	res := SuccessResponse(UpdateCarePlanObjectiveResponse{
		ObjectiveID: objective.ID,
		CarePlanId:  objective.CarePlanID,
	}, "Care plan objective updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteCarePlanObjectiveApi deletes a care plan objective by its ID
// @Summary Delete care plan objective
// @Description Delete a care plan objective by its ID
// @Tags care_plan
// @Produce json
// @Param objective_id path int true "Objective ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /objectives/{objective_id} [delete]
func (server *Server) DeleteCarePlanObjectiveApi(ctx *gin.Context) {
	objectiveId, err := strconv.ParseInt(ctx.Param("objective_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanObjectiveApi", "Invalid objective ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid objective ID")))
		return
	}

	err = server.store.DeleteCarePlanObjective(ctx, objectiveId)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanObjectiveApi", "Failed to delete care plan objective", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan objective")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan objective deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateCarePlanActionsRequest represents the request body for creating a care plan action
type CreateCarePlanActionsRequest struct {
	ActionDescription string `json:"action_description" binding:"required"`
}

// CreateCarePlanActionsResponse represents the response for the CreateCarePlanActions API
type CreateCarePlanActionsResponse struct {
	ActionID          int64  `json:"action_id"`
	ObjectiveID       int64  `json:"objective_id"`
	ActionDescription string `json:"action_description"`
}

// CreateCarePlanActionsApi creates a new care plan action
// @Summary Create care plan action
// @Description Create a new care plan action
// @Tags care_plan
// @Accept json
// @Produce json
// @Param objective_id path int true "Objective ID"
// @Param request body CreateCarePlanActionsRequest true "Request body"
// @Success 201 {object} Response[CreateCarePlanActionsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /objectives/{objective_id}/actions [post]
func (server *Server) CreateCarePlanActionsApi(ctx *gin.Context) {
	objectiveID, err := strconv.ParseInt(ctx.Param("objective_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanActionsApi", "Invalid objective ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid objective ID")))
		return
	}

	var req CreateCarePlanActionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanActionsApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	maxSortOrder, err := server.store.GetCarePlanActionsMaxSortOrder(ctx, objectiveID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanActionsApi", "Failed to get max sort order for objective", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get max sort order for objective")))
		return
	}

	action, err := server.store.CreateCarePlanAction(ctx, db.CreateCarePlanActionParams{
		ObjectiveID:       objectiveID,
		ActionDescription: req.ActionDescription,
		SortOrder:         maxSortOrder + 1, // Increment the max sort order by 1
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanActionsApi", "Failed to create care plan action", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan action")))
		return
	}

	res := SuccessResponse(CreateCarePlanActionsResponse{
		ActionID:          action.ID,
		ObjectiveID:       action.ObjectiveID,
		ActionDescription: action.ActionDescription,
	}, "Care plan action created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// UpdateCarePlanActionsRequest represents the request body for updating a care plan action
type UpdateCarePlanActionsRequest struct {
	ActionDescription *string `json:"action_description"`
}

// UpdateCarePlanActionsResponse represents the response for the UpdateCarePlanActions API
type UpdateCarePlanActionsResponse struct {
	ActionID          int64  `json:"action_id"`
	ObjectiveID       int64  `json:"objective_id"`
	ActionDescription string `json:"action_description"`
}

// UpdateCarePlanActionsApi updates a care plan action by its ID
// @Summary Update care plan action
// @Description Update a care plan action by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param action_id path int true "Action ID"
// @Param request body UpdateCarePlanActionsRequest true "Request body"
// @Success 200 {object} Response[UpdateCarePlanActionsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /actions/{action_id} [put]
func (server *Server) UpdateCarePlanActionsApi(ctx *gin.Context) {
	actionID, err := strconv.ParseInt(ctx.Param("action_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanActionsApi", "Invalid action ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid action ID")))
		return
	}

	var req UpdateCarePlanActionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanActionsApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	action, err := server.store.UpdateCarePlanAction(ctx, db.UpdateCarePlanActionParams{
		ID:                actionID,
		ActionDescription: req.ActionDescription,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanActionsApi", "Failed to update care plan action", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan action")))
		return
	}

	res := SuccessResponse(UpdateCarePlanActionsResponse{
		ActionID:          action.ID,
		ObjectiveID:       action.ObjectiveID,
		ActionDescription: action.ActionDescription,
	}, "Care plan action updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteCarePlanActionApi deletes a care plan action by its ID
// @Summary Delete care plan action
// @Description Delete a care plan action by its ID
// @Tags care_plan
// @Produce json
// @Param action_id path int true "Action ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /actions/{action_id} [delete]
func (server *Server) DeleteCarePlanActionApi(ctx *gin.Context) {
	actionID, err := strconv.ParseInt(ctx.Param("action_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanActionApi", "Invalid action ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid action ID")))
		return
	}

	err = server.store.DeleteCarePlanAction(ctx, actionID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanActionApi", "Failed to delete care plan action", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan action")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan action deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// =========================== Care plan interventions ===========================

// CreateCarePlanInterventionRequest represents the request body for creating a care plan intervention
type CreateCarePlanInterventionRequest struct {
	Frequency               string `json:"frequency" binding:"required,oneof=daily weekly monthly"`
	InterventionDescription string `json:"intervention_description" binding:"required"`
}

// CreateCarePlanInterventionResponse represents the response for the CreateCarePlanIntervention API
type CreateCarePlanInterventionResponse struct {
	InterventionID          int64  `json:"intervention_id"`
	CarePlanID              int64  `json:"care_plan_id"`
	Frequency               string `json:"frequency"`
	InterventionDescription string `json:"intervention_description"`
}

// CreateCarePlanInterventionApi creates a new care plan intervention
// @Summary Create care plan intervention
// @Description Create a new care plan intervention
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body CreateCarePlanInterventionRequest true "Request body"
// @Success 201 {object} Response[CreateCarePlanInterventionResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/interventions [post]
func (server *Server) CreateCarePlanInterventionApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanInterventionApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req CreateCarePlanInterventionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanInterventionApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	intervention, err := server.store.CreateCarePlanIntervention(ctx, db.CreateCarePlanInterventionParams{
		CarePlanID:              carePlanID,
		Frequency:               req.Frequency,
		InterventionDescription: req.InterventionDescription,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanInterventionApi", "Failed to create care plan intervention", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan intervention")))
		return
	}

	res := SuccessResponse(CreateCarePlanInterventionResponse{
		InterventionID:          intervention.ID,
		CarePlanID:              intervention.CarePlanID,
		Frequency:               intervention.Frequency,
		InterventionDescription: intervention.InterventionDescription,
	}, "Care plan intervention created successfully")

	ctx.JSON(http.StatusCreated, res)
}

type Intervention struct {
	InterventionID          int64  `json:"intervention_id"`
	InterventionDescription string `json:"intervention_description"`
}

// GetCarePlanInterventionsResponse represents the response for the GetCarePlanInterventions API
type GetCarePlanInterventionsResponse struct {
	DailyActivities   []Intervention `json:"daily_activities"`
	WeeklyActivities  []Intervention `json:"weekly_activities"`
	MonthlyActivities []Intervention `json:"monthly_activities"`
}

// GetCarePlanInterventionsApi retrieves the care plan interventions for a given care plan ID
// @Summary Get care plan interventions
// @Description Get the care plan interventions for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[GetCarePlanInterventionsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/interventions [get]
func (server *Server) GetCarePlanInterventionsApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanInterventionsApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	interventions, err := server.store.GetCarePlanInterventions(ctx, carePlanID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanInterventionsApi", "Failed to get care plan interventions", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan interventions")))
		return
	}

	response := &GetCarePlanInterventionsResponse{
		DailyActivities:   []Intervention{},
		WeeklyActivities:  []Intervention{},
		MonthlyActivities: []Intervention{},
	}
	for _, intervention := range interventions {
		switch intervention.Frequency {
		case "daily":
			response.DailyActivities = append(response.DailyActivities, Intervention{
				InterventionID:          intervention.ID,
				InterventionDescription: intervention.InterventionDescription,
			})
		case "weekly":
			response.WeeklyActivities = append(response.WeeklyActivities, Intervention{
				InterventionID:          intervention.ID,
				InterventionDescription: intervention.InterventionDescription,
			})
		case "monthly":
			response.MonthlyActivities = append(response.MonthlyActivities, Intervention{
				InterventionID:          intervention.ID,
				InterventionDescription: intervention.InterventionDescription,
			})
		default:
			server.logBusinessEvent(LogLevelError, "GetCarePlanInterventionsApi", "Unknown intervention frequency", zap.String("frequency", intervention.Frequency))
		}

	}
	res := SuccessResponse(response, "Care plan interventions retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanInterventionRequest represents the request body for updating a care plan intervention
type UpdateCarePlanInterventionRequest struct {
	Frequency               *string `json:"frequency" binding:"oneof=daily weekly monthly"`
	InterventionDescription *string `json:"intervention_description"`
}

// UpdateCarePlanInterventionApi updates a care plan intervention by its ID
type UpdateCarePlanInterventionResponse struct {
	InterventionID          int64  `json:"intervention_id"`
	CarePlanID              int64  `json:"care_plan_id"`
	Frequency               string `json:"frequency"`
	InterventionDescription string `json:"intervention_description"`
}

// @Summary Update care plan intervention
// @Description Update a care plan intervention by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param intervention_id path int true "Intervention ID"
// @Param request body UpdateCarePlanInterventionRequest true "Request body"
// @Success 200 {object} Response[UpdateCarePlanInterventionResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /interventions/{intervention_id} [put]
func (server *Server) UpdateCarePlanInterventionApi(ctx *gin.Context) {
	interventionID, err := strconv.ParseInt(ctx.Param("intervention_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanInterventionApi", "Invalid intervention ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid intervention ID")))
		return
	}

	var req UpdateCarePlanInterventionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanInterventionApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	intervention, err := server.store.UpdateCarePlanIntervention(ctx, db.UpdateCarePlanInterventionParams{
		ID:                      interventionID,
		Frequency:               req.Frequency,
		InterventionDescription: req.InterventionDescription,
	})

	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanInterventionApi", "Failed to update care plan intervention", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan intervention")))
		return
	}
	res := SuccessResponse(UpdateCarePlanInterventionResponse{
		InterventionID:          intervention.ID,
		CarePlanID:              intervention.CarePlanID,
		Frequency:               intervention.Frequency,
		InterventionDescription: intervention.InterventionDescription,
	}, "Care plan intervention updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteCarePlanInterventionApi deletes a care plan intervention by its ID
// @Summary Delete care plan intervention
// @Description Delete a care plan intervention by its ID
// @Tags care_plan
// @Produce json
// @Param intervention_id path int true "Intervention ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /interventions/{intervention_id} [delete]
func (server *Server) DeleteCarePlanInterventionApi(ctx *gin.Context) {
	interventionID, err := strconv.ParseInt(ctx.Param("intervention_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanInterventionApi", "Invalid intervention ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid intervention ID")))
		return
	}

	err = server.store.DeleteCarePlanIntervention(ctx, interventionID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanInterventionApi", "Failed to delete care plan intervention", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan intervention")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan intervention deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
// =========================== Care plan success metrics ===========================
//////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////

// CreateCarePlanSuccessMetricsRequest represents the request body for creating a care plan success metric
type CreateCarePlanSuccessMetricsRequest struct {
	MetricName        string  `json:"metric_name" binding:"required"`
	TargetValue       string  `json:"target_value" binding:"required"`
	MeasurementMethod string  `json:"measurement_method" binding:"required"`
	CurrentValue      *string `json:"current_value"` // Optional, can be nil if not set
}

// CreateCarePlanSuccessMetricsResponse represents the response for the CreateCarePlanSuccessMetrics API
type CreateCarePlanSuccessMetricsResponse struct {
	MetricID          int64   `json:"metric_id"`
	MetricName        string  `json:"metric_name"`
	CurrentValue      *string `json:"current_value"`
	TargetValue       string  `json:"target_value"`
	MeasurementMethod string  `json:"measurement_method"`
}

// CreateCarePlanSuccessMetricsApi creates a new care plan success metric
// @Summary Create care plan success metric
// @Description Create a new care plan success metric
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body CreateCarePlanSuccessMetricsRequest true "Request body"
// @Success 201 {object} Response[CreateCarePlanSuccessMetricsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/success_metrics [post]
func (server *Server) CreateCarePlanSuccessMetricsApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanSuccessMetricsApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req CreateCarePlanSuccessMetricsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanSuccessMetricsApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	successMetric, err := server.store.CreateCarePlanSuccessMetric(ctx, db.CreateCarePlanSuccessMetricParams{
		CarePlanID:        carePlanID,
		MetricName:        req.MetricName,
		TargetValue:       req.TargetValue,
		CurrentValue:      req.CurrentValue,
		MeasurementMethod: req.MeasurementMethod,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanSuccessMetricsApi", "Failed to create care plan success metric", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan success metric")))
		return
	}

	res := SuccessResponse(CreateCarePlanSuccessMetricsResponse{
		MetricID:          successMetric.ID,
		MetricName:        successMetric.MetricName,
		CurrentValue:      successMetric.CurrentValue,
		TargetValue:       successMetric.TargetValue,
		MeasurementMethod: successMetric.MeasurementMethod,
	}, "Care plan success metric created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GetCarePlanSuccessMetricsResponse represents the response for the GetCarePlanSuccessMetrics AP
type GetCarePlanSuccessMetricsResponse struct {
	MetricID          int64   `json:"metric_id"`
	MetricName        string  `json:"metric_name"`
	CurrentValue      *string `json:"current_value"`
	TargetValue       string  `json:"target_value"`
	MeasurementMethod string  `json:"measurement_method"`
}

// GetCarePlanSuccessMetricsApi retrieves the success metrics for a given care plan ID
// @Summary Get care plan success metrics
// @Description Get the success metrics for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[[]GetCarePlanSuccessMetricsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/success_metrics [get]
func (server *Server) GetCarePlanSuccessMetricsApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanSuccessMetricsApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	successMetrics, err := server.store.GetCarePlanSuccessMetrics(ctx, carePlanID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanSuccessMetricsApi", "Failed to get care plan success metrics", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan success metrics")))
		return
	}
	if len(successMetrics) == 0 {
		res := SuccessResponse([]GetCarePlanSuccessMetricsResponse{}, "No care plan success metrics found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	response := make([]GetCarePlanSuccessMetricsResponse, len(successMetrics))
	for i, metric := range successMetrics {
		response[i] = GetCarePlanSuccessMetricsResponse{
			MetricID:          metric.ID,
			MetricName:        metric.MetricName,
			CurrentValue:      metric.CurrentValue,
			TargetValue:       metric.TargetValue,
			MeasurementMethod: metric.MeasurementMethod,
		}
	}

	res := SuccessResponse(response, "Care plan success metrics retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanSuccessMetricsRequest represents the request body for updating a care plan success metric
type UpdateCarePlanSuccessMetricsRequest struct {
	MetricName        *string `json:"metric_name"`
	TargetValue       *string `json:"target_value"`
	MeasurementMethod *string `json:"measurement_method"`
	CurrentValue      *string `json:"current_value"` // Optional, can be nil if not set
}

// UpdateCarePlanSuccessMetricsResponse represents the response for the UpdateCarePlanSuccessMetrics API
type UpdateCarePlanSuccessMetricsResponse struct {
	MetricID          int64   `json:"metric_id"`
	MetricName        string  `json:"metric_name"`
	CurrentValue      *string `json:"current_value"`
	TargetValue       string  `json:"target_value"`
	MeasurementMethod string  `json:"measurement_method"`
}

// UpdateCarePlanSuccessMetricsApi updates a care plan success metric by its ID
// @Summary Update care plan success metric
// @Description Update a care plan success metric by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param metric_id path int true "Metric ID"
// @Param request body UpdateCarePlanSuccessMetricsRequest true "Request body"
// @Success 200 {object} Response[UpdateCarePlanSuccessMetricsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /success_metrics/{metric_id} [put]
func (server *Server) UpdateCarePlanSuccessMetricsApi(ctx *gin.Context) {
	metricID, err := strconv.ParseInt(ctx.Param("metric_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanSuccessMetricsApi", "Invalid metric ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid metric ID")))
		return
	}

	var req UpdateCarePlanSuccessMetricsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanSuccessMetricsApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	successMetric, err := server.store.UpdateCarePlanSuccessMetric(ctx, db.UpdateCarePlanSuccessMetricParams{
		ID:                metricID,
		MetricName:        req.MetricName,
		TargetValue:       req.TargetValue,
		MeasurementMethod: req.MeasurementMethod,
		CurrentValue:      req.CurrentValue,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanSuccessMetricsApi", "Failed to update care plan success metric", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan success metric")))
		return
	}

	res := SuccessResponse(UpdateCarePlanSuccessMetricsResponse{
		MetricID:          successMetric.ID,
		MetricName:        successMetric.MetricName,
		CurrentValue:      successMetric.CurrentValue,
		TargetValue:       successMetric.TargetValue,
		MeasurementMethod: successMetric.MeasurementMethod,
	}, "Care plan success metric updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteCarePlanSuccessMetricApi deletes a care plan success metric by its ID
// @Summary Delete care plan success metric
// @Description Delete a care plan success metric by its ID
// @Tags care_plan
// @Produce json
// @Param metric_id path int true "Metric ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /success_metrics/{metric_id} [delete]
func (server *Server) DeleteCarePlanSuccessMetricApi(ctx *gin.Context) {
	metricID, err := strconv.ParseInt(ctx.Param("metric_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanSuccessMetricApi", "Invalid metric ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid metric ID")))
		return
	}

	err = server.store.DeleteCarePlanSuccessMetric(ctx, metricID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanSuccessMetricApi", "Failed to delete care plan success metric", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan success metric")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan success metric deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// =========================== Care plan risks ===========================

// CreateCarePlanRisksRequest represents the request body for creating a care plan risk
type CreateCarePlanRisksRequest struct {
	RiskDescription    string  `json:"risk_description" binding:"required"`
	MitigationStrategy string  `json:"mitigation_strategy" binding:"required"`
	RiskLevel          *string `json:"risk_level"`
}

// CreateCarePlanRisksResponse represents the response for the CreateCarePlanRisks API
type CreateCarePlanRisksResponse struct {
	RiskID             int64   `json:"risk_id"`
	RiskDescription    string  `json:"risk_description"`
	MitigationStrategy string  `json:"mitigation_strategy"`
	RiskLevel          *string `json:"risk_level"`
}

// CreateCarePlanRisksApi creates a new care plan risk
// @Summary Create care plan risk
// @Description Create a new care plan risk
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body CreateCarePlanRisksRequest true "Request body"
// @Success 201 {object} Response[CreateCarePlanRisksResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/risks [post]
func (server *Server) CreateCarePlanRisksApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanRisksApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req CreateCarePlanRisksRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanRisksApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	risk, err := server.store.CreateCarePlanRisk(ctx, db.CreateCarePlanRiskParams{
		CarePlanID:         carePlanID,
		RiskDescription:    req.RiskDescription,
		MitigationStrategy: req.MitigationStrategy,
		RiskLevel:          req.RiskLevel,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanRisksApi", "Failed to create care plan risk", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan risk")))
		return
	}

	res := SuccessResponse(CreateCarePlanRisksResponse{
		RiskID:             risk.ID,
		RiskDescription:    risk.RiskDescription,
		MitigationStrategy: risk.MitigationStrategy,
		RiskLevel:          risk.RiskLevel,
	}, "Care plan risk created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GetCarePlanRisksResponse represents the response for the GetCarePlanRisks API
type GetCarePlanRisksResponse struct {
	RiskID             int64   `json:"risk_id"`
	RiskDescription    string  `json:"risk_description"`
	MitigationStrategy string  `json:"mitigation_strategy"`
	RiskLevel          *string `json:"risk_level"`
}

// GetCarePlanRisksApi retrieves the risks associated with a given care plan ID
// @Summary Get care plan risks
// @Description Get the risks associated with a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[[]GetCarePlanRisksResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plan/{care_plan_id}/risks [get]
func (server *Server) GetCarePlanRisksApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanRisksApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	risks, err := server.store.GetCarePlanRisks(ctx, carePlanID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanRisksApi", "Failed to get care plan risks", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan risks")))
		return
	}

	if len(risks) == 0 {
		res := SuccessResponse([]GetCarePlanRisksResponse{}, "No care plan risks found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	response := make([]GetCarePlanRisksResponse, len(risks))
	for i, risk := range risks {
		response[i] = GetCarePlanRisksResponse{
			RiskID:             risk.ID,
			RiskDescription:    risk.RiskDescription,
			MitigationStrategy: risk.MitigationStrategy,
			RiskLevel:          risk.RiskLevel,
		}
	}

	res := SuccessResponse(response, "Care plan risks retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanRisksRequest represents the request body for updating a care plan risk
type UpdateCarePlanRisksRequest struct {
	RiskDescription    *string `json:"risk_description"`
	MitigationStrategy *string `json:"mitigation_strategy"`
	RiskLevel          *string `json:"risk_level"`
}

// UpdateCarePlanRisksResponse represents the response for the UpdateCarePlanRisks API
type UpdateCarePlanRisksResponse struct {
	RiskID             int64   `json:"risk_id"`
	RiskDescription    string  `json:"risk_description"`
	MitigationStrategy string  `json:"mitigation_strategy"`
	RiskLevel          *string `json:"risk_level"`
}

// UpdateCarePlanRisksApi updates a care plan risk by its ID
// @Summary Update care plan risk
// @Description Update a care plan risk by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param risk_id path int true "Risk ID"
// @Param request body UpdateCarePlanRisksRequest true "Request body"
// @Success 200 {object} Response[UpdateCarePlanRisksResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /risks/{risk_id} [put]
func (server *Server) UpdateCarePlanRisksApi(ctx *gin.Context) {
	riskID, err := strconv.ParseInt(ctx.Param("risk_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanRisksApi", "Invalid risk ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid risk ID")))
		return
	}

	var req UpdateCarePlanRisksRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanRisksApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	risk, err := server.store.UpdateCarePlanRisk(ctx, db.UpdateCarePlanRiskParams{
		ID:                 riskID,
		RiskDescription:    req.RiskDescription,
		MitigationStrategy: req.MitigationStrategy,
		RiskLevel:          req.RiskLevel,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanRisksApi", "Failed to update care plan risk", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan risk")))
		return
	}

	res := SuccessResponse(UpdateCarePlanRisksResponse{
		RiskID:             risk.ID,
		RiskDescription:    risk.RiskDescription,
		MitigationStrategy: risk.MitigationStrategy,
		RiskLevel:          risk.RiskLevel,
	}, "Care plan risk updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteCarePlanRiskApi deletes a care plan risk by its ID
// @Summary Delete care plan risk
// @Description Delete a care plan risk by its ID
// @Tags care_plan
// @Produce json
// @Param risk_id path int true "Risk ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /risks/{risk_id} [delete]
func (server *Server) DeleteCarePlanRiskApi(ctx *gin.Context) {
	riskID, err := strconv.ParseInt(ctx.Param("risk_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanRiskApi", "Invalid risk ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid risk ID")))
		return
	}

	err = server.store.DeleteCarePlanRisk(ctx, riskID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanRiskApi", "Failed to delete care plan risk", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan risk")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan risk deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

////////////////////////////////////////////////////////
////////////////////////////////////////////////////////
///////////////////////////////////////////////////////
// =========================== Care plan support network ===========================
//////////////////////////////////////////////////////
//////////////////////////////////////////////////////
///////////////////////////////////////////////////////
//////////////////////////////////////////////////////

// CreateCarePlanSupportNetworkRequest represents the request body for creating a care plan support network
type CreateCarePlanSupportNetworkRequest struct {
	RoleTitle                 string `json:"role_title" binding:"required"`
	ResponsibilityDescription string `json:"responsibility_description"`
}

// CreateCarePlanSupportNetworkResponse represents the response for the CreateCarePlanSupportNetwork API
type CreateCarePlanSupportNetworkResponse struct {
	SupportNetworkID          int64  `json:"support_network_id"`
	RoleTitle                 string `json:"role_title"`
	ResponsibilityDescription string `json:"responsibility_description"`
}

// CreateCareplanSupportNetworkApi creates a new care plan support network
// @Summary Create care plan support network
// @Description Create a new care plan support network
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body CreateCarePlanSupportNetworkRequest true "Request body"
// @Success 201 {object} Response[CreateCarePlanSupportNetworkResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/support_network [post]
func (server *Server) CreateCareplanSupportNetworkApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCareplanSupportNetworkApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req CreateCarePlanSupportNetworkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCareplanSupportNetworkApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	supportNetwork, err := server.store.CreateCarePlanSupportNetwork(ctx, db.CreateCarePlanSupportNetworkParams{
		CarePlanID:                carePlanID,
		RoleTitle:                 req.RoleTitle,
		ResponsibilityDescription: req.ResponsibilityDescription,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCareplanSupportNetworkApi", "Failed to create care plan support network", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan support network")))
		return
	}

	res := SuccessResponse(CreateCarePlanSupportNetworkResponse{
		SupportNetworkID:          supportNetwork.ID,
		RoleTitle:                 supportNetwork.RoleTitle,
		ResponsibilityDescription: supportNetwork.ResponsibilityDescription,
	}, "Care plan support network created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GetCarePlanSupportNetworkResponse represents the response for the GetCarePlanSupportNetwork API
type GetCarePlanSupportNetworkResponse struct {
	SupportNetworkID          int64   `json:"support_network_id"`
	RoleTitle                 string  `json:"role_title"`
	ResponsibilityDescription *string `json:"responsibility_description"`
}

// @Summary Get care plan support network
// @Description Get the support network for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[[]GetCarePlanSupportNetworkResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/support_network [get]
func (server *Server) GetCarePlanSupportNetworkApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanSupportNetworkApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	supportNetwork, err := server.store.GetCarePlanSupportNetwork(ctx, carePlanID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCareplanSupportNetworkApi", "Failed to get care plan support network", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan support network")))
		return
	}

	if len(supportNetwork) == 0 {
		res := SuccessResponse([]GetCarePlanSupportNetworkResponse{}, "No care plan support network found")
		ctx.JSON(http.StatusOK, res)
		return
	}
	response := make([]GetCarePlanSupportNetworkResponse, len(supportNetwork))
	for i, support := range supportNetwork {
		response[i] = GetCarePlanSupportNetworkResponse{
			SupportNetworkID:          support.ID,
			RoleTitle:                 support.RoleTitle,
			ResponsibilityDescription: &support.ResponsibilityDescription,
		}
	}

	res := SuccessResponse(response, "Care plan support network retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanSupportNetworkRequest represents the request body for updating a care plan support network
type UpdateCarePlanSupportNetworkRequest struct {
	RoleTitle                 *string `json:"role_title"`
	ResponsibilityDescription *string `json:"responsibility_description"`
}

// UpdateCarePlanSupportNetworkResponse represents the response for the UpdateCarePlanSupportNetwork API
type UpdateCarePlanSupportNetworkResponse struct {
	SupportNetworkID          int64  `json:"support_network_id"`
	RoleTitle                 string `json:"role_title"`
	ResponsibilityDescription string `json:"responsibility_description"`
}

// UpdateCarePlanSupportNetworkApi updates a care plan support network by its ID
// @Summary Update care plan support network
// @Description Update a care plan support network by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param support_network_id path int true "Support Network ID"
// @Param request body UpdateCarePlanSupportNetworkRequest true "Request body"
// @Success 200 {object} Response[UpdateCarePlanSupportNetworkResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /support_network/{support_network_id} [put]
func (server *Server) UpdateCarePlanSupportNetworkApi(ctx *gin.Context) {
	supportNetworkID, err := strconv.ParseInt(ctx.Param("support_network_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanSupportNetworkApi", "Invalid support network ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid support network ID")))
		return
	}

	var req UpdateCarePlanSupportNetworkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanSupportNetworkApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	supportNetwork, err := server.store.UpdateCarePlanSupportNetwork(ctx, db.UpdateCarePlanSupportNetworkParams{
		ID:                        supportNetworkID,
		RoleTitle:                 req.RoleTitle,
		ResponsibilityDescription: req.ResponsibilityDescription,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanSupportNetworkApi", "Failed to update care plan support network", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan support network")))
		return
	}

	res := SuccessResponse(UpdateCarePlanSupportNetworkResponse{
		SupportNetworkID:          supportNetwork.ID,
		RoleTitle:                 supportNetwork.RoleTitle,
		ResponsibilityDescription: supportNetwork.ResponsibilityDescription,
	}, "Care plan support network updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteCarePlanSupportNetworkApi deletes a care plan support network by its ID
// @Summary Delete care plan support network
// @Description Delete a care plan support network by its ID
// @Tags care_plan
// @Produce json
// @Param support_network_id path int true "Support Network ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /support_network/{support_network_id} [delete]
func (server *Server) DeleteCarePlanSupportNetworkApi(ctx *gin.Context) {
	supportNetworkID, err := strconv.ParseInt(ctx.Param("support_network_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanSupportNetworkApi", "Invalid support network ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid support network ID")))
		return
	}

	err = server.store.DeleteCarePlanSupportNetwork(ctx, supportNetworkID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanSupportNetworkApi", "Failed to delete care plan support network", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan support network")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan support network deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

/////////////////////////////////////////

//////////////////////////////////////////

// =========================== Care plan resources ===========================

/////////////////////////////////////////////

// /////////////////////////////////////////
// CreateCarePlanResourcesRequest represents the request body for creating a care plan resource
type CreateCarePlanResourcesRequest struct {
	ResourceDescription string     `json:"resource_description" binding:"required"`
	IsObtained          *bool      `json:"is_obtained"`
	ObtainedDate        *time.Time `json:"obtained_date"`
}

// CreateCarePlanResourcesResponse represents the response for the CreateCarePlanResources API
type CreateCarePlanResourcesResponse struct {
	ID                  int64      `json:"id"`
	ResourceDescription string     `json:"resource_description"`
	IsObtained          bool       `json:"is_obtained"`
	ObtainedDate        *time.Time `json:"obtained_date"`
}

// CreateCarePlanResourcesApi creates a new care plan resource
// @Summary Create care plan resource
// @Description Create a new care plan resource
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body CreateCarePlanResourcesRequest true "Request body"
// @Success 201 {object} Response[CreateCarePlanResourcesResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plan/{care_plan_id}/resources [post]
func (server *Server) CreateCarePlanResourcesApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanResourcesApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req CreateCarePlanResourcesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanResourcesApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	resource, err := server.store.CreateCarePlanResources(ctx, db.CreateCarePlanResourcesParams{
		CarePlanID:          carePlanID,
		ResourceDescription: req.ResourceDescription,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanResourcesApi", "Failed to create care plan resource", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan resource")))
		return
	}

	res := SuccessResponse(CreateCarePlanResourcesResponse{
		ID:                  resource.ID,
		ResourceDescription: resource.ResourceDescription,
		IsObtained:          resource.IsObtained,
		ObtainedDate:        &resource.ObtainedDate.Time,
	}, "Care plan resource created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GetCarePlanResourcesResponse represents the response for the GetCarePlanResources API
type GetCarePlanResourcesResponse struct {
	ID                  int64      `json:"id"`
	ResourceDescription string     `json:"resource_description"`
	IsObtained          bool       `json:"is_obtained"`
	ObtainedDate        *time.Time `json:"obtained_date"`
}

// GetCarePlanResourcesApi retrieves the resources for a given care plan ID
// @Summary Get care plan resources
// @Description Get the resources for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[[]GetCarePlanResourcesResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plan/{care_plan_id}/resources [get]
func (server *Server) GetCarePlanResourcesApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanResourcesApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	resources, err := server.store.GetCarePlanResources(ctx, carePlanID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "GetCarePlanResourcesApi", "Failed to get care plan resources", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan resources")))
		return
	}

	if len(resources) == 0 {
		res := SuccessResponse([]GetCarePlanResourcesResponse{}, "No care plan resources found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	response := make([]GetCarePlanResourcesResponse, len(resources))
	for i, resource := range resources {
		response[i] = GetCarePlanResourcesResponse{
			ID:                  resource.ID,
			ResourceDescription: resource.ResourceDescription,
			IsObtained:          resource.IsObtained,
			ObtainedDate:        &resource.ObtainedDate.Time,
		}
	}

	res := SuccessResponse(response, "Care plan resources retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanResourcesRequest represents the request body for updating a care plan resource
type UpdateCarePlanResourcesRequest struct {
	ResourceDescription *string   `json:"resource_description"`
	IsObtained          *bool     `json:"is_obtained"`
	ObtainedDate        time.Time `json:"obtained_date"`
}

// UpdateCarePlanResourcesResponse represents the response for the UpdateCarePlanResources API
type UpdateCarePlanResourcesResponse struct {
	ID                  int64     `json:"id"`
	ResourceDescription string    `json:"resource_description"`
	IsObtained          bool      `json:"is_obtained"`
	ObtainedDate        time.Time `json:"obtained_date"`
}

// UpdateCarePlanResourcesApi updates a care plan resource by its ID
// @Summary Update care plan resource
// @Description Update a care plan resource by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param resource_id path int true "Resource ID"
// @Param request body UpdateCarePlanResourcesRequest true "Request body"
// @Success 200 {object} Response[UpdateCarePlanResourcesResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /resources/{resource_id} [put]
func (server *Server) UpdateCarePlanResourcesApi(ctx *gin.Context) {
	resourceID, err := strconv.ParseInt(ctx.Param("resource_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanResourcesApi", "Invalid resource ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid resource ID")))
		return
	}

	var req UpdateCarePlanResourcesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanResourcesApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	resource, err := server.store.UpdateCarePlanResource(ctx, db.UpdateCarePlanResourceParams{
		ID:                  resourceID,
		ResourceDescription: req.ResourceDescription,
		IsObtained:          req.IsObtained,
		ObtainedDate:        pgtype.Date{Time: req.ObtainedDate, Valid: true},
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanResourcesApi", "Failed to update care plan resource", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan resource")))
		return
	}

	res := SuccessResponse(UpdateCarePlanResourcesResponse{
		ID:                  resource.ID,
		ResourceDescription: resource.ResourceDescription,
		IsObtained:          resource.IsObtained,
		ObtainedDate:        resource.ObtainedDate.Time,
	}, "Care plan resource updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteCarePlanResourcesApi deletes all resources associated with a care plan
// @Summary Delete care plan resources
// @Description Delete all resources associated with a care plan
// @Tags care_plan
// @Produce json
// @Param resource_id path int true "Resource ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /resources/{resource_id} [delete]
func (server *Server) DeleteCarePlanResourcesApi(ctx *gin.Context) {
	resourceID, err := strconv.ParseInt(ctx.Param("resource_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanResourcesApi", "Invalid resource ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid resource ID")))
		return
	}

	err = server.store.DeleteCarePlanResource(ctx, resourceID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanResourcesApi", "Failed to delete care plan resources", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan resources")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan resources deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

//////////////////////

////////////////////// Care Plan Reports /////////

////////////////////////////

// CreateCarePlanReportRequest represents the request body for creating a care plan report
type CreateCarePlanReportRequest struct {
	ReportType    string `json:"report_type" binding:"required" oneof:"progress concern achievement modification"`
	ReportContent string `json:"report_content" binding:"required"`
	IsCritical    bool   `json:"is_critical"`
}

// CreateCarePlanReportResponse represents the response for the CreateCarePlanReport API
type CreateCarePlanReportResponse struct {
	ID            int64     `json:"id"`
	CarePlanID    int64     `json:"care_plan_id"`
	ReportType    string    `json:"report_type"`
	ReportContent string    `json:"report_content"`
	IsCritical    bool      `json:"is_critical"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateCarePlanReportApi creates a new care plan report
// @Summary Create care plan report
// @Description Create a new care plan report
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body CreateCarePlanReportRequest true "Request body"
// @Success 201 {object} Response[CreateCarePlanReportResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/reports [post]
func (server *Server) CreateCarePlanReportApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanReportApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req CreateCarePlanReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanReportApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanReportApi", "Unauthorized access", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}
	employeeID, err := server.store.GetEmployeeIDByUserID(ctx, payload.UserId)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanReportApi", "Failed to get employee ID", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get employee ID")))
		return
	}

	arg := db.CreateCarePlanReportParams{
		CarePlanID:          carePlanID,
		ReportType:          req.ReportType,
		ReportContent:       req.ReportContent,
		IsCritical:          req.IsCritical,
		CreatedByEmployeeID: employeeID,
	}

	report, err := server.store.CreateCarePlanReport(ctx, arg)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "CreateCarePlanReportApi", "Failed to create care plan report", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan report")))
		return
	}

	res := SuccessResponse(CreateCarePlanReportResponse{
		ID:            report.ID,
		CarePlanID:    report.CarePlanID,
		ReportType:    report.ReportType,
		ReportContent: report.ReportContent,
		IsCritical:    report.IsCritical,
		CreatedAt:     report.CreatedAt.Time,
	}, "Care plan report created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListCarePlanReportsRequest represents the request body for listing care plan reports
type ListCarePlanReportsRequest struct {
	pagination.Request
}

// CarePlanReportsResponse represents the response for the ListCarePlanReports API
type ListCarePlanReportsResponse struct {
	ID                 int64     `json:"id"`
	CarePlanID         int64     `json:"care_plan_id"`
	ReportType         string    `json:"report_type"`
	ReportContent      string    `json:"report_content"`
	CreatedByFirstName string    `json:"created_by_first_name"`
	CreatedByLastName  string    `json:"created_by_last_name"`
	IsCritical         bool      `json:"is_critical"`
	CreatedAt          time.Time `json:"created_at"`
}

// ListCarePlanReportsApi retrieves the reports for a given care plan ID
// @Summary List care plan reports
// @Description List all reports for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request query ListCarePlanReportsRequest false "Pagination parameters"
// @Success 200 {object} Response[[]ListCarePlanReportsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/reports [get]
func (server *Server) ListCarePlanReportsApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListCarePlanReportsApi", "Invalid care plan ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req ListCarePlanReportsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "ListCarePlanReportsApi", "Invalid request query", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request query")))
		return
	}

	params := req.GetParams()

	reports, err := server.store.ListCarePlanReports(ctx, db.ListCarePlanReportsParams{
		CarePlanID: carePlanID,
		Limit:      params.Limit,
		Offset:     params.Offset,
	})
	if err != nil {
		server.logBusinessEvent(LogLevelError, "ListCarePlanReportsApi", "Failed to list care plan reports", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list care plan reports")))
		return
	}

	if len(reports) == 0 {
		res := SuccessResponse([]ListCarePlanReportsResponse{}, "No care plan reports found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	response := make([]ListCarePlanReportsResponse, len(reports))
	for i, report := range reports {
		response[i] = ListCarePlanReportsResponse{
			ID:                 report.ID,
			CarePlanID:         report.CarePlanID,
			ReportType:         report.ReportType,
			ReportContent:      report.ReportContent,
			CreatedByFirstName: report.CreatedByFirstName,
			CreatedByLastName:  report.CreatedByLastName,
			IsCritical:         report.IsCritical,
			CreatedAt:          report.CreatedAt.Time,
		}
	}
	pag := pagination.NewResponse(ctx, req.Request, response, reports[0].TotalCount)
	res := SuccessResponse(pag, "Care plan reports retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanReportRequest represents the request body for updating a care plan report
type UpdateCarePlanReportRequest struct {
	ReportType    *string `json:"report_type"`
	ReportContent *string `json:"report_content"`
	IsCritical    *bool   `json:"is_critical"`
}

// UpdateCarePlanReportResponse represents the response for the UpdateCarePlanReport API
type UpdateCarePlanReportResponse struct {
	ID            int64     `json:"id"`
	CarePlanID    int64     `json:"care_plan_id"`
	ReportType    string    `json:"report_type"`
	ReportContent string    `json:"report_content"`
	IsCritical    bool      `json:"is_critical"`
	CreatedAt     time.Time `json:"created_at"`
}

// UpdateCarePlanReportApi updates a care plan report by its ID
// @Summary Update care plan report
// @Description Update a care plan report by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param report_id path int true "Report ID"
// @Param request body UpdateCarePlanReportRequest true "Request body"
// @Success 200 {object} Response[UpdateCarePlanReportResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/reports/{report_id} [put]
func (server *Server) UpdateCarePlanReportApi(ctx *gin.Context) {
	reportID, err := strconv.ParseInt(ctx.Param("report_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanReportApi", "Invalid report ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid report ID")))
		return
	}

	var req UpdateCarePlanReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanReportApi", "Invalid request body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	arg := db.UpdateCarePlanReportParams{
		ID:            reportID,
		ReportType:    req.ReportType,
		ReportContent: req.ReportContent,
		IsCritical:    req.IsCritical,
	}
	report, err := server.store.UpdateCarePlanReport(ctx, arg)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "UpdateCarePlanReportApi", "Failed to update care plan report", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan report")))
		return
	}

	res := SuccessResponse(UpdateCarePlanReportResponse{
		ID:            report.ID,
		CarePlanID:    report.CarePlanID,
		ReportType:    report.ReportType,
		ReportContent: report.ReportContent,
		IsCritical:    report.IsCritical,
		CreatedAt:     report.CreatedAt.Time,
	}, "Care plan report updated successfully")

	ctx.JSON(http.StatusOK, res)
}

// DeleteCarePlanReportApi deletes a care plan report by its ID
// @Summary Delete care plan report
// @Description Delete a care plan report by its ID
// @Tags care_plan
// @Produce json
// @Param report_id path int true "Report ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/reports/{report_id} [delete]
func (server *Server) DeleteCarePlanReportApi(ctx *gin.Context) {
	reportID, err := strconv.ParseInt(ctx.Param("report_id"), 10, 64)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanReportApi", "Invalid report ID", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid report ID")))
		return
	}

	err = server.store.DeleteCarePlanReport(ctx, reportID)
	if err != nil {
		server.logBusinessEvent(LogLevelError, "DeleteCarePlanReportApi", "Failed to delete care plan report", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan report")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan report deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

///////////////////////////////////////////////////////////

// OLD CODE BELOW IS FOR MATURITY MATRIX ASSESSMENT

///////////////////////////////////////////////////////////////////////

// GetClientMaturityMatrixAssessmentResponse represents a response for GetClientMaturityMatrixAssessmentApi
type GetClientMaturityMatrixAssessmentResponse struct {
	ID               int64     `json:"id"`
	ClientID         int64     `json:"client_id"`
	MaturityMatrixID int64     `json:"maturity_matrix_id"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	InitialLevel     int32     `json:"initial_level"`
	CurrentLevel     int32     `json:"current_level"`
	IsActive         bool      `json:"is_active"`
	TopicName        string    `json:"topic_name"`
}

// @Summary Get client maturity matrix assessment
// @Description Get a client maturity matrix assessment
// @Tags maturity_matrix
// @Produce json
// @Param id path int true "Client ID"
// @Param assessment_id path int true "Client maturity matrix assessment ID"
// @Success 200 {object} Response[GetClientMaturityMatrixAssessmentResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/maturity_matrix_assessment/{mma_id} [get]
func (server *Server) GetClientMaturityMatrixAssessmentApi(ctx *gin.Context) {
	mmaID := ctx.Param("assessment_id")
	clientMaturityMatrixAssessmentID, err := strconv.ParseInt(mmaID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	clientMaturityMatrixAssessment, err := server.store.GetClientMaturityMatrixAssessment(ctx, clientMaturityMatrixAssessmentID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(GetClientMaturityMatrixAssessmentResponse{
		ID:               clientMaturityMatrixAssessment.ID,
		ClientID:         clientMaturityMatrixAssessment.ClientID,
		MaturityMatrixID: clientMaturityMatrixAssessment.MaturityMatrixID,
		StartDate:        clientMaturityMatrixAssessment.StartDate.Time,
		EndDate:          clientMaturityMatrixAssessment.EndDate.Time,
		InitialLevel:     clientMaturityMatrixAssessment.InitialLevel,
		CurrentLevel:     clientMaturityMatrixAssessment.CurrentLevel,
		IsActive:         clientMaturityMatrixAssessment.IsActive,
		TopicName:        clientMaturityMatrixAssessment.TopicName,
	}, "Client maturity matrix assessment retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// CreateClientGoalRequest represents a request to create a client goal
type CreateClientGoalRequest struct {
	Description string    `json:"description" binding:"required"`
	TargetDate  time.Time `json:"target_date" binding:"required"`
	StartDate   time.Time `json:"start_date" binding:"required"`
	TargetLevel int32     `json:"target_level" binding:"required"`
	Status      string    `json:"status" binding:"required"`
}

// CreateClientGoalResponse represents a response for CreateClientGoalApi
type CreateClientGoalResponse struct {
	ID                               int64     `json:"id"`
	ClientMaturityMatrixAssessmentID int64     `json:"client_maturity_matrix_assessment_id"`
	Description                      string    `json:"description"`
	Status                           string    `json:"status"`
	TargetLevel                      int32     `json:"target_level"`
	StartDate                        time.Time `json:"start_date"`
	TargetDate                       time.Time `json:"target_date"`
	CompletionDate                   time.Time `json:"completion_date"`
	CreatedAt                        time.Time `json:"created_at"`
}

// @Summary Create client goal
// @Description Create a client goal
// @Tags maturity_matrix
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param assessment_id path int true "Client maturity matrix assessment ID"
// @Param request body CreateClientGoalRequest true "Request body"
// @Success 201 {object} Response[CreateClientGoalResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/maturity_matrix_assessment/{assessment_id}/goals [post]
func (server *Server) CreateClientGoalsApi(ctx *gin.Context) {
	mmaID := ctx.Param("assessment_id")
	clientMaturityMatrixAssessmentID, err := strconv.ParseInt(mmaID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req CreateClientGoalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateClientGoalParams{
		ClientMaturityMatrixAssessmentID: clientMaturityMatrixAssessmentID,
		Description:                      req.Description,
		Status:                           req.Status,
		TargetLevel:                      req.TargetLevel,
		StartDate:                        pgtype.Date{Time: req.StartDate, Valid: true},
		TargetDate:                       pgtype.Date{Time: req.TargetDate, Valid: true},
	}

	clientGoal, err := server.store.CreateClientGoal(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateClientGoalResponse{
		ID:                               clientGoal.ID,
		ClientMaturityMatrixAssessmentID: clientGoal.ClientMaturityMatrixAssessmentID,
		Description:                      clientGoal.Description,
		Status:                           clientGoal.Status,
		TargetLevel:                      clientGoal.TargetLevel,
		StartDate:                        clientGoal.StartDate.Time,
		TargetDate:                       clientGoal.TargetDate.Time,
		CompletionDate:                   clientGoal.CompletionDate.Time,
		CreatedAt:                        clientGoal.CreatedAt.Time,
	}, "Client goal created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListClientGoalsRequest represents a request to list client goals
type ListClientGoalsRequest struct {
	pagination.Request
}

// ListClientGoalsResponse represents a response for ListClientGoalsApi
type ListClientGoalsResponse struct {
	ID                               int64     `json:"id"`
	ClientMaturityMatrixAssessmentID int64     `json:"client_maturity_matrix_assessment_id"`
	Description                      string    `json:"description"`
	Status                           string    `json:"status"`
	TargetLevel                      int32     `json:"target_level"`
	StartDate                        time.Time `json:"start_date"`
	TargetDate                       time.Time `json:"target_date"`
	CompletionDate                   time.Time `json:"completion_date"`
	CreatedAt                        time.Time `json:"created_at"`
}

// @Summary List client goals
// @Description Get a list of client goals
// @Tags maturity_matrix
// @Produce json
// @Param id path int true "Client ID"
// @Param assessment_id path int true "Client maturity matrix assessment ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListClientGoalsResponse]]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/maturity_matrix_assessment/{assessment_id}/goals [get]
func (server *Server) ListClientGoalsApi(ctx *gin.Context) {
	mmaID := ctx.Param("assessment_id")
	clientMaturityMatrixAssessmentID, err := strconv.ParseInt(mmaID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req ListClientGoalsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	params := req.GetParams()

	clientGoals, err := server.store.ListClientGoals(ctx, db.ListClientGoalsParams{
		ClientMaturityMatrixAssessmentID: clientMaturityMatrixAssessmentID,
		Limit:                            params.Limit,
		Offset:                           params.Offset,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if len(clientGoals) == 0 {
		pag := pagination.NewResponse(ctx, req.Request, []ListClientGoalsResponse{}, 0)
		res := SuccessResponse(pag, "No client goals found")
		ctx.JSON(http.StatusOK, res)
		return
	}

	responseClientGoals := make([]ListClientGoalsResponse, len(clientGoals))
	for i, goal := range clientGoals {
		responseClientGoals[i] = ListClientGoalsResponse{
			ID:                               goal.ID,
			ClientMaturityMatrixAssessmentID: goal.ClientMaturityMatrixAssessmentID,
			Description:                      goal.Description,
			Status:                           goal.Status,
			TargetLevel:                      goal.TargetLevel,
			StartDate:                        goal.StartDate.Time,
			TargetDate:                       goal.TargetDate.Time,
			CompletionDate:                   goal.CompletionDate.Time,
			CreatedAt:                        goal.CreatedAt.Time,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, responseClientGoals, clientGoals[0].TotalCount)
	res := SuccessResponse(pag, "Client goals retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// GoalObjectives represents a goal objective
type GoalObjectives struct {
	ID                   int64     `json:"id"`
	ObjectiveDescription string    `json:"objective_description"`
	DueDate              time.Time `json:"due_date"`
	Status               string    `json:"status"`
	CompletionDate       time.Time `json:"completion_date"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// GetClientGoalResponse represents a response for GetClientGoalApi
type GetClientGoalResponse struct {
	ID                               int64            `json:"id"`
	ClientMaturityMatrixAssessmentID int64            `json:"client_maturity_matrix_assessment_id"`
	Description                      string           `json:"description"`
	Status                           string           `json:"status"`
	TargetLevel                      int32            `json:"target_level"`
	StartDate                        time.Time        `json:"start_date"`
	TargetDate                       time.Time        `json:"target_date"`
	CompletionDate                   time.Time        `json:"completion_date"`
	CreatedAt                        time.Time        `json:"created_at"`
	Objectives                       []GoalObjectives `json:"objectives"`
}

// @Summary Get client goal
// @Description Get a client goal
// @Tags maturity_matrix
// @Produce json
// @Param id path int true "Client ID"
// @Param assessment_id path int true "Client maturity matrix assessment ID"
// @Param goal_id path int true "Client goal ID"
// @Success 200 {object} Response[GetClientGoalResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/maturity_matrix_assessment/{assessment_id}/goals/{goal_id} [get]
func (server *Server) GetClientGoalApi(ctx *gin.Context) {
	goalID := ctx.Param("goal_id")
	clientGoalID, err := strconv.ParseInt(goalID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	clientGoal, err := server.store.GetClientGoal(ctx, clientGoalID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	objectives, err := server.store.ListGoalObjectives(ctx, clientGoalID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	responseObjectives := make([]GoalObjectives, len(objectives))
	for i, objective := range objectives {
		responseObjectives[i] = GoalObjectives{
			ID:                   objective.ID,
			ObjectiveDescription: objective.ObjectiveDescription,
			DueDate:              objective.DueDate.Time,
			Status:               objective.Status,
			CompletionDate:       objective.CompletionDate.Time,
			CreatedAt:            objective.CreatedAt.Time,
			UpdatedAt:            objective.UpdatedAt.Time,
		}
	}

	res := SuccessResponse(GetClientGoalResponse{
		ID:                               clientGoal.ID,
		ClientMaturityMatrixAssessmentID: clientGoal.ClientMaturityMatrixAssessmentID,
		Description:                      clientGoal.Description,
		Status:                           clientGoal.Status,
		TargetLevel:                      clientGoal.TargetLevel,
		StartDate:                        clientGoal.StartDate.Time,
		TargetDate:                       clientGoal.TargetDate.Time,
		CompletionDate:                   clientGoal.CompletionDate.Time,
		CreatedAt:                        clientGoal.CreatedAt.Time,
		Objectives:                       responseObjectives,
	}, "Client goal retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// CreateGoalObjectiveRequest represents a request to create a goal objective
type CreateGoalObjectiveRequest struct {
	ObjectiveDescription string    `json:"objective_description" binding:"required"`
	DueDate              time.Time `json:"due_date" binding:"required"`
}

// CreateGoalObjectiveResponse represents a response for CreateGoalObjectiveApi
type CreateGoalObjectiveResponse struct {
	ID                   int64     `json:"id"`
	GoalID               int64     `json:"goal_id"`
	ObjectiveDescription string    `json:"objective_description"`
	Status               string    `json:"status"`
	DueDate              time.Time `json:"due_date"`
}

// @Summary Create goal objective
// @Description Create a goal objective
// @Tags maturity_matrix
// @Accept json
// @Produce json
// @Param goal_id path int true "Client goal ID"
// @Param client_id path int true "Client ID"
// @Param assessment_id path int true "Client maturity matrix assessment ID"
// @Param request body []CreateGoalObjectiveRequest true "Request body"
// @Success 201 {object} Response[[]CreateGoalObjectiveResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/maturity_matrix_assessment/{assessment_id}/goals/{goal_id}/objectives [post]
func (server *Server) CreateGoalObjectiveApi(ctx *gin.Context) {
	goalID := ctx.Param("goal_id")
	clientGoalID, err := strconv.ParseInt(goalID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req []CreateGoalObjectiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := server.store.WithTx(tx)

	responses := make([]CreateGoalObjectiveResponse, 0, len(req))

	for _, objective := range req {
		arg := db.CreateGoalObjectiveParams{
			GoalID:               clientGoalID,
			ObjectiveDescription: objective.ObjectiveDescription,
			Status:               "pending",
			DueDate:              pgtype.Date{Time: objective.DueDate, Valid: true},
		}

		createdObjective, err := qtx.CreateGoalObjective(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		response := CreateGoalObjectiveResponse{
			ID:                   createdObjective.ID,
			GoalID:               createdObjective.GoalID,
			ObjectiveDescription: createdObjective.ObjectiveDescription,
			Status:               createdObjective.Status,
			DueDate:              createdObjective.DueDate.Time,
		}

		responses = append(responses, response)
	}

	err = tx.Commit(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, SuccessResponse(responses, "Goal objectives created successfully"))
}
