package handler

import (
	"fmt"
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/pagination"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterOrganizationRoutes(rg *gin.RouterGroup, handler *OrganizationHandler) {
	rg.GET("/organizations", handler.ListOrganizations)
	rg.GET("/organizations/:id/locations", handler.ListOrganizationLocations)
}

type OrganizationHandler struct {
	service domain.OrganizationService
}

func NewOrganizationHandler(service domain.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{service: service}
}

func (h *OrganizationHandler) ListOrganizations(ctx *gin.Context) {
	var req listOrganizationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	page, err := h.service.ListOrganizations(ctx.Request.Context(), toListOrganizationsParams(req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list organizations")))
		return
	}

	results := make([]listOrganizationsResponse, len(page.Items))
	for i, item := range page.Items {
		results[i] = toListOrganizationsResponse(item)
	}

	response := pagination.NewResponse(ctx, req.Request, results, page.TotalCount)
	ctx.JSON(http.StatusOK, successResponse(response, "Organizations retrieved successfully"))
}

func (h *OrganizationHandler) ListOrganizationLocations(ctx *gin.Context) {
	organizationID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("invalid organization ID")))
		return
	}

	var req listOrganizationLocationsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	page, err := h.service.ListOrganizationLocations(ctx.Request.Context(), toListOrganizationLocationsParams(organizationID, req))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to list locations")))
		return
	}

	results := make([]listOrganizationLocationsResponse, len(page.Items))
	for i, item := range page.Items {
		results[i] = toListOrganizationLocationsResponse(item)
	}

	response := pagination.NewResponse(ctx, req.Request, results, page.TotalCount)
	ctx.JSON(http.StatusOK, successResponse(response, "Locations retrieved successfully"))
}
