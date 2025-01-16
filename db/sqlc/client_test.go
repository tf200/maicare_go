package db

import (
	"context"
	"testing"
	"time"

	"maicare_go/util"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestCreateClientDetails(t *testing.T) {
	user := CreateRandomUser(t)
	location := CreateRandomLocation(t)

	arg := CreateClientDetailsParams{
		UserID:                user.ID,
		FirstName:             util.RandomString(5),
		LastName:              util.RandomString(5),
		Email:                 util.RandomEmail(),
		PhoneNumber:           util.StringPtr("0653316517"),
		DateOfBirth:           pgtype.Date{Time: time.Now().AddDate(-20, 0, 0), Valid: true},
		Identity:              false,
		Status:                util.StringPtr("On Waiting List"),
		Bsn:                   util.StringPtr(util.RandomString(9)),
		Source:                util.StringPtr("Test Source"),
		Birthplace:            util.StringPtr("test city"),
		Organisation:          util.StringPtr("test org"),
		Departement:           util.StringPtr("test dep"),
		Gender:                "Male", // or "Female" or other values as per your requirements
		Filenumber:            int32(util.RandomInt(1, 1000)),
		ProfilePicture:        util.StringPtr("test-profile.jpg"),
		Infix:                 util.StringPtr("van"),
		SenderID:              util.IntPtr(5),
		LocationID:            util.IntPtr(location.ID),
		IdentityAttachmentIds: []byte("[]"),
		DepartureReason:       util.StringPtr("test Reason"),
		DepartureReport:       util.StringPtr("test report"),
		GpsPosition:           []byte("[]"),
		MaturityDomains:       []byte("[]"),
		Addresses:             []byte("[]"),
		LegalMeasure:          util.StringPtr("test measure"),
		HasUntakenMedications: false,
	}

	client, err := testQueries.CreateClientDetails(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, client)
	require.Equal(t, arg.UserID, client.UserID)
	require.Equal(t, arg.FirstName, client.FirstName)
	require.Equal(t, arg.LastName, client.LastName)
	require.Equal(t, arg.Email, client.Email)
	require.Equal(t, arg.PhoneNumber, client.PhoneNumber)
	require.Equal(t, arg.Filenumber, client.Filenumber)
	require.Equal(t, arg.ProfilePicture, client.ProfilePicture)
	require.Equal(t, arg.Infix, client.Infix)
	require.Equal(t, arg.SenderID, client.SenderID)
	require.Equal(t, arg.LocationID, client.LocationID)
	require.Equal(t, arg.IdentityAttachmentIds, client.IdentityAttachmentIds)
	require.Equal(t, arg.DepartureReason, client.DepartureReason)
	require.Equal(t, arg.DepartureReport, client.DepartureReport)
	require.Equal(t, arg.GpsPosition, client.GpsPosition)
	require.Equal(t, arg.MaturityDomains, client.MaturityDomains)
	require.Equal(t, arg.Addresses, client.Addresses)
	require.Equal(t, arg.LegalMeasure, client.LegalMeasure)
	require.Equal(t, arg.HasUntakenMedications, client.HasUntakenMedications)
}


