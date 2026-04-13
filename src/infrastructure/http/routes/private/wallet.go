package private

import (
	"net/http"

	"everlasting/src/domain/wallet"
	"everlasting/src/infrastructure/http/middleware"
	"everlasting/src/infrastructure/http/routes"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg/logger"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

// Get wallet by id
//
//	@Summary	Get wallet by id
//	@Security	BasicAuth
//	@Param		id	path	string	true	"wallet id"
//	@Tags		Wallet
//	@Produce	json
//	@Success	200	{object}	example.WalletResponse
//	@Router		/wallet/{id} [get]
func getWalletByIDHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container        = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		walletRepository = container.Get("persistence.wallet").(*persistence.WalletPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	id := c.Param("id")
	if id == "" {
		return c.String(http.StatusInternalServerError, "ID should not be empty")
	}

	result, err := wallet.WalletID(id).GetDetailFrom(ctx, walletRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, result, "Ok", "ok", 200, nil)
}

// Get wallet by account
//
//	@Summary	Get wallet by account
//	@Security	BasicAuth
//	@Param		type	path	string	true	"account type"	Enums(doctor, user, laboratory)
//	@Param		id		path	string	true	"account id"
//	@Tags		Wallet
//	@Produce	json
//	@Success	200	{object}	example.WalletResponse
//	@Router		/wallet/account/{type}/{id} [get]
func getWalletByAccount(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	account := wallet.Account{
		Type: wallet.AccountType(c.Param("type")),
		ID:   c.Param("id"),
	}

	if err = c.Validate(account); err != nil {
		return err
	}

	// Load container
	var (
		container        = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		walletRepository = container.Get("persistence.wallet").(*persistence.WalletPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	result, err := account.GetDetailFrom(ctx, walletRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, result, "Ok", "ok", 200, nil)
}

// Create wallet
//
//	@Summary		Create new wallet
//	@Security		BasicAuth
//	@Description	### Param (JSON Body)
//	@Tags			Wallet
//	@Accept			json
//	@Produce		json
//	@Param			input	body		wallet.CreateWalletRequest	true	"Wallet"
//	@Success		200		{object}	example.WalletResponse
//	@Router			/wallet [post]
func createWalletHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container        = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		walletRepository = container.Get("persistence.wallet").(*persistence.WalletPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	// Generate request payload
	createWalletRequest := new(wallet.CreateWalletRequest)
	if err := c.Bind(createWalletRequest); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err = c.Validate(createWalletRequest); err != nil {
		return err
	}

	// Process saving
	result, err := createWalletRequest.GetOrSaveTo(ctx, walletRepository)
	if err != nil {
		return err
	}
	return routes.JsonResponse(c, result, "Ok", "ok", 201, nil)
}

// Get list of wallet
//
//	@Summary	Get list of wallet. Paginated
//	@Security	BasicAuth
//	@Param		q			query	string	false	"Search by wallet id, holder name, ai care account id."
//	@Param		ref_type	query	string	false	"Filter by reference type"	Enums(doctor, user, laboratory)
//	@Param		status		query	string	false	"Filter by wallet status"	Enums(active, inactive)
//	@Param		sort_by		query	string	false	"Sort by"					Enums(created_at, name, last_transaction, balance)
//	@Param		order		query	string	false	"Sorting order"				Enums(asc, desc)
//	@Param		page		query	string	false	"Current page, default is 1"
//	@Param		per_page	query	string	false	"Displayed data per page, default is 20"
//	@Tags		Wallet
//	@Produce	json
//	@Success	200	{object}	example.WalletsResponse
//	@Router		/wallet [get]
func getWalletListHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container        = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		walletRepository = container.Get("persistence.wallet").(*persistence.WalletPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	query := new(wallet.Query)
	if err := c.Bind(query); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err = c.Validate(query); err != nil {
		return err
	}

	result, err := query.CollectFrom(ctx, walletRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, result.Collection, "Ok", "ok", 200, map[string]interface{}{
		"pagination": result.Pagination,
	})
}

func RegisterWalletRoutes(container di.Container, server *echo.Group) {
	wallet := server.Group("/wallet")
	wallet.Use(middleware.BasicAuthenticationMiddleware)

	wallet.POST("", createWalletHandler)
	wallet.GET("", getWalletListHandler)
	wallet.GET("/:id", getWalletByIDHandler)
	wallet.GET("/account/:type/:id", getWalletByAccount)
}
