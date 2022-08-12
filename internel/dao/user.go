package dao

import (
	"context"
	"entrytask/internel/constant"
	"entrytask/internel/dao/cache"
	"entrytask/internel/model"
	"entrytask/pkg/utils"
	"gorm.io/gorm"
	"log"
	"time"
)

// CreateUser 创建用户
func (d *Dao) CreateUser(username, password string) (*model.User, error) {
	user := model.User{
		Username: username,
		Password: password,
		Model: gorm.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return user.CreateUser(d.engine)
}

// SelectUserByUsername
// return (nil , err) if the record not found
func (d *Dao) GetUserByName(username string) (*model.User, error) {

	return d.GetUserByNameCache(username)
}

// GetUserByNameCache 封装获取User的缓存模块
func (d *Dao) GetUserByNameCache(username string) (*model.User, error) {

	loadFunction := func(ctx context.Context, key any) (*model.User, error) {
		log.Println("get username cache failed , getting data from database ")
		redisKey, _ := key.(string)
		username := utils.ConvertRedisKeyToString(redisKey)

		user := model.User{
			Username: username,
		}
		userRes, err := user.SelectUserByUsername(d.engine)
		if err != nil {
			return nil, err
		}
		return userRes, nil
	}

	loadableCache := cache.NewLoadableCache[*model.User](loadFunction, d.RedisClient, time.Hour)
	product, err := loadableCache.Get(context.Background(), utils.ConvertStringToRedisKey(constant.USERNAME, username))

	if err != nil {
		return nil, err
	}
	return product, nil

}
