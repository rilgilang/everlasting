package guest

import "context"

type GuestRepository interface {
	GetOneByID(ctx context.Context, id string) (*Guest, error)
	Create(ctx context.Context, guest *Guest) (*Guest, error)
	Update(ctx context.Context, guest *Guest) (*Guest, error)
	Delete(ctx context.Context, id string) error
	GetByQuery(ctx context.Context, query *Query) (*Guests, error)
}
