package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	db "maicare_go/db/sqlc"
	"maicare_go/util"
)

const (
	numOrganisations = 5
	numLocations     = 3 // per organisation
	numUsers         = 20
	numEmployees     = 15
)

func main() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	// Load configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Connect to database
	ctx := context.Background()
	connPool, err := pgxpool.New(ctx, config.DbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}
	defer connPool.Close()

	store := db.NewStore(connPool)

	fmt.Println("🌱 Starting database seeding...")

	// Seed organisations
	organisations, err := seedOrganisations(ctx, store)
	if err != nil {
		log.Fatal("Failed to seed organisations:", err)
	}
	fmt.Printf("✅ Created %d organisations\n", len(organisations))

	// Seed locations
	locations, err := seedLocations(ctx, store, organisations)
	if err != nil {
		log.Fatal("Failed to seed locations:", err)
	}
	fmt.Printf("✅ Created %d locations\n", len(locations))

	// Seed users
	users, err := seedUsers(ctx, store)
	if err != nil {
		log.Fatal("Failed to seed users:", err)
	}
	fmt.Printf("✅ Created %d users\n", len(users))

	// Seed employee profiles
	employees, err := seedEmployeeProfiles(ctx, store, users, locations)
	if err != nil {
		log.Fatal("Failed to seed employee profiles:", err)
	}
	fmt.Printf("✅ Created %d employee profiles\n", len(employees))

	fmt.Println("🎉 Database seeding completed successfully!")
}

func seedOrganisations(ctx context.Context, store *db.Store) ([]db.Organisation, error) {
	organisations := make([]db.Organisation, 0, numOrganisations)

	for i := 0; i < numOrganisations; i++ {
		phone := gofakeit.Phone()
		email := gofakeit.Email()
		kvk := fmt.Sprintf("%08d", i+10000000)
		btw := fmt.Sprintf("NL%09dB01", i+100000000)

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

		organisations = append(organisations, org)
	}

	return organisations, nil
}

func seedLocations(ctx context.Context, store *db.Store, organisations []db.Organisation) ([]db.Location, error) {
	locations := make([]db.Location, 0, len(organisations)*numLocations)
	locationTypes := []db.LocationTypeEnum{
		db.LocationTypeEnumCareHome,
		db.LocationTypeEnumOffice,
		db.LocationTypeEnumOther,
	}

	for _, org := range organisations {
		for j := 0; j < numLocations; j++ {
			capacity := int32(gofakeit.Number(10, 50))
			locationType := locationTypes[j%len(locationTypes)]

			locationName := gofakeit.Word()
			if locationType == db.LocationTypeEnumCareHome {
				locationName += " Care Home"
			} else if locationType == db.LocationTypeEnumOffice {
				locationName += " Office"
			} else {
				locationName += " Facility"
			}

			location, err := store.CreateLocation(ctx, db.CreateLocationParams{
				OrganisationID: org.ID,
				Name:           locationName,
				Address:        gofakeit.Street(),
				Capacity:       &capacity,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to create location: %w", err)
			}

			locations = append(locations, location)
		}
	}

	return locations, nil
}

func seedUsers(ctx context.Context, store *db.Store) ([]db.CustomUser, error) {
	users := make([]db.CustomUser, 0, numUsers)

	// Hash a default password for all users
	defaultPassword := "Password123!"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	for i := 0; i < numUsers; i++ {
		email := fmt.Sprintf("user%d_%s", i, gofakeit.Email())

		user, err := store.CreateUser(ctx, db.CreateUserParams{
			Email:    email,
			Password: string(hashedPassword),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func seedEmployeeProfiles(ctx context.Context, store *db.Store, users []db.CustomUser, locations []db.Location) ([]db.EmployeeProfile, error) {
	employees := make([]db.EmployeeProfile, 0, numEmployees)

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

	// Use first numEmployees users
	for i := 0; i < numEmployees && i < len(users); i++ {
		user := users[i]

		// Random location
		location := locations[i%len(locations)]
		locationID := location.ID

		position := positions[i%len(positions)]
		department := departments[i%len(departments)]
		employeeNumber := fmt.Sprintf("EMP%05d", i+1)
		employmentNumber := fmt.Sprintf("EMPL%06d", i+1000)
		privateEmail := gofakeit.Email()
		authPhone := gofakeit.Phone()
		privatePhone := gofakeit.Phone()
		workPhone := gofakeit.Phone()
		homePhone := gofakeit.Phone()
		isSubcontractor := i%5 == 0

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
			Gender:                    genders[i%len(genders)],
			LocationID:                &locationID,
			ContractType:              contractTypes[i%len(contractTypes)],
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create employee profile: %w", err)
		}
		// Add contract details
		var contractType db.EmployeeContractTypeEnum
		if isSubcontractor {
			contractType = db.EmployeeContractTypeEnumZZP
		} else {
			contractType = db.EmployeeContractTypeEnumLoondienst
		}
		startdate := gofakeit.Date()
		_, err = store.AddEmployeeContractDetails(ctx, db.AddEmployeeContractDetailsParams{
			ContractHours:     util.Float64Ptr(float64(gofakeit.Number(20, 40))),
			ContractStartDate: pgtype.Date{Time: startdate, Valid: true},
			ContractEndDate:   pgtype.Date{Time: startdate.AddDate(1, 0, 0), Valid: true},
			ID:                employee.ID,
			ContractType:      db.NullEmployeeContractTypeEnum{EmployeeContractTypeEnum: contractType, Valid: true},
			ContractRate:      util.Float64Ptr(gofakeit.Price(20, 60)),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add contract details for subcontractor: %w", err)
		}

		employees = append(employees, employee)
	}
	return employees, nil
}
