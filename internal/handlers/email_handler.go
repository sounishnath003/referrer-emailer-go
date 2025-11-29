package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/yuin/goldmark"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
)

// EmailSenderDto hold the DTO for the email sending data payload.
type EmailSenderDto struct {
	From             string   `json:"from"`
	To               []string `json:"to"`
	Sub              string   `json:"subject"`
	Body             string   `json:"body"`
	TailoredResumeID string   `json:"tailoredResumeId,omitempty"`
}

// SendEmailHandler handlers will handle the email sending capability to the multiple users
func SendEmailHandler(c echo.Context) error {
	var emailSenderDto EmailSenderDto

	if err := c.Bind(&emailSenderDto); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	hctx := c.(*HandlerContext)

	// Identify real recipients (excluding the sender)
	var recipients []string
	sender := emailSenderDto.From
	for _, to := range emailSenderDto.To {
		if !strings.EqualFold(to, sender) {
			recipients = append(recipients, to)
		}
	}

	// Helper to prepare resume
	prepareResume := func(ctx context.Context) (string, error) {
		if emailSenderDto.TailoredResumeID != "" {
			objID, objErr := primitive.ObjectIDFromHex(emailSenderDto.TailoredResumeID)
			if objErr != nil {
				return "", fmt.Errorf("invalid tailoredResumeId: %w", objErr)
			}
			tr, tErr := hctx.GetCore().DB.GetTailoredResumeByID(ctx, objID)
			if tErr != nil {
				return "", tErr
			}
			var htmlBuf bytes.Buffer
			if err := goldmark.Convert([]byte(tr.ResumeMarkdown), &htmlBuf); err != nil {
				return "", fmt.Errorf("failed to convert markdown to HTML: %w", err)
			}
			pdfResp, pdfErr := generatePDFfromResume(hctx.GetCore().PdfServiceUri, htmlBuf.String())
			if pdfErr != nil {
				return "", pdfErr
			}
			tmpFile, tmpErr := os.CreateTemp("/tmp", "Sounish_Naths_Resume_*.pdf")
			if tmpErr != nil {
				return "", tmpErr
			}
			defer tmpFile.Close()
			if _, wErr := tmpFile.Write(pdfResp); wErr != nil {
				return "", wErr
			}
			return tmpFile.Name(), nil
		} else {
			u, err := hctx.GetCore().DB.GetProfileByEmail(emailSenderDto.From)
			if err != nil {
				return "", err
			}
			return hctx.GetCore().DownloadObjectFromGCSBucket(u.Resume)
		}
	}

	// BULK MODE: >1 recipients
	if len(recipients) > 1 {
		// Create Job
		job := &repository.BulkEmailJob{
			UserEmail:       sender,
			TotalRecipients: len(recipients),
			SentCount:       0,
			Status:          "PENDING",
		}
		if err := hctx.GetCore().DB.CreateBulkEmailJob(job); err != nil {
			return SendErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("failed to create bulk job: %w", err))
		}

		go func() {
			// Background context
			bgCtx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()

			localDst, err := prepareResume(bgCtx)
			if err != nil {
				hctx.GetCore().Lo.Error("failed to prepare resume for bulk send", "error", err)
				hctx.GetCore().DB.FailBulkEmailJob(job.ID, fmt.Sprintf("failed to prepare resume: %s", err.Error()))
				return
			}
			defer os.Remove(localDst)

			sentCount := 0
			for _, recipient := range recipients {
				// Construct To: [recipient, sender]
				currentTo := []string{recipient, sender}
				
				// Send
				err := hctx.GetCore().InvokeSendMailWithAttachment(
					sender, 
					currentTo, 
					emailSenderDto.Sub, 
					emailSenderDto.Body, 
					emailSenderDto.TailoredResumeID, 
					localDst,
				)
				
				if err != nil {
					hctx.GetCore().Lo.Error("failed to send bulk email", "to", recipient, "error", err)
					// We log error but continue sending to others
				} else {
					hctx.GetCore().Lo.Info("bulk email sent", "to", recipient)
					sentCount++
					hctx.GetCore().DB.UpdateBulkEmailJobProgress(job.ID, sentCount)
				}
				
				// Small delay to be nice to SMTP
				time.Sleep(500 * time.Millisecond)
			}
			hctx.GetCore().DB.CompleteBulkEmailJob(job.ID)
		}()

		return c.JSON(http.StatusAccepted, map[string]any{
			"message": fmt.Sprintf("Processing bulk emails to %d recipients in background.", len(recipients)),
			"jobId":   job.ID.Hex(),
		})
	}

	// SINGLE MODE (Existing synchronous logic)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	localDst, err := prepareResume(ctx)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}
	defer os.Remove(localDst)

	err = hctx.GetCore().InvokeSendMailWithAttachment(emailSenderDto.From, emailSenderDto.To, emailSenderDto.Sub, emailSenderDto.Body, emailSenderDto.TailoredResumeID, localDst)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, emailSenderDto)
}

// Helper to call the PDF service (simulate the logic from craft_resume_handler.go)
func generatePDFfromResume(pdfServiceUri, resumeContent string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(http.MethodPost, pdfServiceUri+"/generate-pdf", strings.NewReader(resumeContent))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "text/plain")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[Failed]: PDF Service returned status: %s", resp.Status)
	}
	pdfData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return pdfData, nil
}

func GetReferralEmailsHandler(c echo.Context) error {
	// Get context
	hctx := c.(*HandlerContext)

	userEmail := c.QueryParam("email")
	emailUuid := c.QueryParam("uuid")
	company := c.QueryParam("company")
	
	// Pagination params
	page := 1
	limit := 10
	if p := c.QueryParam("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if l := c.QueryParam("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	if page < 1 { page = 1 }
	if limit < 1 { limit = 10 }
	offset := (page - 1) * limit

	// Date Range params (YYYY-MM-DD)
	startDateStr := c.QueryParam("startDate")
	endDateStr := c.QueryParam("endDate")

	if len(userEmail) > 0 && !isValidEmail(userEmail) {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid email or no email found"))
	}

	filter := bson.M{}
	if emailUuid != "" {
		filter["uuid"] = emailUuid
	} else if userEmail != "" {
		filter["from"] = userEmail
		if company != "" {
			filter["subject"] = bson.M{"$regex": primitive.Regex{Pattern: company, Options: "i"}}
		}
		
		// Add Date Range Filter
		if startDateStr != "" || endDateStr != "" {
			dateFilter := bson.M{}
			if startDateStr != "" {
				if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
					dateFilter["$gte"] = t
				}
			}
			if endDateStr != "" {
				if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
					// Add 23h 59m 59s to end date to include the whole day
					dateFilter["$lte"] = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
				}
			}
			if len(dateFilter) > 0 {
				filter["createdAt"] = dateFilter
			}
		}
	} else {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("either email or uuid must be provided"))
	}

	emails, totalCount, err := hctx.GetCore().DB.GetLatestEmailsByFilter(filter, limit, offset)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to fetch mails: %w", err))
	}
	
	return c.JSON(http.StatusOK, map[string]any{
		"data": emails,
		"meta": map[string]any{
			"total": totalCount,
			"page":  page,
			"limit": limit,
		},
	})
}

func GetBulkEmailJobStatusHandler(c echo.Context) error {
	hctx := c.(*HandlerContext)
	idStr := c.Param("id")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid job id"))
	}

	job, err := hctx.GetCore().DB.GetBulkEmailJob(id)
	if err != nil {
		return SendErrorResponse(c, http.StatusNotFound, fmt.Errorf("job not found"))
	}

	return c.JSON(http.StatusOK, job)
}

