package handler

import "github.com/gin-gonic/gin"

func RegisterMaturityMatrixRoutes(
	rg *gin.RouterGroup,
	handler *MaturityMatrixHandler,
	auth gin.HandlerFunc,
	requirePermission func(string) gin.HandlerFunc,
) {
	rg.GET("/maturity_matrix", auth, handler.ListMaturityMatrix)
}
