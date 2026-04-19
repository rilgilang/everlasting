package user

import (
	"context"
	"time"

	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/sharedkernel/marshaler"
)

type (
	UserID     string
	UserStatus string
	UserRole   string
)

// User status definition
const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"

	UserRoleuser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

// User Input
type UserInput struct {
	Email Email  `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
	Role  string `json:"role" validate:"required,oneof=user user-assistant"`
}

func (ui *UserInput) SaveTo(ctx context.Context, repo UserRepository) (result *User, err error) {
	return NewUser(ui.Name, ui.Email, UserRole(ui.Role)).SaveTo(ctx, repo)
}

func (ui *UserInput) UpdateTo(ctx context.Context, repo UserRepository, id UserID) (result *User, err error) {
	now := time.Now().UTC()
	user, err := id.GetDetailFrom(ctx, repo)
	if err != nil {
		return result, err
	}

	user.Name = ui.Name
	user.Email = ui.Email
	user.UpdatedAt = marshaler.JsonTime(now)

	return user.UpdateTo(ctx, repo)
}

// User data structure
type User struct {
	ID         identity.ID        `json:"_id"`
	Email      Email              `json:"email" validate:"required,email"`
	Name       string             `json:"name" validate:"required"`
	Role       UserRole           `json:"role"`
	Status     UserStatus         `json:"status"`
	CreatedAt  marshaler.JsonTime `json:"created_at"`
	UpdatedAt  marshaler.JsonTime `json:"updated_at"`
	CipherText CipherText         `json:"-"`
}

func NewUser(name string, email Email, role UserRole) *User {
	now := marshaler.JsonTime(time.Now())
	return &User{
		ID:        identity.NewID(),
		Email:     email,
		Name:      name,
		Role:      role,
		Status:    UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) With(id UserID) *User {
	u.ID = identity.FromStringOrNil(string(id))
	return u
}

func (u *User) UpdateTo(ctx context.Context, repo UserRepository) (user *User, err error) {
	if u.ID.IsNil() {
		return u, errDomain.ErrDataNotFound
	}
	return repo.UpdateByID(ctx, u, u.ID)
}

// By default new account definition does not contain CipherText (hashed version of Password) value. We can defined it here
func (acc *User) SetPassword(password Password) (err error) {
	chiperText, err := NewCipherTextFromPassword(password)
	if err != nil {
		return err
	}
	acc.CipherText = chiperText
	return err
}

// Verify password string to account CipherText
func (acc *User) VerifyPassword(password Password) (err error) {
	err = acc.CipherText.VerifyPassword(password)
	if err != nil {
		return errDomain.ErrorInvalidCredential
	}

	return err
}

func (acc *User) GenerateTokenWith(ctx context.Context, repo TokenRepository) (result *TokenSet, err error) {
	now := time.Now().UTC()
	accessTokenExp := now.Add(time.Duration(2*60*60) * time.Second) // Expired in 2 hours
	accessToken, err := repo.Generate(ctx, TokenSubjectAccessToken, acc.ID, TokenOption{
		ValidAt:   now.Unix(),
		ExpiredAt: accessTokenExp.Unix(),
	})

	if err != nil {
		return nil, err
	}

	refreshTokenExp := now.Add(time.Duration(48*60*60) * time.Second) // Expired in 2 days
	refreshToken, err := repo.Generate(ctx, TokenSubjectRefreshToken, acc.ID, TokenOption{
		ValidAt:   accessTokenExp.Unix(), // Valid once access token expired
		ExpiredAt: refreshTokenExp.Unix(),
	})

	if err != nil {
		return nil, err
	}

	responseToken := &TokenSet{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return responseToken, err
}

func (p *User) SaveTo(ctx context.Context, repo UserRepository) (*User, error) {
	return repo.Create(ctx, p)
}

func (uid UserID) GetDetailFrom(ctx context.Context, repo UserRepository) (user *User, err error) {

	id := identity.FromStringOrNil(string(uid))
	if id.IsNil() {
		return nil, errDomain.ErrDataNotFound
	}
	return repo.GetOneByID(ctx, id)
}

// Get list by query
type (
	UserSortBy string
	UserOrder  string
)

const (
	UserSortByCreatedAt       UserSortBy = "created_at"
	UserSortByName            UserSortBy = "name"
	UserSortByLastTransaction UserSortBy = "last_transaction"

	UserOrderAsc  UserOrder = "asc"
	UserOrderDesc UserOrder = "desc"
)

type (
	Query struct {
		Q       string     `query:"q"`
		Status  UserStatus `query:"status" validate:"omitempty,oneof=active inactive"`
		Role    UserRole   `query:"role" validate:"omitempty,oneof=superadmin admin"`
		SortBy  UserSortBy `query:"sort_by" validate:"omitempty,oneof=created_at name last_transaction"`
		Order   UserOrder  `query:"order" validate:"omitempty,oneof=asc desc"`
		Page    int64      `query:"page"`
		PerPage int64      `query:"per_page"`
	}

	Pagination struct {
		CurrentPage int64 `json:"current_page"`
		MaxPage     int64 `json:"max_page"`
		TotalData   int64 `json:"total_data"`
	}

	Users struct {
		Query      *Query     `json:"-"`
		Collection []User     `json:"collection"`
		Pagination Pagination `json:"pagination"`
	}
)

func (q *Query) CollectFrom(ctx context.Context, repo UserRepository) (Users *Users, err error) {
	return repo.GetByQuery(ctx, q)
}
