package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
)

// ProfileInformationHandler handlers the submission of information
// It will have a upload resume button which will upload the resume also.
func ProfileInformationHandler(c echo.Context) error {
	// Get core.
	hctx := c.(*HandlerContext)

	// Parse form data
	firstName := c.FormValue("firstName")
	lastName := c.FormValue("lastName")
	about := c.FormValue("about")
	email := c.FormValue("email")
	country := c.FormValue("country")
	notificationsStr := c.FormValue("notifications")

	// Parse the notification into the repository.notification struct
	var notification repository.Notification
	err := json.Unmarshal([]byte(notificationsStr), &notification)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("notification json parsing error"))
	}

	// Handle file upload functionality
	file, err := c.FormFile("resume")
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to upload resume"))
	}

	// Check the file type
	if err = isFilePDF(file); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	// Check the file size limit
	if err = isFileUnder2MB(file); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 2)
	dstPathChan := make(chan string, 1)

	// Goroutine to save the resume to disk
	wg.Add(1)
	go func() {
		defer wg.Done()
		dstPath, shouldReturn, err := saveResumeIntoDisk(file, c)
		if shouldReturn {
			errChan <- err
			return
		}
		dstPathChan <- dstPath

		err = hctx.GetCore().ResumeParser(dstPath)
		if err != nil {
			errChan <- err
		}
	}()

	// Goroutine to update profile information in the database
	wg.Add(1)
	go func() {
		defer wg.Done()
		dstPath := <-dstPathChan
		profileInfo := &repository.User{
			Firstname:    firstName,
			LastName:     lastName,
			Resume:       dstPath,
			About:        about,
			Email:        email,
			Country:      country,
			Notification: notification,
		}

		err := hctx.GetCore().DB.UpdateProfileInformation(profileInfo)
		if err != nil {
			errChan <- err
		}
	}()

	// Wait for both goroutines to complete
	wg.Wait()
	close(errChan)
	close(dstPathChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return SendErrorResponse(c, http.StatusInternalServerError, err)
		}
	}

	// Return response
	return c.JSON(http.StatusOK, map[string]string{"message": "Profile information updated successfully"})

}

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

// GetProfileHandler helps to get the information from the provided email from query params.
func GetProfileHandler(c echo.Context) error {
	// Get core
	hctx := c.(*HandlerContext)

	email := c.QueryParam("email")
	hctx.GetCore().Lo.Info("is valid email", "email", email, "valid", isValidEmail(email))
	if len(email) == 0 || !isValidEmail(email) {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("provide valid email address"))
	}

	u, err := hctx.GetCore().DB.GetProfileByEmail(email)
	// Make sure you omit the user's password
	u.Password = ""
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, u)
}
