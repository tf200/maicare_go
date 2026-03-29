package handler

import (
	"time"

	"maicare_go/internal/domain"
	"maicare_go/pagination"

	"github.com/google/uuid"
)

type listOrganizationsRequest struct {
	pagination.Request
	Search *string `form:"search"`
}

type listOrganizationsResponse struct {
	ID                  uuid.UUID `json:"id"`
	Name                string    `json:"name"`
	Street              string    `json:"street"`
	HouseNumber         string    `json:"house_number"`
	HouseNumberAddition *string   `json:"house_number_addition"`
	PostalCode          string    `json:"postal_code"`
	City                string    `json:"city"`
	Email               *string   `json:"email"`
	KvkNumber           *string   `json:"kvk_number"`
	BtwNumber           *string   `json:"btw_number"`
	LocationCount       int64     `json:"location_count"`
}

type listOrganizationLocationsRequest struct {
	pagination.Request
	Search *string `form:"search"`
}

type listOrganizationLocationsResponse struct {
	ID                  uuid.UUID                          `json:"id"`
	Name                string                             `json:"name"`
	Street              string                             `json:"street"`
	HouseNumber         string                             `json:"house_number"`
	HouseNumberAddition *string                            `json:"house_number_addition"`
	PostalCode          string                             `json:"postal_code"`
	City                string                             `json:"city"`
	Capacity            *int32                             `json:"capacity"`
	Occupied            int32                              `json:"occupied"`
	Available           int32                              `json:"available"`
	CreatedAt           time.Time                          `json:"created_at"`
	UpdatedAt           time.Time                          `json:"updated_at"`
	Shifts              []listOrganizationLocationShiftDTO `json:"shifts"`
}

type listOrganizationLocationShiftDTO struct {
	ID         uuid.UUID `json:"id"`
	LocationID uuid.UUID `json:"location_id"`
	Slot       int16     `json:"slot"`
	ShiftName  string    `json:"shift"`
	StartTime  string    `json:"start_time"`
	EndTime    string    `json:"end_time"`
}

func toListOrganizationsParams(req listOrganizationsRequest) domain.ListOrganizationsParams {
	search := ""
	if req.Search != nil {
		search = *req.Search
	}

	return domain.ListOrganizationsParams{
		Limit:  req.PageSize,
		Offset: (req.Page - 1) * req.PageSize,
		Search: search,
	}
}

func toListOrganizationsResponse(organization domain.Organization) listOrganizationsResponse {
	return listOrganizationsResponse{
		ID:                  organization.ID,
		Name:                organization.Name,
		Street:              organization.Street,
		HouseNumber:         organization.HouseNumber,
		HouseNumberAddition: organization.HouseNumberAddition,
		PostalCode:          organization.PostalCode,
		City:                organization.City,
		Email:               organization.Email,
		KvkNumber:           organization.KvkNumber,
		BtwNumber:           organization.BtwNumber,
		LocationCount:       organization.LocationCount,
	}
}

func toListOrganizationLocationsParams(organizationID uuid.UUID, req listOrganizationLocationsRequest) domain.ListOrganizationLocationsParams {
	search := ""
	if req.Search != nil {
		search = *req.Search
	}

	return domain.ListOrganizationLocationsParams{
		OrganizationID: organizationID,
		Limit:          req.PageSize,
		Offset:         (req.Page - 1) * req.PageSize,
		Search:         search,
	}
}

func toListOrganizationLocationsResponse(location domain.OrganizationLocation) listOrganizationLocationsResponse {
	return listOrganizationLocationsResponse{
		ID:                  location.ID,
		Name:                location.Name,
		Street:              location.Street,
		HouseNumber:         location.HouseNumber,
		HouseNumberAddition: location.HouseNumberAddition,
		PostalCode:          location.PostalCode,
		City:                location.City,
		Capacity:            location.Capacity,
		Occupied:            location.Occupied,
		Available:           location.Available,
		CreatedAt:           location.CreatedAt,
		UpdatedAt:           location.UpdatedAt,
		Shifts:              toListOrganizationLocationShiftDTOs(location.Shifts),
	}
}

func toListOrganizationLocationShiftDTOs(shifts []domain.OrganizationLocationShift) []listOrganizationLocationShiftDTO {
	result := make([]listOrganizationLocationShiftDTO, len(shifts))
	for i, shift := range shifts {
		result[i] = listOrganizationLocationShiftDTO{
			ID:         shift.ID,
			LocationID: shift.LocationID,
			Slot:       shift.Slot,
			ShiftName:  shift.ShiftName,
			StartTime:  shift.StartTime,
			EndTime:    shift.EndTime,
		}
	}
	return result
}
