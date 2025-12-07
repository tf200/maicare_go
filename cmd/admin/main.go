package main

import (
	"context"
	"fmt"
	"log"
	db "maicare_go/db/sqlc"
	"maicare_go/util"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func seedOrganisations(ctx context.Context, store *db.Store) (*db.Organisation, error) {

	phone := gofakeit.Phone()
	email := gofakeit.Email()
	kvk := fmt.Sprintf("%08d", 10000000)
	btw := fmt.Sprintf("NL%09dB01", 100000000)

	org, err := store.CreateOrganisation(ctx, db.CreateOrganisationParams{
		Name:        gofakeit.Company() + " Care Organization",
		Address:     gofakeit.Street(),
		PostalCode:  gofakeit.Zip(),
		City:        gofakeit.City(),
		PhoneNumber: &phone,
		Email:       &email,
		KvkNumber:   &kvk,
		BtwNumber:   &btw,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create organisation: %w", err)
	}

	return &org, nil
}

func seedLocations(ctx context.Context, store *db.Store, organisation *db.Organisation) (*db.Location, error) {
	locationTypes := []db.LocationTypeEnum{
		db.LocationTypeEnumCareHome,
		db.LocationTypeEnumOffice,
		db.LocationTypeEnumOther,
	}

	capacity := int32(gofakeit.Number(10, 50))
	locationType := locationTypes[0]

	locationName := gofakeit.Word()
	switch locationType {
	case db.LocationTypeEnumCareHome:
		locationName += " Care Home"
	case db.LocationTypeEnumOffice:
		locationName += " Office"
	default:
		locationName += " Facility"
	}

	location, err := store.CreateLocation(ctx, db.CreateLocationParams{
		OrganisationID: organisation.ID,
		Name:           locationName,
		Address:        gofakeit.Street(),
		Capacity:       &capacity,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return &location, nil
}

func createCustomUser(ctx context.Context, store *db.Store, email string, password string) (*db.CustomUser, error) {

	// Hash a default password for all users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := store.CreateUser(ctx, db.CreateUserParams{
		Email:    email,
		Password: string(hashedPassword),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func seedEmployeeProfiles(ctx context.Context, store *db.Store, user *db.CustomUser, location *db.Location) (*db.EmployeeProfile, error) {

	genders := []db.EmployeeGenderEnum{
		db.EmployeeGenderEnumMale,
		db.EmployeeGenderEnumFemale,
		db.EmployeeGenderEnumNotSpecified,
	}

	contractTypes := []db.EmployeeContractTypeEnum{
		db.EmployeeContractTypeEnumLoondienst,
		db.EmployeeContractTypeEnumZZP,
		db.EmployeeContractTypeEnumNone,
	}

	positions := []string{
		"Care Worker",
		"Senior Care Worker",
		"Team Leader",
		"Case Manager",
		"Therapist",
		"Nurse",
		"Administrative Staff",
		"Coordinator",
	}

	departments := []string{
		"Youth Care",
		"Mental Health",
		"Day Care",
		"Residential Care",
		"Administration",
		"Management",
	}

	// Random location
	locationID := location.ID

	position := positions[0]
	department := departments[0]
	employeeNumber := fmt.Sprintf("EMP%05d", 1)
	employmentNumber := fmt.Sprintf("EMPL%06d", 1000)
	privateEmail := gofakeit.Email()
	authPhone := gofakeit.Phone()
	privatePhone := gofakeit.Phone()
	workPhone := gofakeit.Phone()
	homePhone := gofakeit.Phone()
	isSubcontractor := false

	// Random dates
	dob := time.Now().AddDate(-rand.Intn(35)-25, -rand.Intn(12), -rand.Intn(28))

	employee, err := store.CreateEmployeeProfile(ctx, db.CreateEmployeeProfileParams{
		UserID:                    user.ID,
		FirstName:                 gofakeit.FirstName(),
		LastName:                  gofakeit.LastName(),
		Position:                  &position,
		Department:                &department,
		EmployeeNumber:            &employeeNumber,
		EmploymentNumber:          &employmentNumber,
		PrivateEmailAddress:       &privateEmail,
		Email:                     user.Email,
		AuthenticationPhoneNumber: &authPhone,
		PrivatePhoneNumber:        &privatePhone,
		WorkPhoneNumber:           &workPhone,
		DateOfBirth:               pgtype.Date{Time: dob, Valid: true},
		HomeTelephoneNumber:       &homePhone,
		IsSubcontractor:           &isSubcontractor,
		Gender:                    genders[0],
		LocationID:                &locationID,
		ContractType:              contractTypes[0],
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create employee profile: %w", err)
	}

	return &employee, nil
}

func grantPermissions(ctx context.Context, store *db.Store, userID uuid.UUID) {
	roleID, err := store.GetAdminRoleId(ctx)
	if err != nil {
		log.Fatalf("Failed to get admin role ID: %v", err)
	}
	store.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	})
	store.GrantRolePermissionsToUser(ctx, db.GrantRolePermissionsToUserParams{
		UserID: userID,
		RoleID: roleID,
	})
}

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}
	if config.AdminEmail == "" || config.AdminPassword == "" {
		log.Fatal("ADMIN_EMAIL and ADMIN_PASSWORD must be set in the environment variables")
	}
	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, config.DbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	defer connPool.Close()

	store := db.NewStore(connPool)
	fmt.Println("Seeding Admin User...")

	user, err := createCustomUser(ctx, store, config.AdminEmail, config.AdminPassword)
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	org, err := seedOrganisations(ctx, store)
	if err != nil {
		log.Fatalf("Failed to seed organisation: %v", err)
	}

	location, err := seedLocations(ctx, store, org)
	if err != nil {
		log.Fatalf("Failed to seed location: %v", err)
	}

	_, err = seedEmployeeProfiles(ctx, store, user, location)
	if err != nil {
		log.Fatalf("Failed to seed employee profile: %v", err)
	}
	grantPermissions(ctx, store, user.ID)

	fmt.Println("Admin User Seeded Successfully")

}
