package sender

import (
	"maicare_go/pagination"
	"time"
)

// Contact represents a contact information.
type SenderContact struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"email"`
	PhoneNumber *string `json:"phone_number"`
}

// CreateSenderRequest represents a request to create a new sender.
type CreateSenderRequest struct {
	Types        string          `json:"types" binding:"required,oneof=main_provider local_authority particular_party healthcare_institution"`
	Name         string          `json:"name" binding:"required"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	KVKNumber    *string         `json:"KVKnumber"`
	BTWNumber    *string         `json:"BTWnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	Contacts     []SenderContact `json:"contacts" binding:"dive"`
}

// CreateSenderResponse represents a response to a request to create a new sender.
type CreateSenderResponse struct {
	ID           int64           `json:"id"`
	Types        string          `json:"types"`
	Name         string          `json:"name"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	KVKNumber    *string         `json:"KVKnumber"`
	BTWNumber    *string         `json:"BTWnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	Contacts     []SenderContact `json:"contacts"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

// GetSenderRequest represents a request to get a sender by ID.
type ListSendersRequest struct {
	pagination.Request
	IncludeArchived *bool   `form:"include_archived"`
	Search          *string `form:"search"`
}

// GetSenderResponse represents a response to a request to get a sender by ID.
type ListSendersResponse struct {
	ID           int64           `json:"id"`
	Types        string          `json:"types"`
	Name         string          `json:"name"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	KVKNumber    *string         `json:"KVKnumber"`
	BTWNumber    *string         `json:"BTWnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	Contacts     []SenderContact `json:"contacts"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

type TemplateItem struct {
	ID           int64  `json:"id"`
	ItemTag      string `json:"item_tag"`
	Description  string `json:"description"`
	SourceTable  string `json:"source_table"`
	SourceColumn string `json:"source_column"`
}

// GetSenderByIdResponse represents a response to a request to get a sender by ID.
type GetSenderByIdResponse struct {
	ID                   int64           `json:"id"`
	Types                string          `json:"types"`
	Name                 string          `json:"name"`
	Address              *string         `json:"address"`
	PostalCode           *string         `json:"postal_code"`
	Place                *string         `json:"place"`
	Land                 *string         `json:"land"`
	Kvknumber            *string         `json:"KVKnumber"`
	Btwnumber            *string         `json:"BTWnumber"`
	PhoneNumber          *string         `json:"phone_number"`
	ClientNumber         *string         `json:"client_number"`
	EmailAddress         *string         `json:"email_address"`
	Contacts             []SenderContact `json:"contacts"`
	InvoiceTemplateItems []TemplateItem  `json:"invoice_template_items"`
	IsArchived           bool            `json:"is_archived"`
	CreatedAt            time.Time       `json:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at"`
}

// UpdateSenderRequest represents a request to update a sender.
type UpdateSenderRequest struct {
	Name         *string         `json:"name"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	Kvknumber    *string         `json:"KVKnumber"`
	Btwnumber    *string         `json:"BTWnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	EmailAddress *string         `json:"email_address"`
	Contacts     []SenderContact `json:"contacts"`
	IsArchived   *bool           `json:"is_archived"`
	Types        *string         `json:"types" binding:"omitempty,oneof=main_provider local_authority particular_party healthcare_institution"`
}

// UpdateSenderResponse represents a response to a request to update a sender.
type UpdateSenderResponse struct {
	ID           int64           `json:"id"`
	Types        string          `json:"types"`
	Name         string          `json:"name"`
	Address      *string         `json:"address"`
	PostalCode   *string         `json:"postal_code"`
	Place        *string         `json:"place"`
	Land         *string         `json:"land"`
	Kvknumber    *string         `json:"KVKnumber"`
	Btwnumber    *string         `json:"BTWnumber"`
	PhoneNumber  *string         `json:"phone_number"`
	ClientNumber *string         `json:"client_number"`
	EmailAddress *string         `json:"email_address"`
	Contacts     []SenderContact `json:"contacts"`
	IsArchived   bool            `json:"is_archived"`
	CreatedAt    string          `json:"created_at"`
	UpdatedAt    string          `json:"updated_at"`
}

// CreateSenderInvoiceTemplateRequest represents a request to create a new sender invoice template.
type CreateSenderInvoiceTemplateRequest struct {
	InvoiceTemplateIDs []int64 `json:"invoice_template" binding:"required"`
}
