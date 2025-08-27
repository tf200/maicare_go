package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAuthRoutes(baseRouter *gin.RouterGroup) {
	authGroup := baseRouter.Group("/auth")
	authGroup.POST("/token", server.Login)
	authGroup.POST("/refresh", server.RefreshToken)
	authGroup.POST("/verify_2fa", server.Verify2FAHandler)

	// 2fa setup routes
	authGroup.POST("/setup_2fa", server.AuthMiddleware(), server.Setup2FAHandler)
	authGroup.POST("/enable_2fa", server.AuthMiddleware(), server.Enable2FAHandler)

	authGroup.POST("/change_password", server.AuthMiddleware(), server.ChangePasswordApi)
}
