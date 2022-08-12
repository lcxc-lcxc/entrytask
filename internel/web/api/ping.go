package api

import (
	"entrytask/internel/constant"
	"entrytask/internel/web/response"
	"github.com/gin-gonic/gin"
)

type Ping struct{}

func NewPing() Ping {
	return Ping{}
}

type PingRequest struct {
	Ping int `form:"ping" json:"ping" binding:"isdefault=0"`
}

func (p Ping) Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})

	resp := response.NewResponse(c)

	param := PingRequest{}
	err := c.ShouldBind(&param)
	if err != nil {
		resp.ResponseError(constant.ServerError.GetRetCode(), "")
		return
	}

}
