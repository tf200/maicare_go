package api

import (
	"fmt"
	_ "maicare_go/pagination"
	"net/http"
	"strconv"

	"maicare_go/service/care"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// @Summary List all maturity matrix
// @Description Get a list of all maturity matrix
// @Tags maturity_matrix
// @Produce json
// @Success 200 {object} Response[[]care.ListCarePlanTopics]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /maturity_matrix [get]
func (server *Server) ListMaturityMatrixApi(ctx *gin.Context) {
	result, err := server.businessService.CarePlanService.ListCarePlanTopics(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Maturity matrix retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// @Summary Create client maturity matrix assessment
// @Description Create a client maturity matrix assessment
// @Tags maturity_matrix
// @Accept json
// @Produce json
// @Param id path int true "Client ID"
// @Param request body care.CreateClientCarePlanRequest true "Request body"
// @Success 201 {object} Response[care.CreateClientCarePlanResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/assessments [post]
func (server *Server) CreateClientMaturityMatrixAssessmentApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid client ID")))
		return
	}

	var req care.CreateClientCarePlanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	payload, err := GetAuthPayload(ctx)
	if err != nil {
		server.logger.Error("failed to get auth payload", zap.Error(err))
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized")))
		return
	}

	carePlan, err := server.businessService.CarePlanService.CreateClientCarePlan(
		ctx, clientID, payload.EmployeeID, &req,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create client care plan")))
		return
	}

	res := SuccessResponse(carePlan, "Client maturity matrix assessment created successfully")
	ctx.JSON(http.StatusCreated, res)
}

// @Summary List client maturity matrix assessments
// @Description Get a list of client maturity matrix assessments
// @Tags care_plan
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[care.ListClientCarePlansResponse]]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/assessments [get]
func (server *Server) ListClientMaturityMatrixAssessmentsApi(ctx *gin.Context) {
	id := ctx.Param("id")
	clientID, err := uuid.Parse(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req care.ListClientCarePlansRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := server.businessService.CarePlanService.ListClientCarePlans(ctx, clientID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Client maturity matrix assessments retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// ============================ CarePlan Overview ===========================

// GetCarePlanOverviewApi retrieves the care plan overview for a given assessment ID
// @Summary Get care plan overview
// @Description Get the care plan overview for a given assessment ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[care.GetCarePlanOverviewResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id} [get]
func (server *Server) GetCarePlanOverviewApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	result, err := server.businessService.CarePlanService.GetCarePlanOverview(ctx, carePlanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Care plan overview retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanOverviewApi updates the care plan overview for a given care plan ID
// @Summary Update care plan overview
// @Description Update the care plan overview for a given care plan ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body care.UpdateCarePlanOverviewRequest true "Request body"
// @Success 200 {object} Response[care.UpdateCarePlanOverviewResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id} [put]
func (server *Server) UpdateCarePlanOverviewApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req care.UpdateCarePlanOverviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	result, err := server.businessService.CarePlanService.UpdateCarePlanOverview(ctx, carePlanID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Care plan overview updated successfully")

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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	err = server.businessService.CarePlanService.DeleteCarePlan(ctx, carePlanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Care plan deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// =========================== Care plan objectives and actions ===========================

// CreateCarePlanObjectiveApi creates a new care plan objective
// @Summary Create care plan objective
// @Description Create a new care plan objective
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body care.CreateCarePlanObjectiveRequest true "Request body"
// @Success 201 {object} Response[care.CreateCarePlanObjectiveResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/objectives [post]
func (server *Server) CreateCarePlanObjectiveApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req care.CreateCarePlanObjectiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	objective, err := server.businessService.CarePlanService.CreateCarePlanObjective(ctx, carePlanID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan objective")))
		return
	}

	res := SuccessResponse(objective, "Care plan objective created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GetCarePlanObjectivesApi retrieves the care plan objectives for a given care plan ID
// @Summary Get care plan objectives
// @Description Get the care plan objectives for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[care.GetCarePlanObjectivesResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/objectives [get]
func (server *Server) GetCarePlanObjectivesApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	response, err := server.businessService.CarePlanService.GetCarePlanObjectivesAndActions(ctx, carePlanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan objectives and actions")))
		return
	}

	res := SuccessResponse(response, "Care plan objectives retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update care plan objective
// @Description Update a care plan objective by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param objective_id path int true "Objective ID"
// @Param request body care.UpdateCarePlanObjectiveRequest true "Request body"
// @Success 200 {object} Response[care.UpdateCarePlanObjectiveResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /objectives/{objective_id} [put]
func (server *Server) UpdateCarePlanObjectiveApi(ctx *gin.Context) {
	objectiveId, err := strconv.ParseInt(ctx.Param("objective_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid objective ID")))
		return
	}

	var req care.UpdateCarePlanObjectiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	response, err := server.businessService.CarePlanService.UpdateCarePlanObjective(ctx, objectiveId, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan objective")))
		return
	}

	res := SuccessResponse(response, "Care plan objective updated successfully")

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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid objective ID")))
		return
	}

	err = server.businessService.CarePlanService.DeleteCarePlanObjective(ctx, objectiveId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan objective")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan objective deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// CreateCarePlanActionsApi creates a new care plan action
// @Summary Create care plan action
// @Description Create a new care plan action
// @Tags care_plan
// @Accept json
// @Produce json
// @Param objective_id path int true "Objective ID"
// @Param request body care.CreateCarePlanActionsRequest true "Request body"
// @Success 201 {object} Response[care.CreateCarePlanActionsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /objectives/{objective_id}/actions [post]
func (server *Server) CreateCarePlanActionsApi(ctx *gin.Context) {
	objectiveID, err := strconv.ParseInt(ctx.Param("objective_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid objective ID")))
		return
	}

	var req care.CreateCarePlanActionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	action, err := server.businessService.CarePlanService.CreateCarePlanAction(ctx, objectiveID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan action")))
		return
	}

	res := SuccessResponse(action, "Care plan action created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// UpdateCarePlanActionsApi updates a care plan action by its ID
// @Summary Update care plan action
// @Description Update a care plan action by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param action_id path int true "Action ID"
// @Param request body care.UpdateCarePlanActionsRequest true "Request body"
// @Success 200 {object} Response[care.UpdateCarePlanActionsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /actions/{action_id} [put]
func (server *Server) UpdateCarePlanActionsApi(ctx *gin.Context) {
	actionID, err := strconv.ParseInt(ctx.Param("action_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid action ID")))
		return
	}

	var req care.UpdateCarePlanActionsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	action, err := server.businessService.CarePlanService.UpdateCarePlanAction(ctx, actionID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan action")))
		return
	}

	res := SuccessResponse(action, "Care plan action updated successfully")

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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid action ID")))
		return
	}

	err = server.businessService.CarePlanService.DeleteCarePlanAction(ctx, actionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan action")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan action deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// =========================== Care plan interventions ===========================

// CreateCarePlanInterventionApi creates a new care plan intervention
// @Summary Create care plan intervention
// @Description Create a new care plan intervention
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body care.CreateCarePlanInterventionRequest true "Request body"
// @Success 201 {object} Response[care.CreateCarePlanInterventionResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/interventions [post]
func (server *Server) CreateCarePlanInterventionApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req care.CreateCarePlanInterventionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	intervention, err := server.businessService.CarePlanService.CreateCarePlanIntervention(ctx, carePlanID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan intervention")))
		return
	}

	res := SuccessResponse(intervention, "Care plan intervention created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GetCarePlanInterventionsApi retrieves the care plan interventions for a given care plan ID
// @Summary Get care plan interventions
// @Description Get the care plan interventions for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[care.GetCarePlanInterventionsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/interventions [get]
func (server *Server) GetCarePlanInterventionsApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	response, err := server.businessService.CarePlanService.GetCarePlanInterventions(ctx, carePlanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan interventions")))
		return
	}

	res := SuccessResponse(response, "Care plan interventions retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update care plan intervention
// @Description Update a care plan intervention by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param intervention_id path int true "Intervention ID"
// @Param request body care.UpdateCarePlanInterventionRequest true "Request body"
// @Success 200 {object} Response[care.UpdateCarePlanInterventionResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /interventions/{intervention_id} [put]
func (server *Server) UpdateCarePlanInterventionApi(ctx *gin.Context) {
	interventionID, err := strconv.ParseInt(ctx.Param("intervention_id"), 10, 64)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid intervention ID")))
		return
	}

	var req care.UpdateCarePlanInterventionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	intervention, err := server.businessService.CarePlanService.UpdateCarePlanIntervention(ctx, interventionID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan intervention")))
		return
	}

	res := SuccessResponse(intervention, "Care plan intervention updated successfully")

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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid intervention ID")))
		return
	}

	err = server.businessService.CarePlanService.DeleteCarePlanIntervention(ctx, interventionID)
	if err != nil {
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

// CreateCarePlanSuccessMetricsApi creates a new care plan success metric
// @Summary Create care plan success metric
// @Description Create a new care plan success metric
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body care.CreateCarePlanSuccessMetricsRequest true "Request body"
// @Success 201 {object} Response[care.CreateCarePlanSuccessMetricsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/success_metrics [post]
func (server *Server) CreateCarePlanSuccessMetricsApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req care.CreateCarePlanSuccessMetricsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	successMetric, err := server.businessService.CarePlanService.CreateCarePlanSuccessMetric(ctx, carePlanID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan success metric")))
		return
	}

	res := SuccessResponse(successMetric, "Care plan success metric created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GetCarePlanSuccessMetricsApi retrieves the success metrics for a given care plan ID
// @Summary Get care plan success metrics
// @Description Get the success metrics for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[[]care.GetCarePlanSuccessMetricsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/success_metrics [get]
func (server *Server) GetCarePlanSuccessMetricsApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	response, err := server.businessService.CarePlanService.GetCarePlanSuccessMetrics(ctx, carePlanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan success metrics")))
		return
	}

	res := SuccessResponse(response, "Care plan success metrics retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanSuccessMetricsApi updates a care plan success metric by its ID
// @Summary Update care plan success metric
// @Description Update a care plan success metric by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param metric_id path int true "Metric ID"
// @Param request body care.UpdateCarePlanSuccessMetricsRequest true "Request body"
// @Success 200 {object} Response[care.UpdateCarePlanSuccessMetricsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /success_metrics/{metric_id} [put]
func (server *Server) UpdateCarePlanSuccessMetricsApi(ctx *gin.Context) {
	metricID, err := strconv.ParseInt(ctx.Param("metric_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid metric ID")))
		return
	}

	var req care.UpdateCarePlanSuccessMetricsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	successMetric, err := server.businessService.CarePlanService.UpdateCarePlanSuccessMetric(ctx, metricID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan success metric")))
		return
	}

	res := SuccessResponse(successMetric, "Care plan success metric updated successfully")

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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid metric ID")))
		return
	}

	err = server.businessService.CarePlanService.DeleteCarePlanSuccessMetric(ctx, metricID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan success metric")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan success metric deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// =========================== Care plan risks ===========================

// CreateCarePlanRisksApi creates a new care plan risk
// @Summary Create care plan risk
// @Description Create a new care plan risk
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body care.CreateCarePlanRisksRequest true "Request body"
// @Success 201 {object} Response[care.CreateCarePlanRisksResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/risks [post]
func (server *Server) CreateCarePlanRisksApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req care.CreateCarePlanRisksRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	risk, err := server.businessService.CarePlanService.CreateCarePlanRisk(ctx, carePlanID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan risk")))
		return
	}

	res := SuccessResponse(risk, "Care plan risk created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GetCarePlanRisksApi retrieves the risks associated with a given care plan ID
// @Summary Get care plan risks
// @Description Get the risks associated with a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[[]care.GetCarePlanRisksResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plan/{care_plan_id}/risks [get]
func (server *Server) GetCarePlanRisksApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}
	response, err := server.businessService.CarePlanService.GetCarePlanRisks(ctx, carePlanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan risks")))
		return
	}

	res := SuccessResponse(response, "Care plan risks retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanRisksApi updates a care plan risk by its ID
// @Summary Update care plan risk
// @Description Update a care plan risk by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param risk_id path int true "Risk ID"
// @Param request body care.UpdateCarePlanRisksRequest true "Request body"
// @Success 200 {object} Response[care.UpdateCarePlanRisksResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /risks/{risk_id} [put]
func (server *Server) UpdateCarePlanRisksApi(ctx *gin.Context) {
	riskID, err := strconv.ParseInt(ctx.Param("risk_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid risk ID")))
		return
	}

	var req care.UpdateCarePlanRisksRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	risk, err := server.businessService.CarePlanService.UpdateCarePlanRisk(ctx, riskID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan risk")))
		return
	}

	res := SuccessResponse(risk, "Care plan risk updated successfully")

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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid risk ID")))
		return
	}

	err = server.businessService.CarePlanService.DeleteCarePlanRisk(ctx, riskID)
	if err != nil {
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

// CreateCareplanSupportNetworkApi creates a new care plan support network
// @Summary Create care plan support network
// @Description Create a new care plan support network
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body care.CreateCarePlanSupportNetworkRequest true "Request body"
// @Success 201 {object} Response[care.CreateCarePlanSupportNetworkResponse]
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

	var req care.CreateCarePlanSupportNetworkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	result, err := server.businessService.CarePlanService.CreateCarePlanSupportNetwork(ctx, carePlanID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan support network")))
		return
	}

	res := SuccessResponse(result, "Care plan support network created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// @Summary Get care plan support network
// @Description Get the support network for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[[]care.GetCarePlanSupportNetworkResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/support_network [get]
func (server *Server) GetCarePlanSupportNetworkApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}
	response, err := server.businessService.CarePlanService.GetCarePlanSupportNetwork(ctx, carePlanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan support network")))
		return
	}

	res := SuccessResponse(response, "Care plan support network retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanSupportNetworkApi updates a care plan support network by its ID
// @Summary Update care plan support network
// @Description Update a care plan support network by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param support_network_id path int true "Support Network ID"
// @Param request body care.UpdateCarePlanSupportNetworkRequest true "Request body"
// @Success 200 {object} Response[care.UpdateCarePlanSupportNetworkResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /support_network/{support_network_id} [put]
func (server *Server) UpdateCarePlanSupportNetworkApi(ctx *gin.Context) {
	supportNetworkID, err := strconv.ParseInt(ctx.Param("support_network_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid support network ID")))
		return
	}

	var req care.UpdateCarePlanSupportNetworkRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	result, err := server.businessService.CarePlanService.UpdateCarePlanSupportNetwork(ctx, supportNetworkID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan support network")))
		return
	}

	res := SuccessResponse(result, "Care plan support network updated successfully")

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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid support network ID")))
		return
	}

	err = server.businessService.CarePlanService.DeleteCarePlanSupportNetwork(ctx, supportNetworkID)
	if err != nil {
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

// CreateCarePlanResourcesApi creates a new care plan resource
// @Summary Create care plan resource
// @Description Create a new care plan resource
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body care.CreateCarePlanResourcesRequest true "Request body"
// @Success 201 {object} Response[care.CreateCarePlanResourcesResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plan/{care_plan_id}/resources [post]
func (server *Server) CreateCarePlanResourcesApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req care.CreateCarePlanResourcesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	result, err := server.businessService.CarePlanService.CreateCarePlanResource(ctx, carePlanID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan resource")))
		return
	}

	res := SuccessResponse(result, "Care plan resource created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// GetCarePlanResourcesApi retrieves the resources for a given care plan ID
// @Summary Get care plan resources
// @Description Get the resources for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Success 200 {object} Response[[]care.GetCarePlanResourcesResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plan/{care_plan_id}/resources [get]
func (server *Server) GetCarePlanResourcesApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	response, err := server.businessService.CarePlanService.GetCarePlanResources(ctx, carePlanID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to get care plan resources")))
		return
	}

	res := SuccessResponse(response, "Care plan resources retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanResourcesApi updates a care plan resource by its ID
// @Summary Update care plan resource
// @Description Update a care plan resource by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param resource_id path int true "Resource ID"
// @Param request body care.UpdateCarePlanResourcesRequest true "Request body"
// @Success 200 {object} Response[care.UpdateCarePlanResourcesResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /resources/{resource_id} [put]
func (server *Server) UpdateCarePlanResourcesApi(ctx *gin.Context) {
	resourceID, err := strconv.ParseInt(ctx.Param("resource_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid resource ID")))
		return
	}

	var req care.UpdateCarePlanResourcesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	result, err := server.businessService.CarePlanService.UpdateCarePlanResource(ctx, resourceID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan resource")))
		return
	}

	res := SuccessResponse(result, "Care plan resource updated successfully")

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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid resource ID")))
		return
	}

	err = server.businessService.CarePlanService.DeleteCarePlanResource(ctx, resourceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan resources")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan resources deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

//////////////////////

////////////////////// Care Plan Reports /////////

////////////////////////////

// CreateCarePlanReportApi creates a new care plan report
// @Summary Create care plan report
// @Description Create a new care plan report
// @Tags care_plan
// @Accept json
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request body care.CreateCarePlanReportRequest true "Request body"
// @Success 201 {object} Response[care.CreateCarePlanReportResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/reports [post]
func (server *Server) CreateCarePlanReportApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req care.CreateCarePlanReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("unauthorized access")))
		return
	}

	report, err := server.businessService.CarePlanService.CreateCarePlanReport(ctx, carePlanID, payload.EmployeeID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to create care plan report")))
		return
	}

	res := SuccessResponse(report, "Care plan report created successfully")

	ctx.JSON(http.StatusCreated, res)
}

// ListCarePlanReportsApi retrieves the reports for a given care plan ID
// @Summary List care plan reports
// @Description List all reports for a given care plan ID
// @Tags care_plan
// @Produce json
// @Param care_plan_id path int true "Care Plan ID"
// @Param request query care.ListCarePlanReportsRequest false "Pagination parameters"
// @Success 200 {object} Response[[]care.ListCarePlanReportsResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/{care_plan_id}/reports [get]
func (server *Server) ListCarePlanReportsApi(ctx *gin.Context) {
	carePlanID, err := strconv.ParseInt(ctx.Param("care_plan_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid care plan ID")))
		return
	}

	var req care.ListCarePlanReportsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request query")))
		return
	}

	result, err := server.businessService.CarePlanService.ListCarePlanReports(ctx, carePlanID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list care plan reports")))
		return
	}

	res := SuccessResponse(result, "Care plan reports retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// UpdateCarePlanReportApi updates a care plan report by its ID
// @Summary Update care plan report
// @Description Update a care plan report by its ID
// @Tags care_plan
// @Accept json
// @Produce json
// @Param report_id path int true "Report ID"
// @Param request body care.UpdateCarePlanReportRequest true "Request body"
// @Success 200 {object} Response[care.UpdateCarePlanReportResponse]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /care_plans/reports/{report_id} [put]
func (server *Server) UpdateCarePlanReportApi(ctx *gin.Context) {
	reportID, err := strconv.ParseInt(ctx.Param("report_id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid report ID")))
		return
	}

	var req care.UpdateCarePlanReportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	report, err := server.businessService.CarePlanService.UpdateCarePlanReport(ctx, reportID, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to update care plan report")))
		return
	}

	res := SuccessResponse(report, "Care plan report updated successfully")

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
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid report ID")))
		return
	}
	err = server.businessService.CarePlanService.DeleteCarePlanReport(ctx, reportID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to delete care plan report")))
		return
	}

	res := SuccessResponse[any](nil, "Care plan report deleted successfully")
	ctx.JSON(http.StatusOK, res)
}
