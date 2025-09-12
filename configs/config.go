package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"KYC"`
	Server      Server
	Database    Database
	Debug       bool   `envconfig:"DEBUG" default:"false"`
	Lang        string `default:"fa"`
}

var GlobalConfigInstance *Config

type Server struct {
	Host         string        `envconfig:"SERVER_HOST"`
	Port         int           `envconfig:"SERVER_PORT"`
	WriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
	ReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"10s"`
	Debug        bool          `envconfig:"SERVER_DEBUG" default:"false"`
}

type Database struct {
	Type         string `envconfig:"DATABASE_TYPE" default:"sqlite3"`
	Dns          string `envconfig:"DATABASE_DNS"`
	MaxLifeTime  int    `envconfig:"DATABASE_MAX_LIFETIME"`
	MaxIdleTime  int    `envconfig:"DATABASE_MAX_IDLETIME"`
	MaxIdleConns int    `envconfig:"DATABASE_MAX_IDLECONNS"`
	MaxOpenConns int    `envconfig:"DATABASE_MAX_OPENCONNS"`
}

func Load() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load env variable into config struct: %w", err)
	}
	GlobalConfigInstance = &cfg

	return &cfg, nil
}
