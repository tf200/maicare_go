package organization

import "time"

// ListLocationsResponse represents a location in the list
type ListLocationsResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Capacity  *int32    `json:"capacity"`
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
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity *int32 `json:"capacity"`
}

// UpdateLocationRequest represents a request to update a location
type UpdateLocationRequest struct {
	Name     *string `json:"name"`
	Address  *string `json:"address"`
	Capacity *int32  `json:"capacity"`
}

// UpdateLocationResponse represents a response for UpdateLocationApi
type UpdateLocationResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity *int32 `json:"capacity"`
}

// DeleteLocationResponse represents a response for DeleteLocationApi
type DeleteLocationResponse struct {
	ID int64 `json:"id"`
}

// GetLocationResponse represents a response for GetLocationApi
type GetLocationResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity *int32 `json:"capacity"`
}
