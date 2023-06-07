package intergrationtest

import (
	"bytes"
	api "chatgpt-api-proxy/pkg/api/openai"
	intergrationtest "chatgpt-api-proxy/tests/intergrationtest/common"
	"encoding/json"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"testing"
)

func TestCompletion(t *testing.T) {
	config := intergrationtest.InitConfig()
	ctx := context.Background()
	testURL := "http://localhost:8080/api/openai/completion"

	request := &api.CompletionRequest{
		Model:  "text-davinci-003",
		Prompt: "This is a test",
	}

	client := http.DefaultClient

	reqBody, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, testURL, bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.APIKey)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusTooManyRequests {
		// parse error response
		var body []byte
		if resp.Body != nil {
			body, _ = io.ReadAll(resp.Body)
		}
		t.Fatalf("expected status code %d, got %d, body: %s", http.StatusOK, resp.StatusCode, string(body))
	}

	var completionResponse api.CompletionResponse

	err = json.NewDecoder(resp.Body).Decode(&completionResponse)
	if err != nil {
		t.Fatal(err)
	}
	if completionResponse.ID == "" {
		t.Fatalf("expected completion response ID, got empty string")
	}
}
