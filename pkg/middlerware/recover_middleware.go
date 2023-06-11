package middlerware

import (
	"chatgpt-api-proxy/pkg/httphelper"
	"chatgpt-api-proxy/pkg/logger"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recover is a middleware that recovers from any panics and writes a 500 if there was one.
func Recover(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				stack := string(debug.Stack())
				logger.Errorf("panic: %v\n%s", err, stack)
				httphelper.WrapperError(c, httphelper.NewAPIError(http.StatusInternalServerError, err.Error()))
			} else {
				httphelper.WrapperError(c, httphelper.ErrInternalServerError)
			}
			c.Abort()
			return
		}
	}()
	c.Next()
}
