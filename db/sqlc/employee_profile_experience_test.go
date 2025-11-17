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

func TestAddEmployeeExperience(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) AddEmployeeExperienceParams
		checks func(t *testing.T, experience EmployeeExperience)
	}{
		{
			name: "successful creation with minimal required fields",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeExperienceParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEmployeeExperienceParams{
					EmployeeID:  employee.UserID,
					JobTitle:    "Software Engineer",
					CompanyName: "Tech Corp",
					StartDate: pgtype.Date{
						Time:  time.Date(2018, 1, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: nil,
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience) {
				require.NotZero(t, experience.ID, "Experience ID should be set")
				require.Equal(t, "Software Engineer", experience.JobTitle)
				require.Equal(t, "Tech Corp", experience.CompanyName)
				require.True(t, experience.StartDate.Valid, "StartDate should be valid")
				require.Equal(t, time.Date(2018, 1, 15, 0, 0, 0, 0, time.UTC), experience.StartDate.Time)
				require.True(t, experience.EndDate.Valid, "EndDate should be valid")
				require.Equal(t, time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC), experience.EndDate.Time)
				require.Nil(t, experience.Description, "Description should be nil")
				require.True(t, experience.CreatedAt.Valid, "CreatedAt should be valid")
			},
		},
		{
			name: "successful creation with description",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeExperienceParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEmployeeExperienceParams{
					EmployeeID:  employee.UserID,
					JobTitle:    "Senior Developer",
					CompanyName: "Innovation Labs",
					StartDate: pgtype.Date{
						Time:  time.Date(2021, 3, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: util.StringPtr("Led backend development team, implemented microservices architecture"),
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience) {
				require.Equal(t, "Senior Developer", experience.JobTitle)
				require.Equal(t, "Innovation Labs", experience.CompanyName)
				require.NotNil(t, experience.Description)
				require.Equal(t, "Led backend development team, implemented microservices architecture", *experience.Description)
			},
		},
		{
			name: "successful creation with healthcare role",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeExperienceParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEmployeeExperienceParams{
					EmployeeID:  employee.UserID,
					JobTitle:    "Registered Nurse",
					CompanyName: "City Medical Center",
					StartDate: pgtype.Date{
						Time:  time.Date(2019, 7, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2023, 11, 30, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: util.StringPtr("Critical care nursing, patient assessment and monitoring"),
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience) {
				require.Equal(t, "Registered Nurse", experience.JobTitle)
				require.Equal(t, "City Medical Center", experience.CompanyName)
			},
		},
		{
			name: "successful creation with long description",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeExperienceParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				longDesc := "Responsible for managing the entire software development lifecycle including requirements gathering, architecture design, implementation, testing, and deployment. Worked with cross-functional teams to deliver high-quality solutions."
				return AddEmployeeExperienceParams{
					EmployeeID:  employee.UserID,
					JobTitle:    "Project Manager",
					CompanyName: "Global Consulting Group",
					StartDate: pgtype.Date{
						Time:  time.Date(2017, 5, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2021, 4, 30, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: util.StringPtr(longDesc),
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience) {
				require.NotNil(t, experience.Description)
				require.Contains(t, *experience.Description, "software development lifecycle")
			},
		},
		{
			name: "successful creation with special characters in company name",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeExperienceParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEmployeeExperienceParams{
					EmployeeID:  employee.UserID,
					JobTitle:    "Solutions Architect",
					CompanyName: "O'Reilly & Associates, Inc.",
					StartDate: pgtype.Date{
						Time:  time.Date(2020, 2, 10, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2022, 8, 20, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: nil,
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience) {
				require.Equal(t, "O'Reilly & Associates, Inc.", experience.CompanyName)
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
			experience, err := qtx.AddEmployeeExperience(ctx, params)
			require.NoError(t, err, "AddEmployeeExperience() should not error")

			tt.checks(t, experience)
		})
	}
}

func TestDeleteEmployeeExperience(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, experience EmployeeExperience, err error)
	}{
		{
			name: "successful deletion of existing experience",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				exp := createRandomExperience(ctx, qtx)
				return exp.ID
			},
			checks: func(t *testing.T, experience EmployeeExperience, err error) {
				require.NoError(t, err, "DeleteEmployeeExperience() should not error")
				require.NotZero(t, experience.ID, "Deleted experience should have an ID")
				require.Equal(t, "Software Engineer", experience.JobTitle)
				require.Equal(t, "Tech Corp", experience.CompanyName)
			},
		},
		{
			name: "delete non-existent experience",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				return 99999 // Non-existent ID
			},
			checks: func(t *testing.T, experience EmployeeExperience, err error) {
				require.Error(t, err, "DeleteEmployeeExperience() should error for non-existent ID")
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
			experience, err := qtx.DeleteEmployeeExperience(ctx, id)
			tt.checks(t, experience, err)
		})
	}
}

func TestListEmployeeExperience(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, experiences []EmployeeExperience, err error)
	}{
		{
			name: "list experiences for employee with multiple experiences",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)

				// Create first experience
				_, err := qtx.AddEmployeeExperience(ctx, AddEmployeeExperienceParams{
					EmployeeID:  employee.UserID,
					JobTitle:    "Junior Developer",
					CompanyName: "StartUp Inc",
					StartDate: pgtype.Date{
						Time:  time.Date(2016, 6, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2018, 5, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: util.StringPtr("Frontend development with React"),
				})
				require.NoError(t, err)

				// Create second experience
				_, err = qtx.AddEmployeeExperience(ctx, AddEmployeeExperienceParams{
					EmployeeID:  employee.UserID,
					JobTitle:    "Mid-level Developer",
					CompanyName: "Tech Solutions",
					StartDate: pgtype.Date{
						Time:  time.Date(2018, 7, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: util.StringPtr("Full stack development"),
				})
				require.NoError(t, err)

				// Create third experience
				_, err = qtx.AddEmployeeExperience(ctx, AddEmployeeExperienceParams{
					EmployeeID:  employee.UserID,
					JobTitle:    "Senior Developer",
					CompanyName: "Enterprise Corp",
					StartDate: pgtype.Date{
						Time:  time.Date(2021, 1, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2024, 11, 30, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: util.StringPtr("Team lead, architecture design"),
				})
				require.NoError(t, err)

				return employee.UserID
			},
			checks: func(t *testing.T, experiences []EmployeeExperience, err error) {
				require.NoError(t, err, "ListEmployeeExperience() should not error")
				require.Len(t, experiences, 3, "Should return 3 experience records")
				require.Equal(t, "Junior Developer", experiences[0].JobTitle)
				require.Equal(t, "Mid-level Developer", experiences[1].JobTitle)
				require.Equal(t, "Senior Developer", experiences[2].JobTitle)
			},
		},
		{
			name: "list experiences for employee with no experiences",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return employee.UserID
			},
			checks: func(t *testing.T, experiences []EmployeeExperience, err error) {
				require.NoError(t, err, "ListEmployeeExperience() should not error")
				require.Empty(t, experiences, "Should return empty list for employee with no experiences")
			},
		},
		{
			name: "list experiences with single experience",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)

				_, err := qtx.AddEmployeeExperience(ctx, AddEmployeeExperienceParams{
					EmployeeID:  employee.UserID,
					JobTitle:    "Healthcare Coordinator",
					CompanyName: "Medical Clinic",
					StartDate: pgtype.Date{
						Time:  time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: util.StringPtr("Patient coordination and scheduling"),
				})
				require.NoError(t, err)

				return employee.UserID
			},
			checks: func(t *testing.T, experiences []EmployeeExperience, err error) {
				require.NoError(t, err, "ListEmployeeExperience() should not error")
				require.Len(t, experiences, 1, "Should return 1 experience record")
				require.Equal(t, "Healthcare Coordinator", experiences[0].JobTitle)
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
			experiences, err := qtx.ListEmployeeExperience(ctx, employeeID)
			tt.checks(t, experiences, err)
		})
	}
}

func TestUpdateEmployeeExperience(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateEmployeeExperienceParams
		checks func(t *testing.T, experience EmployeeExperience, err error)
	}{
		{
			name: "successful update with all fields",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeExperienceParams {
				exp := createRandomExperience(ctx, qtx)
				return UpdateEmployeeExperienceParams{
					ID:          exp.ID,
					JobTitle:    util.StringPtr("Senior Architect"),
					CompanyName: util.StringPtr("Fortune 500 Company"),
					StartDate: pgtype.Date{
						Time:  time.Date(2017, 6, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: util.StringPtr("Leading architecture decisions and team management"),
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience, err error) {
				require.NoError(t, err, "UpdateEmployeeExperience() should not error")
				require.Equal(t, "Senior Architect", experience.JobTitle)
				require.Equal(t, "Fortune 500 Company", experience.CompanyName)
				require.NotNil(t, experience.Description)
				require.Equal(t, "Leading architecture decisions and team management", *experience.Description)
				require.Equal(t, time.Date(2017, 6, 1, 0, 0, 0, 0, time.UTC), experience.StartDate.Time)
				require.Equal(t, time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC), experience.EndDate.Time)
			},
		},
		{
			name: "successful update with partial fields (only job title)",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeExperienceParams {
				exp := createRandomExperience(ctx, qtx)
				return UpdateEmployeeExperienceParams{
					ID:          exp.ID,
					JobTitle:    util.StringPtr("Principal Engineer"),
					CompanyName: nil,
					StartDate:   pgtype.Date{Valid: false},
					EndDate:     pgtype.Date{Valid: false},
					Description: nil,
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience, err error) {
				require.NoError(t, err, "UpdateEmployeeExperience() should not error")
				require.Equal(t, "Principal Engineer", experience.JobTitle)
				require.Equal(t, "Tech Corp", experience.CompanyName, "CompanyName should remain unchanged")
			},
		},
		{
			name: "successful update with partial fields (only company name)",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeExperienceParams {
				exp := createRandomExperience(ctx, qtx)
				return UpdateEmployeeExperienceParams{
					ID:          exp.ID,
					JobTitle:    nil,
					CompanyName: util.StringPtr("Updated Corporation"),
					StartDate:   pgtype.Date{Valid: false},
					EndDate:     pgtype.Date{Valid: false},
					Description: nil,
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience, err error) {
				require.NoError(t, err, "UpdateEmployeeExperience() should not error")
				require.Equal(t, "Updated Corporation", experience.CompanyName)
				require.Equal(t, "Software Engineer", experience.JobTitle, "JobTitle should remain unchanged")
			},
		},
		{
			name: "successful update with partial fields (only dates)",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeExperienceParams {
				exp := createRandomExperience(ctx, qtx)
				return UpdateEmployeeExperienceParams{
					ID:          exp.ID,
					JobTitle:    nil,
					CompanyName: nil,
					StartDate: pgtype.Date{
						Time:  time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: nil,
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience, err error) {
				require.NoError(t, err, "UpdateEmployeeExperience() should not error")
				require.Equal(t, time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC), experience.StartDate.Time)
				require.Equal(t, time.Date(2021, 12, 31, 0, 0, 0, 0, time.UTC), experience.EndDate.Time)
				require.Equal(t, "Software Engineer", experience.JobTitle, "JobTitle should remain unchanged")
			},
		},
		{
			name: "successful update with description only",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeExperienceParams {
				exp := createRandomExperience(ctx, qtx)
				return UpdateEmployeeExperienceParams{
					ID:          exp.ID,
					JobTitle:    nil,
					CompanyName: nil,
					StartDate:   pgtype.Date{Valid: false},
					EndDate:     pgtype.Date{Valid: false},
					Description: util.StringPtr("Updated description with new responsibilities and achievements"),
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience, err error) {
				require.NoError(t, err, "UpdateEmployeeExperience() should not error")
				require.NotNil(t, experience.Description)
				require.Equal(t, "Updated description with new responsibilities and achievements", *experience.Description)
				require.Equal(t, "Software Engineer", experience.JobTitle, "JobTitle should remain unchanged")
			},
		},
		{
			name: "update non-existent experience",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeExperienceParams {
				return UpdateEmployeeExperienceParams{
					ID:          99999,
					JobTitle:    util.StringPtr("Senior Architect"),
					CompanyName: util.StringPtr("Fortune 500 Company"),
					StartDate: pgtype.Date{
						Time:  time.Date(2017, 6, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: util.StringPtr("Leading architecture decisions"),
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience, err error) {
				require.Error(t, err, "UpdateEmployeeExperience() should error for non-existent ID")
			},
		},
		{
			name: "update with empty/nil fields",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeExperienceParams {
				exp := createRandomExperience(ctx, qtx)
				return UpdateEmployeeExperienceParams{
					ID:          exp.ID,
					JobTitle:    nil,
					CompanyName: nil,
					StartDate:   pgtype.Date{Valid: false},
					EndDate:     pgtype.Date{Valid: false},
					Description: nil,
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience, err error) {
				require.NoError(t, err, "UpdateEmployeeExperience() should not error")
				require.Equal(t, "Software Engineer", experience.JobTitle)
				require.Equal(t, "Tech Corp", experience.CompanyName)
			},
		},
		{
			name: "successful update removing description",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeExperienceParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				exp, err := qtx.AddEmployeeExperience(ctx, AddEmployeeExperienceParams{
					EmployeeID:  employee.UserID,
					JobTitle:    "Developer",
					CompanyName: "Company",
					StartDate: pgtype.Date{
						Time:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Description: util.StringPtr("Old description"),
				})
				require.NoError(t, err)

				return UpdateEmployeeExperienceParams{
					ID:          exp.ID,
					JobTitle:    nil,
					CompanyName: nil,
					StartDate:   pgtype.Date{Valid: false},
					EndDate:     pgtype.Date{Valid: false},
					Description: util.StringPtr(""),
				}
			},
			checks: func(t *testing.T, experience EmployeeExperience, err error) {
				require.NoError(t, err, "UpdateEmployeeExperience() should not error")
				require.NotNil(t, experience.Description)
				require.Equal(t, "", *experience.Description)
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
			experience, err := qtx.UpdateEmployeeExperience(ctx, params)
			tt.checks(t, experience, err)
		})
	}
}

// Random data generators for tests

func createRandomExperience(ctx context.Context, qtx *Queries) EmployeeExperience {
	employee := createRandomEmployeeProfile(ctx, qtx)

	exp, err := qtx.AddEmployeeExperience(ctx, AddEmployeeExperienceParams{
		EmployeeID:  employee.UserID,
		JobTitle:    "Software Engineer",
		CompanyName: "Tech Corp",
		StartDate: pgtype.Date{
			Time:  time.Date(2018, 1, 15, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		EndDate: pgtype.Date{
			Time:  time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
		Description: nil,
	})
	if err != nil {
		panic("failed to create random experience: " + err.Error())
	}
	return exp
}
