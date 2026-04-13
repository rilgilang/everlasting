package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	b64 "encoding/base64"

	errDomain "everlasting/src/domain/error"
	cacheDomain "everlasting/src/domain/sharedkernel/cache"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/user"
	"everlasting/src/infrastructure/pkg/cache"
	"everlasting/src/infrastructure/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

type TokenClaims struct {
	Identifier string `json:"id"`
	Subject    string `json:"sub"`
	ValidTime  int64  `json:"nbf"`
}

func (c TokenClaims) IsValid(now time.Time) bool {
	return now.Unix() >= int64(c.ValidTime)
}

type TokenProvider struct {
	cache  *cache.RedisCache
	key    string
	logger *logger.AppLogger
}

func NewTokenProvider(cache *cache.RedisCache, key string, logger *logger.AppLogger) (result *TokenProvider, err error) {
	if len(key) != 32 {
		return result, errors.New("key should be 32 bytes sized string")
	}

	return &TokenProvider{
		cache,
		key,
		logger,
	}, nil
}

func (t *TokenProvider) Generate(ctx context.Context, subject user.TokenSubject, identifier identity.ID, option user.TokenOption) (result user.Token, err error) {
	now := time.Now()
	key := fmt.Sprintf("%s-%s-%d", subject, identifier.String(), now.UnixMilli())

	// Generate token
	bytes, err := bcrypt.GenerateFromPassword([]byte(key), 14)
	if err != nil {
		t.logger.Error(ctx, "bcrypt:generate_token", err.Error())
		return result, err
	}

	token := b64.StdEncoding.EncodeToString(bytes)

	// Compose claims
	claims := TokenClaims{
		Identifier: identifier.String(),
		Subject:    string(subject),
		ValidTime:  option.ValidAt,
	}

	// Define TTL
	exp := time.Unix(int64(option.ExpiredAt), 0)

	// Safe it to storage
	err = t.cache.Set(ctx, token, claims, exp.Sub(now))
	if err != nil {
		t.logger.Error(ctx, "cache:safe_token", err.Error())
		return result, err
	}

	return user.Token(token), err
}

func (t *TokenProvider) Verify(ctx context.Context, token user.Token) (result identity.ID, err error) {
	now := time.Now()
	claims := TokenClaims{}
	err = t.cache.Get(ctx, string(token), &claims)
	if err != nil {
		if err == cacheDomain.ErrCacheKeyNotFound {
			return result, errDomain.ErrInvalidAuth
		}

		t.logger.Error(ctx, "cache:get_token", err.Error())
		return result, err
	}

	if !claims.IsValid(now) {
		return result, errDomain.ErrInvalidAuth
	}

	result = identity.FromStringOrNil(claims.Identifier)
	if result.IsNil() {
		return result, errDomain.ErrInvalidAuth
	}

	return result, err
}

func (t *TokenProvider) Revoke(ctx context.Context, token user.Token) (err error) {
	return t.cache.Delete(ctx, string(token))
}
