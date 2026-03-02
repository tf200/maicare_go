package organization

import (
	"context"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"
	"maicare_go/pagination"
	"maicare_go/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *organizationService) ListOrgLocations(ctx *gin.Context, organizationID uuid.UUID, req ListLocationsRequest) (*pagination.Response[ListOrgLocationsResponse], error) {
	search := ""
	if req.Search != nil {
		search = *req.Search
	}

	params := db.ListLocationsPaginatedParams{
		OrganisationID: organizationID,
		Limit:          req.PageSize,
		Offset:         (req.Page - 1) * req.PageSize,
		Column4:        search,
	}

	locations, err := s.Store.ListLocationsPaginated(ctx, params)

	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListOrgLocationsApi", "Failed to list locations", zap.Error(err))
		return nil, fmt.Errorf("failed to list locations")
	}

	response := []ListOrgLocationsResponse{}
	var totalCount int64
	for _, loc := range locations {
		shifts, err := s.Store.GetShiftsByLocationID(ctx, loc.ID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListOrgLocationsApi", "Failed to get shifts by location", zap.Error(err), zap.String("location_id", loc.ID.String()))
			return nil, fmt.Errorf("failed to list locations")
		}

		locationShifts := make([]ListShiftsByLocationIDResponse, len(shifts))
		for i, shift := range shifts {
			locationShifts[i] = ListShiftsByLocationIDResponse{
				ID:         shift.ID,
				LocationID: shift.LocationID,
				Slot:       shift.Slot,
				ShiftName:  shift.ShiftName,
				StartTime:  util.PgTimeToString(shift.StartTime),
				EndTime:    util.PgTimeToString(shift.EndTime),
			}
		}

		totalCount = loc.TotalCount
		response = append(response, ListOrgLocationsResponse{
			ID:                  loc.ID,
			Name:                loc.Name,
			Street:              loc.Street,
			HouseNumber:         loc.HouseNumber,
			HouseNumberAddition: loc.HouseNumberAddition,
			PostalCode:          loc.PostalCode,
			City:                loc.City,
			Capacity:            loc.Capacity,
			Occupied:            int32(loc.ClientCount),
			Available: func() int32 {
				if loc.Capacity != nil {
					available := *loc.Capacity - int32(loc.ClientCount)
					if available < 0 {
						return 0
					}
					return available
				}
				return 0
			}(),
			CreatedAt: loc.CreatedAt.Time,
			UpdatedAt: loc.UpdatedAt.Time,
			Shifts:    locationShifts,
		})
	}

	paginatedResponse := pagination.NewResponse(ctx, req.Request, response, totalCount)
	return &paginatedResponse, nil
}

func (s *organizationService) ListAllLocations(ctx *gin.Context, req ListAllLocationsRequest) (*pagination.Response[ListLocationsResponse], error) {
	search := ""
	if req.Search != nil {
		search = *req.Search
	}

	params := req.GetParams()
	locations, err := s.Store.ListAllLocationsPaginated(ctx, db.ListAllLocationsPaginatedParams{
		Limit:   params.Limit,
		Offset:  params.Offset,
		Column3: search,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListAllLocationsApi", "Failed to list all locations", zap.Error(err))
		return nil, fmt.Errorf("failed to list all locations")
	}

	if len(locations) == 0 {
		paginatedResponse := pagination.NewResponse(ctx, req.Request, []ListLocationsResponse{}, 0)
		return &paginatedResponse, nil
	}

	response := make([]ListLocationsResponse, len(locations))
	for i, loc := range locations {
		shifts, err := s.Store.GetShiftsByLocationID(ctx, loc.ID)
		if err != nil {
			s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "ListAllLocationsApi", "Failed to get shifts by location", zap.Error(err), zap.String("location_id", loc.ID.String()))
			return nil, fmt.Errorf("failed to list all locations")
		}

		locationShifts := make([]ListShiftsByLocationIDResponse, len(shifts))
		for j, shift := range shifts {
			locationShifts[j] = ListShiftsByLocationIDResponse{
				ID:         shift.ID,
				LocationID: shift.LocationID,
				Slot:       shift.Slot,
				ShiftName:  shift.ShiftName,
				StartTime:  util.PgTimeToString(shift.StartTime),
				EndTime:    util.PgTimeToString(shift.EndTime),
			}
		}

		response[i] = ListLocationsResponse{
			ID:                  loc.ID,
			Name:                loc.Name,
			Street:              loc.Street,
			HouseNumber:         loc.HouseNumber,
			HouseNumberAddition: loc.HouseNumberAddition,
			PostalCode:          loc.PostalCode,
			City:                loc.City,
			Capacity:            loc.Capacity,
			Occupied:            int32(loc.ClientCount),
			Available: func() int32 {
				if loc.Capacity != nil {
					available := *loc.Capacity - int32(loc.ClientCount)
					if available < 0 {
						return 0
					}
					return available
				}
				return 0
			}(),
			CreatedAt: loc.CreatedAt.Time,
			UpdatedAt: loc.UpdatedAt.Time,
			Shifts:    locationShifts,
		}
	}

	paginatedResponse := pagination.NewResponse(ctx, req.Request, response, locations[0].TotalCount)
	return &paginatedResponse, nil
}

func (s *organizationService) CreateLocation(ctx context.Context, organizationID uuid.UUID, req CreateLocationRequest) (*CreateLocationResponse, error) {
	location, err := s.Store.CreateLocation(ctx, db.CreateLocationParams{
		OrganisationID:      organizationID,
		Name:                req.Name,
		Street:              req.Street,
		HouseNumber:         req.HouseNumber,
		HouseNumberAddition: req.HouseNumberAddition,
		PostalCode:          req.PostalCode,
		City:                req.City,
		Capacity:            req.Capacity,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "CreateLocationApi", "Failed to create location", zap.Error(err))
		return nil, fmt.Errorf("failed to create location")
	}

	return &CreateLocationResponse{
		ID:                  location.ID,
		Name:                location.Name,
		Street:              location.Street,
		HouseNumber:         location.HouseNumber,
		HouseNumberAddition: location.HouseNumberAddition,
		PostalCode:          location.PostalCode,
		City:                location.City,
		Capacity:            location.Capacity,
	}, nil
}

func (s *organizationService) UpdateLocation(ctx context.Context, locationID uuid.UUID, req UpdateLocationRequest) (*UpdateLocationResponse, error) {
	location, err := s.Store.UpdateLocation(ctx, db.UpdateLocationParams{
		ID:                  locationID,
		Name:                req.Name,
		Street:              req.Street,
		HouseNumber:         req.HouseNumber,
		HouseNumberAddition: req.HouseNumberAddition,
		PostalCode:          req.PostalCode,
		City:                req.City,
		Capacity:            req.Capacity,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "UpdateLocationApi", "Failed to update location", zap.Error(err))
		return nil, fmt.Errorf("failed to update location")
	}

	return &UpdateLocationResponse{
		ID:                  location.ID,
		Name:                location.Name,
		Street:              location.Street,
		HouseNumber:         location.HouseNumber,
		HouseNumberAddition: location.HouseNumberAddition,
		PostalCode:          location.PostalCode,
		City:                location.City,
		Capacity:            location.Capacity,
	}, nil
}

func (s *organizationService) DeleteLocation(ctx context.Context, locationID uuid.UUID) (*DeleteLocationResponse, error) {
	_, err := s.Store.DeleteLocation(ctx, locationID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "DeleteLocationApi", "Failed to delete location", zap.Error(err))
		return nil, fmt.Errorf("failed to delete location")
	}

	return &DeleteLocationResponse{
		ID: locationID,
	}, nil
}

func (s *organizationService) GetLocationByID(ctx context.Context, locationID uuid.UUID) (*GetLocationResponse, error) {
	location, err := s.Store.GetLocation(ctx, locationID)
	if err != nil {
		s.Logger.LogBusinessEvent(ctx, logger.LogLevelError, "GetLocationApi", "Failed to get location by ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get location by ID")
	}

	return &GetLocationResponse{
		ID:                  location.ID,
		Name:                location.Name,
		Street:              location.Street,
		HouseNumber:         location.HouseNumber,
		HouseNumberAddition: location.HouseNumberAddition,
		PostalCode:          location.PostalCode,
		City:                location.City,
		Capacity:            location.Capacity,
	}, nil
}
