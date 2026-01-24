package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HTTPAddress string
	DSN         string
}

func Load() (Config, error) {
	httpAddr := getEnv("HTTP_ADDR", ":8080")
	postgresdsn := strings.TrimSpace(os.Getenv("POSTGRES_DSN"))
	if postgresdsn == "" {
		return Config{}, fmt.Errorf("environment variable POSTGRES_DSN is not set")
	}
	return Config{
		HTTPAddress: httpAddr,
		DSN:         postgresdsn,
	}, nil
}

func getEnv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
