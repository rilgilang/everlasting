package http

import (
	_ "everlasting/docs"
	"everlasting/src/domain/validator"
	"everlasting/src/infrastructure/amqp"
	md "everlasting/src/infrastructure/http/middleware"
	"everlasting/src/infrastructure/http/routes/dashboard"
	"everlasting/src/infrastructure/http/routes/guest"
	"everlasting/src/infrastructure/pkg"
	"everlasting/src/infrastructure/pkg/logger"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RunDashboardAPI(container di.Container, config *pkg.Config) {
	server := echo.New()
	server.IPExtractor = echo.ExtractIPFromRealIPHeader()
	server.Use(middleware.Recover())

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{config.GetCORSAllowedDomain()},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Inject container to echo context
	server.Use(md.InjectContainer(container))

	// Add http error handler middleware
	server.Use(md.HttpErrorHandler())

	// Add context middleware
	server.Use(md.InjectLoggerContext())

	// Add tdr log middleware
	server.Use(middleware.RequestLoggerWithConfig(logger.GenerateTdrLogConfig(
		container.Get("logrus.tdr").(*logrus.Logger),
	)))

	// Define swagger apidocs page
	server.GET("/apidocs/*", echoSwagger.WrapHandler)

	// Define custom validator
	server.Validator = validator.NewCustomValidator()

	server.HTTPErrorHandler = customHttpErrorHandler

	// Register route groups
	api := server.Group("/api/v1")

	dashboard.RegisterEventRoutes(container, api)
	guest.RegisterGuestRoutes(container, api)
	//dashboard.RegisterUserRoutes(container, api)
	//dashboard.RegisterWalletRoutes(container, api)
	//dashboard.RegisterTransactionRoutes(container, api)
	//dashboard.RegisterAuthRoutes(container, api)
	//dashboard.RegisterResetPasswordRoutes(container, api)

	// Register message broker handler
	go amqp.Consume(container, config)

	// Start server
	server.Logger.Fatal(server.Start(fmt.Sprintf(":%d", config.AppPort)))
}

func customHttpErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	md.GenerateHTTPErrorResponse(c, err)
}
