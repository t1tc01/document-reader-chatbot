package configs

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	JWT           JWTConfig
	Observability ObservabilityConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Name         string
	Version      string
	Environment  string
	Port         string
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	URL          string
	Name         string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey      string
	ExpirationTime time.Duration
	Issuer         string
}

// ObservabilityConfig holds observability configuration
type ObservabilityConfig struct {
	ServiceName    string
	ServiceVersion string
	JaegerEndpoint string
	MetricsPort    string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file is optional, so we don't return error if it doesn't exist
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	config := &Config{
		Server: ServerConfig{
			Name:         getEnv("SERVER_NAME", "document-reader-chatbot"),
			Version:      getEnv("SERVER_VERSION", "1.0.0"),
			Environment:  getEnv("ENVIRONMENT", "development"),
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 30),
			IdleTimeout:  getEnvAsInt("IDLE_TIMEOUT", 60),
		},
		Database: DatabaseConfig{
			URL:          getEnv("DATABASE_URL", "postgres://postgres:password@localhost:5432/document_reader_chatbot?sslmode=disable"),
			Name:         getEnv("DATABASE_NAME", "document_reader_chatbot"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxLifetime:  time.Duration(getEnvAsInt("DB_MAX_LIFETIME", 3600)) * time.Second,
		},
		JWT: JWTConfig{
			SecretKey:      getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production"),
			ExpirationTime: time.Duration(getEnvAsInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
			Issuer:         getEnv("JWT_ISSUER", "document-reader-chatbot"),
		},
		Observability: ObservabilityConfig{
			ServiceName:    getEnv("OTEL_SERVICE_NAME", "document-reader-chatbot"),
			ServiceVersion: getEnv("OTEL_SERVICE_VERSION", "1.0.0"),
			JaegerEndpoint: getEnv("JAEGER_ENDPOINT", "http://localhost:14268/api/traces"),
			MetricsPort:    getEnv("METRICS_PORT", "9090"),
		},
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.URL == "" {
		return fmt.Errorf("database URL is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if c.JWT.SecretKey == "" {
		return fmt.Errorf("JWT secret key is required")
	}
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	return nil
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnvAsInt gets an environment variable as integer with a fallback value
func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}
