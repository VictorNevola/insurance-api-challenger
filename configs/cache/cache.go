package cache

import (
	"context"
	"log"
	"main-api/configs/envs"

	"github.com/redis/go-redis/v9"
)

func InitRedisStorage() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     envs.AppConfig.RedisURL,
		DB:       0,
		Password: "",
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}

	return client
}
