package resetpassword

import (
	"context"

	"everlasting/src/domain/user"
)

const (
	VerificationCodePrefix       string            = "RESET_PASSWORD"
	VerificationCodeTTLInSeconds int64             = 2 * 60 * 60
	TokenSubjectVerificationCode user.TokenSubject = "verification_code"
)

type (
	VerificationCode   user.Token
	RedeemTokenRequest struct {
		VerificationCode string `json:"verification_code" example:"my_verification_code" validate:"required"`
	}
)

// Verification code can initated by email
// It needs email, generator object
func GenerateVerificationCodeWith(ctx context.Context, matchedUser *user.User, tokenProvider user.TokenRepository) (result VerificationCode, err error) {
	token, err := tokenProvider.Generate(ctx, TokenSubjectVerificationCode, matchedUser.ID, user.TokenOption{})
	if err != nil {
		return result, err
	}
	return VerificationCode(token), err
}

func (c VerificationCode) VerifyWith(ctx context.Context, tokenProvider user.TokenRepository, repo user.UserRepository) (result *user.User, err error) {
	uid, err := tokenProvider.Verify(ctx, user.Token(c))
	if err != nil {
		return result, err
	}
	err = tokenProvider.Revoke(ctx, user.Token(c))
	if err != nil {
		return result, err
	}
	return user.UserID(uid.String()).GetDetailFrom(ctx, repo)
}
