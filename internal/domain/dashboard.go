package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AdminDashboardClientStats struct {
	Total       int64
	InCare      int64
	WaitingList int64
}

type AdminDashboardEmployeeStats struct {
	Total int64
}

type AdminDashboardIncidentStats struct {
	Today int64
}

type AdminDashboardInvoiceStats struct {
	Overdue int64
}

type AdminDashboardStats struct {
	Clients            AdminDashboardClientStats
	Employees          AdminDashboardEmployeeStats
	Incidents          AdminDashboardIncidentStats
	Invoices           AdminDashboardInvoiceStats
	TopAdminActions    []AdminDashboardActionItem
	RegistrationsToday []AdminDashboardRegistrationTodayItem
}

type AdminDashboardActionItem struct {
	ID          string
	Type        string
	Severity    string
	Title       string
	Subtitle    string
	ActionLabel string
	ActionURL   string
	OccurredAt  *time.Time
	DueDate     *time.Time
	SortAt      time.Time
}

type AdminDashboardRegistrationTodayItem struct {
	ID                   uuid.UUID
	ClientFirstName      string
	ClientLastName       string
	ReferrerFirstName    string
	ReferrerLastName     string
	ReferrerOrganization string
	FormStatus           string
	RiskCount            int32
	SubmittedAt          *time.Time
	CreatedAt            time.Time
	ActionURL            string
}

type DashboardRepository interface {
	GetAdminDashboardStats(ctx context.Context) (*AdminDashboardStats, error)
}

type DashboardService interface {
	GetAdminDashboardStats(ctx context.Context) (*AdminDashboardStats, error)
}
