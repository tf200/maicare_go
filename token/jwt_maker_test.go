package token

import (
	"maicare_go/util"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	accessKey := util.RandomString(32)
	refreshKey := util.RandomString(32)
	maker, err := NewJWTMaker(accessKey, refreshKey)
	require.NoError(t, err)

	userID := util.RandomInt(9999, 5555)
	RoleID := util.RandomInt(9999, 5555)
	duration := time.Minute

	testCases := []struct {
		name      string
		tokenType TokenType
	}{
		{
			name:      "AccessToken",
			tokenType: AccessToken,
		},
		{
			name:      "RefreshToken",
			tokenType: RefreshToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, payload, err := maker.CreateToken(userID, RoleID, duration, tc.tokenType)
			require.NoError(t, err)
			require.NotEmpty(t, token)
			require.NotNil(t, payload)

			payload, err = maker.VerifyToken(token)
			require.NoError(t, err)
			require.NotEmpty(t, payload)

			require.Equal(t, userID, payload.UserId)
			require.Equal(t, RoleID, payload.RoleID)
			require.Equal(t, tc.tokenType, payload.TokenType)
			require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
			require.WithinDuration(t, time.Now().Add(duration), payload.ExpiresAt, time.Second)
		})
	}
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32), util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(util.RandomInt(9999, 5555), util.RandomInt(9999, 5555), -time.Minute, AccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(util.RandomInt(9999, 5555), util.RandomInt(9999, 5555), time.Minute, AccessToken)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32), util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
