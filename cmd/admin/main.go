package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	db "maicare_go/db/sqlc"
	"maicare_go/internal/ctxkeys"
	"maicare_go/pkg/password"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func seedOrganisations(ctx context.Context, store *db.Store) (*db.Organisation, error) {
	orgs, err := store.ListOrganisations(ctx)
	if err == nil && len(orgs) > 0 {
		return &db.Organisation{
			ID:          orgs[0].ID,
			Name:        orgs[0].Name,
			Street:      orgs[0].Street,
			HouseNumber: orgs[0].HouseNumber,
			PostalCode:  orgs[0].PostalCode,
			City:        orgs[0].City,
			PhoneNumber: orgs[0].PhoneNumber,
			Email:       orgs[0].Email,
			KvkNumber:   orgs[0].KvkNumber,
			BtwNumber:   orgs[0].BtwNumber,
			CreatedAt:   orgs[0].CreatedAt,
			UpdatedAt:   orgs[0].UpdatedAt,
		}, nil
	}

	phone := gofakeit.Phone()
	email := gofakeit.Email()
	kvk := fmt.Sprintf("%08d", 10000000)
	btw := fmt.Sprintf("NL%09dB01", 100000000)

	org, err := store.CreateOrganisation(ctx, db.CreateOrganisationParams{
		Name:        gofakeit.Company() + " Care Organization",
		Street:      gofakeit.Street(),
		HouseNumber: fmt.Sprintf("%d", gofakeit.Number(1, 300)),
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
	locs, err := store.ListLocations(ctx, organisation.ID)
	if err == nil && len(locs) > 0 {
		return &db.Location{
			ID:                  locs[0].ID,
			OrganisationID:      locs[0].OrganisationID,
			Name:                locs[0].Name,
			Street:              locs[0].Street,
			HouseNumber:         locs[0].HouseNumber,
			HouseNumberAddition: locs[0].HouseNumberAddition,
			PostalCode:          locs[0].PostalCode,
			City:                locs[0].City,
			Capacity:            locs[0].Capacity,
			CreatedAt:           locs[0].CreatedAt,
			UpdatedAt:           locs[0].UpdatedAt,
		}, nil
	}

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
		OrganisationID:      organisation.ID,
		Name:                locationName,
		Street:              gofakeit.Street(),
		HouseNumber:         fmt.Sprintf("%d", gofakeit.Number(1, 300)),
		HouseNumberAddition: nil,
		PostalCode:          gofakeit.Zip(),
		City:                gofakeit.City(),
		Capacity:            &capacity,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return &location, nil
}

func createCustomUser(ctx context.Context, pool *pgxpool.Pool, store *db.Store, email string, plainPassword string) (*db.CustomUser, error) {
	var existingID uuid.UUID
	err := pool.QueryRow(ctx, "SELECT id FROM custom_user WHERE email = $1", email).Scan(&existingID)
	if err == nil {
		hashedPassword, hashErr := password.HashPassword(plainPassword)
		if hashErr == nil {
			_ = store.UpdatePassword(ctx, db.UpdatePasswordParams{
				ID:       existingID,
				Password: hashedPassword,
			})
		}
		return &db.CustomUser{
			ID:       existingID,
			Email:    email,
			IsActive: true,
		}, nil
	}

	// Hash a default password for all users
	hashedPassword, err := password.HashPassword(plainPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := store.CreateUser(ctx, db.CreateUserParams{
		Email:    email,
		Password: hashedPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func seedEmployeeProfiles(ctx context.Context, pool *pgxpool.Pool, store *db.Store, user *db.CustomUser, location *db.Location) (*db.EmployeeProfile, error) {
	var existingEmpID uuid.UUID
	err := pool.QueryRow(ctx, "SELECT id FROM employee_profile WHERE user_id = $1", user.ID).Scan(&existingEmpID)
	if err == nil {
		return &db.EmployeeProfile{
			ID:     existingEmpID,
			UserID: user.ID,
		}, nil
	}

	genders := []db.GenderEnum{
		db.GenderEnumMale,
		db.GenderEnumFemale,
		db.GenderEnumOther,
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
	privateEmail := gofakeit.Email()
	authPhone := gofakeit.Phone()
	privatePhone := gofakeit.Phone()
	homePhone := gofakeit.Phone()

	// Random dates
	dob := time.Now().AddDate(-rand.Intn(35)-25, -rand.Intn(12), -rand.Intn(28))

	dept, err := store.GetDepartmentByName(ctx, department)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			dept, err = store.CreateDepartment(ctx, db.CreateDepartmentParams{
				Name:                     department,
				DepartmentHeadEmployeeID: nil,
			})
		}
		if err != nil {
			return nil, fmt.Errorf("failed to resolve department: %w", err)
		}
	}
	departmentID := dept.ID

	employee, err := store.CreateEmployeeProfile(ctx, db.CreateEmployeeProfileParams{
		UserID:              user.ID,
		FirstName:           gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		Position:            &position,
		DepartmentID:        &departmentID,
		EmployeeNumber:      &employeeNumber,
		PrivateEmailAddress: &privateEmail,
		WorkEmailAddress:    &user.Email,
		WorkPhoneNumber:     &authPhone,
		PrivatePhoneNumber:  &privatePhone,
		DateOfBirth:         pgtype.Date{Time: dob, Valid: true},
		HomeTelephoneNumber: &homePhone,
		Gender:              genders[0],
		LocationID:          &locationID,
		ContractType:        contractTypes[0],
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create employee profile: %w", err)
	}

	return &employee, nil
}

func seedSystemActor(ctx context.Context, pool *pgxpool.Pool, store *db.Store) error {
	systemUserUUIDStr := strings.TrimSpace(os.Getenv("SYSTEM_ACTOR_USER_ID"))
	if systemUserUUIDStr == "" {
		systemUserUUIDStr = "00000000-0000-0000-0000-000000000001"
	}
	systemUserID, err := uuid.Parse(systemUserUUIDStr)
	if err != nil {
		return fmt.Errorf("invalid SYSTEM_ACTOR_USER_ID %q: %w", systemUserUUIDStr, err)
	}

	systemEmpUUIDStr := strings.TrimSpace(os.Getenv("SYSTEM_ACTOR_EMPLOYEE_ID"))
	if systemEmpUUIDStr == "" {
		systemEmpUUIDStr = "00000000-0000-0000-0000-000000000002"
	}
	systemEmployeeID, err := uuid.Parse(systemEmpUUIDStr)
	if err != nil {
		return fmt.Errorf("invalid SYSTEM_ACTOR_EMPLOYEE_ID %q: %w", systemEmpUUIDStr, err)
	}

	systemEmail := strings.TrimSpace(os.Getenv("SYSTEM_ACTOR_EMAIL"))
	if systemEmail == "" {
		systemEmail = "system@maicare.internal"
	}

	systemPassword, err := password.HashPassword(uuid.NewString() + uuid.NewString())
	if err != nil {
		return fmt.Errorf("hash system password: %w", err)
	}

	fmt.Printf("Seeding System Actor (User: %s, Employee: %s)...\n", systemUserID, systemEmployeeID)

	// 1. Upsert custom_user
	_, err = pool.Exec(ctx, `
		INSERT INTO custom_user (id, email, password, is_active)
		VALUES ($1, $2, $3, true)
		ON CONFLICT (id) DO UPDATE
		SET is_active = true,
		    email = EXCLUDED.email
	`, systemUserID, systemEmail, systemPassword)
	if err != nil {
		return fmt.Errorf("upsert system custom_user: %w", err)
	}

	// 2. Upsert employee_profile
	_, err = pool.Exec(ctx, `
		INSERT INTO employee_profile (
			id, user_id, first_name, last_name, bsn, street, house_number,
			postal_code, city, position, employee_number, work_email_address,
			gender, is_archived, out_of_service, contract_type
		) VALUES (
			$1, $2, 'System', 'Actor', '000000000', 'System Street', '1',
			'0000AA', 'System City', 'System Service Actor', 'SYS00001', $3,
			'other', false, false, 'none'
		)
		ON CONFLICT (id) DO UPDATE
		SET user_id = EXCLUDED.user_id,
		    is_archived = false,
		    out_of_service = false
	`, systemEmployeeID, systemUserID, systemEmail)
	if err != nil {
		return fmt.Errorf("upsert system employee_profile: %w", err)
	}

	// 3. Grant admin role
	grantPermissions(ctx, store, systemUserID)

	// 4. Validate actor identity check
	systemActor := ctxkeys.ActorIdentity{UserID: systemUserID, EmployeeID: systemEmployeeID}
	if err := store.ValidateActorIdentity(ctx, systemActor); err != nil {
		return fmt.Errorf("validate system actor identity: %w", err)
	}

	fmt.Println("System Actor Seeded & Validated Successfully")
	return nil
}

func grantPermissions(ctx context.Context, store *db.Store, userID uuid.UUID) {
	roleID, err := store.GetAdminRoleId(ctx)
	if err != nil {
		log.Fatalf("Failed to get admin role ID: %v", err)
	}
	if err := store.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: userID,
		RoleID: roleID,
	}); err != nil {
		log.Fatalf("Failed to assign role to user: %v", err)
	}
}

func loadEnvFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Printf("[admin] warning: cannot read %s: %v", path, err)
		}
		return
	}

	for _, rawLine := range strings.Split(string(content), "\n") {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		if os.Getenv(key) == "" {
			os.Setenv(key, strings.TrimSpace(value))
		}
	}
}

func main() {
	loadEnvFile("app.env")

	dbSourceFlag := flag.String("db", "", "database connection string (defaults to DB_SOURCE then local default)")
	flag.Parse()

	adminEmail := strings.TrimSpace(os.Getenv("ADMIN_EMAIL"))
	adminPassword := strings.TrimSpace(os.Getenv("ADMIN_PASSWORD"))
	if adminEmail == "" || adminPassword == "" {
		log.Fatal("ADMIN_EMAIL and ADMIN_PASSWORD must be set in the environment variables")
	}

	ctx := context.Background()
	dbSource := strings.TrimSpace(*dbSourceFlag)
	if dbSource == "" {
		dbSource = strings.TrimSpace(os.Getenv("DB_SOURCE"))
	}
	if dbSource == "" {
		dbSource = "postgres://maicare:maicare@localhost:5432/maicare?sslmode=disable"
	}

	connPool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	defer connPool.Close()

	store := db.NewStore(connPool)

	if err := seedSystemActor(ctx, connPool, store); err != nil {
		log.Fatalf("Failed to seed system actor: %v", err)
	}

	fmt.Println("Seeding Admin User...")

	user, err := createCustomUser(ctx, connPool, store, adminEmail, adminPassword)
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

	_, err = seedEmployeeProfiles(ctx, connPool, store, user, location)
	if err != nil {
		log.Fatalf("Failed to seed employee profile: %v", err)
	}
	grantPermissions(ctx, store, user.ID)

	fmt.Println("Admin User Seeded Successfully")

}
