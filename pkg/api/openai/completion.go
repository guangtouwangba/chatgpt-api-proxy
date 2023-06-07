package api

import (
	"bytes"
	"chatgpt-api-proxy/config"
	"chatgpt-api-proxy/internal/constant"
	"chatgpt-api-proxy/pkg/httphelper"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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

const (
	completionBaseURL = "https://api.openai.com/v1/completions"
	defaultBufferSize = 1024
)

func InitCompletionRouter(r *gin.Engine) {
	api := r.Group("/api/openai")
	api.POST("/completion", HandleCompletion)
}

type CompletionRequest struct {
	Model            string         `json:"model"`
	Prompt           any            `json:"prompt,omitempty"`
	Suffix           string         `json:"suffix,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	Temperature      float32        `json:"temperature,omitempty"`
	TopP             float32        `json:"top_p,omitempty"`
	N                int            `json:"n,omitempty"`
	Stream           bool           `json:"stream,omitempty"`
	LogProbs         int            `json:"logprobs,omitempty"`
	Echo             bool           `json:"echo,omitempty"`
	Stop             []string       `json:"stop,omitempty"`
	PresencePenalty  float32        `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32        `json:"frequency_penalty,omitempty"`
	BestOf           int            `json:"best_of,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	User             string         `json:"user,omitempty"`
}

// CompletionChoice represents one of possible completions.
type CompletionChoice struct {
	Text         string        `json:"text"`
	Index        int           `json:"index"`
	FinishReason string        `json:"finish_reason"`
	LogProbs     LogprobResult `json:"logprobs"`
}

// LogprobResult represents logprob result of Choice.
type LogprobResult struct {
	Tokens        []string             `json:"tokens"`
	TokenLogprobs []float32            `json:"token_logprobs"`
	TopLogprobs   []map[string]float32 `json:"top_logprobs"`
	TextOffset    []int                `json:"text_offset"`
}

// CompletionResponse represents a response structure for completion API.
type CompletionResponse struct {
	ID      string             `json:"id"`
	Object  string             `json:"object"`
	Created int64              `json:"created"`
	Model   string             `json:"model"`
	Choices []CompletionChoice `json:"choices"`
	Usage   Usage              `json:"usage"`
}

// Usage Represents the total token usage per request to OpenAI.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func HandleCompletion(c *gin.Context) {
	// parse request
	request := CompletionRequest{}

	if err := c.Bind(&request); err != nil {
		httphelper.WrapperError(c, constant.NewBaseErrorWithMsg(constant.InvalidRequestError, err.Error()))
		return
	}

	// validate request
	if !isSupportedModel(request.Model) {
		httphelper.WrapperError(c, constant.NewBaseErrorWithMsg(constant.InvalidRequestError, "no supported model"))
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
	validateModel(ctx, request)

	payload, err := json.Marshal(request)
	if err != nil {
		httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.InternalError, err.Error()))
		return
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, completionBaseURL, bytes.NewReader(payload))
	if err != nil {
		httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.InternalError, err.Error()))
	}

	c.initRequestHeader(ctx, req)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.InternalError, err.Error()))
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	logrus.Infof("got response code from openai: %+v", resp.StatusCode)
	c.checkResponse(ctx, resp)

	var completionResponse CompletionResponse

	err = json.NewDecoder(resp.Body).Decode(&completionResponse)
	if err != nil {
		// parse response body
		var body []byte
		if resp.Body != nil {
			body, _ = io.ReadAll(resp.Body)
		}
		logrus.Infof("got response body from openai: %+v", string(body))
		httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.JSONUnmarshalError, err.Error()))
		return
	}
	ctx.JSON(http.StatusOK, completionResponse)
}

func (c *Client) checkResponse(ctx *gin.Context, resp *http.Response) {
	if resp.StatusCode != http.StatusOK {
		logrus.Info("resp.StatusCode: ", resp.StatusCode)
		// parse response body
		var body []byte
		if resp.Body != nil {
			body, _ = io.ReadAll(resp.Body)
		}
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.AuthenticationError, string(body)))
			return
		case http.StatusTooManyRequests:
			httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.TooManyRequestsError, string(body)))
			return
		case http.StatusServiceUnavailable:
			httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.ServiceUnavailableError, string(body)))
			return
		case http.StatusGatewayTimeout:
			httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.GatewayTimeoutError, string(body)))
			return
		default:
			httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.InternalError, string(body)))
			return
		}
	}
}

func validateModel(ctx *gin.Context, request CompletionRequest) {
	if request.Model == "" {
		request.Model = string(Davinci)
	}

	if !isSupportedModel(request.Model) {
		httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.InvalidRequestError, "not supported model"))
		return
	}
}

func (c *Client) SendStreamCompletion(ctx *gin.Context, request *CompletionRequest) {
	request.Stream = true
	validateModel(ctx, *request)

	payload, err := json.Marshal(request)
	if err != nil {
		httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.InternalError, err.Error()))
		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, completionBaseURL, bytes.NewReader(payload))
	if err != nil {
		httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.InternalError, err.Error()))
		return
	}

	c.initRequestHeader(ctx, req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		httphelper.WrapperError(ctx, constant.NewBaseErrorWithMsg(constant.InternalError, err.Error()))
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	c.checkResponse(ctx, resp)

	buf := make([]byte, defaultBufferSize)
	for {
		if n, err := resp.Body.Read(buf); errors.Is(err, io.EOF) || n == 0 {
			return
		} else if err != nil {
			log.Println("error while reading respbody: ", err.Error())
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		} else {
			if _, err = ctx.Writer.Write(buf[:n]); err != nil {
				log.Println("error while writing resp: ", err.Error())
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			if f, ok := ctx.Writer.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

func (c *Client) initRequestHeader(ctx *gin.Context, req *http.Request) {
	// copy headers from gin request
	for k, v := range ctx.Request.Header {
		req.Header.Set(k, v[0])
	}

	// if authorization header is not set, set it
	if req.Header.Get("Authorization") == "" {
		req.Header.Add("Authorization", "Bearer "+c.APIKey)
	}
	req.Header.Add("Content-Type", "application/json")
}
