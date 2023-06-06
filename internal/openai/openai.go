package openai

import (
	"chatgpt-api-proxy/config"
	"net/http"
)

// OpenAI is an interface for OpenAI API.
type OpenAI interface {
	Completion(request CompletionRequest) (CompletionResponse, error)
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
