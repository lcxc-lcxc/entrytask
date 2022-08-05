package web

import (
	"entrytask/internel/middleware"
	"entrytask/internel/web/api"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.Default()

	ping := api.NewPing()
	user := api.NewUser()
	product := api.NewProduct()
	comment := api.NewComment()

	router.GET("/ping", ping.Ping)
	router.POST("/api/sessions", user.Login)
	router.POST("/api/users", user.Register)

	productGroup := router.Group("/api/products")
	productGroup.Use(middleware.AuthSessionID)
	{
		productGroup.GET("", product.List)
		productGroup.GET("/search", product.Search)
		productGroup.GET("/:product_id", product.Detail)

		productGroup.GET("/:product_id/comments/:comment_id", comment.Detail)
	}

	commentGroup := router.Group("/api/products/:product_id/comments")
	commentGroup.Use(middleware.AuthUserLogin)
	{
		commentGroup.POST("", comment.Create)
		commentGroup.POST("/:comment_id/reply", comment.Reply)
	}

	return router
}
