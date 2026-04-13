package persistence

import (
	"context"

	domain "everlasting/src/domain/sharedkernel/unitofwork"

	"github.com/jmoiron/sqlx"
)

type UnitOfWork struct {
	db *sqlx.DB
}

func NewUnitOfWork(db *sqlx.DB) *UnitOfWork {
	return &UnitOfWork{
		db,
	}
}

func (u *UnitOfWork) Execute(ctx context.Context, fun func(ctx context.Context) (result *domain.Result, err error)) (result *domain.Result, err error) {
	tx, err := u.db.Beginx()
	if err != nil {
		return result, err
	}

	ctx = context.WithValue(ctx, domain.TransactionContextKey, tx)

	result, err = fun(ctx)
	if err != nil {
		tx.Rollback()
		return result, err
	}

	tx.Commit()
	return result, err
}
