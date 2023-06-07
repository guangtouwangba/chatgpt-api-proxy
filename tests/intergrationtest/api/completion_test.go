//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	api "chatgpt-api-proxy/pkg/api/openai"
	integrationtest "chatgpt-api-proxy/tests/intergrationtest/common"
	"fmt"

	"encoding/json"
	"io"
	"net/http"
	"testing"

	"golang.org/x/net/context"
)

func TestCompletion(t *testing.T) {
	config := integrationtest.InitConfig()
	ctx := context.Background()
	testURL := fmt.Sprintf("%s/api/v1/openai/completion", config.Host)

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

	defer func() {
		_ = resp.Body.Close()
	}()

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
