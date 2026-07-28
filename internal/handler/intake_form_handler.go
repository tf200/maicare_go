package handler

import (
	"net/http"
	"strings"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"
	"maicare_go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type IntakeFormHandler struct {
	service domain.IntakeFormService
}

func NewIntakeFormHandler(service domain.IntakeFormService) *IntakeFormHandler {
	return &IntakeFormHandler{service: service}
}

// CreateIntakeForm creates a new intake form.
// @Summary Create a new intake form
// @Tags intake_forms
// @Accept json
// @Produce json
// @Param request body createIntakeFormRequest true "Create Intake Form Request"
// @Success 200 {object} httpapi.Envelope[intakeFormResponse]
// @Router /intake_forms [post]
func (h *IntakeFormHandler) CreateIntakeForm(ctx *gin.Context) {
	var req createIntakeFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.CreateIntakeForm(ctx.Request.Context(), toCreateIntakeFormParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create intake form", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toIntakeFormResponse(result), "Intake form created successfully"))
}

// ListIntakeForms returns a list of intake forms.
// @Summary List intake forms
// @Tags intake_forms
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param search query string false "Search by client name"
// @Param status query string false "Filter by status"
// @Param sort_order query string false "Sort order (asc/desc)"
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[[]intakeFormListItemResponse]]
// @Router /intake_forms [get]
func (h *IntakeFormHandler) ListIntakeForms(ctx *gin.Context) {
	var req listIntakeFormsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid query params", err.Error()))
		return
	}

	result, err := h.service.ListIntakeForms(ctx.Request.Context(), toListIntakeFormsParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list intake forms", err.Error()))
		return
	}

	items := make([]intakeFormListItemResponse, len(result.Items))
	for i, item := range result.Items {
		items[i] = toIntakeFormListItemResponse(item)
	}

	ctx.JSON(http.StatusOK, httpapi.OK(httpapi.NewPageResponse(ctx, req.PageRequest, items, result.TotalCount), "Intake forms retrieved successfully"))
}

// GetIntakeFormTotals returns totals for intake forms.
// @Summary Get intake form totals
// @Tags intake_forms
// @Produce json
// @Success 200 {object} httpapi.Envelope[intakeFormTotalsResponse]
// @Router /intake_forms/totals [get]
func (h *IntakeFormHandler) GetIntakeFormTotals(ctx *gin.Context) {
	result, err := h.service.GetIntakeFormTotals(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get intake form totals", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toIntakeFormTotalsResponse(result), "Intake form totals retrieved successfully"))
}

// GetIntakeForm returns an intake form by ID.
// @Summary Get an intake form by ID
// @Tags intake_forms
// @Produce json
// @Param id path string true "Intake Form ID"
// @Success 200 {object} httpapi.Envelope[intakeFormDetailResponse]
// @Router /intake_forms/{id} [get]
func (h *IntakeFormHandler) GetIntakeForm(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid intake form ID", err.Error()))
		return
	}

	result, err := h.service.GetIntakeForm(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get intake form", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toIntakeFormDetailResponse(result), "Intake form retrieved successfully"))
}

// UpdateIntakeForm updates an intake form.
// @Summary Update an intake form
// @Tags intake_forms
// @Accept json
// @Produce json
// @Param id path string true "Intake Form ID"
// @Param request body updateIntakeFormRequest true "Update Intake Form Request"
// @Success 200 {object} httpapi.Envelope[intakeFormResponse]
// @Router /intake_forms/{id} [patch]
func (h *IntakeFormHandler) UpdateIntakeForm(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid intake form ID", err.Error()))
		return
	}

	var req updateIntakeFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.UpdateIntakeForm(ctx.Request.Context(), toUpdateIntakeFormParams(id, req))
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			ctx.JSON(http.StatusNotFound, httpapi.Fail("intake form not found", err.Error()))
		case domain.ErrIntakeFormUpdateBlockedByActiveClient, domain.ErrIntakeFormUpdateConflict:
			ctx.JSON(http.StatusConflict, httpapi.Fail("intake form update conflict", err.Error()))
		case domain.ErrNoIntakeFormFieldsToUpdate, domain.ErrInvalidIntakeFormClearField, domain.ErrInvalidIntakeFormUpdate:
			ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid intake form update", err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update intake form", err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toIntakeFormResponse(result), "Intake form updated successfully"))
}

// GenerateIntakeGoals generates intake goals for a topic.
// @Summary Generate intake goals
// @Tags intake_forms
// @Accept json
// @Produce json
// @Param id path string true "Intake Form ID"
// @Param request body generateIntakeGoalsRequest true "Generate Intake Goals Request"
// @Success 200 {object} httpapi.Envelope[generateIntakeGoalsResponse]
// @Router /intake_forms/{id}/generate_goals [post]
func (h *IntakeFormHandler) GenerateIntakeGoals(ctx *gin.Context) {
	intakeFormID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid intake form ID", err.Error()))
		return
	}

	var req generateIntakeGoalsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.GenerateIntakeGoals(ctx.Request.Context(), toGenerateIntakeGoalsParams(intakeFormID, req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to generate intake goals", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toGenerateIntakeGoalsResponse(result), "Intake goals generated successfully"))
}

// ReplaceIntakeFormGoals replaces goals for an intake form.
// @Summary Replace intake form goals
// @Tags intake_forms
// @Accept json
// @Produce json
// @Param id path string true "Intake Form ID"
// @Param request body createIntakeFormGoalsRequest true "Replace Intake Form Goals Request"
// @Success 200 {object} httpapi.Envelope[intakeFormGoalsResponse]
// @Router /intake_forms/{id}/goals [put]
func (h *IntakeFormHandler) ReplaceIntakeFormGoals(ctx *gin.Context) {
	intakeFormID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid intake form ID", err.Error()))
		return
	}

	var req createIntakeFormGoalsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.ReplaceIntakeFormGoals(ctx.Request.Context(), intakeFormID, toReplaceIntakeFormGoalsParams(req))
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			ctx.JSON(http.StatusNotFound, httpapi.Fail("intake form not found", err.Error()))
		case domain.ErrIntakeGoalsUpdateBlockedByActiveClient:
			ctx.JSON(http.StatusConflict, httpapi.Fail("intake goals update blocked by active client", err.Error()))
		default:
			ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to replace intake form goals", err.Error()))
		}
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toIntakeFormGoalsResponse(result), "Intake form goals replaced successfully"))
}

// UpdateIntakeConclusion updates the conclusion of an intake form.
// @Summary Update intake conclusion
// @Tags intake_forms
// @Accept json
// @Produce json
// @Param id path string true "Intake Form ID"
// @Param request body updateIntakeConclusionRequest true "Update Intake Conclusion Request"
// @Success 200 {object} httpapi.Envelope[intakeFormConclusionResponse]
// @Router /intake_forms/{id}/conclusion [patch]
func (h *IntakeFormHandler) UpdateIntakeConclusion(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid intake form ID", err.Error()))
		return
	}

	var req updateIntakeConclusionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.service.UpdateIntakeConclusion(ctx.Request.Context(), id, toUpdateIntakeConclusionParams(req))
	if err != nil {
		if strings.Contains(err.Error(), "invalid decision") {
			ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid decision", err.Error()))
			return
		}
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to update intake conclusion", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toIntakeFormConclusionResponse(result), "Intake conclusion updated successfully"))
}

// PromoteIntakeToClient promotes an intake form to a client.
// @Summary Promote intake to client
// @Tags intake_forms
// @Accept json
// @Produce json
// @Param id path string true "Intake Form ID"
// @Success 200 {object} httpapi.Envelope[promoteIntakeToClientResponse]
// @Router /intake_forms/{id}/promote [post]
func (h *IntakeFormHandler) PromoteIntakeToClient(ctx *gin.Context) {
	payload, ok := middleware.AuthPayloadFromContext(ctx.Request.Context())
	if !ok || payload == nil {
		ctx.JSON(http.StatusUnauthorized, httpapi.Fail("unauthorized access", ""))
		return
	}

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid intake form ID", err.Error()))
		return
	}

	result, err := h.service.PromoteIntakeToClient(ctx.Request.Context(), domain.PromoteIntakeToClientParams{
		IntakeFormID: id,
		EmployeeID:   payload.EmployeeID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to promote intake to client", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toPromoteIntakeToClientResponse(result), "Intake promoted to client successfully"))
}