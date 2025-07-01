package config

import (
	"log/slog"
	"os"
)

type Config struct {
	TTLAnnotation        string
	CheckIntervalSeconds string
}

func ReadConfig() Config {
	var config Config

	config.TTLAnnotation = os.Getenv("TTL_ANNOTATION")
	if config.TTLAnnotation == "" {
		slog.Error("TTL_ANNOTATION env variable not set")
		os.Exit(2)
	}

	config.CheckIntervalSeconds = os.Getenv("CHECK_INTERVAL_SECONDS")
	if config.CheckIntervalSeconds == "" {
		slog.Error("CHECK_INTERVAL_SECONDS env variable not set")
		os.Exit(2)
	}
	return config
}
