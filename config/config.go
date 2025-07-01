package config

import (
	"fmt"
	"os"
)

type Config struct {
	TTLAnnotation string
}

func ReadConfig() Config {
	var config Config

	config.TTLAnnotation = os.Getenv("TTL_ANNOTATION")
	if config.TTLAnnotation == "" {
		fmt.Println("TTL_ANNOTATION env variable not set")
		os.Exit(2)
	}
	return config
}
