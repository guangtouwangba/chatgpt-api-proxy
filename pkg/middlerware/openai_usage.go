package middlerware

import (
	"bytes"
	"chatgpt-api-proxy/internal/db"
	"chatgpt-api-proxy/internal/db/model"
	"chatgpt-api-proxy/internal/db/repository"
	"chatgpt-api-proxy/pkg/logger"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func OpenAIUsage() gin.HandlerFunc {
	return func(c *gin.Context) {
		bodyLogWriter := &bodyLogWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter
		// parse response body
		// convert to json
		c.Next()
		// if response is not ok, return
		if c.Writer.Status() != http.StatusOK {
			return
		}
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
		mod, ok := openAIResponse["data"].(map[string]interface{})["model"].(string)
		if !ok {
			logger.Errorf("error when get openai model")
		}
		database := db.GetDB()
		repo := repository.NewGormOpenAIUsageRepository(database)
		openaiUsage := &model.OpenAIUsage{
			OpenAIID: id,
			Usage:    int64(usage["total_tokens"].(float64)),
			Tokens:   int64(usage["total_tokens"].(float64)),
			Model:    mod,
		}
		err = repo.CreateOrUpdate(openaiUsage)
		if err != nil {
			logger.Errorf("error when create or update openai usage: %v", err)
		}
	}
}
