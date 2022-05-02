package mailer

import (
	"fmt"

	"github.com/overlorddamygod/go-auth/configs"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mailer struct {
	client *mail.SMTPClient
}

// new mailer
func NewMailer() *Mailer {
	// intialized mail
	mailServer := mail.NewSMTPClient()

	mailConfig := configs.GetConfig().Mail

	mailServer.Host = mailConfig.Host
	mailServer.Port = mailConfig.Port
	mailServer.Username = mailConfig.Username
	mailServer.Password = mailConfig.Password
	mailServer.Encryption = mail.EncryptionTLS

	smtpClient, err := mailServer.Connect()
	if err != nil {
		// log.Fatal(err)
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
