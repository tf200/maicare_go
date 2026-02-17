package api

import "github.com/gin-gonic/gin"

func (server *Server) setupContractRoutes(baseRouter *gin.RouterGroup) {
	// Per-client contract listing
	clientGroup := baseRouter.Group("/clients")
	clientGroup.Use(server.AuthMiddleware())
	{
		clientGroup.GET("/:id/contracts", server.RBACMiddleware("CONTRACT.VIEW"), server.ListClientContractsApi)
	}

	// Contract type routes
	baseRouter.POST("/contract_types", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT_TYPE.CREATE"), server.CreateContractTypeApi)
	baseRouter.GET("/contract_types", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT_TYPE.VIEW"), server.ListContractTypesApi)
	baseRouter.DELETE("/contract_types/:id", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT_TYPE.DELETE"), server.DeleteContractTypeApi)

	// Contract routes
	baseRouter.POST("/contracts", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT.CREATE"), server.CreateContractApi)
	baseRouter.GET("/contracts", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT.VIEW"), server.ListContractsApi)
	baseRouter.GET("/contracts/:id", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT.VIEW"), server.GetClientContractApi)
	baseRouter.PUT("/contracts/:id", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT.UPDATE"), server.UpdateContractApi)
	baseRouter.PUT("/contracts/:id/status", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT.UPDATE"), server.UpdateContractStatusApi)
	baseRouter.GET("/contracts/:id/audit", server.AuthMiddleware(), server.RBACMiddleware("CONTRACT.VIEW"), server.GetContractAuditLogApi)
}
