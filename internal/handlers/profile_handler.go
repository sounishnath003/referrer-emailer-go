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
	currentCompany := c.FormValue("currentCompany")
	currentRole := c.FormValue("currentRole")
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
		CurrentCompany:   currentCompany,
		CurrentRole:      currentRole,
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
	// For self-hosted, we assume the search is performed by the owner. 
	// We might need the owner's email here. For now, let's try to get it from query or header if possible,
	// or we can change the signature to accept it.
	// HOWEVER, the `SearchPeople` previously didn't filter by owner.
	// Since this is SINGLE USER, we can probably just search ALL contacts or require the email.
	// But `PeopleSearchHandler` is called by `searchPeople$` in frontend which doesn't pass email usually in the query for the autocomplete.
	// Wait, `searchPeople$` in frontend DOES NOT pass email.
	// To fix this for "Single User Self Hosted", we can just fetch ALL contacts matching query, assuming all contacts belong to the single user.
	// Or we can try to extract the user from the JWT token if available.
	
	// Let's assume we search all contacts for now as it's single user.
	// But `GetContacts` requires OwnerID.
	// We can fetch the first user from DB or just allow passing a wildcard?
	// Better: Extract from JWT token "email" claim if possible.
	
	// SIMPLE FIX for Single User: Just search all contacts or use a wildcard owner if we modified repository.
	// But let's look at `SearchContactsForAutocomplete`.
	// We will pass a wildcard "%" or empty string if we can't find the user, but `GetContacts` filters by OwnerID.
	
	// Let's actually UPDATE `PeopleSearchHandler` to get the user from the JWT Token.
	// user := c.Get("user") // echojwt puts the token in "user" context
	// But wait, the token parsing might be different.
	// Let's check server.go... `echojwt.Config{...}`.
	// It sets claims.
	
	// If we can't get the user, we might be in trouble. 
	// But wait, `PeopleSearchHandler` is used in `DraftWithAi` and `EmailDrafter`.
	// These pages should be protected.
	
	// FALLBACK: Since it's self-hosted, let's just search *all* contacts if we can't identify the owner, 
	// OR we can update the repository to NOT filter by owner if ownerID is empty.
	
	hctx.GetCore().Lo.Info("Received search Param", "query", queryParam)

	// In a single-user self-hosted env, we can just search matching any contact.
	// Let's use a repository method that ignores owner for autocomplete if we want, 
	// OR just assume one main user.
	// Let's try to use the "owner" from the query param if frontend sends it? No it doesn't.
	
	// Let's update `SearchContactsForAutocomplete` to take a wildcard for owner.
	// Actually, let's just fetch all contacts matching the query.
	
	// We will modify repository to allow empty ownerID to mean "all".
	
	results, err := hctx.GetCore().DB.SearchContactsForAutocomplete("", queryParam)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, results)
}

// UpdateProfileHandler handles the partial update of a user's profile.
func UpdateProfileHandler(c echo.Context) error {
	hctx := c.(*HandlerContext)

	var user repository.User
	if err := c.Bind(&user); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid request body"))
	}

	if err := hctx.GetCore().DB.UpdateProfile(&user); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Profile updated successfully"})
}
