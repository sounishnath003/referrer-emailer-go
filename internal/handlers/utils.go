package handlers

import (
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
)

// SendErrorResponse sends the API Errors in JSON formatted standard.
func SendErrorResponse(c echo.Context, statusCode int, err error) error {
	return c.JSON(statusCode, map[string]any{
		"statusCode": statusCode,
		"error":      err.Error(),
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

// isValidEmail checks whether the string s contains any match of the regular expression email.
func isValidEmail(email string) bool {
    var emailRgx = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return emailRgx.MatchString(email)
}