package repository

import (
	"context"
	"fmt"
	"strings"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
)

type dashboardRepository struct {
	store *db.Store
}

func NewDashboardRepository(store *db.Store) domain.DashboardRepository {
	return &dashboardRepository{store: store}
}

func (r *dashboardRepository) GetAdminDashboardStats(ctx context.Context) (*domain.AdminDashboardStats, error) {
	stats, err := r.getAdminDashboardStatCards(ctx)
	if err != nil {
		return nil, err
	}

	actions, err := r.listTopAdminActions(ctx, 4)
	if err != nil {
		return nil, err
	}
	stats.TopAdminActions = actions

	registrationsToday, err := r.listRegistrationsToday(ctx, 5)
	if err != nil {
		return nil, err
	}
	stats.RegistrationsToday = registrationsToday

	return stats, nil
}

func (r *dashboardRepository) getAdminDashboardStatCards(ctx context.Context) (*domain.AdminDashboardStats, error) {
	row, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.GetAdminDashboardStatCardsRow, error) {
		return q.GetAdminDashboardStatCards(ctx)
	})
	if err != nil {
		return nil, err
	}

	stats := domain.AdminDashboardStats{
		Clients: domain.AdminDashboardClientStats{
			InCare:      row.ClientsInCare,
			WaitingList: row.ClientsWaitingList,
		},
		Employees: domain.AdminDashboardEmployeeStats{Total: row.TotalEmployees},
		Incidents: domain.AdminDashboardIncidentStats{Today: row.IncidentsToday},
		Invoices:  domain.AdminDashboardInvoiceStats{Overdue: row.OverdueInvoices},
	}
	stats.Clients.Total = stats.Clients.InCare + stats.Clients.WaitingList
	return &stats, nil
}

func (r *dashboardRepository) listTopAdminActions(ctx context.Context, limit int32) ([]domain.AdminDashboardActionItem, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListTopAdminActionsRow, error) {
		return q.ListTopAdminActions(ctx, limit)
	})
	if err != nil {
		return nil, err
	}

	items := make([]domain.AdminDashboardActionItem, 0, len(rows))
	for _, row := range rows {
		item := domain.AdminDashboardActionItem{
			ID:          row.ID,
			Type:        row.Type,
			Severity:    row.Severity,
			Title:       strings.TrimSpace(row.Title),
			Subtitle:    strings.TrimSpace(row.Subtitle),
			ActionLabel: row.ActionLabel,
			ActionURL:   row.ActionUrl,
		}
		if row.OccurredAt.Valid {
			item.OccurredAt = &row.OccurredAt.Time
		}
		if row.DueDate.Valid {
			date := row.DueDate.Time
			item.DueDate = &date
		}
		if row.SortAt.Valid {
			item.SortAt = row.SortAt.Time
		}
		if item.Title == "" {
			item.Title = fmt.Sprintf("%s needs attention", strings.ReplaceAll(item.Type, "_", " "))
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *dashboardRepository) listRegistrationsToday(ctx context.Context, limit int32) ([]domain.AdminDashboardRegistrationTodayItem, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListDashboardRegistrationsTodayRow, error) {
		return q.ListDashboardRegistrationsToday(ctx, limit)
	})
	if err != nil {
		return nil, err
	}

	items := make([]domain.AdminDashboardRegistrationTodayItem, 0, len(rows))
	for _, row := range rows {
		item := domain.AdminDashboardRegistrationTodayItem{
			ID:                   row.ID,
			ClientFirstName:      row.ClientFirstName,
			ClientLastName:       row.ClientLastName,
			ReferrerFirstName:    row.ReferrerFirstName,
			ReferrerLastName:     row.ReferrerLastName,
			ReferrerOrganization: row.ReferrerOrganization,
			FormStatus:           string(row.FormStatus),
			RiskCount:            row.RiskCount,
			ActionURL:            row.ActionUrl,
		}
		if row.SubmittedAt.Valid {
			item.SubmittedAt = &row.SubmittedAt.Time
		}
		if row.CreatedAt.Valid {
			item.CreatedAt = row.CreatedAt.Time
		}
		items = append(items, item)
	}

	return items, nil
}
