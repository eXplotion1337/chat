package postgres

import (
	"chat/pkg/wrap"
	"database/sql"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

type Config struct {
	Dsn                string
	DBPoolSize         int
	ServiceName        string
	MaxOpenConnections int
	MaxIdleConnections int
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
}

func (c Config) Validate() error {
	err := validation.ValidateStruct(&c,
		validation.Field(&c.Dsn, validation.Required),
		validation.Field(&c.DBPoolSize, validation.Required),
		validation.Field(&c.ServiceName, validation.Required),
		validation.Field(&c.MaxOpenConnections, validation.Required),
		validation.Field(&c.MaxIdleConnections, validation.Required),
		validation.Field(&c.ReadTimeout, validation.Required),
		validation.Field(&c.WriteTimeout, validation.Required),
	)

	if err != nil {
		return wrap.Errorf("failed to validate database config: %w", err)
	}

	return nil
}

func NewConnectDB(
	dsn,
	appName string,
	poolSize int,
	readTimeout time.Duration,
	writeTimeout time.Duration,
) (*bun.DB, error) {
	sqlDB := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithDSN(dsn),
			pgdriver.WithTLSConfig(nil),
			pgdriver.WithApplicationName(fmt.Sprintf("[%s]", appName)),
			pgdriver.WithReadTimeout(readTimeout),
			pgdriver.WithWriteTimeout(writeTimeout),
		),
	)

	sqlDB.SetMaxOpenConns(poolSize)
	sqlDB.SetMaxIdleConns(poolSize)

	db := bun.NewDB(sqlDB, pgdialect.New(), bun.WithDiscardUnknownColumns())

	return db, nil
}
