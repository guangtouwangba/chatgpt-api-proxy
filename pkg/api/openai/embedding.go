package api

import (
	"chatgpt-api-proxy/pkg/httphelper"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

func InitEmbeddingRouter(r *gin.Engine) {
	api := r.Group("/api/openai")
	api.POST("/embedding", HandleEmbedding)
}

var models = map[string]openai.EmbeddingModel{
	"text-similarity-ada-001":       openai.AdaSimilarity,
	"text-similarity-babbage-001":   openai.BabbageSimilarity,
	"text-similarity-curie-001":     openai.CurieSimilarity,
	"text-similarity-davinci-001":   openai.DavinciSimilarity,
	"text-search-ada-doc-001":       openai.AdaSearchDocument,
	"text-search-ada-query-001":     openai.AdaSearchQuery,
	"text-search-babbage-doc-001":   openai.BabbageSearchDocument,
	"text-search-babbage-query-001": openai.BabbageSearchQuery,
	"text-search-curie-doc-001":     openai.CurieSearchDocument,
	"text-search-curie-query-001":   openai.CurieSearchQuery,
	"text-search-davinci-doc-001":   openai.DavinciSearchDocument,
	"text-search-davinci-query-001": openai.DavinciSearchQuery,
	"code-search-ada-code-001":      openai.AdaCodeSearchCode,
	"code-search-ada-text-001":      openai.AdaCodeSearchText,
	"code-search-babbage-code-001":  openai.BabbageCodeSearchCode,
	"code-search-babbage-text-001":  openai.BabbageCodeSearchText,
	"text-embedding-ada-002":        openai.AdaEmbeddingV2,
}

// EmbeddingResponse is the response from a Create embeddings request.
type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

// Embedding is an individual embedding vector.
type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
}

// EmbeddingRequest is the input to a Create embeddings request.
type EmbeddingRequest struct {
	// Input is a slice of strings for which you want to generate an Embedding vector.
	// Each input must not exceed 2048 tokens in length.
	// OpenAPI suggests replacing newlines (\n) in your input with a single space, as they
	// have observed inferior results when newlines are present.
	// E.g.
	//	"The food was delicious and the waiter..."
	Input []string `json:"input"`
	// ID of the model to use. You can use the List models API to see all of your available models,
	// or see our Model overview for descriptions of them.
	Model string `json:"model"`
	// A unique identifier representing your end-user, which will help OpenAI to monitor and detect abuse.
	User string `json:"user"`
}

func HandleEmbedding(c *gin.Context) {
	request := EmbeddingRequest{}
	if request.Model == "" {
		request.Model = "ada"
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		httphelper.WrapperError(c, httphelper.ErrInternalServerError)
		return
	}

	client := openai.NewClient(getOpenAIAPIKey(c))
	response, err := client.CreateEmbeddings(c, openai.EmbeddingRequest{
		Input: request.Input,
		Model: models[request.Model],
	})

	if err != nil {
		httphelper.WrapperError(c, httphelper.ErrInternalServerError)
		return
	}

	httphelper.WrapperSuccess(c, convertEmbeddingResponse(response))
}

func convertEmbeddingResponse(response openai.EmbeddingResponse) EmbeddingResponse {
	data := make([]Embedding, len(response.Data))
	for i, embedding := range response.Data {
		data[i] = Embedding{
			Object:    embedding.Object,
			Embedding: embedding.Embedding,
			Index:     embedding.Index,
		}
	}
	return EmbeddingResponse{
		Object: response.Object,
		Data:   data,
		Model:  response.Model.String(),
	}
}
