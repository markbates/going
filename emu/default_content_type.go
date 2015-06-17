package emu

import "github.com/labstack/echo"

func DefaultContentType(s string) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if c.Request().Header.Get("Content-Type") == "" {
				c.Request().Header.Set("Content-Type", s)
			}
			return h(c)
		}
	}
}
