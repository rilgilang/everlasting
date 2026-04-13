package transaction

import (
	"context"
	"strings"
	"time"

	constant "everlasting/src/domain/error"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/sharedkernel/marshaler"
	"everlasting/src/domain/wallet"
)

type (
	TransactionStatus  string
	TransactionType    string
	TransactionRefType string
)

// Transaction attributes definition
const (
	TransactionStatusActive   TransactionStatus = "active"
	TransactionStatusInactive TransactionStatus = "inactive"

	TransactionTypeDebt   TransactionType = "debt"
	TransactionTypeCredit TransactionType = "credit"

	TransactionRefTypeOrder  TransactionRefType = "order"
	TransactionRefTypePayout TransactionRefType = "payout"
	TransactionRefTypeRefund TransactionRefType = "refund"
)

// Get Transaction by it's Transaction ID
type TransactionID string

func (p TransactionID) GetDetailFrom(ctx context.Context, repo TransactionRepository) (transaction *Transaction, err error) {
	id := identity.FromStringOrNil(string(p))
	if id.IsNil() {
		return nil, constant.ErrTransactionNotFound
	}
	return repo.GetOneByID(ctx, id)
}

type CreateTransactionRequest struct {
	Type     TransactionType            `json:"type" validate:"required,oneof=debt credit" example:"debt|credit"`
	Amount   int64                      `json:"amount" validate:"required,numeric" example:"220000"`
	WalletID string                     `json:"wallet_id" validate:"omitempty,uuid" example:"8f364610-6f12-4bed-b7d1-7ea1892803c7"`
	Wallet   wallet.CreateWalletRequest `json:"wallet" validate:"omitempty,required_without=WalletID"`
	RefType  TransactionRefType         `json:"ref_type" validate:"required,oneof=order payout refund" example:"order|payout|refund"`
	RefID    string                     `json:"ref_id" example:""`
	Notes    string                     `json:"notes" example:""`
}

func (ctr CreateTransactionRequest) VerifyAndTransform(ctx context.Context, walletRepo wallet.WalletRepository, transactionRepo TransactionRepository) (oTrans *Transaction, oWallet *wallet.Wallet, err error) {
	now := marshaler.JsonTime(time.Now().UTC())
	// Validate reference
	refID := strings.Trim(ctr.RefID, " ")
	_, err = transactionRepo.GetOneByReference(ctx, ctr.RefType, refID)
	if err != nil {
		if err != constant.ErrTransactionNotFound {
			return oTrans, oWallet, err
		}
	} else {
		// If transaction with same reference is exists
		return oTrans, oWallet, constant.ErrTransactionReferenceNotAvailable
	}

	if ctr.WalletID != "" {
		// If user use wallet id as wallet identifier
		oWallet, err = wallet.WalletID(ctr.WalletID).GetDetailFrom(ctx, walletRepo)
	} else {
		// If user use wallet input as wallet identifier
		oWallet, err = ctr.Wallet.GetOrSaveTo(ctx, walletRepo)
	}
	if err != nil {
		return oTrans, oWallet, err
	}

	if ctr.Type == TransactionTypeDebt && oWallet.Balance < ctr.Amount {
		return nil, oWallet, constant.ErrInsufficientBalance
	}

	return &Transaction{
		ID:        identity.NewID(),
		Type:      ctr.Type,
		Amount:    ctr.Amount,
		WalletID:  oWallet.ID,
		RefType:   ctr.RefType,
		RefID:     refID,
		Notes:     strings.Trim(ctr.Notes, " "),
		CreatedAt: now,
		UpdatedAt: now,
		Status:    TransactionStatusActive,
	}, oWallet, err
}

// Transaction data structure
type Transaction struct {
	ID        identity.ID        `json:"_id"`
	Type      TransactionType    `json:"type"`
	Amount    int64              `json:"amount"`
	WalletID  identity.ID        `json:"wallet_id"`
	RefType   TransactionRefType `json:"ref_type"`
	RefID     string             `json:"ref_id"`
	Notes     string             `json:"notes"`
	Status    TransactionStatus  `json:"status"`
	CreatedAt marshaler.JsonTime `json:"created_at"`
	UpdatedAt marshaler.JsonTime `json:"updated_at"`

	Meta                   interface{} `json:"-"`
	UpdatedToWalletBalance bool        `json:"-"`
}

// Save Transaction to repository
func (p *Transaction) SaveTo(ctx context.Context, repo TransactionRepository, walletRepo wallet.WalletRepository) (transaction *Transaction, err error) {
	return repo.Create(ctx, p)
}

func (p *Transaction) UpdateBalanceToWallet(ctx context.Context, matchedWallet *wallet.Wallet, repo TransactionRepository, walletRepo wallet.WalletRepository) (err error) {
	// If transaction already added to wallet balance, return immediately
	if p.UpdatedToWalletBalance {
		return err
	}

	// Update wallet balance
	matchedWallet.LastTransaction = p.CreatedAt
	if p.Type == TransactionTypeCredit {
		matchedWallet.Balance = matchedWallet.Balance + p.Amount
	} else {
		matchedWallet.Balance = matchedWallet.Balance - p.Amount
	}
	matchedWallet.LastTransaction = p.CreatedAt
	matchedWallet.UpdatedAt = p.CreatedAt

	_, err = walletRepo.Update(ctx, matchedWallet)
	if err != nil {
		return err
	}

	// // Set updated to wallet balance to true
	p.UpdatedToWalletBalance = true
	_, err = repo.Update(ctx, p)
	if err != nil {
		return err
	}

	return err
}

// Get list by query
type (
	Query struct {
		Type      TransactionType    `query:"type" validate:"omitempty,oneof=debt credit"`
		RefType   TransactionRefType `query:"ref_type" validate:"omitempty,oneof=order payout refund"`
		RefID     string             `query:"ref_id"`
		DateFrom  string             `query:"date_from" validate:"date"`
		DateUntil string             `query:"date_until" validate:"date"`
		Cursor    int64              `query:"cursor"`
		PerPage   int64              `query:"per_page"`
	}

	Pagination struct {
		NextCursor int64 `json:"next_cursor"`
	}

	Transactions struct {
		Query      *Query        `json:"-"`
		Collection []Transaction `json:"collection"`
		Pagination Pagination    `json:"pagination"`
	}
)

func (q *Query) CollectFrom(ctx context.Context, walletId wallet.WalletID, repo TransactionRepository) (transactions *Transactions, err error) {
	return repo.GetByQuery(ctx, walletId, q)
}
