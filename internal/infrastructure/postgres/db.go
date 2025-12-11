package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	PoolSize int
}

func ConnectDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure the connection pool
	db.SetMaxOpenConns(cfg.PoolSize)     // Use config setting for Max connections
	db.SetMaxIdleConns(cfg.PoolSize / 2) // Set idle to half of max
	db.SetConnMaxLifetime(5 * time.Minute)

	// Health Check and Retry Loop
	const maxRetries = 5
	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err = db.PingContext(ctx); err == nil {
			cancel()
			log.Println("✅ Successfully connected to PostgreSQL (standard *sql.DB)")
			return db, nil
		}
		cancel()

		log.Printf("⚠️ Failed to ping database (attempt %d/%d): %v. Retrying in 2 seconds...", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d retries", maxRetries)
}
