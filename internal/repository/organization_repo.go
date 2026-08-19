package repository

import (
	"context"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/domain"
	"maicare_go/pkg/conv"

	"github.com/google/uuid"
)

type OrganizationRepository struct {
	store *db.Store
}

func NewOrganizationRepository(store *db.Store) domain.OrganizationRepository {
	return &OrganizationRepository{store: store}
}

func (r *OrganizationRepository) CreateOrganization(ctx context.Context, params domain.CreateOrganizationParams) (*domain.Organization, error) {
	organization, err := r.store.CreateOrganisation(ctx, db.CreateOrganisationParams{
		Name:                params.Name,
		Street:              params.Street,
		HouseNumber:         params.HouseNumber,
		HouseNumberAddition: params.HouseNumberAddition,
		PostalCode:          params.PostalCode,
		City:                params.City,
		PhoneNumber:         nil,
		Email:               params.Email,
		KvkNumber:           params.KvkNumber,
		BtwNumber:           params.BtwNumber,
	})
	if err != nil {
		return nil, err
	}

	return toDomainOrganizationFromOrganisation(organization), nil
}

func (r *OrganizationRepository) UpdateOrganization(ctx context.Context, organizationID uuid.UUID, params domain.UpdateOrganizationParams) (*domain.Organization, error) {
	organization, err := r.store.UpdateOrganisation(ctx, db.UpdateOrganisationParams{
		ID:                  organizationID,
		Name:                params.Name,
		Street:              params.Street,
		HouseNumber:         params.HouseNumber,
		HouseNumberAddition: params.HouseNumberAddition,
		PostalCode:          params.PostalCode,
		City:                params.City,
		PhoneNumber:         nil,
		Email:               params.Email,
		KvkNumber:           params.KvkNumber,
		BtwNumber:           params.BtwNumber,
	})
	if err != nil {
		return nil, err
	}

	return toDomainOrganizationFromOrganisation(organization), nil
}

func (r *OrganizationRepository) DeleteOrganization(ctx context.Context, organizationID uuid.UUID) error {
	_, err := r.store.DeleteOrganisation(ctx, organizationID)
	return err
}

func (r *OrganizationRepository) CreateOrganizationLocation(ctx context.Context, organizationID uuid.UUID, params domain.CreateOrganizationLocationParams) (*domain.OrganizationLocation, error) {
	location, err := r.store.CreateLocation(ctx, db.CreateLocationParams{
		OrganisationID:      organizationID,
		Name:                params.Name,
		Street:              params.Street,
		HouseNumber:         params.HouseNumber,
		HouseNumberAddition: params.HouseNumberAddition,
		PostalCode:          params.PostalCode,
		City:                params.City,
		Capacity:            params.Capacity,
	})
	if err != nil {
		return nil, err
	}

	result := toDomainOrganizationLocationFromLocation(location)
	return &result, nil
}

func (r *OrganizationRepository) GetLocationByID(ctx context.Context, locationID uuid.UUID) (*domain.OrganizationLocation, error) {
	location, err := r.store.GetLocation(ctx, locationID)
	if err != nil {
		return nil, err
	}

	result := toDomainOrganizationLocationFromLocation(location)
	return &result, nil
}

func (r *OrganizationRepository) UpdateLocation(ctx context.Context, locationID uuid.UUID, params domain.UpdateOrganizationLocationParams) (*domain.OrganizationLocation, error) {
	location, err := r.store.UpdateLocation(ctx, db.UpdateLocationParams{
		ID:                  locationID,
		Name:                params.Name,
		Street:              params.Street,
		HouseNumber:         params.HouseNumber,
		HouseNumberAddition: params.HouseNumberAddition,
		PostalCode:          params.PostalCode,
		City:                params.City,
		Capacity:            params.Capacity,
	})
	if err != nil {
		return nil, err
	}

	result := toDomainOrganizationLocationFromLocation(location)
	return &result, nil
}

func (r *OrganizationRepository) DeleteLocation(ctx context.Context, locationID uuid.UUID) error {
	_, err := r.store.DeleteLocation(ctx, locationID)
	return err
}

func (r *OrganizationRepository) CreateShift(ctx context.Context, params domain.CreateShiftParams) (*domain.OrganizationLocationShift, error) {
	startTime, err := conv.PgTimeFromString(params.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := conv.PgTimeFromString(params.EndTime)
	if err != nil {
		return nil, err
	}

	shift, err := r.store.CreateShift(ctx, db.CreateShiftParams{
		LocationID: params.LocationID,
		Slot:       params.Slot,
		ShiftName:  params.ShiftName,
		StartTime:  startTime,
		EndTime:    endTime,
	})
	if err != nil {
		return nil, err
	}

	result := domain.OrganizationLocationShift{
		ID:         shift.ID,
		LocationID: shift.LocationID,
		Slot:       shift.Slot,
		ShiftName:  shift.ShiftName,
		StartTime:  conv.StringFromPgTime(shift.StartTime),
		EndTime:    conv.StringFromPgTime(shift.EndTime),
	}

	return &result, nil
}

func (r *OrganizationRepository) UpdateShift(ctx context.Context, shiftID uuid.UUID, params domain.UpdateShiftParams) (*domain.OrganizationLocationShift, error) {
	startTime, err := conv.PgTimeFromString(params.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := conv.PgTimeFromString(params.EndTime)
	if err != nil {
		return nil, err
	}

	shift, err := r.store.UpdateShift(ctx, db.UpdateShiftParams{
		ID:        shiftID,
		ShiftName: params.ShiftName,
		StartTime: startTime,
		EndTime:   endTime,
	})
	if err != nil {
		return nil, err
	}

	result := domain.OrganizationLocationShift{
		ID:         shift.ID,
		LocationID: shift.LocationID,
		Slot:       shift.Slot,
		ShiftName:  shift.ShiftName,
		StartTime:  conv.StringFromPgTime(shift.StartTime),
		EndTime:    conv.StringFromPgTime(shift.EndTime),
	}

	return &result, nil
}

func (r *OrganizationRepository) DeleteShift(ctx context.Context, shiftID uuid.UUID) error {
	return r.store.DeleteShift(ctx, shiftID)
}

func (r *OrganizationRepository) GetShiftsByLocationID(ctx context.Context, locationID uuid.UUID) ([]domain.OrganizationLocationShift, error) {
	shifts, err := r.store.GetShiftsByLocationID(ctx, locationID)
	if err != nil {
		return nil, err
	}

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

	return result, nil
}

func (r *OrganizationRepository) GetOrganizationCounts(ctx context.Context, organizationID uuid.UUID) (*domain.OrganizationCounts, error) {
	counts, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.GetOrganisationCountsRow, error) {
		return q.GetOrganisationCounts(ctx, organizationID)
	})
	if err != nil {
		return nil, err
	}

	return &domain.OrganizationCounts{
		OrganizationID:   counts.OrganisationID,
		OrganizationName: counts.OrganisationName,
		LocationCount:    counts.LocationCount,
		ClientCount:      counts.ClientCount,
		EmployeeCount:    counts.EmployeeCount,
	}, nil
}

func (r *OrganizationRepository) GetGlobalOrganizationCounts(ctx context.Context) (*domain.GlobalOrganizationCounts, error) {
	counts, err := actorQuery(ctx, r.store, func(q *db.Queries) (db.GetGlobalOrganisationCountsRow, error) {
		return q.GetGlobalOrganisationCounts(ctx)
	})
	if err != nil {
		return nil, err
	}

	return &domain.GlobalOrganizationCounts{
		TotalLocations: counts.TotalLocations,
		TotalCapacity:  counts.TotalCapacity,
	}, nil
}

func (r *OrganizationRepository) GetOrganizationByID(ctx context.Context, organizationID uuid.UUID) (*domain.Organization, error) {
	organization, err := r.store.GetOrganisation(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	return toDomainOrganizationDetail(organization), nil
}

func (r *OrganizationRepository) ListOrganizations(ctx context.Context, params domain.ListOrganizationsParams) (*domain.OrganizationPage, error) {
	rows, err := r.store.ListOrganisationsPaginated(ctx, db.ListOrganisationsPaginatedParams{
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
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListLocationsPaginatedRow, error) {
		return q.ListLocationsPaginated(ctx, db.ListLocationsPaginatedParams{
			OrganisationID: params.OrganizationID,
			Limit:          params.Limit,
			Offset:         params.Offset,
			Column4:        params.Search,
		})
	})
	if err != nil {
		return nil, err
	}

	page := &domain.OrganizationLocationPage{
		Items: make([]domain.OrganizationLocation, 0, len(rows)),
	}

	for _, row := range rows {
		shifts, err := r.store.GetShiftsByLocationID(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		page.Items = append(page.Items, toDomainOrganizationLocation(row, shifts))
		page.TotalCount = row.TotalCount
	}

	return page, nil
}

func (r *OrganizationRepository) ListAllLocations(ctx context.Context, params domain.ListAllLocationsParams) (*domain.OrganizationLocationPage, error) {
	rows, err := actorQuery(ctx, r.store, func(q *db.Queries) ([]db.ListAllLocationsPaginatedRow, error) {
		return q.ListAllLocationsPaginated(ctx, db.ListAllLocationsPaginatedParams{
			Limit:   params.Limit,
			Offset:  params.Offset,
			Column3: params.Search,
		})
	})
	if err != nil {
		return nil, err
	}

	page := &domain.OrganizationLocationPage{
		Items: make([]domain.OrganizationLocation, 0, len(rows)),
	}

	for _, row := range rows {
		shifts, err := r.store.GetShiftsByLocationID(ctx, row.ID)
		if err != nil {
			return nil, err
		}

		page.Items = append(page.Items, toDomainOrganizationLocationFromAllLocations(row, shifts))
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

func toDomainOrganizationFromOrganisation(organization db.Organisation) *domain.Organization {
	return &domain.Organization{
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
		CreatedAt:           conv.TimeFromPgTimestamptz(organization.CreatedAt),
		UpdatedAt:           conv.TimeFromPgTimestamptz(organization.UpdatedAt),
	}
}

func toDomainOrganizationDetail(row db.GetOrganisationRow) *domain.Organization {
	return &domain.Organization{
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
		CreatedAt:           conv.TimeFromPgTimestamptz(row.CreatedAt),
		UpdatedAt:           conv.TimeFromPgTimestamptz(row.UpdatedAt),
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
		CreatedAt:           conv.TimeFromPgTimestamptz(row.CreatedAt),
		UpdatedAt:           conv.TimeFromPgTimestamptz(row.UpdatedAt),
		Shifts:              toDomainOrganizationLocationShifts(shifts),
	}
}

func toDomainOrganizationLocationFromLocation(row db.Location) domain.OrganizationLocation {
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
		CreatedAt:           conv.TimeFromPgTimestamptz(row.CreatedAt),
		UpdatedAt:           conv.TimeFromPgTimestamptz(row.UpdatedAt),
	}
}

func toDomainOrganizationLocationFromAllLocations(row db.ListAllLocationsPaginatedRow, shifts []db.LocationShift) domain.OrganizationLocation {
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

var _ domain.OrganizationRepository = (*OrganizationRepository)(nil)
