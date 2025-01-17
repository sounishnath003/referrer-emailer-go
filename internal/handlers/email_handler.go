package handlers

import (
	"net/http"

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

	err = hctx.GetCore().InvokeSendMail(emailSenderDto.From, emailSenderDto.To, emailSenderDto.Sub, emailSenderDto.Body)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, emailSenderDto)
}
