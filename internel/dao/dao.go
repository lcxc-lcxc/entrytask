package dao

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Dao struct {
	engine      *gorm.DB
	RedisClient *redis.Client
}

func NewDAO(engine *gorm.DB, client *redis.Client) *Dao {
	return &Dao{
		engine:      engine,
		RedisClient: client,
	}
}
