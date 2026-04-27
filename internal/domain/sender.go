package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Status: scaffolded, unwired

// ==================== Domain Types ====================

// Named with "Entity" suffix to avoid collision with domain.Sender in client.go
// which represents a client network/referral sender.

type SenderContactInfo struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

type InvoiceTemplateItem struct {
	ID           uuid.UUID
	ItemTag      string
	Description  string
	SourceTable  string
	SourceColumn string
}

type SenderEntity struct {
	ID                  uuid.UUID
	Types               string
	Name                string
	Street              *string
	HouseNumber         *string
	HouseNumberAddition *string
	PostalCode          *string
	City                *string
	Land                *string
	KVKNumber           *string
	BTWNumber           *string
	PhoneNumber         *string
	ClientNumber        *string
	EmailAddress        *string
	Contacts            []SenderContactInfo
	InvoiceTemplate     []uuid.UUID
	IsArchived          bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type SenderEntityWithTemplate struct {
	SenderEntity
	InvoiceTemplateItems []InvoiceTemplateItem
}

// ==================== Parameter Structs ====================

type CreateSenderParams struct {
	Types               string
	Name                string
	Street              *string
	HouseNumber         *string
	HouseNumberAddition *string
	PostalCode          *string
	City                *string
	Land                *string
	KVKNumber           *string
	BTWNumber           *string
	PhoneNumber         *string
	ClientNumber        *string
	Contacts            []SenderContactInfo
}

type UpdateSenderParams struct {
	ID                  uuid.UUID
	Name                *string
	Street              *string
	HouseNumber         *string
	HouseNumberAddition *string
	PostalCode          *string
	City                *string
	Land                *string
	KVKNumber           *string
	BTWNumber           *string
	PhoneNumber         *string
	ClientNumber        *string
	EmailAddress        *string
	Contacts            []SenderContactInfo
	IsArchived          *bool
	Types               *string
}

type ListSendersParams struct {
	Limit           int32
	Offset          int32
	Search          *string
	IncludeArchived *bool
}

// ==================== Repository Interface ====================

type SenderRepository interface {
	Create(ctx context.Context, params CreateSenderParams) (*SenderEntity, error)
	List(ctx context.Context, params ListSendersParams) ([]SenderEntity, int64, error)
	GetByID(ctx context.Context, id uuid.UUID) (*SenderEntity, error)
	Update(ctx context.Context, params UpdateSenderParams) (*SenderEntity, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CreateInvoiceTemplate(ctx context.Context, senderID uuid.UUID, templateIDs []uuid.UUID) error
	GetTemplateItemsBySourceTable(ctx context.Context, ids []uuid.UUID) ([]InvoiceTemplateItem, error)
	GetTemplateItemsByIDs(ctx context.Context, ids []uuid.UUID) ([]uuid.UUID, error)
}

// ==================== Service Interface ====================

type SenderService interface {
	CreateSender(ctx context.Context, params CreateSenderParams) (*SenderEntity, error)
	ListSenders(ctx context.Context, params ListSendersParams) (*ListResult[SenderEntity], error)
	GetSenderByID(ctx context.Context, id uuid.UUID) (*SenderEntityWithTemplate, error)
	UpdateSender(ctx context.Context, params UpdateSenderParams) (*SenderEntity, error)
	DeleteSender(ctx context.Context, id uuid.UUID) error
	CreateSenderInvoiceTemplate(ctx context.Context, senderID uuid.UUID, templateIDs []uuid.UUID) error
}
