package resetpassword_test

import (
	"context"
	"errors"
	"testing"

	"everlasting/src/domain/mocks"
	"everlasting/src/domain/sharedkernel/renderer"
	"everlasting/src/domain/user"
	"everlasting/src/domain/user/resetpassword"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RenderAndSendTestSuite struct {
	suite.Suite
	rendererMock *mocks.HTMLRenderer
	smtpMock     *mocks.SMTP
}

func (suite *RenderAndSendTestSuite) SetupTest() {
	suite.rendererMock = mocks.NewHTMLRenderer(suite.T())
	suite.smtpMock = mocks.NewSMTP(suite.T())
}

func (suite *RenderAndSendTestSuite) TestRenderAndSendWithErrorResult() {
	email := user.Email("xxx@yyy.zzz")
	code := resetpassword.VerificationCode("asolole")
	emailMessage := resetpassword.NewEmailMessageFrom(email, code)
	rendererError := errors.New("any error")

	suite.rendererMock.On("Render", renderer.ResetPasswordEmailTemplate, map[string]any{
		"verification_code": code,
		"base_url":          "http://base.url.com",
	}).Return("", rendererError)

	_, err := emailMessage.RenderMessageWith(suite.rendererMock, "ini_email_aicare_asolole@gmail.com", "http://base.url.com")
	suite.Equal(rendererError, err)
}

func (suite *RenderAndSendTestSuite) TestRenderAndSendWithSuccessfulResult() {
	email := user.Email("xxx@yyy.zzz")
	code := resetpassword.VerificationCode("asolole")
	emailMessage := resetpassword.NewEmailMessageFrom(email, code)

	suite.rendererMock.On("Render", renderer.ResetPasswordEmailTemplate, map[string]any{
		"verification_code": code,
		"base_url":          "http://base.url.com",
	}).Return("<html>asolole</html>", nil)

	suite.smtpMock.On("Send", context.Background(), mock.Anything).Return(nil)

	result, err := emailMessage.RenderMessageWith(suite.rendererMock, "ini_email_aicare_asolole@gmail.com", "http://base.url.com")
	suite.Equal("<html>asolole</html>", result.Message)
	suite.Equal(nil, err)
	ctx := context.Background()
	result.SendWith(ctx, suite.smtpMock)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestGenerateAndRenderAndSendTestSuite(t *testing.T) {
	suite.Run(t, new(RenderAndSendTestSuite))
}
