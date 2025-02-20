package db

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomIntakeFormToken(t *testing.T) IntakeFormToken {
	arg := CreateIntakeFormTokenParams{
		Token: uuid.New().String(),
		ExpiresAt: pgtype.Timestamp{
			Time:  time.Now().Add(time.Hour),
			Valid: true,
		},
	}

	token, err := testQueries.CreateIntakeFormToken(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	return token
}

func TestCreateIntakeFormToken(t *testing.T) {
	createRandomIntakeFormToken(t)
}

func TestGetIntakeFormToken(t *testing.T) {
	token1 := createRandomIntakeFormToken(t)

	token2, err := testQueries.GetIntakeFormToken(context.Background(), token1.Token)
	require.NoError(t, err)
	require.NotEmpty(t, token2)

	require.Equal(t, token1.Token, token2.Token)
	require.Equal(t, token1.ExpiresAt, token2.ExpiresAt)
	require.Equal(t, token1.IsRevoked, token2.IsRevoked)
	require.WithinDuration(t, token1.CreatedAt.Time, token2.CreatedAt.Time, time.Second)
}

func TestRevokedIntakeFormToken(t *testing.T) {
	token1 := createRandomIntakeFormToken(t)

	_, err := testQueries.RevokedIntakeFormToken(context.Background(), token1.Token)
	require.NoError(t, err)

	token2, err := testQueries.GetIntakeFormToken(context.Background(), token1.Token)
	require.NoError(t, err)
	require.NotEmpty(t, token2)

	require.Equal(t, token1.Token, token2.Token)
	require.Equal(t, token1.ExpiresAt, token2.ExpiresAt)
	require.True(t, token2.IsRevoked)
	require.WithinDuration(t, token1.CreatedAt.Time, token2.CreatedAt.Time, time.Second)
}

func TestCreateIntakeForm(t *testing.T) {
	token := createRandomIntakeFormToken(t)

	arg := CreateIntakeFormParams{
		IntakeFormToken: token.Token,
		FirstName:       faker.FirstName(),
		LastName:        faker.LastName(),
		DateOfBirth: pgtype.Date{
			Time:  time.Now(),
			Valid: true,
		},
		Gender:                     faker.Gender(),
		PlaceOfBirth:               faker.GetRealAddress().Address,
		RepresentativeFirstName:    faker.FirstName(),
		RepresentativeLastName:     faker.LastName(),
		RepresentativeEmail:        faker.Email(),
		RepresentativePhoneNumber:  faker.Phonenumber(),
		RepresentativeRelationship: faker.WORD,
		RepresentativeAddress:      faker.GetRealAddress().Address,
		AttachementIds:             []uuid.UUID{},
	}

	form, err := testQueries.CreateIntakeForm(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, form)
	require.NotEmpty(t, form.ID)
	require.Equal(t, arg.IntakeFormToken, form.IntakeFormToken)
	require.Equal(t, arg.FirstName, form.FirstName)
	require.Equal(t, arg.LastName, form.LastName)

}
