package api

import (
	"fmt"
	"maicare_go/service/settings"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary List departments
// @Tags settings
// @Produce json
// @Success 200 {object} Response[any]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /settings/departments [get]
func (server *Server) ListDepartmentsApi(ctx *gin.Context) {
	res, err := server.businessService.SettingsService.ListDepartments(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Departments retrieved"))
}

// @Summary Create department
// @Tags settings
// @Accept json
// @Produce json
// @Param request body settings.CreateDepartmentRequest true "Department payload"
// @Success 201 {object} Response[settings.CreateDepartmentResponse]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /settings/departments [post]
func (server *Server) CreateDepartmentApi(ctx *gin.Context) {
	var req settings.CreateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}
	res, err := server.businessService.SettingsService.CreateDepartment(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, SuccessResponse(res, "Department created"))
}

// @Summary Update department
// @Tags settings
// @Accept json
// @Produce json
// @Param id path uuid true "Department ID"
// @Param request body settings.UpdateDepartmentRequest true "Department payload"
// @Success 200 {object} Response[settings.UpdateDepartmentResponse]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /settings/departments/{id} [put]
func (server *Server) UpdateDepartmentApi(ctx *gin.Context) {
	departmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid department id")))
		return
	}

	var req settings.UpdateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	res, err := server.businessService.SettingsService.UpdateDepartment(ctx, departmentID, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Department updated"))
}

// @Summary Get organization profile
// @Tags settings
// @Produce json
// @Success 200 {object} Response[settings.GetOrganizationProfileResponse]
// @Failure 401,403,500 {object} Response[any]
// @Router /settings/organization-profile [get]
func (server *Server) GetOrganizationProfileApi(ctx *gin.Context) {
	res, err := server.businessService.SettingsService.GetOrganizationProfile(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Organization profile retrieved"))
}

// @Summary Update organization profile
// @Tags settings
// @Accept json
// @Produce json
// @Param request body settings.UpdateOrganizationProfileRequest true "Organization profile payload"
// @Success 200 {object} Response[settings.GetOrganizationProfileResponse]
// @Failure 400,401,403,500 {object} Response[any]
// @Router /settings/organization-profile [put]
func (server *Server) UpdateOrganizationProfileApi(ctx *gin.Context) {
	var req settings.UpdateOrganizationProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid request body")))
		return
	}

	res, err := server.businessService.SettingsService.UpdateOrganizationProfile(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, SuccessResponse(res, "Organization profile updated"))
}
