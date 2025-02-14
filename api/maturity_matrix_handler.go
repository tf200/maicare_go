package api

import (
	db "maicare_go/db/sqlc"
	"maicare_go/pagination"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

// MatrixAssessment represents a matrix assessment
type MatrixAssessment struct {
	ID           int64 `json:"id"`
	InitialLevel int32 `json:"initial_level"`
}

// CreateClientMaturityMatrixAssessmentRequest represents a request to create a client maturity matrix assessment
type CreateClientMaturityMatrixAssessmentRequest struct {
	Assessments []MatrixAssessment `json:"assessment"`
	StartDate   time.Time          `json:"start_date"`
	EndDate     time.Time          `json:"end_date"`
}

// CreateClientMaturityMatrixAssessmentResponse represents a response for CreateClientMaturityMatrixAssessmentApi
type CreateClientMaturityMatrixAssessmentResponse struct {
	ClientID    int64              `json:"client_id"`
	Assessments []MatrixAssessment `json:"assessment"`
	StartDate   time.Time          `json:"start_date"`
	EndDate     time.Time          `json:"end_date"`
}

// @Summary Create client maturity matrix assessment
// @Description Create a client maturity matrix assessment
// @Tags client
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
			MaturityMatrixID: assessment.ID,
			StartDate:        pgtype.Date{Time: req.StartDate, Valid: true},
			EndDate:          pgtype.Date{Time: req.EndDate, Valid: true},
			InitialLevel:     assessment.InitialLevel,
			CurrentLevel:     assessment.InitialLevel,
		}

		clientAssessments, err := qtx.CreateClientMaturityMatrixAssessment(ctx, arg)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		assessmentRes[i] = MatrixAssessment{
			ID:           clientAssessments.MaturityMatrixID,
			InitialLevel: clientAssessments.InitialLevel,
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
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}, "Client maturity matrix assessment created successfully")
	ctx.JSON(http.StatusCreated, res)
}


type ListClientMaturityMatrixAssessmentsRequest struct {
	pagination.Request
}

type ListClientMaturityMatrixAssessmentsResponse struct {
	ID               int64       `json:"id"`
	ClientID         int64       `json:"client_id"`
	StartDate        pgtype.Date `json:"start_date"`
	EndDate          pgtype.Date `json:"end_date"`
	InitialLevel     int32       `json:"initial_level"`
	CurrentLevel     int32       `json:"current_level"`
	IsActive         bool        `json:"is_active"`
	TopicName        string      `json:"topic_name"`
}

func (server *Server)TestListClientMaturityMatrixAssessmentsApi(t *testing.T){

}