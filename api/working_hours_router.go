package api

import "github.com/gin-gonic/gin"

func (server *Server) setupWorkingHours(baseRouter *gin.RouterGroup) {
	workingHours := baseRouter.Group("/employees")
	workingHours.Use(AuthMiddleware(server.tokenMaker))
	{
		workingHours.GET("/:id/working_hours", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.ListWorkingHours)
	}

}
