package service

import (
	"context"

	"maicare_go/internal/domain"

	"go.uber.org/zap"
)

type DashboardService struct {
	repository domain.DashboardRepository
	logger     domain.Logger
	audit      domain.AuditLogger
}

func NewDashboardService(repository domain.DashboardRepository, logger domain.Logger, audit domain.AuditLogger) domain.DashboardService {
	return &DashboardService{repository: repository, logger: logger, audit: audit}
}

func (s *DashboardService) GetAdminDashboardStats(ctx context.Context) (*domain.AdminDashboardStats, error) {
	stats, err := s.repository.GetAdminDashboardStats(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.LogError(ctx, "DashboardService.GetAdminDashboardStats", "failed to get admin dashboard stats", err)
		}
		return nil, err
	}

	if s.logger != nil {
		s.logger.LogInfo(ctx, "DashboardService.GetAdminDashboardStats", "admin dashboard stats retrieved successfully")
	}

	if s.audit != nil {
		if auditErr := s.audit.Log(ctx, domain.AuditEvent{
			EventType:   "record_access",
			Action:      "read",
			Result:      "success",
			SubjectType: "dashboard",
			SubjectID:   "admin",
			AccessRule:  strPtr("DASHBOARD.VIEW"),
			Details: map[string]any{
				"metric":                 "admin_dashboard_stats",
				"clients_total":          stats.Clients.Total,
				"clients_in_care":        stats.Clients.InCare,
				"clients_waiting_list":   stats.Clients.WaitingList,
				"employees_total":        stats.Employees.Total,
				"incidents_today":        stats.Incidents.Today,
				"overdue_invoices_total": stats.Invoices.Overdue,
			},
		}); auditErr != nil && s.logger != nil {
			s.logger.LogError(ctx, "DashboardService.GetAdminDashboardStats", "audit log failed", auditErr, zap.String("dashboard", "admin"))
		}
	}

	return stats, nil
}
