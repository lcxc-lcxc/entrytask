package web

import (
	"entrytask/internel/web/api"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	ping := api.NewPing()

	router.GET("/ping", ping.Ping)

	return router
}
