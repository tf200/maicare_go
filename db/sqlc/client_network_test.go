package db

import (
	"context"
	"maicare_go/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEmergencyContact(t *testing.T, clientID int64) ClientEmergencyContact {

	arg := CreateEmemrgencyContactParams{
		ClientID:         clientID,
		FirstName:        util.StringPtr(util.RandomString(5)),
		LastName:         util.StringPtr(util.RandomString(5)),
		Email:            util.StringPtr(util.RandomEmail()),
		PhoneNumber:      util.StringPtr(util.RandomString(4)),
		Address:          util.StringPtr(util.RandomString(5)),
		Relationship:     util.StringPtr(util.RandomString(5)),
		RelationStatus:   util.StringPtr("Primary Relationship"),
		MedicalReports:   util.RandomBool(),
		IncidentsReports: util.RandomBool(),
		GoalsReports:     util.RandomBool(),
	}

	contact, err := testQueries.CreateEmemrgencyContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact)
	return contact
}

func TestCreateEmemrgencyContact(t *testing.T) {
	client := createRandomClientDetails(t)
	createRandomEmergencyContact(t, client.ID)
}

func TestListEmergencyContacts(t *testing.T) {
	client := createRandomClientDetails(t)
	for i := 0; i < 10; i++ {
		createRandomEmergencyContact(t, client.ID)
	}
	arg := ListEmergencyContactsParams{
		ClientID: client.ID,
		Limit:    5,
		Offset:   5,
	}
	contacts, err := testQueries.ListEmergencyContacts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, contacts, 5)
}

func TestGetEmergencyContact(t *testing.T) {
	client := createRandomClientDetails(t)
	contact1 := createRandomEmergencyContact(t, client.ID)
	contact2, err := testQueries.GetEmergencyContact(context.Background(), contact1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, contact2)
	require.Equal(t, contact1.ID, contact2.ID)
	require.Equal(t, contact1.ClientID, contact2.ClientID)
	require.Equal(t, contact1.FirstName, contact2.FirstName)
	require.Equal(t, contact1.LastName, contact2.LastName)
	require.Equal(t, contact1.Email, contact2.Email)
	require.Equal(t, contact1.PhoneNumber, contact2.PhoneNumber)
	require.Equal(t, contact1.Address, contact2.Address)
	require.Equal(t, contact1.Relationship, contact2.Relationship)
	require.Equal(t, contact1.RelationStatus, contact2.RelationStatus)
	require.Equal(t, contact1.MedicalReports, contact2.MedicalReports)
	require.Equal(t, contact1.IncidentsReports, contact2.IncidentsReports)
	require.Equal(t, contact1.GoalsReports, contact2.GoalsReports)
}

func TestUpdateEmergencyContact(t *testing.T) {
	client := createRandomClientDetails(t)
	contact1 := createRandomEmergencyContact(t, client.ID)
	arg := UpdateEmergencyContactParams{
		ID:        contact1.ID,
		FirstName: util.StringPtr(util.RandomString(5)),
	}
	contact2, err := testQueries.UpdateEmergencyContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact2)
	require.Equal(t, contact1.ID, contact2.ID)
	require.NotEqual(t, contact1.FirstName, contact2.FirstName)
}

func TestDeleteEmergencyContact(t *testing.T) {
	client := createRandomClientDetails(t)
	contact1 := createRandomEmergencyContact(t, client.ID)
	_, err := testQueries.DeleteEmergencyContact(context.Background(), contact1.ID)
	require.NoError(t, err)
	contact2, err := testQueries.GetEmergencyContact(context.Background(), contact1.ID)
	require.Error(t, err)
	require.Empty(t, contact2)
}
