package api

import (
	clientp "maicare_go/service/client"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create Registration Form
// @Description Create a new registration form
// @Tags Registration Form
// @Accept json
// @Produce json
// @Param request body CreateRegistrationFormRequest true "Create Registration Form Request"
// @Success 200 {object} CreateRegistrationFormResponse
// @Failure 400 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_form [post]
func (server *Server) CreateRegistrationFormApi(ctx *gin.Context) {
	var req clientp.CreateRegistrationFormRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
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
// @Param status query string false "Form status" Enums(pending, approved, rejected)
// @Param risk_aggressive_behavior query bool false "Risk aggressive behavior"
// @Param risk_suicidal_selfharm query bool false "Risk suicidal self-harm"
// @Param risk_substance_abuse query bool false "Risk substance abuse"
// @Param risk_psychiatric_issues query bool false "Risk psychiatric issues"
// @Param risk_criminal_history query bool false "Risk criminal history"
// @Param risk_flight_behavior query bool false "Risk flight behavior"
// @Param risk_weapon_possession query bool false "Risk weapon possession"
// @Param risk_sexual_behavior query bool false "Risk sexual behavior"
// @Param risk_day_night_rhythm query bool false "Risk day-night rhythm"
// @Success 200 {object} Response[pagination.Response[ListRegistrationFormsResponse]]
// @Failure 400 {object}  Response[any]
// @Failure 500 {object}  Response[any]
// @Router /registration_form [get]
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
// @Param id path int true "Registration Form ID"
// @Success 200 {object} Response[GetRegistrationFormResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_form/{id} [get]
func (server *Server) GetRegistrationFormApi(ctx *gin.Context) {
	rfId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response, err := server.businessService.ClientService.GetRegistrationFormB(ctx, rfId)
	if err != nil {
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
// @Param id path int true "Registration Form ID"
// @Param request body UpdateRegistrationFormRequest true "Update Registration Form Request"
// @Success 200 {object} Response[UpdateRegistrationFormResponse]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_form/{id} [put]
func (server *Server) UpdateRegistrationFormApi(ctx *gin.Context) {
	rfId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
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
// @Param id path int true "Registration Form ID"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_form/{id} [delete]
func (server *Server) DeleteRegistrationFormApi(ctx *gin.Context) {
	rfId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
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
// @Description Update the status of a registration form by ID
// @Tags Registration Form
// @Produce json
// @Param id path int true "Registration Form ID"
// @Param request body UpdateRegistrationFormStatusRequest true "Update Registration Form Status Request"
// @Success 200 {object} Response[any]
// @Failure 400 {object} Response[any]
// @Failure 404 {object} Response[any]
// @Failure 500 {object} Response[any]
// @Router /registration_form/{id}/status [post]
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

	rfId, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
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
