package db

import (
	"context"
	"maicare_go/util"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func TestAssignSender(t *testing.T) {
	client := createRandomClientDetails(t)
	sender := createRandomSenders(t)

	newClient, err := testQueries.AssignSender(context.Background(), AssignSenderParams{
		ID:       client.ID,
		SenderID: &sender.ID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, newClient)
	require.Equal(t, client.ID, newClient.ID)
	require.Equal(t, &sender.ID, newClient.SenderID)
}
func TestGetClientSender(t *testing.T) {
	client := createRandomClientDetails(t)

	sender, err := testQueries.GetClientSender(context.Background(), client.ID)
	require.NoError(t, err)
	require.NotEmpty(t, sender)
	require.Equal(t, client.SenderID, &sender.ID)
}

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

func assignRandomEmployee(t *testing.T, clientID int64, employeeID int64) AssignedEmployee {
	arg := AssignEmployeeParams{
		ClientID:   clientID,
		EmployeeID: employeeID,
		StartDate:  pgtype.Date{Time: time.Now(), Valid: true},
		Role:       "Primary Caregiver",
	}
	assign, err := testQueries.AssignEmployee(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, assign)
	require.Equal(t, arg.ClientID, assign.ClientID)
	require.Equal(t, arg.EmployeeID, assign.EmployeeID)

	return assign
}

func TestAssignEmployee(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	assignRandomEmployee(t, client.ID, employee.ID)

}

func TestListAssignedEmployees(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	for i := 0; i < 10; i++ {
		assignRandomEmployee(t, client.ID, employee.ID)
	}
	arg := ListAssignedEmployeesParams{
		ClientID: client.ID,
		Limit:    5,
		Offset:   5,
	}
	assigns, err := testQueries.ListAssignedEmployees(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, assigns, 5)
}

func TestGetAssignEmployee(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	assign1 := assignRandomEmployee(t, client.ID, employee.ID)
	assign2, err := testQueries.GetAssignedEmployee(context.Background(), assign1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, assign2)
	require.Equal(t, assign1.ID, assign2.ID)
	require.Equal(t, assign1.ClientID, assign2.ClientID)
}

func TestUpdateAssignEmployee(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	assign1 := assignRandomEmployee(t, client.ID, employee.ID)
	arg := UpdateAssignedEmployeeParams{
		ID:   assign1.ID,
		Role: util.StringPtr("Secondary Caregiver"),
	}
	assign2, err := testQueries.UpdateAssignedEmployee(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, assign2)
	require.Equal(t, assign1.ID, assign2.ID)
	require.NotEqual(t, assign1.Role, assign2.Role)
}

func TestGetClientRelatedEmails(t *testing.T) {
	client := createRandomClientDetails(t)
	employee, _ := createRandomEmployee(t)
	_ = assignRandomEmployee(t, client.ID, employee.ID)
	emergencyContact := createRandomEmergencyContact(t, client.ID)

	emails, err := testQueries.GetClientRelatedEmails(context.Background(), client.ID)
	require.NoError(t, err)
	require.NotEmpty(t, emails)
	require.Contains(t, emails, employee.Email)
	require.NotNil(t, emergencyContact.Email)
	require.Contains(t, emails, *emergencyContact.Email)
}
