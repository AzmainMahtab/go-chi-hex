// Package config
// config package is used to manage and import all the configs
package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	PoolSize int
}

type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func LoadConfig() (*Config, error) {
	if os.Getenv("GO_ENV") == "" || os.Getenv("GO_ENV") == "development" {
		_ = godotenv.Load()
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("APP_PORT", "8080"),
		},

		DB: DatabaseConfig{
			Host:     getEnv("_DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			DBName:   getEnv("DB_NAME", "docpad_db"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}
	// Validate required configuration (Password is critical)
	if cfg.DB.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD must be set in the environment or .env file")
	}

	// Load and parse PoolSize with validation and fallback
	poolSizeStr := getEnv("DB_POOL_SIZE", "25") // Fallback: DB_POOL_SIZE defaults to 25
	poolSize, err := strconv.Atoi(poolSizeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_POOL_SIZE value: %w", err)
	}
	cfg.DB.PoolSize = poolSize

	return cfg, nil
}
