package user

import (
	"context"

	"everlasting/src/domain/sharedkernel/identity"

	"github.com/golang-jwt/jwt"
)

type TokenSubject string

const (
	TokenSubjectAccessToken   TokenSubject = "access_token"
	TokenSubjectRefreshToken  TokenSubject = "refresh_token"
	TokenSubjectResetPassword TokenSubject = "reset_password"
)

type (
	TokenOption struct {
		ValidAt   int64
		ExpiredAt int64
	}

	Token string

	// Generated token structure
	TokenSet struct {
		AccessToken  Token `json:"access_token"`
		RefreshToken Token `json:"refresh_token"`
	}

	Claims struct {
		ID      string `json:"id"`
		Subject string `json:"subject"`
		jwt.StandardClaims
	}
)

func (t Token) VerifyWith(ctx context.Context, tokenRepo TokenRepository, userRepo UserRepository) (user *User, err error) {
	id, err := tokenRepo.Verify(ctx, t)
	if err != nil {
		return user, err
	}

	return UserID(id.String()).GetDetailFrom(ctx, userRepo)
}

type TokenRepository interface {
	Generate(ctx context.Context, subject TokenSubject, identifier identity.ID, option TokenOption) (result Token, err error)
	Verify(ctx context.Context, token Token) (result identity.ID, err error)
	Revoke(ctx context.Context, token Token) (err error)
}
