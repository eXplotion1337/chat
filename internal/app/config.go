package app

import (
	"chat/internal/infrastructure/postgres"
	"chat/pkg/wrap"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultServiceName          = "chat"
	DefaultDBPoolSizeMultiplier = 10
	DefaultDBMaxOpenConnections = 10
	DefaultDBMaxIdleConnections = 10
	DefaultDBReadTimeout        = 5 * time.Second
	DefaultDBWriteTimeout       = 5 * time.Second
)

var ROOT_PATH string
var WEB_PATH string

type Config struct {
	HttpConfing    *Http
	PostgresConfig *postgres.Config
}

type Http struct {
	Port string
	Host string
}

func InitConfig() (*Config, error) {
	config := Config{
		HttpConfing: &Http{
			Port: getEnvAsStr("HTTP_PORT", ":8080"),
			Host: getEnvAsStr("HTTP_HOST", "localhost"),
		},
		PostgresConfig: &postgres.Config{
			Dsn:                getEnvAsStr("POSTGRES_DSN", ""),
			DBPoolSize:         DefaultDBPoolSizeMultiplier * runtime.NumCPU(),
			ServiceName:        DefaultServiceName,
			MaxOpenConnections: getEnvAsInt("POSTGRES_MAX_OPEN_CONNECTIONS", DefaultDBMaxOpenConnections),
			MaxIdleConnections: getEnvAsInt("POSTGRES_MAX_IDLE_CONNECTIONS", DefaultDBMaxIdleConnections),
			ReadTimeout:        getEnvAsDuration("POSTGRES_READ_TIMEOUT", DefaultDBReadTimeout),
			WriteTimeout:       getEnvAsDuration("POSTGRES_WRITE_TIMEOUT", DefaultDBWriteTimeout),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	ROOT_PATH = filepath.Dir(".")
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
	err := validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
	)

	if err != nil {
		return wrap.Errorf("failed to validate http config: %w", err)
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

func getEnvAsInt(name string, defaultValue int) int {
	envVal := os.Getenv(name)

	result, err := strconv.Atoi(envVal)
	if err != nil {
		return defaultValue
	}

	return result
}

func getEnvAsDuration(name string, defaultValue time.Duration) time.Duration {
	envVal := os.Getenv(name)

	result, err := time.ParseDuration(envVal)
	if err != nil {
		return defaultValue
	}

	return result
}
