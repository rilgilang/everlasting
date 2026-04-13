package middleware

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	errDomain "everlasting/src/domain/error"
	"everlasting/src/infrastructure/http/routes"

	validatorLib "github.com/go-playground/validator"
)

type ErrorAttributes struct {
	Status  int
	Code    string
	Message string
}

func NewErrorAttributes(status int, code, message string) *ErrorAttributes {
	return &ErrorAttributes{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

var ErrorMap map[error]ErrorAttributes = map[error]ErrorAttributes{
	errDomain.ErrValidation:                       *NewErrorAttributes(400, "invalid_input", "Please check your input"),
	errDomain.ErrWalletNotFound:                   *NewErrorAttributes(404, "wallet_not_found", "We can't find your wallet"),
	errDomain.ErrWalletAlreadyExists:              *NewErrorAttributes(409, "wallet_exists", "Wallet already exists"),
	errDomain.ErrInsufficientBalance:              *NewErrorAttributes(409, "insufficient balance", "Insufficient balance"),
	errDomain.ErrInvalidCredential:                *NewErrorAttributes(400, "invalid_credential", "Invalid Credential"),
	errDomain.ErrUserNotFound:                     *NewErrorAttributes(404, "user_not_found", "User not found"),
	errDomain.ErrInvalidAuth:                      *NewErrorAttributes(400, "invalid_token", "Invalid authentication token"),
	errDomain.ErrTransactionReferenceNotAvailable: *NewErrorAttributes(409, "transaction_already_exist", "Transaction reference not available"),
}

func HttpErrorHandler() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err != nil {
				return GenerateHTTPErrorResponse(c, err)

			}
			return err
		}
	}
}

func GenerateHTTPErrorResponse(c echo.Context, err error) error {
	// Handle if error is instance of echo HttpError
	if he, ok := err.(*echo.HTTPError); ok {
		var message string
		if message, ok = he.Message.(string); !ok {
			message = "An error occurred"
		}
		return c.JSON(he.Code, routes.NewResponse(nil, message, strconv.Itoa(he.Code), he.Code, nil))
	}

	// Handle if error type is validation error
	if he, ok := err.(validatorLib.ValidationErrors); ok {
		payload := routes.NewErrorResponseFromValidator(he)
		return c.JSON(payload.Meta.Status, payload)
	}

	// Handle if error is defined in map
	if found, ok := ErrorMap[err]; ok {
		return c.JSON(int(found.Status), routes.NewResponse(nil, found.Message, found.Code, int(found.Status), nil))
	}
	// Otherwise, return 500 error
	return c.JSON(
		http.StatusInternalServerError,
		routes.NewResponse(
			nil,
			"Internal server error",
			"server_error",
			http.StatusInternalServerError,
			nil,
		),
	)

}
