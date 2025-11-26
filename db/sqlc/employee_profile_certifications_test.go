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

func TestAddEmployeeCertification(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) AddEmployeeCertificationParams
		checks func(t *testing.T, cert Certification)
	}{
		{
			name: "successful creation with minimal required fields",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeCertificationParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEmployeeCertificationParams{
					EmployeeID: employee.ID,
					Name:       "AWS Certified Solutions Architect",
					IssuedBy:   "Amazon Web Services",
					DateIssued: pgtype.Date{
						Time:  time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, cert Certification) {
				require.NotZero(t, cert.ID, "Certification ID should be set")
				require.Equal(t, "AWS Certified Solutions Architect", cert.Name)
				require.Equal(t, "Amazon Web Services", cert.IssuedBy)
				require.True(t, cert.DateIssued.Valid, "DateIssued should be valid")
				require.Equal(t, time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC), cert.DateIssued.Time)
				require.True(t, cert.CreatedAt.Valid, "CreatedAt should be valid")
			},
		},
		{
			name: "successful creation with different certification",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeCertificationParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEmployeeCertificationParams{
					EmployeeID: employee.ID,
					Name:       "Certified Nursing Assistant",
					IssuedBy:   "National Healthcare Board",
					DateIssued: pgtype.Date{
						Time:  time.Date(2022, 3, 20, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, cert Certification) {
				require.Equal(t, "Certified Nursing Assistant", cert.Name)
				require.Equal(t, "National Healthcare Board", cert.IssuedBy)
			},
		},
		{
			name: "successful creation with special characters in name",
			setup: func(ctx context.Context, qtx *Queries) AddEmployeeCertificationParams {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return AddEmployeeCertificationParams{
					EmployeeID: employee.ID,
					Name:       "Microsoft Certified: Azure Solutions Architect Expert",
					IssuedBy:   "Microsoft Learning",
					DateIssued: pgtype.Date{
						Time:  time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, cert Certification) {
				require.Equal(t, "Microsoft Certified: Azure Solutions Architect Expert", cert.Name)
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
			cert, err := qtx.AddEmployeeCertification(ctx, params)
			require.NoError(t, err, "AddEmployeeCertification() should not error")

			tt.checks(t, cert)
		})
	}
}

func TestDeleteEmployeeCertification(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) int64
		checks func(t *testing.T, cert Certification, err error)
	}{
		{
			name: "successful deletion of existing certification",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				cert := createRandomCertification(ctx, qtx)
				return cert.ID
			},
			checks: func(t *testing.T, cert Certification, err error) {
				require.NoError(t, err, "DeleteEmployeeCertification() should not error")
				require.NotZero(t, cert.ID, "Deleted certification should have an ID")
				require.Equal(t, "AWS Certified Solutions Architect", cert.Name)
			},
		},
		{
			name: "delete non-existent certification",
			setup: func(ctx context.Context, qtx *Queries) int64 {
				return 99999 // Non-existent ID
			},
			checks: func(t *testing.T, cert Certification, err error) {
				require.Error(t, err, "DeleteEmployeeCertification() should error for non-existent ID")
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
			cert, err := qtx.DeleteEmployeeCertification(ctx, id)
			tt.checks(t, cert, err)
		})
	}
}

func TestListEmployeeCertifications(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) uuid.UUID
		checks func(t *testing.T, certs []Certification, err error)
	}{
		{
			name: "list certifications for employee with multiple certs",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)

				// Create multiple certifications
				_, err := qtx.AddEmployeeCertification(ctx, AddEmployeeCertificationParams{
					EmployeeID: employee.ID,
					Name:       "AWS Certified Solutions Architect",
					IssuedBy:   "Amazon Web Services",
					DateIssued: pgtype.Date{
						Time:  time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				})
				require.NoError(t, err)

				_, err = qtx.AddEmployeeCertification(ctx, AddEmployeeCertificationParams{
					EmployeeID: employee.ID,
					Name:       "Kubernetes Administrator",
					IssuedBy:   "Cloud Native Computing Foundation",
					DateIssued: pgtype.Date{
						Time:  time.Date(2023, 9, 20, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				})
				require.NoError(t, err)

				return employee.ID
			},
			checks: func(t *testing.T, certs []Certification, err error) {
				require.NoError(t, err, "ListEmployeeCertifications() should not error")
				require.Len(t, certs, 2, "Should return 2 certifications")
				require.Equal(t, "AWS Certified Solutions Architect", certs[0].Name)
				require.Equal(t, "Kubernetes Administrator", certs[1].Name)
			},
		},
		{
			name: "list certifications for employee with no certifications",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)
				return employee.ID
			},
			checks: func(t *testing.T, certs []Certification, err error) {
				require.NoError(t, err, "ListEmployeeCertifications() should not error")
				require.Empty(t, certs, "Should return empty list for employee with no certifications")
			},
		},
		{
			name: "list certifications with single certification",
			setup: func(ctx context.Context, qtx *Queries) uuid.UUID {
				employee := createRandomEmployeeProfile(ctx, qtx)

				_, err := qtx.AddEmployeeCertification(ctx, AddEmployeeCertificationParams{
					EmployeeID: employee.ID,
					Name:       "Certified Nursing Assistant",
					IssuedBy:   "National Healthcare Board",
					DateIssued: pgtype.Date{
						Time:  time.Date(2022, 3, 20, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				})
				require.NoError(t, err)

				return employee.ID
			},
			checks: func(t *testing.T, certs []Certification, err error) {
				require.NoError(t, err, "ListEmployeeCertifications() should not error")
				require.Len(t, certs, 1, "Should return 1 certification")
				require.Equal(t, "Certified Nursing Assistant", certs[0].Name)
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
			certs, err := qtx.ListEmployeeCertifications(ctx, employeeID)
			tt.checks(t, certs, err)
		})
	}
}

func TestUpdateEmployeeCertification(t *testing.T) {
	tests := []struct {
		name   string
		setup  func(ctx context.Context, qtx *Queries) UpdateEmployeeCertificationParams
		checks func(t *testing.T, cert Certification, err error)
	}{
		{
			name: "successful update with new name and issued_by",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeCertificationParams {
				cert := createRandomCertification(ctx, qtx)
				return UpdateEmployeeCertificationParams{
					ID:       cert.ID,
					Name:     util.StringPtr("AWS Certified Solutions Architect Professional"),
					IssuedBy: util.StringPtr("Amazon Web Services Professional"),
					DateIssued: pgtype.Date{
						Time:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, cert Certification, err error) {
				require.NoError(t, err, "UpdateEmployeeCertification() should not error")
				require.Equal(t, "AWS Certified Solutions Architect Professional", cert.Name)
				require.Equal(t, "Amazon Web Services Professional", cert.IssuedBy)
				require.Equal(t, time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC), cert.DateIssued.Time)
			},
		},
		{
			name: "successful update with partial fields (only name)",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeCertificationParams {
				cert := createRandomCertification(ctx, qtx)
				return UpdateEmployeeCertificationParams{
					ID:         cert.ID,
					Name:       util.StringPtr("Updated Certification Name"),
					IssuedBy:   nil,
					DateIssued: pgtype.Date{Valid: false},
				}
			},
			checks: func(t *testing.T, cert Certification, err error) {
				require.NoError(t, err, "UpdateEmployeeCertification() should not error")
				require.Equal(t, "Updated Certification Name", cert.Name)
				require.Equal(t, "Amazon Web Services", cert.IssuedBy, "IssuedBy should remain unchanged")
			},
		},
		{
			name: "successful update with new date only",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeCertificationParams {
				cert := createRandomCertification(ctx, qtx)
				return UpdateEmployeeCertificationParams{
					ID:         cert.ID,
					Name:       nil,
					IssuedBy:   nil,
					DateIssued: pgtype.Date{Time: time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), Valid: true},
				}
			},
			checks: func(t *testing.T, cert Certification, err error) {
				require.NoError(t, err, "UpdateEmployeeCertification() should not error")
				require.Equal(t, "AWS Certified Solutions Architect", cert.Name, "Name should remain unchanged")
				require.Equal(t, time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC), cert.DateIssued.Time)
			},
		},
		{
			name: "update non-existent certification",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeCertificationParams {
				return UpdateEmployeeCertificationParams{
					ID:       99999,
					Name:     util.StringPtr("Updated Name"),
					IssuedBy: util.StringPtr("Updated Issuer"),
					DateIssued: pgtype.Date{
						Time:  time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
			},
			checks: func(t *testing.T, cert Certification, err error) {
				require.Error(t, err, "UpdateEmployeeCertification() should error for non-existent ID")
			},
		},
		{
			name: "update with empty/nil date",
			setup: func(ctx context.Context, qtx *Queries) UpdateEmployeeCertificationParams {
				cert := createRandomCertification(ctx, qtx)
				return UpdateEmployeeCertificationParams{
					ID:         cert.ID,
					Name:       nil,
					IssuedBy:   nil,
					DateIssued: pgtype.Date{Valid: false},
				}
			},
			checks: func(t *testing.T, cert Certification, err error) {
				require.NoError(t, err, "UpdateEmployeeCertification() should not error")
				require.Equal(t, "AWS Certified Solutions Architect", cert.Name)
				require.Equal(t, "Amazon Web Services", cert.IssuedBy)
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
			cert, err := qtx.UpdateEmployeeCertification(ctx, params)
			tt.checks(t, cert, err)
		})
	}
}

// Random data generators for tests

func createRandomCertification(ctx context.Context, qtx *Queries) Certification {
	employee := createRandomEmployeeProfile(ctx, qtx)

	cert, err := qtx.AddEmployeeCertification(ctx, AddEmployeeCertificationParams{
		EmployeeID: employee.ID,
		Name:       "AWS Certified Solutions Architect",
		IssuedBy:   "Amazon Web Services",
		DateIssued: pgtype.Date{
			Time:  time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
			Valid: true,
		},
	})
	if err != nil {
		panic("failed to create random certification: " + err.Error())
	}
	return cert
}
