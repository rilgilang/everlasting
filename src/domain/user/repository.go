package user

import (
	"context"

	"everlasting/src/domain/sharedkernel/identity"
)

type UserRepository interface {
	GetOneByEmail(ctx context.Context, email Email) (*User, error)
	GetOneByID(ctx context.Context, id identity.ID) (*User, error)
	UpdateByID(ctx context.Context, user *User, id identity.ID) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	GetByQuery(ctx context.Context, query *Query) (*Users, error)
}
