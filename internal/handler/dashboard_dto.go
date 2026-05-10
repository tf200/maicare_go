package handler

import (
	"time"

	"maicare_go/internal/domain"
)

type adminDashboardResponse struct {
	StatCards          adminDashboardStatCardsResponse          `json:"stat_cards"`
	TopAdminActions    adminDashboardTopActionsResponse         `json:"top_admin_actions"`
	RegistrationsToday adminDashboardRegistrationsTodayResponse `json:"registrations_today"`
}

type adminDashboardStatCardsResponse struct {
	Clients   adminDashboardClientStatsResponse   `json:"clients"`
	Employees adminDashboardEmployeeStatsResponse `json:"employees"`
	Incidents adminDashboardIncidentStatsResponse `json:"incidents"`
	Invoices  adminDashboardInvoiceStatsResponse  `json:"invoices"`
}

type adminDashboardClientStatsResponse struct {
	Total       int64 `json:"total"`
	InCare      int64 `json:"in_care"`
	WaitingList int64 `json:"waiting_list"`
}

type adminDashboardEmployeeStatsResponse struct {
	Total int64 `json:"total"`
}

type adminDashboardIncidentStatsResponse struct {
	Today int64 `json:"today"`
}

type adminDashboardInvoiceStatsResponse struct {
	Overdue int64 `json:"overdue"`
}

type adminDashboardTopActionsResponse struct {
	Items []adminDashboardActionItemResponse `json:"items"`
}

type adminDashboardActionItemResponse struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Severity    string     `json:"severity"`
	Title       string     `json:"title"`
	Subtitle    string     `json:"subtitle"`
	ActionLabel string     `json:"action_label"`
	ActionURL   string     `json:"action_url"`
	OccurredAt  *time.Time `json:"occurred_at,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	SortAt      time.Time  `json:"sort_at"`
}

type adminDashboardRegistrationsTodayResponse struct {
	Items []adminDashboardRegistrationTodayItemResponse `json:"items"`
}

type adminDashboardRegistrationTodayItemResponse struct {
	ID                   string     `json:"id"`
	ClientFirstName      string     `json:"client_first_name"`
	ClientLastName       string     `json:"client_last_name"`
	ReferrerFirstName    string     `json:"referrer_first_name"`
	ReferrerLastName     string     `json:"referrer_last_name"`
	ReferrerOrganization string     `json:"referrer_organization"`
	FormStatus           string     `json:"form_status"`
	RiskCount            int32      `json:"risk_count"`
	SubmittedAt          *time.Time `json:"submitted_at,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	ActionURL            string     `json:"action_url"`
}

func toAdminDashboardResponse(stats *domain.AdminDashboardStats) adminDashboardResponse {
	actions := make([]adminDashboardActionItemResponse, len(stats.TopAdminActions))
	for i, item := range stats.TopAdminActions {
		actions[i] = adminDashboardActionItemResponse{
			ID:          item.ID,
			Type:        item.Type,
			Severity:    item.Severity,
			Title:       item.Title,
			Subtitle:    item.Subtitle,
			ActionLabel: item.ActionLabel,
			ActionURL:   item.ActionURL,
			OccurredAt:  item.OccurredAt,
			DueDate:     item.DueDate,
			SortAt:      item.SortAt,
		}
	}
	registrationsToday := make([]adminDashboardRegistrationTodayItemResponse, len(stats.RegistrationsToday))
	for i, item := range stats.RegistrationsToday {
		registrationsToday[i] = adminDashboardRegistrationTodayItemResponse{
			ID:                   item.ID.String(),
			ClientFirstName:      item.ClientFirstName,
			ClientLastName:       item.ClientLastName,
			ReferrerFirstName:    item.ReferrerFirstName,
			ReferrerLastName:     item.ReferrerLastName,
			ReferrerOrganization: item.ReferrerOrganization,
			FormStatus:           item.FormStatus,
			RiskCount:            item.RiskCount,
			SubmittedAt:          item.SubmittedAt,
			CreatedAt:            item.CreatedAt,
			ActionURL:            item.ActionURL,
		}
	}

	return adminDashboardResponse{
		StatCards: adminDashboardStatCardsResponse{
			Clients: adminDashboardClientStatsResponse{
				Total:       stats.Clients.Total,
				InCare:      stats.Clients.InCare,
				WaitingList: stats.Clients.WaitingList,
			},
			Employees: adminDashboardEmployeeStatsResponse{Total: stats.Employees.Total},
			Incidents: adminDashboardIncidentStatsResponse{Today: stats.Incidents.Today},
			Invoices:  adminDashboardInvoiceStatsResponse{Overdue: stats.Invoices.Overdue},
		},
		TopAdminActions:    adminDashboardTopActionsResponse{Items: actions},
		RegistrationsToday: adminDashboardRegistrationsTodayResponse{Items: registrationsToday},
	}
}
