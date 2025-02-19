package core

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/smtp"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"cloud.google.com/go/vertexai/genai"
	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
)

type CoreOpts struct {
	Port       int
	MailAddr   string
	MailSecret string
	SmtpAddr   string
	MongoDbUri string

	ModelName        string
	GcpProjectID     string
	GcpLocation      string
	GcpStorageBucket string
}

// Core defines the core construct of the service.
type Core struct {
	Port int
	opts *CoreOpts

	smtpAuth      smtp.Auth
	storageClient *storage.Client
	llm           *genai.GenerativeModel
	DB            *repository.MongoDBClient
	Lo            *slog.Logger
}

// configureIndexesDB helps to configure database level constraints and checks.
// If complaints and rules on the collections as per the defined actions to make data integrity stronger.
func (co *Core) configureIndexesDB() {
	co.configureUsersIndexes()
}

func NewCore(opts *CoreOpts) *Core {

	co := &Core{
		opts: opts,
		Port: opts.Port,
		Lo:   slog.Default(),
	}

	// Initialize the database
	mdb, err := intiializeDatabase(co.opts.MongoDbUri)
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
		co.opts.MailAddr,
		co.opts.MailSecret,
		co.opts.SmtpAddr,
	)

	// Configure indexes managements
	go co.configureIndexesDB()

	// Initialize LLM model (Gemini).
	if err := co.initializeLLM(); err != nil {
		panic(err)
	}

	// Initialize Storage Client (GCS Bucket).
	if err := co.initializeGCSClient(); err != nil {
		co.Lo.Error("error occured during GCS storage client inisialization:", "error", fmt.Errorf("Unable to create GCS storage client: %w\n", err))
		panic(err)
	}

	return co
}

// initializeGCSClient initializes a new Google Cloud Storage (GCS) client and assigns it to the Core struct.
// It sets a context with a timeout of 10 seconds for the client creation process.
// If the client creation fails, it returns an error with a descriptive message.
// The function returns nil if the client is successfully created.
func (co *Core) initializeGCSClient() error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	storageClient, err := storage.NewClient(ctx, option.WithServiceAccountFile("/Users/sounishnath/sounish-cloud-workstation-ac143dfffa26.json"))

	if err != nil {
		return fmt.Errorf("Unable to create GCS storage client: %w\n", err)
	}

	co.storageClient = storageClient

	return nil
}

func (co *Core) UploadFileToGCSBucket(file *multipart.FileHeader) (string, error) {

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %w\n", err)
	}
	defer src.Close()

	ctx, cancel := getContextWithTimeout(10) // waits for 10 seconds
	defer cancel()

	objectName := fmt.Sprintf("referrer-uploads/%s", file.Filename)
	wc := co.storageClient.Bucket(co.opts.GcpStorageBucket).Object(objectName).NewWriter(ctx)

	if _, err = io.Copy(wc, src); err != nil {
		return "", fmt.Errorf("failed to upload resume to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close the GCS writer: %w", err)
	}

	return fmt.Sprintf("gs://%s/%s", co.opts.GcpStorageBucket, objectName), nil
}
