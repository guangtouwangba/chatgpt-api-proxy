package api

import (
	"bytes"
	"chatgpt-api-proxy/config"
	"chatgpt-api-proxy/internal/constant"
	"chatgpt-api-proxy/pkg/httphelper"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

var completionModels = []string{
	"text-davinci-003",
	"babage",
}

type CompletionModel string

const (
	Davinci CompletionModel = "text-davinci-003"
	Babbage CompletionModel = "babbage"
)

const completionBaseURL = "https://api.openai.com/v1/completions"

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
	Stop             string         `json:"stop"`
}

// CompletionResponse is a response from the completion endpoint.
type CompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Usage is a usage from the completion endpoint.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Choice is a choice from the completion endpoint.
type Choice struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

func HandleCompletion(c *gin.Context) {
	// parse request
	request := CompletionRequest{}

	if err := c.Bind(&request); err != nil {
		httphelper.WrapperError(c, constant.NewBaseError(constant.InvalidRequestError, err.Error()))
		return
	}

	// validate request
	if !isSupportedModel(request.Model) {
		httphelper.WrapperError(c, constant.NewBaseError(constant.InvalidRequestError, "no supported model"))
		return
	}

	// call openai client
	client := NewClient()
	client.SendCompletionRequest(c, request)
}

func isSupportedModel(model string) bool {
	for _, m := range completionModels {
		if m == model {
			return true
		}
	}
	return false
}

// Client is a client for OpenAI API.
type Client struct {
	httpClient *http.Client
	APIKey     string
}

// NewClient returns a new OpenAI client.
func NewClient() *Client {
	apiKey := config.Store.OpenAI.APIKey
	httpClient := &http.Client{}

	return &Client{
		APIKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *Client) SendCompletionRequest(ctx *gin.Context, request CompletionRequest) {
	if request.Model == "" {
		request.Model = string(Davinci)
	}

	if !isSupportedModel(request.Model) {
		httphelper.WrapperError(ctx, constant.NewBaseError(constant.InvalidRequestError, "not supported model"))
		return
	}

	payload, err := json.Marshal(request)
	if err != nil {
		httphelper.WrapperError(ctx, constant.NewBaseError(constant.InternalError, err.Error()))
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, completionBaseURL, bytes.NewReader(payload))
	if err != nil {
		httphelper.WrapperError(ctx, constant.NewBaseError(constant.InternalError, err.Error()))
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.APIKey)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		httphelper.WrapperError(ctx, constant.NewBaseError(constant.InternalError, err.Error()))
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			httphelper.WrapperError(ctx, constant.NewBaseError(constant.AuthenticationError, err.Error()))
		case http.StatusTooManyRequests:
			httphelper.WrapperError(ctx, constant.NewBaseError(constant.TooManyRequestsError, err.Error()))
		case http.StatusServiceUnavailable:
			httphelper.WrapperError(ctx, constant.NewBaseError(constant.ServiceUnavailableError, err.Error()))
		case http.StatusGatewayTimeout:
			httphelper.WrapperError(ctx, constant.NewBaseError(constant.GatewayTimeoutError, err.Error()))
		default:
			httphelper.WrapperError(ctx, constant.NewBaseError(constant.InternalError, err.Error()))
		}
	}

	var completionResponse CompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&completionResponse)
	if err != nil {
		httphelper.WrapperError(ctx, constant.NewBaseError(constant.JSONUnmarshalError, err.Error()))
	}

	httphelper.WrapperSuccess(ctx, resp)
}
