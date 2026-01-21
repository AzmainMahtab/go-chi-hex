// Package ports
// This one contains the ports for cache
package ports

import (
	"context"
	"time"
)

type CacheRepo interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
}
