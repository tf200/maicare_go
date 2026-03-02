package api

import (
	"errors"
	"fmt"
	_ "maicare_go/pagination"
	clientp "maicare_go/service/client"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

	result, err := s.businessService.ClientService.CreateIntakeForm(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Intake Form created successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary List Intake Forms
// @Description Retrieve a paginated list of intake forms with optional search and sorting.
// @Tags Intake Forms
// @Accept json
// @Produce json
// @Param search query string false "Search term to filter intake forms"
// @Param status query string false "Intake status (intake_conclusion) filter"
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

	result, err := s.businessService.ClientService.ListIntakeForms(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Intake Forms listed successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Intake Form Totals
// @Description Retrieve totals for further investigation intakes and intakes without goals.
// @Tags Intake Forms
// @Produce json
// @Success 200 {object} Response[clientp.GetIntakeFormTotalsResponse]
// @Failure 500 {object} Response[any]
// @Router /intake_forms/totals [get]
func (s *Server) GetIntakeFormTotalsApi(ctx *gin.Context) {
	result, err := s.businessService.ClientService.GetIntakeFormTotals(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Intake Form totals retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Intake Form
// @Description Retrieve an intake form by ID with related details.
// @Tags Intake Forms
// @Produce json
// @Param id path uuid true "Intake Form ID"
// @Success 200 {object} Response[clientp.GetIntakeFormResponse]
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /intake_forms/{id} [get]
func (s *Server) GetIntakeFormApi(ctx *gin.Context) {
	intakeFormID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	result, err := s.businessService.ClientService.GetIntakeForm(ctx, intakeFormID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Intake Form retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update Intake Form
// @Description Partially update editable intake form fields.
// @Tags Intake Forms
// @Accept json
// @Produce json
// @Param id path uuid true "Intake Form ID"
// @Param request body clientp.UpdateIntakeFormRequest true "Intake form update payload"
// @Success 200 {object} Response[clientp.UpdateIntakeFormResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 409 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /intake_forms/{id} [patch]
func (s *Server) UpdateIntakeFormApi(ctx *gin.Context) {
	intakeFormID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.UpdateIntakeFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	result, err := s.businessService.ClientService.UpdateIntakeForm(ctx, intakeFormID, &req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("intake form not found")))
			return
		}
		if errors.Is(err, clientp.ErrIntakeFormUpdateBlockedByActiveClient) || errors.Is(err, clientp.ErrIntakeFormUpdateConflict) {
			ctx.JSON(http.StatusConflict, errorResponse(err))
			return
		}
		if errors.Is(err, clientp.ErrNoIntakeFormFieldsToUpdate) || errors.Is(err, clientp.ErrInvalidIntakeFormClearField) || strings.Contains(err.Error(), "invalid intake form update") {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Intake Form updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Replace Intake Form Goals
// @Description Replace intake maturity assessments (goals) for an intake form.
// @Tags Intake Forms
// @Accept json
// @Produce json
// @Param id path uuid true "Intake Form ID"
// @Param request body clientp.CreateIntakeFormGoalsRequest true "Goals payload"
// @Success 200 {object} Response[clientp.CreateIntakeFormGoalsResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 409 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /intake_forms/{id}/goals [put]
func (s *Server) CreateIntakeFormGoalsApi(ctx *gin.Context) {
	intakeFormID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.CreateIntakeFormGoalsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	result, err := s.businessService.ClientService.CreateIntakeFormGoals(ctx, intakeFormID, &req)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("intake form not found")))
			return
		}
		if errors.Is(err, clientp.ErrIntakeGoalsUpdateBlockedByActiveClient) {
			ctx.JSON(http.StatusConflict, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Intake Form goals replaced successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update Intake Conclusion
// @Description Accept or refuse an intake by updating intake conclusion.
// @Tags Intake Forms
// @Accept json
// @Produce json
// @Param id path uuid true "Intake Form ID"
// @Param request body clientp.UpdateIntakeConclusionRequest true "Conclusion update payload"
// @Success 200 {object} Response[clientp.UpdateIntakeConclusionResponse]
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /intake_forms/{id}/conclusion [patch]
func (s *Server) UpdateIntakeConclusionApi(ctx *gin.Context) {
	intakeFormID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.UpdateIntakeConclusionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	result, err := s.businessService.ClientService.UpdateIntakeConclusion(ctx, intakeFormID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "invalid decision") {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(result, "Intake conclusion updated successfully")
	ctx.JSON(http.StatusOK, res)
}
