package core

import (
	"log/slog"
	"net/smtp"

	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
	"github.com/sounishnath003/customgo-mailer-service/internal/utils"
)

// Core defines the core construct of the service.
type Core struct {
	Port       int
	mailAddr   string
	mailSecret string
	smtpAddr   string
	mongoDbUri string
	smtpAuth   smtp.Auth

	DB *repository.MongoDBClient
	Lo *slog.Logger
}

// configureIndexesDB helps to configure database level constraints and checks.
// If complaints and rules on the collections as per the defined actions to make data integrity stronger.
func (co *Core) configureIndexesDB() {
	co.configureUsersIndexes()
}

func NewCore() *Core {

	co := &Core{
		smtpAddr:   "smtp.gmail.com",
		Port:       utils.GetNumberFromEnv("PORT", 3000),
		mailAddr:   utils.GetStringFromEnv("MAIL_ADDR", "flock.sinasini@gmail.com"),
		mailSecret: utils.GetStringFromEnv("MAIL_SECRET", "P@55w0Rd5!"),
		mongoDbUri: utils.GetStringFromEnv("MONGO_DB_URI", "localhost"),

		Lo: slog.Default(),
	}

	// Initialize the database
	mdb, err := intiializeDatabase(co.mongoDbUri)
	if err != nil {
		co.Lo.Error("not able to connect to mongoDB", slog.Any("mdb_err", err.Error()))
		panic(err)
	}
	co.DB = &repository.MongoDBClient{
		Client: mdb,
	}

	// Configure indexes managements
	go co.configureIndexesDB()

	// Initialize the SMTP Auth instance to be reused.
	co.smtpAuth = smtp.PlainAuth(
		"",
		co.mailAddr,
		co.mailSecret,
		co.smtpAddr,
	)

	return co
}
