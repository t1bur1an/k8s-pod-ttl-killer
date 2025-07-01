package config

import (
	"log/slog"
	"os"
	"strconv"
)

type Config struct {
	TTLAnnotation        string
	CheckIntervalSeconds int
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
	return config
}
