package logger

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

type AppLogger struct {
	log *logrus.Logger
}

func NewAppLogger(log *logrus.Logger) *AppLogger {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(formatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the info severity or above.
	log.SetLevel(logrus.InfoLevel)

	return &AppLogger{
		log: log,
	}
}

func (l *AppLogger) Debug(ctx context.Context, eventName, message string) {
	fields := getFields(ctx, eventName)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("recovered: %v", err)
			}
		}()
		l.log.WithFields(fields).Debug(message)
	}()
}

func (l *AppLogger) Info(ctx context.Context, eventName, message string) {
	fields := getFields(ctx, eventName)
	go func() {
		if err := recover(); err != nil {
			fmt.Printf("recovered: %v", err)
		}
		l.log.WithFields(fields).Info(message)
	}()
}

func (l *AppLogger) Warn(ctx context.Context, eventName, message string) {
	fields := getFields(ctx, eventName)
	go func() {
		if err := recover(); err != nil {
			fmt.Printf("recovered: %v", err)
		}
		l.log.WithFields(fields).Warn(message)
	}()
}

func (l *AppLogger) Error(ctx context.Context, eventName, message string) {
	fields := getFields(ctx, eventName)
	go func() {
		if err := recover(); err != nil {
			fmt.Printf("recovered: %v", err)
		}
		l.log.WithFields(fields).Error(message)
	}()
}

func (l *AppLogger) Fatal(ctx context.Context, eventName, message string) {
	fields := getFields(ctx, eventName)
	go func() {
		if err := recover(); err != nil {
			fmt.Printf("recovered: %v", err)
		}
		l.log.WithFields(fields).Fatal(message)
	}()
}

func (l *AppLogger) Panic(ctx context.Context, eventName, message string) {
	fields := getFields(ctx, eventName)
	go func() {
		if err := recover(); err != nil {
			fmt.Printf("recovered: %v", err)
		}
		l.log.WithFields(fields).Panic(message)
	}()
}

func getFields(ctx context.Context, eventName string) (result logrus.Fields) {
	var (
		proto   RequestProtocol
		request map[string]interface{}
		event   map[string]interface{}
	)

	loggerContext, ok := ctx.Value(AppLoggerContextKey).(*AppLoggerContext)
	if !ok {
		return result
	}

	proto = loggerContext.Proto
	request = map[string]interface{}{
		"ip":       loggerContext.IP,
		"endpoint": loggerContext.RequestPath,
		"method":   loggerContext.RequestMethod,
	}

	pc, file, lineNo, ok := runtime.Caller(2)
	if !ok {
		return result
	}
	functionName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file) // The Base function returns the last element of the path
	event = map[string]interface{}{
		"key":         eventName,
		"file":        fileName,
		"function":    functionName,
		"line_number": lineNo,
	}

	return logrus.Fields{"proto": proto, "request": request, "event": event}
}
