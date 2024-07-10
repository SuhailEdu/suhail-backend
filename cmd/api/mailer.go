package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"github.com/SuhailEdu/suhail-backend/internal/database/schema"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/wneessen/go-mail"
	"html/template"
	"time"
)

func sendWelcomeEmail(user schema.User, c context.Context, config *Config) {

	url, err := createVerificationUrl(user, c, config)
	if err != nil {
		config.logger.Fatal(err)
	}

	data := struct {
		User schema.User
		Url  string
	}{
		User: user,
		Url:  url,
	}

	tmpl := template.Must(template.ParseFiles("mail_templates/welcome_email.tmpl"))

	var parsedTmpl bytes.Buffer

	err = tmpl.Execute(&parsedTmpl, data)
	if err != nil {
		config.logger.Fatal(err)
	}

	m := mail.NewMsg()
	if err := m.From("no-reply@rafiqi-uk.com"); err != nil {
		config.logger.Fatalf("failed to set From address: %s", err)
		return
	}
	if err := m.To(user.Email); err != nil {
		config.logger.Fatalf("failed to set To address: %s", err)
		return
	}

	m.Subject("مرحبا في في سهيل !")
	m.SetBodyString(mail.TypeTextHTML, parsedTmpl.String())

	client, err := mail.NewClient("smtp.zoho.com", mail.WithSSL(), mail.WithPort(465), mail.WithSMTPAuth(mail.SMTPAuthLogin),
		mail.WithUsername("no-reply@rafiqi-uk.com"), mail.WithPassword("rqbTaPpJdBwM1R9A@"))
	if err != nil {

		config.logger.Fatalf("failed to create mail client: %s", err)
	}
	if err := client.DialAndSend(m); err != nil {
		config.logger.Fatalf("failed to send mail: %s", err)
	}
}

func createVerificationUrl(user schema.User, c context.Context, config *Config) (string, error) {

	randomBytes := make([]byte, 16)

	// Create empty byte string
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	plainText := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(plainText))

	_, err = config.db.StoreVerificationCode(c, schema.StoreVerificationCodeParams{
		Hash:   hash[:],
		UserID: user.ID,
		Expiry: pgtype.Timestamp(pgtype.Timestamptz{Time: time.Now().Add(time.Hour * 2), Valid: true}),
		Scope:  "login",
	})

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://localhost:3000/auth/verify?token=%s", plainText), nil

}
