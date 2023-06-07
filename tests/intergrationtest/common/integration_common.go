package intergrationtest

import "os"

type IntegrationTestConfig struct {
	APIKey string
}

func InitConfig() *IntegrationTestConfig {
	apiKey := os.Getenv("OPENAI_API_KEY")

	return &IntegrationTestConfig{
		APIKey: apiKey,
	}
}
