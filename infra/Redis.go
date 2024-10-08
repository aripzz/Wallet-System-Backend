package infra

import (
	"Wallet-System-Backend/utils"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {

	utils.LoadEnv()
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Get(key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

func (r *RedisClient) Set(key string, value string) error {
	return r.client.Set(ctx, key, value, 5*time.Minute).Err()
}

func (r *RedisClient) Delete(key string) error {
	return r.client.Del(ctx, key).Err()
}
