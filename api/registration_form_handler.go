package api

import (
	"errors"
	"io"
	"net/http"
	"strings"

	_ "maicare_go/pagination"
	clientp "maicare_go/service/client"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// @Summary Create Registration Form
// @Description Create a new registration form
// @Tags Registration Form
// @Accept json
// @Produce json
// @Param request body clientp.CreateRegistrationFormRequest true "Create Registration Form Request"
// @Success 200 {object} clientp.CreateRegistrationFormResponse
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_forms [post]
func (server *Server) CreateRegistrationFormApi(ctx *gin.Context) {
	var req clientp.CreateRegistrationFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil && !errors.Is(err, io.EOF) {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.businessService.ClientService.CreateRegistrationForm(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(response, "Registration form created successfully")
	ctx.JSON(http.StatusCreated, res)
}

// @Summary List Registration Forms
// @Description List all registration forms
// @Tags Registration Form
// @Produce json
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Form status" Enums(pending, processed, rejected)
// @Param risk_aggressive_behavior query bool false "Risk aggressive behavior"
// @Param risk_suicidal_selfharm query bool false "Risk suicidal self-harm"
// @Param risk_substance_abuse query bool false "Risk substance abuse"
// @Param risk_psychiatric_issues query bool false "Risk psychiatric issues"
// @Param risk_criminal_history query bool false "Risk criminal history"
// @Param risk_flight_behavior query bool false "Risk flight behavior"
// @Param risk_weapon_possession query bool false "Risk weapon possession"
// @Param risk_sexual_behavior query bool false "Risk sexual behavior"
// @Param risk_day_night_rhythm query bool false "Risk day-night rhythm"
// @Success 200 {object} Response[pagination.Response[clientp.ListRegistrationFormsResponse]]
// @Failure 400 {object}  Response[any]
// @Failure 500 {object}  Response[any]
// @Router /registration_forms [get]
func (server *Server) ListRegistrationFormsApi(ctx *gin.Context) {
	var req clientp.ListRegistrationFormsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.businessService.ClientService.ListRegistrationForms(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(response, "Registration forms retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Registration Form
// @Description Get a registration form by ID
// @Tags Registration Form
// @Produce json
// @Param id path uuid true "Registration Form ID"
// @Success 200 {object} Response[clientp.GetRegistrationFormResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_forms/{id} [get]
func (server *Server) GetRegistrationFormApi(ctx *gin.Context) {
	rfId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.businessService.ClientService.GetRegistrationFormB(ctx, rfId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(response, "Registration form retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update Registration Form
// @Description Update a registration form by ID
// @Tags Registration Form
// @Accept json
// @Produce json
// @Param id path uuid true "Registration Form ID"
// @Param request body clientp.UpdateRegistrationFormRequest true "Update Registration Form Request"
// @Success 200 {object} Response[clientp.UpdateRegistrationFormResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_forms/{id} [put]
func (server *Server) UpdateRegistrationFormApi(ctx *gin.Context) {
	rfId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req clientp.UpdateRegistrationFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.businessService.ClientService.UpdateRegistrationForm(ctx, &req, rfId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(response, "Registration form updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Registration Form
// @Description Delete a registration form by ID
// @Tags Registration Form
// @Produce json
// @Param id path uuid true "Registration Form ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_forms/{id} [delete]
func (server *Server) DeleteRegistrationFormApi(ctx *gin.Context) {
	rfId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.businessService.ClientService.DeleteRegistrationForm(ctx, rfId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Registration form deleted successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Update Registration Form Status
// @Description Process a registration form by sending intake proposals
// @Tags Registration Form
// @Produce json
// @Param id path uuid true "Registration Form ID"
// @Param request body clientp.ProcessRegistrationFormRequest true "Process Registration Form Request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_forms/{id}/process [post]
func (server *Server) ProcessRegistrationFormApi(ctx *gin.Context) {
	var req clientp.ProcessRegistrationFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	rfId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.businessService.ClientService.ProcessRegistrationForm(ctx, &req, rfId, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Registration form processed and intake options sent")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Public Intake Options
// @Description Get intake options for a registration form via token
// @Tags Public Intake
// @Produce json
// @Param token path string true "Intake Token"
// @Success 200 {object} Response[clientp.PublicIntakeOptionsResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /public/intake-options/{token} [get]
func (server *Server) GetPublicIntakeOptionsApi(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(http.ErrNoLocation))
		return
	}

	response, err := server.businessService.ClientService.GetPublicIntakeOptions(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
		return
	}

	res := SuccessResponse(response, "Intake options retrieved successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Select Intake Date
// @Description Select an intake date for a registration form via token
// @Tags Public Intake
// @Produce json
// @Param token path string true "Intake Token"
// @Param request body clientp.SelectIntakeDateRequest true "Select Intake Date Request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /public/intake-options/{token}/confirm [post]
func (server *Server) SelectIntakeDateApi(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(http.ErrNoLocation))
		return
	}

	var req clientp.SelectIntakeDateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := server.businessService.ClientService.SelectIntakeDate(ctx, token, &req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Intake date selected successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Description Update the status of a registration form by ID
// @Tags Registration Form
// @Produce json
// @Param id path uuid true "Registration Form ID"
// @Param request body clientp.UpdateRegistrationFormStatusRequest true "Update Registration Form Status Request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_forms/{id}/status [post]
func (server *Server) UpdateRegistrationFormStatusApi(ctx *gin.Context) {
	var req clientp.UpdateRegistrationFormStatusRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload, err := GetAuthPayload(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	rfId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	err = server.businessService.ClientService.UpdateRegistrationFormStatus(ctx, &req, rfId, payload.EmployeeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse[any](nil, "Registration form status updated successfully")
	ctx.JSON(http.StatusOK, res)
}

// @Summary Promote Intake to Client
// @Description Promote an intake form and all its assessments to a full client record with a care plan
// @Tags Registration Form
// @Produce json
// @Param id path uuid true "Intake Form ID"
// @Success 200 {object} Response[clientp.PromoteIntakeToClientResponse]
// @Router /intake_forms/{id}/promote [post]
func (server *Server) PromoteIntakeToClientApi(ctx *gin.Context) {
	intakeFormID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	req := clientp.PromoteIntakeToClientRequest{
		IntakeFormID: intakeFormID,
	}

	response, err := server.businessService.ClientService.PromoteIntakeToClient(ctx, &req)
	if err != nil {
		if strings.Contains(err.Error(), "only intakes with conclusion 'suitable' can be promoted") {
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := SuccessResponse(response, "Intake successfully promoted to client")
	ctx.JSON(http.StatusOK, res)
}
