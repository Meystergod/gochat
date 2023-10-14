package config

import (
	"github.com/kelseyhightower/envconfig"
	"time"
)

type Config struct {
	Log struct {
		LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
		LogFile  string `envconfig:"LOG_FILE"`
		LogSize  int    `envconfig:"LOG_SIZE" default:"10"`
		LogAge   int    `envconfig:"LOG_AGE" default:"28"`
	}

	HTTPServer struct {
		Host         string        `envconfig:"HTTP_HOST" default:"0.0.0.0"`
		Port         string        `envconfig:"HTTP_PORT" default:"8000"`
		WriteTimeout time.Duration `envconfig:"HTTP_WRITE_TIMEOUT" default:"0"`
		ReadTimeout  time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"0"`
	}

	Database struct {
		Host     string `envconfig:"DB_HOST" default:"localhost"`
		Port     string `envconfig:"DB_PORT" default:"27017"`
		Username string `envconfig:"DB_USERNAME"`
		Password string `envconfig:"DB_PASSWORD"`
		Auth     string `envconfig:"DB_AUTH"`
		Name     string `envconfig:"DB_NAME" default:"gochat"`
	}

	Application struct {
		Name    string `envconfig:"APP_NAME" default:"gochat"`
		Version string `envconfig:"APP_VERSION" default:"v0.0.1"`
	}
}

func GetConfig() *Config {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil
	}

	return &cfg
}
