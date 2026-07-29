package handler

import (
	"maicare_go/internal/domain"

	"github.com/gin-gonic/gin"
)

func RegisterIntakeFormRoutes(
	rg *gin.RouterGroup,
	handler *IntakeFormHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	g := rg.Group("/intake_forms")
	g.POST("", auth, requirePermission(domain.PermIntakeFormCreate.String()), handler.CreateIntakeForm)
	g.GET("", auth, requirePermission(domain.PermIntakeFormView.String()), handler.ListIntakeForms)
	g.GET("/totals", auth, requirePermission(domain.PermIntakeFormView.String()), handler.GetIntakeFormTotals)
	g.GET("/:id", auth, requirePermission(domain.PermIntakeFormView.String()), handler.GetIntakeForm)
	g.PATCH("/:id", auth, requirePermission(domain.PermIntakeFormUpdate.String()), handler.UpdateIntakeForm)
	g.POST("/:id/generate_goals", auth, requirePermission(domain.PermIntakeFormCreate.String()), handler.GenerateIntakeGoals)
	g.PUT("/:id/goals", auth, requirePermission(domain.PermIntakeFormCreate.String()), handler.ReplaceIntakeFormGoals)
	g.PATCH("/:id/conclusion", auth, requirePermission(domain.PermIntakeFormCreate.String()), handler.UpdateIntakeConclusion)
	g.POST("/:id/promote", auth, requirePermission(domain.PermIntakeFormCreate.String()), handler.PromoteIntakeToClient)
	g.DELETE("/:id", auth, requirePermission(domain.PermIntakeFormDelete.String()), handler.DeleteIntakeForm)
}
