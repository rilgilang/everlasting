package event

import (
	"context"
	"everlasting/src/domain/sharedkernel/identity"
)

type EventRepository interface {
	GetOneByID(ctx context.Context, id identity.ID) (*Event, error)
	GetOneByTitle(ctx context.Context, title string) (*Event, error)
	Create(ctx context.Context, event *Event) (*Event, error)
	Update(ctx context.Context, event *Event) (*Event, error)
	GetByQuery(ctx context.Context, query *Query) (*Events, error)
}

type WishingWallMessageRepository interface {
	Create(ctx context.Context, message *WishingWallMessage) (*WishingWallMessage, error)
	GetAllByEventID(ctx context.Context, id string) ([]WishingWallMessage, error)
}
