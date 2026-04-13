package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

func InjectContainer(container di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(string(MiddlewareValueContainer), container)
			return next(c)
		}
	}
}
