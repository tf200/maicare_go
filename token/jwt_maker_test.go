package token

import (
	"testing"
	"time"

	"maicare_go/util"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	accessKey := util.RandomString(32)
	refreshKey := util.RandomString(32)
	twoFATokenKey := util.RandomString(32)
	maker, err := NewJWTMaker(accessKey, refreshKey, twoFATokenKey)
	require.NoError(t, err)

	userID := uuid.New()
	employeeID := uuid.New()
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
			token, payload, err := maker.CreateToken(userID, employeeID, duration, tc.tokenType)
			require.NoError(t, err)
			require.NotEmpty(t, token)
			require.NotNil(t, payload)

			payload, err = maker.VerifyToken(token)
			require.NoError(t, err)
			require.NotEmpty(t, payload)

			require.Equal(t, userID, payload.UserId)
			require.Equal(t, tc.tokenType, payload.TokenType)
			require.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
			require.WithinDuration(t, time.Now().Add(duration), payload.ExpiresAt, time.Second)
		})
	}
}

func TestExpiredJWTToken(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32), util.RandomString(32), util.RandomString(32))
	require.NoError(t, err)

	token, payload, err := maker.CreateToken(uuid.New(), uuid.New(), -time.Minute, AccessToken)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, payload)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTTokenAlgNone(t *testing.T) {
	payload, err := NewPayload(uuid.New(), uuid.New(), time.Minute, AccessToken)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(util.RandomString(32), util.RandomString(32), util.RandomString(32))
	require.NoError(t, err)

	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}

func TestJWTMakerCreateTokenWithSessionID(t *testing.T) {
	maker, err := NewJWTMaker(util.RandomString(32), util.RandomString(32), util.RandomString(32))
	require.NoError(t, err)

	userID := uuid.New()
	employeeID := uuid.New()
	sessionID := uuid.New()

	token, payload, err := maker.CreateTokenWithSessionID(userID, employeeID, time.Minute, AccessToken, sessionID)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotNil(t, payload)
	require.Equal(t, sessionID, payload.SessionID)

	verifiedPayload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.Equal(t, sessionID, verifiedPayload.SessionID)
}
