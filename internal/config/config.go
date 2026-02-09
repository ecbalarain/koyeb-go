package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	Port        string
	DatabaseURL string
	CORSOrigin  string
	AdminSecret string
	JWTSecret   string
	JWTExpiry   int // JWT token expiry in hours
	Environment string
	EmailHost   string
	EmailPort   int
	EmailUser   string
	EmailPass   string
	EmailFrom   string
}

// Load reads configuration from environment variables (and .env file if present).
func Load() *Config {
	// Load .env file if it exists; ignore error in production where env vars are set externally.
	_ = godotenv.Load()

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if jwtExpiry <= 0 {
		jwtExpiry = 24
	}

	emailPort, _ := strconv.Atoi(getEnv("EMAIL_PORT", "465"))
	if emailPort <= 0 {
		emailPort = 465
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:3000"),
		AdminSecret: getEnv("ADMIN_SECRET", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		JWTExpiry:   jwtExpiry,
		Environment: getEnv("ENVIRONMENT", "development"),
		EmailHost:   getEnv("EMAIL_HOST", "smtp-relay.brevo.com"),
		EmailPort:   emailPort,
		EmailUser:   getEnv("EMAIL_HOST_USER", ""),
		EmailPass:   getEnv("EMAIL_HOST_PASSWORD", ""),
		EmailFrom:   getEnv("DEFAULT_FROM_EMAIL", "Bhomanshah <store@bhomanshah.com>"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
