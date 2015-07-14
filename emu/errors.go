package emu

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/markbates/going/imt"
)

type ErrorHandler func(error, *echo.Context)

var Errors = map[error]ErrorHandler{}

func init() {
	Errors[sql.ErrNoRows] = func(err error, c *echo.Context) {
		if c.Request().Header.Get("Content-Type") == imt.Application.JSON {
			JSONErrorHandler(echo.NewHTTPError(404, "Record not found"), c)
		} else {
			c.Render(404, "404", err)
		}
	}
}

var JSONErrorHandler = func(err error, c *echo.Context) {
	code := http.StatusInternalServerError
	msg := http.StatusText(code)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code()
		msg = he.Error()
	}
	c.JSON(code, map[string]string{"error": msg})
}

func HandleErrors(err error, c *echo.Context) {
	f := Errors[err]
	if f != nil {
		f(err, c)
		return
	}
	if c.Request().Header.Get("Content-Type") == imt.Application.JSON {
		JSONErrorHandler(echo.NewHTTPError(500, err.Error()), c)
	} else {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code()
		}
		c.Render(code, strconv.Itoa(code), err)
	}
}
