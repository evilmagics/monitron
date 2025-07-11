package utils

import (
	"fmt"
	"net/smtp"

	"monitron-server/config"
)

// SendEmail sends an email using the configured SMTP server
func SendEmail(cfg *config.Config, to, subject, body string) error {
	from := cfg.Email.From
	password := cfg.Email.Password

	msg := []byte(fmt.Sprintf("To: %s\r\n" +
		"From: %s\r\n" +
		"Subject: %s\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n%s\r\n",
		to, from, subject, body))

	auth := smtp.PlainAuth("", from, password, cfg.Email.Host)

	err := smtp.SendMail(fmt.Sprintf("%s:%d", cfg.Email.Host, cfg.Email.Port), auth, from, []string{to}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

