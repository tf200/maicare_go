package email

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	email := NewSmtpConf("config", "dev@maicare.online", "u,Q4(;9^$tzWZjm", "mail.privateemail.com", 587)
	err := email.Send("Test", "Test", []string{"farjiataha@gmail.com"})
	require.NoError(t, err)
}

func TestSendCredentials(t *testing.T) {
	err := testBrevo.SendCredentials(context.Background(), []string{"g.louz@diversitaszorg.nl"}, Credentials{
		Name:     "Guila Louz",
		Email:    "g.louz@diversitaszorg.nl",
		Password: "Secr3tPass!",
	})
	t.Log(err)
	require.NoError(t, err)

	// err = testBrevo.SendCredentials(context.Background(), []string{"h.belhadj@diversitaszorg.nl"}, Credentials{
	// 	Name:     "Hamza Belhadj",
	// 	Email:    "h.belhadj@diversitaszorg.nl",
	// 	Password: "Secr3tPass!",
	// })
	// t.Log(err)
	// require.NoError(t, err)

	err = testBrevo.SendCredentials(context.Background(), []string{"q.zeeman@diversitaszorg.nl"}, Credentials{
		Name:     "Quinten Zeeman",
		Email:    "q.zeeman@diversitaszorg.nl",


		
		Password: "Secr3tPass!",
	})
	t.Log(err)
	require.NoError(t, err)

	// err = testBrevo.SendCredentials(context.Background(), []string{"Q.zeeman@diversitaszorg.nl"}, Credentials{
	// 	Name:     "Quinten Zeeman",
	// 	Email:    "Q.zeeman@diversitaszorg.nl",
	// 	Password: "Secr3tPass!",
	// })
	// t.Log(err)
	// require.NoError(t, err)

}

func TestSendIncident(t *testing.T) {
	arg := Incident{
		IncidentID:   1,
		ReportedBy:   "John Doe",
		ClientName:   "Johny Doe",
		IncidentType: "workplace_accident",
		Severity:     "serious",
		Location:     "Main Building - Floor 3",
		DocumentLink: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	err := testBrevo.SendIncident(context.Background(), []string{"farjiataha@gmail.com"}, arg)
	require.NoError(t, err)
}

func TestSendAcceptedRegistrationForm(t *testing.T) {
	arg := AcceptedRegitrationForm{
		ReferrerName:        "Jane Smith",
		ChildName:           "Alice Doe",
		ChildBSN:            "123456789",
		AppointmentDate:     "2023-10-01 10:00:00Z",
		AppointmentLocation: "Main Office - Room 101",
	}
	err := testBrevo.SendAcceptedRegistrationForm(context.Background(), []string{"farjiataha@gmail.com"}, arg)
	require.NoError(t, err)
}

func TestSendClientContractReminder(t *testing.T) {
	arg := ClientContractReminder{
		ClientID:           1,
		ClientFirstName:    "Alice",
		ClientLastName:     "Doe",
		ContractID:         12345,
		CareType:           "ambulante",
		ContractStartDate:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
		ContractEndDate:    time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
		ContractStatus:     "active",
		ReminderType:       "initial",
		LastReminderSentAt: nil,
		CurrentDate:        time.Now(),
		CurrentYear:        time.Now().Year(),
	}
	err := testBrevo.SendClientContractReminder(context.Background(), []string{"farjiataha@gmail.com"}, arg)
	require.NoError(t, err)
}
