package dashboard

import (
	"net/http"

	"everlasting/src/domain/user"
	"everlasting/src/infrastructure/http/middleware"
	"everlasting/src/infrastructure/http/routes"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

// Get user id
//
//	@Security	BearerAuth
//	@Summary	Get user by id
//	@Param		id	path	string	true	"user id"
//	@Tags		User
//	@Produce	json
//	@Success	200	{object}	example.UserResponse
//	@Router		/user/{id} [get]
func getUserByIDHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container      = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		userRepository = container.Get("persistence.user").(*persistence.UserPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusInternalServerError, "ID should not be empty")
	}
	result, err := user.UserID(id).GetDetailFrom(ctx, userRepository)

	if err != nil {
		return err
	}
	return routes.JsonResponse(c, result, "Ok", "ok", 200, nil)
}

// Create user
//
//	@Security		BearerAuth
//	@Summary		Create new user account
//	@Description	### Param (JSON Body)
//	@Description	- name (string) : User Full Name
//	@Description	- email (string) : User email
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			input	body		user.UserInput	true	"User"
//	@Success		200		{object}	example.UserResponse
//	@Router			/user [post]
func createUserHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container      = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		userRepository = container.Get("persistence.user").(*persistence.UserPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	newUser := new(user.UserInput)
	if err := c.Bind(newUser); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err = c.Validate(newUser); err != nil {
		return err
	}

	result, err := newUser.SaveTo(ctx, userRepository)

	if err != nil {
		return err
	}
	return routes.JsonResponse(c, result, "Ok", "ok", 201, nil)
}

// Update user
//
//	@Security		BearerAuth
//	@Summary		Update user account
//	@Description	### Param (JSON Body)
//	@Description	- name (string) : User Full Name
//	@Description	- email (string) : User email
//	@Param			id	path	string	true	"user id"
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			input	body		user.UserInput	true	"User"
//	@Success		200		{object}	example.UserResponse
//	@Router			/user/{id} [put]
func updateUserHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container      = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		userRepository = container.Get("persistence.user").(*persistence.UserPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusInternalServerError, "ID should not be empty")
	}

	newUser := new(user.UserInput)
	if err := c.Bind(newUser); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err = c.Validate(newUser); err != nil {
		return err
	}

	result, err := newUser.UpdateTo(ctx, userRepository, user.UserID(id))

	if err != nil {
		return err
	}
	return routes.JsonResponse(c, result, "Ok", "ok", 201, nil)
}

// Get list of user
//
//	@Security	BearerAuth
//	@Summary	Get list of user. Paginated
//	@Param		q			query	string	false	"Search by name."
//	@Param		role		query	string	false	"Filter by user role"	Enums(superadmin, admin)
//	@Param		status		query	string	false	"Filter by user status"	Enums(active, inactive)
//	@Param		sort_by		query	string	false	"Sort by"				Enums(created_at, name, last_transaction)
//	@Param		order		query	string	false	"Sorting order"			Enums(asc, desc)
//	@Param		page		query	string	false	"Current page, default is 1"
//	@Param		per_page	query	string	false	"Displayed data per page, default is 20"
//	@Tags		User
//	@Produce	json
//	@Success	200	{object}	example.UsersResponse
//	@Router		/user [get]
func getUserListHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container      = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		userRepository = container.Get("persistence.user").(*persistence.UserPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	query := new(user.Query)
	if err := c.Bind(query); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err = c.Validate(query); err != nil {
		return err
	}

	result, err := query.CollectFrom(ctx, userRepository)

	if err != nil {
		return err
	}
	return routes.JsonResponse(c, result.Collection, "Ok", "ok", 200, map[string]interface{}{
		"pagination": result.Pagination,
	})
}

func RegisterUserRoutes(container di.Container, server *echo.Group) {
	user := server.Group("/user")
	user.Use(middleware.BearerAuthenticationMiddleware)

	user.GET("/:id", getUserByIDHandler)
	user.PUT("/:id", updateUserHandler)
	user.POST("", createUserHandler)
	user.GET("", getUserListHandler)
}
