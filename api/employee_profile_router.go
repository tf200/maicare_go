// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupEmployeeRoutes(baseRouter *gin.RouterGroup) {
	employeeGroup := baseRouter.Group("/employees")
	employeeGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		employeeGroup.POST("/employees_create/", server.CreateEmployeeProfileApi)
		employeeGroup.GET("/employees_list/", server.ListEmployeeProfileApi)
		employeeGroup.GET("/profile", server.GetEmployeeProfileApi)
	}
	// Add other auth routes
	// auth.POST("/refresh", server.RefreshToken)
	// auth.POST("/logout", server.Logout)
}
