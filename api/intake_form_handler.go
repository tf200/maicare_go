package api

import (
	"fmt"
	_ "maicare_go/pagination"
	clientp "maicare_go/service/client"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create Intake Form
// @Description Create a new intake form with the provided details.
// @Tags Intake Forms
// @Accept json
// @Produce json
// @Param intake_form body clientp.CreateIntakeFormRequest true "Intake Form Details"
// @Success 200 {object} clientp.CreateIntakeFormResponse
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /intake_forms [post]
func (s *Server) CreateIntakeFormApi(ctx *gin.Context) {
	var req clientp.CreateIntakeFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	res, err := s.businessService.ClientService.CreateIntakeForm(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary List Intake Forms
// @Description Retrieve a paginated list of intake forms with optional search and sorting.
// @Tags Intake Forms
// @Accept json
// @Produce json
// @Param search query string false "Search term to filter intake forms"
// @Param sort_order query string false "Sort order (asc or desc)"
// @Param page query int false "Page number for pagination"
// @Param page_size query int false "Number of items per page for pagination"
// @Success 200 {object} pagination.Response[clientp.ListIntakeFormsResponse]
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /intake_forms [get]
func (s *Server) ListIntakeFormsApi(ctx *gin.Context) {
	var req clientp.ListIntakeFormsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid query parameters")))
		return
	}

	res, err := s.businessService.ClientService.ListIntakeForms(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Complete Intake Form
// @Description Complete an intake form and move accepted client to waiting list.
// @Tags Intake Forms
// @Accept json
// @Produce json
// @Param id path uuid true "Intake Form ID"
// @Param request body clientp.CompleteIntakeFormRequest true "Complete Intake Form Request"
// @Success 200 {object} Response[clientp.CompleteIntakeFormResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /intake_form/{id}/outcome [put]
func (s *Server) CompleteIntakeFormApi(ctx *gin.Context) {
	intakeID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid intake form id")))
		return
	}

	var req clientp.CompleteIntakeFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	res, err := s.businessService.ClientService.CompleteIntakeForm(ctx, &req, intakeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, SuccessResponse(res, "Intake completed"))
}
