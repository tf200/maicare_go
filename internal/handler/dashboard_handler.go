package handler

import (
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service domain.DashboardService
}

func NewDashboardHandler(service domain.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func RegisterDashboardRoutes(
	rg *gin.RouterGroup,
	handler *DashboardHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	dashboardGroup := rg.Group("/dashboard")
	{
		dashboardGroup.GET("/admin", auth, requirePermission("DASHBOARD.VIEW"), handler.GetAdminDashboard)
	}
}

// GetAdminDashboard gets admin dashboard statistics.
// @Summary Get admin dashboard statistics
// @Tags dashboard
// @Produce json
// @Success 200 {object} httpapi.Envelope[adminDashboardResponse]
// @Failure 401,403,500 {object} httpapi.Envelope[any]
// @Router /dashboard/admin [get]
func (h *DashboardHandler) GetAdminDashboard(ctx *gin.Context) {
	stats, err := h.service.GetAdminDashboardStats(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to get admin dashboard stats", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(toAdminDashboardResponse(stats), "Admin dashboard stats fetched successfully"))
}
