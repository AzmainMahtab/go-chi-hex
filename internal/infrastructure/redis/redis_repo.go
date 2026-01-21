// Package redis
// This is the repository for redis
package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisAdapter(client *redis.Client) *RedisRepo {
	return &RedisRepo{
		client: client,
	}
}

func (r *RedisRepo) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *RedisRepo) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}
