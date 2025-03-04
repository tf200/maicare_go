package email

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/wneessen/go-mail"
)

type SmtpConf struct {
	Name          string
	Address       string
	Athentication string
	SmtpHost      string
	SmtpPort      int
}

type Credentials struct {
	Email    string
	Password string
}

type Incident struct {
	IncidentID   int64
	ReportedBy   string
	ClientName   string
	IncidentType string
	Severity     string
	Location     string
	DocumentLink string
}

func NewSmtpConf(name, address, authentication, smtpHost string, smtpPort int) *SmtpConf {
	return &SmtpConf{
		Name:          name,
		Address:       address,
		Athentication: authentication,
		SmtpHost:      smtpHost,
		SmtpPort:      smtpPort,
	}
}

func (e *SmtpConf) Send(subject, body string, to []string) error {
	message := mail.NewMsg()
	if err := message.From("dev@maicare.online"); err != nil {
		log.Fatalf("failed to set From address: %s", err)
	}
	if err := message.To(to[0]); err != nil {
		log.Fatalf("failed to set To address: %s", err)
	}
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextPlain, body)
	client, err := mail.NewClient(e.SmtpHost, mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(e.Address), mail.WithPassword(e.Athentication), mail.WithPort(e.SmtpPort))
	if err != nil {
		log.Fatalf("failed to create mail client: %s", err)
	}
	if err = client.DialAndSend(message); err != nil {
		fmt.Printf("failed to send mail: %s\n", err)
		os.Exit(1)
	}
	return nil
}

//go:embed templates/credentials.html
var credentialsTemplateFS embed.FS

func (s *SmtpConf) SendCredentials(ctx context.Context, to []string, data Credentials) error {

	if len(to) == 0 {
		return errors.New("no recipient addresses provided")
	}
	if s.SmtpHost == "" || s.SmtpPort == 0 {
		return errors.New("invalid SMTP configuration")
	}

	tmpl, err := template.ParseFS(credentialsTemplateFS, "templates/credentials.html")
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	message := mail.NewMsg()
	if err := message.From(fmt.Sprintf("%s <%s>", s.Name, s.Address)); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}
	if err := message.To(to...); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}
	message.Subject("Welcome to Maicare!")
	message.SetBodyString(mail.TypeTextHTML, body.String())

	client, err := mail.NewClient(
		s.SmtpHost,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(s.Address),
		mail.WithPassword(s.Athentication),
		mail.WithPort(s.SmtpPort),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	if err := client.DialAndSendWithContext(ctx, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Email sent to %s", to)

	return nil
}

func (s *SmtpConf) SendIncident(ctx context.Context, to []string, data Incident) error {

	if len(to) == 0 {
		return errors.New("no recipient addresses provided")
	}
	if s.SmtpHost == "" || s.SmtpPort == 0 {
		return errors.New("invalid SMTP configuration")
	}

	tmpl, err := template.ParseFiles("templates/incident.html")
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	message := mail.NewMsg()
	if err := message.From(fmt.Sprintf("%s <%s>", s.Name, s.Address)); err != nil {
		return fmt.Errorf("failed to set From address: %w", err)
	}
	if err := message.To(to...); err != nil {
		return fmt.Errorf("failed to set To address: %w", err)
	}
	message.Subject("Incident!")
	message.SetBodyString(mail.TypeTextHTML, body.String())

	client, err := mail.NewClient(
		s.SmtpHost,
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(s.Address),
		mail.WithPassword(s.Athentication),
		mail.WithPort(s.SmtpPort),
	)
	if err != nil {
		return fmt.Errorf("failed to create mail client: %w", err)
	}

	if err := client.DialAndSendWithContext(ctx, message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
