package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAttachementRoutes(baseRouter *gin.RouterGroup) {

	baseRouter.POST("/upload", server.UploadHandlerApi)

}
