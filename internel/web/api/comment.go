package api

import (
	"entrytask/internel/constant"
	http_service "entrytask/internel/service/http-service"
	"entrytask/internel/web/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type Comment struct {
}

func NewComment() *Comment {
	return &Comment{}
}

func (cm *Comment) Detail(c *gin.Context) {
	resp := response.NewResponse(c)

	productId, commentId := cast.ToUint(c.Param(constant.PRODUCT_ID)), cast.ToUint(c.Param(constant.COMMENT_ID))
	if productId == 0 || commentId == 0 {
		resp.ResponseError(constant.InvalidParams.GetRetCode())
		return
	}
	param := &http_service.CommentDetailRequest{
		ProductId: productId,
		CommentId: commentId,
	}
	service := http_service.NewService(c.Request.Context())
	commentDetailResponse, err := service.CommentDetail(param)
	if err != nil {
		resp.ResponseError(constant.CommentDetailGetFailed.GetRetCode())
		return
	}
	resp.ResponseOK(commentDetailResponse)

}

func (cm *Comment) Create(c *gin.Context) {
	resp := response.NewResponse(c)

	//1 获取当前用户id和username
	usernameAny, usernameExists := c.Get(constant.USERNAME)
	userIdAny, userIdExists := c.Get(constant.USER_ID)
	if !usernameExists || !userIdExists {
		resp.ResponseError(constant.UserLoginRequired.GetRetCode())
		return
	}
	userId := cast.ToUint(userIdAny)
	username := cast.ToString(usernameAny)
	if userId == 0 || username == "" {
		resp.ResponseError(constant.UserLoginRequired.GetRetCode())
		return
	}

	// 2 获取产品id
	productId := cast.ToUint(c.Param(constant.PRODUCT_ID))
	if productId == 0 {
		resp.ResponseError(constant.InvalidParams.GetRetCode())
		return
	}

	// 3 获取评论内容
	param := http_service.CommentCreateRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		resp.ResponseError(constant.InvalidParams.GetRetCode())
		return
	}
	param.UserId = userId
	param.Username = username
	param.ProductId = productId

	// 4 调用service
	service := http_service.NewService(c.Request.Context())
	commentCreateResponse, err := service.CommentCreate(&param)
	if err != nil {
		resp.ResponseError(constant.CommentCreateFailed.GetRetCode())
		return
	}
	resp.ResponseOK(commentCreateResponse)

}

func (cm *Comment) Reply(c *gin.Context) {
	resp := response.NewResponse(c)

	//1 获取当前用户id和username
	usernameAny, usernameExists := c.Get(constant.USERNAME)
	userIdAny, userIdExists := c.Get(constant.USER_ID)
	if !usernameExists || !userIdExists {
		resp.ResponseError(constant.UserLoginRequired.GetRetCode())
		return
	}
	userId := cast.ToUint(userIdAny)
	username := cast.ToString(usernameAny)
	if userId == 0 || username == "" {
		resp.ResponseError(constant.UserLoginRequired.GetRetCode())
		return
	}

	// 2 获取产品id和主评论id
	productId := cast.ToUint(c.Param(constant.PRODUCT_ID))
	commentId := cast.ToUint(c.Param(constant.COMMENT_ID))
	if productId == 0 || commentId == 0 {
		resp.ResponseError(constant.InvalidParams.GetRetCode())
		return
	}

	// 3 获取回复对象 和 内容
	param := http_service.CommentReplyRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		resp.ResponseError(constant.InvalidParams.GetRetCode())
		return
	}

	param.UserId = userId
	param.Username = username
	param.ProductId = productId
	param.CommentId = commentId

	// 4 调用service
	service := http_service.NewService(c.Request.Context())
	commentReplyResponse, err := service.CommentReply(&param)

	if err != nil {
		resp.ResponseError(constant.CommentReplyCreateFailed.GetRetCode())
		return
	}
	resp.ResponseOK(commentReplyResponse)

}
