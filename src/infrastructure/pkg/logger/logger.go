package logger

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"everlasting/src/infrastructure/pkg"

	splunkHook "github.com/Franco-Poveda/logrus-splunk-hook"
	"github.com/sirupsen/logrus"
	rotateFileHook "github.com/snowzach/rotatefilehook"
)

type (
	LogType string
)

const (
	LogTypeApp LogType = "app"
	LogTypeTdr LogType = "tdr"
)

var formatter = &logrus.JSONFormatter{
	FieldMap: logrus.FieldMap{
		logrus.FieldKeyMsg:  "message",
		logrus.FieldKeyTime: "timestamp",
	},
	TimestampFormat: time.RFC3339,
}

func NewLogrus(config *pkg.Config, logType LogType) (log *logrus.Logger) {
	log = logrus.New()

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(formatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the info severity or above.
	log.SetLevel(logrus.InfoLevel)

	if config.LogRotateActive {
		registerRotateFileHook(log, config, logType)
	}

	if config.LogSplunkActive {
		registerSplunkHook(log, config, logType)
	}

	return log
}

func registerSplunkHook(log *logrus.Logger, config *pkg.Config, logType LogType) {
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	sourceType := config.LogSplunkSourceTypeApp
	if logType == LogTypeTdr {
		sourceType = config.LogSplunkSourceTypeTdr
	}

	httpClient := &http.Client{Transport: httpTransport}

	client := splunkHook.NewClient(
		httpClient,
		fmt.Sprintf("%s?channel=%s", config.LogSplunkUrl, config.LogSplunkChannel),
		config.LogSplunkToken,
		config.LogSplunkSource,
		sourceType,
		config.LogSplunkIndex,
	)

	hook := splunkHook.NewHook(client, logrus.AllLevels)

	log.AddHook(hook)
}

func registerRotateFileHook(log *logrus.Logger, config *pkg.Config, logType LogType) {
	// Register Rotate File Hook
	logfile := config.LogRotateAppFile
	if logType == LogTypeTdr {
		logfile = config.LogRotateTdrFile
	}

	hook, err := rotateFileHook.NewRotateFileHook(rotateFileHook.RotateFileConfig{
		Filename:   logfile,
		MaxSize:    50, // megabytes
		MaxBackups: 3,  // amouts
		MaxAge:     28, // days
		Level:      logrus.InfoLevel,
		Formatter:  formatter,
	})

	if err != nil {
		log.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	log.AddHook(hook)
}
