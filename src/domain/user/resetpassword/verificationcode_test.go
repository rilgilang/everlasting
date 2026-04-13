package resetpassword_test

import (
	"context"
	"testing"

	"everlasting/src/domain/mocks"
	"everlasting/src/domain/sharedkernel/identity"
	"everlasting/src/domain/user"
	"everlasting/src/domain/user/resetpassword"

	"github.com/stretchr/testify/suite"
)

type GenerateAndVerifyTestSuite struct {
	suite.Suite
	keyGeneratorMock *mocks.KeyGenerator
}

func (suite *GenerateAndVerifyTestSuite) SetupTest() {
	suite.keyGeneratorMock = mocks.NewKeyGenerator(suite.T())
}

func (suite *GenerateAndVerifyTestSuite) TestGenerateAndVerify() {
	// define params
	ctx := context.Background()
	uid := identity.NewID()
	matchedUser := &user.User{
		ID: uid,
	}
	generatedToken := "generated"

	userRepo := mocks.NewUserRepository(suite.T())
	userRepo.
		On("GetOneByID", ctx, uid).
		Return(matchedUser, nil)

	tokenRepo := mocks.NewTokenRepository(suite.T())
	tokenRepo.
		On("Generate", ctx, resetpassword.TokenSubjectVerificationCode, uid, user.TokenOption{}).
		Return(user.Token(generatedToken), nil)

	tokenRepo.
		On("Verify", ctx, user.Token(generatedToken)).
		Return(uid, nil)

	tokenRepo.
		On("Revoke", ctx, user.Token(generatedToken)).
		Return(nil)

	veriCode, err := resetpassword.GenerateVerificationCodeWith(ctx, matchedUser, tokenRepo)
	tokenRepo.AssertCalled(suite.T(), "Generate", ctx, resetpassword.TokenSubjectVerificationCode, uid, user.TokenOption{})
	suite.Nil(err)
	suite.NotNil(veriCode)

	usr, err := veriCode.VerifyWith(ctx, tokenRepo, userRepo)
	tokenRepo.AssertCalled(suite.T(), "Verify", ctx, user.Token(generatedToken))
	tokenRepo.AssertCalled(suite.T(), "Revoke", ctx, user.Token(generatedToken))
	userRepo.AssertCalled(suite.T(), "GetOneByID", ctx, uid)
	suite.Nil(err)
	suite.Equal(matchedUser, usr)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestVerificationCodeTestSuite(t *testing.T) {
	suite.Run(t, new(GenerateAndVerifyTestSuite))
}
