package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

// Embed the templates directory.
//
//go:embed templates/*
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

// Initialize a new mail.Dialer with the given SMTP server settings and configure it to use a 5-second timeout for sending emails.
func New(host string, port int, username, password, sender string) Mailer {
	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

func (m Mailer) Send(recipient, templateFile string, data interface{}) error {
	// Ensure the templateFile path matches the structure in your embedded filesystem.
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	// Execute the "subject" template and store the result in a bytes.Buffer variable.
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	// Execute the "plainBody" template and store the result in plainBody.
	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	// Execute the "htmlBody" template and store the result in htmlBody.
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	// Initialize a new mail.Message instance.
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String()) // AddAlternative() after SetBody().

	// Send the message via dialer.DialAndSend().
	// err = m.dialer.DialAndSend(msg)
	// if err != nil {
	// 	return err
	// }

	// return nil

	// Retry logic: attempt to send the email up to 3 times
	for i := 1; i <= 3; i++ {
		err = m.dialer.DialAndSend(msg)
		// If everything worked, return nil.
		if err == nil {
			return nil
		}
		// If it didn't work, sleep for a short time and retry.
		time.Sleep(500 * time.Millisecond)
	}
	return err
}
