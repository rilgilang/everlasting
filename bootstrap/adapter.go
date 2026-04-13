package bootstrap

import (
	"context"
	"fmt"
	"time"

	"log"

	"everlasting/src/infrastructure/pkg"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/sarulabs/di"
)

func loadAdapter(builder *di.Builder, config *pkg.Config) {
	builder.Add([]di.Def{
		{
			Name:  "adapter.postgres",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				// Generate DSN string from config
				var generateConnectionString = func() string {
					return fmt.Sprintf(
						"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s application_name=%s",
						config.PostgresHost,
						config.PostgresPort,
						config.PostgresDB,
						config.PostgresUser,
						config.PostgresPass,
						config.PostgresSSL,
						config.AppName,
					)
				}

				db, err := sqlx.Connect("postgres", generateConnectionString())
				if err != nil {
					log.Printf("Error while initialize db provider. Detail: %s", err.Error())
					return nil, err
				}
				db.SetMaxOpenConns(5)
				db.SetConnMaxLifetime(time.Minute * 15)
				db.SetMaxIdleConns(5)
				return db, err
			},
			Close: func(obj interface{}) error {
				return obj.(*sqlx.DB).Close()
			},
		},
		{
			Name:  "provider.lock",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				client := goredislib.NewClient(&goredislib.Options{
					Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
					Password: config.RedisPass,
				})
				err := client.Set(context.Background(), "key", "value", 0).Err()
				if err != nil {
					log.Printf("Error while initialize redis provider. Detail: %s", err.Error())
					return nil, err
				}

				pool := goredis.NewPool(client)
				rs := redsync.New(pool)
				return rs, err
			},
		},
		{
			Name:  "provider.tz",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				timezone := config.AppTimeZone
				if timezone == "" {
					timezone = "Asia/Jakarta"
				}
				return time.LoadLocation(timezone)
			},
		},
		{
			Name:  "config.auth",
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return config.BasicAuthAccount, nil
			},
		},
	}...)
}
