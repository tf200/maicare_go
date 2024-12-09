package db

import (
	"context"
	"maicare_go/util"
	"testing"
	"time"

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
		Position:                  util.RandomPgText(),
		Department:                util.RandomPgText(),
		EmployeeNumber:            util.RandomPgText(),
		EmploymentNumber:          util.RandomPgText(),
		PrivateEmailAddress:       util.RandomPgText(),
		EmailAddress:              util.RandomPgText(),
		AuthenticationPhoneNumber: util.RandomPgText(),
		PrivatePhoneNumber:        util.RandomPgText(),
		WorkPhoneNumber:           util.RandomPgText(),
		DateOfBirth: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		HomeTelephoneNumber: util.RandomPgText(),
		IsSubcontractor:     util.RandomPgBool(),
		Gender:              util.RandomPgText(),
		LocationID: pgtype.Int8{
			Int64: location.ID,
			Valid: true,
		},
		HasBorrowed:  false,
		OutOfService: util.RandomPgBool(),
		IsArchived:   util.RandomPgBool(),
	}
	userArg := CreateUserParams{
		Username:    util.RandomString(5),
		FirstName:   employeeArg.FirstName,
		LastName:    employeeArg.LastName,
		Email:       employeeArg.EmailAddress.String,
		Password:    "password123",
		IsActive:    true,
		IsStaff:     false,
		IsSuperuser: false,
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
	require.Equal(t, arg.CreateUserParams.IsStaff, result.User.IsStaff)
	require.Equal(t, arg.CreateUserParams.IsSuperuser, result.User.IsSuperuser)

	// Check employee profile was created properly
	require.NotEmpty(t, result.Employee)
	require.Equal(t, result.User.ID, result.Employee.UserID)
	require.Equal(t, arg.CreateEmployeeParams.FirstName, result.Employee.FirstName)
	require.Equal(t, arg.CreateEmployeeParams.LastName, result.Employee.LastName)
	require.Equal(t, arg.CreateEmployeeParams.Position, result.Employee.Position)
	require.Equal(t, arg.CreateEmployeeParams.Department, result.Employee.Department)
	require.Equal(t, arg.CreateEmployeeParams.LocationID, result.Employee.LocationID)
	require.Equal(t, arg.CreateEmployeeParams.HasBorrowed, result.Employee.HasBorrowed)

	// Verify IDs were generated
	require.NotZero(t, result.User.ID)
	require.NotZero(t, result.Employee.ID)

}
