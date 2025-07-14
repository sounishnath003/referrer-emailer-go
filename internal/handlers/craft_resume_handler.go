package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TailorResumeWithJobDescriptionHandler handles tailoring a resume based on a job description and user extracted content.
func TailorResumeWithJobDescriptionHandler(c echo.Context) error {
	hctx := c.(*HandlerContext)

	type requestDto struct {
		JobDescription string `json:"jobDescription"`
		UserEmail      string `json:"userEmail"`
	}
	var req requestDto
	if err := c.Bind(&req); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}
	if len(req.JobDescription) == 0 || len(req.UserEmail) == 0 {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("jobDescription and userEmail are required"))
	}

	u, err := hctx.GetCore().DB.GetProfileByEmail(req.UserEmail)
	if err != nil || u == nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("unable to fetch user: %w", err))
	}
	if len(u.ExtractedContent) == 0 {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("user has no extracted resume content"))
	}

	resumeMarkdown, err := hctx.GetCore().TailorResumeWithJobDescriptionLLM(req.JobDescription, u.ExtractedContent)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	// Store tailored resume in MongoDB
	tr := &repository.TailoredResume{
		UserID:         u.ID.Hex(),
		JobDescription: req.JobDescription,
		ResumeMarkdown: resumeMarkdown,
	}
	ctx := context.Background()
	insertedID, err := hctx.GetCore().DB.CreateTailoredResume(ctx, tr)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to store tailored resume: %w", err))
	}

	return c.JSON(http.StatusOK, map[string]any{"id": insertedID.Hex()})
}

// GetTailoredResumeByIDHandler fetches a tailored resume by its ID
func GetTailoredResumeByIDHandler(c echo.Context) error {
	hctx := c.(*HandlerContext)
	idStr := c.Param("id")
	if idStr == "" {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("missing tailored resume id"))
	}
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid tailored resume id"))
	}
	ctx := context.Background()
	tr, err := hctx.GetCore().DB.GetTailoredResumeByID(ctx, id)
	if err != nil {
		return SendErrorResponse(c, http.StatusNotFound, fmt.Errorf("tailored resume not found: %w", err))
	}
	return c.JSON(http.StatusOK, tr)
}
