package repository

import (
	"context"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
	"maicare_go/pkg/conv"

	"github.com/google/uuid"
)

type organizationLocationQuerier interface {
	ListOrganisationsPaginated(ctx context.Context, arg db.ListOrganisationsPaginatedParams) ([]db.ListOrganisationsPaginatedRow, error)
	ListLocationsPaginated(ctx context.Context, arg db.ListLocationsPaginatedParams) ([]db.ListLocationsPaginatedRow, error)
	GetShiftsByLocationID(ctx context.Context, locationID uuid.UUID) ([]db.LocationShift, error)
}

type OrganizationRepository struct {
	queries organizationLocationQuerier
}

func NewOrganizationRepository(queries organizationLocationQuerier) *OrganizationRepository {
	return &OrganizationRepository{queries: queries}
}

func (r *OrganizationRepository) ListOrganizations(ctx context.Context, params domain.ListOrganizationsParams) (*domain.OrganizationPage, error) {
	rows, err := r.queries.ListOrganisationsPaginated(ctx, db.ListOrganisationsPaginatedParams{
		Limit:   params.Limit,
		Offset:  params.Offset,
		Column3: params.Search,
	})
	if err != nil {
		return nil, err
	}

	page := &domain.OrganizationPage{
		Items: make([]domain.Organization, 0, len(rows)),
	}

	if len(rows) > 0 {
		page.TotalCount = rows[0].TotalCount
	}

	for _, row := range rows {
		page.Items = append(page.Items, toDomainOrganization(row))
	}

	return page, nil
}

func (r *OrganizationRepository) ListOrganizationLocations(ctx context.Context, params domain.ListOrganizationLocationsParams) (*domain.OrganizationLocationPage, error) {
	rows, err := r.queries.ListLocationsPaginated(ctx, db.ListLocationsPaginatedParams{
		OrganisationID: params.OrganizationID,
		Limit:          params.Limit,
		Offset:         params.Offset,
		Column4:        params.Search,
	})
	if err != nil {
		return nil, err
	}

	page := &domain.OrganizationLocationPage{
		Items: make([]domain.OrganizationLocation, 0, len(rows)),
	}

	for _, row := range rows {
		shifts, err := r.queries.GetShiftsByLocationID(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		page.Items = append(page.Items, toDomainOrganizationLocation(row, shifts))
		page.TotalCount = row.TotalCount
	}

	return page, nil
}

func toDomainOrganization(row db.ListOrganisationsPaginatedRow) domain.Organization {
	return domain.Organization{
		ID:                  row.ID,
		Name:                row.Name,
		Street:              row.Street,
		HouseNumber:         row.HouseNumber,
		HouseNumberAddition: row.HouseNumberAddition,
		PostalCode:          row.PostalCode,
		City:                row.City,
		Email:               row.Email,
		KvkNumber:           row.KvkNumber,
		BtwNumber:           row.BtwNumber,
		LocationCount:       row.LocationCount,
	}
}

func toDomainOrganizationLocation(row db.ListLocationsPaginatedRow, shifts []db.LocationShift) domain.OrganizationLocation {
	return domain.OrganizationLocation{
		ID:                  row.ID,
		OrganizationID:      row.OrganisationID,
		Name:                row.Name,
		Street:              row.Street,
		HouseNumber:         row.HouseNumber,
		HouseNumberAddition: row.HouseNumberAddition,
		PostalCode:          row.PostalCode,
		City:                row.City,
		Capacity:            row.Capacity,
		Occupied:            int32(row.ClientCount),
		Available:           availableCapacity(row.Capacity, row.ClientCount),
		CreatedAt:           conv.TimeFromPgTimestamptz(row.CreatedAt),
		UpdatedAt:           conv.TimeFromPgTimestamptz(row.UpdatedAt),
		Shifts:              toDomainOrganizationLocationShifts(shifts),
	}
}

func toDomainOrganizationLocationShifts(shifts []db.LocationShift) []domain.OrganizationLocationShift {
	result := make([]domain.OrganizationLocationShift, len(shifts))
	for i, shift := range shifts {
		result[i] = domain.OrganizationLocationShift{
			ID:         shift.ID,
			LocationID: shift.LocationID,
			Slot:       shift.Slot,
			ShiftName:  shift.ShiftName,
			StartTime:  conv.StringFromPgTime(shift.StartTime),
			EndTime:    conv.StringFromPgTime(shift.EndTime),
		}
	}
	return result
}

func availableCapacity(capacity *int32, occupied int64) int32 {
	if capacity == nil {
		return 0
	}

	available := *capacity - int32(occupied)
	if available < 0 {
		return 0
	}

	return available
}
