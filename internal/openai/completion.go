package openai

import (
	"bytes"
	"chatgpt-api-proxy/internal/constant"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const completionBaseURL = "https://api.openai.com/v1/completions"

type CompletionModel string

const (
	Davinci CompletionModel = "text-davinci-003"
	Babbage CompletionModel = "babbage"
)

var completionModels = []CompletionModel{
	Davinci,
	Babbage,
}

type CompletionRequest struct {
	Model       CompletionModel `json:"model"`
	MaxTokens   int             `json:"max_tokens"`
	Prompt      string          `json:"prompt"`
	Temperature float64         `json:"temperature"`
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

// Choice is a choice from the completion endpoint.
type Choice struct {
	Text         string      `json:"text"`
	Index        int         `json:"index"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

// Usage is a usage from the completion endpoint.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type StreamCompletionResponse struct {
}

// Completion returns a completion response.
func (c *Client) Completion(ctx context.Context, request *CompletionRequest) (*CompletionResponse, error) {
	if request.Model == "" {
		request.Model = Davinci
	}

	if !isSupportedModel(request.Model) {
		return nil, errors.Wrap(constant.Error(constant.InvalidRequestError), "unsupported model")
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, errors.Wrap(constant.Error(constant.JSONMarshalError), err.Error())
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, completionBaseURL, bytes.NewReader(payload))
	if err != nil {
		return nil, errors.Wrap(constant.Error(constant.HTTPRequestError), err.Error())
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+c.APIKey)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, errors.Wrap(constant.Error(constant.HTTPRequestError), err.Error())
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusUnauthorized:
			return nil, errors.Wrap(constant.Error(constant.AuthenticationError), err.Error())
		case http.StatusTooManyRequests:
			return nil, errors.Wrap(constant.Error(constant.TooManyRequestsError), err.Error())
		case http.StatusServiceUnavailable:
			return nil, errors.Wrap(constant.Error(constant.ServiceUnavailableError), err.Error())
		case http.StatusGatewayTimeout:
			return nil, errors.Wrap(constant.Error(constant.GatewayTimeoutError), err.Error())
		default:
			return nil, errors.Wrap(constant.Error(constant.HTTPRequestError), err.Error())
		}
	}

	var completionResponse CompletionResponse
	err = json.NewDecoder(resp.Body).Decode(&completionResponse)
	if err != nil {
		return nil, errors.Wrap(constant.Error(constant.JSONUnmarshalError), err.Error())
	}

	return &completionResponse, nil
}

func isSupportedModel(model CompletionModel) bool {
	for _, m := range completionModels {
		if model == m {
			return true
		}
	}

	return false
}
