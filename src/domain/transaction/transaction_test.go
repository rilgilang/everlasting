package transaction_test

import (
	"context"
	"testing"

	constant "everlasting/src/domain/error"
	"everlasting/src/domain/mocks"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/transaction"

	"github.com/stretchr/testify/suite"
)

type CreateTransactionRequestTestSuite struct {
	suite.Suite
}

type GetDetailFromRepositoryTestSuite struct {
	suite.Suite
}

func (suite *GetDetailFromRepositoryTestSuite) TestSuccessfulGetDetailFromRepository() {
	// Define input
	ctx := context.Background()
	id := "f2fcea76-88d2-448f-b040-55922eee3c94"
	transactionID := transaction.TransactionID(id)
	transactionUUID := identity.FromStringOrNil(id)

	// arguments mock result
	repository := mocks.NewTransactionRepository(suite.T())
	repository.On("GetOneByID", ctx, transactionUUID).Return(nil, constant.ErrTransactionNotFound)

	// run test object
	result, err := transactionID.GetDetailFrom(ctx, repository)

	// assert result
	repository.AssertCalled(suite.T(), "GetOneByID", ctx, transactionUUID)
	suite.Nil(result)
	suite.Equal(err, constant.ErrTransactionNotFound)
}

func (suite *GetDetailFromRepositoryTestSuite) TestFailedGetDetailFromRepository() {
	// Define input
	ctx := context.Background()
	id := "invalid uuid string"
	transactionID := transaction.TransactionID(id)

	// arguments mock result
	repository := mocks.NewTransactionRepository(suite.T())

	// run test object
	result, err := transactionID.GetDetailFrom(ctx, repository)

	// assert result
	suite.Nil(result)
	suite.Equal(err, constant.ErrTransactionNotFound)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTransactiontTestSuite(t *testing.T) {
	suite.Run(t, new(CreateTransactionRequestTestSuite))
	suite.Run(t, new(GetDetailFromRepositoryTestSuite))
}
