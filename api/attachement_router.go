package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAttachementRoutes(baseRouter *gin.RouterGroup) {
	attachmentRouter := baseRouter.Group("/attachment").Use(AuthMiddleware(server.tokenMaker))
	{
		attachmentRouter.POST("/upload", server.UploadHandlerApi)
		attachmentRouter.GET("/:id", server.GetAttachmentByIdApi)
	}
}
