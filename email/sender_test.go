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
	email := NewSmtpConf("config", "dev@maicare.online", "u,Q4(;9^$tzWZjm", "mail.privateemail.com", 465)
	err := email.SendCredentials(context.Background(), []string{"farjiataha@gmail.com"}, Credentials{
		Name:     "John Doe",
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
		DocumentLink: "https://www.w3.org/WAI/ER/tests/xhtml/testfiles/resources/pdf/dummy.pdf",
	}
	err := email.SendIncident(context.Background(), []string{"farjiataha@gmail.com"}, arg)
	require.NoError(t, err)
}
