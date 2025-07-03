package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	TTLAnnotation        string
	CheckIntervalSeconds int
	HTTPListenAddress string
	HTTPListenPort       string
}

func ReadConfig() Config {
	var config Config

	config.TTLAnnotation = os.Getenv("TTL_ANNOTATION")
	if config.TTLAnnotation == "" {
		slog.Error("TTL_ANNOTATION env variable not set")
		os.Exit(2)
	}

	checkIntervalSecondsInt, err := strconv.Atoi(os.Getenv("CHECK_INTERVAL_SECONDS"))
	if err != nil {
		slog.Error("CHECK_INTERVAL_SECONDS env variable not set", "error", err.Error())
		os.Exit(2)
	}
	config.CheckIntervalSeconds = checkIntervalSecondsInt

	config.HTTPListenPort = os.Getenv("HTTP_LISTEN_PORT")
	if config.HTTPListenPort == "" {
		slog.Error("HTTP_LISTEN_PORT env variable not set")
		os.Exit(2)
	}

	config.HTTPListenAddress = os.Getenv("HTTP_LISTEN_ADDRESS")
	if config.HTTPListenAddress == "" {
		slog.Error("HTTP_LISTEN_ADDRESS env variable not set")
		os.Exit(2)
	}

	return config
}
