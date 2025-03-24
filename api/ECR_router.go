// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupECRRoutes(baseRouter *gin.RouterGroup) {

	ECRGroup := baseRouter.Group("/ecr")
	ECRGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		ECRGroup.GET("/discharge_overview", server.DischargeOverviewApi)
		ECRGroup.GET("/total_discharge_count", server.TotalDischargeCountApi)
		ECRGroup.GET("/urgent_cases_count", server.UrgentCasesCountApi)
		ECRGroup.GET("/status_change_count", server.StatusChangeCountApi)
		ECRGroup.GET("/contract_end_count", server.ContractEndCountApi)
	}

}
