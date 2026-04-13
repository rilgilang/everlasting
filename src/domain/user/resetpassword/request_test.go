package resetpassword_test

import (
	"context"
	"errors"
	"testing"

	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/mocks"
	"everlasting/src/domain/user"
	"everlasting/src/domain/user/resetpassword"

	"github.com/stretchr/testify/suite"
)

type (
	produce struct {
		called bool
		err    error
	}

	putInstructionQueueInTestCases struct {
		testName string
		produce  produce
		err      error
	}

	PutInstructionQueueInTestSuite struct {
		suite.Suite
		testCases []putInstructionQueueInTestCases
	}
)

func (suite *PutInstructionQueueInTestSuite) SetupTest() {
	err := errors.New("any error")
	suite.testCases = []putInstructionQueueInTestCases{
		{
			testName: "test with error result",
			produce:  produce{called: true, err: err},
			err:      err,
		},
		{
			testName: "test with successful result",
			produce:  produce{called: true, err: nil},
			err:      nil,
		},
	}
}

func (suite *PutInstructionQueueInTestSuite) TestPutInQueue() {
	for _, test := range suite.testCases {
		suite.Run(test.testName, func() {
			// args preparation
			ctx := context.Background()
			email := user.Email("xxx@yyy.com")
			request := resetpassword.NewResetPasswordRequest(email)

			messageBroker := mocks.NewMessageBroker(suite.T())
			if test.produce.called {
				messageBroker.
					On("Produce", ctx, resetpassword.TaskSendResetPasswordRequest, request).
					Return(
						test.produce.err,
					)
			}
			// running function
			err := request.PutInstructionQueueIn(ctx, messageBroker)

			// post running assertion
			if test.produce.called {
				messageBroker.AssertCalled(suite.T(), "Produce", ctx, resetpassword.TaskSendResetPasswordRequest, request)
			}

			suite.Equal(test.err, err)
		})
	}
}

type (
	getOneByEmail struct {
		result *user.User
		err    error
	}

	isHasMatchedUserInTestCase struct {
		testName      string
		getOneByEmail getOneByEmail
		result        bool
		err           error
	}

	IsHasMatchedUserInTestSuite struct {
		suite.Suite
		testCases []isHasMatchedUserInTestCase
	}
)

func (suite *IsHasMatchedUserInTestSuite) SetupTest() {
	err := errors.New("any error")
	suite.testCases = []isHasMatchedUserInTestCase{
		{
			testName:      "test with error result",
			getOneByEmail: getOneByEmail{result: nil, err: err},
			result:        false,
			err:           err,
		},
		{
			testName:      "test with not found result",
			getOneByEmail: getOneByEmail{result: nil, err: errDomain.ErrDataNotFound},
			result:        false,
			err:           errDomain.ErrDataNotFound,
		},
		{
			testName:      "test with successful result",
			getOneByEmail: getOneByEmail{result: user.NewUser("Markonah", user.Email("xxx@yyy.zzz")), err: nil},
			result:        true,
			err:           nil,
		},
	}
}

func (suite *IsHasMatchedUserInTestSuite) TestIsHasMatchedUserIn() {
	for _, test := range suite.testCases {
		suite.Run(test.testName, func() {
			// args preparation
			ctx := context.Background()
			email := user.Email("xxx@yyy.com")

			// Define user repo mock behaviour
			userRepoMock := mocks.NewUserRepository(suite.T())
			userRepoMock.On("GetOneByEmail", ctx, email).Return(test.getOneByEmail.result, test.getOneByEmail.err)

			request := resetpassword.NewResetPasswordRequest(email)
			result, err := request.IsHasMatchedUserIn(ctx, userRepoMock)

			suite.Equal(test.result, result)
			suite.Equal(test.err, err)
		})
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestRequestTestSuite(t *testing.T) {
	suite.Run(t, new(PutInstructionQueueInTestSuite))
	suite.Run(t, new(IsHasMatchedUserInTestSuite))
}
