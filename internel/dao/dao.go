package dao

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Dao struct {
	engine      *gorm.DB
	RedisClient *redis.Client
}

// NewDAO 返回对mysql和redis客户端对封装，便于后续调用
func NewDAO(engine *gorm.DB, client *redis.Client) *Dao {
	return &Dao{
		engine:      engine,
		RedisClient: client,
	}
}
