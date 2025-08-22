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
		Email:          util.RandomEmail(),
		IsActive:       true,
		ProfilePicture: util.StringPtr(util.GetRandomImageURL()),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.IsActive, user.IsActive)
	require.Equal(t, arg.ProfilePicture, user.ProfilePicture)

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
	require.Equal(t, created_user.Password, retrieved_user.Password)
	require.Equal(t, created_user.Email, retrieved_user.Email)
	require.Equal(t, created_user.IsActive, retrieved_user.IsActive)
	require.Equal(t, created_user.ProfilePicture, retrieved_user.ProfilePicture)
	require.Equal(t, created_user.DateJoined, retrieved_user.DateJoined)
}
