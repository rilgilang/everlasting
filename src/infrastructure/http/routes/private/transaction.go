package private

import (
	"context"
	"net/http"

	"everlasting/src/domain/sharedkernel/unitofwork"
	"everlasting/src/domain/transaction"
	"everlasting/src/domain/wallet"
	"everlasting/src/infrastructure/http/middleware"
	"everlasting/src/infrastructure/http/routes"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg/logger"

	"github.com/go-redsync/redsync/v4"
	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

// Create transaction
//
//	@Summary		Create new transaction account
//	@Security		BasicAuth
//	@Description	### Param (JSON Body)
//	@Description	- amount (number) : Transaction amount
//	@Description	- notes (string) : Transaction amount
//	@Description	- ref_type (string) : Transaction reference type. "order" | "payout" | "refund"
//	@Description	- ref_id (string, optional) : Transaction reference id
//	@Description	- type (string) : Transaction type. "debt" | "credit"
//	@Description	- wallet_id (string, uuid) : Target wallet id
//	@Description	- wallet {object} : "Target create wallet. Mandatory if wallet id is empty. Otherwise, this field is omited"
//
// @Tags			Transaction
// @Accept			json
// @Produce		json
// @Param			input	body		transaction.CreateTransactionRequest	true	"Transaction"
// @Success		200		{object}	example.TransactionResponse
// @Router			/transaction [post]
func createTransactionHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container             = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		transactionRepository = container.Get("persistence.transaction").(*persistence.TransactionPersistence)
		walletRepository      = container.Get("persistence.wallet").(*persistence.WalletPersistence)
		uow                   = container.Get("persistence.uow").(*persistence.UnitOfWork)
		lock                  = container.Get("provider.lock").(*redsync.Redsync)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	// Get, validate and cast payload
	transactionRequest := new(transaction.CreateTransactionRequest)
	if err := c.Bind(transactionRequest); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err = c.Validate(transactionRequest); err != nil {
		return err
	}

	payload, wallet, err := transactionRequest.VerifyAndTransform(ctx, walletRepository, transactionRepository)
	if err != nil {
		return err
	}

	// Operation begin
	var result *transaction.Transaction

	// Initiate distributed lock
	mutex := lock.NewMutex(payload.WalletID.String())
	if err := mutex.Lock(); err != nil {
		return err
	}

	defer func() {
		mutex.Unlock()
	}()

	// Run unit of work (transaction)
	uowResult, err := uow.Execute(ctx, func(ctx context.Context) (result *unitofwork.Result, err error) {
		transaction, err := payload.SaveTo(ctx, transactionRepository, walletRepository)
		if err != nil {
			return result, err
		}

		err = transaction.UpdateBalanceToWallet(ctx, wallet, transactionRepository, walletRepository)
		if err != nil {
			return result, err
		}

		return &unitofwork.Result{
			Body: transaction,
		}, err
	})

	if uowResult != nil && uowResult.Body != nil {
		if t, ok := uowResult.Body.(*transaction.Transaction); ok {
			result = t
		}
	}

	if err != nil {
		return err
	}
	return routes.JsonResponse(c, result, "Ok", "ok", 201, nil)
}

// Get list of transaction of a walletID
//
//	@Summary	Get list of transaction. Paginated
//	@Security	BasicAuth
//	@Param		wallet_id	path	string	true	"wallet id"
//	@Param		type		query	string	false	"Filter by transaction type"	Enums(debt, credit)
//	@Param		ref_type	query	string	false	"Reference type"				Enums(order, payout, refund)
//	@Param		cursor		query	number	false	"Cursor value obtained from previous page. Left it empty or 0 for the first page"
//	@Param		date_from	query	string	false	"Time window filter:from in YYYY-mm-dd format"
//	@Param		date_until	query	string	false	"Time window filter:until in YYYY-mm-dd format"
//	@Param		per_page	query	string	false	"Displayed data per page, default is 20"
//	@Tags		Transaction
//	@Produce	json
//	@Success	200	{object}	example.WalletsResponse
//	@Router		/transaction/{wallet_id} [get]
func getWalletTransactionListHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container             = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		transactionRepository = container.Get("persistence.transaction").(*persistence.TransactionPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	walletID := c.Param("wallet_id")
	if walletID == "" {
		return c.String(http.StatusInternalServerError, "Wallet ID should not be empty")
	}
	query := new(transaction.Query)
	if err := c.Bind(query); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	if err = c.Validate(query); err != nil {
		return err
	}

	result, err := query.CollectFrom(ctx, wallet.WalletID(walletID), transactionRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, result.Collection, "Ok", "ok", 200, map[string]interface{}{
		"pagination": result.Pagination,
	})
}

func RegisterTransactionRoutes(container di.Container, server *echo.Group) {
	transaction := server.Group("/transaction")
	transaction.Use(middleware.BasicAuthenticationMiddleware)

	transaction.POST("", createTransactionHandler)
	transaction.GET("/:wallet_id", getWalletTransactionListHandler)
}
