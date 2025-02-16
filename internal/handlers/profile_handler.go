package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type NotificationType struct {
	Offers           bool   `json:"offers"`
	PushNotification string `json:"pushNotifications"`
	ReceiveEmails    bool   `json:"receiveEmails"`
}

type ProfileInformationRequestDto struct {
	Firstname    string           `json:"firstName"`
	LastName     string           `json:"lastName"`
	Resume       string           `json:"resume"`
	About        string           `json:"about"`
	Email        string           `json:"email"`
	Notification NotificationType `json:"notifications"`
}

// ProfileInformationHandler handlers the submission of information
// It will have a upload resume button which will upload the resume also.
func ProfileInformationHandler(c echo.Context) error {
	// Parse form data
	firstName := c.FormValue("firstName")
	lastName := c.FormValue("lastName")
	about := c.FormValue("about")
	email := c.FormValue("email")

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
	if err = isFileUnder1MB(file); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	dstPath, shouldReturn, err := saveResumeIntoDisk(file, c)
	if shouldReturn {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	// Create response DTO
	profileInfo := ProfileInformationRequestDto{
		Firstname: firstName,
		LastName:  lastName,
		Resume:    dstPath,
		About:     about,
		Email:     email,
	}

	// Return response
	return c.JSON(http.StatusOK, profileInfo)

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

func isFileUnder1MB(file *multipart.FileHeader) error {
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
