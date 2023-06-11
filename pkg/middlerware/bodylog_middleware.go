package middlerware

import (
	"bytes"
	"chatgpt-api-proxy/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	res, err := w.ResponseWriter.Write(b)
	return res, errors.Wrap(err, "error when write")
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	res, err := w.ResponseWriter.WriteString(s)
	return res, errors.Wrap(err, "error when write string")
}

func BodyLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter
		c.Next()
		logger.Infof("request body: %s", bodyLogWriter.body.String())
	}
}
