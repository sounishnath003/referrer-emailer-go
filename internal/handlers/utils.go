package handlers

import (
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