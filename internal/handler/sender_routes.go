package handler

import "github.com/gin-gonic/gin"

func RegisterSenderRoutes(
	rg *gin.RouterGroup,
	handler *SenderHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	senders := rg.Group("/senders")
	{
		senders.POST("", auth, requirePermission("SENDER.CREATE"), handler.CreateSender)
		senders.GET("", auth, requirePermission("SENDER.VIEW"), handler.ListSenders)
		senders.GET("/:id", auth, requirePermission("SENDER.VIEW"), handler.GetSenderByID)
		senders.PUT("/:id", auth, requirePermission("SENDER.UPDATE"), handler.UpdateSender)
		senders.DELETE("/:id", auth, requirePermission("SENDER.DELETE"), handler.DeleteSender)
		senders.POST("/:id/invoice_template", auth, requirePermission("SENDER.CREATE"), handler.CreateSenderInvoiceTemplate)
	}
}
