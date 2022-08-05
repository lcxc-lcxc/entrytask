package http_service

import (
	"entrytask/internel/dao"
)

type CommentDetailRequest struct {
	ProductId uint
	CommentId uint
}

type CommentDetailResponse struct {
	CommentDetail *dao.CommentDetail `json:"comment"`
}

type CommentCreateRequest struct {
	UserId    uint   `binding:"-"`
	Username  string `binding:"-"`
	ProductId uint   `binding:"-"`
	Content   string `json:"comment_content" form:"comment_content" binding:"required,max=512"`
}

type CommentCreateResponse struct {
	CommentId uint `json:"comment_id"`
}

type CommentReplyRequest struct {
	UserId       uint   `binding:"-"`
	Username     string `binding:"-"`
	ProductId    uint   `binding:"-"`
	CommentId    uint   `binding:"-"`
	ReplyToId    uint   `json:"reply_to_id" form:"reply_to_id" binding:"required"`
	ReplyToName  string `json:"reply_to_name" form:"reply_to_name" binding:"required,max=32"`
	ReplyContent string `json:"reply_content" form:"reply_content" binding:"required,max=512"`
}

type CommentReplyResponse struct {
	CommentReplyId uint `json:"comment_reply_id"`
}

func (svc *Service) CommentDetail(request *CommentDetailRequest) (*CommentDetailResponse, error) {
	commentDetail, err := svc.dao.GetCommentDetail(request.ProductId, request.CommentId)
	if err != nil {
		return nil, err
	}
	return &CommentDetailResponse{CommentDetail: commentDetail}, nil
}

func (svc *Service) CommentCreate(request *CommentCreateRequest) (*CommentCreateResponse, error) {
	commentInfoId, err := svc.dao.CreateCommentInfo(request.UserId, request.Username, request.ProductId, request.Content)
	if err != nil {
		return nil, err
	}
	return &CommentCreateResponse{CommentId: commentInfoId}, nil
}

func (svc *Service) CommentReply(request *CommentReplyRequest) (*CommentReplyResponse, error) {

	replyId, err := svc.dao.CreateCommentReply(request.UserId, request.Username, request.ReplyToId, request.ReplyToName, request.ProductId, request.CommentId, request.ReplyContent)
	if err != nil {
		return nil, err
	}
	return &CommentReplyResponse{CommentReplyId: replyId}, nil

}
