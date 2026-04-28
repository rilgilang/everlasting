package event

import (
	"context"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/sharedkernel/messagebroker"
	"time"
)

const TaskSendWishingWallMessage messagebroker.TaskName = "send_wishing_wall_message"

type (
	EventID string
)

type (
	Query struct {
		UserId    string `json:"user_id"`
		Category  string `json:"category"`
		Status    string `json:"status"`
		Organizer string `json:"organizer"`
		Location  string `json:"location"`
		DateFrom  string `query:"date_from" validate:"date"`
		DateUntil string `query:"date_until" validate:"date"`
		Cursor    int64  `query:"cursor"`
		PerPage   int64  `query:"per_page"`
	}

	EventInput struct {
		Title       string `json:"title" validate:"required"`
		Description string `json:"description" validate:"required"`
		Date        string `json:"date" validate:"required"`
		Time        string `json:"time" validate:"required"`
		Location    string `json:"location" validate:"required"`
		Category    string `json:"category" validate:"required"`
		MaxMessages int    `json:"max_messages" validate:"required"`
		Status      string `json:"status" validate:"required"`
		Organizer   string `json:"organizer" validate:"required"`
	}

	Event struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Date        string    `json:"date"`
		Time        string    `json:"time"`
		Location    string    `json:"location"`
		Category    string    `json:"category"`
		Messages    int       `json:"messages"`
		MaxMessages int       `json:"max_messages"`
		Image       string    `json:"image"`
		Status      string    `json:"status"`
		Organizer   string    `json:"organizer"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	Pagination struct {
		NextCursor int64 `json:"next_cursor"`
	}

	Events struct {
		Query      *Query     `json:"-"`
		Collection []Event    `json:"collection"`
		Pagination Pagination `json:"pagination"`
	}
)

func (e *EventInput) SaveEvent(ctx context.Context, repository EventRepository) (*Event, error) {
	eventId := identity.NewID()

	now := time.Now()

	return repository.Create(ctx, &Event{
		ID:          eventId.String(),
		Title:       e.Title,
		Description: e.Description,
		Date:        e.Date,
		Time:        e.Time,
		Location:    e.Location,
		Category:    e.Category,
		Messages:    0,
		MaxMessages: e.MaxMessages,
		Image:       "",
		Status:      e.Status,
		Organizer:   e.Organizer,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

func (e *EventInput) UpdateTo(ctx context.Context, repository EventRepository, eventId EventID) (*Event, error) {
	now := time.Now()

	return repository.Update(ctx, &Event{
		ID:          string(eventId),
		Title:       e.Title,
		Description: e.Description,
		Date:        e.Date,
		Time:        e.Time,
		Location:    e.Location,
		Category:    e.Category,
		MaxMessages: e.MaxMessages,
		Image:       "",
		Status:      e.Status,
		Organizer:   e.Organizer,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
}

func (e *EventID) GetByID(ctx context.Context, repository EventRepository) (*Event, error) {
	return repository.GetOneByID(ctx, identity.FromStringOrNil(string(*e)))
}

func (q *Query) CollectFrom(ctx context.Context, repository EventRepository) (*Events, error) {
	return repository.GetByQuery(ctx, q)
}
