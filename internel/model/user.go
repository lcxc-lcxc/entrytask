package model

import (
	"entrytask/global"
	"entrytask/internel/constant"
	"errors"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type User struct {
	UserId    int64 `gorm:"primarykey"`
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// CreateUser
func (u User) CreateUser(db *gorm.DB) (*User, error) {
	// 0 获取雪花id ，并找到序号
	userId := global.SnowFlakeNode1.Generate().Int64()
	tabSeq := userId % 100
	// 1 开启事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 2 user_index插入数据
		userIndex := UserIndex{
			Username: u.Username,
			UserId:   userId,
		}
		err := db.Create(&userIndex).Error
		if err != nil {
			log.Println("Insert UserIndex tab failed ")
			return err
		}

		// 3 往user_tab_xxx插入数据，返回User
		err = db.Exec("INSERT INTO "+constant.USER_TAB+"_"+strconv.FormatInt(tabSeq, 10)+
			"(`user_id`,`username`,`password`,`created_at`,`updated_at`) values (?,?,?,?,?)",
			userId,
			u.Username,
			u.Password,
			time.Now(),
			time.Now(),
		).Error
		if err != nil {
			log.Printf("Insert user_tab_%d tab failed", tabSeq)
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	u.UserId = userId

	return &u, nil
}

func (u User) SelectUserByUsername(db *gorm.DB) (*User, error) {
	// 1 先用username获取userid
	userIndex := UserIndex{
		Username: u.Username,
	}
	err := db.First(&userIndex, "username = ?", userIndex.Username).Error
	if err != nil {
		return nil, err
	}

	// 2 获得表序号
	tabSeq := userIndex.UserId % 100

	// 3 根据userId  获得真正的user
	var user User
	err = db.Raw("SELECT * FROM "+constant.USER_TAB+"_"+strconv.FormatInt(tabSeq, 10)+
		" where user_id = ?", userIndex.UserId,
	).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	if user.UserId == 0 {
		return nil, errors.New("Get User from database Failed")
	}

	return &user, nil
}

func (u User) TableName() string {
	return constant.USER_TAB
}

// UserIndex 用于映射username -> userId
type UserIndex struct {
	Username string `gorm:"primarykey"`
	UserId   int64
}

func (ui UserIndex) TableName() string {
	return constant.USER_INDEX_TAB
}
