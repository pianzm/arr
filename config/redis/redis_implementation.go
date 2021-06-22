package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pianzm/arr/config"
)

// Conn base struct
type Conn struct {
	Client *redis.Client
}

// ConnectRedis init redis
func ConnectRedis(cfg *config.Config) (Client, error) {
	cl, err := setupRedis(cfg.RedisHost, cfg.RedisTLS, cfg.RedisPassword, cfg.RedisPort, cfg.RedisDB)
	if err != nil {
		return nil, err
	}
	return Conn{
		Client: cl,
	}, nil
}

func (r Conn) Publish(ctx context.Context, channel string, message interface{}) error {
	return r.Client.Publish(ctx, channel, message).Err()
}

func (r Conn) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}
func (r Conn) Del(ctx context.Context, key string) (int64, error) {
	return 0, nil
}

func (r Conn) Ping(ctx context.Context) (string, error) {
	return r.Client.Ping(ctx).Result()
}

func (r Conn) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return r.Client.Subscribe(ctx, channels...)
}

// Set redis
func (r Conn) Set(ctx context.Context, key string, value interface{}, exp time.Duration) (string, error) {
	return r.Client.Set(ctx, key, value, exp).Result()
}

// SetupRedis
func setupRedis(redisHost, redisTLS, redisPassword, redisPort, redisDB string) (*redis.Client, error) {
	tlsSecured, err := strconv.ParseBool(redisTLS)
	if err != nil {
		return nil, err
	}

	var conf *tls.Config

	if tlsSecured {
		conf = &tls.Config{
			InsecureSkipVerify: tlsSecured,
		}
	}

	useDB, _ := strconv.Atoi(redisDB)
	cl := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%v:%v", redisHost, redisPort),
		Password:  redisPassword,
		DB:        useDB,
		TLSConfig: conf,
	})

	return cl, nil
}
