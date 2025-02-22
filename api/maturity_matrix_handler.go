package api

import (
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
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

// MatrixAssessment represents a matrix assessment
type MatrixAssessment struct {
	ID               int64     `json:"id"`
	MaturityMatrixID int64     `json:"maturity_matrix_id"`
	InitialLevel     int32     `json:"initial_level"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
}

// CreateClientMaturityMatrixAssessmentRequest represents a request to create a client maturity matrix assessment
type CreateClientMaturityMatrixAssessmentRequest struct {
	Assessments []MatrixAssessment `json:"assessment" binding:"required"`
}

// CreateClientMaturityMatrixAssessmentResponse represents a response for CreateClientMaturityMatrixAssessmentApi
type CreateClientMaturityMatrixAssessmentResponse struct {
	ClientID    int64              `json:"client_id"`
	Assessments []MatrixAssessment `json:"assessment"`
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
// @Router /clients/{id}/maturity_matrix_assessment [post]
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

	tx, err := server.store.ConnPool.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	defer tx.Rollback(ctx)

	qtx := server.store.WithTx(tx)

	assessmentRes := make([]MatrixAssessment, len(req.Assessments))

	for i, assessment := range req.Assessments {
		arg := db.CreateClientMaturityMatrixAssessmentParams{
			ClientID:         clientID,
			MaturityMatrixID: assessment.MaturityMatrixID,
			StartDate:        pgtype.Date{Time: assessment.StartDate, Valid: true},
			EndDate:          pgtype.Date{Time: assessment.EndDate, Valid: true},
			InitialLevel:     assessment.InitialLevel,
			CurrentLevel:     assessment.InitialLevel,
		}

		clientAssessments, err := qtx.CreateClientMaturityMatrixAssessment(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		assessmentRes[i] = MatrixAssessment{
			ID:               clientAssessments.ID,
			MaturityMatrixID: clientAssessments.MaturityMatrixID,
			InitialLevel:     clientAssessments.InitialLevel,
			StartDate:        clientAssessments.StartDate.Time,
			EndDate:          clientAssessments.EndDate.Time,
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateClientMaturityMatrixAssessmentResponse{
		ClientID:    clientID,
		Assessments: assessmentRes,
	}, "Client maturity matrix assessment created successfully")
	ctx.JSON(http.StatusCreated, res)
}

// ListClientMaturityMatrixAssessmentsRequest represents a request to list client maturity matrix assessments
type ListClientMaturityMatrixAssessmentsRequest struct {
	pagination.Request
}

// ListClientMaturityMatrixAssessmentsResponse represents a response for ListClientMaturityMatrixAssessmentsApi
type ListClientMaturityMatrixAssessmentsResponse struct {
	ID                 int64       `json:"id"`
	MatrixAssessmentID int64       `json:"matrix_assessment_id"`
	ClientID           int64       `json:"client_id"`
	StartDate          pgtype.Date `json:"start_date"`
	EndDate            pgtype.Date `json:"end_date"`
	InitialLevel       int32       `json:"initial_level"`
	CurrentLevel       int32       `json:"current_level"`
	IsActive           bool        `json:"is_active"`
	TopicName          string      `json:"topic_name"`
}

// @Summary List client maturity matrix assessments
// @Description Get a list of client maturity matrix assessments
// @Tags maturity_matrix
// @Produce json
// @Param id path int true "Client ID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} Response[pagination.Response[ListClientMaturityMatrixAssessmentsResponse]]
// @Failure 400 {object} Response[any] "Bad request"
// @Failure 401 {object} Response[any] "Unauthorized"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /clients/{id}/maturity_matrix_assessment [get]
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
			ID:                 assessment.ID,
			MatrixAssessmentID: assessment.MaturityMatrixID,
			TopicName:          assessment.TopicName,
			ClientID:           assessment.ClientID,
			StartDate:          assessment.StartDate,
			EndDate:            assessment.EndDate,
			InitialLevel:       assessment.InitialLevel,
			CurrentLevel:       assessment.CurrentLevel,
			IsActive:           assessment.IsActive,
		}
	}

	pag := pagination.NewResponse(ctx, req.Request, responseClientAssessments, clientAssessments[0].TotalCount)
	res := SuccessResponse(pag, "Client maturity matrix assessments retrieved successfully")

	ctx.JSON(http.StatusOK, res)

}

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
	Status               string    `json:"status" binding:"required"`
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
// @Param request body CreateGoalObjectiveRequest true "Request body"
// @Success 201 {object} Response[CreateGoalObjectiveResponse]
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

	var req CreateGoalObjectiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateGoalObjectiveParams{
		GoalID:               clientGoalID,
		ObjectiveDescription: req.ObjectiveDescription,
		Status:               req.Status,
		DueDate:              pgtype.Date{Time: req.DueDate, Valid: true},
	}

	goalObjective, err := server.store.CreateGoalObjective(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(CreateGoalObjectiveResponse{
		ID:                   goalObjective.ID,
		GoalID:               goalObjective.GoalID,
		ObjectiveDescription: goalObjective.ObjectiveDescription,
		Status:               goalObjective.Status,
		DueDate:              goalObjective.DueDate.Time,
	}, "Goal objective created successfully")

	ctx.JSON(http.StatusCreated, res)
}
