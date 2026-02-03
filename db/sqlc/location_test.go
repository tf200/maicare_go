package db

import (
	"context"
	"testing"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateOrganisation(t *testing.T) {
	tests := []struct {
		name   string
		params CreateOrganisationParams
		checks func(t *testing.T, organisation Organisation)
	}{
		{
			name: "successful creation with minimal fields",
			params: CreateOrganisationParams{
				Name:        "Test Organisation",
				Street:      "Test Street",
				HouseNumber: "123",
				PostalCode:  "12345",
				City:        "Test City",
			},
			checks: func(t *testing.T, organisation Organisation) {
				require.Equal(t, "Test Organisation", organisation.Name)
				require.Equal(t, "Test Street", organisation.Street)
				require.Equal(t, "123", organisation.HouseNumber)
				require.Nil(t, organisation.HouseNumberAddition)
				require.Equal(t, "12345", organisation.PostalCode)
				require.Equal(t, "Test City", organisation.City)
				require.Nil(t, organisation.PhoneNumber)
				require.Nil(t, organisation.Email)
				require.Nil(t, organisation.KvkNumber)
				require.Nil(t, organisation.BtwNumber)
			},
		},
		{
			name: "successful creation with all fields",
			params: CreateOrganisationParams{
				Name:                "Full Organisation",
				Street:              "Full Ave",
				HouseNumber:         "456",
				HouseNumberAddition: util.StringPtr("A"),
				PostalCode:          "67890",
				City:                "Full City",
				PhoneNumber:         util.StringPtr("+1234567890"),
				Email:               util.StringPtr("full@example.com"),
				KvkNumber:           util.StringPtr("KVK123456"),
				BtwNumber:           util.StringPtr("BTW789012"),
			},
			checks: func(t *testing.T, organisation Organisation) {
				require.Equal(t, "Full Organisation", organisation.Name)
				require.Equal(t, "Full Ave", organisation.Street)
				require.Equal(t, "456", organisation.HouseNumber)
				require.NotNil(t, organisation.HouseNumberAddition)
				require.Equal(t, "A", *organisation.HouseNumberAddition)
				require.Equal(t, "67890", organisation.PostalCode)
				require.Equal(t, "Full City", organisation.City)
				require.NotNil(t, organisation.PhoneNumber)
				require.Equal(t, "+1234567890", *organisation.PhoneNumber)
				require.NotNil(t, organisation.Email)
				require.Equal(t, "full@example.com", *organisation.Email)
				require.NotNil(t, organisation.KvkNumber)
				require.Equal(t, "KVK123456", *organisation.KvkNumber)
				require.NotNil(t, organisation.BtwNumber)
				require.Equal(t, "BTW789012", *organisation.BtwNumber)
			},
		},
		{
			name: "successful creation with some optional fields",
			params: CreateOrganisationParams{
				Name:        "Partial Organisation",
				Street:      "Partial St",
				HouseNumber: "789",
				PostalCode:  "13579",
				City:        "Partial City",
				Email:       util.StringPtr("partial@example.com"),
				KvkNumber:   util.StringPtr("KVK111111"),
			},
			checks: func(t *testing.T, organisation Organisation) {
				require.Equal(t, "Partial Organisation", organisation.Name)
				require.Equal(t, "Partial St", organisation.Street)
				require.Equal(t, "789", organisation.HouseNumber)
				require.Nil(t, organisation.HouseNumberAddition)
				require.NotNil(t, organisation.Email)
				require.Equal(t, "partial@example.com", *organisation.Email)
				require.NotNil(t, organisation.KvkNumber)
				require.Equal(t, "KVK111111", *organisation.KvkNumber)
				require.Nil(t, organisation.PhoneNumber)
				require.Nil(t, organisation.BtwNumber)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			organisation, err := qtx.CreateOrganisation(ctx, tt.params)
			require.NoError(t, err, "CreateOrganisation() should not error")

			tt.checks(t, organisation)
		})
	}
}

func TestListOrganisations(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries)
		checks func(t *testing.T, organisations []ListOrganisationsRow, err error)
	}{
		{
			name: "list organisations with multiple entries",
			setup: func(ctx context.Context, qtx *Queries) {
				// Create multiple organisations
				for i := 0; i < 3; i++ {
					_, err := qtx.CreateOrganisation(ctx, CreateOrganisationParams{
						Name:        util.RandomString(10),
						Street:      util.RandomString(20),
						HouseNumber: util.RandomString(4),
						PostalCode:  util.RandomString(5),
						City:        util.RandomString(10),
					})
					require.NoError(t, err)
				}
			},
			checks: func(t *testing.T, organisations []ListOrganisationsRow, err error) {
				require.NoError(t, err, "ListOrganisations() should not error")
				require.GreaterOrEqual(t, len(organisations), 3, "should have at least 3 organisations")
				// Check that they are ordered by name
				for i := 1; i < len(organisations); i++ {
					require.LessOrEqual(t, organisations[i-1].Name, organisations[i].Name, "organisations should be ordered by name")
				}
			},
		},
		{
			name: "list organisations when empty",
			setup: func(ctx context.Context, qtx *Queries) {
				// No setup - keep it empty
			},
			checks: func(t *testing.T, organisations []ListOrganisationsRow, err error) {
				require.NoError(t, err, "ListOrganisations() should not error")
				require.Empty(t, organisations, "should return empty slice when no organisations exist")
			},
		},
		{
			name: "list organisations with location counts",
			setup: func(ctx context.Context, qtx *Queries) {
				org, err := qtx.CreateOrganisation(ctx, CreateOrganisationParams{
					Name:        "Org With Locations",
					Street:      "Test St",
					HouseNumber: "123",
					PostalCode:  "12345",
					City:        "Test City",
				})
				require.NoError(t, err)

				// Create locations for the organisation
				for i := 0; i < 2; i++ {
					_, err := qtx.CreateLocation(ctx, CreateLocationParams{
						OrganisationID: org.ID,
						Name:           util.RandomString(10),
						Street:         util.RandomString(20),
						HouseNumber:    util.RandomString(2),
						PostalCode:     util.RandomString(6),
						City:           util.RandomString(10),
					})
					require.NoError(t, err)
				}
			},
			checks: func(t *testing.T, organisations []ListOrganisationsRow, err error) {
				require.NoError(t, err, "ListOrganisations() should not error")
				var orgWithLocations *ListOrganisationsRow
				for _, org := range organisations {
					if org.Name == "Org With Locations" {
						orgWithLocations = &org
						break
					}
				}
				require.NotNil(t, orgWithLocations, "should find organisation with locations")
				require.Equal(t, int64(2), orgWithLocations.LocationCount, "should have correct location count")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			tt.setup(ctx, qtx)
			organisations, err := qtx.ListOrganisations(ctx)
			tt.checks(t, organisations, err)
		})
	}
}

func TestGetOrganisation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, organisation GetOrganisationRow, err error)
	}{
		{
			name: "get existing organisation by ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				organisation, err := qtx.CreateOrganisation(ctx, CreateOrganisationParams{
					Name:        "Get Test Organisation",
					Street:      "Get St",
					HouseNumber: "456",
					PostalCode:  "54321",
					City:        "Get City",
				})
				require.NoError(t, err)
				return organisation.ID
			},
			checks: func(t *testing.T, organisation GetOrganisationRow, err error) {
				require.NoError(t, err, "GetOrganisation() should not error")
				require.Equal(t, "Get Test Organisation", organisation.Name)
				require.Equal(t, int64(0), organisation.LocationCount, "new organisation should have 0 locations")
			},
		},
		{
			name: "get non-existent organisation by ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New() // Use a very high ID that likely doesn't exist
			},
			checks: func(t *testing.T, organisation GetOrganisationRow, err error) {
				require.Error(t, err, "GetOrganisation() should error for non-existent ID")
			},
		},
		{
			name: "get organisation with location count",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				org, err := qtx.CreateOrganisation(ctx, CreateOrganisationParams{
					Name:        "Org With Locations",
					Street:      "Count St",
					HouseNumber: "789",
					PostalCode:  "98765",
					City:        "Count City",
				})
				require.NoError(t, err)

				// Create a location for the organisation
				_, err = qtx.CreateLocation(ctx, CreateLocationParams{
					OrganisationID: org.ID,
					Name:           "Test Location",
					Street:         "123 Location St",
					HouseNumber:    "1",
					PostalCode:     "1234AB",
					City:           "Test City",
				})
				require.NoError(t, err)

				return org.ID
			},
			checks: func(t *testing.T, organisation GetOrganisationRow, err error) {
				require.NoError(t, err, "GetOrganisation() should not error")
				require.Equal(t, int64(1), organisation.LocationCount, "should have 1 location")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			organisation, err := qtx.GetOrganisation(ctx, id)
			tt.checks(t, organisation, err)
		})
	}
}

func TestGetOrganisationCounts(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, counts GetOrganisationCountsRow, err error)
	}{
		{
			name: "get counts for organisation with no entities",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				org, err := qtx.CreateOrganisation(ctx, CreateOrganisationParams{
					Name:        "Empty Org",
					Street:      "Empty St",
					HouseNumber: "123",
					PostalCode:  "12345",
					City:        "Empty City",
				})
				require.NoError(t, err)
				return org.ID
			},
			checks: func(t *testing.T, counts GetOrganisationCountsRow, err error) {
				require.NoError(t, err, "GetOrganisationCounts() should not error")
				require.Equal(t, int64(0), counts.LocationCount, "should have 0 locations")
				require.Equal(t, int64(0), counts.ClientCount, "should have 0 clients")
				require.Equal(t, int64(0), counts.EmployeeCount, "should have 0 employees")
			},
		},
		{
			name: "get counts for organisation with locations",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				org, err := qtx.CreateOrganisation(ctx, CreateOrganisationParams{
					Name:        "Org With Locations",
					Street:      "Location St",
					HouseNumber: "456",
					PostalCode:  "54321",
					City:        "Location City",
				})
				require.NoError(t, err)

				// Create 2 locations
				for i := 0; i < 2; i++ {
					_, err := qtx.CreateLocation(ctx, CreateLocationParams{
						OrganisationID: org.ID,
						Name:           util.RandomString(10),
						Street:         util.RandomString(20),
						HouseNumber:    util.RandomString(2),
						PostalCode:     util.RandomString(6),
						City:           util.RandomString(10),
					})
					require.NoError(t, err)
				}

				return org.ID
			},
			checks: func(t *testing.T, counts GetOrganisationCountsRow, err error) {
				require.NoError(t, err, "GetOrganisationCounts() should not error")
				require.Equal(t, int64(2), counts.LocationCount, "should have 2 locations")
				require.Equal(t, int64(0), counts.ClientCount, "should have 0 clients")
				require.Equal(t, int64(0), counts.EmployeeCount, "should have 0 employees")
			},
		},
		{
			name: "get counts for non-existent organisation",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New() // Use a very high ID that likely doesn't exist
			},
			checks: func(t *testing.T, counts GetOrganisationCountsRow, err error) {
				require.Error(t, err, "GetOrganisationCounts() should error for non-existent ID")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			counts, err := qtx.GetOrganisationCounts(ctx, id)
			tt.checks(t, counts, err)
		})
	}
}

func TestUpdateOrganisation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateOrganisationParams
		checks func(t *testing.T, err error)
	}{
		{
			name: "successful update with name",
			setup: func(ctx context.Context, qtx *Queries) UpdateOrganisationParams {
				org := createRandomOrganisation(ctx, qtx)
				return UpdateOrganisationParams{
					ID:   org.ID,
					Name: util.StringPtr("Updated Organisation Name"),
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "UpdateOrganisation() should not error")
			},
		},
		{
			name: "successful update with all fields",
			setup: func(ctx context.Context, qtx *Queries) UpdateOrganisationParams {
				org := createRandomOrganisation(ctx, qtx)
				return UpdateOrganisationParams{
					ID:                  org.ID,
					Name:                util.StringPtr("Fully Updated Org"),
					Street:              util.StringPtr("Updated Ave"),
					HouseNumber:         util.StringPtr("456"),
					HouseNumberAddition: util.StringPtr("B"),
					PostalCode:          util.StringPtr("65432"),
					City:                util.StringPtr("Updated City"),
					PhoneNumber:         util.StringPtr("+9876543210"),
					Email:               util.StringPtr("updated@example.com"),
					KvkNumber:           util.StringPtr("KVK999999"),
					BtwNumber:           util.StringPtr("BTW888888"),
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "UpdateOrganisation() should not error")
			},
		},
		{
			name: "update with nil values (no change)",
			setup: func(ctx context.Context, qtx *Queries) UpdateOrganisationParams {
				org := createRandomOrganisation(ctx, qtx)
				return UpdateOrganisationParams{
					ID: org.ID,
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "UpdateOrganisation() should not error even with no updates")
			},
		},
		{
			name: "update non-existent organisation",
			setup: func(ctx context.Context, qtx *Queries) UpdateOrganisationParams {
				return UpdateOrganisationParams{
					ID:   uuid.New(),
					Name: util.StringPtr("Non-existent Org"),
				}
			},
			checks: func(t *testing.T, err error) {
				require.Error(t, err, "UpdateOrganisation() should error for non-existent organisation")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			_, err = qtx.UpdateOrganisation(ctx, params)
			tt.checks(t, err)
		})
	}
}

func TestDeleteOrganisation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, organisation Organisation, err error)
	}{
		{
			name: "successful deletion",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				org := createRandomOrganisation(ctx, qtx)
				return org.ID
			},
			checks: func(t *testing.T, organisation Organisation, err error) {
				require.NoError(t, err, "DeleteOrganisation() should not error")
				require.NotEmpty(t, organisation.ID, "deleted organisation should have ID")
			},
		},
		{
			name: "delete non-existent organisation",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, organisation Organisation, err error) {
				require.Error(t, err, "DeleteOrganisation() should error for non-existent ID")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			organisation, err := qtx.DeleteOrganisation(ctx, id)
			tt.checks(t, organisation, err)
		})
	}
}

func TestCreateLocation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateLocationParams
		checks func(t *testing.T, location Location)
	}{
		{
			name: "successful creation without capacity",
			setup: func(ctx context.Context, qtx *Queries) CreateLocationParams {
				org := createRandomOrganisation(ctx, qtx)
				return CreateLocationParams{
					OrganisationID: org.ID,
					Name:           "Test Location",
					Street:         "123 Location St",
					HouseNumber:    "1",
					PostalCode:     "1234AB",
					City:           "Test City",
					Capacity:       nil,
				}
			},
			checks: func(t *testing.T, location Location) {
				require.Equal(t, "Test Location", location.Name)
				require.Equal(t, "123 Location St", location.Street)
				require.Equal(t, "1", location.HouseNumber)
				require.Equal(t, "1234AB", location.PostalCode)
				require.Equal(t, "Test City", location.City)
				require.Nil(t, location.Capacity, "capacity should be nil when not provided")
			},
		},
		{
			name: "successful creation with capacity",
			setup: func(ctx context.Context, qtx *Queries) CreateLocationParams {
				org := createRandomOrganisation(ctx, qtx)
				capacity := int32(100)
				return CreateLocationParams{
					OrganisationID: org.ID,
					Name:           "Location With Capacity",
					Street:         "456 Capacity Ave",
					HouseNumber:    "2",
					PostalCode:     "5678CD",
					City:           "Capacity City",
					Capacity:       &capacity,
				}
			},
			checks: func(t *testing.T, location Location) {
				require.Equal(t, "Location With Capacity", location.Name)
				require.Equal(t, "456 Capacity Ave", location.Street)
				require.NotNil(t, location.Capacity, "capacity should not be nil")
				require.Equal(t, int32(100), *location.Capacity, "capacity should match")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			location, err := qtx.CreateLocation(ctx, params)
			require.NoError(t, err, "CreateLocation() should not error")

			tt.checks(t, location)
		})
	}
}

func TestListLocations(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, locations []ListLocationsRow, err error)
	}{
		{
			name: "list locations for organisation with multiple locations",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				org := createRandomOrganisation(ctx, qtx)

				// Create multiple locations
				for i := 0; i < 3; i++ {
					_, err := qtx.CreateLocation(ctx, CreateLocationParams{
						OrganisationID: org.ID,
						Name:           util.RandomString(10),
						Street:         util.RandomString(20),
						HouseNumber:    util.RandomString(2),
						PostalCode:     util.RandomString(6),
						City:           util.RandomString(10),
					})
					require.NoError(t, err)
				}

				return org.ID
			},
			checks: func(t *testing.T, locations []ListLocationsRow, err error) {
				require.NoError(t, err, "ListLocations() should not error")
				require.Len(t, locations, 3, "should have exactly 3 locations")
			},
		},
		{
			name: "list locations for organisation with no locations",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				org := createRandomOrganisation(ctx, qtx)
				return org.ID
			},
			checks: func(t *testing.T, locations []ListLocationsRow, err error) {
				require.NoError(t, err, "ListLocations() should not error")
				require.Empty(t, locations, "should return empty slice when no locations exist")
			},
		},
		{
			name: "list locations for non-existent organisation",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, locations []ListLocationsRow, err error) {
				require.NoError(t, err, "ListLocations() should not error even for non-existent organisation")
				require.Empty(t, locations, "should return empty slice for non-existent organisation")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			orgID := tt.setup(ctx, qtx)
			locations, err := qtx.ListLocations(ctx, orgID)
			tt.checks(t, locations, err)
		})
	}
}

func TestGetLocation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, location Location, err error)
	}{
		{
			name: "get existing location by ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				location := createRandomLocation(ctx, qtx)
				return location.ID
			},
			checks: func(t *testing.T, location Location, err error) {
				require.NoError(t, err, "GetLocation() should not error")
				require.NotEmpty(t, location.ID, "location should have valid ID")
				require.NotEmpty(t, location.Name, "location should have name")
			},
		},
		{
			name: "get non-existent location by ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, location Location, err error) {
				require.Error(t, err, "GetLocation() should error for non-existent ID")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			location, err := qtx.GetLocation(ctx, id)
			tt.checks(t, location, err)
		})
	}
}

func TestUpdateLocation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateLocationParams
		checks func(t *testing.T, err error)
	}{
		{
			name: "successful update with name",
			setup: func(ctx context.Context, qtx *Queries) UpdateLocationParams {
				location := createRandomLocation(ctx, qtx)
				return UpdateLocationParams{
					ID:   location.ID,
					Name: util.StringPtr("Updated Location Name"),
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "UpdateLocation() should not error")
			},
		},
		{
			name: "successful update with all fields",
			setup: func(ctx context.Context, qtx *Queries) UpdateLocationParams {
				location := createRandomLocation(ctx, qtx)
				capacity := int32(200)
				return UpdateLocationParams{
					ID:       location.ID,
					Name:     util.StringPtr("Fully Updated Location"),
					Street:   util.StringPtr("789 Updated Blvd"),
					Capacity: &capacity,
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "UpdateLocation() should not error")
			},
		},
		{
			name: "update with nil values (no change)",
			setup: func(ctx context.Context, qtx *Queries) UpdateLocationParams {
				location := createRandomLocation(ctx, qtx)
				return UpdateLocationParams{
					ID: location.ID,
				}
			},
			checks: func(t *testing.T, err error) {
				require.NoError(t, err, "UpdateLocation() should not error even with no updates")
			},
		},
		{
			name: "update non-existent location",
			setup: func(ctx context.Context, qtx *Queries) UpdateLocationParams {
				return UpdateLocationParams{
					ID:   uuid.New(),
					Name: util.StringPtr("Non-existent Location"),
				}
			},
			checks: func(t *testing.T, err error) {
				require.Error(t, err, "UpdateLocation() should error for non-existent location")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			params := tt.setup(ctx, qtx)
			_, err = qtx.UpdateLocation(ctx, params)
			tt.checks(t, err)
		})
	}
}

func TestDeleteLocation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, location Location, err error)
	}{
		{
			name: "successful deletion",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				location := createRandomLocation(ctx, qtx)
				return location.ID
			},
			checks: func(t *testing.T, location Location, err error) {
				require.NoError(t, err, "DeleteLocation() should not error")
				require.NotEmpty(t, location.ID, "deleted location should have ID")
			},
		},
		{
			name: "delete non-existent location",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, location Location, err error) {
				require.Error(t, err, "DeleteLocation() should error for non-existent ID")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			id := tt.setup(ctx, qtx)
			location, err := qtx.DeleteLocation(ctx, id)
			tt.checks(t, location, err)
		})
	}
}

func TestListAllLocations(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries)
		checks func(t *testing.T, locations []ListAllLocationsRow, err error)
	}{
		{
			name: "list all locations with multiple organisations",
			setup: func(ctx context.Context, qtx *Queries) {
				// Create multiple organisations with locations
				for i := 0; i < 2; i++ {
					org := createRandomOrganisation(ctx, qtx)
					for j := 0; j < 2; j++ {
						_, err := qtx.CreateLocation(ctx, CreateLocationParams{
							OrganisationID: org.ID,
							Name:           util.RandomString(10),
							Street:         util.RandomString(20),
							HouseNumber:    util.RandomString(2),
							PostalCode:     util.RandomString(6),
							City:           util.RandomString(10),
						})
						require.NoError(t, err)
					}
				}
			},
			checks: func(t *testing.T, locations []ListAllLocationsRow, err error) {
				require.NoError(t, err, "ListAllLocations() should not error")
				require.GreaterOrEqual(t, len(locations), 4, "should have at least 4 locations")
				// Check that they are ordered by name
				for i := 1; i < len(locations); i++ {
					require.LessOrEqual(t, locations[i-1].Name, locations[i].Name, "locations should be ordered by name")
				}
			},
		},
		{
			name: "list all locations when empty",
			setup: func(ctx context.Context, qtx *Queries) {
				// No setup - keep it empty
			},
			checks: func(t *testing.T, locations []ListAllLocationsRow, err error) {
				require.NoError(t, err, "ListAllLocations() should not error")
				require.Empty(t, locations, "should return empty slice when no locations exist")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tx, err := testDB.Begin(ctx)
			require.NoError(t, err, "failed to begin transaction")
			defer tx.Rollback(ctx)

			qtx := testQueries.WithTx(tx)

			tt.setup(ctx, qtx)
			locations, err := qtx.ListAllLocations(ctx)
			tt.checks(t, locations, err)
		})
	}
}

// Random data generators for tests

func createRandomOrganisation(ctx context.Context, qtx *Queries) Organisation {
	org, err := qtx.CreateOrganisation(ctx, CreateOrganisationParams{
		Name:        util.RandomString(15),
		Street:      util.RandomString(25),
		HouseNumber: util.RandomString(4),
		PostalCode:  util.RandomString(6),
		City:        util.RandomString(12),
	})
	if err != nil {
		panic("failed to create random organisation: " + err.Error())
	}
	return org
}

func createRandomLocation(ctx context.Context, qtx *Queries) Location {
	org := createRandomOrganisation(ctx, qtx)
	location, err := qtx.CreateLocation(ctx, CreateLocationParams{
		OrganisationID: org.ID,
		Name:           util.RandomString(12),
		Street:         util.RandomString(22),
		HouseNumber:    util.RandomString(2),
		PostalCode:     util.RandomString(6),
		City:           util.RandomString(10),
		Capacity:       nil,
	})
	if err != nil {
		panic("failed to create random location: " + err.Error())
	}
	return location
}
