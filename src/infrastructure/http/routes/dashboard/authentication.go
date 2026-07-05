package dashboard

import (
	"encoding/json"
	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/sharedkernel/marshaler"
	"net/http"

	"everlasting/src/domain/user"
	"everlasting/src/infrastructure/http/middleware"
	"everlasting/src/infrastructure/http/routes"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg/logger"
	"everlasting/src/infrastructure/pkg/token"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

// Authentication Login Handler
//
//	@Summary		Handle user login
//	@Description	Authenticates a user using provided credentials.
//	@Tags			Authentication
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		user.Credential	true	"User credentials for login"
//	@Success		200			{object}	example.SuccessfulLoginResponse
//	@Failure		401			{object}	example.InvalidCredentialResponse
//	@Router			/auth/login [post]
func loginHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container       = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		userRepository  = container.Get("persistence.user").(*persistence.UserPersistence)
		tokenRepository = container.Get("pkg.token").(*token.TokenProvider)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	credential := new(user.Credential)
	err = json.NewDecoder(c.Request().Body).Decode(credential)
	if err != nil {
		return err
	}

	// Verify credential
	account, err := credential.VerifyWith(ctx, userRepository)
	if err != nil {
		return err
	}

	// If valid, generate token
	result, err := account.GenerateTokenWith(ctx, tokenRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, result, "Ok", "ok", 200, nil)
}

func generateHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	var (
		container      = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		userRepository = container.Get("persistence.user").(*persistence.UserPersistence)
		cc             = c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
		ctx            = cc.GetContext()
	)

	input := new(user.UserInput)
	if err := c.Bind(input); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	generated, err := user.NewCipherTextFromPassword(input.Password)
	if err != nil {
		return err
	}

	_, err = input.Email.GetMatchedUserIn(ctx, userRepository)
	if err != nil && err != errDomain.ErrUserNotFound {
		return err
	}

	user := user.User{
		ID:         identity.NewID(),
		Email:      input.Email,
		Name:       input.Name,
		Role:       user.UserRole(input.Role),
		Status:     "active",
		CreatedAt:  marshaler.JsonTime{},
		UpdatedAt:  marshaler.JsonTime{},
		CipherText: generated,
	}

	if _, err := user.SaveTo(ctx, userRepository); err != nil {
		return err
	}

	return routes.JsonResponse(c, input, "Ok", "ok", 201, nil)
}

func RegisterAuthRoutes(container di.Container, server *echo.Group) {
	auth := server.Group("/auth")
	auth.POST("/login", loginHandler)
	auth.POST("/generate", generateHandler)
}
