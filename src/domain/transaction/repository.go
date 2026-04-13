package transaction

import (
	"context"

	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/wallet"
)

type TransactionRepository interface {
	GetOneByID(ctx context.Context, id identity.ID) (*Transaction, error)
	GetOneByReference(ctx context.Context, refType TransactionRefType, refID string) (*Transaction, error)
	Create(ctx context.Context, trx *Transaction) (*Transaction, error)
	Update(ctx context.Context, trx *Transaction) (*Transaction, error)
	GetByQuery(ctx context.Context, walletID wallet.WalletID, query *Query) (*Transactions, error)
}
