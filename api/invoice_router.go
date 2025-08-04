package api

import "github.com/gin-gonic/gin"

func (server *Server) setupInvoiceRoutes(baseRouter *gin.RouterGroup) {
	invoiceGroup := baseRouter.Group("/invoices")
	invoiceGroup.Use(AuthMiddleware(server.tokenMaker))
	{
		invoiceGroup.POST("", RBACMiddleware(server.store, "INVOICE.CREATE"), server.CreateInvoiceApi)
		invoiceGroup.POST("/generate", RBACMiddleware(server.store, "INVOICE.CREATE"), server.GenerateInvoiceApi)
		invoiceGroup.GET("", RBACMiddleware(server.store, "INVOICE.VIEW"), server.ListInvoicesApi)
		invoiceGroup.GET("/:id", RBACMiddleware(server.store, "INVOICE.VIEW"), server.GetInvoiceByIDApi)
		invoiceGroup.PUT("/:id", RBACMiddleware(server.store, "INVOICE.UPDATE"), server.UpdateInvoiceApi)
		invoiceGroup.DELETE("/:id", RBACMiddleware(server.store, "INVOICE.DELETE"), server.DeleteInvoiceApi)
		invoiceGroup.POST("/:id/credit", RBACMiddleware(server.store, "INVOICE.CREDIT"), server.CreditInvoiceApi)
		invoiceGroup.GET("/:id/generate_pdf", RBACMiddleware(server.store, "INVOICE.VIEW"), server.GenerateInvoicePdfApi)

		invoiceGroup.POST("/:id/payments", RBACMiddleware(server.store, "INVOICE.PAYMENT.CREATE"), server.CreatePaymentApi)
		invoiceGroup.GET("/:id/payments", RBACMiddleware(server.store, "INVOICE.PAYMENT.VIEW"), server.ListPaymentsApi)
		invoiceGroup.GET("/:id/payments/:payment_id", RBACMiddleware(server.store, "INVOICE.PAYMENT.VIEW"), server.GetPaymentByIDApi)
		invoiceGroup.PUT("/:id/payments/:payment_id", RBACMiddleware(server.store, "INVOICE.PAYMENT.UPDATE"), server.UpdatePaymentApi)
		invoiceGroup.DELETE("/:id/payments/:payment_id", RBACMiddleware(server.store, "INVOICE.PAYMENT.DELETE"), server.DeletePaymentApi)

		invoiceGroup.GET("/:id/audit", RBACMiddleware(server.store, "INVOICE.VIEW"), server.GetInvoiceAuditLogApi)

		invoiceGroup.GET("/template_items", RBACMiddleware(server.store, "INVOICE.VIEW"), server.GetInvoiceTemplateItemsApi)
	}
}
