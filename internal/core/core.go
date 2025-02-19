package core

import (
	"log/slog"
	"net/smtp"

	"cloud.google.com/go/vertexai/genai"
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

	modelName    string
	gcpProjectID string
	gcpLocation  string
	llm          *genai.GenerativeModel

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
		smtpAddr:     "smtp.gmail.com",
		Port:         utils.GetNumberFromEnv("PORT", 3000),
		mailAddr:     utils.GetStringFromEnv("MAIL_ADDR", "flock.sinasini@gmail.com"),
		mailSecret:   utils.GetStringFromEnv("MAIL_SECRET", "P@55w0Rd5!"),
		mongoDbUri:   utils.GetStringFromEnv("MONGO_DB_URI", "localhost"),
		gcpProjectID: utils.GetStringFromEnv("GCP_PROJECT_ID", "sounish-cloud-workstation"),
		gcpLocation:  utils.GetStringFromEnv("GCP_PROJECT_LOCATION", "asia-south1"),
		modelName:    utils.GetStringFromEnv("GCP_VERTEX_AI_LLM", "gemini-1.5-flash-002"),

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

	// Initialize the SMTP Auth instance to be reused.
	co.smtpAuth = smtp.PlainAuth(
		"",
		co.mailAddr,
		co.mailSecret,
		co.smtpAddr,
	)

	// Configure indexes managements
	go co.configureIndexesDB()

	// Initialize LLM model (Gemini).
	if err := co.initializeLLM(); err != nil {
		panic(err)
	}

	return co
}
