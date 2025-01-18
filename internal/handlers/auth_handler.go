package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

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
