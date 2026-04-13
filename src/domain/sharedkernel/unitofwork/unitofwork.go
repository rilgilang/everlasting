package unitofwork

import (
	"context"

	"everlasting/src/domain/sharedkernel/datatype"
)

const TransactionContextKey datatype.ContextKey = "transaction_context"

type (
	Result struct {
		Body interface{}
	}
	UnitOfWork interface {
		Execute(ctx context.Context, fun func(ctx context.Context) (result *Result, err error)) (result *Result, err error)
	}
)
