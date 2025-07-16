package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
		CompanyName    string `json:"companyName"`
		JobRole        string `json:"jobRole"`
	}
	var req requestDto
	if err := c.Bind(&req); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}
	if len(req.JobDescription) == 0 || len(req.UserEmail) == 0 || len(req.CompanyName) == 0 || len(req.JobRole) == 0 {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("jobDescription, userEmail, companyName, and jobRole are required"))
	}

	u, err := hctx.GetCore().DB.GetProfileByEmail(req.UserEmail)
	if err != nil || u == nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("unable to fetch user: %w", err))
	}
	if len(u.ExtractedContent) == 0 {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("user has no extracted resume content"))
	}

	resumeMarkdown, err := hctx.GetCore().TailorResumeWithJobDescriptionLLM(req.JobDescription, u.ExtractedContent, req.CompanyName, req.JobRole)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	// Store tailored resume in MongoDB
	tr := &repository.TailoredResume{
		UserID:         u.ID.Hex(),
		JobDescription: req.JobDescription,
		ResumeMarkdown: resumeMarkdown,
		CompanyName:    req.CompanyName,
		JobRole:        req.JobRole,
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

// UpdateTailoredResumeHandler updates the resumeMarkdown of a tailored resume by ID
func UpdateTailoredResumeHandler(c echo.Context) error {
	hctx := c.(*HandlerContext)
	type updateDto struct {
		ID             string `json:"id"`
		ResumeMarkdown string `json:"resumeMarkdown"`
	}
	var req updateDto
	if err := c.Bind(&req); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}
	if req.ID == "" || req.ResumeMarkdown == "" {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("id and resumeMarkdown are required"))
	}
	id, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid id"))
	}
	ctx := context.Background()
	err = hctx.GetCore().DB.UpdateTailoredResumeMarkdown(ctx, id, req.ResumeMarkdown)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]any{"success": true})
}

// Generate the PDF using NodeJS Puppeteer service
func GeneratePDFHandler(c echo.Context) error {
	// Http client with 10 second time out
	client := &http.Client{Timeout: 10 * time.Second}
	// Take out the resumeContent from JSON request body
	var reqBody struct {
		ResumeContent string `json:"resumeContent"`
	}
	if err := c.Bind(&reqBody); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}
	// NOTE: The PDF service expects HTML in the body as plain text
	req, err := http.NewRequest(http.MethodPost, "http://0.0.0.0:3001/generate-pdf", strings.NewReader(reqBody.ResumeContent))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("[Failed]: PDF Service returned status: %s", resp.Status)
	}

	// Read the PDF buffer
	pdfData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return c.Blob(http.StatusOK, "application/pdf", pdfData)
}
