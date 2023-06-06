package router

import (
	api "chatgpt-api-proxy/pkg/api/openai"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	raw := r.Group("/")
	raw.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	apiGroup := r.Group("/api")
	apiGroup.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api.InitCompletionRouter(r)

	return r
}
