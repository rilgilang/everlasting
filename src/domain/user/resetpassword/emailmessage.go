package resetpassword

import (
	"context"

	"everlasting/src/domain/sharedkernel/renderer"
	"everlasting/src/domain/sharedkernel/smtp"
	userDomain "everlasting/src/domain/user"
)

const (
	EmailSubject string = "Reset Password Instruction"
)

type (
	EmailMessage struct {
		email userDomain.Email
		code  VerificationCode
	}

	RenderedMessage struct {
		smtp.Payload
	}
)

func NewEmailMessageFrom(email userDomain.Email, code VerificationCode) (result *EmailMessage) {
	return &EmailMessage{
		email,
		code,
	}
}

func (e *EmailMessage) RenderMessageWith(rendererLib renderer.HTMLRenderer, emailFrom string, clientBaseUrl string) (result *RenderedMessage, err error) {
	html, err := rendererLib.Render(renderer.ResetPasswordEmailTemplate, map[string]any{
		"base_url":          clientBaseUrl,
		"verification_code": e.code,
	})
	if err != nil {
		return result, err
	}

	return &RenderedMessage{
		Payload: smtp.Payload{
			From:       emailFrom,
			Recipients: []string{string(e.email)},
			Subject:    EmailSubject,
			Message:    html,
		},
	}, err
}

func (r *RenderedMessage) SendWith(ctx context.Context, s smtp.SMTP) (err error) {
	return s.Send(ctx, r.Payload)
}
