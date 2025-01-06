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
