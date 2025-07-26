package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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

	dstPath, err := hctx.GetCore().UploadFileToGCSBucket(file)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	// Push the resume update as new job
	if err = hctx.GetCore().SubmitResumeToJobQueue(email, dstPath); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	profileInfo := &repository.User{
		Firstname:        firstName,
		LastName:         lastName,
		Resume:           dstPath,
		About:            about,
		Email:            email,
		Country:          country,
		ProfileSummary:   "",
		ExtractedContent: "",
		Notification:     notification,
	}

	err = hctx.GetCore().DB.UpdateProfileInformation(profileInfo)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	// Return response
	return c.JSON(http.StatusOK, map[string]string{"message": "Profile information updated successfully"})

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

// PeopleSearchHandler helps to find out persons/users who belongs to certain filters.
// Filters like companyName, Email
func PeopleSearchHandler(c echo.Context) error {
	// Get the core
	hctx := c.(*HandlerContext)

	queryParam := c.QueryParam("query")
	hctx.GetCore().Lo.Info("Received search Param", "query", queryParam)

	emails, err := hctx.GetCore().DB.SearchPeople(queryParam)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string][]string{"users": emails})
}
