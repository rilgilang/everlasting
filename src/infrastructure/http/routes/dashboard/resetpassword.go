package dashboard

import (
	"context"
	"encoding/json"

	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/user/resetpassword"

	"everlasting/src/domain/user"
	"everlasting/src/infrastructure/amqp"
	"everlasting/src/infrastructure/http/middleware"
	"everlasting/src/infrastructure/http/routes"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg/logger"
	"everlasting/src/infrastructure/pkg/token"

	"github.com/labstack/echo/v4"
	"github.com/sarulabs/di"
)

// Reset password request
//
//	@Summary	Request reset password email
//	@Description
//	@Tags		Reset Password
//	@Accept		json
//	@Produce	json
//	@Param		input	body		resetpassword.ResetPasswordRequest	true	"resetpassword.ResetPasswordRequest"
//	@Success	200		{object}	example.SuccessfulResetPassword
//	@Router		/reset-password/request [post]
func resetPasswordRequestHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container      = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		userRepository = container.Get("persistence.user").(*persistence.UserPersistence)
		messagebroker  = container.Get("pkg.messagebroker.amqp").(*amqp.MessageBroker)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	resetPasswordRequest := new(resetpassword.ResetPasswordRequest)
	err = json.NewDecoder(c.Request().Body).Decode(&resetPasswordRequest)
	if err != nil {
		return err
	}

	err = c.Validate(resetPasswordRequest)
	if err != nil {
		return err
	}

	hasMatchedUser, err := resetPasswordRequest.IsHasMatchedUserIn(ctx, userRepository)
	if err != nil {
		if !hasMatchedUser {
			return err
		}

		return nil
	}

	err = resetPasswordRequest.PutInstructionQueueIn(ctx, messagebroker)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, nil, "Ok", "ok", 200, nil)
}

// Reset password - redeem token
//
//	@Summary	Redeem token from verification code
//	@Description
//	@Tags		Reset Password
//	@Accept		json
//	@Produce	json
//	@Param		input	body		resetpassword.RedeemTokenRequest	true	"resetpassword.RedeemTokenRequest"
//	@Success	200		{object}	example.SuccessfulRedeemToken
//	@Router		/reset-password/redeem-token [post]
func resetPasswordRedeemTokenHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container       = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		userRepository  = container.Get("persistence.user").(*persistence.UserPersistence)
		tokenRepository = container.Get("pkg.token").(*token.TokenProvider)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	payload := new(resetpassword.RedeemTokenRequest)

	err = json.NewDecoder(c.Request().Body).Decode(&payload)
	if err != nil {
		return err
	}

	err = c.Validate(payload)
	if err != nil {
		return err
	}

	verificationCode := resetpassword.VerificationCode(payload.VerificationCode)

	var token resetpassword.Token
	user, err := verificationCode.VerifyWith(ctx, tokenRepository, userRepository)
	if err != nil {
		return err
	}

	token, err = resetpassword.GenerateTokenWith(ctx, user, tokenRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, token, "Ok", "ok", 200, nil)
}

// Reset password - submit
//
//	@Security	BearerAuth
//	@Summary	Submit new password
//	@Description
//	@Tags		Reset Password
//	@Accept		json
//	@Produce	json
//	@Param		input	body		resetpassword.SubmitPasswordRequest	true	"resetpassword.SubmitPasswordRequest"
//	@Success	200		{object}	example.SuccessfulRedeemToken
//	@Router		/reset-password/submit [put]
func resetPasswordSubmitHandler(c echo.Context) (err error) {
	defer c.Request().Body.Close()

	// Load container
	var (
		container      = c.Get(string(middleware.MiddlewareValueContainer)).(di.Container)
		userRepository = container.Get("persistence.user").(*persistence.UserPersistence)
	)

	cc := c.Get(string(middleware.MiddlewareValueAppLoggerContext)).(*logger.AppLoggerContext)
	ctx := cc.GetContext()

	newPassword := new(resetpassword.SubmitPasswordRequest)

	err = json.NewDecoder(c.Request().Body).Decode(newPassword)
	if err != nil {
		return err
	}

	err = c.Validate(newPassword)
	if err != nil {
		return err
	}

	uid, ok := c.Get(string(middleware.MiddlewareValueUserID)).(user.UserID)
	if !ok {
		return errDomain.ErrInvalidResetPasswordAToken
	}

	ctx = context.WithValue(ctx, resetpassword.UIDContextKey, uid)
	err = newPassword.SaveTo(ctx, userRepository)
	if err != nil {
		return err
	}

	return routes.JsonResponse(c, nil, "Ok", "ok", 200, nil)
}

func RegisterResetPasswordRoutes(container di.Container, server *echo.Group) {
	auth := server.Group("/reset-password")

	auth.POST("/request", resetPasswordRequestHandler)
	auth.POST("/redeem-token", resetPasswordRedeemTokenHandler)
	auth.PUT("/submit", resetPasswordSubmitHandler, middleware.ResetPasswordAuthenticationMiddleware)
}
