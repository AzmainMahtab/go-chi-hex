// Package redis
// redis connection and new instance fuc here
package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func NewRedisClient(cfg RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancle := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancle()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, fmt.Errorf("Redis connection failed: %w", err)
	}

	log.Print("Redis conneted !")

	return client, nil
}
