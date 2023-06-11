package middlerware

import (
	"bytes"
	"chatgpt-api-proxy/pkg/logger"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func OpenAIUsage() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter
		// parse response body
		// convert to json
		c.Next()
		data := bodyLogWriter.body.Bytes()
		// TODO: temp use map[string]interface{} to unmarshal, since cycle import with api response
		var openAIResponse map[string]interface{}
		err := json.Unmarshal(data, &openAIResponse)
		if err != nil {
			logger.Errorf("error when unmarshal openai response: %v", err)
		}
		usage, ok := openAIResponse["data"].(map[string]interface{})["usage"].(map[string]interface{})
		if !ok {
			logger.Errorf("error when get openai usage")
		}
		id, ok := openAIResponse["data"].(map[string]interface{})["id"].(string)
		if !ok {
			logger.Errorf("error when get openai id")
		}
		logger.Infof("openai usage: %v", usage)
		logger.Infof("openai id: %v", id)
	}
}
