package api

import "github.com/gin-gonic/gin"

func (server *Server) setupSenderRoutes(baseRouter *gin.RouterGroup) {
	senders := baseRouter.Group("/senders")
	senders.Use(server.AuthMiddleware())
	{
		senders.POST("", server.RBACMiddleware("SENDER.CREATE"), server.CreateSenderApi) // POST /senders
		senders.GET("", server.RBACMiddleware("SENDER.VIEW"), server.ListSendersAPI)
		senders.GET("/:id", server.RBACMiddleware("SENDER.VIEW"), server.GetSenderByIdAPI)
		senders.PUT("/:id", server.RBACMiddleware("SENDER.UPDATE"), server.UpdateSenderApi) // PUT /senders/:id
		// GET /senders
		// Future endpoints:
		// senders.GET("/:id", server.GetSenderAPI)    // GET /senders/:id
		// senders.PUT("/:id", server.UpdateSenderAPI) // PUT /senders/:id
		senders.DELETE("/:id", server.RBACMiddleware("SENDER.DELETE"), server.DeleteSenderApi)                               // DELETE /senders/:id
		senders.POST("/:id/invoice_template", server.RBACMiddleware("SENDER.CREATE"), server.CreateSenderInvoiceTemplateApi) // POST /senders/:id/invoice_template
	}
}
