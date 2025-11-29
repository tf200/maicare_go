package organization

import "github.com/google/uuid"

// CreateShiftApi creates a new shift for a specific location
type CreateShiftApiRequest struct {
	ShiftName string `json:"shift"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// CreateShiftApiResponse represents the response structure for creating a shift
type CreateShiftApiResponse struct {
	ID         uuid.UUID `json:"id"`
	LocationID uuid.UUID `json:"location_id"`
	ShiftName  string    `json:"shift"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
}

// UpdateShiftApiRequest represents the request structure for updating a shift
type UpdateShiftApiRequest struct {
	ShiftName string `json:"shift"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

// UpdateShiftApiResponse represents the response structure for updating a shift
type UpdateShiftApiResponse struct {
	ID         uuid.UUID `json:"id"`
	LocationID uuid.UUID `json:"location_id"`
	ShiftName  string    `json:"shift"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
}

// ListShiftsByLocationIDResponse represents the response structure for listing shifts by location ID
type ListShiftsByLocationIDResponse struct {
	ID         uuid.UUID `json:"id"`
	LocationID uuid.UUID `json:"location_id"`
	ShiftName  string    `json:"shift"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
}
