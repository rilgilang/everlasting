package middleware

type MiddlewareValue string

const (
	MiddlewareValueUserID           MiddlewareValue = "uid"
	MiddlewareValueContainer        MiddlewareValue = "container"
	MiddlewareValueAppLoggerContext MiddlewareValue = "request_context"
)
