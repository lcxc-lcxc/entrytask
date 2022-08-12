package dao

import (
	"entrytask/internel/model"
	"time"
)

type CommentReplyBrief struct {
	ReplyId       uint      `json:"reply_id"`
	ReplyFromId   uint      `json:"reply_from_id"`
	ReplyFromName string    `json:"reply_from_name"`
	ReplyContent  string    `json:"reply_content"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateCommentReply 创建用户回复
func (d *Dao) CreateCommentReply(userId uint, username string, replyToId uint, replyToName string, productId uint, commentId uint, replyContent string) (uint, error) {
	commentReply := model.CommentReply{
		ReplyFromId:   userId,
		ReplyFromName: username,
		ReplyToId:     replyToId,
		ReplyToName:   replyToName,
		CommentId:     commentId,
		Content:       replyContent,
	}
	replyId, err := commentReply.CreateCommentReply(d.engine)
	if err != nil {
		return 0, err
	}
	return replyId, nil

}
