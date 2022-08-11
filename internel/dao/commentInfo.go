package dao

import (
	"context"
	"entrytask/internel/constant"
	"entrytask/internel/model"
	"entrytask/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type CommentBrief struct {
	CommentId uint      `json:"comment_id"`
	FromName  string    `json:"from_name"`
	Content   string    `json:"comment_content"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentDetail struct {
	CommentId      uint                `json:"comment_id"`
	FromName       string              `json:"from_name"`
	CommentContent string              `json:"comment_content"`
	CreatedAt      time.Time           `json:"created_at"`
	ReplyList      []CommentReplyBrief `json:"reply_list,omitempty"`
}

func (d *Dao) GetCommentDetail(productId uint, commentId uint) (*CommentDetail, error) {

	// 1 获取comment info
	c := model.CommentInfo{
		Model: gorm.Model{
			ID: commentId,
		},
	}
	commentInfo, err := c.SelectCommentInfo(d.engine)
	if err != nil {
		return nil, err
	}

	// 2 获取 reply 列表
	r := model.CommentReply{}
	commentReplyList, err := r.SelectCommentReplyList(d.engine, commentId)
	if err != nil {
		return nil, err
	}

	// 3 组装数据
	commentDetail := &CommentDetail{
		CommentId:      commentInfo.ID,
		FromName:       commentInfo.FromName,
		CommentContent: commentInfo.Content,
		CreatedAt:      commentInfo.CreatedAt,
	}
	var commentReplyBriefList []CommentReplyBrief
	for _, reply := range commentReplyList {
		commentReplyBriefList = append(commentReplyBriefList, CommentReplyBrief{
			ReplyId:       reply.ID,
			ReplyFromId:   reply.ReplyFromId,
			ReplyFromName: reply.ReplyToName,
			ReplyContent:  reply.Content,
			CreatedAt:     reply.CreatedAt,
		})
	}
	commentDetail.ReplyList = commentReplyBriefList
	return commentDetail, nil

}

func (d *Dao) CreateCommentInfo(userId uint, username string, productId uint, content string) (uint, error) {
	commentInfo := model.CommentInfo{
		ProductId: productId,
		FromId:    userId,
		FromName:  username,
		Content:   content,
	}
	info, err := commentInfo.CreateCommentInfo(d.engine)
	d.RedisClient.Del(context.Background(), utils.ConvertUintIdToRedisKey(constant.PRODUCT_ID, info.ProductId))
	if err != nil {
		return 0, err
	}
	return info.ID, nil
}
