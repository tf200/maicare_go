package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	db "maicare_go/db/sqlc"
	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

func ensureAdminAccountOnStartup(ctx context.Context, store *db.Store, config util.Config) error {
	adminEmail := strings.TrimSpace(config.AdminEmail)
	adminPassword := strings.TrimSpace(config.AdminPassword)

	if adminEmail == "" || adminPassword == "" {
		log.Println("ADMIN_EMAIL or ADMIN_PASSWORD is empty; skipping admin bootstrap")
		return nil
	}

	userID, err := findUserIDByEmail(ctx, store, adminEmail)
	if err != nil {
		return err
	}

	createdAny := false
	if userID == uuid.Nil {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
		if hashErr != nil {
			return fmt.Errorf("failed to hash admin password: %w", hashErr)
		}

		user, createErr := store.CreateUser(ctx, db.CreateUserParams{
			Email:    adminEmail,
			Password: string(hashedPassword),
			IsActive: true,
		})
		if createErr != nil {
			if isConflictError(createErr) {
				userID, err = findUserIDByEmail(ctx, store, adminEmail)
				if err != nil {
					return err
				}
				if userID == uuid.Nil {
					return fmt.Errorf("admin user conflict detected but user lookup failed: %w", createErr)
				}
			} else {
				return fmt.Errorf("failed to create admin user: %w", createErr)
			}
		}

		if userID == uuid.Nil {
			userID = user.ID
			createdAny = true
		}
	}

	_, err = store.GetEmployeeProfileByUserID(ctx, userID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to check admin employee profile: %w", err)
		}

		locationID, locationErr := ensureBootstrapLocation(ctx, store)
		if locationErr != nil {
			return locationErr
		}

		if createProfileErr := createAdminEmployeeProfile(ctx, store, userID, adminEmail, locationID); createProfileErr != nil {
			if !isConflictError(createProfileErr) {
				return createProfileErr
			}
		}

		createdAny = true
	}

	if err := ensureAdminRoleAndPermissions(ctx, store, userID); err != nil {
		return err
	}

	if createdAny {
		log.Printf("Admin bootstrap ready for %s", adminEmail)
	} else {
		log.Printf("Admin account already present for %s", adminEmail)
	}

	return nil
}

func findUserIDByEmail(ctx context.Context, store *db.Store, email string) (uuid.UUID, error) {
	var userID uuid.UUID
	err := store.ConnPool.QueryRow(ctx, "SELECT id FROM custom_user WHERE email = $1 LIMIT 1", email).Scan(&userID)
	if err == nil {
		return userID, nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, nil
	}

	return uuid.Nil, fmt.Errorf("failed to lookup admin user: %w", err)
}

func ensureBootstrapLocation(ctx context.Context, store *db.Store) (*uuid.UUID, error) {
	locations, err := store.ListAllLocations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list locations for admin bootstrap: %w", err)
	}
	if len(locations) > 0 {
		locationID := locations[0].ID
		return &locationID, nil
	}

	org, err := store.CreateOrganisation(ctx, db.CreateOrganisationParams{
		Name:        "MaiCare Default Organisation",
		Street:      "Default Street",
		HouseNumber: "1",
		PostalCode:  "0000AA",
		City:        "Default City",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create default organisation for admin bootstrap: %w", err)
	}

	location, err := store.CreateLocation(ctx, db.CreateLocationParams{
		OrganisationID: org.ID,
		Name:           "Main Location",
		Street:         "Default Street",
		HouseNumber:    "1",
		PostalCode:     "0000AA",
		City:           "Default City",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create default location for admin bootstrap: %w", err)
	}

	return &location.ID, nil
}

func createAdminEmployeeProfile(
	ctx context.Context,
	store *db.Store,
	userID uuid.UUID,
	adminEmail string,
	locationID *uuid.UUID,
) error {
	_, err := store.CreateEmployeeProfile(ctx, db.CreateEmployeeProfileParams{
		UserID:            userID,
		FirstName:         "Admin",
		LastName:          "User",
		Bsn:               "000000000",
		Street:            "Default Street",
		HouseNumber:       "1",
		PostalCode:        "0000AA",
		City:              "Default City",
		WorkEmailAddress:  &adminEmail,
		DateOfBirth:       pgtype.Date{Valid: false},
		Gender:            db.GenderEnumOther,
		LocationID:        locationID,
		ContractEndDate:   pgtype.Date{Valid: false},
		ContractStartDate: pgtype.Date{Valid: false},
		ContractType:      db.EmployeeContractTypeEnumNone,
	})
	if err != nil {
		return fmt.Errorf("failed to create admin employee profile: %w", err)
	}

	return nil
}

func ensureAdminRoleAndPermissions(ctx context.Context, store *db.Store, userID uuid.UUID) error {
	adminRoleID, err := store.GetAdminRoleId(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch admin role id: %w", err)
	}

	if err := store.AssignRoleToUser(ctx, db.AssignRoleToUserParams{UserID: userID, RoleID: adminRoleID}); err != nil {
		return fmt.Errorf("failed to assign admin role: %w", err)
	}

	rolePerms, err := store.ListAllRolePermissions(ctx, adminRoleID)
	if err != nil {
		return fmt.Errorf("failed to list admin role permissions: %w", err)
	}

	permissionIDs := make([]uuid.UUID, 0, len(rolePerms))
	for _, perm := range rolePerms {
		permissionIDs = append(permissionIDs, perm.PermissionID)
	}

	if len(permissionIDs) == 0 {
		return nil
	}

	if err := store.GrantUserPermissions(ctx, db.GrantUserPermissionsParams{
		UserID:        userID,
		PermissionIds: permissionIDs,
	}); err != nil {
		return fmt.Errorf("failed to grant admin permissions: %w", err)
	}

	return nil
}

func isConflictError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return strings.HasPrefix(pgErr.Code, "23")
	}

	return false
}
