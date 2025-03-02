package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sounishnath003/customgo-mailer-service/internal/core"
	"github.com/sounishnath003/customgo-mailer-service/internal/handlers"
	"github.com/sounishnath003/customgo-mailer-service/internal/utils"
	"golang.org/x/time/rate"
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
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{echo.POST, echo.GET},
		AllowCredentials: true,
		AllowHeaders:     []string{"X-API-TrackerId", echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization, echo.HeaderContentLength},
		MaxAge:           time.Now().Add(1 * time.Hour).Second(),
	}))
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(utils.GetStringFromEnv("JWT_SECRET_KEY", "Sec%&!*RT#*!@(89231%&!*RT#12345")),
		Skipper: func(c echo.Context) bool {
			if c.Path() == "/api/auth/login" || c.Path() == "/api/auth/signup" {
				return true
			}
			return true
		},
	}))
	e.Use(middleware.Gzip())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(5))))

	// Serve static files from web/dist directory as per in docker container
	e.Static("/", "web/dist")

	// Add routes
	e.Add("GET", "health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hellow. Ok!")
	})

	// Add API routes.
	api := e.Group("/api")

	// Auth endpoints.
	api.Add("POST", "/auth/signup", handlers.SignupHandler)
	api.Add("POST", "/auth/login", handlers.LoginHandler)
	// Profile endpoints.
	api.Add("GET", "/profile", handlers.GetProfileHandler)
	api.Add("POST", "/profile/information", handlers.ProfileInformationHandler)
	// Draft Coldmails Ai endpoints.
	api.Add("POST", "/draft-with-ai", handlers.DraftReferralEmailWithAiHandler)
	// Email endpoints.
	api.Add("GET", "/sent-referrals", handlers.GetReferralEmailsHandler)
	api.Add("POST", "/send-email", handlers.SendEmailHandler)

	// Log that server started
	s.co.Lo.Info("server has been started on port ", slog.Any("API", fmt.Sprintf("http://localhost:%d", s.co.Port)))
	return e.Start(fmt.Sprintf(":%d", s.co.Port))
}
