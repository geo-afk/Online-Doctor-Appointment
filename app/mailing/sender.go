package mailing

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"

	_ "github.com/joho/godotenv/autoload"
)

// EmailData holds the information to pass into the template
type EmailData struct {
	Name      string
	ResetLink string
}

// SMTP Config (adjust this)
var (
	smtpHost       = os.Getenv("SMTP_HOST")
	smtpPort       = os.Getenv("SMTP_PORT")
	senderEmail    = os.Getenv("SMTP_EMAIL")
	senderPassword = os.Getenv("SMTP_PASSWORD")
)

func verify() {

	resetToken := "abc123token"
	resetURL := fmt.Sprintf("https://localhost:%s/reset-password?token=%s", os.Getenv("PORT"), resetToken)

	data := EmailData{
		Name:      "John Doe",
		ResetLink: resetURL,
	}

	tmplPath := filepath.Join("templates", "password_reset.tmpl")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		panic(err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		panic(err)
	}

	to := []string{"xobimis971@miracle3.com"}
	msg := []byte("From: " + senderEmail + "\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: Password Reset Request\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body.String())

	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, to, msg)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Password reset email sent successfully!")
}
