package middlerware

import (
	"chatgpt-api-proxy/internal/constant"
	"chatgpt-api-proxy/pkg/httphelper"

	"github.com/gin-gonic/gin"
)

// Recover is a middleware that recovers from any panics and writes a 500 if there was one.
func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				httphelper.WrapperError(c, constant.NewBaseError(constant.InternalServerError, err.Error()))
			} else {
				httphelper.WrapperError(c, constant.NewBaseError(constant.InternalServerError, "internal server error"))
			}
			c.Abort()
			return
		}
	}()
	c.Next()
}
