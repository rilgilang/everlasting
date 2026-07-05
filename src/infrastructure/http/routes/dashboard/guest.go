package dashboard

import (
	"net/http"

	"everlasting/src/domain/guest"
	"everlasting/src/infrastructure/http/middleware"
	"everlasting/src/infrastructure/http/routes"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

func createGuestHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		container       = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		guestRepository = container.Get("persistence.guest").(*persistence.GuestPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	eventId := c.Param("event_id")
	input := new(guest.GuestInput)
	if err := c.Bind(input); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	input.EventId = eventId

	if err := c.Validate(input); err != nil {
		return err
	}

	result, err := input.SaveTo(ctx, guestRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, result, "Ok", "ok", 201, nil)
}

func getGuestByIDHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		container       = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		guestRepository = container.Get("persistence.guest").(*persistence.GuestPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusBadRequest, "ID should not be empty")
	}

	result, err := guest.GuestID(id).GetDetailFrom(ctx, guestRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, result, "Ok", "ok", 200, nil)
}

func updateGuestHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		container       = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		guestRepository = container.Get("persistence.guest").(*persistence.GuestPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusBadRequest, "ID should not be empty")
	}

	input := new(guest.GuestInput)
	if err := c.Bind(input); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err := c.Validate(input); err != nil {
		return err
	}

	result, err := input.UpdateTo(ctx, guestRepository, id)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, result, "Ok", "ok", 200, nil)
}

func deleteGuestHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		container       = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		guestRepository = container.Get("persistence.guest").(*persistence.GuestPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusBadRequest, "ID should not be empty")
	}

	err := guest.GuestID(id).DeleteFrom(ctx, guestRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, nil, "Ok", "ok", 200, nil)
}

func getGuestListHandler(c echo.Context) error {
	defer c.Request().Body.Close()

	var (
		container       = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		guestRepository = container.Get("persistence.guest").(*persistence.GuestPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	eventId := c.Param("event_id")
	query := new(guest.Query)
	if err := c.Bind(query); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	query.EventId = eventId

	if err := c.Validate(query); err != nil {
		return err
	}

	result, err := query.CollectFrom(ctx, guestRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, result.Collection, "Ok", "ok", 200, map[string]interface{}{
		"pagination": result.Pagination,
	})
}

func RegisterGuestRoutes(container di.Container, server *echo.Group) {
	g := server.Group("/:event_id/guest")
	g.Use(middleware.BearerAuthenticationMiddleware)
	g.Use(middleware.UserCheckEventMiddleware)

	g.POST("", createGuestHandler)
	g.GET("", getGuestListHandler)
	g.GET("/:id", getGuestByIDHandler)
	g.PUT("/:id", updateGuestHandler)
	g.DELETE("/:id", deleteGuestHandler)
}
