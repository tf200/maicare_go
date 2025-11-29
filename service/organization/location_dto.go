package organization

import (
	"time"

	"github.com/google/uuid"
)

// ListLocationsResponse represents a location in the list
type ListLocationsResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Capacity  *int32    `json:"capacity"`
	Occupied  int32     `json:"occupied"`
	Available int32     `json:"available"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateLocationRequest represents a request to create a location
type CreateLocationRequest struct {
	Name     string `json:"name" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Capacity *int32 `json:"capacity"`
}

// CreateLocationResponse represents a response for CreateLocationApi
type CreateLocationResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Capacity *int32    `json:"capacity"`
}

// UpdateLocationRequest represents a request to update a location
type UpdateLocationRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	Capacity *int32  `json:"capacity"`
}

// UpdateLocationResponse represents a response for UpdateLocationApi
type UpdateLocationResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Capacity *int32    `json:"capacity"`
}

// DeleteLocationResponse represents a response for DeleteLocationApi
type DeleteLocationResponse struct {
	ID uuid.UUID `json:"id"`
}

// GetLocationResponse represents a response for GetLocationApi
type GetLocationResponse struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Address  string    `json:"address"`
	Capacity *int32    `json:"capacity"`
}
