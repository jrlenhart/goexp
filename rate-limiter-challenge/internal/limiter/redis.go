package limiter

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Storage interface {
	Increment(ctx context.Context, key string) (int, error)
	Expire(ctx context.Context, key string, duration time.Duration) error
	BlockKey(ctx context.Context, key string, duration time.Duration) error
	IsBlocked(ctx context.Context, key string) (bool, error)
}

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(host string, port string) *RedisStorage {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	return &RedisStorage{client: client}
}

func (rs *RedisStorage) Increment(ctx context.Context, key string) (int, error) {
	count, err := rs.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (rs *RedisStorage) Expire(ctx context.Context, key string, duration time.Duration) error {
	_, err := rs.client.Expire(ctx, key, duration).Result()
	return err
}

func (rs *RedisStorage) BlockKey(ctx context.Context, key string, duration time.Duration) error {
	_, err := rs.client.Set(ctx, "blocked:"+key, 1, duration).Result()
	return err
}

func (rs *RedisStorage) IsBlocked(ctx context.Context, key string) (bool, error) {
	_, err := rs.client.Get(ctx, "blocked:"+key).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, err
}
