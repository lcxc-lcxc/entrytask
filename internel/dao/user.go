package dao

import (
	"entrytask/internel/model"
	"gorm.io/gorm"
	"time"
)

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
	user := model.User{
		Username: username,
	}
	return user.SelectUserByUsername(d.engine)
}
