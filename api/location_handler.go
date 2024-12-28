package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) ListLocationsApi(ctx *gin.Context) {
	locations, err := server.store.ListLocations(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, locations)
}
