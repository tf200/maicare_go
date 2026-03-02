package organization

import (
	"time"

	"maicare_go/pagination"

	"github.com/google/uuid"
)

// ListAllLocationsRequest represents a request to list all locations with pagination
type ListAllLocationsRequest struct {
	pagination.Request
	Search *string `form:"search"`
}

// ListLocationsRequest represents a request to list locations with pagination
type ListLocationsRequest struct {
	pagination.Request
	Search *string `form:"search"`
}

// ListLocationsResponse represents a location in the list
type ListLocationsResponse struct {
	ID                  uuid.UUID                        `json:"id"`
	Name                string                           `json:"name"`
	Street              string                           `json:"street"`
	HouseNumber         string                           `json:"house_number"`
	HouseNumberAddition *string                          `json:"house_number_addition"`
	PostalCode          string                           `json:"postal_code"`
	City                string                           `json:"city"`
	Capacity            *int32                           `json:"capacity"`
	Occupied            int32                            `json:"occupied"`
	Available           int32                            `json:"available"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           time.Time                        `json:"updated_at"`
	Shifts              []ListShiftsByLocationIDResponse `json:"shifts"`
}

// ListOrgLocationsResponse represents an organization-scoped location in the list.
// It includes location shifts for each location.
type ListOrgLocationsResponse struct {
	ID                  uuid.UUID                        `json:"id"`
	Name                string                           `json:"name"`
	Street              string                           `json:"street"`
	HouseNumber         string                           `json:"house_number"`
	HouseNumberAddition *string                          `json:"house_number_addition"`
	PostalCode          string                           `json:"postal_code"`
	City                string                           `json:"city"`
	Capacity            *int32                           `json:"capacity"`
	Occupied            int32                            `json:"occupied"`
	Available           int32                            `json:"available"`
	CreatedAt           time.Time                        `json:"created_at"`
	UpdatedAt           time.Time                        `json:"updated_at"`
	Shifts              []ListShiftsByLocationIDResponse `json:"shifts"`
}

// CreateLocationRequest represents a request to create a location
type CreateLocationRequest struct {
	Name                string  `json:"name" binding:"required"`
	Street              string  `json:"street" binding:"required"`
	HouseNumber         string  `json:"house_number" binding:"required"`
	HouseNumberAddition *string `json:"house_number_addition"`
	PostalCode          string  `json:"postal_code" binding:"required"`
	City                string  `json:"city" binding:"required"`
	Capacity            *int32  `json:"capacity"`
}

// CreateLocationResponse represents a response for CreateLocationApi
type CreateLocationResponse struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Street              string    `json:"street"`
	HouseNumber         string    `json:"house_number"`
	HouseNumberAddition *string   `json:"house_number_addition"`
	PostalCode          string    `json:"postal_code"`
	City                string    `json:"city"`
	Capacity            *int32    `json:"capacity"`
}

// UpdateLocationRequest represents a request to update a location
type UpdateLocationRequest struct {
	Name                *string `json:"name"`
	Street              *string `json:"street"`
	HouseNumber         *string `json:"house_number"`
	HouseNumberAddition *string `json:"house_number_addition"`
	PostalCode          *string `json:"postal_code"`
	City                *string `json:"city"`
	Capacity            *int32  `json:"capacity"`
}

// UpdateLocationResponse represents a response for UpdateLocationApi
type UpdateLocationResponse struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Street              string    `json:"street"`
	HouseNumber         string    `json:"house_number"`
	HouseNumberAddition *string   `json:"house_number_addition"`
	PostalCode          string    `json:"postal_code"`
	City                string    `json:"city"`
	Capacity            *int32    `json:"capacity"`
}

// DeleteLocationResponse represents a response for DeleteLocationApi
type DeleteLocationResponse struct {
	ID uuid.UUID `json:"id"`
}

// GetLocationResponse represents a response for GetLocationApi
type GetLocationResponse struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Street              string    `json:"street"`
	HouseNumber         string    `json:"house_number"`
	HouseNumberAddition *string   `json:"house_number_addition"`
	PostalCode          string    `json:"postal_code"`
	City                string    `json:"city"`
	Capacity            *int32    `json:"capacity"`
}
