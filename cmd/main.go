package main

import (
	"github.com/sounishnath003/customgo-mailer-service/internal/core"
	"github.com/sounishnath003/customgo-mailer-service/internal/server"
	"github.com/sounishnath003/customgo-mailer-service/internal/utils"
	"github.com/sounishnath003/customgo-mailer-service/internal/workerpool"
)

func main() {

	co := core.NewCore(&core.CoreOpts{
		Port:          utils.GetNumberFromEnv("PORT", 3000),
		SmtpAddr:      "smtp.gmail.com",
		MailAddr:      utils.GetStringFromEnv("MAIL_ADDR", "sounish.nath17@gmail.com"),
		MailSecret:    utils.GetStringFromEnv("MAIL_SECRET", "P@55w0Rd5!"),
		MongoDbUri:    utils.GetStringFromEnv("MONGO_DB_URI", "localhost"),
		PdfServiceUri: utils.GetStringFromEnv("PDF_SERVICE_URI", "http://0.0.0.0:3001"),

		GcpProjectID:     utils.GetStringFromEnv("GCP_PROJECT_ID", "sounish-cloud-workstation"),
		GcpLocation:      utils.GetStringFromEnv("GCP_PROJECT_LOCATION", "asia-south1"),
		ModelName:        utils.GetStringFromEnv("GCP_VERTEX_AI_LLM", "gemini-1.5-flash-002"),
		GcpStorageBucket: utils.GetStringFromEnv("GCP_STORAGE_BUCKET", "sounish-cloud-workstation"),
	})

	// Initialize worker pool
	// Buffer Queue = 10 x Concurrency
	wp := workerpool.NewWorkerPool(co, 5)

	go func() {
		// Start worker pool
		go wp.StartWorkers()
		// Attach the mongo db client
		wp.ListenForThePendingJobs()
		// Wait for the execution
		wp.Wait()
	}()

	server := server.NewServer(co)
	panic(server.Start())
}
