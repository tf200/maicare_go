package db

import (
	"context"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateEmployeeWithAccountTx(t *testing.T) {
	store := NewStore(testDB)
	location := CreateRandomLocation(t)

	// Create params for the user

	// Create params for the employee
	employeeArg := CreateEmployeeProfileParams{
		FirstName:                 util.RandomString(5),
		LastName:                  util.RandomString(5),
		Position:                  util.StringPtr(util.RandomString(5)),
		Department:                util.StringPtr(util.RandomString(5)),
		EmployeeNumber:            util.StringPtr(util.RandomString(5)),
		EmploymentNumber:          util.StringPtr(util.RandomString(5)),
		PrivateEmailAddress:       util.StringPtr(util.RandomString(5)),
		Email:                     util.RandomEmail(),
		AuthenticationPhoneNumber: util.StringPtr(util.RandomString(5)),
		PrivatePhoneNumber:        util.StringPtr(util.RandomString(5)),
		WorkPhoneNumber:           util.StringPtr(util.RandomString(5)),
		DateOfBirth: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		HomeTelephoneNumber: util.StringPtr(util.RandomString(5)),
		IsSubcontractor:     util.BoolPtr(false),
		Gender:              util.StringPtr(util.RandomString(5)),
		LocationID:          util.IntPtr(location.ID),
		ContractType:        util.StringPtr("loondienst"),
	}
	userArg := CreateUserParams{
		Email:    employeeArg.Email,
		Password: "password123",
		IsActive: true,
	}

	// Create transaction params
	arg := CreateEmployeeWithAccountTxParams{
		CreateUserParams:     userArg,
		CreateEmployeeParams: employeeArg,
	}

	// Execute transaction
	result, err := store.CreateEmployeeWithAccountTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	// Check user account was created properly
	require.NotEmpty(t, result.User)
	require.Equal(t, arg.CreateUserParams.Email, result.User.Email)
	require.Equal(t, arg.CreateUserParams.IsActive, result.User.IsActive)

	// Check employee profile was created properly
	require.NotEmpty(t, result.Employee)
	require.Equal(t, result.User.ID, result.Employee.UserID)
	require.Equal(t, arg.CreateEmployeeParams.FirstName, result.Employee.FirstName)
	require.Equal(t, arg.CreateEmployeeParams.LastName, result.Employee.LastName)
	require.Equal(t, arg.CreateEmployeeParams.Position, result.Employee.Position)
	require.Equal(t, arg.CreateEmployeeParams.Department, result.Employee.Department)
	require.Equal(t, arg.CreateEmployeeParams.LocationID, result.Employee.LocationID)

	// Verify IDs were generated
	require.NotZero(t, result.User.ID)
	require.NotZero(t, result.Employee.ID)
}

// func TestCreateClientwithAccountTx(t *testing.T) {
// 	store := NewStore(testDB)
// 	location := CreateRandomLocation(t)
// 	clientArg := CreateClientDetailsParams{
// 		FirstName:             util.RandomString(6),
// 		LastName:              util.RandomString(6),
// 		Email:                 util.RandomEmail(),
// 		DateOfBirth:           pgtype.Date{Time: time.Now(), Valid: true},
// 		Identity:              false,
// 		Status:                util.StringPtr("active"),
// 		Bsn:                   util.StringPtr(util.RandomString(9)),
// 		Source:                util.StringPtr("web"),
// 		Birthplace:            util.StringPtr(util.RandomString(6)),
// 		PhoneNumber:           util.StringPtr(util.RandomString(10)),
// 		Organisation:          util.StringPtr(util.RandomString(6)),
// 		Departement:           util.StringPtr(util.RandomString(6)),
// 		Gender:                "male",
// 		Filenumber:            1235,
// 		ProfilePicture:        util.StringPtr(util.GetRandomImageURL()),
// 		Infix:                 util.StringPtr(util.RandomString(6)),
// 		SenderID:              util.IntPtr(1),
// 		LocationID:            util.IntPtr(location.ID),
// 		IdentityAttachmentIds: []byte(util.RandomString(5)),
// 		DepartureReason:       util.StringPtr(util.RandomString(6)),
// 		DepartureReport:       util.StringPtr(util.RandomString(6)),
// 		GpsPosition:           []byte(util.RandomString(5)),
// 		MaturityDomains:       []byte(util.RandomString(5)),
// 		Addresses:             []byte(util.RandomString(5)),
// 		LegalMeasure:          util.StringPtr(util.RandomString(6)),
// 		HasUntakenMedications: false,
// 	}

// 	userArg := CreateUserParams{
// 		Username:    util.StringPtr(util.RandomString(5)),
// 		FirstName:   clientArg.FirstName,
// 		LastName:    clientArg.LastName,
// 		Email:       clientArg.Email,
// 		Password:    "password123",
// 		IsActive:    true,
// 		IsStaff:     false,
// 		IsSuperuser: false,
// 	}

// 	arg := CreateClientWithAccountTxParams{
// 		CreateClientParams: clientArg,
// 		CreateUserParams:   userArg,
// 	}
// 	client, err := store.CreateClientWithAccountTx(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, client)
// 	require.Equal(t, clientArg.FirstName, client.Client.FirstName)
// 	require.Equal(t, clientArg.LastName, client.Client.LastName)
// 	require.Equal(t, clientArg.Email, client.Client.Email)
// 	require.Equal(t, clientArg.PhoneNumber, client.Client.PhoneNumber)
// 	require.Equal(t, clientArg.DateOfBirth, client.Client.DateOfBirth)
// 	require.Equal(t, clientArg.Gender, client.Client.Gender)
// 	require.Equal(t, clientArg.IdentityAttachmentIds, client.Client.IdentityAttachmentIds)
// 	require.Equal(t, clientArg.DepartureReason, client.Client.DepartureReason)
// 	require.Equal(t, clientArg.DepartureReport, client.Client.DepartureReport)
// 	require.Equal(t, clientArg.GpsPosition, client.Client.GpsPosition)
// 	require.Equal(t, clientArg.MaturityDomains, client.Client.MaturityDomains)
// }
