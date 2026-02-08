package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	Port        string
	DatabaseURL string
	CORSOrigin  string
	AdminSecret string
	Environment string
}

// Load reads configuration from environment variables (and .env file if present).
func Load() *Config {
	// Load .env file if it exists; ignore error in production where env vars are set externally.
	_ = godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:3000"),
		AdminSecret: getEnv("ADMIN_SECRET", ""),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
