package response

import (
	"entrytask/internel/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func (r *Response) ResponseError(retcode int, data string) {

	r.Ctx.JSON(GetStatusCode(retcode), gin.H{
		"retcode": retcode,
		"msg":     constant.RetcodeMap[retcode].GetMsg(),
		"data":    data,
	})
}
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
