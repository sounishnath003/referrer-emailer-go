package handlers

import (
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
