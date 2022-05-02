package mailer

import (
	"fmt"
	"log"

	mail "github.com/xhit/go-simple-mail/v2"
)

type Mailer struct {
	client *mail.SMTPClient
}

// new mailer
func NewMailer() *Mailer {
	// intialized mail
	mailServer := mail.NewSMTPClient()

	mailServer.Host = "smtp.mailgun.org"
	mailServer.Port = 587
	mailServer.Username = "postmaster@sandboxd2841c95e0b341629e9b9f2cc63fd90e.mailgun.org"
	mailServer.Password = "2e2a9d65ef908ff91aca7bc9166a2b3e-1831c31e-db49bda2"
	mailServer.Encryption = mail.EncryptionTLS

	smtpClient, err := mailServer.Connect()
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to SMTP server")
	}

	return &Mailer{
		client: smtpClient,
	}
}

type MailParams struct {
	// From    string
	To      string
	Subject string
	Body    string
}

// send mail
func (m *Mailer) Send(params MailParams) error {
	email := mail.NewMSG()
	email.SetFrom("GoAuth <no-reply@goauth.com>")
	email.AddTo(params.To)
	email.SetSubject(params.Subject)

	email.SetBody(mail.TextHTML, params.Body)

	// Send email
	return email.Send(m.client)
}

func (m *Mailer) SendConfirmationMail(email string, name string, link string) error {
	if email == "" {
		return fmt.Errorf("email is empty")
	}
	body := `
	<div>Hi %s,</div>
	<br>
	<div>We just need to verify your email address before you can access the site.</div>
	<br>
	<div>Verify your email address <a href="%s">here</a></div>
	<br>
	<div>Thanks! – The Go Auth team</div>
	`
	err := m.Send(MailParams{
		To:      email,
		Subject: "Confirm your account",
		Body:    fmt.Sprintf(body, name, link),
	})
	return err
}
