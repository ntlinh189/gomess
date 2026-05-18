package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RedisInterface interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
}

type Redis struct {
	client *goredis.Client
}

func NewRedis(addr string) *Redis {
	client := goredis.NewClient(&goredis.Options{
		Addr: addr,
	})

	return &Redis{client: client}
}

func (r *Redis) Set(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}