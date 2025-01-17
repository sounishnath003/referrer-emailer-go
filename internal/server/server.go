package server

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sounishnath003/customgo-mailer-service/internal/core"
	"github.com/sounishnath003/customgo-mailer-service/internal/handlers"
)

// InitServer - run a ever blocking server spawn to run the webservice backend.
// Make sure if you have other I/O bounded tasks, are runs on goroutines.
func InitServer(co *core.Core) error {
	e := echo.New()

	// Add context
	// Declare the custom context in the route handler
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &handlers.HandlerContext{
				Context: c, Co: co,
			}
			return next(cc)
		}
	})

	// Add middlewares
	e.Use(middleware.RemoveTrailingSlash())
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Add routes
	e.Add("GET", "health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hellow. Ok!")
	})

	// Add API routes.
	api := e.Group("/api")
	api.Add("POST", "/send-email", handlers.SendEmailHandler)

	return e.Start(fmt.Sprintf(":%d", co.Port))
}
