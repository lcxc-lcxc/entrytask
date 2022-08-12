package response

import (
	"entrytask/internel/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
    Response
	统一的响应结果工具
*/
type Response struct {
	Ctx *gin.Context
}

func NewResponse(ctx *gin.Context) *Response {
	return &Response{Ctx: ctx}
}

func (r *Response) ResponseOK(data interface{}) {
	response := gin.H{"retcode": constant.Success.GetRetCode(), "msg": constant.Success.GetMsg(), "data": data}
	r.Ctx.JSON(http.StatusOK, response)
}

// ResponseError 响应错误时，可带上具体的错误信息。也可不带，直接使用""
func (r *Response) ResponseError(retcode int, data string) {

	r.Ctx.JSON(GetStatusCode(retcode), gin.H{
		"retcode": retcode,
		"msg":     constant.RetcodeMap[retcode].GetMsg(),
		"data":    data,
	})
}

// GetStatusCode 根据不同的retcode来决定不同的http状态码
func GetStatusCode(retcode int) int {
	switch {
	case retcode == constant.ServerError.GetRetCode():
		return http.StatusInternalServerError
	case retcode == constant.InvalidParams.GetRetCode():
		return http.StatusBadRequest
	case retcode == constant.NotFound.GetRetCode():
		return http.StatusNotFound
	case retcode >= 20000000:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
