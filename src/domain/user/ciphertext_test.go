package user_test

import (
	"testing"

	"everlasting/src/domain/user"

	"github.com/stretchr/testify/suite"
)

type (
	passwordVerificationTestCase struct {
		testName        string
		input           string
		plainText       string
		shouldBeMatched bool
	}

	PasswordVerificationTestSuite struct {
		suite.Suite
		testCases []passwordVerificationTestCase
	}
)

func (suite *PasswordVerificationTestSuite) SetupTest() {
	suite.testCases = []passwordVerificationTestCase{
		{
			testName:        "Unmatched Password",
			input:           "matched_password",
			plainText:       "unmatched_password",
			shouldBeMatched: false,
		},
		{
			testName:        "Matched Password",
			input:           "1234qweR!",
			plainText:       "1234qweR!",
			shouldBeMatched: true,
		},
	}
}

func (suite *PasswordVerificationTestSuite) TestVerifyPassword() {
	for _, test := range suite.testCases {
		suite.Run(test.testName, func() {
			plainText := user.Password(test.input)
			cipherText, _ := user.NewCipherTextFromPassword(plainText)
			err := cipherText.VerifyPassword(user.Password(test.plainText))
			suite.Equal(test.shouldBeMatched, (err == nil))
		})
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestPasswordTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordVerificationTestSuite))
}
