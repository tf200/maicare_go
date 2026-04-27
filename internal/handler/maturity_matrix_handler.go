package handler

import (
	"net/http"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/gin-gonic/gin"
)

type MaturityMatrixHandler struct {
	service domain.MaturityMatrixService
}

func NewMaturityMatrixHandler(service domain.MaturityMatrixService) *MaturityMatrixHandler {
	return &MaturityMatrixHandler{service: service}
}

func (h *MaturityMatrixHandler) ListMaturityMatrix(ctx *gin.Context) {
	result, err := h.service.ListMaturityMatrix(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpapi.Fail("failed to list maturity matrix", ""))
		return
	}

	ctx.JSON(http.StatusOK, httpapi.OK(result, "Maturity matrix retrieved successfully"))
}
