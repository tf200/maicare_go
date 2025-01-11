package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetPermissionsByRoleID(t *testing.T) {
	permissions, err := testQueries.GetPermissionsByRoleID(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, permissions)
	for _, permission := range permissions {
		require.NotEmpty(t, permission)
	}

}

func TestGetRoleByID(t *testing.T) {
	role, err := testQueries.GetRoleByID(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, role)
	require.Equal(t, role.ID, int32(1))
}

func TestAssignRoleToUser(t *testing.T) {
	user := CreateRandomUser(t)
	updatedUser, err := testQueries.AssignRoleToUser(context.Background(),
		AssignRoleToUserParams{
			RoleID: 2,
			ID:     user.ID,
		})
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.Equal(t, updatedUser.ID, user.ID)
	require.Equal(t, updatedUser.RoleID, int32(2))

}
