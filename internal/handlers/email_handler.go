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

	// Create a context with a timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Channel to receive the result of the email sending
	errChan := make(chan error, 1)

	go func() {
		var localDst string

		if emailSenderDto.TailoredResumeID != "" {
			// Convert string to ObjectID
			objID, objErr := primitive.ObjectIDFromHex(emailSenderDto.TailoredResumeID)
			if objErr != nil {
				errChan <- fmt.Errorf("invalid tailoredResumeId: %w", objErr)
				return
			}
			// Fetch tailored resume and generate PDF
			tr, tErr := hctx.GetCore().DB.GetTailoredResumeByID(ctx, objID)
			if tErr != nil {
				errChan <- tErr
				return
			}
			// Convert Markdown to HTML using goldmark
			var htmlBuf bytes.Buffer
			if err := goldmark.Convert([]byte(tr.ResumeMarkdown), &htmlBuf); err != nil {
				errChan <- fmt.Errorf("failed to convert markdown to HTML: %w", err)
				return
			}
			htmlContent := htmlBuf.String()
			// Call PDF service to generate PDF from HTML
			pdfResp, pdfErr := generatePDFfromResume(hctx.GetCore().PdfServiceUri, htmlContent)
			if pdfErr != nil {
				errChan <- pdfErr
				return
			}
			// Save PDF to a temp file
			tmpFile, tmpErr := os.CreateTemp("/tmp", "Sounish_Naths_Resume_*.pdf")
			if tmpErr != nil {
				errChan <- tmpErr
				return
			}
			defer tmpFile.Close()
			_, wErr := tmpFile.Write(pdfResp)
			if wErr != nil {
				errChan <- wErr
				return
			}
			localDst = tmpFile.Name()
		} else {
			// get the `from` user from DB
			u, err := hctx.GetCore().DB.GetProfileByEmail(emailSenderDto.From)
			if err != nil {
				errChan <- err
				return
			}
			// Download the file into LocalDisk
			localDst, err = hctx.GetCore().DownloadObjectFromGCSBucket(u.Resume)
			if err != nil {
				errChan <- err
				return
			}
		}

		// Invoke the SendMail with Attachment.
		errChan <- hctx.GetCore().InvokeSendMailWithAttachment(emailSenderDto.From, emailSenderDto.To, emailSenderDto.Sub, emailSenderDto.Body, localDst)

		// Purge the Local file
		if err := os.Remove(localDst); err != nil {
			hctx.GetCore().Lo.Error("error deleting the file:", "error", err)
		}
	}()

	select {
	case <-ctx.Done():
		// Context timeout
		return SendErrorResponse(c, http.StatusRequestTimeout, ctx.Err())
	case err := <-errChan:
		if err != nil {
			return SendErrorResponse(c, http.StatusInternalServerError, err)
		}
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

	if len(userEmail) > 0 && !isValidEmail(userEmail) {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid email or no email found"))
	}

	emails, err := hctx.GetCore().DB.GetLatestEmailsByFilter(bson.M{
		"$or": []bson.M{
			bson.M{"from": userEmail},
			bson.M{"uuid": emailUuid},
		},
	})
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("failed to fetch mails: %w", err))
	}
	return c.JSON(http.StatusOK, emails)
}
