package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"time"
)

type RedisClient struct {
	redisClient *redis.Client
}

func NewCache(redisClient *redis.Client) *RedisClient {
	return &RedisClient{redisClient: redisClient}
}

func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	value, err := rc.redisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

//缓存存储
func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expireTime time.Duration) error {

	err := rc.redisClient.Set(ctx, key, value, expireTime).Err()

	if err != nil {
		fmt.Printf("Redis Set Fail: %v", err)
	}
	return err
}
