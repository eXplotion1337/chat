package model

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:public.users,alias:u"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Nick          string    `bun:"nick"`
	PasswordHash  string    `bun:"passwordHash"`
	CreatedAt     time.Time `bun:"createdAt"`
}

type UserRepo struct {
	db *bun.DB
}

func NewUserRepo(
	dsn,
	appName string,
	poolSize int,
	readTimeout time.Duration,
	writeTimeout time.Duration,
) (*UserRepo, error) {
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

	return &UserRepo{db}, nil
}
