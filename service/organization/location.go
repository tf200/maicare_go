package organization

import (
	"context"
	"fmt"

	db "maicare_go/db/sqlc"
	"maicare_go/logger"

	"go.uber.org/zap"
)

func (s *organizationService) ListOrgLocations(ctx context.Context, organizationID int64) ([]ListLocationsResponse, error) {
	locations, err := s.Store.ListLocations(ctx, organizationID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListOrgLocationsApi", "Failed to list locations", zap.Error(err))
		return nil, fmt.Errorf("failed to list locations")
	}

	response := []ListLocationsResponse{}
	for _, loc := range locations {
		response = append(response, ListLocationsResponse{
			ID:        loc.ID,
			Name:      loc.Name,
			Address:   loc.Address,
			Capacity:  loc.Capacity,
			CreatedAt: loc.CreatedAt.Time,
			UpdatedAt: loc.UpdatedAt.Time,
		})
	}

	return response, nil
}

func (s *organizationService) ListAllLocations(ctx context.Context) ([]ListLocationsResponse, error) {
	locations, err := s.Store.ListAllLocations(ctx)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "ListAllLocationsApi", "Failed to list all locations", zap.Error(err))
		return nil, fmt.Errorf("failed to list all locations")
	}

	response := []ListLocationsResponse{}
	for _, loc := range locations {
		response = append(response, ListLocationsResponse{
			ID:       loc.ID,
			Name:     loc.Name,
			Address:  loc.Address,
			Capacity: loc.Capacity,
			Occupied: int32(loc.ClientCount),
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
		})
	}

	return response, nil
}

func (s *organizationService) CreateLocation(ctx context.Context, organizationID int64, req CreateLocationRequest) (*CreateLocationResponse, error) {
	location, err := s.Store.CreateLocation(ctx, db.CreateLocationParams{
		OrganisationID: organizationID,
		Name:           req.Name,
		Address:        req.Address,
		Capacity:       req.Capacity,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "CreateLocationApi", "Failed to create location", zap.Error(err))
		return nil, fmt.Errorf("failed to create location")
	}

	return &CreateLocationResponse{
		ID:       location.ID,
		Name:     location.Name,
		Address:  location.Address,
		Capacity: location.Capacity,
	}, nil
}

func (s *organizationService) UpdateLocation(ctx context.Context, locationID int64, req UpdateLocationRequest) (*UpdateLocationResponse, error) {
	location, err := s.Store.UpdateLocation(ctx, db.UpdateLocationParams{
		Name:     req.Name,
		Address:  req.Address,
		Capacity: req.Capacity,
		ID:       locationID,
	})
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "UpdateLocationApi", "Failed to update location", zap.Error(err))
		return nil, fmt.Errorf("failed to update location")
	}

	return &UpdateLocationResponse{
		ID:       location.ID,
		Name:     location.Name,
		Address:  location.Address,
		Capacity: location.Capacity,
	}, nil
}

func (s *organizationService) DeleteLocation(ctx context.Context, locationID int64) (*DeleteLocationResponse, error) {
	_, err := s.Store.DeleteLocation(ctx, locationID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "DeleteLocationApi", "Failed to delete location", zap.Error(err))
		return nil, fmt.Errorf("failed to delete location")
	}

	return &DeleteLocationResponse{
		ID: locationID,
	}, nil
}

func (s *organizationService) GetLocationByID(ctx context.Context, locationID int64) (*GetLocationResponse, error) {
	location, err := s.Store.GetLocation(ctx, locationID)
	if err != nil {
		s.Logger.LogBusinessEvent(logger.LogLevelError, "GetLocationApi", "Failed to get location by ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get location by ID")
	}

	return &GetLocationResponse{
		ID:       location.ID,
		Name:     location.Name,
		Address:  location.Address,
		Capacity: location.Capacity,
	}, nil
}
