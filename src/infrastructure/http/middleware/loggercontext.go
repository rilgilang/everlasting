package middleware

import (
	"bytes"
	"io"

	"everlasting/src/infrastructure/pkg/logger"

	"github.com/labstack/echo/v4"
)

// RequestBodyLogger Set request_body
func RequestBodyLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Read request body
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		bodyString := string(bodyBytes)

		// Restore the request body so it can be read again later
		c.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Store request body in context to be logged
		c.Set("request_body", bodyString)

		return next(c)
	}
}

func InjectLoggerContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(string(MiddlewareValueAppLoggerContext), logger.NewAppLoggerContextFromEchoContext(c))
			return next(c)
		}
	}
}
