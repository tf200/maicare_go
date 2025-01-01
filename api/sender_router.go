package api

import "github.com/gin-gonic/gin"

func (server *Server) setupSenderRoutes(baseRouter *gin.RouterGroup) {
	senders := baseRouter.Group("/senders")
	senders.Use(AuthMiddleware(server.tokenMaker))
	{
		senders.POST("", server.CreateSenderApi) // POST /senders
		senders.GET("", server.ListSendersAPI)   // GET /senders
		// Future endpoints:
		// senders.GET("/:id", server.GetSenderAPI)    // GET /senders/:id
		// senders.PUT("/:id", server.UpdateSenderAPI) // PUT /senders/:id
		// senders.DELETE("/:id", server.DeleteSenderAPI) // DELETE /senders/:id
	}
}
