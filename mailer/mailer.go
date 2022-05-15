package mailer

import (
	"errors"
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
		println(err.Error())
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

// send password recovery mail
func (m *Mailer) SendPasswordRecoveryMail(email string, name string, code string) error {
	if email == "" {
		return errors.New("email required")
	}

	body := `
	<h1>Hi %s,</h1>
	<br>
	<p>You have requested to reset your password. Here is the code below to reset your password.</p>
	<p>Code: %s</p>
	<br>
	<div>Thanks! – The Go Auth team</div>
	`

	return m.Send(MailParams{
		To:      email,
		Subject: "Password recovery",
		Body:    fmt.Sprintf(body, name, code),
	})
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
	return m.Send(MailParams{
		To:      email,
		Subject: "Confirm your account",
		Body:    fmt.Sprintf(body, name, link),
	})
}

func (m *Mailer) SendMagicLink(email string, name string, link string) error {
	if email == "" {
		return fmt.Errorf("email is empty")
	}
	body := `
	<div>Hi %s,</div>
	<br>
	<div>You have requested a magic link to access the site.</div>
	<br>
	<div>Click <a href="%s">here</a> to access the site.</div>
	<br>
	<div>Thanks! – The Go Auth team</div>
	`

	return m.Send(MailParams{
		To:      email,
		Subject: "Magic link",
		Body:    fmt.Sprintf(body, name, link),
	})
}
