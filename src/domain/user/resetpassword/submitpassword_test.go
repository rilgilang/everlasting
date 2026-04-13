package resetpassword_test

import (
	"context"
	"errors"
	"testing"

	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/mocks"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/user"
	"everlasting/src/domain/user/resetpassword"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type (
	SubmitPasswordRequestTestCase struct {
		uid  string
		user *user.User
	}

	SubmitPasswordRequestSuite struct {
		suite.Suite
		testCase SubmitPasswordRequestTestCase
	}
)

func (suite *SubmitPasswordRequestSuite) SetupTest() {
	uid := "3cf02954-ac39-4fc3-a818-abce93eb3991"
	suite.testCase = SubmitPasswordRequestTestCase{
		uid: uid,
		user: &user.User{
			ID: identity.FromStringOrNil(uid),
		},
	}
}

func (suite *SubmitPasswordRequestSuite) TestWithNotFoundUser() {
	expectedUID := user.UserID(suite.testCase.uid)
	ctx := context.WithValue(context.Background(), resetpassword.UIDContextKey, expectedUID)
	uidUUID := identity.FromStringOrNil(suite.testCase.uid)

	userRepo := mocks.NewUserRepository(suite.T())
	userRepo.On("GetOneByID", ctx, uidUUID).Return(nil, errDomain.ErrDataNotFound)

	input := &resetpassword.SubmitPasswordRequest{
		Password:        "~1234QweR!",
		PasswordConfirm: "~1234QweR!",
	}

	err := input.SaveTo(ctx, userRepo)
	userRepo.AssertCalled(suite.T(), "GetOneByID", ctx, uidUUID)
	suite.Equal(errDomain.ErrDataNotFound, err)
}

func (suite *SubmitPasswordRequestSuite) TestWithFailedUpdate() {
	expectedUID := user.UserID(suite.testCase.uid)
	ctx := context.WithValue(context.Background(), resetpassword.UIDContextKey, expectedUID)
	uidUUID := identity.FromStringOrNil(suite.testCase.uid)
	expectedError := errors.New("any error")

	userRepo := mocks.NewUserRepository(suite.T())
	userRepo.On("GetOneByID", ctx, uidUUID).Return(suite.testCase.user, nil)
	userRepo.On("UpdateByID", ctx, mock.Anything, uidUUID).Return(nil, expectedError)

	input := &resetpassword.SubmitPasswordRequest{
		Password:        "~1234QweR!",
		PasswordConfirm: "~1234QweR!",
	}

	err := input.SaveTo(ctx, userRepo)
	userRepo.AssertCalled(suite.T(), "GetOneByID", ctx, uidUUID)
	userRepo.AssertCalled(suite.T(), "UpdateByID", ctx, mock.Anything, uidUUID)
	suite.Equal(expectedError, err)
}

func (suite *SubmitPasswordRequestSuite) TestWithSuccesfulResult() {
	expectedUID := user.UserID(suite.testCase.uid)
	ctx := context.WithValue(context.Background(), resetpassword.UIDContextKey, expectedUID)
	uidUUID := identity.FromStringOrNil(suite.testCase.uid)

	userRepo := mocks.NewUserRepository(suite.T())
	userRepo.On("GetOneByID", ctx, uidUUID).Return(suite.testCase.user, nil)
	userRepo.On("UpdateByID", ctx, mock.Anything, uidUUID).Return(suite.testCase.user, nil)

	input := &resetpassword.SubmitPasswordRequest{
		Password:        "~1234QweR!",
		PasswordConfirm: "~1234QweR!",
	}

	err := input.SaveTo(ctx, userRepo)
	userRepo.AssertCalled(suite.T(), "GetOneByID", ctx, uidUUID)
	userRepo.AssertCalled(suite.T(), "UpdateByID", ctx, mock.Anything, uidUUID)
	suite.Nil(err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSubmitPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(SubmitPasswordRequestSuite))
}
