package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAuthRoutes(baseRouter *gin.RouterGroup) {
	authGroup := baseRouter.Group("/auth")
	authGroup.POST("/token", server.Login)
	authGroup.POST("/refresh", server.RefreshToken)

	authGroup.POST("/change_password", AuthMiddleware(server.tokenMaker), server.ChangePasswordApi)
}
