package config

import (
	"errors"
	"os"
	"strings"
)

const (
	defaultAppEnv             = "development"
	defaultHTTPAddr           = ":8080"
	defaultCORSAllowedOrigins = "http://localhost:5173"
)

type Config struct {
	AppEnv             string
	HTTPAddr           string
	DatabaseURL        string
	CORSAllowedOrigins []string
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:             getEnv("APP_ENV", defaultAppEnv),
		HTTPAddr:           getEnv("HTTP_ADDR", defaultHTTPAddr),
		DatabaseURL:        strings.TrimSpace(os.Getenv("DATABASE_URL")),
		CORSAllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", defaultCORSAllowedOrigins)),
	}

	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			items = append(items, item)
		}
	}
	return items
}
