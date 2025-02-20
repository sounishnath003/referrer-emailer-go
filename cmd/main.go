package main

import (
	"github.com/sounishnath003/customgo-mailer-service/internal/core"
	"github.com/sounishnath003/customgo-mailer-service/internal/server"
	"github.com/sounishnath003/customgo-mailer-service/internal/utils"
)

func main() {

	co := core.NewCore(&core.CoreOpts{
		Port:             utils.GetNumberFromEnv("PORT", 3000),
		SmtpAddr:         "smtp.gmail.com",
		MailAddr:         utils.GetStringFromEnv("MAIL_ADDR", "flock.sinasini@gmail.com"),
		MailSecret:       utils.GetStringFromEnv("MAIL_SECRET", "P@55w0Rd5!"),
		MongoDbUri:       utils.GetStringFromEnv("MONGO_DB_URI", "localhost"),
		GcpProjectID:     utils.GetStringFromEnv("GCP_PROJECT_ID", "sounish-cloud-workstation"),
		GcpLocation:      utils.GetStringFromEnv("GCP_PROJECT_LOCATION", "asia-south1"),
		ModelName:        utils.GetStringFromEnv("GCP_VERTEX_AI_LLM", "gemini-1.5-flash-002"),
		GcpStorageBucket: utils.GetStringFromEnv("GCP_STORAGE_BUCKET", "sounish-cloud-workstation"),

		Conucrrency: 5,
	})

	server := server.NewServer(co)
	panic(server.Start())
}
