package api

import "github.com/gin-gonic/gin"

func (server *Server) setupInvoiceRoutes(baseRouter *gin.RouterGroup) {
	invoiceGroup := baseRouter.Group("/invoices")
	invoiceGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		invoiceGroup.POST("/generate", RBACMiddleware(server.store, "INVOICE.CREATE"), server.GenerateInvoiceApi)
		invoiceGroup.GET("", RBACMiddleware(server.store, "INVOICE.VIEW"), server.ListInvoicesApi)
		invoiceGroup.GET("/:id", RBACMiddleware(server.store, "INVOICE.VIEW"), server.GetInvoiceByIDApi)
		invoiceGroup.PUT("/:id", RBACMiddleware(server.store, "INVOICE.UPDATE"), server.UpdateInvoiceApi)
		invoiceGroup.DELETE("/:id", RBACMiddleware(server.store, "INVOICE.DELETE"), server.DeleteInvoiceApi)
	}
}
