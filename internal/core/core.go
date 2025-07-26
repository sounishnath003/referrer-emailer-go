package core

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"

	"cloud.google.com/go/vertexai/genai"
	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
)

type CoreOpts struct {
	Port          int
	MailAddr      string
	MailSecret    string
	SmtpAddr      string
	MongoDbUri    string
	PdfServiceUri string

	ModelName        string
	GcpProjectID     string
	GcpLocation      string
	GcpStorageBucket string
}

// Core defines the core construct of the service.
type Core struct {
	Port          int
	DB            *repository.MongoDBClient
	Lo            *slog.Logger
	PdfServiceUri string

	opts          *CoreOpts
	smtpAuth      smtp.Auth
	storageClient *storage.Client
	llm           *genai.GenerativeModel
}

// configureIndexesDB helps to configure database level constraints and checks.
// If complaints and rules on the collections as per the defined actions to make data integrity stronger.
func (co *Core) configureIndexesDB() {
	co.createIndexHelper("users", "email", true)
	co.createIndexHelper("job_queues", "userEmailAddress", false)
	co.createIndexHelper("referral_mailbox", "from", false)
	co.createIndexHelper("referral_mailbox", "createdAt", false)
	co.createIndexHelper("ai_email_drafts", "userEmailAddress", false)
	co.createIndexHelper("ai_email_drafts", "from", false)
	co.createIndexHelper("ai_email_drafts", "companyName", false)
}

func NewCore(opts *CoreOpts) *Core {

	co := &Core{
		opts:          opts,
		Port:          opts.Port,
		Lo:            slog.Default(),
		PdfServiceUri: opts.PdfServiceUri,
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
		co.Lo.Error("error occured during GCS storage client inisialization:", "error", fmt.Errorf("unable to create GCS storage client: %w", err))
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

	storageClient, err := storage.NewClient(ctx)

	if err != nil {
		return fmt.Errorf("unable to create GCS storage client: %w", err)
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
		return "", fmt.Errorf("error opening file: %w", err)
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

// DownloadObjectFromGCSBucket downloads an object from GCS and stores it in the `storage` directory.
// It returns the local file path of the downloaded object.
func (co *Core) DownloadObjectFromGCSBucket(objectAddress string) (string, error) {
	// Parse the object address to get the bucket name and object name
	bucketName, objectName, err := parseGCSObjectAddress(objectAddress)
	if err != nil {
		return "", fmt.Errorf("invalid GCS object address: %w", err)
	}

	// Create the storage directory if it doesn't exist
	storageDir := "storage"
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Create the local file path
	localFilePath := filepath.Join(storageDir, filepath.Base(objectName))

	// Open the local file for writing
	localFile, err := os.Create(localFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer localFile.Close()

	// Get the object from GCS
	ctx, cancel := getContextWithTimeout(50) // waits for 50 seconds
	defer cancel()

	rc, err := co.storageClient.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create GCS object reader: %w", err)
	}
	defer rc.Close()

	// Copy the object data to the local file
	if _, err := io.Copy(localFile, rc); err != nil {
		return "", fmt.Errorf("failed to copy GCS object to local file: %w", err)
	}

	return localFilePath, nil
}

// parseGCSObjectAddress parses a GCS object address and returns the bucket name and object name.
func parseGCSObjectAddress(objectAddress string) (string, string, error) {
	// Example object address: gs://bucket-name/object-name
	if !strings.HasPrefix(objectAddress, "gs://") {
		return "", "", fmt.Errorf("invalid GCS object address: %s", objectAddress)
	}

	parts := strings.SplitN(objectAddress[5:], "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid GCS object address: %s", objectAddress)
	}

	return parts[0], parts[1], nil
}

func (co *Core) SubmitResumeToJobQueue(userEmailAddress, resumeGCSPath string) error {
	// Get context.
	ctx, cancel := getContextWithTimeout(10)
	defer cancel()

	collection := co.DB.Database("referrer").Collection("job_queues")

	// Create a job
	job := repository.JobQueue{
		UserEmailAddress: userEmailAddress,
		JobType:          repository.EXTRACT_CONTENT,
		Status:           "PENDING",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		Payload: repository.Payload{
			ResumeURL: resumeGCSPath,
		},
	}
	_, err := collection.InsertOne(ctx, job)
	if err != nil {
		return fmt.Errorf("not able to process resume: %w", err)
	}

	return nil
}
