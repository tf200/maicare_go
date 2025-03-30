// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupEmployeeRoutes(baseRouter *gin.RouterGroup) {
	employeeGroup := baseRouter.Group("/employees")
	employeeGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		employeeGroup.POST("", RBACMiddleware(server.store, "EMPLOYEE.CREATE"), server.CreateEmployeeProfileApi)
		employeeGroup.GET("", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.ListEmployeeProfileApi)
		employeeGroup.GET("/counts", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.GetEmployeeCountsApi)
		employeeGroup.GET("/:id", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.GetEmployeeProfileByIDApi)
		employeeGroup.PUT("/:id", RBACMiddleware(server.store, "EMPLOYEE.UPDATE"), server.UpdateEmployeeProfileApi)
		employeeGroup.GET("/profile", server.GetEmployeeProfileApi)
		employeeGroup.PUT("/:id/profile_picture", RBACMiddleware(server.store, "EMPLOYEE.UPDATE"), server.SetEmployeeProfilePictureApi)

		employeeGroup.POST("/:id/education", RBACMiddleware(server.store, "EMPLOYEE.CREATE"), server.AddEducationToEmployeeProfileApi)
		employeeGroup.GET("/:id/education", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.ListEmployeeEducationApi)
		employeeGroup.PUT("/:id/education/:education_id", RBACMiddleware(server.store, "EMPLOYEE.UPDATE"), server.UpdateEmployeeEducationApi)
		employeeGroup.DELETE("/:id/education/:education_id", RBACMiddleware(server.store, "EMPLOYEE.DELETE"), server.DeleteEmployeeEducationApi)

		employeeGroup.POST("/:id/experience", RBACMiddleware(server.store, "EMPLOYEE.CREATE"), server.AddEmployeeExperienceApi)
		employeeGroup.GET("/:id/experience", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.ListEmployeeExperienceApi)
		employeeGroup.PUT("/:id/experience/:experience_id", RBACMiddleware(server.store, "EMPLOYEE.UPDATE"), server.UpdateEmployeeExperienceApi)
		employeeGroup.DELETE("/:id/experience/:experience_id", RBACMiddleware(server.store, "EMPLOYEE.DELETE"), server.DeleteEmployeeExperienceApi)

		employeeGroup.POST("/:id/certification", RBACMiddleware(server.store, "EMPLOYEE.CREATE"), server.AddEmployeeCertificationApi)
		employeeGroup.GET("/:id/certification", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.ListEmployeeCertificationApi)
		employeeGroup.PUT("/:id/certification/:certification_id", RBACMiddleware(server.store, "EMPLOYEE.UPDATE"), server.UpdateEmployeeCertificationApi)
		employeeGroup.DELETE("/:id/certification/:certification_id", RBACMiddleware(server.store, "EMPLOYEE.DELETE"), server.DeleteEmployeeCertificationApi)

		employeeGroup.GET("/emails", RBACMiddleware(server.store, "EMPLOYEE.VIEW"), server.SearchEmployeesByNameOrEmailApi)

	}
	// Add other auth routes
	// auth.POST("/refresh", server.RefreshToken)
	// auth.POST("/logout", server.Logout)
}
