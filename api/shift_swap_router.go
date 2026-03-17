package api

import "github.com/gin-gonic/gin"

func (server *Server) setupShiftSwapRoutes(baseRouter *gin.RouterGroup) {
	shiftSwap := baseRouter.Group("")
	shiftSwap.Use(server.AuthMiddleware())
	{
		shiftSwap.POST("/shift-swaps", server.RBACMiddleware("SCHEDULE_SWAP.REQUEST"), server.CreateShiftSwapRequestApi)
		shiftSwap.POST("/shift-swaps/:id/respond", server.RBACMiddleware("SCHEDULE_SWAP.RESPOND"), server.RespondShiftSwapRequestApi)
		shiftSwap.POST("/shift-swaps/:id/admin-decision", server.RBACMiddleware("SCHEDULE_SWAP.APPROVE"), server.AdminDecisionShiftSwapRequestApi)
		shiftSwap.GET("/shift-swaps", server.RBACMiddleware("SCHEDULE_SWAP.APPROVE"), server.ListShiftSwapRequestsApi)
		shiftSwap.GET("/shift-swaps/my", server.RBACMiddleware("SCHEDULE_SWAP.VIEW"), server.ListMyShiftSwapRequestsApi)
	}
}
