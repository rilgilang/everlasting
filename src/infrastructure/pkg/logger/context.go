package logger

import (
	"context"

	messagebrokerDomain "everlasting/src/domain/sharedkernel/messagebroker"

	"github.com/labstack/echo/v4"
)

type (
	ContextKey      string
	RequestProtocol string
)

const (
	AppLoggerContextKey ContextKey = "request_context_key"

	RequestProtocolHTTP = "http"
	RequestProtocolAMQP = "amqp"
)

type AppLoggerContext struct {
	Proto         RequestProtocol
	RequestID     string `json:"request_id"`
	IP            string `json:"src_ip"`
	RequestPath   string `json:"request_uri"`
	RequestMethod string `json:"request_method"`
}

func NewAppLoggerContextFromEchoContext(cc echo.Context) *AppLoggerContext {
	return &AppLoggerContext{
		Proto:         RequestProtocolHTTP,
		RequestID:     cc.Request().Header.Get("X-Request-ID"),
		IP:            cc.RealIP(),
		RequestPath:   cc.Request().RequestURI,
		RequestMethod: cc.Request().Method,
	}
}

func NewAppLoggerContextFromAMQPEvent(event messagebrokerDomain.TaskName) *AppLoggerContext {
	return &AppLoggerContext{
		Proto:       RequestProtocolAMQP,
		RequestPath: string(event),
	}
}

func (c *AppLoggerContext) GetContext() context.Context {
	return context.WithValue(context.Background(), AppLoggerContextKey, c)
}
