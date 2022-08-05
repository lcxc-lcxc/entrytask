package model

import (
	"entrytask/internel/constant"
	"gorm.io/gorm"
	"log"
)

type CommentInfo struct {
	ProductId uint
	FromId    uint
	FromName  string
	Content   string
	gorm.Model
}

func (c CommentInfo) SelectCommentInfoList(db *gorm.DB, productId uint) ([]CommentInfo, error) {
	var commentInfoList []CommentInfo
	err := db.Model(&c).Where("product_id = ?", productId).Find(&commentInfoList).Error
	if err != nil {
		log.Println("get comment info list failed ")
		return nil, err
	}
	return commentInfoList, nil

}

func (c CommentInfo) SelectCommentInfo(db *gorm.DB) (*CommentInfo, error) {
	var commentInfo CommentInfo
	err := db.Model(&c).First(&commentInfo).Error
	if err != nil {
		log.Println("get comment info failed")
		return nil, err
	}
	return &commentInfo, nil
}

func (c CommentInfo) CreateCommentInfo(db *gorm.DB) (*CommentInfo, error) {
	err := db.Create(&c).Error
	if err != nil {
		log.Println("insert comment info failed")
		return nil, err
	}
	return &c, nil
}

func (c CommentInfo) TableName() string {
	return constant.COMMENT_INFO_TAB
}
