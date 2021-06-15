package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/mebr0/tiny-url/internal/domain"
	"time"
)

type URLs interface {
	Set(ctx context.Context, url domain.URL) error
	Get(ctx context.Context, alias string) (domain.URL, error)
	Delete(ctx context.Context, alias string) error
}

type Caches struct {
	URLs URLs
}

func NewCaches(client *redis.Client, defaultTTL time.Duration) *Caches {
	return &Caches{
		URLs: newURLsCache(client, defaultTTL),
	}
}
