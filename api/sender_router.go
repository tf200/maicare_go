package api

import "github.com/gin-gonic/gin"

func (server *Server) setupSenderRoutes(baseRouter *gin.RouterGroup) {
	senders := baseRouter.Group("/senders")
	senders.Use(AuthMiddleware(server.tokenMaker))
	{
		senders.POST("", RBACMiddleware(server.store, "SENDER.CREATE"), server.CreateSenderApi) // POST /senders
		senders.GET("", RBACMiddleware(server.store, "SENDER.VIEW"), server.ListSendersAPI)
		senders.GET("/:id", RBACMiddleware(server.store, "SENDER.VIEW"), server.GetSenderByIdAPI)
		senders.PUT("/:id", RBACMiddleware(server.store, "SENDER.UPDATE"), server.UpdateSenderApi) // PUT /senders/:id
		// GET /senders
		// Future endpoints:
		// senders.GET("/:id", server.GetSenderAPI)    // GET /senders/:id
		// senders.PUT("/:id", server.UpdateSenderAPI) // PUT /senders/:id
		senders.DELETE("/:id", RBACMiddleware(server.store, "SENDER.VIEW"), server.DeleteSenderApi)                                 // DELETE /senders/:id
		senders.POST("/:id/invoice_template", RBACMiddleware(server.store, "SENDER.CREATE"), server.CreateSenderInvoiceTemplateApi) // POST /senders/:id/invoice_template
	}
}
