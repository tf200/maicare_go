// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupEmployeeRoutes(baseRouter *gin.RouterGroup) {
	employeeGroup := baseRouter.Group("/employees")
	employeeGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		employeeGroup.POST("", RBACMiddleware(server.store, "EMPLOYEE.CREATE"), server.CreateEmployeeProfileApi)
		employeeGroup.GET("", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.ListEmployeeProfileApi)
		employeeGroup.GET("/:id", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.GetEmployeeProfileByIDApi)
		employeeGroup.PUT("/:id", RBACMiddleware(server.store, "EMPLOYEE.UPDATE"), server.UpdateEmployeeProfileApi)
		employeeGroup.GET("/profile", server.GetEmployeeProfileApi)

		employeeGroup.POST("/:id/education", RBACMiddleware(server.store, "EMPLOYEE.CREATE"), server.AddEducationToEmployeeProfileApi)
		employeeGroup.GET("/:id/education", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.ListEmployeeEducationApi)
		employeeGroup.PUT("/:id/education/:education_id", RBACMiddleware(server.store, "EMPLOYEE.UPDATE"), server.UpdateEmployeeEducationApi)

		employeeGroup.POST("/:id/experience", RBACMiddleware(server.store, "EMPLOYEE.CREATE"), server.AddEmployeeExperienceApi)
		employeeGroup.GET("/:id/experience", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.ListEmployeeExperienceApi)
		employeeGroup.PUT("/:id/experience/:experience_id", RBACMiddleware(server.store, "EMPLOYEE.UPDATE"), server.UpdateEmployeeExperienceApi)

	}
	// Add other auth routes
	// auth.POST("/refresh", server.RefreshToken)
	// auth.POST("/logout", server.Logout)
}
