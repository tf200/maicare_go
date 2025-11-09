package api

import (
	"maicare_go/service/audit"
	"net/http"

	"github.com/gin-gonic/gin"
)

// @Summary List Audit Logs
// @Description Retrieve a list of audit logs based on query parameters
// @Tags audit
// @Accept json
// @Produce json
// @Param request query audit.ListAuditRecordsRequest true "Audit log query parameters"
// @Success 200 {object} Response[[]audit.AuditRecord] "Successfully retrieved audit logs"
// @Failure 400 {object} Response[any] "Bad request - Invalid input"
// @Failure 500 {object} Response[any] "Internal server error"
// @Router /audit/logs [get]
// @Security ApiKeyAuth
func (s *Server) ListAuditLogs(ctx *gin.Context) {
	var req audit.ListAuditRecordsRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	records, err := s.businessService.AuditService.ListAuditRecords(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := SuccessResponse(records, "Audit records retrieved successfully")

	ctx.JSON(http.StatusOK, res)
}
