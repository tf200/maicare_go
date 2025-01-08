// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupEmployeeRoutes(baseRouter *gin.RouterGroup) {
	employeeGroup := baseRouter.Group("/employees")
	employeeGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		employeeGroup.POST("", RBACMiddleware(server.store, "EMPLOYEE_CREATE"), server.CreateEmployeeProfileApi)
		employeeGroup.GET("", RBACMiddleware(server.store, "EMPLOYEE_VIEW"), server.ListEmployeeProfileApi)
		employeeGroup.GET("/:id", RBACMiddleware(server.store, "EMPLOYEE_VIEW"), server.GetEmployeeProfileByIDApi)
		employeeGroup.PUT("/:id", RBACMiddleware(server.store, "EMPLOYEE_UPDATE"), server.UpdateEmployeeProfileApi)
		employeeGroup.GET("/profile", server.GetEmployeeProfileApi)
	}
	// Add other auth routes
	// auth.POST("/refresh", server.RefreshToken)
	// auth.POST("/logout", server.Logout)
}
