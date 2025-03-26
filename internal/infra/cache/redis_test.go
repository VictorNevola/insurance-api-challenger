package cache_test

import (
	"main-api/internal/infra/cache"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisStorage(t *testing.T) {
	t.Parallel()

	testRedis := miniredis.RunT(t)
	rdb := redis.NewClient(&redis.Options{
		Addr: testRedis.Addr(),
	})

	redisStorage := cache.NewRedisCacheAdapter(rdb)

	t.Run("Set value in cache and should not return error", func(t *testing.T) {
		err := redisStorage.Set(t.Context(), "test-01", "test", 0)

		assert.NoError(t, err)
	})

	t.Run("Get value from cache and should not return error", func(t *testing.T) {
		_ = redisStorage.Set(t.Context(), "test-02", "fake-value", 10*time.Second)

		value, err := redisStorage.Get(t.Context(), "test-02")

		assert.NoError(t, err)
		assert.Equal(t, "fake-value", value)
	})

	t.Run("Get value from cache and should return error when not found key", func(t *testing.T) {
		_, err := redisStorage.Get(t.Context(), "test-03")

		assert.Error(t, err)
		assert.Equal(t, cache.ErrCacheMiss, err)
	})
}
