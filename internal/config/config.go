// Package config provides application configuration loading and constants.
//
// It is responsible for:
//   - Loading configuration from environment variables
//   - Providing compile-time constants for API endpoints, limits, etc.
//   - Central place for all configuration-related values
//
// Usage:
//
//	cfg := config.Load()
//	producer := producer.NewHTTPProducer(..., config.APIBaseURL, ...)
package config

import (
	"fmt"
	"os"

	"github.com/nisemenov/etl_service/internal/validation"
)

const (
	FetchPaymentsPath = "/payments/not-exported/"
	FetchYookassaPath = "/payments/yookassa/not-exported/"
	AckPaymentsPath   = "/payments/mark-exported/"
	AckYookassaPath   = "/payments/yookassa/mark-exported/"
)

type Config struct {
	DBPath     string `validate:"required"`
	APIBaseURL string `validate:"required"`
}

func Load() *Config {
	dbPath := getEnv("DB_PATH", "")
	apiBaseURL := getEnv("API_BASE_URL", "")

	config := &Config{
		DBPath:     dbPath,
		APIBaseURL: apiBaseURL,
	}

	if err := validation.Validate.Struct(config); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
