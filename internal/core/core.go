package core

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/smtp"
	"time"

	"cloud.google.com/go/storage"

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

	Conucrrency int
}

// Core defines the core construct of the service.
type Core struct {
	Port int
	DB   *repository.MongoDBClient
	Lo   *slog.Logger

	opts          *CoreOpts
	smtpAuth      smtp.Auth
	storageClient *storage.Client
	llm           *genai.GenerativeModel
	workerPool    *WorkerPool
}

// configureIndexesDB helps to configure database level constraints and checks.
// If complaints and rules on the collections as per the defined actions to make data integrity stronger.
func (co *Core) configureIndexesDB() {
	co.createIndexHelper("users", "email", true)
	co.createIndexHelper("job_queues", "userEmailAddress", false)
	co.createIndexHelper("resumes", "emailAddress", true)
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

	// Initialize worker pool
	// Buffer Queue = 10 x Concurrency
	wp := NewWorkerPool(co.DB, opts.Conucrrency, 10*opts.Conucrrency)
	co.workerPool = wp

	go func() {
		// Start worker pool
		go co.workerPool.StartWorkers()
		// Attach the mongo db client
		co.workerPool.ListenForThePendingJobs()
		// Wait for the execution
		co.workerPool.Wait()
	}()

	return co
}

// initializeGCSClient initializes a new Google Cloud Storage (GCS) client and assigns it to the Core struct.
// It sets a context with a timeout of 10 seconds for the client creation process.
// If the client creation fails, it returns an error with a descriptive message.
// The function returns nil if the client is successfully created.
func (co *Core) initializeGCSClient() error {
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	storageClient, err := storage.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("Unable to create GCS storage client: %w\n", err)
	}

	co.storageClient = storageClient

	return nil
}

// UploadFileToGCSBucket uploads a file to a Google Cloud Storage (GCS) bucket.
// It takes a multipart.FileHeader as input and returns the URL of the uploaded file
// or an error if the upload fails.
//
// Parameters:
//   - file: A pointer to a multipart.FileHeader representing the file to be uploaded.
//
// Returns:
//   - string: The URL of the uploaded file in the GCS bucket.
//   - error: An error if the upload fails.
//
// Example:
//
//	url, err := co.UploadFileToGCSBucket(fileHeader)
//	if err != nil {
//	    log.Fatalf("Failed to upload file: %v", err)
//	}
//	fmt.Printf("File uploaded to: %s\n", url)
func (co *Core) UploadFileToGCSBucket(file *multipart.FileHeader) (string, error) {

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("error opening file: %w\n", err)
	}
	defer src.Close()

	ctx, cancel := getContextWithTimeout(50) // waits for 50 seconds
	defer cancel()

	objectName := fmt.Sprintf("referrer-uploads/%s", file.Filename)
	dstPath := fmt.Sprintf("gs://%s/%s", co.opts.GcpStorageBucket, objectName)
	wc := co.storageClient.Bucket(co.opts.GcpStorageBucket).Object(objectName).NewWriter(ctx)

	if _, err = io.Copy(wc, src); err != nil {
		return "", fmt.Errorf("failed to upload resume to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close the GCS writer: %w", err)
	}

	return dstPath, nil
}

func (co *Core) SubmitResumeToJobQueue(userEmailAddress, resumeGCSPath string) error {
	// Get context.
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := co.DB.Database("referrer").Collection("job_queues")

	// Create a job
	job := repository.JobQueue{
		UserEmailAddress: userEmailAddress,
		JobType:          "EXTRACT_CONTENT",
		Status:           "PENDING",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Payload: repository.Payload{
			ResumeURL: resumeGCSPath,
		},
	}
	_, err := collection.InsertOne(ctx, job)
	if err != nil {
		return fmt.Errorf("not able to process resume: %w\n", err)
	}

	return nil
}
