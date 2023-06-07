package router

import (
	api "chatgpt-api-proxy/pkg/api/openai"
	"chatgpt-api-proxy/pkg/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(server.Recover)
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
