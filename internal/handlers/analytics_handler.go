package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func ProfileAnalyticsHandler(c echo.Context) error {
	// Get Context.
	hctx := c.(*HandlerContext)

	userEmail := c.QueryParam("email")
	if len(userEmail) == 0 || !isValidEmail(userEmail) {
		return SendErrorResponse(c, http.StatusBadRequest, fmt.Errorf("invalid email or no email found"))
	}

	analytics, err := hctx.GetCore().DB.GetProfileAnalytics(userEmail)
	if err != nil {
		return SendErrorResponse(c, http.StatusInternalServerError, fmt.Errorf("something went wrong: %w", err))
	}

	return c.JSON(http.StatusOK, analytics)
}
