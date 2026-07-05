package guest

import (
	"context"
	"time"

	"everlasting/src/domain/sharedkernel/identity"
)

const (
	GuestStatusVIP     = "vip"
	GuestStatusReguler = "regular"
)

type Guest struct {
	ID                 string     `json:"id"`
	EventId            string     `json:"user_id"`
	Name               string     `json:"name"`
	PhoneNumber        string     `json:"phone_number"`
	Address            string     `json:"address"`
	Status             string     `json:"status"`
	IsInvitationSent   bool       `json:"invitation_sended"`
	LastInvitationSent *time.Time `json:"last_invitation_sent"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}

type GuestInput struct {
	EventId     string `json:"user_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Status      string `json:"status" validate:"omitempty,oneof=vip regular"`
}

func (g *GuestInput) SaveTo(ctx context.Context, repo GuestRepository) (*Guest, error) {
	now := time.Now().UTC()
	guest := &Guest{
		ID:                 identity.NewID().String(),
		EventId:            g.EventId,
		Name:               g.Name,
		PhoneNumber:        g.PhoneNumber,
		Address:            g.Address,
		Status:             g.Status,
		IsInvitationSent:   false,
		LastInvitationSent: nil,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if g.Status == "" {
		guest.Status = GuestStatusReguler
	}

	return repo.Create(ctx, guest)
}

func (g *GuestInput) UpdateTo(ctx context.Context, repo GuestRepository, guestId string) (*Guest, error) {
	existing, err := repo.GetOneByID(ctx, guestId)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	existing.EventId = g.EventId
	existing.Name = g.Name
	existing.PhoneNumber = g.PhoneNumber
	existing.Address = g.Address
	existing.Status = g.Status
	existing.UpdatedAt = now

	return repo.Update(ctx, existing)
}

func (id GuestID) GetDetailFrom(ctx context.Context, repo GuestRepository) (*Guest, error) {
	return repo.GetOneByID(ctx, string(id))
}

type GuestID string

func (id GuestID) DeleteFrom(ctx context.Context, repo GuestRepository) error {
	return repo.Delete(ctx, string(id))
}

type (
	GuestSortBy string
	GuestOrder  string
)

const (
	GuestSortByCreatedAt GuestSortBy = "created_at"
	GuestSortByName      GuestSortBy = "name"

	GuestOrderAsc  GuestOrder = "asc"
	GuestOrderDesc GuestOrder = "desc"
)

type (
	Query struct {
		EventId string      `query:"event_id"`
		Status  string      `query:"status" validate:"omitempty,oneof=vip regular"`
		Q       string      `query:"q"`
		SortBy  GuestSortBy `query:"sort_by" validate:"omitempty,oneof=created_at name"`
		Order   GuestOrder  `query:"order" validate:"omitempty,oneof=asc desc"`
		Page    int64       `query:"page"`
		PerPage int64       `query:"per_page"`
	}

	Pagination struct {
		CurrentPage int64 `json:"current_page"`
		MaxPage     int64 `json:"max_page"`
		TotalData   int64 `json:"total_data"`
	}

	Guests struct {
		Query      *Query     `json:"-"`
		Collection []Guest    `json:"collection"`
		Pagination Pagination `json:"pagination"`
	}
)

func (q *Query) CollectFrom(ctx context.Context, repo GuestRepository) (*Guests, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PerPage < 1 {
		q.PerPage = 20
	}
	return repo.GetByQuery(ctx, q)
}
