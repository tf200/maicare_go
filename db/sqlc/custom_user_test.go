package db

import (
	"context"
	"testing"

	"maicare_go/util"

	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) *CustomUser {
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
	arg := CreateUserParams{
		Password:       hashedPassword,
		Username:       util.StringPtr(util.RandomString(5)),
		Email:          "testemail@gmail.com",
		FirstName:      util.RandomString(5),
		LastName:       util.RandomString(5),
		IsSuperuser:    true,
		IsStaff:        true,
		IsActive:       true,
		ProfilePicture: util.StringPtr(util.GetRandomImageURL()),
		PhoneNumber:    util.IntPtr(456),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.FirstName, user.FirstName)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.IsSuperuser, user.IsSuperuser)
	require.Equal(t, arg.IsStaff, user.IsStaff)
	require.Equal(t, arg.IsActive, user.IsActive)
	require.Equal(t, arg.ProfilePicture, user.ProfilePicture)
	require.Equal(t, arg.PhoneNumber, user.PhoneNumber)

	// Verify timestamps and ID are set
	require.NotZero(t, user.ID)

	return &user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUserByID(t *testing.T) {
	created_user := CreateRandomUser(t)
	id := created_user.ID

	retrieved_user, err := testQueries.GetUserByID(context.Background(), id)
	require.NoError(t, err)
	require.NotEmpty(t, retrieved_user)

	// Verify all fields match
	require.Equal(t, created_user.ID, retrieved_user.ID)
	require.Equal(t, created_user.Username, retrieved_user.Username)
	require.Equal(t, created_user.Password, retrieved_user.Password)
	require.Equal(t, created_user.FirstName, retrieved_user.FirstName)
	require.Equal(t, created_user.LastName, retrieved_user.LastName)
	require.Equal(t, created_user.Email, retrieved_user.Email)
	require.Equal(t, created_user.IsSuperuser, retrieved_user.IsSuperuser)
	require.Equal(t, created_user.IsStaff, retrieved_user.IsStaff)
	require.Equal(t, created_user.IsActive, retrieved_user.IsActive)
	require.Equal(t, created_user.ProfilePicture, retrieved_user.ProfilePicture)
	require.Equal(t, created_user.PhoneNumber, retrieved_user.PhoneNumber)
	require.Equal(t, created_user.DateJoined, retrieved_user.DateJoined)
}

func TestGetUserByUsername(t *testing.T) {
	created_user := CreateRandomUser(t)
	username := created_user.Username

	retrieved_user, err := testQueries.GetUserByUsername(context.Background(), username)
	require.NoError(t, err)
	require.NotEmpty(t, retrieved_user)

	// Verify all fields match
	require.Equal(t, created_user.ID, retrieved_user.ID)
	require.Equal(t, created_user.Username, retrieved_user.Username)
	require.Equal(t, created_user.Password, retrieved_user.Password)
	require.Equal(t, created_user.FirstName, retrieved_user.FirstName)
	require.Equal(t, created_user.LastName, retrieved_user.LastName)
	require.Equal(t, created_user.Email, retrieved_user.Email)
	require.Equal(t, created_user.IsSuperuser, retrieved_user.IsSuperuser)
	require.Equal(t, created_user.IsStaff, retrieved_user.IsStaff)
	require.Equal(t, created_user.IsActive, retrieved_user.IsActive)
	require.Equal(t, created_user.ProfilePicture, retrieved_user.ProfilePicture)
	require.Equal(t, created_user.PhoneNumber, retrieved_user.PhoneNumber)
	require.Equal(t, created_user.DateJoined, retrieved_user.DateJoined)
}
