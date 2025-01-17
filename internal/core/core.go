package core

import (
	"log/slog"
	"net/smtp"

	"github.com/sounishnath003/customgo-mailer-service/internal/utils"
)

// Core defines the core construct of the service.
type Core struct {
	Port       int
	mailAddr   string
	mailSecret string
	smtpAddr   string
	smtpAuth   smtp.Auth
	Lo         *slog.Logger
}

func NewCore() *Core {

	co := &Core{
		Port:       utils.GetNumberFromEnv("PORT", 3000),
		mailAddr:   utils.GetStringFromEnv("MAIL_ADDR", "flock.sinasini@gmail.com"),
		mailSecret: utils.GetStringFromEnv("MAIL_SECRET", "P@55w0Rd5!"),
		smtpAddr:   "smtp.gmail.com",
		Lo:         slog.Default(),
	}

	// Initialize the SMTP Auth instance to be reused.
	co.smtpAuth = smtp.PlainAuth(
		"",
		co.mailAddr,
		co.mailSecret,
		co.smtpAddr,
	)

	return co
}
