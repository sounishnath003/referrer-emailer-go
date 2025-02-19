package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
)

// SendErrorResponse sends the API Errors in JSON formatted standard.
func SendErrorResponse(c echo.Context, statusCode int, err error) error {
	return c.JSON(statusCode, map[string]any{
		"statusCode": statusCode,
		"error":      err.Error(),
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

// isValidEmail checks whether the string s contains any match of the regular expression email.
func isValidEmail(email string) bool {
	var emailRgx = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRgx.MatchString(email)
}

// saveResumeIntoDisk saves an uploaded file to the disk.
//
// Parameters:
//   - file: A pointer to the multipart.FileHeader representing the uploaded file.
//   - c: The echo.Context for the current request context.
//
// Returns:
//   - string: The path where the file is saved.
//   - bool: A boolean indicating if an error occurred.
//   - error: An error object if an error occurred, otherwise nil.
//
// The function performs the following steps:
//  1. Opens the uploaded file.
//  2. Creates a destination file in the "storage" directory with the same filename.
func saveResumeIntoDisk(file *multipart.FileHeader, c echo.Context) (string, bool, error) {
	src, err := file.Open()
	if err != nil {
		return "", true, SendErrorResponse(c, http.StatusInternalServerError, err)
	}
	defer src.Close()

	// Create destination file
	dstPath := filepath.Join("storage", file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", true, SendErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("Failed to upload resume. Something went wrong."))
	}
	defer dst.Close()

	// Copy file content
	if _, err = io.Copy(dst, src); err != nil {
		return "", true, SendErrorResponse(c, http.StatusInternalServerError, err)
	}
	return dstPath, false, nil
}

// isFileUnder2MB checks the file size under 2MB
func isFileUnder2MB(file *multipart.FileHeader) error {
	const MAX_LIMIT = (2 * 1024 * 1024) // 2 MB
	if file.Size > MAX_LIMIT {
		return fmt.Errorf("file exceeds 2 MB limit.")
	}
	return nil
}

// isFilePDF check the file type of the upload file object
func isFilePDF(file *multipart.FileHeader) error {
	if file.Header.Get("Content-Type") != "application/pdf" {
		return fmt.Errorf("only pdf files are allowed.")
	}
	return nil
}
