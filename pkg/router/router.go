package router

import (
	api "chatgpt-api-proxy/pkg/api/openai"
	"chatgpt-api-proxy/pkg/middlerware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middlerware.Recover)
	raw := r.Group("/")
	raw.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	apiGroup := r.Group("/api")
	apiGroup.GET("ping", middlerware.BodyLog(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api.InitCompletionRouter(r)
	api.InitChatRouter(r)

	return r
}
