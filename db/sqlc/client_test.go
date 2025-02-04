package db

import (
	"context"
	"testing"
	"time"

	"github.com/goccy/go-json"

	"maicare_go/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomClientDetails(t *testing.T) ClientDetail {
	location := CreateRandomLocation(t)
	sender := createRandomSenders(t)

	arg := CreateClientDetailsParams{
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
		Filenumber:            "testfile",
		ProfilePicture:        util.StringPtr(util.GetRandomImageURL()),
		Infix:                 util.StringPtr("van"),
		SenderID:              sender.ID,
		LocationID:            util.IntPtr(location.ID),
		IdentityAttachmentIds: []byte("[]"),
		DepartureReason:       util.StringPtr("test Reason"),
		DepartureReport:       util.StringPtr("test report"),
		Addresses:             []byte("[]"),
		LegalMeasure:          util.StringPtr("test measure"),
	}

	client, err := testQueries.CreateClientDetails(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, client)
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
	require.Equal(t, arg.Addresses, client.Addresses)
	require.Equal(t, arg.LegalMeasure, client.LegalMeasure)
	return client
}

func TestCreateClientDetails(t *testing.T) {
	createRandomClientDetails(t)
}

func TestCreateClientDetailsTx(t *testing.T) {
	store := NewStore(testDB)
	attachement := createRandomAttachmentFile(t)

	attachementList := []uuid.UUID{attachement.Uuid}
	identityAttachmentIds, err := json.Marshal(attachementList)
	require.NoError(t, err)

	CreateClientParams := CreateClientDetailsParams{
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
		Filenumber:            "testfile",
		ProfilePicture:        util.StringPtr("test-profile.jpg"),
		Infix:                 util.StringPtr("van"),
		SenderID:              1,
		LocationID:            util.IntPtr(1),
		IdentityAttachmentIds: identityAttachmentIds,
		DepartureReason:       util.StringPtr("test Reason"),
		DepartureReport:       util.StringPtr("test report"),
		Addresses:             []byte("[]"),
		LegalMeasure:          util.StringPtr("test measure"),
	}

	arg := CreateClientDetailsTxParams{
		CreateClientParams:  CreateClientParams,
		IdentityAttachments: attachementList,
	}

	client, err := store.CreateClientDetailsTx(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, client)
	require.Equal(t, CreateClientParams.FirstName, client.Client.FirstName)
	require.Equal(t, CreateClientParams.LastName, client.Client.LastName)

}

func TestListClientDetails(t *testing.T) {
	var clients []ClientDetail
	for i := 0; i < 20; i++ {
		_ = append(clients, createRandomClientDetails(t))
	}
	testCases := []struct {
		name  string
		arg   ListClientDetailsParams
		check func(t *testing.T, clients []ListClientDetailsRow)
	}{
		{
			name: "base case",
			arg: ListClientDetailsParams{
				Limit:  5,
				Offset: 0,
			},
			check: func(t *testing.T, clients []ListClientDetailsRow) {
				require.NotEmpty(t, clients)
				require.Len(t, clients, 5)
			},
		},
		{
			name: "with offset",
			arg: ListClientDetailsParams{
				Limit:  5,
				Offset: 5,
			},
			check: func(t *testing.T, clients []ListClientDetailsRow) {
				require.NotEmpty(t, clients)
				require.Len(t, clients, 5)
			},
		},
		{
			name: "with search",
			arg: ListClientDetailsParams{
				Limit:  5,
				Offset: 0,
				Search: util.StringPtr("a"),
			},
			check: func(t *testing.T, clients []ListClientDetailsRow) {
				require.NotEmpty(t, clients)

			},
		},
		{
			name: "with status",
			arg: ListClientDetailsParams{
				Limit:  5,
				Offset: 0,
				Status: util.StringPtr("On Waiting List"),
			},
			check: func(t *testing.T, clients []ListClientDetailsRow) {
				require.NotEmpty(t, clients)
				require.Len(t, clients, 5)
				require.Equal(t, util.StringPtr("On Waiting List"), clients[0].Status)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			clients, err := testQueries.ListClientDetails(context.Background(), tc.arg)
			require.NoError(t, err)
			tc.check(t, clients)
		})
	}

}

func TestGetClientDetails(t *testing.T) {
	client := createRandomClientDetails(t)
	client1, err := testQueries.GetClientDetails(context.Background(), client.ID)
	require.NoError(t, err)
	require.NotEmpty(t, client1)
	require.Equal(t, client.ID, client1.ID)
	require.Equal(t, client.FirstName, client1.FirstName)
	require.Equal(t, client.LastName, client1.LastName)
	require.Equal(t, client.Email, client1.Email)
	require.Equal(t, client.PhoneNumber, client1.PhoneNumber)
	require.Equal(t, client.DateOfBirth, client1.DateOfBirth)
	require.Equal(t, client.Identity, client1.Identity)
	require.Equal(t, client.Status, client1.Status)
	require.Equal(t, client.Bsn, client1.Bsn)
	require.Equal(t, client.Source, client1.Source)
	require.Equal(t, client.Birthplace, client1.Birthplace)
	require.Equal(t, client.Organisation, client1.Organisation)
	require.Equal(t, client.Departement, client1.Departement)
}

func TestSetClientProfilePictureTx(t *testing.T) {
	store := NewStore(testDB)
	client := createRandomClientDetails(t)
	attachment := createRandomAttachmentFile(t)

	arg := SetClientProfilePictureTxParams{
		ClientID:     client.ID,
		AttachmentID: attachment.Uuid,
	}

	clientUp, err := store.SetClientProfilePictureTx(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, clientUp)
}

func createRandomClientAllergy(t *testing.T, clientID int64) ClientAllergy {

	arg := CreateClientAllergyParams{
		ClientID:      clientID,
		AllergyTypeID: 1,
		Severity:      "Mild",
		Reaction:      "test reaction",
		Notes:         util.StringPtr("test note"),
	}

	allergy, err := testQueries.CreateClientAllergy(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, allergy)
	require.Equal(t, arg.ClientID, allergy.ClientID)
	return allergy
}

func TestCreateClientAllergy(t *testing.T) {
	client := createRandomClientDetails(t)
	createRandomClientAllergy(t, client.ID)
}

func TestListClientAllergies(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 20; i++ {
		_ = createRandomClientAllergy(t, client.ID)
	}
	testCases := []struct {
		name  string
		arg   ListClientAllergiesParams
		check func(t *testing.T, allergies []ListClientAllergiesRow)
	}{
		{
			name: "base case",
			arg: ListClientAllergiesParams{
				ClientID: client.ID,
				Limit:    5,
				Offset:   0,
			},
			check: func(t *testing.T, allergies []ListClientAllergiesRow) {
				require.NotEmpty(t, allergies)
				require.Len(t, allergies, 5)
				require.Equal(t, int64(20), allergies[0].TotalAllergies)
			},
		},
		{
			name: "with offset",
			arg: ListClientAllergiesParams{
				ClientID: client.ID,
				Limit:    5,
				Offset:   5,
			},
			check: func(t *testing.T, allergies []ListClientAllergiesRow) {
				require.NotEmpty(t, allergies)
				require.Len(t, allergies, 5)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			allergies, err := testQueries.ListClientAllergies(context.Background(), tc.arg)
			require.NoError(t, err)
			tc.check(t, allergies)
		})
	}
}
