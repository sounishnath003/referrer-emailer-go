package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddContactHandler handles creating a new contact
func AddContactHandler(c echo.Context) error {
	hctx := c.(*HandlerContext)

	// Get logged in user email/id (In this simple self-hosted setup, we use email as ownerID usually, 
	// or we can fetch the user ID. Let's assume we use the email from the token/context if available, 
	// or passed in. For consistency with auth, let's look for the user.)
	
	// NOTE: In the current setup, we might need to extract the user from the token or request.
	// Assuming the frontend sends the 'ownerEmail' or similar, OR we extract it from JWT.
	// Since we are "self-hosted single user", we can also just rely on the passed email 
	// or use a default if we want to be strict.
	// Let's rely on the client sending 'ownerEmail' for now to match the pattern, or extract from context if auth middleware sets it.
	// The auth middleware in server.go sets claims. Let's use that if possible, but for simplicity here 
	// matching the other handlers:
	
	type requestDto struct {
		OwnerEmail string `json:"ownerEmail"`
		Name       string `json:"name"`
		Email      string `json:"email"`
		Company    string `json:"company"`
		Role       string `json:"role"`
		LinkedIn   string `json:"linkedin"`
		Notes      string `json:"notes"`
	}
	var req requestDto
	if err := c.Bind(&req); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	if req.Name == "" || req.OwnerEmail == "" {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("name and ownerEmail are required"))
	}

	contact := &repository.Contact{
		OwnerID:  req.OwnerEmail, // Using Email as OwnerID for simplicity in this project structure
		Name:     req.Name,
		Email:    req.Email,
		Company:  req.Company,
		Role:     req.Role,
		LinkedIn: req.LinkedIn,
		Notes:    req.Notes,
	}

	if err := hctx.GetCore().DB.CreateContact(contact); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, contact)
}

// GetContactsHandler fetches contacts
func GetContactsHandler(c echo.Context) error {
	hctx := c.(*HandlerContext)
	ownerEmail := c.QueryParam("email")
	query := c.QueryParam("query")

	if ownerEmail == "" {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("owner email is required"))
	}

	contacts, err := hctx.GetCore().DB.GetContacts(ownerEmail, query)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, contacts)
}

// DeleteContactHandler deletes a contact
func DeleteContactHandler(c echo.Context) error {
	hctx := c.(*HandlerContext)
	idStr := c.Param("id")
	ownerEmail := c.QueryParam("email")

	if idStr == "" || ownerEmail == "" {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("id and owner email are required"))
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid id"))
	}

	if err := hctx.GetCore().DB.DeleteContact(id, ownerEmail); err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Contact deleted"})
}

// SyncContactsHandler triggers the import of contacts from drafts
func SyncContactsHandler(c echo.Context) error {
	hctx := c.(*HandlerContext)
	ownerEmail := c.QueryParam("email")

	if ownerEmail == "" {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("owner email is required"))
	}

	count, err := hctx.GetCore().DB.SyncContactsFromDrafts(ownerEmail)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": fmt.Sprintf("Successfully imported %d contacts", count),
		"count":   count,
	})
}
