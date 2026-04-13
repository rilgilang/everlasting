package validator_test

import (
	"testing"

	"everlasting/src/domain/user"
	"everlasting/src/domain/validator"

	"github.com/stretchr/testify/suite"

	libValidator "github.com/go-playground/validator"
)

type (
	passwordValidationTestCase struct {
		testName      string
		input         string
		shouldBeValid bool
	}

	PasswordValidationTestSuite struct {
		suite.Suite
		testCases []passwordValidationTestCase
	}
)

func (suite *PasswordValidationTestSuite) SetupTest() {
	suite.testCases = []passwordValidationTestCase{
		{
			testName:      "No uppercase password",
			input:         "nouppercase123!",
			shouldBeValid: false,
		},
		{
			testName:      "No number password",
			input:         "NoNumber!",
			shouldBeValid: false,
		},
		{
			testName:      "No special char password",
			input:         "NoSpecialChar123",
			shouldBeValid: false,
		},
		{
			testName:      "Valid Password",
			input:         "ValidPassword123!",
			shouldBeValid: true,
		},
	}
}

func (suite *PasswordValidationTestSuite) TestValidatePassword() {
	for _, test := range suite.testCases {
		suite.Run(test.testName, func() {
			myPassword := user.Password(test.input)
			validate := libValidator.New()
			validate.RegisterValidation("password_custom_validator", validator.PasswordCustomValidator)
			error := validate.Var(myPassword, "required,password_custom_validator")
			suite.Equal(test.shouldBeValid, error == nil)
		})
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestCredentialTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordValidationTestSuite))
}
