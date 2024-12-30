
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupAuthRoutes(baseRouter *gin.RouterGroup) {

	baseRouter.POST("/token", server.Login)
	baseRouter.POST("/refresh", server.RefreshToken)

}
