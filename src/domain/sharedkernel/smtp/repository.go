package smtp

import "context"

type Payload struct {
	From       string
	Recipients []string
	Subject    string
	Message    string
}

type SMTP interface {
	Send(ctx context.Context, payload Payload) (err error)
}
