package user

import (
	"context"

	errDomain "everlasting/src/domain/error"
)

type Credential struct {
	Email    Email    `json:"email" example:"admin@klondike.id" validate:"required,email"`
	Password Password `json:"password" example:"1234qweR!" validate:"required,password_custom_validator"`
}

func (c *Credential) VerifyWith(ctx context.Context, repo UserRepository) (result *User, err error) {

	// Find a matched account with credential email
	result, err = c.Email.GetMatchedUserIn(ctx, repo)
	if err != nil {
		// If there is no matched email in data storage. Expected persistence should return errDomain.NotFoundEntityError
		if err == errDomain.ErrUserNotFound {
			return result, errDomain.ErrInvalidCredential
		}
		// Otherwise ...
		return result, err
	}

	// If matched account doesn't have valid password
	if result.VerifyPassword(c.Password) != nil {
		return result, errDomain.ErrInvalidCredential
	}

	return result, err
}
