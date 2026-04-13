package resetpassword

import (
	"context"
	"time"

	userDomain "everlasting/src/domain/user"
)

type Token userDomain.Token

func GenerateTokenWith(ctx context.Context, usr *userDomain.User, tokenRepo userDomain.TokenRepository) (token Token, err error) {
	now := time.Now().UTC()
	tokenExp := now.Add(time.Duration(2*60*60) * time.Second) // Expired in 2 hours
	option := userDomain.TokenOption{
		ValidAt:   now.Unix(),
		ExpiredAt: tokenExp.Unix(),
	}

	userToken, err := tokenRepo.Generate(ctx, userDomain.TokenSubjectResetPassword, usr.ID, option)
	if err != nil {
		return token, err
	}
	return Token(userToken), err
}

func (t Token) VerifyWith(ctx context.Context, tokenRepo userDomain.TokenRepository) (uid userDomain.UserID, err error) {
	identifier, err := tokenRepo.Verify(ctx, userDomain.Token(t))
	if err != nil {
		return uid, err
	}

	err = tokenRepo.Revoke(ctx, userDomain.Token(t))
	if err != nil {
		return uid, err
	}

	return userDomain.UserID(identifier.String()), err
}
