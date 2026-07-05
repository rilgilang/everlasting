package middleware

import (
	"encoding/base64"
	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/sharedkernel/identity"
	"net/http"
	"slices"
	"strings"

	userDomain "everlasting/src/domain/user"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg/token"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

const MESSAGE_MISSING_HEADER = "Missing Authorization header"
const MESSAGE_INVALID_HEADER = "Invalid Authorization header"
const MESSAGE_INVALID_EVENT_ID = "Invalid Event ID"
const MESSAGE_SERVER_ERROR = "Internal Server Error"
const MESSAGE_NOT_FOUND = "Not Found"

func BearerAuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			container       = c.Get(string(MiddlewareValueContainer)).(di.Container)
			userRepository  = container.Get("persistence.user").(*persistence.UserPersistence)
			tokenRepository = container.Get("pkg.token").(*token.TokenProvider)
		)

		ctx := c.Request().Context()

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_MISSING_HEADER)
		}

		tokenizedHeader := strings.Split(authHeader, " ")
		if len(tokenizedHeader) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		if strings.ToLower(tokenizedHeader[0]) != "bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		tokenString := strings.TrimSpace(tokenizedHeader[1])
		user, err := userDomain.Token(tokenString).VerifyWith(ctx, tokenRepository, userRepository)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		c.Set("user_id", user.ID.String())
		c.Set("role", user.Role)

		// Continue to the next handler
		return next(c)
	}
}

func UserCheckEventMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			container           = c.Get(string(MiddlewareValueContainer)).(di.Container)
			userEventRepository = container.Get("persistence.user_event").(*persistence.UserEventPersistence)
			userId              = c.Get("user_id").(string)
			eventId             = c.Param("event_id")
		)

		ctx := c.Request().Context()

		userEvents, err := userEventRepository.GetOneByUserID(ctx, identity.FromStringOrNil(userId))
		if err != nil {
			if err == errDomain.ErrEventNotFound {
				return echo.NewHTTPError(http.StatusNotFound, MESSAGE_NOT_FOUND)
			}
			return echo.NewHTTPError(http.StatusInternalServerError, MESSAGE_SERVER_ERROR)
		}

		if !slices.Contains(userEvents.EventId, eventId) {
			return echo.NewHTTPError(http.StatusNotFound, MESSAGE_NOT_FOUND)
		}

		// Continue to the next handler
		return next(c)
	}
}

func BearerAdminAuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			container       = c.Get(string(MiddlewareValueContainer)).(di.Container)
			userRepository  = container.Get("persistence.user").(*persistence.UserPersistence)
			tokenRepository = container.Get("pkg.token").(*token.TokenProvider)
		)

		ctx := c.Request().Context()

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_MISSING_HEADER)
		}

		tokenizedHeader := strings.Split(authHeader, " ")
		if len(tokenizedHeader) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		if strings.ToLower(tokenizedHeader[0]) != "bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		tokenString := strings.TrimSpace(tokenizedHeader[1])
		user, err := userDomain.Token(tokenString).VerifyWith(ctx, tokenRepository, userRepository)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		if user.Role != userDomain.UserRoleAdmin {
			return echo.NewHTTPError(http.StatusForbidden, MESSAGE_INVALID_HEADER)
		}

		c.Set("user_id", user)

		// Continue to the next handler
		return next(c)
	}
}

func BearerUserAuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			container       = c.Get(string(MiddlewareValueContainer)).(di.Container)
			userRepository  = container.Get("persistence.user").(*persistence.UserPersistence)
			tokenRepository = container.Get("pkg.token").(*token.TokenProvider)
		)

		ctx := c.Request().Context()

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_MISSING_HEADER)
		}

		tokenizedHeader := strings.Split(authHeader, " ")
		if len(tokenizedHeader) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		if strings.ToLower(tokenizedHeader[0]) != "bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		tokenString := strings.TrimSpace(tokenizedHeader[1])
		user, err := userDomain.Token(tokenString).VerifyWith(ctx, tokenRepository, userRepository)

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		if user.Role != userDomain.UserRoleAdmin {
			return echo.NewHTTPError(http.StatusForbidden, MESSAGE_INVALID_HEADER)
		}

		c.Set("user_id", user)

		// Continue to the next handler
		return next(c)
	}
}

func BasicAuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			container = c.Get(string(MiddlewareValueContainer)).(di.Container)
		)

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_MISSING_HEADER)
		}

		tokenizedHeader := strings.Split(authHeader, " ")
		if len(tokenizedHeader) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		if strings.ToLower(tokenizedHeader[0]) != "basic" {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		tokenString := strings.TrimSpace(tokenizedHeader[1])
		authAccountByte, _ := base64.StdEncoding.DecodeString(tokenString)
		authAccountString := string(authAccountByte)
		savedAccount, _ := container.Get("config.auth").(string)
		if authAccountString != savedAccount {
			return echo.NewHTTPError(http.StatusUnauthorized, MESSAGE_INVALID_HEADER)
		}

		// Continue to the next handler
		return next(c)
	}
}
