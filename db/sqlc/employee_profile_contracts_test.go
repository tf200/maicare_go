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

func TestAddEmployeeContractDetails(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) AddEmployeeContractDetailsParams
		checks func(t *testing.T, profile EmployeeProfile, err error)
	}{
		{
			name: "successful creation with all contract fields",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeContractDetailsParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				hours := 40.0
				rate := 25.50
				return AddEmployeeContractDetailsParams{
					ID:                employee.ID,
					ContractHours:     &hours,
					ContractStartDate: pgtype.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					ContractEndDate:   pgtype.Date{Time: time.Date(2025, 1, 14, 0, 0, 0, 0, time.UTC), Valid: true},
					ContractType:      NullEmployeeContractTypeEnum{EmployeeContractTypeEnum: EmployeeContractTypeEnumLoondienst, Valid: true},
					ContractRate:      &rate,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "AddEmployeeContractDetails() should not error")
				require.NotNil(t, profile.ContractHours)
				require.Equal(t, 40.0, *profile.ContractHours)
				require.True(t, profile.ContractStartDate.Valid)
				require.Equal(t, 2024, profile.ContractStartDate.Time.Year())
				require.Equal(t, 1, int(profile.ContractStartDate.Time.Month()))
				require.Equal(t, 15, profile.ContractStartDate.Time.Day())
				require.True(t, profile.ContractEndDate.Valid)
				require.Equal(t, 2025, profile.ContractEndDate.Time.Year())
				require.NotNil(t, profile.ContractType)
				require.Equal(t, EmployeeContractTypeEnumLoondienst, profile.ContractType)
				require.NotNil(t, profile.ContractRate)
				require.Equal(t, 25.50, *profile.ContractRate)
			},
		},
		{
			name: "successful update with partial contract fields",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeContractDetailsParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				hours := 30.5
				return AddEmployeeContractDetailsParams{
					ID:            employee.ID,
					ContractHours: &hours,
					// Other fields are nil, so they won't be updated
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "AddEmployeeContractDetails() should not error")
				require.NotNil(t, profile.ContractHours)
				require.Equal(t, 30.5, *profile.ContractHours)
			},
		},
		{
			name: "successful update with only contract type",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeContractDetailsParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEmployeeContractDetailsParams{
					ID:           employee.ID,
					ContractType: NullEmployeeContractTypeEnum{EmployeeContractTypeEnum: EmployeeContractTypeEnumZZP, Valid: true},
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "AddEmployeeContractDetails() should not error")
				require.NotNil(t, profile.ContractType)
				require.Equal(t, EmployeeContractTypeEnumZZP, profile.ContractType)
			},
		},
		{
			name: "successful update with contract rate only",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeContractDetailsParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				rate := 50.0
				return AddEmployeeContractDetailsParams{
					ID:           employee.ID,
					ContractRate: &rate,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "AddEmployeeContractDetails() should not error")
				require.NotNil(t, profile.ContractRate)
				require.Equal(t, 50.0, *profile.ContractRate)
			},
		},
		{
			name: "successful update with both start and end dates",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeContractDetailsParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEmployeeContractDetailsParams{
					ID:                employee.ID,
					ContractStartDate: pgtype.Date{Time: time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC), Valid: true},
					ContractEndDate:   pgtype.Date{Time: time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC), Valid: true},
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "AddEmployeeContractDetails() should not error")
				require.True(t, profile.ContractStartDate.Valid)
				require.Equal(t, 2023, profile.ContractStartDate.Time.Year())
				require.Equal(t, 6, int(profile.ContractStartDate.Time.Month()))
				require.True(t, profile.ContractEndDate.Valid)
				require.Equal(t, 2023, profile.ContractEndDate.Time.Year())
				require.Equal(t, 12, int(profile.ContractEndDate.Time.Month()))
			},
		},
		{
			name: "update non-existent employee should not error",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeContractDetailsParams {
				hours := 40.0
				return AddEmployeeContractDetailsParams{
					ID:            uuid.New(),
					ContractHours: &hours,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				// UUID.Nil check - non-existent employee returns zero value profile
				require.Error(t, err, "AddEmployeeContractDetails() should error for non-existent employee")
				require.Equal(t, uuid.Nil, profile.ID)
			},
		},
		{
			name: "successful update with all fields as nil",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeContractDetailsParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEmployeeContractDetailsParams{
					ID: employee.ID,
					// All contract fields are nil, so nothing should change
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "AddEmployeeContractDetails() should not error")
				require.NotEqual(t, uuid.Nil, profile.ID, "employee ID should be set")
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
			profile, err := qtx.AddEmployeeContractDetails(ctx, params)
			tt.checks(t, profile, err)
		})
	}
}

func TestGetEmployeeContractDetails(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, details GetEmployeeContractDetailsRow, err error)
	}{
		{
			name: "successfully retrieve contract details for existing employee",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)
				hours := 40.0
				rate := 25.50
				_, err := qtx.AddEmployeeContractDetails(ctx, AddEmployeeContractDetailsParams{
					ID:                employee.ID,
					ContractHours:     &hours,
					ContractStartDate: pgtype.Date{Time: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), Valid: true},
					ContractEndDate:   pgtype.Date{Time: time.Date(2025, 1, 14, 0, 0, 0, 0, time.UTC), Valid: true},
					ContractType:      NullEmployeeContractTypeEnum{EmployeeContractTypeEnum: EmployeeContractTypeEnumLoondienst, Valid: true},
					ContractRate:      &rate,
				})
				require.NoError(t, err)
				return employee.ID
			},
			checks: func(t *testing.T, details GetEmployeeContractDetailsRow, err error) {
				require.NoError(t, err, "GetEmployeeContractDetails() should not error for existing employee")
				require.NotNil(t, details.ContractHours)
				require.Equal(t, 40.0, *details.ContractHours)
				require.True(t, details.ContractStartDate.Valid)
				require.Equal(t, 2024, details.ContractStartDate.Time.Year())
				require.NotNil(t, details.ContractType)
				require.Equal(t, EmployeeContractTypeEnumLoondienst, details.ContractType)
				require.NotNil(t, details.ContractRate)
				require.Equal(t, 25.50, *details.ContractRate)
			},
		},
		{
			name: "successfully retrieve contract details with nil optional fields",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return employee.ID
			},
			checks: func(t *testing.T, details GetEmployeeContractDetailsRow, err error) {
				require.NoError(t, err, "GetEmployeeContractDetails() should not error")
				// require.Nil(t, details.ContractHours)
				// require.Nil(t, details.ContractType)
				// require.Nil(t, details.ContractRate)
			},
		},
		{
			name: "error when retrieving contract details for non-existent employee",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				return uuid.New()
			},
			checks: func(t *testing.T, details GetEmployeeContractDetailsRow, err error) {
				require.Error(t, err, "GetEmployeeContractDetails() should error for non-existent employee")
			},
		},
		{
			name: "successfully retrieve is_subcontractor field",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return employee.ID
			},
			checks: func(t *testing.T, details GetEmployeeContractDetailsRow, err error) {
				require.NoError(t, err, "GetEmployeeContractDetails() should not error")
				require.NotNil(t, details.IsSubcontractor)
				require.True(t, *details.IsSubcontractor)
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
			details, err := qtx.GetEmployeeContractDetails(ctx, id)
			tt.checks(t, details, err)
		})
	}
}

func TestUpdateEmployeeIsSubcontractor(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateEmployeeIsSubcontractorParams
		checks func(t *testing.T, profile EmployeeProfile, err error)
	}{
		{
			name: "successful update to mark employee as subcontractor",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeIsSubcontractorParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return UpdateEmployeeIsSubcontractorParams{
					ID:              employee.ID,
					IsSubcontractor: util.BoolPtr(true),
					ContractType:    EmployeeContractTypeEnumZZP,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "UpdateEmployeeIsSubcontractor() should not error")
				require.NotNil(t, profile.IsSubcontractor)
				require.True(t, *profile.IsSubcontractor)
				require.NotNil(t, profile.ContractType)
				require.Equal(t, EmployeeContractTypeEnumZZP, profile.ContractType)
			},
		},
		{
			name: "successful update to mark employee as NOT a subcontractor",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeIsSubcontractorParams {
				employee := createRandomEmployeeProfileIsSubcontractor(ctx, qtx, true)
				return UpdateEmployeeIsSubcontractorParams{
					ID:              employee.ID,
					IsSubcontractor: util.BoolPtr(false),
					ContractType:    EmployeeContractTypeEnumLoondienst,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "UpdateEmployeeIsSubcontractor() should not error")
				require.NotNil(t, profile.IsSubcontractor)
				require.False(t, *profile.IsSubcontractor)
				require.NotNil(t, profile.ContractType)
				require.Equal(t, EmployeeContractTypeEnumLoondienst, profile.ContractType)
			},
		},
		{
			name: "successful update with only is_subcontractor field",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeIsSubcontractorParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return UpdateEmployeeIsSubcontractorParams{
					ID:              employee.ID,
					IsSubcontractor: util.BoolPtr(true),
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.Error(t, err, "UpdateEmployeeIsSubcontractor() should error")
			},
		},
		{
			name: "successful update with only contract type",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeIsSubcontractorParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return UpdateEmployeeIsSubcontractorParams{
					ID:           employee.ID,
					ContractType: EmployeeContractTypeEnumLoondienst,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "UpdateEmployeeIsSubcontractor() should not error")
				require.NotNil(t, profile.ContractType)
				require.Equal(t, EmployeeContractTypeEnumLoondienst, profile.ContractType)
			},
		},
		{
			name: "successful update with both is_subcontractor and contract_type",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeIsSubcontractorParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return UpdateEmployeeIsSubcontractorParams{
					ID:              employee.ID,
					IsSubcontractor: util.BoolPtr(true),
					ContractType:    EmployeeContractTypeEnumLoondienst,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "UpdateEmployeeIsSubcontractor() should not error")
				require.NotNil(t, profile.IsSubcontractor)
				require.True(t, *profile.IsSubcontractor)
				require.NotNil(t, profile.ContractType)
				require.Equal(t, EmployeeContractTypeEnumLoondienst, profile.ContractType)
			},
		},
		{
			name: "update non-existent employee should not error",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeIsSubcontractorParams {
				return UpdateEmployeeIsSubcontractorParams{
					ID:              uuid.New(),
					IsSubcontractor: util.BoolPtr(true),
					ContractType:    EmployeeContractTypeEnumZZP,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				// Non-existent employee returns zero value profile
				require.Error(t, err, "UpdateEmployeeIsSubcontractor() should error for non-existent employee")
			},
		},
		{
			name: "update contract type to various types",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeIsSubcontractorParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return UpdateEmployeeIsSubcontractorParams{
					ID:           employee.ID,
					ContractType: EmployeeContractTypeEnumLoondienst,
				}
			},
			checks: func(t *testing.T, profile EmployeeProfile, err error) {
				require.NoError(t, err, "UpdateEmployeeIsSubcontractor() should not error")
				require.NotNil(t, profile.ContractType)
				require.Equal(t, EmployeeContractTypeEnumLoondienst, profile.ContractType)
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
			profile, err := qtx.UpdateEmployeeIsSubcontractor(ctx, params)
			tt.checks(t, profile, err)
		})
	}
}
