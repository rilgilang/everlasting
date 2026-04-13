package event

import (
	"context"
	"encoding/json"

	"everlasting/src/domain/user/resetpassword"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg"
	"everlasting/src/infrastructure/pkg/logger"
	"everlasting/src/infrastructure/pkg/smtp"
	"everlasting/src/infrastructure/pkg/token"

	"github.com/sarulabs/di"
)

type ResetPassword struct {
	container di.Container
	config    *pkg.Config
}

func (u *ResetPassword) Handle(ctx context.Context, body []byte) (err error) {
	var (
		userRepository  = u.container.Get("persistence.user").(*persistence.UserPersistence)
		tokenRepository = u.container.Get("pkg.token").(*token.TokenProvider)
		renderer        = u.container.Get("pkg.renderer.html").(pkg.HtmlRenderer)
		smtp            = u.container.Get("pkg.smtp.client").(*smtp.SmtpClients)
		logger          = u.container.Get("logger.app").(*logger.AppLogger)
	)

	resetPasswordRequest := new(resetpassword.ResetPasswordRequest)
	err = json.Unmarshal(body, resetPasswordRequest)
	if err != nil {
		logger.Error(ctx, "reset_password", err.Error())
		return err
	}

	// Getting matched user
	matchedUser, err := resetPasswordRequest.Email.GetMatchedUserIn(ctx, userRepository)
	if err != nil {
		return nil
	}

	// Generate verification code
	verificationCode, err := resetpassword.GenerateVerificationCodeWith(ctx, matchedUser, tokenRepository)
	if err != nil {
		return err
	}

	email, err := resetpassword.NewEmailMessageFrom(resetPasswordRequest.Email, verificationCode).
		RenderMessageWith(renderer, u.config.SMTPEmailFrom, u.config.AppClientBaseURL)

	if err != nil {
		return err
	}

	return email.SendWith(ctx, smtp)

}

func NewResetPassword(container di.Container, config *pkg.Config) *ResetPassword {
	return &ResetPassword{container: container, config: config}
}
