package config

import (
	"log/slog"
	"os"

	"github.com/codingconcepts/env"
)

type Config struct {
	TTLAnnotation        string `env:"TTL_ANNOTATION"`
	CheckIntervalSeconds int    `env:"CHECK_INTERVAL_SECONDS"`
	HTTPListenAddress    string `env:"HTTP_LISTEN_ADDRESS"`
	HTTPListenPort       string `env:"HTTP_LISTEN_PORT"`
}

func ReadConfig() Config {
	var config Config

	if err := env.Set(&config); err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}

	return config
}
