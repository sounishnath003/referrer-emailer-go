package core

import (
	"log/slog"

	"github.com/sounishnath003/customgo-mailer-service/internal/utils"
)

// Core defines the core construct of the service.
type Core struct {
	Port       int
	mailAddr   string
	mailSecret string
	smtpAddr   string
	Lo         *slog.Logger
}

func NewCore() *Core {

	return &Core{
		Port:       utils.GetNumberFromEnv("PORT", 3000),
		mailAddr:   utils.GetStringFromEnv("MAIL_ADDR", "flock.sinasini@gmail.com"),
		mailSecret: utils.GetStringFromEnv("MAIL_SECRET", "P@55w0Rd5!"),
		smtpAddr:   "smtp.gmail.com",
		Lo:         slog.Default(),
	}
}
