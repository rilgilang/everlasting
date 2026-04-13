package middleware

import (
	"net/http"
	"strings"

	"everlasting/src/domain/user/resetpassword"
	"everlasting/src/infrastructure/pkg/token"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

func ResetPasswordAuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			container       = c.Get(string(MiddlewareValueContainer)).(di.Container)
			tokenRepository = container.Get("pkg.token").(*token.TokenProvider)
		)

		ctx := c.Request().Context()

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}

		tokenizedHeader := strings.Split(authHeader, " ")
		if len(tokenizedHeader) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header")
		}

		if strings.ToLower(tokenizedHeader[0]) != "bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header")
		}

		tokenString := strings.TrimSpace(tokenizedHeader[1])
		token := resetpassword.Token(tokenString)
		uid, err := token.VerifyWith(ctx, tokenRepository)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		c.Set(string(MiddlewareValueUserID), uid)

		// Continue to the next handler
		return next(c)
	}
}
