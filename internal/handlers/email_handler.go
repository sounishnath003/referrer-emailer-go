package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
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

	// Create a context with a timeout of 5 seonds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Channel to receive the result of the email sending
	resultChan := make(chan error, 1)

	// Start a goroutine to send the email
	go func() {
		resultChan <- hctx.GetCore().InvokeSendMail(emailSenderDto.From, emailSenderDto.To, emailSenderDto.Sub, emailSenderDto.Body)
	}()

	select {
	case <-ctx.Done():
		// Context timeout
		return SendErrorResponse(c, http.StatusRequestTimeout, ctx.Err())
	case err := <-resultChan:
		if err != nil {
			return SendErrorResponse(c, http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, emailSenderDto)
}
