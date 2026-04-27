package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterIntakeFormRoutes(
	rg *gin.RouterGroup,
	handler *IntakeFormHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	g := rg.Group("/intake_forms")
	g.POST("", auth, requirePermission("REGISTRATION_FORM.UPDATE"), handler.CreateIntakeForm)
	g.GET("", auth, requirePermission("REGISTRATION_FORM.VIEW"), handler.ListIntakeForms)
	g.GET("/totals", auth, requirePermission("REGISTRATION_FORM.VIEW"), handler.GetIntakeFormTotals)
	g.GET("/:id", auth, requirePermission("REGISTRATION_FORM.VIEW"), handler.GetIntakeForm)
	g.PATCH("/:id", auth, requirePermission("REGISTRATION_FORM.UPDATE"), handler.UpdateIntakeForm)
	g.POST("/:id/generate_goals", auth, requirePermission("REGISTRATION_FORM.UPDATE"), handler.GenerateIntakeGoals)
	g.PUT("/:id/goals", auth, requirePermission("REGISTRATION_FORM.UPDATE"), handler.ReplaceIntakeFormGoals)
	g.PATCH("/:id/conclusion", auth, requirePermission("REGISTRATION_FORM.UPDATE"), handler.UpdateIntakeConclusion)
	g.POST("/:id/promote", auth, requirePermission("REGISTRATION_FORM.UPDATE"), handler.PromoteIntakeToClient)
}