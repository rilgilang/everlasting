package bootstrap

import (
	"context"
	"crypto/tls"
	minioStorage "everlasting/src/infrastructure/pkg/storage"
	websocketPkg "everlasting/src/infrastructure/pkg/websocket"
	"fmt"
	"github.com/coder/websocket"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"

	"everlasting/src/infrastructure/amqp"
	"everlasting/src/infrastructure/persistence"
	"everlasting/src/infrastructure/pkg"
	"everlasting/src/infrastructure/pkg/cache"
	"everlasting/src/infrastructure/pkg/logger"
	"everlasting/src/infrastructure/pkg/smtp"
	"everlasting/src/infrastructure/pkg/token"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"gopkg.in/mail.v2"
)

func loadPersistence(builder *di.Builder, config *pkg.Config) {
	builder.Add([]di.Def{
		{
			Name:  "persistence.event",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				// Connect to persistence
				transactionPersistence := persistence.NewEventPersistence(
					ctn.Get("adapter.postgres").(*sqlx.DB),
					ctn.Get("logger.app").(*logger.AppLogger),
				)

				return transactionPersistence, nil
			},
		},
		{
			Name:  "persistence.user",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				// Connect to persistence
				transactionPersistence := persistence.NewUserPersistence(
					ctn.Get("adapter.postgres").(*sqlx.DB),
					ctn.Get("logger.app").(*logger.AppLogger),
				)

				return transactionPersistence, nil
			},
		},
		{
			Name:  "persistence.wishing_wall_message",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				// Connect to persistence
				transactionPersistence := persistence.NewWishingWallMessagePersistence(
					ctn.Get("adapter.postgres").(*sqlx.DB),
					ctn.Get("logger.app").(*logger.AppLogger),
				)

				return transactionPersistence, nil
			},
		},
		{
			Name:  "persistence.uow",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				// Connect to persistence
				uow := persistence.NewUnitOfWork(
					ctn.Get("adapter.postgres").(*sqlx.DB),
				)

				return uow, nil
			},
		},
	}...)
}

func loadPkg(builder *di.Builder, config *pkg.Config) {
	builder.Add([]di.Def{
		{
			Name:  "logrus.app",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				logger := logger.NewLogrus(config, logger.LogTypeApp)
				return logger, nil
			},
		},
		{
			Name:  "logrus.tdr",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				logger := logger.NewLogrus(config, logger.LogTypeTdr)
				return logger, nil
			},
		},
		{
			Name:  "logger.app",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				logger := logger.NewAppLogger(
					ctn.Get("logrus.app").(*logrus.Logger),
				)
				return logger, nil
			},
		},
		{
			Name:  "pkg.cache.redis",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				client := redis.NewClient(&redis.Options{
					Addr:     fmt.Sprintf("%s:%d", config.AppHost, config.RedisPort),
					Password: config.RedisPass,
				})
				err := client.Set(context.Background(), "key", "value", 0).Err()
				if err != nil {
					log.Printf("Error while initialize redis adapter. Detail: %s", err.Error())
					return nil, err
				}
				return cache.NewRedisCache(client), nil
			},
		},
		{
			Name:  "pkg.token",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return token.NewTokenProvider(
					ctn.Get("pkg.cache.redis").(*cache.RedisCache),
					config.KeyTokenGenerator,
					ctn.Get("logger.app").(*logger.AppLogger),
				)
			},
		},
		{
			Name:  "pkg.messagebroker.amqp",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				a := amqp.NewMessageBroker(
					config,
					ctn.Get("logger.app").(*logger.AppLogger),
				)
				return a, nil
			},
		},
		{
			Name:  "pkg.smtp.client",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				dialer := mail.NewDialer(config.SMTPHost, int(config.SMTPPort), config.SMTPUsername, config.SMTPPassword)
				dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
				logger := ctn.Get("logger.app").(*logger.AppLogger)

				smtp := smtp.NewSendEmail(dialer, logger)
				return smtp, nil
			},
		},
		{
			Name:  "pkg.renderer.html",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				templateDir := config.TemplateDir
				if templateDir == "" {
					templateDir = "template/"
				}
				renderer := pkg.NewHtmlRenderer(templateDir)
				return renderer, nil
			},
		},
		{
			Name:  "pkg.storage.minio",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				logger := ctn.Get("logger.app").(*logger.AppLogger)

				endpoint := config.MinioHost
				if config.MinioSVCPort != "" {
					endpoint = fmt.Sprintf("%s:%s", config.MinioHost, config.MinioSVCPort)
				}

				accessKey := config.MinioAccessKey
				secretAccessKey := config.MinioAccessSecret

				// Initialize minio client object.
				minioClient, err := minio.New(endpoint, &minio.Options{
					Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
					Secure: config.MinioSecure,
				})

				minioClient.MakeBucket(context.Background(), "wishing-wall", minio.MakeBucketOptions{
					Region:        "",
					ObjectLocking: false,
					ForceCreate:   false,
				})

				if err != nil {
					log.Printf("Error while initialize minio storage adapter. Detail: %s", err.Error())
					return nil, err
				}
				storage := minioStorage.NewMinioStorage(minioClient, logger)
				return storage, nil
			},
		},
		{
			Name:  "pkg.websocket.client",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				logger := ctn.Get("logger.app").(*logger.AppLogger)

				// Connect to the server
				url := "wss://s13783.blr1.piesocket.com/v3/1?api_key=7SEqHklfXLf4YSvF8OmgAd147ewDT0RT2tZrCE3f&notify_self=1"
				c, _, err := websocket.Dial(context.Background(), url, nil)
				if err != nil {
					log.Fatal(err)
				}
				//defer c.Close(websocket.StatusNormalClosure, "")

				client := websocketPkg.NewWebSocketClient(c, logger)
				return client, nil
			},
		},
	}...)
}
