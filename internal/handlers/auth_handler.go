package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sounishnath003/customgo-mailer-service/internal/repository"
)

// SignupHandler handles the new user signup feature.
// User come to setup the account holder for the application to consume.
func SignupHandler(c echo.Context) error {
	// Get the context
	hctx := c.(*HandlerContext)

	// Bind
	u := &repository.User{}

	if err := c.Bind(u); err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	// Check validation.
	if len(u.Email) == 0 || len(u.Password) == 0 || !isValidEmail(u.Email) {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid email or password."))
	}

	// Save the user.
	err := hctx.GetCore().DB.CreateUser(u)
	// It may failed due to duplicate email as email has to be unique
	if err != nil {
		hctx.GetCore().Lo.Error("duplicate key found", "error", err.Error())
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, "user created")
}

type LoginUserDto struct {
	Email  string `json:"email"`
	Passwd string `json:"password"`
}

func LoginHandler(c echo.Context) error {
	// Get the core
	hctx := c.(*HandlerContext)

	var user LoginUserDto
	err := c.Bind(&user)

	if len(user.Email) == 0 || !isValidEmail(user.Email) {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid email address"))
	}

	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	// Check the user credentails and other into Dbs.
	u, err := hctx.GetCore().DB.FindUserByEmailAndPassword(user.Email, user.Passwd)

	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	resp := map[string]any{
		// "user":        u,
		"accessToken": u.Token,
		"success":     true,
	}

	return c.JSON(http.StatusOK, resp)
}
