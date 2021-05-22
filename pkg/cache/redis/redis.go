package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

const timeout = 10 * time.Second

// NewClient established connection to a redis instance using provided URI and auth credentials
func NewClient(uri, password string, db int) (*redis.Client, error) {
	opts := &redis.Options{
		Addr:     uri,
		Password: password,
		DB:       db,
	}

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}
