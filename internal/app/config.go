package app

import (
	"chat/pkg/wrap"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"os"
	"path/filepath"
	"strings"
)

var ROOT_PATH string
var WEB_PATH string

type Config struct {
	HttpConfing    *Http
	PostgresConfig *Postgres
}

type Http struct {
	Port string
	Host string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	DbDriver string
	SslMode  string
}

func InitConfig() (*Config, error) {
	config := Config{
		HttpConfing: &Http{
			Port: getEnvAsStr("HTTP_PORT", "8080"),
			Host: getEnvAsStr("HTTP_HOST", "localhost"),
		},
		PostgresConfig: &Postgres{
			Host:     getEnvAsStr("DB_HOST", "localhost"),
			Port:     getEnvAsStr("DB_PORT", "5432"),
			User:     getEnvAsStr("DB_USER", "root"),
			Password: getEnvAsStr("DB_PASSWORD", "123123"),
			DbName:   getEnvAsStr("DB_NAME", "chat"),
			DbDriver: getEnvAsStr("DB_DRIVER", "postgres"),
			SslMode:  getEnvAsStr("SSL_MODE", "disable"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	ROOT_PATH = filepath.Dir("../../")
	WEB_PATH = filepath.Join(ROOT_PATH, "/web/")

	return &config, nil
}

func (c Config) Validate() error {
	if err := c.HttpConfing.Validate(); err != nil {
		return wrap.Errorf("invalid http config: %w", err)
	}

	if err := c.PostgresConfig.Validate(); err != nil {
		return wrap.Errorf("invalid postgres config: %w", err)
	}

	return nil
}

func (c Http) Validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
	)

	if err != nil {
		return wrap.Errorf("failed to validate http config: %w", err)
	}

	return nil
}

func (c Postgres) Validate() error {
	err := validation.ValidateStruct(c,
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.User, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.DbName, validation.Required),
		validation.Field(&c.DbDriver, validation.Required),
		validation.Field(&c.SslMode, validation.Required),
	)

	if err != nil {
		return wrap.Errorf("failed to validate database config: %w", err)
	}

	return nil
}

func getEnvAsStr(name string, defaultValue string) string {
	valStr := os.Getenv(name)
	if strings.TrimSpace(valStr) == "" {
		return defaultValue
	}

	return valStr
}
