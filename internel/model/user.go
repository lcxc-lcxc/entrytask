package model

import (
	"entrytask/internel/constant"
	"gorm.io/gorm"
)

type User struct {
	Username string
	Password string
	gorm.Model
}

func (u User) CreateUser(db *gorm.DB) (*User, error) {
	err := db.Create(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (u User) SelectUserByUsername(db *gorm.DB) (*User, error) {
	var user User
	err := db.Where("username = ?", u.Username).Select("ID", "username", "password").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u User) TableName() string {
	return constant.USER_TAB
}
