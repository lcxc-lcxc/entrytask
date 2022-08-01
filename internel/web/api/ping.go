package api

import (
	"entrytask/internel/web/response"
	"github.com/gin-gonic/gin"
)

type Ping struct{}

func NewPing() Ping {
	return Ping{}
}

func (p Ping) Ping(c *gin.Context) {
	//c.JSON(200, gin.H{"message": "pong"})
	//
	resp := response.NewResponse(c)
	resp.ResponseOK("pong")
}
