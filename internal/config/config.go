package config

import (
	"time"

	"github.com/caarlos0/env/v10"
)

var envPrefix = "TMPL_SRV_"

type (
	Config struct {
		Logger Logger `env:"LOG_"`
		App    App    `env:"APP_"`
	}

	Logger struct {
		Level string `env:"LEVEL" envDefault:"info"`
	}

	App struct {
		HTTP HTTPServer `env:"HTTP_"`
	}

	HTTPServer struct {
		Port            int           `env:"PORT" envDefault:"8080"`
		ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT" envDefault:"10s"`
		ReadTimeout     time.Duration `env:"READ_TIMEOUT" envDefault:"30s"`
	}
)

func New() (*Config, error) {
	opts := env.Options{
		Prefix:                envPrefix,
		UseFieldNameByDefault: true,
	}

	var cfg Config
	if err := env.ParseWithOptions(&cfg, opts); err != nil {
		return nil, err
	}

	return &cfg, nil
}
