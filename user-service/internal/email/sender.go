package email

import (
	"fmt"
	"net/smtp"
	"os"
)

type EmailSender struct {
	host     string
	port     string
	username string
	password string
}

func NewEmailSender() *EmailSender {
	return &EmailSender{
		host:     os.Getenv("SMTP_HOST"),
		port:     os.Getenv("SMTP_PORT"),
		username: os.Getenv("SMTP_USER"),
		password: os.Getenv("SMTP_PASS"),
	}
}

func (s *EmailSender) Send(to, subject, body string) error {
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" + body)

	return smtp.SendMail(addr, auth, s.username, []string{to}, msg)
}
