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

func TestAddEducationToEmployeeProfile(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) AddEducationToEmployeeProfileParams
		checks func(t *testing.T, education EmployeeEducation)
	}{
		{
			name: "successful creation with all fields",
			setup: func(ctx context.Context, qtx *Queries) AddEducationToEmployeeProfileParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEducationToEmployeeProfileParams{
					EmployeeID:      employee.UserID,
					InstitutionName: "Stanford University",
					Degree:          "Bachelor of Science",
					FieldOfStudy:    "Computer Science",
					StartDate: pgtype.Date{
						Time:  time.Date(2016, 9, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2020, 5, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, education EmployeeEducation) {
				require.NotZero(t, education.ID, "Education ID should be set")
				require.Equal(t, "Stanford University", education.InstitutionName)
				require.Equal(t, "Bachelor of Science", education.Degree)
				require.Equal(t, "Computer Science", education.FieldOfStudy)
				require.True(t, education.StartDate.Valid, "StartDate should be valid")
				require.Equal(t, time.Date(2016, 9, 1, 0, 0, 0, 0, time.UTC), education.StartDate.Time)
				require.True(t, education.EndDate.Valid, "EndDate should be valid")
				require.Equal(t, time.Date(2020, 5, 31, 0, 0, 0, 0, time.UTC), education.EndDate.Time)
				require.True(t, education.CreatedAt.Valid, "CreatedAt should be valid")
			},
		},
		{
			name: "successful creation with master's degree",
			setup: func(ctx context.Context, qtx *Queries) AddEducationToEmployeeProfileParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEducationToEmployeeProfileParams{
					EmployeeID:      employee.UserID,
					InstitutionName: "MIT",
					Degree:          "Master of Science",
					FieldOfStudy:    "Electrical Engineering",
					StartDate: pgtype.Date{
						Time:  time.Date(2020, 9, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2022, 5, 20, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, education EmployeeEducation) {
				require.Equal(t, "MIT", education.InstitutionName)
				require.Equal(t, "Master of Science", education.Degree)
				require.Equal(t, "Electrical Engineering", education.FieldOfStudy)
			},
		},
		{
			name: "successful creation with different institution",
			setup: func(ctx context.Context, qtx *Queries) AddEducationToEmployeeProfileParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEducationToEmployeeProfileParams{
					EmployeeID:      employee.UserID,
					InstitutionName: "Harvard University",
					Degree:          "Ph.D.",
					FieldOfStudy:    "Biology",
					StartDate: pgtype.Date{
						Time:  time.Date(2019, 1, 10, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, education EmployeeEducation) {
				require.Equal(t, "Harvard University", education.InstitutionName)
				require.Equal(t, "Ph.D.", education.Degree)
				require.Equal(t, "Biology", education.FieldOfStudy)
			},
		},
		{
			name: "successful creation with special characters in field of study",
			setup: func(ctx context.Context, qtx *Queries) AddEducationToEmployeeProfileParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEducationToEmployeeProfileParams{
					EmployeeID:      employee.UserID,
					InstitutionName: "Oxford University",
					Degree:          "Bachelor of Arts",
					FieldOfStudy:    "English Literature & Linguistics",
					StartDate: pgtype.Date{
						Time:  time.Date(2018, 10, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2021, 6, 30, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, education EmployeeEducation) {
				require.Equal(t, "English Literature & Linguistics", education.FieldOfStudy)
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
			education, err := qtx.AddEducationToEmployeeProfile(ctx, params)
			require.NoError(t, err, "AddEducationToEmployeeProfile() should not error")

			tt.checks(t, education)
		})
	}
}

func TestDeleteEmployeeEducation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, education EmployeeEducation, err error)
	}{
		{
			name: "successful deletion of existing education",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				edu := createRandomEducation(ctx, qtx)
				return edu.ID
			},
			checks: func(t *testing.T, education EmployeeEducation, err error) {
				require.NoError(t, err, "DeleteEmployeeEducation() should not error")
				require.NotZero(t, education.ID, "Deleted education should have an ID")
				require.Equal(t, "Stanford University", education.InstitutionName)
			},
		},
		{
			name: "delete non-existent education",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				return 99999 // Non-existent ID
			},
			checks: func(t *testing.T, education EmployeeEducation, err error) {
				require.Error(t, err, "DeleteEmployeeEducation() should error for non-existent ID")
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
			education, err := qtx.DeleteEmployeeEducation(ctx, id)
			tt.checks(t, education, err)
		})
	}
}

func TestListEducations(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, educations []EmployeeEducation, err error)
	}{
		{
			name: "list educations for employee with multiple education records",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)

				// Create first education
				_, err := qtx.AddEducationToEmployeeProfile(ctx, AddEducationToEmployeeProfileParams{
					EmployeeID:      employee.UserID,
					InstitutionName: "Stanford University",
					Degree:          "Bachelor of Science",
					FieldOfStudy:    "Computer Science",
					StartDate: pgtype.Date{
						Time:  time.Date(2016, 9, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2020, 5, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				})
				require.NoError(t, err)

				// Create second education
				_, err = qtx.AddEducationToEmployeeProfile(ctx, AddEducationToEmployeeProfileParams{
					EmployeeID:      employee.UserID,
					InstitutionName: "MIT",
					Degree:          "Master of Science",
					FieldOfStudy:    "Electrical Engineering",
					StartDate: pgtype.Date{
						Time:  time.Date(2020, 9, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2022, 5, 20, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				})
				require.NoError(t, err)

				return employee.UserID
			},
			checks: func(t *testing.T, educations []EmployeeEducation, err error) {
				require.NoError(t, err, "ListEducations() should not error")
				require.Len(t, educations, 2, "Should return 2 education records")
				require.Equal(t, "Stanford University", educations[0].InstitutionName)
				require.Equal(t, "MIT", educations[1].InstitutionName)
			},
		},
		{
			name: "list educations for employee with no education records",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return employee.UserID
			},
			checks: func(t *testing.T, educations []EmployeeEducation, err error) {
				require.NoError(t, err, "ListEducations() should not error")
				require.Empty(t, educations, "Should return empty list for employee with no education records")
			},
		},
		{
			name: "list educations with single education record",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)

				_, err := qtx.AddEducationToEmployeeProfile(ctx, AddEducationToEmployeeProfileParams{
					EmployeeID:      employee.UserID,
					InstitutionName: "Harvard University",
					Degree:          "Ph.D.",
					FieldOfStudy:    "Biology",
					StartDate: pgtype.Date{
						Time:  time.Date(2019, 1, 10, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2023, 12, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				})
				require.NoError(t, err)

				return employee.UserID
			},
			checks: func(t *testing.T, educations []EmployeeEducation, err error) {
				require.NoError(t, err, "ListEducations() should not error")
				require.Len(t, educations, 1, "Should return 1 education record")
				require.Equal(t, "Harvard University", educations[0].InstitutionName)
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

			employeeID := tt.setup(ctx, qtx)
			educations, err := qtx.ListEducations(ctx, employeeID)
			tt.checks(t, educations, err)
		})
	}
}

func TestUpdateEmployeeEducation(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateEmployeeEducationParams
		checks func(t *testing.T, education EmployeeEducation, err error)
	}{
		{
			name: "successful update with all fields",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeEducationParams {
				edu := createRandomEducation(ctx, qtx)
				return UpdateEmployeeEducationParams{
					ID:              edu.ID,
					InstitutionName: util.StringPtr("Updated University"),
					Degree:          util.StringPtr("Master of Science"),
					FieldOfStudy:    util.StringPtr("Data Science"),
					StartDate: pgtype.Date{
						Time:  time.Date(2017, 9, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, education EmployeeEducation, err error) {
				require.NoError(t, err, "UpdateEmployeeEducation() should not error")
				require.Equal(t, "Updated University", education.InstitutionName)
				require.Equal(t, "Master of Science", education.Degree)
				require.Equal(t, "Data Science", education.FieldOfStudy)
				require.Equal(t, time.Date(2017, 9, 1, 0, 0, 0, 0, time.UTC), education.StartDate.Time)
				require.Equal(t, time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC), education.EndDate.Time)
			},
		},
		{
			name: "successful update with partial fields (only institution)",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeEducationParams {
				edu := createRandomEducation(ctx, qtx)
				return UpdateEmployeeEducationParams{
					ID:              edu.ID,
					InstitutionName: util.StringPtr("New Institution Name"),
					Degree:          nil,
					FieldOfStudy:    nil,
					StartDate:       pgtype.Date{Valid: false},
					EndDate:         pgtype.Date{Valid: false},
				}
			},
			checks: func(t *testing.T, education EmployeeEducation, err error) {
				require.NoError(t, err, "UpdateEmployeeEducation() should not error")
				require.Equal(t, "New Institution Name", education.InstitutionName)
				require.Equal(t, "Bachelor of Science", education.Degree, "Degree should remain unchanged")
				require.Equal(t, "Computer Science", education.FieldOfStudy, "FieldOfStudy should remain unchanged")
			},
		},
		{
			name: "successful update with partial fields (only dates)",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeEducationParams {
				edu := createRandomEducation(ctx, qtx)
				return UpdateEmployeeEducationParams{
					ID:              edu.ID,
					InstitutionName: nil,
					Degree:          nil,
					FieldOfStudy:    nil,
					StartDate: pgtype.Date{
						Time:  time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, education EmployeeEducation, err error) {
				require.NoError(t, err, "UpdateEmployeeEducation() should not error")
				require.Equal(t, "Stanford University", education.InstitutionName, "InstitutionName should remain unchanged")
				require.Equal(t, time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC), education.StartDate.Time)
				require.Equal(t, time.Date(2019, 12, 31, 0, 0, 0, 0, time.UTC), education.EndDate.Time)
			},
		},
		{
			name: "successful update with degree only",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeEducationParams {
				edu := createRandomEducation(ctx, qtx)
				return UpdateEmployeeEducationParams{
					ID:              edu.ID,
					InstitutionName: nil,
					Degree:          util.StringPtr("Associate Degree"),
					FieldOfStudy:    nil,
					StartDate:       pgtype.Date{Valid: false},
					EndDate:         pgtype.Date{Valid: false},
				}
			},
			checks: func(t *testing.T, education EmployeeEducation, err error) {
				require.NoError(t, err, "UpdateEmployeeEducation() should not error")
				require.Equal(t, "Associate Degree", education.Degree)
				require.Equal(t, "Stanford University", education.InstitutionName, "InstitutionName should remain unchanged")
			},
		},
		{
			name: "update non-existent education",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeEducationParams {
				return UpdateEmployeeEducationParams{
					ID:              99999,
					InstitutionName: util.StringPtr("Updated University"),
					Degree:          util.StringPtr("Master of Science"),
					FieldOfStudy:    util.StringPtr("Data Science"),
					StartDate: pgtype.Date{
						Time:  time.Date(2017, 9, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, education EmployeeEducation, err error) {
				require.Error(t, err, "UpdateEmployeeEducation() should error for non-existent ID")
			},
		},
		{
			name: "update with empty/nil dates",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeEducationParams {
				edu := createRandomEducation(ctx, qtx)
				return UpdateEmployeeEducationParams{
					ID:              edu.ID,
					InstitutionName: nil,
					Degree:          nil,
					FieldOfStudy:    nil,
					StartDate:       pgtype.Date{Valid: false},
					EndDate:         pgtype.Date{Valid: false},
				}
			},
			checks: func(t *testing.T, education EmployeeEducation, err error) {
				require.NoError(t, err, "UpdateEmployeeEducation() should not error")
				require.Equal(t, "Stanford University", education.InstitutionName)
				require.Equal(t, "Bachelor of Science", education.Degree)
				require.Equal(t, "Computer Science", education.FieldOfStudy)
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
			education, err := qtx.UpdateEmployeeEducation(ctx, params)
			tt.checks(t, education, err)
		})
	}
}

// Random data generators for tests

func createRandomEducation(ctx context.Context, qtx *Queries) EmployeeEducation {
	employee := createRandomEmployeeProfile(ctx, qtx)

	edu, err := qtx.AddEducationToEmployeeProfile(ctx, AddEducationToEmployeeProfileParams{
		EmployeeID:      employee.UserID,
		InstitutionName: "Stanford University",
		Degree:          "Bachelor of Science",
		FieldOfStudy:    "Computer Science",
		StartDate: pgtype.Date{
			Time:  time.Date(2016, 9, 1, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  time.Date(2020, 5, 31, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	})
	if err != nil {
		panic("failed to create random education: " + err.Error())
	}
	return edu
}
