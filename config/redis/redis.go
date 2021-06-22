package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client interface contract
type Client interface {
	Publish(ctx context.Context, channel string, message interface{}) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) (int64, error)
	Ping(ctx context.Context) (string, error)
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
	Set(ctx context.Context, key string, value interface{}, exp time.Duration) (string, error)
}
