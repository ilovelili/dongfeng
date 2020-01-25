package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// Healthcheck GET /helthcheck
func Healthcheck(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
