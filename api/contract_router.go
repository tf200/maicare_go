package api

import "github.com/gin-gonic/gin"

func (server *Server) setupContractRoutes(baseRouter *gin.RouterGroup) {

	clientGroup := baseRouter.Group("/clients")
	clientGroup.Use(server.AuthMiddleware())
	{
		clientGroup.POST("/:id/contracts", server.RBACMiddleware("CONTRACT.CREATE"), server.CreateContractApi)
		clientGroup.GET("/:id/contracts", server.RBACMiddleware("CONTRACT.VIEW"), server.ListClientContractsApi)
		clientGroup.GET("/:id/contracts/:contract_id", server.RBACMiddleware("CONTRACT.VIEW"), server.GetClientContractApi)
	}

	// Routes without /client prefix
	baseRouter.POST("/contract_types", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT_TYPE.CREATE"), server.CreateContractTypeApi)
	baseRouter.GET("/contract_types", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT_TYPE.VIEW"), server.ListContractTypesApi)
	baseRouter.DELETE("/contract_types/:id", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT_TYPE.DELETE"), server.DeleteContractTypeApi)

	baseRouter.GET("/contracts", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT.VIEW"), server.ListContractsApi)
	baseRouter.PUT("/contracts/:id", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT.UPDATE"), server.UpdateContractApi)

	baseRouter.GET("/contracts/:id/audit", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT.VIEW"), server.GetContractAuditLogApi)

}
