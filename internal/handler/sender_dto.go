package handler

import (
	"time"

	"maicare_go/internal/domain"
	"maicare_go/internal/httpapi"

	"github.com/google/uuid"
)

// ==================== Request DTOs ====================

type createSenderRequest struct {
	Types               string              `json:"types" binding:"required,oneof=main_provider local_authority particular_party healthcare_institution"`
	Name                string              `json:"name" binding:"required"`
	Street              *string             `json:"street"`
	HouseNumber         *string             `json:"house_number"`
	HouseNumberAddition *string             `json:"house_number_addition"`
	PostalCode          *string             `json:"postal_code"`
	City                *string             `json:"city"`
	Land                *string             `json:"land"`
	KVKNumber           *string             `json:"KVKnumber"`
	BTWNumber           *string             `json:"BTWnumber"`
	PhoneNumber         *string             `json:"phone_number"`
	ClientNumber        *string             `json:"client_number"`
	Contacts            []senderContactDTO  `json:"contacts" binding:"dive"`
}

type senderContactDTO struct {
	Name        *string `json:"name"`
	Email       *string `json:"email" binding:"email"`
	PhoneNumber *string `json:"phone_number"`
}

type updateSenderRequest struct {
	Name                *string            `json:"name"`
	Street              *string            `json:"street"`
	HouseNumber         *string            `json:"house_number"`
	HouseNumberAddition *string            `json:"house_number_addition"`
	PostalCode          *string            `json:"postal_code"`
	City                *string            `json:"city"`
	Land                *string            `json:"land"`
	KVKNumber           *string            `json:"KVKnumber"`
	BTWNumber           *string            `json:"BTWnumber"`
	PhoneNumber         *string            `json:"phone_number"`
	ClientNumber        *string            `json:"client_number"`
	EmailAddress        *string            `json:"email_address"`
	Contacts            []senderContactDTO `json:"contacts"`
	IsArchived          *bool              `json:"is_archived"`
	Types               *string            `json:"types" binding:"omitempty,oneof=main_provider local_authority particular_party healthcare_institution"`
}

type listSendersRequest struct {
	httpapi.PageRequest
	Search          *string `form:"search" binding:"omitempty"`
	IncludeArchived *bool   `form:"include_archived"`
}

type createSenderInvoiceTemplateRequest struct {
	InvoiceTemplateIDs []uuid.UUID `json:"invoice_template" binding:"required"`
}

// ==================== Response DTOs ====================

type senderResponse struct {
	ID                  uuid.UUID           `json:"id"`
	Types               string              `json:"types"`
	Name                string              `json:"name"`
	Street              *string             `json:"street"`
	HouseNumber         *string             `json:"house_number"`
	HouseNumberAddition *string             `json:"house_number_addition"`
	PostalCode          *string             `json:"postal_code"`
	City                *string             `json:"city"`
	Land                *string             `json:"land"`
	KVKNumber           *string             `json:"KVKnumber"`
	BTWNumber           *string             `json:"BTWnumber"`
	PhoneNumber         *string             `json:"phone_number"`
	ClientNumber        *string             `json:"client_number"`
	EmailAddress        *string             `json:"email_address,omitempty"`
	Contacts            []senderContactDTO  `json:"contacts"`
	IsArchived          bool                `json:"is_archived,omitempty"`
	ClientsCount        *int                `json:"clients_count,omitempty"`
	CreatedAt           string              `json:"created_at"`
	UpdatedAt           string              `json:"updated_at"`
}

type templateItemResponse struct {
	ID           uuid.UUID `json:"id"`
	ItemTag      string    `json:"item_tag"`
	Description  string    `json:"description"`
	SourceTable  string    `json:"source_table"`
	SourceColumn string    `json:"source_column"`
}

type senderDetailResponse struct {
	ID                   uuid.UUID            `json:"id"`
	Types                string               `json:"types"`
	Name                 string               `json:"name"`
	Street               *string              `json:"street"`
	HouseNumber          *string              `json:"house_number"`
	HouseNumberAddition  *string              `json:"house_number_addition"`
	PostalCode           *string              `json:"postal_code"`
	City                 *string              `json:"city"`
	Land                 *string              `json:"land"`
	KVKNumber            *string              `json:"KVKnumber"`
	BTWNumber            *string              `json:"BTWnumber"`
	PhoneNumber          *string              `json:"phone_number"`
	ClientNumber         *string              `json:"client_number"`
	EmailAddress         *string              `json:"email_address"`
	Contacts             []senderContactDTO   `json:"contacts"`
	InvoiceTemplateItems []templateItemResponse `json:"invoice_template_items"`
	IsArchived           bool                 `json:"is_archived"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

// ==================== Mappers ====================

func toSenderResponse(s *domain.SenderEntity) senderResponse {
	contacts := make([]senderContactDTO, len(s.Contacts))
	for i, c := range s.Contacts {
		contacts[i] = senderContactDTO{
			Name:        c.Name,
			Email:       c.Email,
			PhoneNumber: c.PhoneNumber,
		}
	}

	return senderResponse{
		ID:                  s.ID,
		Types:               s.Types,
		Name:                s.Name,
		Street:              s.Street,
		HouseNumber:         s.HouseNumber,
		HouseNumberAddition: s.HouseNumberAddition,
		PostalCode:          s.PostalCode,
		City:                s.City,
		Land:                s.Land,
		KVKNumber:           s.KVKNumber,
		BTWNumber:           s.BTWNumber,
		PhoneNumber:         s.PhoneNumber,
		ClientNumber:        s.ClientNumber,
		Contacts:            contacts,
		CreatedAt:           s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:           s.UpdatedAt.Format(time.RFC3339),
	}
}

func toSenderDetailResponse(s *domain.SenderEntityWithTemplate) senderDetailResponse {
	contacts := make([]senderContactDTO, len(s.Contacts))
	for i, c := range s.Contacts {
		contacts[i] = senderContactDTO{
			Name:        c.Name,
			Email:       c.Email,
			PhoneNumber: c.PhoneNumber,
		}
	}

	templateItems := make([]templateItemResponse, len(s.InvoiceTemplateItems))
	for i, item := range s.InvoiceTemplateItems {
		templateItems[i] = templateItemResponse{
			ID:           item.ID,
			ItemTag:      item.ItemTag,
			Description:  item.Description,
			SourceTable:  item.SourceTable,
			SourceColumn: item.SourceColumn,
		}
	}

	return senderDetailResponse{
		ID:                   s.ID,
		Types:                s.Types,
		Name:                 s.Name,
		Street:               s.Street,
		HouseNumber:          s.HouseNumber,
		HouseNumberAddition:  s.HouseNumberAddition,
		PostalCode:           s.PostalCode,
		City:                 s.City,
		Land:                 s.Land,
		KVKNumber:            s.KVKNumber,
		BTWNumber:            s.BTWNumber,
		PhoneNumber:          s.PhoneNumber,
		ClientNumber:         s.ClientNumber,
		EmailAddress:         s.EmailAddress,
		Contacts:             contacts,
		InvoiceTemplateItems: templateItems,
		IsArchived:           s.IsArchived,
		CreatedAt:            s.CreatedAt,
		UpdatedAt:            s.UpdatedAt,
	}
}

func toCreateSenderParams(req createSenderRequest) domain.CreateSenderParams {
	contacts := make([]domain.SenderContactInfo, len(req.Contacts))
	for i, c := range req.Contacts {
		contacts[i] = domain.SenderContactInfo{
			Name:        c.Name,
			Email:       c.Email,
			PhoneNumber: c.PhoneNumber,
		}
	}

	return domain.CreateSenderParams{
		Types:               req.Types,
		Name:                req.Name,
		Street:              req.Street,
		HouseNumber:         req.HouseNumber,
		HouseNumberAddition: req.HouseNumberAddition,
		PostalCode:          req.PostalCode,
		City:                req.City,
		Land:                req.Land,
		KVKNumber:           req.KVKNumber,
		BTWNumber:           req.BTWNumber,
		PhoneNumber:         req.PhoneNumber,
		ClientNumber:        req.ClientNumber,
		Contacts:            contacts,
	}
}

func toUpdateSenderParams(id uuid.UUID, req updateSenderRequest) domain.UpdateSenderParams {
	params := domain.UpdateSenderParams{
		ID:                  id,
		Name:                req.Name,
		Street:              req.Street,
		HouseNumber:         req.HouseNumber,
		HouseNumberAddition: req.HouseNumberAddition,
		PostalCode:          req.PostalCode,
		City:                req.City,
		Land:                req.Land,
		KVKNumber:           req.KVKNumber,
		BTWNumber:           req.BTWNumber,
		PhoneNumber:         req.PhoneNumber,
		ClientNumber:        req.ClientNumber,
		EmailAddress:        req.EmailAddress,
		IsArchived:          req.IsArchived,
		Types:               req.Types,
	}

	if req.Contacts != nil {
		contacts := make([]domain.SenderContactInfo, len(req.Contacts))
		for i, c := range req.Contacts {
			contacts[i] = domain.SenderContactInfo{
				Name:        c.Name,
				Email:       c.Email,
				PhoneNumber: c.PhoneNumber,
			}
		}
		params.Contacts = contacts
	}

	return params
}

func toListSendersParams(req listSendersRequest) domain.ListSendersParams {
	pageParams := req.Params()
	return domain.ListSendersParams{
		Limit:           pageParams.Limit,
		Offset:          pageParams.Offset,
		Search:          req.Search,
		IncludeArchived: req.IncludeArchived,
	}
}
