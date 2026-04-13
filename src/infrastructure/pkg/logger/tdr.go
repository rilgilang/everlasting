package logger

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

func GenerateTdrLogConfig(log *logrus.Logger) middleware.RequestLoggerConfig {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(formatter)

	return middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogHost:     true,
		LogMethod:   true,
		LogURIPath:  true,
		LogRemoteIP: true,
		LogReferer:  true,
		LogLatency:  true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			go func() {
				if err := recover(); err != nil {
					fmt.Printf("recovered: %v", err)
				}

				var payload string
				payload, _ = c.Get("request_body").(string)

				log.WithFields(logrus.Fields{
					"status":                  c.Response().Status,
					"host":                    values.Host,
					"method":                  values.Method,
					"uri_path":                values.URIPath,
					"src_ip":                  values.RemoteIP,
					"referer":                 values.Referer,
					"start_time":              values.StartTime,
					"latency_in_microseconds": values.Latency.Microseconds(),
					"body_payload":            payload,
				}).Info("request")
			}()
			return nil
		},
	}
}
