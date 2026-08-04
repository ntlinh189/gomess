package redis

import (
	"context"
	"errors"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type RedisInterface interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	GetDel(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
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

func (r *Redis) GetDel(ctx context.Context, key string) (string, error) {
	val, err := r.client.GetDel(ctx, key).Result()
	if errors.Is(err, goredis.Nil) {
		return "", ErrNotFound
	}
	return val, err
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *Redis) Incr(ctx context.Context, key string) (int64, error) {
	return r.client.Incr(ctx, key).Result()
}

func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}