package dao

import (
	"context"
	"entrytask/internel/constant"
	"entrytask/internel/model"
	"entrytask/pkg/utils"
	"errors"
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
	ProductId      uint                `json:"product_id"`
	FromName       string              `json:"from_name"`
	CommentContent string              `json:"comment_content"`
	CreatedAt      time.Time           `json:"created_at"`
	ReplyList      []CommentReplyBrief `json:"reply_list,omitempty"`
}

// GetCommentDetail 获取用户的评论详细信息（包含该评论的回复）
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
	if commentInfo.ProductId != productId { //说明当前评论并不属于productId
		return nil, errors.New("")
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
		ProductId:      commentInfo.ProductId,
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

// CreateCommentInfo 创建用户评论
func (d *Dao) CreateCommentInfo(userId int64, username string, productId uint, content string) (uint, error) {
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
