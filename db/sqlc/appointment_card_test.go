package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomAppointmentCard(t *testing.T, clientID int64) AppointmentCard {

	arg := CreateAppointmentCardParams{
		ClientID:               clientID,
		GeneralInformation:     []string{"Client is doing well", "No concerns raised"},
		ImportantContacts:      []string{"Mother - 555-123-4567", "Case Worker - 555-987-6543"},
		HouseholdInfo:          []string{"Lives with mother and younger sibling", "Stable home environment"},
		OrganizationAgreements: []string{"Signed agreement with YMCA", "Participating in job training program"},
		YouthOfficerAgreements: []string{"Compliant with curfew", "Attending mandatory check-ins"},
		TreatmentAgreements:    []string{"Attending therapy sessions", "Taking prescribed medication"},
		SmokingRules:           []string{"Agreed to smoke outside only", "Limited to 5 cigarettes per day"},
		Work:                   []string{"Employed at local grocery store", "Works 20 hours per week"},
		SchoolInternship:       []string{"Interning at a law firm", "Gaining valuable experience"},
		Travel:                 []string{"Planning a trip to visit family", "Needs assistance with travel arrangements"},
		Leave:                  []string{"Took a week off for vacation", "Returned to work on August 1st"},
	}
	appointmentCard, err := testQueries.CreateAppointmentCard(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, appointmentCard)
	require.Equal(t, arg.ClientID, appointmentCard.ClientID)
	require.Equal(t, arg.GeneralInformation, appointmentCard.GeneralInformation)
	require.Equal(t, arg.ImportantContacts, appointmentCard.ImportantContacts)
	require.Equal(t, arg.HouseholdInfo, appointmentCard.HouseholdInfo)
	return appointmentCard
}

func TestCreateAppointmentCard(t *testing.T) {
	client := createRandomClientDetails(t)
	createRandomAppointmentCard(t, client.ID)
}

func TestGetAppointmentCard(t *testing.T) {
	client := createRandomClientDetails(t)
	appointmentCard1 := createRandomAppointmentCard(t, client.ID)
	appointmentCard2, err := testQueries.GetAppointmentCard(context.Background(), client.ID)
	require.NoError(t, err)
	require.NotEmpty(t, appointmentCard2)
	require.Equal(t, appointmentCard1.ID, appointmentCard2.ID)
}

func TestUpdateAppointmentCard(t *testing.T) {
	client := createRandomClientDetails(t)
	appointmentCard1 := createRandomAppointmentCard(t, client.ID)
	arg := UpdateAppointmentCardParams{
		ClientID:           client.ID,
		GeneralInformation: []string{"Updated: Client is doing well", "Updated: No concerns raised"},
		ImportantContacts:  []string{"Updated: Mother - 555-123-4567", "Updated: Case Worker - 555-987-6543"},
		HouseholdInfo:      []string{"Updated: Lives with mother and younger sibling", "Updated: Stable home environment"},
	}
	appointmentCard2, err := testQueries.UpdateAppointmentCard(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, appointmentCard2)
	require.Equal(t, appointmentCard1.ID, appointmentCard2.ID)
	require.NotEqual(t, appointmentCard1.GeneralInformation, appointmentCard2.GeneralInformation)
	require.NotEqual(t, appointmentCard1.ImportantContacts, appointmentCard2.ImportantContacts)
	require.NotEqual(t, appointmentCard1.HouseholdInfo, appointmentCard2.HouseholdInfo)
	require.Equal(t, appointmentCard1.OrganizationAgreements, appointmentCard2.OrganizationAgreements)
}
