package organization

import (
	"time"

	"github.com/google/uuid"
)

// CreateOrganisationRequest represents a request to create an organisation
type CreateOrganisationRequest struct {
	Name       string  `json:"name" binding:"required"`
	Address    string  `json:"address" binding:"required"`
	PostalCode string  `json:"postal_code" binding:"required"`
	City       string  `json:"city" binding:"required"`
	Email      *string `json:"email"`
	KvkNumber  *string `json:"kvk_number"`
	BtwNumber  *string `json:"btw_number"`
}

// CreateOrganisationResponse represents a response for CreateOrganisationApi
type CreateOrganisationResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
	City       string    `json:"city"`
	Email      *string   `json:"email"`
	KvkNumber  *string   `json:"kvk_number"`
	BtwNumber  *string   `json:"btw_number"`
}

// ListOrganisationsResponse represents an organisation in the list
type ListOrganisationsResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	PostalCode    string    `json:"postal_code"`
	City          string    `json:"city"`
	Email         *string   `json:"email"`
	KvkNumber     *string   `json:"kvk_number"`
	BtwNumber     *string   `json:"btw_number"`
	LocationCount int64     `json:"location_count"`
}

// GetOrganisationResponse represents a response for GetOrganisationApi
type GetOrganisationResponse struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`
	Address       string    `json:"address"`
	PostalCode    string    `json:"postal_code"`
	City          string    `json:"city"`
	Email         *string   `json:"email"`
	KvkNumber     *string   `json:"kvk_number"`
	BtwNumber     *string   `json:"btw_number"`
	LocationCount int64     `json:"location_count"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GetOrganisationCountResponse represents a response for GetOrganisationCountApi
type GetOrganisationCountResponse struct {
	OrganisationID   uuid.UUID `json:"organisation_id"`
	OrganisationName string    `json:"organisation_name"`
	LocationCount    int64     `json:"location_count"`
	ClientCount      int64     `json:"client_count"`
	EmployeeCount    int64     `json:"employee_count"`
}

// UpdateOrganisationRequest represents a request to update an organisation
type UpdateOrganisationRequest struct {
	Name       *string `json:"name"`
	Address    *string `json:"address"`
	PostalCode *string `json:"postal_code"`
	City       *string `json:"city"`
	Email      *string `json:"email"`
	KvkNumber  *string `json:"kvk_number"`
	BtwNumber  *string `json:"btw_number"`
}

// UpdateOrganisationResponse represents a response for UpdateOrganisationApi
type UpdateOrganisationResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	PostalCode string    `json:"postal_code"`
	City       string    `json:"city"`
	Email      *string   `json:"email"`
	KvkNumber  *string   `json:"kvk_number"`
	BtwNumber  *string   `json:"btw_number"`
}

// DeleteOrganisationResponse represents a response for DeleteOrganisationApi
type DeleteOrganisationResponse struct {
	ID uuid.UUID `json:"id"`
}
