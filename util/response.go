package util

import (
	"github.com/labstack/echo"
)

// Response contains custom status code, message and data
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseError renders JSON response for error
func ResponseError(c echo.Context, httpStatus int, code string, message string, err error) error {
	c.JSON(httpStatus, Response{Code: code, Message: message})
	return err
}
