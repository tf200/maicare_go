package api

import "github.com/gin-gonic/gin"

func (server *Server) setupInvoiceRoutes(baseRouter *gin.RouterGroup) {
	invoiceGroup := baseRouter.Group("/invoices")
	invoiceGroup.Use(server.AuthMiddleware())
	{
		invoiceGroup.POST("", server.RBACMiddleware("INVOICE.CREATE"), server.CreateInvoiceApi)
		invoiceGroup.POST("/generate", server.RBACMiddleware("INVOICE.CREATE"), server.GenerateInvoiceApi)
		invoiceGroup.GET("", server.RBACMiddleware("INVOICE.VIEW"), server.ListInvoicesApi)
		invoiceGroup.GET("/:id", server.RBACMiddleware("INVOICE.VIEW"), server.GetInvoiceByIDApi)
		invoiceGroup.PUT("/:id", server.RBACMiddleware("INVOICE.UPDATE"), server.UpdateInvoiceApi)
		invoiceGroup.DELETE("/:id", server.RBACMiddleware("INVOICE.DELETE"), server.DeleteInvoiceApi)
		invoiceGroup.POST("/:id/credit", server.RBACMiddleware("INVOICE.UPDATE"), server.CreditInvoiceApi)
		invoiceGroup.GET("/:id/generate_pdf", server.RBACMiddleware("INVOICE.VIEW"), server.GenerateInvoicePdfApi)

		invoiceGroup.POST("/:id/payments", server.RBACMiddleware("INVOICE.PAYMENT.CREATE"), server.CreatePaymentApi)
		invoiceGroup.GET("/:id/payments", server.RBACMiddleware("INVOICE.PAYMENT.VIEW"), server.ListPaymentsApi)
		invoiceGroup.GET("/:id/payments/:payment_id", server.RBACMiddleware("INVOICE.PAYMENT.VIEW"), server.GetPaymentByIDApi)
		invoiceGroup.PUT("/:id/payments/:payment_id", server.RBACMiddleware("INVOICE.PAYMENT.UPDATE"), server.UpdatePaymentApi)
		invoiceGroup.DELETE("/:id/payments/:payment_id", server.RBACMiddleware("INVOICE.PAYMENT.DELETE"), server.DeletePaymentApi)

		invoiceGroup.GET("/:id/audit", server.RBACMiddleware("INVOICE.VIEW"), server.GetInvoiceAuditLogApi)

		invoiceGroup.GET("/template_items", server.RBACMiddleware("INVOICE.VIEW"), server.GetInvoiceTemplateItemsApi)
	}
}
