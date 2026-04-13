package resetpassword_test

import (
	"context"
	"errors"
	"testing"

	"everlasting/src/domain/mocks"
	"everlasting/src/domain/sharedkernel/identity"
	userDomain "everlasting/src/domain/user"
	"everlasting/src/domain/user/resetpassword"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TokenTestSuite struct {
	suite.Suite
	tokenRepoMock *mocks.TokenRepository
}

func (suite *TokenTestSuite) SetupTest() {
	suite.tokenRepoMock = mocks.NewTokenRepository(suite.T())
}

func (suite *TokenTestSuite) TestWithErrorResult() {
	// Initiate input
	ctx := context.Background()
	uid := identity.NewID()
	user := &userDomain.User{
		ID: uid,
	}
	errRepo := errors.New("any error")

	// Define token repo mock behaviour
	suite.
		tokenRepoMock.
		On("Generate", ctx, userDomain.TokenSubjectResetPassword, uid, mock.Anything).
		Return(userDomain.Token(""), errRepo)

	token, err := resetpassword.GenerateTokenWith(ctx, user, suite.tokenRepoMock)
	suite.Equal(resetpassword.Token(""), token)
	suite.Equal(errRepo, err)
}

func (suite *TokenTestSuite) TestGenerateAndVerify() {
	// Initiate input
	ctx := context.Background()
	uid := identity.NewID()
	user := &userDomain.User{
		ID: uid,
	}
	generatedToken := "asolole"

	// Define token repo mock behaviour
	suite.
		tokenRepoMock.
		On("Generate", ctx, userDomain.TokenSubjectResetPassword, uid, mock.Anything).
		Return(userDomain.Token(generatedToken), nil)

	suite.
		tokenRepoMock.
		On("Verify", ctx, userDomain.Token(generatedToken)).
		Return(uid, nil)

	suite.
		tokenRepoMock.
		On("Revoke", ctx, userDomain.Token(generatedToken)).
		Return(nil)

	token, err := resetpassword.GenerateTokenWith(ctx, user, suite.tokenRepoMock)
	suite.Equal(resetpassword.Token(generatedToken), token)
	suite.Equal(nil, err)

	result, err := token.VerifyWith(ctx, suite.tokenRepoMock)
	suite.Equal(userDomain.UserID(uid.String()), result)
	suite.Equal(nil, err)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestTokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}
