package router

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	api.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}
