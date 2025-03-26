package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	RedisCacheAdapter struct {
		redisClient *redis.Client
	}
)

func NewRedisCacheAdapter(redisClient *redis.Client) *RedisCacheAdapter {
	return &RedisCacheAdapter{
		redisClient: redisClient,
	}
}

func (r *RedisCacheAdapter) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	_, err := r.redisClient.Set(ctx, key, value, ttl).Result()

	return err
}

// func (r *RedisCacheAdapter) GetAsStruct(ctx context.Context, key string, dest interface{}) error {
// 	val, err := r.redisClient.Get(ctx, key).Result()
// 	if err == redis.Nil {
// 		return ErrCacheMiss
// 	}

// 	if err != nil {
// 		return err
// 	}

// 	err = json.Unmarshal([]byte(val), dest)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (r *RedisCacheAdapter) Get(ctx context.Context, key string) (string, error) {
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}

	if err != nil {
		return "", err
	}

	return val, nil
}
