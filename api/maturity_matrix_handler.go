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
