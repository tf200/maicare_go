package email

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	email := NewSmtpConf("config", "dev@maicare.online", "u,Q4(;9^$tzWZjm", "mail.privateemail.com", 587)
	err := email.Send("Test", "Test", []string{"farjiataha@gmail.com"})
	require.NoError(t, err)
}

func TestSendCredentials(t *testing.T) {
	email := NewSmtpConf("config", "dev@maicare.online", "u,Q4(;9^$tzWZjm", "mail.privateemail.com", 587)
	err := email.SendCredentials(context.Background(), []string{"farjiataha@gmail.com"}, Credentials{
		Email:    "farjiataha@gmail.com",
		Password: "password",
	})
	require.NoError(t, err)
}

func TestSendIncident(t *testing.T) {
	email := NewSmtpConf("config", "dev@maicare.online", "u,Q4(;9^$tzWZjm", "mail.privateemail.com", 587)
	arg := Incident{
		IncidentID:   1,
		ReportedBy:   "John Doe",
		ClientName:   "Johny Doe",
		IncidentType: "workplace_accident",
		Severity:     "serious",
		Location:     "Main Building - Floor 3",
	}
	err := email.SendIncident(context.Background(), []string{"farjiataha@gmail.com"}, arg)
	require.NoError(t, err)
}
