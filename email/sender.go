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

	brevo "github.com/getbrevo/brevo-go/lib"

	"github.com/wneessen/go-mail"
)

type SmtpConf struct {
	Name          string
	Address       string
	Athentication string
	SmtpHost      string
	SmtpPort      int
}

type BrevoConf struct {
	SenderName  string
	Senderemail string
	ApiKey      string
	client      *brevo.APIClient
}

func NewBrevoConf(senderName, senderEmail, apiKey string) *BrevoConf {
	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)
	return &BrevoConf{
		SenderName:  senderName,
		Senderemail: senderEmail,
		ApiKey:      apiKey,
		client:      brevo.NewAPIClient(cfg),
	}
}

type Credentials struct {
	Name     string
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

func (b *BrevoConf) SendCredentials(ctx context.Context, to []string, data Credentials) error {

	if len(to) == 0 {
		return errors.New("no recipient addresses provided")
	}

	if b.SenderName == "" || b.Senderemail == "" {
		return errors.New("invalid sender configuration")
	}
	if b.ApiKey == "" {
		return errors.New("invalid API key")
	}

	tmpl, err := template.ParseFS(credentialsTemplateFS, "templates/credentials.html")
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	htmlContent := body.String()

	sender := brevo.SendSmtpEmailSender{
		Name:  b.SenderName,
		Email: b.Senderemail,
	}

	recipients := make([]brevo.SendSmtpEmailTo, 0, len(to))
	for _, recipient := range to {
		recipients = append(recipients, brevo.SendSmtpEmailTo{
			Email: recipient,
			Name:  recipient,
		})
	}

	emailContent := brevo.SendSmtpEmail{
		Sender:      &sender,
		To:          recipients,
		Subject:     "Maicare Credentials",
		HtmlContent: htmlContent,
	}

	result, response, err := b.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailContent)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	if response.StatusCode != 201 {
		return fmt.Errorf("failed to send email, status code: %d", response.StatusCode)
	}
	log.Printf("Email sent to %s", to)
	log.Printf("Response: %s", result)
	log.Printf("Response Status Code: %d", response.StatusCode)
	log.Printf("Response Headers: %v", response.Header)
	log.Printf("Response Body: %s", response.Body)

	return nil
}

//go:embed templates/incident.html
var incidentTemplateFS embed.FS

func (b *BrevoConf) SendIncident(ctx context.Context, to []string, data Incident) error {

	if len(to) == 0 {
		return errors.New("no recipient addresses provided")
	}
	if b.SenderName == "" || b.Senderemail == "" {
		return errors.New("invalid sender configuration")
	}
	if b.ApiKey == "" {
		return errors.New("invalid API key")
	}

	tmpl, err := template.ParseFS(incidentTemplateFS, "templates/incident.html")
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	htmlContent := body.String()
	sender := brevo.SendSmtpEmailSender{
		Name:  b.SenderName,
		Email: b.Senderemail,
	}
	recipients := make([]brevo.SendSmtpEmailTo, 0, len(to))
	for _, recipient := range to {
		recipients = append(recipients, brevo.SendSmtpEmailTo{
			Email: recipient,
			Name:  recipient,
		})
	}
	emailContent := brevo.SendSmtpEmail{
		Sender:      &sender,
		To:          recipients,
		Subject:     "Incident Report",
		HtmlContent: htmlContent,
	}
	result, response, err := b.client.TransactionalEmailsApi.SendTransacEmail(ctx, emailContent)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	if response.StatusCode != 201 {
		return fmt.Errorf("failed to send email, status code: %d", response.StatusCode)
	}
	log.Printf("Email sent to %s", to)
	log.Printf("Response: %s", result)
	log.Printf("Response Status Code: %d", response.StatusCode)
	log.Printf("Response Headers: %v", response.Header)
	log.Printf("Response Body: %s", response.Body)

	return nil
}
