package error

import (
	"errors"
)

// User based error defintion
var (
	// Authentication
	ErrInvalidCredential error = errors.New("invalid credential")
	ErrUserLockout       error = errors.New("user locked")

	ErrValidation error = errors.New("validation error")

	ErrEventNotFound error = errors.New("event not found")

	ErrWalletNotFound      error = errors.New("wallet not found")
	ErrWalletAlreadyExists error = errors.New("wallet already exists")

	ErrTransactionNotFound              error = errors.New("transaction not found")
	ErrTransactionReferenceNotAvailable error = errors.New("transaction reference not available")

	ErrInsufficientBalance error = errors.New("insufficient balance")

	ErrInvalidAuth error = errors.New("invalid authentication token")

	ErrNotFoundEntity error = errors.New("entity not found")

	ErrUserNotFound          error = errors.New("user not found")
	ErrUserAlreadyRegistered error = errors.New("user already registered")
	ErrInvalidInput          error = errors.New("invalid input")
	ErrBadRequest            error = errors.New("bad request")
	ErrorInvalidCredential   error = errors.New("invalid credential")
	ErrDataNotFound          error = errors.New("data not found")

	ErrResetPasswordGenerate      error = errors.New("error generate reset password")
	ErrInvalidVerificationCode    error = errors.New("Invalid verification code")
	ErrInvalidResetPasswordAToken error = errors.New("Invalid reset password token")
	ErrResetPasswordSendingEmail  error = errors.New("error while process reset password request")
)
