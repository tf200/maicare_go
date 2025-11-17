package db

import (
	"context"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateEmployeeProfile(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CreateEmployeeProfileParams
		checks func(t *testing.T, profile EmployeeProfile)
	}{
		{
			name: "successful creation with minimal required fields",
			setup: func(ctx context.Context, qtx *Queries) CreateEmployeeProfileParams {
				user := createRandomUser(ctx, qtx)
				location := createRandomLocation(ctx, qtx)
				return CreateEmployeeProfileParams{
					UserID:          user.ID,
					FirstName:       "John",
					LastName:        "Doe",
					Email:           util.RandomEmail(),
					DateOfBirth:     pgtype.Date{Time: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					WorkPhoneNumber: util.StringPtr("+1234567890"),
					LocationID:      &location.ID,
					ContractType:    EmployeeContractTypeEnumLoondienst,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile) {
				require.Equal(t, "John", profile.FirstName)
				require.Equal(t, "Doe", profile.LastName)
				require.NotNil(t, profile.WorkPhoneNumber)
				require.Equal(t, "+1234567890", *profile.WorkPhoneNumber)
				require.NotNil(t, profile.LocationID)
				require.Equal(t, int64(1), *profile.LocationID)
				require.NotNil(t, profile.ContractType)
				require.Equal(t, EmployeeContractTypeEnumLoondienst, profile.ContractType)
				require.False(t, profile.IsArchived)
				require.Nil(t, profile.OutOfService)
				require.Nil(t, profile.HasBorrowed)
				require.False(t, profile.IsSubcontractor != nil && *profile.IsSubcontractor)
			},
		},
		{
			name: "successful creation with all optional fields",
			setup: func(ctx context.Context, qtx *Queries) CreateEmployeeProfileParams {
				user := createRandomUser(ctx, qtx)
				location := createRandomLocation(ctx, qtx)
				return CreateEmployeeProfileParams{
					UserID:                    user.ID,
					FirstName:                 "Jane",
					LastName:                  "Smith",
					Position:                  util.StringPtr("Software Engineer"),
					Department:                util.StringPtr("Engineering"),
					EmployeeNumber:            util.StringPtr("EMP001"),
					EmploymentNumber:          util.StringPtr("EMPL001"),
					PrivateEmailAddress:       util.StringPtr("jane.private@example.com"),
					Email:                     util.RandomEmail(),
					AuthenticationPhoneNumber: util.StringPtr("+0987654321"),
					PrivatePhoneNumber:        util.StringPtr("+1122334455"),
					WorkPhoneNumber:           util.StringPtr("+1234567890"),
					DateOfBirth:               pgtype.Date{Time: time.Date(1985, 5, 20, 0, 0, 0, 0, time.UTC), Valid: true},
					HomeTelephoneNumber:       util.StringPtr("+5544332211"),
					IsSubcontractor:           util.BoolPtr(true),
					Gender:                    EmployeeGenderEnumFemale,
					LocationID:                &location.ID,
					ContractType:              EmployeeContractTypeEnumLoondienst,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile) {
				require.Equal(t, "Jane", profile.FirstName)
				require.Equal(t, "Smith", profile.LastName)
				require.NotNil(t, profile.Position)
				require.Equal(t, "Software Engineer", *profile.Position)
				require.NotNil(t, profile.Department)
				require.Equal(t, "Engineering", *profile.Department)
				require.NotNil(t, profile.EmployeeNumber)
				require.Equal(t, "EMP001", *profile.EmployeeNumber)
				require.NotNil(t, profile.EmploymentNumber)
				require.Equal(t, "EMPL001", *profile.EmploymentNumber)
				require.NotNil(t, profile.PrivateEmailAddress)
				require.Equal(t, "jane.private@example.com", *profile.PrivateEmailAddress)
				require.NotNil(t, profile.AuthenticationPhoneNumber)
				require.Equal(t, "+0987654321", *profile.AuthenticationPhoneNumber)
				require.NotNil(t, profile.PrivatePhoneNumber)
				require.Equal(t, "+1122334455", *profile.PrivatePhoneNumber)
				require.NotNil(t, profile.HomeTelephoneNumber)
				require.Equal(t, "+5544332211", *profile.HomeTelephoneNumber)
				require.NotNil(t, profile.IsSubcontractor)
				require.True(t, *profile.IsSubcontractor)
				require.NotNil(t, profile.Gender)
				require.Equal(t, EmployeeGenderEnumFemale, profile.Gender)
				require.NotNil(t, profile.LocationID)
				require.Equal(t, int64(2), *profile.LocationID)
				require.NotNil(t, profile.ContractType)
				require.Equal(t, EmployeeContractTypeEnumLoondienst, profile.ContractType)
			},
		},
		{
			name: "creation with nil optional fields",
			setup: func(ctx context.Context, qtx *Queries) CreateEmployeeProfileParams {
				user := createRandomUser(ctx, qtx)
				return CreateEmployeeProfileParams{
					UserID:      user.ID,
					FirstName:   "Bob",
					LastName:    "Wilson",
					Email:       util.RandomEmail(),
					DateOfBirth: pgtype.Date{Time: time.Date(1975, 3, 10, 0, 0, 0, 0, time.UTC), Valid: true},
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile) {
				require.Equal(t, "Bob", profile.FirstName)
				require.Equal(t, "Wilson", profile.LastName)
				require.Nil(t, profile.Position)
				require.Nil(t, profile.Department)
				require.Nil(t, profile.EmployeeNumber)
				require.Nil(t, profile.EmploymentNumber)
				require.Nil(t, profile.PrivateEmailAddress)
				require.Nil(t, profile.AuthenticationPhoneNumber)
				require.Nil(t, profile.PrivatePhoneNumber)
				require.Nil(t, profile.WorkPhoneNumber)
				require.Nil(t, profile.HomeTelephoneNumber)
				require.Nil(t, profile.IsSubcontractor)
				require.Nil(t, profile.Gender)
				require.Nil(t, profile.LocationID)
				require.Nil(t, profile.ContractType)
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
			profile, err := qtx.CreateEmployeeProfile(ctx, params)
			require.NoError(t, err, "CreateEmployeeProfile() should not error")

			tt.checks(t, profile)
		})
	}
}

func TestGetEmployeeProfileByID(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, profile GetEmployeeProfileByIDRow, err error)
	}{
		{
			name: "get existing employee profile by ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				profile := createRandomEmployeeProfile(ctx, qtx)
				return profile.ID
			},
			checks: func(t *testing.T, profile GetEmployeeProfileByIDRow, err error) {
				require.NoError(t, err, "GetEmployeeProfileByID() should not error")
				require.NotEmpty(t, profile.ID)
				require.NotEmpty(t, profile.UserID)
				require.NotEmpty(t, profile.FirstName)
				require.NotEmpty(t, profile.LastName)
				require.NotEmpty(t, profile.Email)
			},
		},
		{
			name: "get non-existent employee profile by ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, profile GetEmployeeProfileByIDRow, err error) {
				require.Error(t, err, "GetEmployeeProfileByID() should error for non-existent ID")
			},
		},
		{
			name: "get employee profile with profile picture",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				profile := createRandomEmployeeProfileWithProfilePicture(ctx, qtx)
				return profile.ID
			},
			checks: func(t *testing.T, profile GetEmployeeProfileByIDRow, err error) {
				require.NoError(t, err, "GetEmployeeProfileByID() should not error")
				require.NotNil(t, profile.ProfilePicture)
				require.NotEmpty(t, *profile.ProfilePicture)
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
			profile, err := qtx.GetEmployeeProfileByID(ctx, id)
			tt.checks(t, profile, err)
		})
	}
}

func TestGetEmployeeProfileByUserID(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, profile GetEmployeeProfileByUserIDRow, err error)
	}{
		{
			name: "get existing employee profile by user ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				profile := createRandomEmployeeProfile(ctx, qtx)
				return profile.UserID
			},
			checks: func(t *testing.T, profile GetEmployeeProfileByUserIDRow, err error) {
				require.NoError(t, err, "GetEmployeeProfileByUserID() should not error")
				require.NotEmpty(t, profile.UserID)
				require.NotEmpty(t, profile.EmployeeID)
				require.NotEmpty(t, profile.FirstName)
				require.NotEmpty(t, profile.LastName)
				require.NotEmpty(t, profile.Email)
			},
		},
		{
			name: "get non-existent employee profile by user ID",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, profile GetEmployeeProfileByUserIDRow, err error) {
				require.Error(t, err, "GetEmployeeProfileByUserID() should error for non-existent user ID")
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

			userID := tt.setup(ctx, qtx)
			profile, err := qtx.GetEmployeeProfileByUserID(ctx, userID)
			tt.checks(t, profile, err)
		})
	}
}

func TestListEmployeeProfile(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) ListEmployeeProfileParams
		checks func(t *testing.T, profiles []ListEmployeeProfileRow, err error)
	}{
		{
			name: "list employee profiles with multiple entries",
			setup: func(ctx context.Context, qtx *Queries) ListEmployeeProfileParams {
				// Create multiple employee profiles
				for i := 0; i < 3; i++ {
					createRandomEmployeeProfile(ctx, qtx)
				}
				return ListEmployeeProfileParams{
					Limit:               10,
					Offset:              0,
					IncludeArchived:     util.BoolPtr(false),
					IncludeOutOfService: util.BoolPtr(false),
				}
			},
			checks: func(t *testing.T, profiles []ListEmployeeProfileRow, err error) {
				require.NoError(t, err, "ListEmployeeProfile() should not error")
				require.GreaterOrEqual(t, len(profiles), 3, "should have at least 3 employee profiles")
				// Check that they are ordered by created_at DESC
				for i := 1; i < len(profiles); i++ {
					require.False(t, profiles[i-1].CreatedAt.Time.Before(profiles[i].CreatedAt.Time), "should be ordered by created_at DESC")
				}
			},
		},
		{
			name: "list employee profiles with department filter",
			setup: func(ctx context.Context, qtx *Queries) ListEmployeeProfileParams {
				// Create employees in specific departments
				for i := 0; i < 2; i++ {
					createRandomEmployeeProfileWithDepartment(ctx, qtx, "Engineering")
				}
				// Create employees in other departments
				createRandomEmployeeProfileWithDepartment(ctx, qtx, "Sales")
				return ListEmployeeProfileParams{
					Limit:      10,
					Offset:     0,
					Department: util.StringPtr("Engineering"),
				}
			},
			checks: func(t *testing.T, profiles []ListEmployeeProfileRow, err error) {
				require.NoError(t, err, "ListEmployeeProfile() should not error")
				for _, profile := range profiles {
					if profile.Department != nil {
						require.Equal(t, "Engineering", *profile.Department, "should only return Engineering department profiles")
					}
				}
			},
		},
		{
			name: "list employee profiles with search term",
			setup: func(ctx context.Context, qtx *Queries) ListEmployeeProfileParams {
				user := createRandomUser(ctx, qtx)
				_, err := qtx.CreateEmployeeProfile(ctx, CreateEmployeeProfileParams{
					UserID:      user.ID,
					FirstName:   "Searchable",
					LastName:    "Name",
					Email:       "searchable@example.com",
					DateOfBirth: pgtype.Date{Time: time.Date(1992, 7, 25, 0, 0, 0, 0, time.UTC), Valid: true},
				})
				require.NoError(t, err)
				return ListEmployeeProfileParams{
					Limit:  10,
					Offset: 0,
					Search: util.StringPtr("Searchable"),
				}
			},
			checks: func(t *testing.T, profiles []ListEmployeeProfileRow, err error) {
				require.NoError(t, err, "ListEmployeeProfile() should not error")
				found := false
				for _, profile := range profiles {
					if profile.FirstName == "Searchable" || profile.LastName == "Searchable" {
						found = true
						break
					}
				}
				require.True(t, found, "should find the searchable employee")
			},
		},
		{
			name: "list employee profiles with pagination",
			setup: func(ctx context.Context, qtx *Queries) ListEmployeeProfileParams {
				// Create multiple employee profiles
				for i := 0; i < 5; i++ {
					createRandomEmployeeProfile(ctx, qtx)
				}
				return ListEmployeeProfileParams{
					Limit:  2,
					Offset: 0,
				}
			},
			checks: func(t *testing.T, profiles []ListEmployeeProfileRow, err error) {
				require.NoError(t, err, "ListEmployeeProfile() should not error")
				require.LessOrEqual(t, len(profiles), 2, "should return at most 2 profiles due to limit")
			},
		},
		{
			name: "list employee profiles when empty",
			setup: func(ctx context.Context, qtx *Queries) ListEmployeeProfileParams {
				return ListEmployeeProfileParams{
					Limit:  10,
					Offset: 0,
				}
			},
			checks: func(t *testing.T, profiles []ListEmployeeProfileRow, err error) {
				require.NoError(t, err, "ListEmployeeProfile() should not error")
				require.Empty(t, profiles, "should return empty slice when no employee profiles exist")
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
			profiles, err := qtx.ListEmployeeProfile(ctx, params)
			tt.checks(t, profiles, err)
		})
	}
}

func TestCountEmployeeProfile(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) CountEmployeeProfileParams
		checks func(t *testing.T, count int64, err error)
	}{
		{
			name: "count employee profiles with multiple entries",
			setup: func(ctx context.Context, qtx *Queries) CountEmployeeProfileParams {
				// Create multiple employee profiles
				for i := 0; i < 3; i++ {
					createRandomEmployeeProfile(ctx, qtx)
				}
				return CountEmployeeProfileParams{
					IncludeArchived:     util.BoolPtr(false),
					IncludeOutOfService: util.BoolPtr(false),
				}
			},
			checks: func(t *testing.T, count int64, err error) {
				require.NoError(t, err, "CountEmployeeProfile() should not error")
				require.GreaterOrEqual(t, count, int64(3), "should count at least 3 employee profiles")
			},
		},
		{
			name: "count employee profiles with department filter",
			setup: func(ctx context.Context, qtx *Queries) CountEmployeeProfileParams {
				// Create employees in specific departments
				for i := 0; i < 2; i++ {
					createRandomEmployeeProfileWithDepartment(ctx, qtx, "Engineering")
				}
				createRandomEmployeeProfileWithDepartment(ctx, qtx, "Sales")
				return CountEmployeeProfileParams{
					Department: util.StringPtr("Engineering"),
				}
			},
			checks: func(t *testing.T, count int64, err error) {
				require.NoError(t, err, "CountEmployeeProfile() should not error")
				require.GreaterOrEqual(t, count, int64(2), "should count at least 2 Engineering employees")
			},
		},
		{
			name: "count employee profiles when empty",
			setup: func(ctx context.Context, qtx *Queries) CountEmployeeProfileParams {
				return CountEmployeeProfileParams{}
			},
			checks: func(t *testing.T, count int64, err error) {
				require.NoError(t, err, "CountEmployeeProfile() should not error")
				require.Equal(t, int64(0), count, "should count 0 when no employee profiles exist")
			},
		},
		{
			name: "count employee profiles with location filter",
			setup: func(ctx context.Context, qtx *Queries) CountEmployeeProfileParams {
				location := createRandomLocation(ctx, qtx)
				for i := 0; i < 2; i++ {
					createRandomEmployeeProfileWithLocation(ctx, qtx, location.ID)
				}
				location2 := createRandomLocation(ctx, qtx)
				createRandomEmployeeProfileWithLocation(ctx, qtx, location2.ID)
				return CountEmployeeProfileParams{
					LocationID: &location.ID,
				}
			},
			checks: func(t *testing.T, count int64, err error) {
				require.NoError(t, err, "CountEmployeeProfile() should not error")
				require.GreaterOrEqual(t, count, int64(2), "should count at least 2 employees at location 1")
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
			count, err := qtx.CountEmployeeProfile(ctx, params)
			tt.checks(t, count, err)
		})
	}
}

func TestUpdateEmployeeProfile(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateEmployeeProfileParams
		checks func(t *testing.T, profile EmployeeProfile, err error)
	}{
		{
			name: "successful update with name",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeProfileParams {
				profile := createRandomEmployeeProfile(ctx, qtx)
				return UpdateEmployeeProfileParams{
					ID:        profile.ID,
					FirstName: util.StringPtr("Updated Name"),
					LastName:  util.StringPtr("Updated Last"),
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "UpdateEmployeeProfile() should not error")
				require.Equal(t, "Updated Name", profile.FirstName)
				require.Equal(t, "Updated Last", profile.LastName)
			},
		},
		{
			name: "successful update with all fields",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeProfileParams {
				profile := createRandomEmployeeProfile(ctx, qtx)
				location := createRandomLocation(ctx, qtx)
				return UpdateEmployeeProfileParams{
					ID:                        profile.ID,
					FirstName:                 util.StringPtr("Fully"),
					LastName:                  util.StringPtr("Updated"),
					Position:                  util.StringPtr("Senior Developer"),
					Department:                util.StringPtr("Tech"),
					EmployeeNumber:            util.StringPtr("EMP999"),
					EmploymentNumber:          util.StringPtr("EMPL999"),
					PrivateEmailAddress:       util.StringPtr("private@example.com"),
					Email:                     util.StringPtr("updated@example.com"),
					AuthenticationPhoneNumber: util.StringPtr("+999888777"),
					PrivatePhoneNumber:        util.StringPtr("+666555444"),
					WorkPhoneNumber:           util.StringPtr("+333222111"),
					DateOfBirth:               pgtype.Date{Time: time.Date(1990, 12, 31, 0, 0, 0, 0, time.UTC), Valid: true},
					HomeTelephoneNumber:       util.StringPtr("+111222333"),
					IsSubcontractor:           util.BoolPtr(true),
					Gender:                    NullEmployeeGenderEnum{EmployeeGenderEnum: EmployeeGenderEnumNotSpecified, Valid: true},
					LocationID:                &location.ID,
					HasBorrowed:               util.BoolPtr(true),
					OutOfService:              util.BoolPtr(false),
					IsArchived:                util.BoolPtr(false),
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "UpdateEmployeeProfile() should not error")
				require.Equal(t, "Fully", profile.FirstName)
				require.Equal(t, "Updated", profile.LastName)
				require.NotNil(t, profile.Position)
				require.Equal(t, "Senior Developer", *profile.Position)
				require.NotNil(t, profile.Department)
				require.Equal(t, "Tech", *profile.Department)
				require.NotNil(t, profile.EmployeeNumber)
				require.Equal(t, "EMP999", *profile.EmployeeNumber)
				require.NotNil(t, profile.EmploymentNumber)
				require.Equal(t, "EMPL999", *profile.EmploymentNumber)
				require.NotNil(t, profile.PrivateEmailAddress)
				require.Equal(t, "private@example.com", *profile.PrivateEmailAddress)
				require.Equal(t, "updated@example.com", profile.Email)
				require.NotNil(t, profile.AuthenticationPhoneNumber)
				require.Equal(t, "+999888777", *profile.AuthenticationPhoneNumber)
				require.NotNil(t, profile.PrivatePhoneNumber)
				require.Equal(t, "+666555444", *profile.PrivatePhoneNumber)
				require.NotNil(t, profile.WorkPhoneNumber)
				require.Equal(t, "+333222111", *profile.WorkPhoneNumber)
				require.NotNil(t, profile.HomeTelephoneNumber)
				require.Equal(t, "+111222333", *profile.HomeTelephoneNumber)
				require.NotNil(t, profile.IsSubcontractor)
				require.True(t, *profile.IsSubcontractor)
				require.Equal(t, EmployeeGenderEnumNotSpecified, profile.Gender)
				require.NotNil(t, profile.LocationID)
				require.Equal(t, int64(99), *profile.LocationID)
				require.NotNil(t, profile.HasBorrowed)
				require.True(t, profile.HasBorrowed)
				require.NotNil(t, profile.OutOfService)
				require.False(t, *profile.OutOfService)
				require.False(t, profile.IsArchived)
			},
		},
		{
			name: "update with nil values (no change)",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeProfileParams {
				profile := createRandomEmployeeProfile(ctx, qtx)
				return UpdateEmployeeProfileParams{
					ID: profile.ID,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "UpdateEmployeeProfile() should not error even with no updates")
			},
		},
		{
			name: "update non-existent employee profile",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeProfileParams {
				return UpdateEmployeeProfileParams{
					ID:        uuid.New(),
					FirstName: util.StringPtr("Non-existent"),
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "UpdateEmployeeProfile() should not error even for non-existent profile")
			},
		},
		{
			name: "update with archived status",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeProfileParams {
				profile := createRandomEmployeeProfile(ctx, qtx)
				return UpdateEmployeeProfileParams{
					ID:         profile.ID,
					IsArchived: util.BoolPtr(true),
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "UpdateEmployeeProfile() should not error when archiving")
				require.True(t, profile.IsArchived, "profile should be archived")
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
			profile, err := qtx.UpdateEmployeeProfile(ctx, params)
			tt.checks(t, profile, err)
		})
	}
}

// Random data generators for tests

func createRandomEmployeeProfile(ctx context.Context, qtx *Queries) EmployeeProfile {
	user := createRandomUser(ctx, qtx)
	location := createRandomLocation(ctx, qtx)
	profile, err := qtx.CreateEmployeeProfile(ctx, CreateEmployeeProfileParams{
		UserID:                    user.ID,
		FirstName:                 util.RandomString(8),
		LastName:                  util.RandomString(10),
		Position:                  util.StringPtr(util.RandomString(15)),
		Department:                util.StringPtr(util.RandomString(12)),
		EmployeeNumber:            util.StringPtr("EMP" + util.RandomString(5)),
		EmploymentNumber:          util.StringPtr("EMPL" + util.RandomString(5)),
		PrivateEmailAddress:       util.StringPtr(util.RandomEmail()),
		Email:                     util.RandomEmail(),
		AuthenticationPhoneNumber: util.StringPtr("+1" + util.RandomString(10)),
		PrivatePhoneNumber:        util.StringPtr("+1" + util.RandomString(10)),
		WorkPhoneNumber:           util.StringPtr("+1" + util.RandomString(10)),
		DateOfBirth:               pgtype.Date{Time: time.Date(1990, 4, 4, 0, 0, 0, 0, time.UTC), Valid: true},
		HomeTelephoneNumber:       util.StringPtr("+1" + util.RandomString(10)),
		IsSubcontractor:           util.BoolPtr(util.RandomInt(0, 1) == 1),
		Gender:                    EmployeeGenderEnumMale,
		LocationID:                &location.ID,
		ContractType:              EmployeeContractTypeEnumLoondienst,
	})
	if err != nil {
		panic("failed to create random employee profile: " + err.Error())
	}
	return profile
}

func createRandomEmployeeProfileWithDepartment(ctx context.Context, qtx *Queries, department string) EmployeeProfile {
	user := createRandomUser(ctx, qtx)
	profile, err := qtx.CreateEmployeeProfile(ctx, CreateEmployeeProfileParams{
		UserID:      user.ID,
		FirstName:   util.RandomString(8),
		LastName:    util.RandomString(10),
		Email:       util.RandomEmail(),
		Department:  util.StringPtr(department),
		DateOfBirth: pgtype.Date{Time: time.Date(1990, 4, 4, 0, 0, 0, 0, time.UTC), Valid: true},
	})
	if err != nil {
		panic("failed to create random employee profile with department: " + err.Error())
	}
	return profile
}

func createRandomEmployeeProfileWithLocation(ctx context.Context, qtx *Queries, locationID int64) EmployeeProfile {
	user := createRandomUser(ctx, qtx)
	profile, err := qtx.CreateEmployeeProfile(ctx, CreateEmployeeProfileParams{
		UserID:      user.ID,
		FirstName:   util.RandomString(8),
		LastName:    util.RandomString(10),
		Email:       util.RandomEmail(),
		LocationID:  &locationID,
		DateOfBirth: pgtype.Date{Time: time.Date(1990, 4, 4, 0, 0, 0, 0, time.UTC), Valid: true},
	})
	if err != nil {
		panic("failed to create random employee profile with location: " + err.Error())
	}
	return profile
}

func createRandomEmployeeProfileWithProfilePicture(ctx context.Context, qtx *Queries) EmployeeProfile {
	user := createRandomUserWithProfilePicture(ctx, qtx)
	profile, err := qtx.CreateEmployeeProfile(ctx, CreateEmployeeProfileParams{
		UserID:      user.ID,
		FirstName:   util.RandomString(8),
		LastName:    util.RandomString(10),
		Email:       util.RandomEmail(),
		DateOfBirth: pgtype.Date{Time: time.Date(1990, 4, 4, 0, 0, 0, 0, time.UTC), Valid: true},
	})
	if err != nil {
		panic("failed to create random employee profile with profile picture: " + err.Error())
	}
	return profile
}

func createRandomUserWithProfilePicture(ctx context.Context, qtx *Queries) CustomUser {
	user, err := qtx.CreateUser(ctx, CreateUserParams{
		Password:       "randomhashedpassword",
		Email:          util.RandomEmail(),
		IsActive:       true,
		ProfilePicture: util.StringPtr("https://example.com/profile/" + util.RandomString(10) + ".jpg"),
	})
	if err != nil {
		panic("failed to create random user with profile picture: " + err.Error())
	}
	return user
}
