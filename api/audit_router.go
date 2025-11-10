package api

import "github.com/gin-gonic/gin"

func (s *Server) setupAuditRoutes(router *gin.RouterGroup) {
	auditRoutes := router.Group("/audit")
	{
		auditRoutes.GET("/logs", s.RBACMiddleware("AUDIT_LOGS.VIEW"), s.ListAuditLogs)
	}
}
