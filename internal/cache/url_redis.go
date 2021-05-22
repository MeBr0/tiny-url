package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/mebr0/tiny-url/internal/domain"
	"time"
)

type URLsCache struct {
	client     *redis.Client
	defaultTTL time.Duration
}

func newURLsCache(client *redis.Client, defaultTTL time.Duration) *URLsCache {
	return &URLsCache{
		client:     client,
		defaultTTL: defaultTTL,
	}
}

func (c *URLsCache) Set(ctx context.Context, url domain.URL) error {
	return c.client.Set(ctx, url.Alias, url, c.defaultTTL).Err()
}

func (c *URLsCache) Get(ctx context.Context, alias string) (domain.URL, error) {
	var url domain.URL

	if err := c.client.Get(ctx, alias).Scan(&url); err != nil {
		return url, err
	}

	return url, nil
}
