package dashboard

import (
	"everlasting/src/domain/event"
	"everlasting/src/infrastructure/http/middleware"
	"everlasting/src/infrastructure/http/routes"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
	"net/http"
)

func createEventHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container       = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		eventRepository = container.Get("persistence.event").(*persistence.EventPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	newEvent := new(event.EventInput)
	if err := c.Bind(newEvent); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err = c.Validate(newEvent); err != nil {
		return err
	}

	result, err := newEvent.SaveEvent(ctx, eventRepository)

	if err != nil {
		return err
	}
	return routes.JsonResponse(c, result, "Ok", "ok", 201, nil)
}

func updateEventHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container       = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		eventRepository = container.Get("persistence.user").(*persistence.EventPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusInternalServerError, "ID should not be empty")
	}

	updateEvent := new(event.EventInput)
	if err := c.Bind(updateEvent); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err = c.Validate(updateEvent); err != nil {
		return err
	}

	result, err := updateEvent.UpdateTo(ctx, eventRepository, event.EventID(id))

	if err != nil {
		return err
	}
	return routes.JsonResponse(c, result, "Ok", "ok", 201, nil)
}

func getEventList(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container       = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		eventRepository = container.Get("persistence.event").(*persistence.EventPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	query := new(event.Query)
	if err := c.Bind(query); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err = c.Validate(query); err != nil {
		return err
	}
	result, err := query.CollectFrom(ctx, eventRepository)

	if err != nil {
		return err
	}
	return routes.JsonResponse(c, result.Collection, "Ok", "ok", 200, map[string]interface{}{
		"pagination": result.Pagination,
	})
}

func RegisterEventRoutes(container di.Container, server *echo.Group) {
	event := server.Group("/event")
	event.Use(middleware.BearerAdminAuthenticationMiddleware)

	//event.GET("/:id", getUserByIDHandler)
	event.PUT("/:id", updateEventHandler)
	event.POST("", createEventHandler)
	event.GET("", getEventList)
}
