package pkg

import "fmt"

type Config struct {
	AppName          string `mapstructure:"APP_NAME"`
	AppPort          uint32 `mapstructure:"APP_PORT"`
	AppHost          string `mapstructure:"APP_HOST"`
	AppBaseUrl       string `mapstructure:"APP_BASEURL"`
	AppTimeZone      string `mapstructure:"APP_TIMEZONE"`
	AppClientBaseURL string `mapstructure:"APP_CLIENT_BASE_URL"`
	AppVersion       string `mapstructure:"APP_VERSION"`

	PostgresHost string `mapstructure:"POSTGRES_HOST"`
	PostgresPort uint32 `mapstructure:"POSTGRES_PORT"`
	PostgresUser string `mapstructure:"POSTGRES_USER"`
	PostgresPass string `mapstructure:"POSTGRES_PASS"`
	PostgresDB   string `mapstructure:"POSTGRES_DB"`
	PostgresSSL  string `mapstructure:"POSTGRES_SSL"`

	RedisPort uint32 `mapstructure:"REDIS_PORT"`
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPass string `mapstructure:"REDIS_PASS"`

	RabbitMQPort      uint32 `mapstructure:"RABBITMQ_PORT"`
	RabbitMQAdminPort uint32 `mapstructure:"RABBITMQ_ADMINPORT"`
	RabbitMQHost      string `mapstructure:"RABBITMQ_HOST"`
	RabbitMQUser      string `mapstructure:"RABBITMQ_USER"`
	RabbitMQPass      string `mapstructure:"RABBITMQ_PASS"`
	RabbitMQVhost     string `mapstructure:"RABBITMQ_VHOST"`
	RabbitMQQueueName string `mapstructure:"RABBITMQ_QUEUENAME"`

	MinioHost         string `mapstructure:"MINIO_HOST"`
	MinioSVCPort      string `mapstructure:"MINIO_SVC_PORT"`
	MinioAccessKey    string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioAccessSecret string `mapstructure:"MINIO_ACCESS_SECRET"`
	MinioSecure       bool   `mapstructure:"MINIO_SECURE"`

	SMTPEmailFrom string `mapstructure:"SMTP_EMAIL_FROM"`
	SMTPHost      string `mapstructure:"SMTP_HOST"`
	SMTPPort      uint32 `mapstructure:"SMTP_PORT"`
	SMTPEmail     string `mapstructure:"SMTP_EMAIL"`
	SMTPUsername  string `mapstructure:"SMTP_USERNAME"`
	SMTPPassword  string `mapstructure:"SMTP_PASSWORD"`

	LogRotateActive  bool   `mapstructure:"LOG_ROTATE_ACTIVE"`
	LogRotateAppFile string `mapstructure:"LOG_ROTATE_APP_FILE"`
	LogRotateTdrFile string `mapstructure:"LOG_ROTATE_TDR_FILE"`

	LogSplunkActive        bool   `mapstructure:"LOG_SPLUNK_ACTIVE"`
	LogSplunkUrl           string `mapstructure:"LOG_SPLUNK_URL"`
	LogSplunkChannel       string `mapstructure:"LOG_SPLUNK_CHANNEL"`
	LogSplunkToken         string `mapstructure:"LOG_SPLUNK_TOKEN"`
	LogSplunkIndex         string `mapstructure:"LOG_SPLUNK_INDEX"`
	LogSplunkSource        string `mapstructure:"LOG_SPLUNK_SOURCE"`
	LogSplunkSourceTypeApp string `mapstructure:"LOG_SPLUNK_SOURCE_TYPE_APP"`
	LogSplunkSourceTypeTdr string `mapstructure:"LOG_SPLUNK_SOURCE_TYPE_TDR"`

	WebsocketURL string `mapstructure:"WEBSOCKET_URL"`

	CORSAllowedDomain string `mapstructure:"CORS_ALLOWED_DOMAIN"`

	KeyTokenGenerator string `mapstructure:"KEY_TOKEN_GENERATOR"`

	TemplateDir string `mapstructure:"TEMPLATE_DIR"`

	BasicAuthAccount string `mapstructure:"BASIC_AUTH_ACCOUNT"`
}

func (c *Config) GetCORSAllowedDomain() string {
	if c.CORSAllowedDomain == "" {
		return "*"
	}

	return c.CORSAllowedDomain
}

func (c *Config) GenerateAMQPConnectionString() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%d/%s",
		c.RabbitMQUser,
		c.RabbitMQPass,
		c.RabbitMQHost,
		c.RabbitMQPort,
		c.RabbitMQVhost,
	)
}
