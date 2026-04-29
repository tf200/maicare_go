package handler

import (
	"github.com/gin-gonic/gin"
)

func RegisterInvoiceRoutes(
	rg *gin.RouterGroup,
	handler *InvoiceHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	invoices := rg.Group("/invoices")
	{
		invoices.POST("", auth, requirePermission("INVOICE.CREATE"), handler.CreateInvoice)
		invoices.POST("/generate", auth, requirePermission("INVOICE.CREATE"), handler.GenerateInvoice)
		invoices.GET("", auth, requirePermission("INVOICE.VIEW"), handler.ListInvoices)
		invoices.GET("/template_items", auth, requirePermission("INVOICE.VIEW"), handler.GetInvoiceTemplateItems)
		invoices.GET("/:id", auth, requirePermission("INVOICE.VIEW"), handler.GetInvoiceByID)
		invoices.PUT("/:id", auth, requirePermission("INVOICE.UPDATE"), handler.UpdateInvoice)
		invoices.DELETE("/:id", auth, requirePermission("INVOICE.DELETE"), handler.DeleteInvoice)
		invoices.POST("/:id/credit", auth, requirePermission("INVOICE.UPDATE"), handler.CreditInvoice)
		invoices.GET("/:id/generate_pdf", auth, requirePermission("INVOICE.VIEW"), handler.GenerateInvoicePDF)
		invoices.POST("/:id/send_reminder", auth, requirePermission("INVOICE.CREATE"), handler.SendInvoiceReminder)
		invoices.GET("/:id/audit", auth, requirePermission("INVOICE.VIEW"), handler.GetInvoiceAuditLog)

		// Payment routes nested under invoices
		invoices.POST("/:id/payments", auth, requirePermission("INVOICE.PAYMENT.CREATE"), handler.CreatePayment)
		invoices.GET("/:id/payments", auth, requirePermission("INVOICE.PAYMENT.VIEW"), handler.ListPayments)
		invoices.GET("/:id/payments/:payment_id", auth, requirePermission("INVOICE.PAYMENT.VIEW"), handler.GetPaymentByID)
		invoices.PUT("/:id/payments/:payment_id", auth, requirePermission("INVOICE.PAYMENT.UPDATE"), handler.UpdatePayment)
		invoices.DELETE("/:id/payments/:payment_id", auth, requirePermission("INVOICE.PAYMENT.DELETE"), handler.DeletePayment)
	}
}
