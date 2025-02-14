// auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupMaturityMatrixRoutes(baseRouter *gin.RouterGroup) {

	mmGroup := baseRouter.Group("/clients")
	mmGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		mmGroup.POST("/:id/maturity_matrix_assessment", RBACMiddleware(server.store, "CLIENT.CREATE"), server.CreateClientMaturityMatrixAssessmentApi)
	}

}
