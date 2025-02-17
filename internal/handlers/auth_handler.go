package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	if len(u.Email) == 0 || len(u.Password) == 0 {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid email or password."))
	}

	// Save the user.
	db := hctx.GetCore().DB
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	collection := db.Database("referrer").Collection("users")
	_, err := collection.InsertOne(ctx, u)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, "user created")
}

type LoginUserDto struct {
	Email  string `json:"email"`
	Passwd string `json:"password"`
}

func LoginHandler(c echo.Context) error {
	var user LoginUserDto
	err := c.Bind(&user)

	if err != nil {
		return SendErrorResponse(c, http.StatusBadRequest, err)
	}

	// Check the user credentails and other into Dbs.

	// return Success or Failure
	token := "token-12424356789"

	resp := map[string]any{
		"token":  token,
		"sucess": true,
	}

	return c.JSON(http.StatusOK, resp)
}
