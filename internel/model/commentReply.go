package model

import (
	"entrytask/internel/constant"
	"gorm.io/gorm"
	"log"
)

type CommentReply struct {
	CommentId     uint
	ReplyToId     uint
	ReplyToName   string
	ReplyFromId   uint
	ReplyFromName string
	Content       string
	gorm.Model
}

func (r CommentReply) SelectCommentReplyList(db *gorm.DB, commentId uint) ([]CommentReply, error) {
	var commentReplyList []CommentReply
	err := db.Model(&r).Where("comment_id = ?", commentId).Find(&commentReplyList).Error
	if err != nil {
		log.Println("select comment reply list failed ")
		return nil, err
	}
	return commentReplyList, nil

}

func (r CommentReply) CreateCommentReply(db *gorm.DB) (uint, error) {
	err := db.Create(&r).Error
	if err != nil {
		log.Println("create comment reply failed")
		return 0, err
	}
	return r.ID, nil

}

func (r CommentReply) TableName() string {
	return constant.COMMENT_REPLY_TAB
}
