package config

import (
	"log/slog"
	"os"

	"github.com/codingconcepts/env"
)

type Config struct {
	TTLAnnotation        string `env:"TTL_ANNOTATION" default:"k8s-pod-ttl-killer"`
	CheckIntervalSeconds int    `env:"CHECK_INTERVAL_SECONDS" default:"60"`
	HTTPListenAddress    string `env:"HTTP_LISTEN_ADDRESS" default:"0.0.0.0"`
	HTTPListenPort       string `env:"HTTP_LISTEN_PORT" default:"8080"`
}

func ReadConfig() Config {
	var config Config

	if err := env.Set(&config); err != nil {
		slog.Error(err.Error())
		os.Exit(2)
	}

	return config
}
