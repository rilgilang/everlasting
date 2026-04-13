package user_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	errDomain "everlasting/src/domain/error"
	"everlasting/src/domain/mocks"
	"everlasting/src/domain/user"

	"github.com/stretchr/testify/suite"
)

type (
	verifyTestCase struct {
		testName      string
		input         *user.Credential
		getOneByEmail getOneByEmailResult
		err           error
	}

	VerifyTestSuite struct {
		suite.Suite
		testCases []verifyTestCase
	}
)

func (suite *VerifyTestSuite) SetupTest() {
	validEmail := user.Email("markonah@digitalsekuriti.id")
	validPassword := user.Password("Asolole123!")

	validUser := user.NewUser("Markonah", user.Email("markonah@digitalsekuriti.id"))
	validUser.SetPassword(validPassword)
	lockoutRemaining := 5 * time.Second
	errLoginAttemps := fmt.Errorf("User is locked out. Remaining time: %s", lockoutRemaining)
	err := errors.New("any error")
	suite.testCases = []verifyTestCase{
		{
			testName: "Test with not found email",
			input: &user.Credential{
				Email:    validEmail,
				Password: validPassword,
			},
			getOneByEmail: getOneByEmailResult{
				result: nil,
				err:    errDomain.ErrDataNotFound,
			},
		},
		{
			testName: "Test with persistence error",
			input: &user.Credential{
				Email:    validEmail,
				Password: validPassword,
			},
			getOneByEmail: getOneByEmailResult{
				result: nil,
				err:    err,
			},
		},
		{
			testName: "Test with error login attemps",
			input: &user.Credential{
				Email:    validEmail,
				Password: validPassword,
			},
			getOneByEmail: getOneByEmailResult{
				result: nil,
				err:    errLoginAttemps,
			},
		},
		{
			testName: "Test with invalid password",
			input: &user.Credential{
				Email:    validEmail,
				Password: "InvalidPassword!234",
			},
			getOneByEmail: getOneByEmailResult{
				result: validUser,
				err:    nil,
			},
		},
		{
			testName: "Test with successful verify",
			input: &user.Credential{
				Email:    validEmail,
				Password: validPassword,
			},
			getOneByEmail: getOneByEmailResult{
				result: validUser,
				err:    nil,
			},
		},
	}
}

func (suite *VerifyTestSuite) TestVerify() {
	for _, test := range suite.testCases {
		suite.Run(test.testName, func() {
			// Define input
			ctx := context.Background()

			repository := mocks.NewUserRepository(suite.T())
			repository.On("GetOneByEmail", ctx, test.input.Email).Return(
				test.getOneByEmail.result,
				test.getOneByEmail.err,
			)
			_, err := test.input.VerifyWith(ctx, repository)
			if test.err != nil {
				suite.Equal(err.Error(), test.err.Error())
			}
		})
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestVerify(t *testing.T) {
	suite.Run(t, new(VerifyTestSuite))
}
