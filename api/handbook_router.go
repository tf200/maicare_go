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
	handbook.GET("/departments/:department_id/templates", server.RBACMiddleware("HANDBOOK.TEMPLATE.VIEW"), server.ListHandbookTemplatesByDepartmentApi)

	handbook.POST("/steps", server.RBACMiddleware("HANDBOOK.STEP.CREATE"), server.CreateHandbookStepApi)
	handbook.GET("/templates/:template_id/steps", server.RBACMiddleware("HANDBOOK.STEP.VIEW"), server.ListHandbookStepsByTemplateApi)

	handbook.POST("/assignments", server.RBACMiddleware("HANDBOOK.ASSIGN"), server.AssignHandbookTemplateToEmployeeApi)
}
