package api

import (
	"chatgpt-api-proxy/config"
	"strings"

	"github.com/gin-gonic/gin"
)

func getOpenAIAPIKey(c *gin.Context) string {
	// check if header contains key
	key := c.GetHeader("Authorization")
	if key != "" {
		return strings.TrimPrefix(key, "Bearer ")
	}

	// return api key from config
	return config.Store.GetOpenAIApiKey()
}
