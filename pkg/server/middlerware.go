package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Recover 中间件
func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "服务器内部错误",
			})
			c.Abort()
			return
		}
	}()
	c.Next()
}
