package handlers

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sounishnath003/customgo-mailer-service/internal/core"
)

// Declare custom context.
type HandlerContext struct {
	echo.Context
	Co *core.Core
}

// GetCore handles to get the *core options for the application
func (hc *HandlerContext) GetCore() *core.Core {
	return hc.Co
}

// getEmailIDFromToken get the user's email from JWT secret key from the header
func getEmailIDFromToken(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["email"].(string)
}
