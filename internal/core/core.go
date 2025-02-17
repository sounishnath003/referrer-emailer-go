package core

import (
	"context"
	"log/slog"
	"net/smtp"
	"time"

	"github.com/sounishnath003/customgo-mailer-service/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Core defines the core construct of the service.
type Core struct {
	Port       int
	mailAddr   string
	mailSecret string
	smtpAddr   string
	mongoDbUri string
	smtpAuth   smtp.Auth

	DB *mongo.Client
	Lo *slog.Logger
}

// initDB tries to connect with mongo DB database within 10 second context
func initDB(dbUri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	slog.Info("mongo db client instance ping checkked and connected")
	return client, nil
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
	mdb, err := initDB(co.mongoDbUri)
	if err != nil {
		co.Lo.Error("not able to connect to mongoDB", slog.Any("mdb_err", err.Error()))
		panic(err)
	}
	co.DB = mdb

	// Initialize the SMTP Auth instance to be reused.
	co.smtpAuth = smtp.PlainAuth(
		"",
		co.mailAddr,
		co.mailSecret,
		co.smtpAddr,
	)

	return co
}
