package api

import "github.com/gin-gonic/gin"

func (server *Server) setupLeaveRoutes(baseRouter *gin.RouterGroup) {
	leaveRouter := baseRouter.Group("/leave-requests")
	leaveRouter.Use(server.AuthMiddleware())
	{
		leaveRouter.POST("", server.RBACMiddleware("LEAVE.REQUEST.CREATE"), server.CreateLeaveRequestApi)
		leaveRouter.PUT("/:id", server.RBACMiddleware("LEAVE.REQUEST.UPDATE"), server.UpdateLeaveRequestApi)
		leaveRouter.PUT("/:id/admin", server.RBACMiddleware("LEAVE.REQUEST.UPDATE_ALL"), server.UpdateLeaveRequestByAdminApi)
		leaveRouter.GET("", server.RBACMiddleware("LEAVE.REQUEST.VIEW_ALL"), server.ListLeaveRequestsApi)
		leaveRouter.GET("/my", server.RBACMiddleware("LEAVE.REQUEST.VIEW"), server.ListMyLeaveRequestsApi)
	}
}
