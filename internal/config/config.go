package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv     string
	AppPort    string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	DBMaxConns string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "devpulse"),
		DBUser:     getEnv("DB_USER", "devpulse"),
		DBPassword: getEnv("DB_PASSWORD", "devpulse"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		DBMaxConns: getEnv("DB_MAX_CONNS", "10"),
	}

	if cfg.AppPort == "" {
		return nil, fmt.Errorf("APP_PORT is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "nil" {
		return fallback
	}

	return value
}
