package handler

import (
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SettingsHandler struct {
	departmentService      domain.DepartmentService
	organizationProfileService domain.OrganizationProfileService
}

func NewSettingsHandler(
	departmentService domain.DepartmentService,
	organizationProfileService domain.OrganizationProfileService,
) *SettingsHandler {
	return &SettingsHandler{
		departmentService:      departmentService,
		organizationProfileService: organizationProfileService,
	}
}

// ListDepartments returns all departments.
// @Summary List departments
// @Tags settings
// @Produce json
// @Success 200 {object} httpapi.Envelope[httpapi.PageResponse[departmentResponse]]
// @Failure 401,403,500 {object} httpapi.Envelope[any]
// @Router /settings/departments [get]
func (h *SettingsHandler) ListDepartments(ctx *gin.Context) {
	departments, err := h.departmentService.ListDepartments(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list departments", err.Error()))
		return
	}

	items := make([]departmentResponse, len(departments))
	for i, d := range departments {
		items[i] = toDepartmentResponse(d)
	}

	pageResp := httpapi.PageResponse[departmentResponse]{
		Next:     nil,
		Previous: nil,
		Count:    int64(len(items)),
		PageSize: int32(len(items)),
		Results:  items,
	}

	ctx.JSON(http.StatusOK, httpapi.OK(pageResp, "Departments retrieved"))
}

// CreateDepartment creates a new department.
// @Summary Create department
// @Tags settings
// @Accept json
// @Produce json
// @Param request body createDepartmentRequest true "Department payload"
// @Success 201 {object} httpapi.Envelope[departmentResponse]
// @Failure 400,401,403,500 {object} httpapi.Envelope[any]
// @Router /settings/departments [post]
func (h *SettingsHandler) CreateDepartment(ctx *gin.Context) {
	var req createDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.departmentService.CreateDepartment(ctx.Request.Context(), toCreateDepartmentParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to create department", err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, httpapi.OK(toDepartmentResponse(*result), "Department created"))
}

// UpdateDepartment updates an existing department.
// @Summary Update department
// @Tags settings
// @Accept json
// @Produce json
// @Param id path uuid true "Department ID"
// @Param request body updateDepartmentRequest true "Department payload"
// @Success 200 {object} httpapi.Envelope[departmentResponse]
// @Failure 400,401,403,500 {object} httpapi.Envelope[any]
// @Router /settings/departments/{id} [put]
func (h *SettingsHandler) UpdateDepartment(ctx *gin.Context) {
	departmentID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid department ID", err.Error()))
		return
	}

	var req updateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.departmentService.UpdateDepartment(ctx.Request.Context(), toUpdateDepartmentParams(departmentID, req))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("failed to update department", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toDepartmentResponse(*result), "Department updated"))
}

// GetOrganizationProfile returns the organization profile.
// @Summary Get organization profile
// @Tags settings
// @Produce json
// @Success 200 {object} httpapi.Envelope[organizationProfileResponse]
// @Failure 401,403,500 {object} httpapi.Envelope[any]
// @Router /settings/organization-profile [get]
func (h *SettingsHandler) GetOrganizationProfile(ctx *gin.Context) {
	profile, err := h.organizationProfileService.GetOrganizationProfile(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get organization profile", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toOrganizationProfileResponse(profile), "Organization profile retrieved"))
}

// UpdateOrganizationProfile updates the organization profile.
// @Summary Update organization profile
// @Tags settings
// @Accept json
// @Produce json
// @Param request body updateOrganizationProfileRequest true "Organization profile payload"
// @Success 200 {object} httpapi.Envelope[organizationProfileResponse]
// @Failure 400,401,403,500 {object} httpapi.Envelope[any]
// @Router /settings/organization-profile [put]
func (h *SettingsHandler) UpdateOrganizationProfile(ctx *gin.Context) {
	var req updateOrganizationProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("invalid request body", err.Error()))
		return
	}

	result, err := h.organizationProfileService.UpdateOrganizationProfile(ctx.Request.Context(), toUpdateOrganizationProfileParams(req))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpapi.Fail("failed to update organization profile", err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toOrganizationProfileResponse(result), "Organization profile updated"))
}
