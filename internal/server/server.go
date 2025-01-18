package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sounishnath003/customgo-mailer-service/internal/core"
	"github.com/sounishnath003/customgo-mailer-service/internal/handlers"
)

type Server struct {
	co *core.Core
}

func NewServer(co *core.Core) *Server {
	return &Server{
		co: co,
	}
}

// Start - run a ever blocking server spawn to run the webservice backend.
// Make sure if you have other I/O bounded tasks, are runs on goroutines.
func (s *Server) Start() error {
	e := echo.New()

	// Add context
	// Declare the custom context in the route handler
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &handlers.HandlerContext{
				Context: c, Co: s.co,
			}
			return next(cc)
		}
	})

	// Add middlewares
	e.Use(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200", "localhost:4200"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Content-Type", "Content-Length", "Authorization", "X-API-TrackerId"},
		MaxAge:       time.Now().Add(1 * time.Hour).Second(),
	}))
	e.Use(middleware.Gzip())

	// Add routes
	e.Add("GET", "health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hellow. Ok!")
	})

	// Add API routes.
	api := e.Group("/api")
	api.Add("POST", "/send-email", handlers.SendEmailHandler)

	// Log that server started
	s.co.Lo.Info("server has been started on port ", slog.Any("API", fmt.Sprintf("http://localhost:%d", s.co.Port)))
	return e.Start(fmt.Sprintf(":%d", s.co.Port))
}
