// // auth_routes.go
package api

import "github.com/gin-gonic/gin"

func (server *Server) setupMaturityMatrixRoutes(baseRouter *gin.RouterGroup) {

	baseRouter.GET("/maturity_matrix", server.ListMaturityMatrixApi)
}
