package config

import (
	"os"
)

type Config struct {
	DBPath string
}

func Load() *Config {
	dbPath := getEnv("DB_PATH", "")
	if dbPath == "" {
		panic("DB_PATH is required")
	}

	return &Config{
		DBPath: dbPath,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
