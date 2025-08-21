// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupEmployeeRoutes(baseRouter *gin.RouterGroup) {
	employeeGroup := baseRouter.Group("/employees")
	employeeGroup.Use(server.AuthMiddleware())
	{
		employeeGroup.POST("", server.RBACMiddleware("EMPLOYEE.CREATE"), server.CreateEmployeeProfileApi)
		employeeGroup.GET("", server.RBACMiddleware("EMPLOYEE.VIEW"), server.ListEmployeeProfileApi)
		employeeGroup.GET("/counts", server.RBACMiddleware("EMPLOYEE.VIEW"), server.GetEmployeeCountsApi)
		employeeGroup.GET("/:id", server.RBACMiddleware("EMPLOYEE.VIEW"), server.GetEmployeeProfileByIDApi)
		employeeGroup.PUT("/:id", server.RBACMiddleware("EMPLOYEE.UPDATE"), server.UpdateEmployeeProfileApi)
		employeeGroup.GET("/profile", server.GetEmployeeProfileApi)
		employeeGroup.PUT("/:id/profile_picture", server.RBACMiddleware("EMPLOYEE.UPDATE"), server.SetEmployeeProfilePictureApi)
		employeeGroup.PUT("/:id/contract_details", server.RBACMiddleware("EMPLOYEE.UPDATE"), server.AddEmployeeContractDetailsApi)
		employeeGroup.GET("/:id/contract_details", server.RBACMiddleware("EMPLOYEE.VIEW"), server.GetEmployeeContractDetailsApi)

		employeeGroup.POST("/:id/education", server.RBACMiddleware("EMPLOYEE.CREATE"), server.AddEducationToEmployeeProfileApi)
		employeeGroup.GET("/:id/education", server.RBACMiddleware("EMPLOYEE.VIEW"), server.ListEmployeeEducationApi)
		employeeGroup.PUT("/:id/education/:education_id", server.RBACMiddleware("EMPLOYEE.UPDATE"), server.UpdateEmployeeEducationApi)
		employeeGroup.DELETE("/:id/education/:education_id", server.RBACMiddleware("EMPLOYEE.DELETE"), server.DeleteEmployeeEducationApi)

		employeeGroup.POST("/:id/experience", server.RBACMiddleware("EMPLOYEE.CREATE"), server.AddEmployeeExperienceApi)
		employeeGroup.GET("/:id/experience", server.RBACMiddleware("EMPLOYEE.VIEW"), server.ListEmployeeExperienceApi)
		employeeGroup.PUT("/:id/experience/:experience_id", server.RBACMiddleware("EMPLOYEE.UPDATE"), server.UpdateEmployeeExperienceApi)
		employeeGroup.DELETE("/:id/experience/:experience_id", server.RBACMiddleware("EMPLOYEE.DELETE"), server.DeleteEmployeeExperienceApi)

		employeeGroup.POST("/:id/certification", server.RBACMiddleware("EMPLOYEE.CREATE"), server.AddEmployeeCertificationApi)
		employeeGroup.GET("/:id/certification", server.RBACMiddleware("EMPLOYEE.VIEW"), server.ListEmployeeCertificationApi)
		employeeGroup.PUT("/:id/certification/:certification_id", server.RBACMiddleware("EMPLOYEE.UPDATE"), server.UpdateEmployeeCertificationApi)
		employeeGroup.DELETE("/:id/certification/:certification_id", server.RBACMiddleware("EMPLOYEE.DELETE"), server.DeleteEmployeeCertificationApi)

		employeeGroup.POST("/:id/appointments", server.RBACMiddleware("EMPLOYEE.CREATE"), server.ListAppointmentsForEmployee)

		employeeGroup.GET("/emails", server.RBACMiddleware("EMPLOYEE.VIEW"), server.SearchEmployeesByNameOrEmailApi)

	}
	// Add other auth routes
	// auth.POST("/refresh", server.RefreshToken)
	// auth.POST("/logout", server.Logout)
}
