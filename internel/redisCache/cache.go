package redisCache

import "github.com/go-redis/redis/v8"

import "github.com/eko/gocache/v3/store"

type RedisClient struct {
	RedisStore *store.RedisStore
}

func NewCache(redisClient *redis.Client) *RedisClient {
	return &RedisClient{RedisStore: store.NewRedis(redisClient)}
}
