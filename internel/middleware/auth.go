package middleware

import (
	"entrytask/internel/constant"
	http_service "entrytask/internel/service/http-service"
	"entrytask/internel/web/response"
	"github.com/gin-gonic/gin"
)

// AuthSessionID
// 检查是否存在session_id这个cookie.
// 		1.存在：继续执行handler
//		2.不存在：response error并终止后面的handler
func AuthSessionID(c *gin.Context) {
	sessionId, err := c.Cookie(constant.SESSION_ID)
	if err != nil || sessionId == "" {
		// 没有session_id意味着没有登录
		response.NewResponse(c).ResponseError(constant.UserLoginRequired.GetRetCode())
		c.Abort()
		return
	}
	c.Next()
}

// AuthUserLogin
// 验证用户是否真正登录
func AuthUserLogin(c *gin.Context) {
	sessionId, err := c.Cookie(constant.SESSION_ID)
	if err != nil || sessionId == "" {
		// 没有session_id意味着没有登录
		response.NewResponse(c).ResponseError(constant.UserLoginRequired.GetRetCode())
		c.Abort()
		return
	}
	service := http_service.NewService(c.Request.Context())
	authResponse, err := service.AuthUser(&http_service.UserAuthRequest{SessionId: sessionId})
	if err != nil {
		response.NewResponse(c).ResponseError(constant.SessionError.GetRetCode())
		c.Abort()
		return
	}
	c.Set(constant.USER_ID, authResponse.UserID)
	c.Set(constant.USERNAME, authResponse.Username)
	c.Next()
}
