package api

import "github.com/gin-gonic/gin"

func (server *Server) setupContractRoutes(baseRouter *gin.RouterGroup) {

	clientGroup := baseRouter.Group("/clients")
	clientGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		clientGroup.POST("/:id/contracts", server.CreateContractApi)
		clientGroup.GET("/:id/contracts", server.GetClientContractApi)
	}

	// Routes without /client prefix
	baseRouter.POST("/contract_types", server.CreateContractTypeApi)
	baseRouter.GET("/contract_types", server.ListContractTypesApi)

	baseRouter.GET("/contracts", server.ListContractsApi)

}
