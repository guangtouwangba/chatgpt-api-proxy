package api

import (
	"chatgpt-api-proxy/config"
	"chatgpt-api-proxy/internal/constant"
	"chatgpt-api-proxy/pkg/httphelper"
	"chatgpt-api-proxy/pkg/logger"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
)

func InitChatRouter(r *gin.Engine) *gin.Engine {
	api := r.Group("/api/openai")
	api.POST("/chat", HandleChat)
	return r
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`

	// This property isn't in the official documentation, but it's in
	// the documentation for the official library for python:
	// - https://github.com/openai/openai-python/blob/main/chatml.md
	// - https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	Name string `json:"name,omitempty"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	MaxTokens        int                     `json:"max_tokens,omitempty"`
	Temperature      float32                 `json:"temperature,omitempty"`
	TopP             float32                 `json:"top_p,omitempty"`
	N                int                     `json:"n,omitempty"`
	Stream           bool                    `json:"stream,omitempty"`
	Stop             []string                `json:"stop,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int          `json:"logit_bias,omitempty"`
	User             string                  `json:"user,omitempty"`
}

type ChatCompletionChoice struct {
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
	FinishReason string                `json:"finish_reason"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   Usage                  `json:"usage"`
}

func HandleChat(c *gin.Context) {
	request := ChatCompletionRequest{}
	if request.Model == "" {
		request.Model = openai.GPT3Dot5Turbo
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		httphelper.WrapperError(c, constant.NewBaseErrorWithMsg(constant.InvalidRequestError, err.Error()))
		return
	}

	client := openai.NewClient(config.Store.OpenAI.APIKey)

	if request.Stream {
		streamChat(c, &request, client)
		return
	}

	chat(c, &request, client)
}

func convertMassages(msg []ChatCompletionMessage) []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, 0, len(msg))
	for _, m := range msg {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
			Name:    m.Name,
		})
	}
	return messages
}

func chat(c *gin.Context, request *ChatCompletionRequest, client *openai.Client) {
	response, err := client.CreateChatCompletion(c, openai.ChatCompletionRequest{
		Model:            request.Model,
		Messages:         convertMassages(request.Messages),
		MaxTokens:        request.MaxTokens,
		Temperature:      request.Temperature,
		TopP:             request.TopP,
		N:                request.N,
		Stream:           request.Stream,
		Stop:             request.Stop,
		PresencePenalty:  request.PresencePenalty,
		FrequencyPenalty: request.FrequencyPenalty,
		LogitBias:        request.LogitBias,
		User:             request.User,
	})

	if err != nil {
		httphelper.WrapperError(c, constant.NewBaseErrorWithMsg(constant.InvalidRequestError, err.Error()))
		return
	}

	httphelper.WrapperSuccess(c, response)
}

func streamChat(c *gin.Context, request *ChatCompletionRequest, client *openai.Client) {
	response, err := client.CreateChatCompletionStream(c, openai.ChatCompletionRequest{
		Model:            request.Model,
		Messages:         convertMassages(request.Messages),
		MaxTokens:        request.MaxTokens,
		Temperature:      request.Temperature,
		TopP:             request.TopP,
		N:                request.N,
		Stream:           request.Stream,
		Stop:             request.Stop,
		PresencePenalty:  request.PresencePenalty,
		FrequencyPenalty: request.FrequencyPenalty,
		LogitBias:        request.LogitBias,
		User:             request.User,
	})

	if err != nil {
		httphelper.WrapperError(c, constant.NewBaseErrorWithMsg(constant.InvalidRequestError, err.Error()))
		return
	}

	for {
		resp, err := response.Recv()
		if err != nil && !errors.Is(err, io.EOF) {
			httphelper.WrapperError(c, constant.NewBaseErrorWithMsg(constant.InternalError, err.Error()))
			return
		}
		// stream chat response
		if errors.Is(err, io.EOF) {
			return
		}
		_, err = c.Writer.Write([]byte(resp.Choices[0].Delta.Content))
		if err != nil && !errors.Is(err, io.EOF) {
			logger.Infof("stream chat error: %v", err)
			httphelper.WrapperError(c, constant.NewBaseErrorWithMsg(constant.InternalServerError, err.Error()))
			return
		}
	}
}
