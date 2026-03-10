package api

import "github.com/gin-gonic/gin"

func (server *Server) setupHandbookRoutes(baseRouter *gin.RouterGroup) {
	handbook := baseRouter.Group("/handbook")
	handbook.Use(server.AuthMiddleware())

	// Employee self-service
	handbook.GET("/me", server.RBACMiddleware("HANDBOOK.SELF.VIEW"), server.GetMyActiveHandbookApi)
	handbook.POST("/me/start", server.RBACMiddleware("HANDBOOK.SELF.UPDATE"), server.StartMyHandbookApi)
	handbook.POST("/me/steps/:step_id/complete", server.RBACMiddleware("HANDBOOK.SELF.UPDATE"), server.CompleteMyHandbookStepApi)

	// Admin/manager
	handbook.POST("/templates", server.RBACMiddleware("HANDBOOK.TEMPLATE.CREATE"), server.CreateHandbookTemplateApi)
	handbook.POST("/templates/clone", server.RBACMiddleware("HANDBOOK.TEMPLATE.CREATE"), server.CloneHandbookTemplateApi)
	handbook.PATCH("/templates/:template_id", server.RBACMiddleware("HANDBOOK.TEMPLATE.UPDATE"), server.UpdateHandbookTemplateApi)
	handbook.POST("/templates/publish", server.RBACMiddleware("HANDBOOK.TEMPLATE.PUBLISH"), server.PublishHandbookTemplateApi)
	handbook.GET("/departments/:department_id/templates", server.RBACMiddleware("HANDBOOK.TEMPLATE.VIEW"), server.ListHandbookTemplatesByDepartmentApi)

	handbook.POST("/steps", server.RBACMiddleware("HANDBOOK.STEP.CREATE"), server.CreateHandbookStepApi)
	handbook.PATCH("/steps/:step_id", server.RBACMiddleware("HANDBOOK.STEP.UPDATE"), server.UpdateHandbookStepApi)
	handbook.DELETE("/steps/:step_id", server.RBACMiddleware("HANDBOOK.STEP.DELETE"), server.DeleteHandbookStepApi)
	handbook.GET("/templates/:template_id/steps", server.RBACMiddleware("HANDBOOK.STEP.VIEW"), server.ListHandbookStepsByTemplateApi)
	handbook.POST("/templates/:template_id/steps/reorder", server.RBACMiddleware("HANDBOOK.STEP.UPDATE"), server.ReorderHandbookStepsApi)

	handbook.POST("/assignments", server.RBACMiddleware("HANDBOOK.ASSIGN"), server.AssignHandbookTemplateToEmployeeApi)
}
