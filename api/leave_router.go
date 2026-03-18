package api

import "github.com/gin-gonic/gin"

func (server *Server) setupLeaveRoutes(baseRouter *gin.RouterGroup) {
	leaveRequestsRouter := baseRouter.Group("/leave-requests")
	leaveRequestsRouter.Use(server.AuthMiddleware())
	{
		leaveRequestsRouter.POST("", server.RBACMiddleware("LEAVE.REQUEST.CREATE"), server.CreateLeaveRequestApi)
		leaveRequestsRouter.POST("/:id/decision", server.RBACMiddleware("LEAVE.REQUEST.DECIDE"), server.DecideLeaveRequestByAdminApi)
		leaveRequestsRouter.PUT("/:id", server.RBACMiddleware("LEAVE.REQUEST.UPDATE"), server.UpdateLeaveRequestApi)
		leaveRequestsRouter.PUT("/:id/admin", server.RBACMiddleware("LEAVE.REQUEST.UPDATE_ALL"), server.UpdateLeaveRequestByAdminApi)
		leaveRequestsRouter.GET("", server.RBACMiddleware("LEAVE.REQUEST.VIEW_ALL"), server.ListLeaveRequestsApi)
		leaveRequestsRouter.GET("/my", server.RBACMiddleware("LEAVE.REQUEST.VIEW"), server.ListMyLeaveRequestsApi)
	}

	leaveBalancesRouter := baseRouter.Group("/leave-balances")
	leaveBalancesRouter.Use(server.AuthMiddleware())
	{
		leaveBalancesRouter.GET("", server.RBACMiddleware("LEAVE.BALANCE.VIEW_ALL"), server.ListLeaveBalancesApi)
		leaveBalancesRouter.GET("/my", server.RBACMiddleware("LEAVE.BALANCE.VIEW"), server.ListMyLeaveBalancesApi)
		leaveBalancesRouter.POST("/adjust", server.RBACMiddleware("LEAVE.BALANCE.ADJUST"), server.AdjustLeaveBalanceApi)
	}
}
