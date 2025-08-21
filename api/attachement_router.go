package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAttachementRoutes(baseRouter *gin.RouterGroup) {
	attachmentRouter := baseRouter.Group("/attachments").Use(server.AuthMiddleware())
	{
		attachmentRouter.POST("/upload", server.UploadHandlerApi)
		attachmentRouter.GET("/:id", server.GetAttachmentByIdApi)
	}
}
