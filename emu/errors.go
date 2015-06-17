package emu

import (
	"net/http"

	"github.com/labstack/echo"
)

var JSONErrorHandler = func(err error, c *echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code()
		msg = he.Error()
	}
	c.JSON(code, map[string]string{"error": msg})
}

func Errors(h echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		err := h(c)
		if err != nil {
			c.Error(err)
		}
		return nil
	}
}
