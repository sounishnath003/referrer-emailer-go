package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

// EmailSenderDto hold the DTO for the email sending data payload.
type EmailSenderDto struct {
	From string   `json:"from"`
	To   []string `json:"to"`
	Sub  string   `json:"subject"`
	Body string   `json:"body"`
}

// SendEmailHandler handlers will handle the email sending capability to the multiple users
func SendEmailHandler(c echo.Context) error {
	var emailSenderDto EmailSenderDto

	err := c.Bind(&emailSenderDto)
	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	hctx := c.(*HandlerContext)

	// Create a context with a timeout of 10 seonds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Channel to receive the result of the email sending
	errChan := make(chan error, 1)

	// Start a goroutine to send the email
	go func() {
		// resultChan <- hctx.GetCore().InvokeSendMail(emailSenderDto.From, emailSenderDto.To, emailSenderDto.Sub, emailSenderDto.Body)

		// get the `from` user from DB
		u, err := hctx.GetCore().DB.GetProfileByEmail(emailSenderDto.From)
		if err != nil {
			errChan <- err
		}

		// Download the file into LocalDisk
		localDst, err := hctx.GetCore().DownloadObjectFromGCSBucket(u.Resume)
		if err != nil {
			errChan <- err
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

func GetReferralEmailsHandler(c echo.Context) error {
	// Get context
	hctx := c.(*HandlerContext)

	userEmail := c.QueryParam("email")
	emailUuid := c.QueryParam("uuid")

	if len(userEmail) == 0 || !isValidEmail(userEmail) {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid email or no email found."))
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
