package smtp

import (
	"context"

	"everlasting/src/domain/sharedkernel/smtp"
	"everlasting/src/infrastructure/pkg/logger"

	"gopkg.in/mail.v2"
)

type SmtpClients struct {
	dialer *mail.Dialer
	logger *logger.AppLogger
}

func NewSendEmail(dialer *mail.Dialer, logger *logger.AppLogger) *SmtpClients {
	return &SmtpClients{
		dialer,
		logger,
	}
}

func (smtp *SmtpClients) Send(ctx context.Context, payload smtp.Payload) error {
	m := mail.NewMessage()
	m.SetHeader("From", payload.From)
	m.SetHeader("To", payload.Recipients...)
	m.SetHeader("Subject", payload.Subject)
	m.SetBody("text/html", payload.Message)

	// Send the email
	err := smtp.dialer.DialAndSend(m)
	if err != nil {
		smtp.logger.Error(ctx, "smtp:failed_to_send", err.Error())
	}
	return err
}
