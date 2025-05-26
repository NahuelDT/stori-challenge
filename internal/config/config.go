package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/NahuelDT/stori-challenge/internal/infrastructure/database"
	"github.com/NahuelDT/stori-challenge/internal/infrastructure/email"
)

type Config struct {
	Server   ServerConfig
	Email    email.SMTPConfig
	Database database.PostgresConfig
	File     FileConfig
}

type ServerConfig struct {
	Port        string
	Environment string
	LogLevel    string
}

type FileConfig struct {
	WatchDirectory string
	ProcessedDir   string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:        getEnvOrDefault("PORT", "8080"),
			Environment: getEnvOrDefault("ENVIRONMENT", "development"),
			LogLevel:    getEnvOrDefault("LOG_LEVEL", "info"),
		},
		Email: email.SMTPConfig{
			Host:     getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnvOrDefault("SMTP_PORT", "587"),
			Username: os.Getenv("SMTP_USERNAME"),
			Password: os.Getenv("SMTP_PASSWORD"),
			From:     getEnvOrDefault("SMTP_FROM", "noreply@stori.com"),
		},
		Database: database.PostgresConfig{
			Host:     getEnvOrDefault("DB_HOST", "localhost"),
			Port:     getEnvOrDefault("DB_PORT", "5432"),
			User:     getEnvOrDefault("DB_USER", "postgres"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   getEnvOrDefault("DB_NAME", "stori_challenge"),
			SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		},
		File: FileConfig{
			WatchDirectory: getEnvOrDefault("WATCH_DIRECTORY", "/data"),
			ProcessedDir:   getEnvOrDefault("PROCESSED_DIRECTORY", "/data/processed"),
		},
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.Server.Environment) == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return strings.ToLower(c.Server.Environment) == "production"
}

// DatabaseEnabled returns true if database configuration is complete
func (c *Config) DatabaseEnabled() bool {
	return c.Database.Host != "" &&
		c.Database.User != "" &&
		c.Database.Password != "" &&
		c.Database.DBName != ""
}

func (c *Config) validate() error {
	var errors []string

	if c.Email.Username == "" {
		errors = append(errors, "SMTP_USERNAME is required")
	}
	if c.Email.Password == "" {
		errors = append(errors, "SMTP_PASSWORD is required")
	}

	if c.File.WatchDirectory == "" {
		errors = append(errors, "WATCH_DIRECTORY is required")
	}

	if port, err := strconv.Atoi(c.Server.Port); err != nil || port <= 0 || port > 65535 {
		errors = append(errors, "PORT must be a valid port number (1-65535)")
	}

	validLogLevels := []string{"debug", "info", "warn", "error", "fatal"}
	if !contains(validLogLevels, strings.ToLower(c.Server.LogLevel)) {
		errors = append(errors, fmt.Sprintf("LOG_LEVEL must be one of: %s", strings.Join(validLogLevels, ", ")))
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
