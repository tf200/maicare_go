package invoice

import (
	"context"
	"encoding/json"
	"fmt"
	db "maicare_go/db/sqlc"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomLocation(t *testing.T) *db.Location {
	arg := db.CreateLocationParams{
		Name:     util.RandomString(5),
		Address:  util.RandomString(8),
		Capacity: util.Int32Ptr(25),
	}

	location, err := testStore.CreateLocation(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, location)

	// Check if the returned location matches the input
	require.Equal(t, arg.Name, location.Name)
	require.Equal(t, arg.Address, location.Address)
	require.Equal(t, arg.Capacity, location.Capacity)

	// Verify ID is generated
	require.NotZero(t, location.ID)
	return &location
}
func CreateRandomUser(t *testing.T) *db.CustomUser {
	hashedPassword, err := util.HashPassword("t2aha000")
	require.NoError(t, err)
	// arg := CreateUserParams{
	// 	Password:       hashedPassword,
	// 	Username:       util.StringPtr(util.RandomString(5)),
	// 	Email:          util.RandomEmail(),
	// 	FirstName:      util.RandomString(5),
	// 	LastName:       util.RandomString(5),
	// 	IsSuperuser:    true,
	// 	IsStaff:        true,
	// 	IsActive:       true,
	// 	ProfilePicture: util.StringPtr(util.GetRandomImageURL()),
	// 	PhoneNumber:    util.IntPtr(456),
	// }
	arg := db.CreateUserParams{
		Password:       hashedPassword,
		Email:          util.RandomEmail(),
		IsActive:       true,
		ProfilePicture: util.StringPtr(util.GetRandomImageURL()),
		RoleID:         1,
	}

	user, err := testStore.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.IsActive, user.IsActive)
	require.Equal(t, arg.ProfilePicture, user.ProfilePicture)

	// Verify timestamps and ID are set
	require.NotZero(t, user.ID)

	return &user
}

func createRandomEmployee(t *testing.T) (db.EmployeeProfile, *db.CustomUser) {
	// Create prerequisite records
	location := createRandomLocation(t)
	user := CreateRandomUser(t)

	arg := db.CreateEmployeeProfileParams{
		UserID:                    user.ID,
		FirstName:                 faker.FirstName(),
		LastName:                  faker.LastName(),
		Position:                  util.StringPtr(util.RandomString(5)),
		Department:                util.StringPtr(util.RandomString(5)),
		EmployeeNumber:            util.StringPtr(util.RandomString(5)),
		EmploymentNumber:          util.StringPtr(util.RandomString(5)),
		PrivateEmailAddress:       util.StringPtr(util.RandomString(5)),
		Email:                     util.RandomEmail(),
		AuthenticationPhoneNumber: util.StringPtr(util.RandomString(5)),
		PrivatePhoneNumber:        util.StringPtr(util.RandomString(5)),
		WorkPhoneNumber:           util.StringPtr(util.RandomString(5)),
		DateOfBirth:               pgtype.Date{Time: time.Now(), Valid: true},
		HomeTelephoneNumber:       util.StringPtr(util.RandomString(5)),
		IsSubcontractor:           util.BoolPtr(true),
		Gender:                    util.StringPtr("male"),
		LocationID:                util.IntPtr(location.ID),
		HasBorrowed:               false,
		OutOfService:              util.BoolPtr(util.RandomBool()),
		IsArchived:                util.RandomBool(),
	}

	employee, err := testStore.CreateEmployeeProfile(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, employee)

	// Verify all fields match
	require.Equal(t, arg.UserID, employee.UserID)
	require.Equal(t, arg.FirstName, employee.FirstName)
	require.Equal(t, arg.LastName, employee.LastName)
	require.Equal(t, arg.Position, employee.Position)
	require.Equal(t, arg.Department, employee.Department)
	require.Equal(t, arg.EmployeeNumber, employee.EmployeeNumber)
	require.Equal(t, arg.EmploymentNumber, employee.EmploymentNumber)
	require.Equal(t, arg.PrivateEmailAddress, employee.PrivateEmailAddress)
	require.Equal(t, arg.Email, employee.Email)
	require.Equal(t, arg.AuthenticationPhoneNumber, employee.AuthenticationPhoneNumber)
	require.Equal(t, arg.PrivatePhoneNumber, employee.PrivatePhoneNumber)
	require.Equal(t, arg.WorkPhoneNumber, employee.WorkPhoneNumber)
	require.Equal(t, arg.DateOfBirth.Time.Format("2006-01-02"), employee.DateOfBirth.Time.Format("2006-01-02"))
	require.Equal(t, arg.HomeTelephoneNumber, employee.HomeTelephoneNumber)
	require.Equal(t, arg.IsSubcontractor, employee.IsSubcontractor)
	require.Equal(t, arg.Gender, employee.Gender)
	require.Equal(t, arg.LocationID, employee.LocationID)
	require.Equal(t, arg.HasBorrowed, employee.HasBorrowed)
	require.Equal(t, arg.OutOfService, employee.OutOfService)
	require.Equal(t, arg.IsArchived, employee.IsArchived)

	require.NotZero(t, employee.ID)
	require.NotZero(t, employee.CreatedAt)

	require.Equal(t, util.IntPtr(location.ID), employee.LocationID)

	return employee, user
}
func createRandomSenders(t *testing.T) db.Sender {
	arg := db.CreateSenderParams{
		Types:        "main_provider",
		Name:         util.RandomString(5),
		Address:      util.StringPtr("test"),
		PostalCode:   util.StringPtr("test"),
		Place:        util.StringPtr("test"),
		Land:         util.StringPtr("test"),
		Kvknumber:    util.StringPtr("test"),
		Btwnumber:    util.StringPtr("test"),
		PhoneNumber:  util.StringPtr("test"),
		ClientNumber: util.StringPtr("test"),
		EmailAddress: util.StringPtr("test"),
		Contacts:     []byte(`[{"name": "Test Contact", "email": "test@example.com", "phone": "1234567890"}]`),
	}

	sender, err := testStore.CreateSender(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, sender)
	require.Equal(t, arg.Types, sender.Types)
	require.Equal(t, arg.Name, sender.Name)
	require.Equal(t, arg.Address, sender.Address)
	require.Equal(t, arg.PostalCode, sender.PostalCode)
	require.Equal(t, arg.Place, sender.Place)
	require.Equal(t, arg.Land, sender.Land)
	require.Equal(t, arg.Kvknumber, sender.Kvknumber)
	require.Equal(t, arg.Btwnumber, sender.Btwnumber)
	require.Equal(t, arg.PhoneNumber, sender.PhoneNumber)
	require.Equal(t, arg.ClientNumber, sender.ClientNumber)
	require.Equal(t, arg.EmailAddress, sender.EmailAddress)
	require.Equal(t, arg.Contacts, sender.Contacts)
	return sender
}

func createRandomClientDetails(t *testing.T) db.ClientDetail {
	location := createRandomLocation(t)
	sender := createRandomSenders(t)

	arg := db.CreateClientDetailsParams{
		FirstName:       faker.FirstName(),
		LastName:        faker.LastName(),
		Email:           faker.Email(),
		PhoneNumber:     util.StringPtr(faker.Phonenumber()),
		DateOfBirth:     pgtype.Date{Time: time.Now().AddDate(-20, 0, 0), Valid: true},
		Identity:        false,
		Bsn:             util.StringPtr(util.RandomString(9)),
		BsnVerifiedBy:   nil, // Assuming employee is created and has an ID
		Source:          util.StringPtr("Test Source"),
		Birthplace:      util.StringPtr("test city"),
		Organisation:    util.StringPtr("test org"),
		Departement:     util.StringPtr("test dep"),
		Gender:          "male", // or "Female" or other values as per your requirements
		Filenumber:      "testfile",
		ProfilePicture:  util.StringPtr(util.GetRandomImageURL()),
		Infix:           util.StringPtr("van"),
		SenderID:        &sender.ID,
		LocationID:      util.IntPtr(location.ID),
		DepartureReason: util.StringPtr("test Reason"),
		DepartureReport: util.StringPtr("test report"),
		Addresses:       []byte("[]"),
		LegalMeasure:    util.StringPtr("test measure"),
	}

	client, err := testStore.CreateClientDetails(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, client)
	require.Equal(t, arg.FirstName, client.FirstName)
	require.Equal(t, arg.LastName, client.LastName)
	require.Equal(t, arg.Email, client.Email)
	require.Equal(t, arg.PhoneNumber, client.PhoneNumber)
	require.Equal(t, arg.Filenumber, client.Filenumber)
	require.Equal(t, arg.ProfilePicture, client.ProfilePicture)
	require.Equal(t, arg.Infix, client.Infix)
	require.Equal(t, arg.SenderID, client.SenderID)
	require.Equal(t, arg.LocationID, client.LocationID)
	require.Equal(t, arg.DepartureReason, client.DepartureReason)
	require.Equal(t, arg.DepartureReport, client.DepartureReport)
	require.Equal(t, arg.Addresses, client.Addresses)
	require.Equal(t, arg.LegalMeasure, client.LegalMeasure)
	return client
}

func createRandomContractType(t *testing.T) db.ContractType {
	// Create a random contract type

	contractType, err := testStore.CreateContractType(context.Background(), "Test Contract Type")
	require.NoError(t, err)
	require.NotEmpty(t, contractType)
	require.NotEmpty(t, contractType.ID)
	require.Equal(t, "Test Contract Type", contractType.Name)
	return contractType
}

func createRandomAppointment(t *testing.T, employeeID *int64) db.ScheduledAppointment {
	arg := db.CreateAppointmentParams{
		CreatorEmployeeID: employeeID,
		StartTime:         pgtype.Timestamp{Time: time.Date(2025, time.August, 6, 12, 0, 0, 0, time.UTC), Valid: true},
		EndTime:           pgtype.Timestamp{Time: time.Date(2025, time.August, 6, 18, 0, 0, 0, time.UTC), Valid: true},
		Location:          util.StringPtr("Test Location"),
		Description:       util.StringPtr("Test Description"),
	}

	appointment, err := testStore.CreateAppointment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, appointment)
	return appointment
}

func TestGenerateInvoice(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)

	// priceFrequency := []string{"minute", "hourly", "daily", "weekly", "monthly"
	// careType := []string{"ambulante", "accommodation"}
	financingAct := []string{"WMO", "ZVW", "WLZ", "JW", "WPG"}
	financingOption := []string{"ZIN", "PGB"}

	contractType := createRandomContractType(t)

	arg1 := db.CreateContractParams{
		TypeID:          &contractType.ID,
		StartDate:       pgtype.Timestamptz{Time: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		EndDate:         pgtype.Timestamptz{Time: time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC), Valid: true},
		ReminderPeriod:  10,
		Vat:             util.Int32Ptr(20),
		Status:          "approved",
		Price:           89,
		PriceTimeUnit:   "daily", // util.RandomEnum(priceFrequency),
		Hours:           nil,
		HoursType:       nil,
		CareName:        "Test Care",
		CareType:        "accommodation",
		ClientID:        client.ID,
		SenderID:        client.SenderID,
		FinancingAct:    util.RandomEnum(financingAct),
		FinancingOption: util.RandomEnum(financingOption),
		AttachmentIds:   []uuid.UUID{},
	}

	contract, err := testStore.CreateContract(context.Background(), arg1)
	require.NoError(t, err)
	require.NotEmpty(t, contract)

	arg2 := db.CreateContractParams{
		TypeID:          &contractType.ID,
		StartDate:       pgtype.Timestamptz{Time: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true},
		EndDate:         pgtype.Timestamptz{Time: time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC), Valid: true},
		ReminderPeriod:  10,
		Vat:             util.Int32Ptr(20),
		Status:          "approved",
		Price:           58,
		PriceTimeUnit:   "hourly", // util.RandomEnum(priceFrequency),
		Hours:           util.Float64Ptr(40),
		HoursType:       util.StringPtr("weekly"),
		CareName:        "Test Care",
		CareType:        "ambulante", // util.RandomEnum(careType),
		ClientID:        client.ID,
		SenderID:        client.SenderID,
		FinancingAct:    util.RandomEnum(financingAct),
		FinancingOption: util.RandomEnum(financingOption),
		AttachmentIds:   []uuid.UUID{},
	}
	contract2, err := testStore.CreateContract(context.Background(), arg2)
	require.NoError(t, err)
	require.NotEmpty(t, contract2)

	appointement := createRandomAppointment(t, &employee.ID)
	err = testStore.BulkAddAppointmentClients(context.Background(), db.BulkAddAppointmentClientsParams{
		AppointmentID: appointement.ID,
		ClientIds:     []int64{client.ID},
	})

	require.NoError(t, err)

	invoiceData := InvoiceParams{
		ClientID:  client.ID,
		StartDate: time.Date(2025, time.August, 1, 0, 0, 0, 0, time.UTC),
		EndDate:   time.Date(2025, time.August, 28, 23, 59, 59, 0, time.UTC),
	}

	invoice, warningCount, err := GenerateInvoice(testStore, invoiceData, context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, invoice)
	fmt.Println("Invoice generated successfully:", invoice)

	invoiceDetailsBytes, err := json.Marshal(invoice.InvoiceDetails)
	require.NoError(t, err)

	insertedInvoice, err := testStore.CreateInvoice(context.Background(), db.CreateInvoiceParams{
		InvoiceNumber:  invoice.InvoiceNumber,
		ClientID:       invoiceData.ClientID,
		DueDate:        pgtype.Date{Time: time.Now().AddDate(0, 0, 30), Valid: true},
		TotalAmount:    invoice.TotalAmount,
		IssueDate:      pgtype.Date{Time: time.Now(), Valid: true},
		ExtraContent:   nil, // Assuming no extra content for simplicity
		InvoiceDetails: invoiceDetailsBytes,
		SenderID:       client.SenderID,
		WarningCount:   int32(warningCount),
	})
	require.NoError(t, err)
	require.NotEmpty(t, insertedInvoice)
}

// To do other methods

func TestCreatePaymentApi(t *testing.T) {

}
