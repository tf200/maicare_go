// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupECRRoutes(baseRouter *gin.RouterGroup) {

	ECRGroup := baseRouter.Group("/ecr")
	ECRGroup.Use(server.AuthMiddleware())
	{
		ECRGroup.GET("/discharge_overview", server.DischargeOverviewApi)
		ECRGroup.GET("/total_discharge_count", server.TotalDischargeCountApi)

	}

}
