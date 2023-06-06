package api

import (
	"chatgpt-api-proxy/internal/constant"
	"chatgpt-api-proxy/internal/openai"
	"chatgpt-api-proxy/pkg/httphelper"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var completionModels = []string{
	"text-davinci-003",
	"babage",
}

func InitCompletionRouter(r *gin.Engine) {
	api := r.Group("/api/openai")
	api.POST("/completion", HandleCompletion)
}

type CompletionRequest struct {
	Model            string         `json:"model" binding:"required"`
	MaxTokens        int            `json:"max_tokens"`
	Prompt           string         `json:"prompt"`
	Temperature      float64        `json:"temperature"`
	TopP             float64        `json:"top_p"`
	N                int            `json:"n"`
	Stream           bool           `json:"stream"`
	LogProbs         int            `json:"logprobs"`
	Echo             bool           `json:"echo"`
	PresencePenalty  float64        `json:"presence_penalty"`
	FrequencyPenalty float64        `json:"frequency_penalty"`
	BestOf           int            `json:"best_of"`
	LogitBias        map[string]int `json:"logit_bias"`
	UserID           string         `json:"user"`
}

func HandleCompletion(c *gin.Context) {
	// parse request
	request := CompletionRequest{}
	err := c.Bind(&request)
	if err != nil {
		httphelper.WrapperError(c, constant.NewBaseError(constant.InvalidRequestError, err.Error()))
		return
	}

	// validate request
	if !isSupportedModel(request.Model) {
		httphelper.WrapperError(c, constant.NewBaseError(constant.InvalidRequestError, "no supported model"))
		return
	}

	// call openai client
	client := openai.NewClient()
	response, err := client.Completion(c, &openai.CompletionRequest{
		Model:       openai.CompletionModel(request.Model),
		MaxTokens:   request.MaxTokens,
		Prompt:      request.Prompt,
		Temperature: request.Temperature,
	})
	if err != nil {
		if errors.Is(err, constant.Error(constant.GatewayTimeoutError)) {
			httphelper.WrapperError(c, constant.NewBaseError(constant.GatewayTimeoutError, err.Error()))
			return
		}
		httphelper.WrapperError(c, constant.NewBaseError(constant.InternalError, err.Error()))
		return
	}

	httphelper.WrapperSuccess(c, response)
}

func isSupportedModel(model string) bool {
	for _, m := range completionModels {
		if m == model {
			return true
		}
	}
	return false
}
