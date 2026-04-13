package user_test

import (
	"context"
	"errors"
	"testing"

	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/mocks"
	"everlasting/src/domain/user"

	"github.com/stretchr/testify/suite"
)

type (
	getOneByEmailResult struct {
		result *user.User
		err    error
	}

	getMatchedUserInTestCase struct {
		testName      string
		input         user.Email
		getOneByEmail getOneByEmailResult
		err           error
	}

	GetMatchedUserInTestSuite struct {
		suite.Suite
		testCases []getMatchedUserInTestCase
	}
)

func (suite *GetMatchedUserInTestSuite) SetupTest() {
	validEmail := user.Email("markonah@digitalsekuriti.id")
	validPassword := user.Password("Asolole123!")

	validUser := user.NewUser("Markonah", user.Email("markonah@digitalsekuriti.id"))
	validUser.SetPassword(validPassword)

	err := errors.New("any error")
	suite.testCases = []getMatchedUserInTestCase{
		{
			testName: "Test with not found email",
			input:    validEmail,
			getOneByEmail: getOneByEmailResult{
				result: nil,
				err:    errDomain.ErrDataNotFound,
			},
			err: errDomain.ErrDataNotFound,
		},
		{
			testName: "Test with persistence error",
			input:    validEmail,
			getOneByEmail: getOneByEmailResult{
				result: nil,
				err:    err,
			},
			err: err,
		},
		{
			testName: "Test with successful getMatchedUserIn",
			input:    validEmail,
			getOneByEmail: getOneByEmailResult{
				result: validUser,
				err:    nil,
			},
			err: nil,
		},
	}
}

func (suite *GetMatchedUserInTestSuite) TestGetMatchedUserIn() {
	for _, test := range suite.testCases {
		suite.Run(test.testName, func() {
			// Define input
			ctx := context.Background()

			repository := mocks.NewUserRepository(suite.T())
			repository.On("GetOneByEmail", ctx, test.input).Return(
				test.getOneByEmail.result,
				test.getOneByEmail.err,
			)

			_, err := test.input.GetMatchedUserIn(ctx, repository)
			repository.AssertCalled(suite.T(), "GetOneByEmail", ctx, test.input)
			if test.err != nil {
				suite.ErrorIs(err, test.err)
			}
		})
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCredentialTestSuite(t *testing.T) {
	suite.Run(t, new(GetMatchedUserInTestSuite))
}
