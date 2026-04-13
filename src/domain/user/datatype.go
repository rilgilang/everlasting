package user

import (
	"context"
)

// Following is non behavioral data type definition
type (
	Email    string
	Password string
)

func (e Email) GetMatchedUserIn(ctx context.Context, repo UserRepository) (result *User, err error) {
	// Find a matched account with credential email
	result, err = repo.GetOneByEmail(ctx, e)
	if err != nil {
		// Otherwise ...
		return result, err
	}

	return result, err
}
