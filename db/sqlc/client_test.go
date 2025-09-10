package db

import (
	"context"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"

	"maicare_go/util"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomClientDetails(t *testing.T) ClientDetail {
	location := CreateRandomLocation(t)
	employee, _ := createRandomEmployee(t) // Uncomment if you want to use employee for BsnVerifiedBy
	sender := createRandomSenders(t)

	arg := CreateClientDetailsParams{
		FirstName:       faker.FirstName(),
		LastName:        faker.LastName(),
		Email:           faker.Email(),
		PhoneNumber:     util.StringPtr(faker.Phonenumber()),
		DateOfBirth:     pgtype.Date{Time: time.Now().AddDate(-20, 0, 0), Valid: true},
		Identity:        false,
		Bsn:             util.StringPtr(util.RandomString(9)),
		BsnVerifiedBy:   &employee.ID, // Assuming employee is created and has an ID
		Source:          util.StringPtr("Test Source"),
		Birthplace:      util.StringPtr("test city"),
		Organisation:    util.StringPtr("test org"),
		Departement:     util.StringPtr("test dep"),
		Gender:          "male", // or "Female" or other values as per your requirements
		Filenumber:      "testfile",
		ProfilePicture:  util.StringPtr(util.GetRandomImageURL()),
		Infix:           util.StringPtr("van"),
		SenderID:        &sender.ID,
		LocationID:      util.IntPtr(location.ID),
		DepartureReason: util.StringPtr("test Reason"),
		DepartureReport: util.StringPtr("test report"),
		Addresses:       []byte("[]"),
		LegalMeasure:    util.StringPtr("test measure"),
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
	require.Equal(t, arg.DepartureReason, client.DepartureReason)
	require.Equal(t, arg.DepartureReport, client.DepartureReport)
	require.Equal(t, arg.Addresses, client.Addresses)
	require.Equal(t, arg.LegalMeasure, client.LegalMeasure)
	return client
}

func TestCreateClientDetails(t *testing.T) {
	createRandomClientDetails(t)
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

func TestUpdateClientDetails(t *testing.T) {
	client := createRandomClientDetails(t)
	arg := UpdateClientDetailsParams{
		ID:        client.ID,
		FirstName: util.StringPtr(faker.FirstName()),
	}

	clientUp, err := testQueries.UpdateClientDetails(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, clientUp)
	require.Equal(t, arg.ID, clientUp.ID)
	require.Equal(t, arg.FirstName, &clientUp.FirstName)

}
func TestUpdateClientStatus(t *testing.T) {
	client := createRandomClientDetails(t)
	arg := UpdateClientStatusParams{
		ID:     client.ID,
		Status: util.StringPtr("In Care"),
	}

	clientUp, err := testQueries.UpdateClientStatus(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, clientUp)
	require.Equal(t, arg.ID, clientUp.ID)
	require.Equal(t, arg.Status, clientUp.Status)
}

func TestCreateClientStatusHistory(t *testing.T) {
	client := createRandomClientDetails(t)
	arg := CreateClientStatusHistoryParams{
		ClientID:  client.ID,
		OldStatus: util.StringPtr("On Waiting List"),
		NewStatus: "In Care",
		Reason:    util.StringPtr("Test Reason"),
	}

	clientStatus, err := testQueries.CreateClientStatusHistory(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, clientStatus)
	require.Equal(t, arg.ClientID, clientStatus.ClientID)
	require.Equal(t, arg.OldStatus, clientStatus.OldStatus)
}

func TestListClientStatusHistory(t *testing.T) {
	client := createRandomClientDetails(t)
	arg := CreateClientStatusHistoryParams{
		ClientID:  client.ID,
		OldStatus: util.StringPtr("On Waiting List"),
		NewStatus: "In Care",
		Reason:    util.StringPtr("Test Reason"),
	}

	_, err := testQueries.CreateClientStatusHistory(context.Background(), arg)
	require.NoError(t, err)

	arg2 := CreateClientStatusHistoryParams{
		ClientID:  client.ID,
		OldStatus: util.StringPtr("In Care"),
		NewStatus: "Out Of Care",
		Reason:    util.StringPtr("Test Reason"),
	}

	_, err = testQueries.CreateClientStatusHistory(context.Background(), arg2)
	require.NoError(t, err)

	historyList, err := testQueries.ListClientStatusHistory(context.Background(), ListClientStatusHistoryParams{
		ClientID: client.ID,
		Limit:    5,
		Offset:   0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, historyList)
	require.Len(t, historyList, 2)
}

func TestCreateSchedueledClientStatusChange(t *testing.T) {
	client := createRandomClientDetails(t)
	arg := CreateSchedueledClientStatusChangeParams{
		ClientID:      client.ID,
		NewStatus:     util.StringPtr("In Care"),
		Reason:        util.StringPtr("Test Reason"),
		ScheduledDate: pgtype.Date{Time: time.Now().AddDate(0, 0, 7), Valid: true},
	}

	clientStatus, err := testQueries.CreateSchedueledClientStatusChange(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, clientStatus)
	require.Equal(t, arg.ClientID, clientStatus.ClientID)
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

func TestCreateClientDocument(t *testing.T) {
	client := createRandomClientDetails(t)
	attachment := createRandomAttachmentFile(t)

	arg := CreateClientDocumentParams{
		ClientID:       client.ID,
		AttachmentUuid: &attachment.Uuid,
		Label:          "registration_form",
	}

	clientDoc, err := testQueries.CreateClientDocument(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, clientDoc)
	require.Equal(t, arg.ClientID, clientDoc.ClientID)
	require.Equal(t, arg.AttachmentUuid, clientDoc.AttachmentUuid)
}

func addRandomClientDocument(t *testing.T, ClientID int64) ClientDocument {
	store := NewStore(testDB)
	attachment := createRandomAttachmentFile(t)

	arg := AddClientDocumentTxParams{
		ClientID:     ClientID,
		AttachmentID: attachment.Uuid,
		Label:        "registration_form",
	}

	clientDoc, err := store.AddClientDocumentTx(context.Background(), arg)
	require.NoError(t, err)

	require.NotEmpty(t, clientDoc)
	require.Equal(t, arg.ClientID, clientDoc.ClientDocument.ClientID)
	return clientDoc.ClientDocument
}

func TestAddClientDocumentTx(t *testing.T) {
	client := createRandomClientDetails(t)
	addRandomClientDocument(t, client.ID)
}

func TestListClientDocuments(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 10; i++ {
		addRandomClientDocument(t, client.ID)
	}

	testCases := []struct {
		name  string
		arg   ListClientDocumentsParams
		check func(t *testing.T, clientDocs []ListClientDocumentsRow)
	}{
		{
			name: "base case",
			arg: ListClientDocumentsParams{
				ClientID: client.ID,
				Limit:    5,
				Offset:   0,
			},
			check: func(t *testing.T, clientDocs []ListClientDocumentsRow) {
				require.NotEmpty(t, clientDocs)
				require.Len(t, clientDocs, 5)
			},
		},
		{
			name: "with offset",
			arg: ListClientDocumentsParams{
				ClientID: client.ID,
				Limit:    5,
				Offset:   5,
			},
			check: func(t *testing.T, clientDocs []ListClientDocumentsRow) {
				require.NotEmpty(t, clientDocs)
				require.Len(t, clientDocs, 5)
			},
		},
	}
	for i := range testCases {
		t.Run(testCases[i].name, func(t *testing.T) {
			clientDocs, err := testQueries.ListClientDocuments(context.Background(), testCases[i].arg)
			require.NoError(t, err)
			testCases[i].check(t, clientDocs)
		})
	}
}

func TestDeleteClientDocument(t *testing.T) {
	client := createRandomClientDetails(t)
	clientDoc := addRandomClientDocument(t, client.ID)

	_, err := testQueries.DeleteClientDocument(context.Background(), clientDoc.AttachmentUuid)
	require.NoError(t, err)

	clientDocs, err := testQueries.ListClientDocuments(context.Background(), ListClientDocumentsParams{
		ClientID: client.ID,
	})
	require.NoError(t, err)
	require.Empty(t, clientDocs)
}

func TestDeleteClientDocumentTx(t *testing.T) {
	client := createRandomClientDetails(t)
	clientDoc := addRandomClientDocument(t, client.ID)

	store := NewStore(testDB)
	_, err := store.DeleteClientDocumentTx(context.Background(), DeleteClientDocumentParams{
		AttachmentID: *clientDoc.AttachmentUuid,
	})
	require.NoError(t, err)

	clientDocs, err := testQueries.ListClientDocuments(context.Background(), ListClientDocumentsParams{
		ClientID: client.ID,
	})
	require.NoError(t, err)
	require.Empty(t, clientDocs)
}

func TestGetMissingClientDocuments(t *testing.T) {
	client := createRandomClientDetails(t)
	clientDoc := addRandomClientDocument(t, client.ID)

	missingDocs, err := testQueries.GetMissingClientDocuments(context.Background(), client.ID)
	require.NoError(t, err)
	require.NotEmpty(t, missingDocs)
	require.NotContains(t, missingDocs, clientDoc.Label)
}
