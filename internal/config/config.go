// Package config
// config package is used to manage and import all the configs
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

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

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	RedisDB  int
}

type JWTConfig struct {
	PrivateKeypath string
	PublicKeyPath  string
	AccessTTL      time.Duration
	RefreshTTL     time.Duration
	Issuer         string
}

type NATSConfig struct {
	URL string
}

type Config struct {
	Server ServerConfig
	DB     DatabaseConfig
	JWT    JWTConfig
	Redis  RedisConfig
	NATS   NATSConfig
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
			Host:     getEnv("DB_HOST", "localhost"), //RUNNING LOCAL DEV REMOVE _DB BEFORE DOCKER BUILD
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			DBName:   getEnv("DB_NAME", "docpad_db"),
			Password: os.Getenv("DB_PASSWORD"),
		},

		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", "hehePassRedis$t0nk"),
		},

		JWT: JWTConfig{
			PrivateKeypath: getEnv("AUTH_PRIVATE_KEY_PATH", "./certs/private.pem"),
			PublicKeyPath:  getEnv("AUTH_PUBLIC_KEY_PATH", "./certs/public.pem"),
			Issuer:         getEnv("AUTH_ISSUER", "appName-api"),
		},

		NATS: NATSConfig{
			URL: getEnv("NATS_URL", "nats://localhost:4222"),
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

	// REDIS DB ENV

	redisDBStr := getEnv("REDIS_DB", "o")
	redisDb, err := strconv.Atoi(redisDBStr)
	if err != nil {
		return nil, fmt.Errorf("Invaild redis DB value: %w", err)
	}
	cfg.Redis.RedisDB = redisDb

	// Hnadle TTL for access and RefreshTTL
	accessTTL, err := time.ParseDuration(getEnv("AUTH_ACCESS_TTL", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid AUTH_ACCESS_TTL: %w", err)
	}
	cfg.JWT.AccessTTL = accessTTL

	refreshTTL, err := time.ParseDuration(getEnv("AUTH_REFRESH_TTL", "168h")) // Default 7 days
	if err != nil {
		return nil, fmt.Errorf("invalid AUTH_REFRESH_TTL: %w", err)
	}
	cfg.JWT.RefreshTTL = refreshTTL

	return cfg, nil
}
