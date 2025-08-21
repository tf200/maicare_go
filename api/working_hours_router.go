package api

import "github.com/gin-gonic/gin"

func (server *Server) setupWorkingHours(baseRouter *gin.RouterGroup) {
	workingHours := baseRouter.Group("/employees")
	workingHours.Use(server.AuthMiddleware())
	{
		workingHours.GET("/:id/working_hours", server.RBACMiddleware("EMPLOYEE.VIEW"), server.ListWorkingHours)
	}

}
