package api

import (
	"entrytask/internel/constant"
	http_service "entrytask/internel/service/http-service"
	"entrytask/internel/web/response"
	"github.com/gin-gonic/gin"
	"log"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

// Register 用户注册
func (u *User) Register(c *gin.Context) {
	resp := response.NewResponse(c)
	param := http_service.UserRegisterRequest{}
	if err := c.ShouldBind(&param); err != nil {
		resp.ResponseError(constant.InvalidParams.GetRetCode(), "校验失败（请输入8-32个字符的用户名和密码）")
		return
	}

	service := http_service.NewService(c.Request.Context())
	userRegisterResp, err := service.UserRegister(&param)
	if err != nil {
		log.Printf("service UserRegister Failed : %v", err)
		resp.ResponseError(constant.UserRegisterFailed.GetRetCode(), err.Error())
		return
	}
	resp.ResponseOK(userRegisterResp)
}

// Login 用户登录
func (u *User) Login(c *gin.Context) {
	resp := response.NewResponse(c)
	param := http_service.UserLoginRequest{}
	if err := c.ShouldBind(&param); err != nil {
		resp.ResponseError(constant.InvalidParams.GetRetCode(), "校验失败（请输入8-32个字符的用户名和密码）")
		return
	}
	service := http_service.NewService(c.Request.Context())
	userLoginResp, err := service.UserLogin(&param)
	if err != nil {
		log.Printf("service UserLogin Failed : %v", err)
		resp.ResponseError(constant.UserLoginFailed.GetRetCode(), err.Error())
		return
	}
	c.SetCookie(constant.SESSION_ID, userLoginResp.SessionId, 3600, "/", "127.0.0.1", false, true)
	resp.ResponseOK(userLoginResp)
}
